package game

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/actor"
	"github.com/lujingwei002/gira/facade"
	"github.com/lujingwei002/gira/framework/smallgame/gen/grpc/hall_grpc"
	"github.com/lujingwei002/gira/log"
	"github.com/lujingwei002/gira/sproto"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func newHallServer(proto *sproto.Sproto, hallHandler HallHandler, playerHandler *sproto.SprotoHandler) *hall_server {
	return &hall_server{
		hallHandler:   hallHandler,
		proto:         proto,
		playerHandler: playerHandler,
		Actor:         actor.NewActor(16),
	}
}

func (hall *hall_server) Register(server *grpc.Server) error {
	// 注册grpc
	hall_grpc.RegisterUpstreamServer(server, upstream_server{
		hall: hall,
	})
	hall_grpc.RegisterHallServer(server, hall)
	go hall.serve()
	return nil
}

func (hall *hall_server) serve() {
	log.Printf(" %p\n", hall)
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

func (hall *hall_server) Push(ctx context.Context, memberId string, req sproto.SprotoPush) error {
	if v, ok := hall.SessionDict.Load(memberId); !ok {
		return gira.ErrNoSession
	} else {
		session, _ := v.(*hall_sesssion)
		session.chPush <- req
		return nil
	}
}

type hall_server struct {
	hall_grpc.UnsafeHallServer

	*actor.Actor
	ctx           context.Context
	cancelFunc    context.CancelFunc
	SessionDict   sync.Map
	hallHandler   HallHandler
	playerHandler *sproto.SprotoHandler
	proto         *sproto.Sproto
}

func (self hall_server) UserInstead(ctx context.Context, req *hall_grpc.UserInsteadRequest) (*hall_grpc.UserInsteadResponse, error) {
	resp := &hall_grpc.UserInsteadResponse{}
	userId := req.UserId
	log.Infow("user instead", "user_id", userId)
	if v, ok := self.SessionDict.Load(userId); !ok {
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
		if err := session.Call_instead(timeoutCtx, fmt.Sprintf("在%s登录", req.Address)); err != nil {
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

type upstream_server struct {
	hall_grpc.UnsafeUpstreamServer
	hall *hall_server
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
	session, err = self.hall.createSession(self.hall.ctx, sessionId, memberId)
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
		defer cancelFunc()
		defer session.OnStreamClose()
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
	log.Infow("stream quit", "error", err, "session_id", sessionId)
	return nil
}
