package gate

import (
	"context"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/lujingwei002/gira/log"
	"golang.org/x/sync/errgroup"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/framework/smallgame/gen/grpc/hall_grpc"
)

type client_session struct {
	ctx             context.Context
	cancelFunc      context.CancelFunc
	sessionId       uint64
	memberId        string
	client          gira.GateConn
	stream          hall_grpc.Hall_ClientStreamClient
	hall            *hall_server
	hall1           *hall_server
	pendingRequests []gira.GateRequest
}

func newSession(hall *hall_server, sessionId uint64, memberId string) *client_session {
	return &client_session{
		hall:            hall,
		sessionId:       sessionId,
		memberId:        memberId,
		pendingRequests: make([]gira.GateRequest, 0),
	}
}

func (session *client_session) serve(client gira.GateConn, req gira.GateRequest, dataReq gira.ProtoRequest) error {
	sessionId := session.sessionId
	var err error
	var stream hall_grpc.Hall_ClientStreamClient
	hall := session.hall
	log.Infow("session open", "session_id", sessionId)
	server := hall.SelectPeer()
	if server == nil {
		hall.loginErrResponse(req, dataReq, gira.ErrUpstreamUnavailable)
		return gira.ErrUpstreamUnavailable
	}
	// stream绑定server
	streamCtx, streamCancelFunc := context.WithCancel(server.ctx)
	defer func() {
		streamCancelFunc()
	}()
	stream, err = server.NewClientStream(streamCtx)
	if err != nil {
		session.hall.loginErrResponse(req, dataReq, gira.ErrUpstreamUnreachable)
		return err
	}
	session.stream = stream
	session.client = client
	// session 绑定 hall
	session.ctx, session.cancelFunc = context.WithCancel(hall.ctx)
	errGroup, _ := errgroup.WithContext(session.ctx)
	atomic.AddInt64(&session.hall.sessionCount, 1)
	defer func() {
		log.Infow("session close", "session_id", sessionId)
		atomic.AddInt64(&session.hall.sessionCount, -1)
		session.cancelFunc()
	}()
	// 将上游消息转发到客户端
	errGroup.Go(func() (err error) {
		defer func() {
			log.Infow("upstream=>client goroutine exit", "session_id", sessionId)
			if e := recover(); e != nil {
				log.Error(e)
				debug.PrintStack()
				err = e.(error)
				session.cancelFunc()
			}
		}()
		for {
			var resp *hall_grpc.StreamDataResponse
			// 上游关闭时，stream并不会返回，会一直阻塞
			if resp, err = stream.Recv(); err == nil {
				session.processStreamResponse(resp)
			} else {
				select {
				case <-session.ctx.Done():
					return session.ctx.Err()
				default:
				}
				log.Infow("上游连接关闭", "session_id", sessionId, "error", err)
				session.stream = nil
				// 重新选择节点
				for {
					// log.Infow("重新选择节点", "session_id", sessionId)
					server = hall.SelectPeer()
					if server != nil {
						// log.Infow("重新选择节点", "session_id", sessionId, "full_name", server.FullName)
						streamCancelFunc()
						streamCtx, streamCancelFunc = context.WithCancel(server.ctx)
						stream, err = server.NewClientStream(streamCtx)
						if err != nil {
							streamCancelFunc()
							select {
							case <-session.ctx.Done():
								return session.ctx.Err()
							default:
								time.Sleep(1 * time.Second)
							}
						} else {
							session.stream = stream
							log.Infow("重新选择节点, 连接成功", "session_id", sessionId, "full_name", server.FullName)
							break
						}
					} else {
						select {
						case <-session.ctx.Done():
							return session.ctx.Err()
						default:
							time.Sleep(1 * time.Second)
						}
					}
				}
			}
		}
	})
	errGroup.Go(func() (err error) {
		select {
		case <-session.ctx.Done():
			client.Close()
			streamCancelFunc()
			return session.ctx.Err()
		}
	})
	// 转发消息协程
	errGroup.Go(func() (err error) {
		defer func() {
			log.Infow("client=>upstream goroutine exit", "session_id", sessionId)
			if e := recover(); e != nil {
				log.Error(e)
				debug.PrintStack()
				err = e.(error)
				session.cancelFunc()
			}
		}()
		// 接收客户端消息
		if err = session.processClientRequest(req); err != nil {
			log.Infow("client=>upstream request fail", "session_id", sessionId, "error", err)
			session.pendingRequests = append(session.pendingRequests, req)
		}
		for {
			req, err = client.Recv(session.ctx)
			if err != nil {
				log.Infow("recv fail", "session_id", sessionId, "error", err)
				session.cancelFunc()
				return err
			}
			for len(session.pendingRequests) > 0 {
				p := session.pendingRequests[0]
				if err = session.processClientRequest(p); err != nil {
					break
				} else {
					log.Infow("补发消息", "session_id", sessionId, "req_id", req.ReqId())
					session.pendingRequests = session.pendingRequests[1:]
				}
			}
			if len(session.pendingRequests) > 0 {
				session.pendingRequests = append(session.pendingRequests, req)
				continue
			}
			if err = session.processClientRequest(req); err != nil {
				log.Warnw("client=>upstream request fail", "session_id", sessionId, "error", err)
				session.pendingRequests = append(session.pendingRequests, req)
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
		log.Warnw("当前服务器不可以用，无法转发", "req_id", req.ReqId())
		return gira.ErrTodo
	} else {
		data := &hall_grpc.StreamDataRequest{
			MemberId:  memberId,
			SessionId: sessionId,
			ReqId:     req.ReqId(),
			Data:      req.Payload(),
		}
		if err := self.stream.Send(data); err != nil {
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
