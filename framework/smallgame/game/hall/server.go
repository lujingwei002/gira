package hall

// grpc server 服务

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/codes"
	"github.com/lujingwei002/gira/errors"
	"github.com/lujingwei002/gira/facade"
	"github.com/lujingwei002/gira/framework/smallgame/gen/service/hall_grpc"
	"github.com/lujingwei002/gira/log"
	"golang.org/x/sync/errgroup"
)

type hall_server struct {
	hall_grpc.UnimplementedHallServer
	hall *hall_service
}

// 踢用户下线
func (self *hall_server) Kick(ctx context.Context, req *hall_grpc.KickRequest) (*hall_grpc.KickResponse, error) {
	resp := &hall_grpc.KickResponse{}
	if err := self.hall.Kick(ctx, req.UserId, req.Reason); err != nil {
		return nil, err
	}
	return resp, nil
}

// 顶号下线
func (self *hall_server) UserInstead(ctx context.Context, req *hall_grpc.UserInsteadRequest) (*hall_grpc.UserInsteadResponse, error) {
	resp := &hall_grpc.UserInsteadResponse{}
	reason := fmt.Sprintf("账号在%s登录", req.Address)
	if err := self.hall.Instead(ctx, req.UserId, reason); err != nil {
		resp.ErrorCode = codes.Code(err)
		resp.ErrorMsg = codes.Msg(err)
		return resp, nil
	} else {

		return resp, nil
	}
}

func (self *hall_server) SendMessage(context.Context, *hall_grpc.SendMessageRequest) (*hall_grpc.SendMessageResponse, error) {
	resp := &hall_grpc.SendMessageResponse{}
	return resp, nil
}

func (self *hall_server) CallMessage(context.Context, *hall_grpc.CallMessageRequest) (*hall_grpc.CallMessageResponse, error) {
	resp := &hall_grpc.CallMessageResponse{}
	return resp, nil
}

// push消息流
// func (self *hall_server) PushStream(server hall_grpc.Hall_PushStreamServer) error {
// 	for {
// 		req, err := server.Recv()
// 		if err == io.EOF {
// 			return nil
// 		}
// 		if err != nil {
// 			return err
// 		}
// 		userId := req.UserId
// 		_, _, push, err := self.hall.proto.PushDecode(req.Data)
// 		if err != nil {
// 			log.Warnw("user push decode fail", "error", err)
// 			continue
// 		}
// 		if v, ok := self.hall.SessionDict.Load(userId); ok {
// 			session, _ := v.(*hall_sesssion)
// 			// WARN: chPush有可能已经关闭
// 			session.chPeerPush <- push
// 		}
// 	}
// }

func (self *hall_server) MustPush(ctx context.Context, req *hall_grpc.MustPushRequest) (resp *hall_grpc.MustPushResponse, err error) {
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
	if v, ok := self.hall.sessionDict.Load(userId); !ok {
		err = errors.ErrNoSession
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

func (self *hall_server) Info(ctx context.Context, req *hall_grpc.InfoRequest) (*hall_grpc.InfoResponse, error) {
	resp := &hall_grpc.InfoResponse{
		BuildTime:    facade.GetBuildTime(),
		AppVersion:   facade.GetAppVersion(),
		SessionCount: self.hall.sessionCount,
	}
	return resp, nil
}

func (self *hall_server) HealthCheck(ctx context.Context, req *hall_grpc.HealthCheckRequest) (*hall_grpc.HealthCheckResponse, error) {
	resp := &hall_grpc.HealthCheckResponse{
		PlayerCount: self.hall.sessionCount,
	}
	resp.Status = self.hall.status
	return resp, nil
}

func (self *hall_server) GateStream(client hall_grpc.Hall_GateStreamServer) error {
	// ticker := time.NewTicker(time.Duration(self.hall.config.Framework.Hall.GatewayReportInterval) * time.Second)
	defer func() {
		// ticker.Stop()
		log.Infow("gate stream exit")
	}()
	// if self.hall.isDestory {
	// 	return gira.ErrAlreadyDestory
	// }
	errGroup, errCtx := errgroup.WithContext(self.hall.gateStreamCtx)
	// 定时发送在线人数
	// errGroup.Go(func() error {
	// 	for {
	// 		select {
	// 		case <-self.hall.gateStreamCtx.Done():
	// 			return self.hall.gateStreamCtx.Err()
	// 		case <-errCtx.Done():
	// 			return errCtx.Err()
	// 		}
	// 	}
	// })
	// 循环接收gate的命令
	errGroup.Go(func() error {
		for {
			select {
			case <-self.hall.gateStreamCtx.Done():
				return self.hall.gateStreamCtx.Err()
			case <-errCtx.Done():
				return errCtx.Err()
				// default:
				// 	req, err := client.Recv()
				// 	if err != nil {
				// 		log.Infow("gate recv fail", "error", err)
				// 		// 连接关闭，返回错误，会触发errCtx的cancel
				// 		return err
				// 	} else {
				// 		log.Infow("gate recv", "req", req)
				// 	}
			}
		}
	})
	return errGroup.Wait()
}

func (self *hall_server) ClientStream(client hall_grpc.Hall_ClientStreamServer) error {
	var req *hall_grpc.ClientMessageRequest
	var err error
	var memberId string
	// if self.hall.isDestory {
	// 	return gira.ErrAlreadyDestory
	// }
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
		return errors.ErrInvalidMemberId
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
	session.chClientResponse = make(chan *hall_grpc.ClientMessageResponse, hall.config.ResponseBufferSize)
	session.chClientRequest = make(chan *hall_grpc.ClientMessageRequest, hall.config.RequestBufferSize)

	// 绑定到hall
	clientCtx, clientCancelFunc := context.WithCancel(self.hall.backgroundCtx)
	session.clientCancelFunc = clientCancelFunc
	// recv ctrl context
	defer func() {
		clientCancelFunc()
		if err := session.Close(self.hall.backgroundCtx); err != nil {
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
