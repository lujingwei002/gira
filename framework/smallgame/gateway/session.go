package gateway

import (
	"context"
	"io"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/lujingwei002/gira/log"
	"github.com/lujingwei002/gira/service/hall/hall_grpc"
	"golang.org/x/sync/errgroup"

	"github.com/lujingwei002/gira"
)

type client_session struct {
	ctx             context.Context
	cancelFunc      context.CancelFunc
	sessionId       uint64
	memberId        string
	client          gira.GatewayConn
	stream          hall_grpc.Hall_ClientStreamClient
	hall            *hall_server
	pendingMessages []gira.GatewayMessage
}

func newSession(hall *hall_server, sessionId uint64, memberId string) *client_session {
	return &client_session{
		hall:            hall,
		sessionId:       sessionId,
		memberId:        memberId,
		pendingMessages: make([]gira.GatewayMessage, 0),
	}
}

func (session *client_session) serve(client gira.GatewayConn, message gira.GatewayMessage, dataReq gira.ProtoRequest) error {
	sessionId := session.sessionId
	var err error
	var stream hall_grpc.Hall_ClientStreamClient
	hall := session.hall
	log.Infow("session open", "session_id", sessionId)
	server := hall.SelectPeer()
	if server == nil {
		hall.loginErrResponse(message, dataReq, gira.ErrUpstreamUnavailable)
		return gira.ErrUpstreamUnavailable
	}
	// stream绑定server
	streamCtx, streamCancelFunc := context.WithCancel(server.ctx)
	defer func() {
		streamCancelFunc()
	}()
	stream, err = server.NewClientStream(streamCtx)
	if err != nil {
		session.hall.loginErrResponse(message, dataReq, gira.ErrUpstreamUnreachable)
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
			var resp *hall_grpc.ClientMessageResponse
			// 上游关闭时，stream并不会返回，会一直阻塞
			if resp, err = stream.Recv(); err == nil {
				session.processStreamMessage(resp)
			} else if err != io.EOF {
				log.Infow("上游连接异常关闭", "session_id", sessionId, "error", err)
				session.cancelFunc()
				return err
			} else {
				select {
				case <-session.ctx.Done():
					return session.ctx.Err()
				default:
				}
				log.Infow("上游连接正常关闭", "session_id", sessionId, "error", err)
				session.stream = nil
				client.SendServerSuspend("")
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
							client.SendServerResume("")
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
		if err = session.processClientMessage(message); err != nil {
			log.Infow("client=>upstream request fail", "session_id", sessionId, "error", err)
			session.pendingMessages = append(session.pendingMessages, message)
		}
		for {
			message, err = client.Recv(session.ctx)
			if err != nil {
				log.Infow("recv fail", "session_id", sessionId, "error", err)
				session.cancelFunc()
				return err
			}
			for len(session.pendingMessages) > 0 {
				p := session.pendingMessages[0]
				if err = session.processClientMessage(p); err != nil {
					break
				} else {
					log.Infow("补发消息", "session_id", sessionId, "req_id", message.ReqId())
					session.pendingMessages = session.pendingMessages[1:]
				}
			}
			if len(session.pendingMessages) > 0 {
				session.pendingMessages = append(session.pendingMessages, message)
				continue
			}
			if err = session.processClientMessage(message); err != nil {
				log.Warnw("client=>upstream request fail", "session_id", sessionId, "error", err)
				session.pendingMessages = append(session.pendingMessages, message)
			}
		}
	})
	err = errGroup.Wait()
	log.Infow("session wait", "error", err)
	return err
}

// 处理客户端的消息
func (self *client_session) processClientMessage(message gira.GatewayMessage) error {
	sessionId := self.sessionId
	memberId := self.memberId
	log.Infow("client=>upstream", "session_id", sessionId, "len", len(message.Payload()), "req_id", message.ReqId())
	if self.stream == nil {
		log.Warnw("当前服务器不可以用，无法转发", "req_id", message.ReqId())
		return gira.ErrTodo
	} else {
		data := &hall_grpc.ClientMessageRequest{
			MemberId:  memberId,
			SessionId: sessionId,
			ReqId:     message.ReqId(),
			Data:      message.Payload(),
		}
		if err := self.stream.Send(data); err != nil {
			return err
		}
		return nil
	}
}

// 处理上游的消息
func (session *client_session) processStreamMessage(message *hall_grpc.ClientMessageResponse) error {
	sessionId := session.sessionId
	log.Infow("upstream=>client", "session_id", sessionId, "type", message.Type, "route", message.Route, "len", len(message.Data), "req_id", message.ReqId)

	switch message.Type {
	case hall_grpc.PacketType_DATA:
		if message.ReqId != 0 {
			session.client.Response(message.ReqId, message.Data)
		} else {
			session.client.Push("", message.Data)
		}
	case hall_grpc.PacketType_USER_INSTEAD:
		session.client.Kick(string(message.Data))
	case hall_grpc.PacketType_KICK:
		session.client.Kick(string(message.Data))
	}
	return nil
}
