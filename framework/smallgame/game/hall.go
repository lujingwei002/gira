package game

import (
	"context"
	"fmt"
	"io"
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

type hall struct {
	*actor.Actor
	ctx           context.Context
	cancelFunc    context.CancelFunc
	SessionDict   sync.Map
	hallHandler   HallHandler
	playerHandler gira.ProtoHandler
	proto         gira.Proto
}

func newHall(proto gira.Proto, hallHandler HallHandler, playerHandler gira.ProtoHandler) *hall {
	return &hall{
		proto:         proto,
		hallHandler:   hallHandler,
		playerHandler: playerHandler,
		Actor:         actor.NewActor(16),
	}
}

func (hall *hall) Register(server *grpc.Server) error {
	// 注册grpc
	hall_grpc.RegisterUpstreamServer(server, upstream_server{
		hall: hall,
	})
	hall_grpc.RegisterHallServer(server, &hall_server{
		hall: hall,
	})
	go hall.serve()
	return nil
}

func (hall *hall) serve() {
	hall.ctx, hall.cancelFunc = context.WithCancel(facade.Context())
	defer hall.cancelFunc()
	for {
		select {
		case r := <-hall.Inbox():
			r.Call()
		case <-hall.ctx.Done():
			log.Infow("hall exit")
			return
		}
	}
}

func (hall *hall) Push(ctx context.Context, userId string, req gira.ProtoPush) error {
	var err error
	if v, ok := hall.SessionDict.Load(userId); !ok {
		var data []byte
		var peer *gira.Peer
		data, err = hall.proto.PushEncode(req)
		if err != nil {
			return err
		}
		peer, err = facade.WhereIsUser(userId)
		if err != nil {
			return err
		}
		if err = rpc.Hall.Push(ctx, peer, userId, data); err != nil {
			return err
		}
		return nil
	} else {
		session, _ := v.(*hall_sesssion)
		session.chPush <- req
		return nil
	}
}

func (hall *hall) MustPush(ctx context.Context, userId string, req gira.ProtoPush) error {
	var err error
	if v, ok := hall.SessionDict.Load(userId); !ok {
		var data []byte
		var peer *gira.Peer
		data, err = hall.proto.PushEncode(req)
		if err != nil {
			return err
		}
		peer, err = facade.WhereIsUser(userId)
		if err != nil {
			return err
		}
		if _, err = rpc.Hall.MustPush(ctx, peer, userId, data); err != nil {
			return err
		}
		return nil
	} else {
		session, _ := v.(*hall_sesssion)
		session.chPush <- req
		return nil
	}
}

type hall_server struct {
	hall_grpc.UnsafeHallServer
	hall *hall
}

func (self hall_server) UserInstead(ctx context.Context, req *hall_grpc.UserInsteadRequest) (*hall_grpc.UserInsteadResponse, error) {
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

func (self hall_server) PushStream(server hall_grpc.Hall_PushStreamServer) error {
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

func (self hall_server) MustPush(ctx context.Context, req *hall_grpc.MustPushRequest) (*hall_grpc.MustPushResponse, error) {
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
type upstream_server struct {
	hall_grpc.UnsafeUpstreamServer
	hall *hall
}

func (h upstream_server) SayHello(ctx context.Context, in *hall_grpc.HelloRequest) (*hall_grpc.HelloResponse, error) {
	resp := new(hall_grpc.HelloResponse)
	resp.Replay = in.Name
	log.Debugw("upstream server recv", "name", in.Name)
	return resp, nil
}

// client.Recv无法解除阻塞，要通知网送开断开连接
// 有两种断开方式，
// 网关连接断开，1.recv goroutine结束，2.recv ctrl goroutine结束 3.response goroutine马上结束, 4.ctrl goroutine马上结束 5.stream结束(释放session)
// session自己主动断开, 1.ctrl goroutine结束, 关闭response channel, 2.response goroutine处理完剩余消息后结束, 3.recv ctrl goroutinue马上结束 4.stream结束(释放session) 5.recv goroutinue结束
func (self upstream_server) DataStream(client hall_grpc.Upstream_DataStreamServer) error {
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
	// TODO 增加timeout
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
		session.Call_close(self.hall.ctx, actor.WithCallTimeOut(5*time.Second))
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
			//expect
		}
		return nil
	})
	err = errGroup.Wait()
	// 尝试关闭 response channel
	select {
	case <-session.chResponse:
	default:
		close(session.chResponse)
	}
	log.Infow("stream exit", "error", err, "session_id", sessionId, "member_id", memberId)
	return nil
}
