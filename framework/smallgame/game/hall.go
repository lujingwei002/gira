package game

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
	"github.com/lujingwei002/gira/framework/smallgame/common/rpc"
	"github.com/lujingwei002/gira/framework/smallgame/gen/grpc/hall_grpc"
	"github.com/lujingwei002/gira/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type hall_server struct {
	*actor.Actor
	framework     *Framework
	ctx           context.Context
	cancelFunc    context.CancelFunc
	SessionDict   sync.Map
	SessionCount  int64
	hallHandler   HallHandler
	playerHandler gira.ProtoHandler
	proto         gira.Proto
	config        *Config
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

func (hall *hall_server) OnAwake() error {
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

func (hall *hall_server) serve() {
	//
	// 1.服务器关闭时保存数据后再退出
	// 2.处理actor调用
	//
	hall.ctx, hall.cancelFunc = context.WithCancel(facade.Context())
	defer hall.cancelFunc()
	facade.Go(func() error {
		select {
		case <-hall.ctx.Done():
			// 服务器停止，保存数据
			for {
				log.Infow("hall停止中", "session_count", hall.SessionCount)
				if hall.SessionCount <= 0 {
					break
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
		return nil
	})
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
		if err = rpc.Hall.Push(ctx, peer, userId, data); err != nil {
			return
		}
		return
	} else {
		session, _ := v.(*hall_sesssion)
		// WARN: chPush有可能已经关闭
		session.chPush <- push
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
		if _, err = rpc.Hall.MustPush(ctx, peer, userId, data); err != nil {
			return
		}
		return
	} else {
		session, _ := v.(*hall_sesssion)
		session.chPush <- push
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
			session.chPush <- push
		}
	}
}

func (self grpc_hall_server) MustPush(ctx context.Context, req *hall_grpc.MustPushRequest) (*hall_grpc.MustPushResponse, error) {
	resp := &hall_grpc.MustPushResponse{}
	userId := req.UserId
	_, _, push, err := self.hall.proto.PushDecode(req.Data)
	if err != nil {
		return resp, err
	}
	if v, ok := self.hall.SessionDict.Load(userId); !ok {
		return resp, gira.ErrNoSession
	} else {
		session, _ := v.(*hall_sesssion)
		session.chPush <- push
		return resp, nil
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

func (self grpc_hall_server) GateStream(client hall_grpc.Hall_GateStreamServer) error {
	ticker := time.NewTicker(time.Duration(self.hall.config.Framework.GateReportInterval) * time.Second)
	defer func() {
		ticker.Stop()
		log.Infow("gate stream exit")
	}()
	errGroup, errCtx := errgroup.WithContext(self.hall.ctx)
	// 定时发送在线人数
	errGroup.Go(func() error {
		for {
			select {
			case <-self.hall.ctx.Done():
				return self.hall.ctx.Err()
			case <-errCtx.Done():
				return errCtx.Err()
			case <-ticker.C:
				resp := &hall_grpc.HallDataPush{}
				resp.Type = hall_grpc.HallDataType_SERVER_REPORT
				resp.Report = &hall_grpc.HallReportPush{
					PlayerCount: self.hall.SessionCount,
				}
				if err := client.Send(resp); err != nil {
					log.Infow("gate send fail", "error", err)
					return err
				}
			}
		}
	})
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

// client.Recv无法解除阻塞，要通知网送开断开连接
// 有两种断开方式，
// 网关连接断开，1.recv goroutine结束，2.recv ctrl goroutine结束 3.response goroutine马上结束, 4.ctrl goroutine马上结束 5.stream结束(释放session)
// session自己主动断开, 1.ctrl goroutine结束, 关闭response channel, 2.response goroutine处理完剩余消息后结束, 3.recv ctrl goroutinue马上结束 4.stream结束(释放session) 5.recv goroutinue结束
func (self grpc_hall_server) ClientStream(client hall_grpc.Hall_ClientStreamServer) error {
	var req *hall_grpc.StreamDataRequest
	var err error
	var memberId string

	if req, err = client.Recv(); err != nil {
		log.Errorw("recv fail", "error", err)
		return err
	}
	sessionId := req.SessionId
	// 会话的第一条消息一定要带memberid来完成登录
	memberId = req.MemberId
	log.Infow("bind memberid", "session_id", sessionId, "member_id", memberId)
	if memberId == "" {
		log.Errorw("invalid memberid", "session_id", sessionId, "member_id", memberId)
		return gira.ErrTodo
	}
	// 申请建立会话
	var session *hall_sesssion
	// TODO: 增加timeout
	session, err = newSession(self.hall, sessionId, memberId)
	if err != nil {
		log.Errorw("create session fail", "session_id", sessionId, "error", err)
		return err
	}
	go session.serve()

	// 建立成功， 函数结束后通知释放
	// response管理
	session.chResponse = make(chan *hall_grpc.StreamDataResponse, 10)
	session.chRequest = make(chan *hall_grpc.StreamDataRequest, 10)

	errGroup, errCtx := errgroup.WithContext(self.hall.ctx)
	// ctrl context
	cancelCtx, cancelFunc := context.WithCancel(self.hall.ctx)
	session.streamCancelFunc = cancelFunc
	// recv ctrl context
	defer func() {
		cancelFunc()
		if err := session.Call_close(self.hall.ctx, actor.WithCallTimeOut(5*time.Second)); err != nil {
			log.Errorw("session close fail", "error", err)
		}
	}()

	errGroup.Go(func() error {
		defer func() {
			log.Infow("response goroutine exit", "session_id", sessionId)
		}()
		for {
			select {
			case resp := <-session.chResponse:
				if resp == nil {
					// session 要求关闭，返回错误，关闭stream
					log.Infow("response chan close", "session_id", sessionId)
					return nil
				}
				if err := client.Send(resp); err != nil {
					log.Infow("client response fail", "session_id", sessionId, "error", err)
				} else {
					log.Infow("client response success", "session_id", sessionId, "data", len(resp.Data))
				}
			case <-errCtx.Done():
				return nil
			}
		}
	})

	// 这个协程不由errgroup控制，因为这个函数需要在stream函数返回后，还活着， stream不提供close函数，因此需要在stream返回后，recv才能解除阻塞
	// errgroup不提供中断方法，也只能由内部返咽错误来关闭，因为要用cancelFunc来中断errGroup
	go func() {
		defer func() {
			defer close(session.chRequest)
			log.Infow("read goroutine exit", "session_id", sessionId)
		}()
		session.chRequest <- req
		for {
			req, err = client.Recv()
			if err != nil {
				log.Infow("client recv fail", "error", err, "session_id", sessionId)
				cancelFunc()
				// 连接关闭，返回错误，会触发errCtx的cancel
				return
			} else {
				session.chRequest <- req
			}
		}
	}()
	errGroup.Go(func() error {
		// 等待被动结束
		defer func() {
			log.Infow("ctrl goroutine exit", "session_id", sessionId)
		}()
		select {
		case <-cancelCtx.Done():
			// 关闭response管道，不再接收消息
			close(session.chResponse)
			return cancelCtx.Err()
		case <-errCtx.Done():
			return errCtx.Err()
		}
	})
	err = errGroup.Wait()
	// 尝试关闭 response channel
	select {
	case <-session.chResponse:
	default:
		close(session.chResponse)
	}
	log.Infow("stream exit", "error", err, "session_id", sessionId, "member_id", memberId)
	return err
}
