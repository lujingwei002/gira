package hall

import (
	"context"
	"fmt"
	"io"
	"runtime/debug"
	"sync"
	"time"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/actor"
	"github.com/lujingwei002/gira/facade"
	"github.com/lujingwei002/gira/log"
	"github.com/lujingwei002/gira/service/hall/hall_grpc"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type hall_server struct {
	*actor.Actor
	framework     *Framework
	ctx           context.Context
	cancelFunc    context.CancelFunc
	SessionDict   sync.Map
	sessionCount  int64
	hallHandler   HallHandler
	playerHandler gira.ProtoHandler
	proto         gira.Proto
	config        *Config
	isDestory     bool
}

func newHall(framework *Framework, proto gira.Proto, config *Config, hallHandler HallHandler, playerHandler gira.ProtoHandler) *hall_server {
	return &hall_server{
		framework:     framework,
		proto:         proto,
		config:        config,
		hallHandler:   hallHandler,
		playerHandler: playerHandler,
		Actor:         actor.NewActor(16),
	}
}

func (hall *hall_server) onAwake() error {
	// 注册grpc
	facade.RegisterGrpc(func(server *grpc.Server) error {
		hall_grpc.RegisterHallServer(server, &grpc_hall_server{
			hall: hall,
		})
		return nil
	})
	go hall.serve()
	return nil
}

func (hall *hall_server) onDestory() {
	hall.isDestory = true
	for {
		log.Infow("hall停止中", "session_count", hall.sessionCount)
		sessions := make([]*hall_sesssion, 0)
		hall.SessionDict.Range(func(key, value any) bool {
			session, _ := value.(*hall_sesssion)
			sessions = append(sessions, session)
			return true
		})
		for _, session := range sessions {
			session.Call_close(hall.ctx, actor.WithCallTimeOut(5*time.Second))
		}
		log.Infow("hall停止中", "session_count", hall.sessionCount)
		if hall.sessionCount <= 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (hall *hall_server) serve() {
	//
	// 1.服务器关闭时保存数据后再退出
	// 2.处理actor调用
	//
	hall.ctx, hall.cancelFunc = context.WithCancel(facade.Context())
	defer hall.cancelFunc()
	// facade.Go(func() error {
	// 	select {
	// 	case <-hall.ctx.Done():
	// 		// 服务器停止，保存数据
	// 		for {
	// 			log.Infow("hall停止中", "session_count", hall.sessionCount)
	// 			if hall.sessionCount <= 0 {
	// 				break
	// 			}
	// 			time.Sleep(100 * time.Millisecond)
	// 		}
	// 	}
	// 	return nil
	// })
	for {
		select {
		// 处理actor请求
		case r := <-hall.Inbox():
			r.Call()
		case <-hall.ctx.Done():
			log.Infow("hall exit")
			return
		}
	}
}

func (hall *hall_server) SessionCount() int64 {
	return hall.sessionCount
}

// 推送消息给其他玩家
func (hall *hall_server) Push(ctx context.Context, userId string, push gira.ProtoPush) (err error) {
	defer func() {
		log.Warnw("hall_server push panic", "user_id", userId, "route", push.GetPushName())
		if e := recover(); e != nil {
			log.Error(e)
			debug.PrintStack()
			err = e.(error)
		}
	}()

	if v, ok := hall.SessionDict.Load(userId); !ok {
		var data []byte
		var peer *gira.Peer
		data, err = hall.proto.PushEncode(push)
		if err != nil {
			return
		}
		peer, err = facade.WhereIsUser(userId)
		if err != nil {
			return
		}
		if _, err = hall_grpc.DefaultHallClients.WithUnicast().WithPeer(peer).MustPush(ctx, &hall_grpc.MustPushRequest{
			UserId: userId,
			Data:   data,
		}); err != nil {
			return
		}
		return
	} else {
		session, _ := v.(*hall_sesssion)
		// WARN: chPush有可能已经关闭
		session.chPeerPush <- push
		return
	}
}

// 推送消息给其他玩家
func (hall *hall_server) MustPush(ctx context.Context, userId string, push gira.ProtoPush) (err error) {
	defer func() {
		log.Warnw("hall_server must push panic", "user_id", userId, "route", push.GetPushName())
		if e := recover(); e != nil {
			log.Error(e)
			debug.PrintStack()
			err = e.(error)
		}
	}()
	if v, ok := hall.SessionDict.Load(userId); !ok {
		var data []byte
		var peer *gira.Peer
		data, err = hall.proto.PushEncode(push)
		if err != nil {
			return
		}
		peer, err = facade.WhereIsUser(userId)
		if err != nil {
			return
		}
		if _, err = hall_grpc.DefaultHallClients.WithUnicast().WithPeer(peer).MustPush(ctx, &hall_grpc.MustPushRequest{
			UserId: userId,
			Data:   data,
		}); err != nil {
			return
		}
		return
	} else {
		session, _ := v.(*hall_sesssion)
		// WARN: chPush有可能已经关闭
		session.chPeerPush <- push
		return
	}
}

type grpc_hall_server struct {
	hall_grpc.UnsafeHallServer
	hall *hall_server
}

// 收到顶号下线的请求
func (self grpc_hall_server) UserInstead(ctx context.Context, req *hall_grpc.UserInsteadRequest) (*hall_grpc.UserInsteadResponse, error) {
	resp := &hall_grpc.UserInsteadResponse{}
	userId := req.UserId
	log.Infow("user instead", "user_id", userId)
	if v, ok := self.hall.SessionDict.Load(userId); !ok {
		// 偿试解锁
		peer, err := facade.UnlockLocalUser(userId)
		log.Infow("unlock local user return", "user_id", userId, "peer", peer, "err", err)
		if err != nil {
			resp.ErrorCode = gira.ErrCode(err)
			resp.ErrorMsg = gira.ErrMsg(err)
			return resp, nil
		} else {
			return resp, nil
		}
	} else {
		session := v.(*hall_sesssion)
		timeoutCtx, timeoutFunc := context.WithTimeout(ctx, 10*time.Second)
		defer timeoutFunc()
		if err := session.Call_instead(timeoutCtx, fmt.Sprintf("在%s登录", req.Address), actor.WithCallTimeOut(1*time.Second)); err != nil {
			log.Errorw("user instead fail", "user_id", userId, "error", err)
			resp.ErrorCode = gira.ErrCode(err)
			resp.ErrorMsg = gira.ErrMsg(err)
			return resp, nil
		} else {
			log.Infow("user instead success", "user_id", userId)
			return resp, nil
		}
	}
}

// push消息流
func (self grpc_hall_server) PushStream(server hall_grpc.Hall_PushStreamServer) error {
	for {
		req, err := server.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		userId := req.UserId
		_, _, push, err := self.hall.proto.PushDecode(req.Data)
		if err != nil {
			log.Warnw("user push decode fail", "error", err)
			continue
		}
		if v, ok := self.hall.SessionDict.Load(userId); ok {
			session, _ := v.(*hall_sesssion)
			// WARN: chPush有可能已经关闭
			session.chPeerPush <- push
		}
	}
}

func (self grpc_hall_server) MustPush(ctx context.Context, req *hall_grpc.MustPushRequest) (resp *hall_grpc.MustPushResponse, err error) {
	resp = &hall_grpc.MustPushResponse{}
	var push gira.ProtoPush
	userId := req.UserId
	_, _, push, err = self.hall.proto.PushDecode(req.Data)
	if err != nil {
		return
	}
	defer func() {
		log.Warnw("grpc_hall_server must push panic", "user_id", req.UserId, "route", push.GetPushName())
		if e := recover(); e != nil {
			log.Error(e)
			debug.PrintStack()
			err = e.(error)
		}
	}()
	if v, ok := self.hall.SessionDict.Load(userId); !ok {
		err = gira.ErrNoSession
		return
	} else {
		session, _ := v.(*hall_sesssion)
		// WARN: chPush有可能已经关闭
		session.chPeerPush <- push
		return
	}
}

// func (self hall_server) UserPush(ctx context.Context, req *hall_grpc.UserPushRequest) (*hall_grpc.UserPushResponse, error) {
// 	resp := &hall_grpc.UserPushResponse{}
// 	userId := req.UserId
// 	_, _, _, push, err := self.hall.proto.PushDecode(req.Data)
// 	if err != nil {
// 		return resp, err
// 	}
// 	log.Println("ffffffffffffffffff1")
// 	if v, ok := self.hall.SessionDict.Load(userId); !ok {
// 		return resp, gira.ErrNoSession
// 	} else {
// 		session, _ := v.(*hall_sesssion)
// 		session.chPush <- push
// 		return resp, nil
// 	}
// }

// 处理网关发过来的消息
// type grpc_upstream_server struct {
// 	hall_grpc.UnsafeUpstreamServer
// 	hall *hall_server
// }

func (self grpc_hall_server) Info(ctx context.Context, req *hall_grpc.InfoRequest) (*hall_grpc.InfoResponse, error) {
	resp := &hall_grpc.InfoResponse{
		BuildTime:    facade.GetBuildTime(),
		BuildVersion: facade.GetBuildVersion(),
	}
	return resp, nil
}

func (self grpc_hall_server) Heartbeat(ctx context.Context, req *hall_grpc.HeartbeatRequest) (*hall_grpc.HeartbeatResponse, error) {
	resp := &hall_grpc.HeartbeatResponse{
		PlayerCount: self.hall.sessionCount,
	}
	return resp, nil
}

func (self grpc_hall_server) GateStream(client hall_grpc.Hall_GateStreamServer) error {
	// ticker := time.NewTicker(time.Duration(self.hall.config.Framework.Hall.GatewayReportInterval) * time.Second)
	defer func() {
		// ticker.Stop()
		log.Infow("gate stream exit")
	}()
	if self.hall.isDestory {
		return gira.ErrAlreadyDestory
	}
	errGroup, errCtx := errgroup.WithContext(self.hall.ctx)
	// 定时发送在线人数
	errGroup.Go(func() error {
		for {
			select {
			case <-self.hall.ctx.Done():
				return self.hall.ctx.Err()
			case <-errCtx.Done():
				return errCtx.Err()
			}
		}
	})
	// 循环接收gate的命令
	errGroup.Go(func() error {
		for {
			select {
			case <-self.hall.ctx.Done():
				return self.hall.ctx.Err()
			case <-errCtx.Done():
				return errCtx.Err()
			default:
				req, err := client.Recv()
				if err != nil {
					log.Infow("gate recv fail", "error", err)
					// 连接关闭，返回错误，会触发errCtx的cancel
					return err
				} else {
					log.Infow("gate recv", "req", req)
				}
			}
		}
	})
	return errGroup.Wait()
}

func (self grpc_hall_server) ClientStream(client hall_grpc.Hall_ClientStreamServer) error {
	var req *hall_grpc.ClientMessageNotify
	var err error
	var memberId string
	if self.hall.isDestory {
		return gira.ErrAlreadyDestory
	}
	if req, err = client.Recv(); err != nil {
		log.Errorw("recv fail", "error", err)
		return err
	}
	hall := self.hall
	sessionId := req.SessionId
	// 会话的第一条消息一定要带memberid来完成登录
	memberId = req.MemberId
	log.Infow("bind memberid", "session_id", sessionId, "member_id", memberId)
	if memberId == "" {
		log.Errorw("invalid memberid", "session_id", sessionId, "member_id", memberId)
		return gira.ErrInvalidMemberId
	}
	// 申请建立会话
	var session *hall_sesssion
	session, err = newSession(self.hall, sessionId, memberId)
	if err != nil {
		log.Errorw("create session fail", "session_id", sessionId, "error", err)
		return err
	}
	go session.serve()

	// 建立成功， 函数结束后通知释放
	// response管理
	session.chClientResponse = make(chan *hall_grpc.ClientMessagePush, hall.config.Framework.Hall.ResponseBufferSize)
	session.chClientRequest = make(chan *hall_grpc.ClientMessageNotify, hall.config.Framework.Hall.RequestBufferSize)

	// 绑定到hall
	clientCtx, clientCancelFunc := context.WithCancel(self.hall.ctx)
	session.clientCancelFunc = clientCancelFunc
	// recv ctrl context
	defer func() {
		clientCancelFunc()
		if err := session.Call_close(self.hall.ctx, actor.WithCallTimeOut(5*time.Second)); err != nil {
			log.Errorw("session close fail", "session_id", sessionId, "error", err)
		}
	}()
	// 读取client消息，塞进管道
	chRecvErr := make(chan struct{})
	go func() {
		defer func() {
			defer close(session.chClientRequest)
			log.Infow("read goroutine exit", "session_id", sessionId)
		}()
		session.chClientRequest <- req
		for {
			req, err = client.Recv()
			if clientCtx.Err() != nil {
				return
			}
			if err != nil {
				log.Infow("client recv fail", "error", err, "session_id", sessionId)
				close(chRecvErr)
				return
			} else {
				session.chClientRequest <- req
			}
		}
	}()
	// 将response消息转发给网关
	func() {
		defer func() {
			select {
			case <-session.chClientResponse:
			default:
				close(session.chClientResponse)
			}
			log.Infow("response goroutine exit", "session_id", sessionId)
		}()
		for {
			select {
			case resp := <-session.chClientResponse:
				if resp == nil {
					log.Infow("response chan close", "session_id", sessionId)
					return
				}
				if err := client.Send(resp); err != nil {
					log.Warnw("client response fail", "session_id", sessionId, "name", resp.Route, "error", err)
				} else {
					log.Infow("response", "session_id", sessionId, "name", resp.Route, "req_id", resp.ReqId, "data", len(resp.Data))
				}
			case <-clientCtx.Done():
				// 发完剩下的消息
				close(session.chClientResponse)
				goto FLUSH_RESPONSE
			case <-chRecvErr:
				// 已经断开了，马上结束
				return
			}
		}
	FLUSH_RESPONSE:
		for {
			select {
			case resp := <-session.chClientResponse:
				if resp == nil {
					log.Infow("response chan close", "session_id", sessionId)
					return
				}
				if err := client.Send(resp); err != nil {
					log.Warnw("client response fail", "session_id", sessionId, "name", resp.Route, "error", err)
				} else {
					log.Infow("response", "session_id", sessionId, "name", resp.Route, "req_id", resp.ReqId, "data", len(resp.Data))
				}
			case <-chRecvErr:
				// 已经断开了，马上结束
				return
			}
		}
	}()
	log.Infow("client stream exit", "error", err, "session_id", sessionId, "member_id", memberId)
	return err
}
