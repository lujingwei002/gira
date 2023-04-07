package gate

import (
	"context"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/lujingwei002/gira/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/framework/smallgame/gen/grpc/hall_grpc"
)

type client_session struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	sessionId  uint64
	memberId   string
	client     gira.GateConn
	stream     hall_grpc.Upstream_DataStreamClient
	hall       *hall_server
	hall1      *hall_server
}

func newSession(hall *hall_server, sessionId uint64, memberId string) *client_session {
	return &client_session{
		hall:      hall,
		sessionId: sessionId,
		memberId:  memberId,
	}
}

func (session *client_session) serve(ctx context.Context, server *upstream_peer, client gira.GateConn, req gira.GateRequest, dataReq gira.ProtoRequest) error {
	sessionId := session.sessionId
	var err error
	log.Infow("session open", "session_id", sessionId)
	// 连接上游服务器
	conn, err := grpc.Dial(server.Address, grpc.WithInsecure())
	if err != nil {
		log.Errorw("server dail fail", "error", err, "address", server.Address)
		session.hall.loginErrResponse(req, dataReq, gira.ErrUpstreamUnreachable)
		return err
	}
	streamCtx, streamCancelFunc := context.WithCancel(ctx)
	defer streamCancelFunc()
	grpcClient := hall_grpc.NewUpstreamClient(conn)
	stream, err := grpcClient.DataStream(streamCtx)
	if err != nil {
		session.hall.loginErrResponse(req, dataReq, gira.ErrUpstreamUnreachable)
		return err
	}
	session.stream = stream
	session.client = client
	session.ctx, session.cancelFunc = context.WithCancel(ctx)
	errGroup, _ := errgroup.WithContext(ctx)
	atomic.AddInt64(&session.hall.sessionCount, 1)
	defer func() {
		log.Infow("session close", "session_id", sessionId)
		atomic.AddInt64(&session.hall.sessionCount, -1)
		stream.CloseSend()
		session.cancelFunc()
	}()
	// 将上游消息转发到客户端
	errGroup.Go(func() (err error) {
		defer func() {
			log.Infow("upstream=>client goroutine close", "session_id", sessionId)
			if e := recover(); e != nil {
				log.Error(e)
				debug.PrintStack()
				err = e.(error)
				client.Close()
			}
		}()
		for {
			var resp *hall_grpc.StreamDataResponse
			// 上游关闭时，stream并不会返回，会一直阻塞
			if resp, err = stream.Recv(); err == nil {
				session.processStreamResponse(resp)
			} else {
				log.Infow("上游连接关闭", "session_id", sessionId, "error", err)
				serverErr, ok := status.FromError(err)
				if ok {
					log.Println("ccccccccccc", serverErr.Code())
					log.Println("ccccccccccc", serverErr.Err())
					if serverErr.Code() == codes.Code(333) {
						log.Println("dddddddddddddd")
						session.stream = nil
						// 再次连接上游服务器
						for {
							conn, err := grpc.Dial(server.Address, grpc.WithInsecure())
							if err != nil {
								log.Errorw("server dail fail", "error", err, "address", server.Address)
								time.Sleep(1 * time.Second)
								continue
							}
							streamCtx, streamCancelFunc = context.WithCancel(ctx)
							defer streamCancelFunc()
							grpcClient = hall_grpc.NewUpstreamClient(conn)
							stream, err = grpcClient.DataStream(streamCtx)
							if err != nil {
								log.Errorw("server dail fail", "error", err, "address", server.Address)
								time.Sleep(1 * time.Second)
								continue

							}
							session.stream = stream
							break
						}
					} else {
						client.Close()
						return
					}
				} else {
					client.Close()
					return

				}
			}
		}
	})
	// 转发消息协程
	errGroup.Go(func() (err error) {
		defer func() {
			log.Infow("client=>upstream goroutine close", "session_id", sessionId)
			if e := recover(); e != nil {
				log.Error(e)
				debug.PrintStack()
				err = e.(error)
				client.Close()
			}
		}()
		// 接收客户端消息
		if err = session.processClientRequest(req); err != nil {
			log.Infow("request fail", "session_id", sessionId, "error", err)
		}
		for {
			req, err = client.Recv(session.ctx)
			if err != nil {
				log.Infow("recv fail", "session_id", sessionId, "error", err)
				stream.CloseSend()
				return err
			}
			if err = session.processClientRequest(req); err != nil {
				log.Warnw("request fail", "session_id", sessionId, "error", err)
			}
		}
	})
	err = errGroup.Wait()
	log.Infow("session wait", "error", err)
	return err
}

// 处理客户端的消息
func (self *client_session) processClientRequest(req gira.GateRequest) error {
	sessionId := self.sessionId
	memberId := self.memberId
	log.Infow("client=>upstream", "session_id", sessionId, "len", len(req.Payload()), "req_id", req.ReqId())
	if self.stream == nil {
		log.Warnw("当前服务器不可以用，无法转发")
		return nil
	} else {
		if err := self.stream.Send(&hall_grpc.StreamDataRequest{
			MemberId:  memberId,
			SessionId: sessionId,
			ReqId:     req.ReqId(),
			Data:      req.Payload(),
		}); err != nil {
			log.Errorw("client=>upstream fail", "session_id", sessionId, "error", err)
			return err
		}
		return nil
	}
}

// 处理上游的消息
func (session *client_session) processStreamResponse(resp *hall_grpc.StreamDataResponse) error {
	sessionId := session.sessionId
	log.Infow("upstream=>client", "session_id", sessionId, "type", resp.Type, "route", resp.Route, "len", len(resp.Data), "req_id", resp.ReqId)
	switch resp.Type {
	case hall_grpc.PacketType_DATA:
		if resp.ReqId != 0 {
			session.client.Response(resp.ReqId, resp.Data)
		} else {
			session.client.Push("", resp.Data)
		}
	case hall_grpc.PacketType_USER_INSTEAD:
		session.client.Kick("账号在其他地方登录")
	case hall_grpc.PacketType_SERVER_DOWN:
		session.client.Kick("服务器关闭")
	}
	return nil
}
