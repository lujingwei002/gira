package gate

import (
	"context"
	"sync/atomic"

	"github.com/lujingwei002/gira/log"
	"golang.org/x/sync/errgroup"

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
	hall       *hall
}

func newSession(hall *hall, sessionId uint64, memberId string) *client_session {
	return &client_session{
		hall:      hall,
		sessionId: sessionId,
		memberId:  memberId,
	}
}

func (session *client_session) serve(ctx context.Context, stream hall_grpc.Upstream_DataStreamClient, client gira.GateConn, req gira.GateRequest) error {
	sessionId := session.sessionId
	var err error
	session.stream = stream
	session.client = client
	log.Infow("session open", "session_id", sessionId)
	defer func() {
		log.Infow("session close", "session_id", sessionId)
	}()
	session.ctx, session.cancelFunc = context.WithCancel(ctx)
	errGroup, _ := errgroup.WithContext(ctx)
	atomic.AddInt64(&session.hall.sessionCount, 1)
	defer atomic.AddInt64(&session.hall.sessionCount, -1)
	// 将上游消息转发到客户端
	errGroup.Go(func() error {
		defer func() {
			log.Infow("upstream=>client goroutine close", "session_id", sessionId)
		}()
		for {
			// 上游关闭时，stream并不会返回，会一直阻塞
			if resp, err := stream.Recv(); err == nil {
				session.streamResponse(resp)
			} else {
				log.Infow("上游连接关闭", "session_id", sessionId, "error", err)
				client.Close()
				return err
			}
		}
	})
	// 转发消息协程
	errGroup.Go(func() error {
		defer func() {
			log.Infow("client=>upstream goroutine close", "session_id", sessionId)
		}()
		// 接收客户端消息
		if err := session.clientRequest(req); err != nil {
			log.Infow("request fail", "session_id", sessionId, "error", err)
		}
		for {
			req, err = client.Recv(session.ctx)
			if err != nil {
				log.Infow("recv fail", "session_id", sessionId, "error", err)
				stream.CloseSend()
				return err
			}
			if err := session.clientRequest(req); err != nil {
				log.Warnw("request fail", "session_id", sessionId, "error", err)

			}
		}
	})
	err = errGroup.Wait()
	return err
}

// 收到客户端的消息
func (self *client_session) clientRequest(req gira.GateRequest) error {
	sessionId := self.sessionId
	memberId := self.memberId
	log.Infow("client=>upstream", "session_id", sessionId, "len", len(req.Payload()), "req_id", req.ReqId())
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

// 收到上游的消息
func (self *client_session) streamResponse(resp *hall_grpc.StreamDataResponse) error {
	sessionId := self.sessionId
	log.Infow("upstream=>client", "session_id", sessionId, "type", resp.Type, "len", len(resp.Data), "req_id", resp.ReqId, "data", resp.Data)
	switch resp.Type {
	case hall_grpc.PacketType_DATA:
		if resp.ReqId != 0 {
			self.client.Response(resp.ReqId, resp.Data)
		} else {
			self.client.Push("", resp.Data)
		}
	case hall_grpc.PacketType_USER_INSTEAD:
		self.client.Kick("账号在其他地方登录")
	}
	return nil
}
