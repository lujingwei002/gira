package gate

import (
	"context"
	"time"

	"github.com/lujingwei002/gira/facade"
	"github.com/lujingwei002/gira/log"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/framework/smallgame/account/jwt"
	"github.com/lujingwei002/gira/framework/smallgame/gen/grpc/hall_grpc"
	"google.golang.org/grpc"
)

type hall struct {
	proto           gira.Proto
	facade          gira.ApplicationFacade
	peer            *upstream_peer
	ctx             context.Context
	cancelFunc      context.CancelFunc
	sessionCount    int64 // 会话数量
	connectionCount int64 // 连接数量
	handler         GateHandler
	config          *Config
}

type upstream_peer struct {
	Id       int32
	FullName string
	Address  string
}

func newHall(facade gira.ApplicationFacade, proto gira.Proto, config *Config) *hall {
	return &hall{
		facade: facade,
		proto:  proto,
		config: config,
	}
}

func (h *hall) OnAwake() error {
	if handler, ok := h.facade.(GateHandler); !ok {
		return gira.ErrGateHandlerNotImplement
	} else {
		h.handler = handler
	}
	h.ctx, h.cancelFunc = context.WithCancel(facade.Context())
	return nil
}

func (h *hall) SelectPeer() *upstream_peer {
	return h.peer
}

func (h *hall) loginErrResponse(r gira.GateRequest, req gira.ProtoRequest, err error) error {
	log.Info(err)
	resp, err := h.proto.NewResponse(req)
	if err == nil {
		return err
	}
	resp.SetErrorCode(gira.ErrCode(err))
	resp.SetErrorMsg(gira.ErrMsg(err))
	if data, err := h.proto.ResponseEncode("Login", int32(r.ReqId()), resp); err != nil {
		return err
	} else {
		r.Response(data)
	}
	return nil
}

func (self *hall) OnGateStream(client gira.GateConn) {
	sessionId := client.ID()
	log.Infow("stream open", "session_id", sessionId)
	defer func() {
		log.Infow("stream close", "session_id", sessionId)
	}()

	var req gira.GateRequest
	var err error
	var memberId string
	// 需要在一定时间内发登录协议
	timeoutCtx, timeoutFunc := context.WithTimeout(self.ctx, 10*time.Second)
	defer timeoutFunc()
	req, err = client.Recv(timeoutCtx)
	if err != nil {
		log.Infow("recv fail", "session_id", sessionId, "error", err)
		return
	}
	// log.Infow("client=>gate request", "session_id", sessionId, "req_id", req.ReqId)
	if name, _, dataReq, err := self.proto.RequestDecode(req.Payload()); err != nil {
		log.Infow("client=>gate request decode fail", "session_id", sessionId, "error", err)
		return
	} else if name != "Login" {
		log.Infow("expeted login request", "name", name, "session_id", sessionId)
		return
	} else {
		loginReq, ok := dataReq.(LoginRequest)
		if !ok {
			log.Infow("login request cast fail", "session_id", sessionId)
			return
		}
		memberId = loginReq.GetMemberId()
		token := loginReq.GetToken()
		log.Infow("client login", "session_id", sessionId, "member_id", memberId)
		log.Infow("token", "session_id", sessionId, "token", token)
		// 验证memberid和token
		if claims, err := jwt.ParseJwtToken(token, self.config.Jwt.Secret); err != nil {
			log.Errorw("token无效", "session_id", sessionId, "token", token, "error", err)
			self.loginErrResponse(req, dataReq, err)
			return
		} else if claims.MemberId != memberId {
			log.Errorw("memberid无效", "token", token, "member_id", memberId, "expected_member_id", claims.MemberId)
			self.loginErrResponse(req, dataReq, gira.ErrInvalidJwt)
			return
		} else {
			memberId = claims.MemberId
		}
		server := self.SelectPeer()
		if server == nil {
			self.loginErrResponse(req, dataReq, gira.ErrUpstreamUnavailable)
			return
		}
		conn, err := grpc.Dial(server.Address, grpc.WithInsecure())
		if err != nil {
			log.Errorw("server dail fail", "error", err, "address", server.Address)
			self.loginErrResponse(req, dataReq, gira.ErrUpstreamUnreachable)
			return
		}
		streamCtx, streamCancelFunc := context.WithCancel(self.ctx)
		defer streamCancelFunc()
		grpcClient := hall_grpc.NewUpstreamClient(conn)
		stream, err := grpcClient.DataStream(streamCtx)
		if err != nil {
			self.loginErrResponse(req, dataReq, gira.ErrUpstreamUnreachable)
			return
		}
		defer stream.CloseSend()
		session := newSession(self, sessionId, memberId)
		session.serve(self.ctx, stream, client, req)
	}
}

func (self *hall) OnPeerAdd(peer *gira.Peer) {
	if peer.Name != self.config.Upstream.Name {
		return
	}
	if peer.Id != self.config.Upstream.Id {
		return
	}
	log.Infow("add upstream", "id", peer.Id, "fullname", peer.FullName, "grpc_addr", peer.GrpcAddr)
	server := &upstream_peer{
		Id:       peer.Id,
		FullName: peer.FullName,
		Address:  peer.GrpcAddr,
	}
	self.peer = server
}

func (self *hall) OnPeerDelete(peer *gira.Peer) {
	log.Debugw("OnPeerDelete")
	if peer.Name != self.config.Upstream.Name {
		return
	}
	if peer.Id != self.config.Upstream.Id {
		return
	}
	self.peer = nil
}

func (self *hall) OnPeerUpdate(peer *gira.Peer) {
	log.Debugw("OnPeerUpdate")
	if peer.Name != self.config.Upstream.Name {
		return
	}
	if peer.Id != self.config.Upstream.Id {
		return
	}
	server := &upstream_peer{
		Id:       peer.Id,
		FullName: peer.FullName,
		Address:  peer.GrpcAddr,
	}
	self.peer = server
}
