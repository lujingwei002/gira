package gate

import (
	"context"
	"time"

	"github.com/lujingwei002/gira/facade"
	"github.com/lujingwei002/gira/log"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/framework/smallgame/account/jwt"
)

type hall_server struct {
	app             *GateApplication
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

func newHall(app *GateApplication, facade gira.ApplicationFacade, proto gira.Proto, config *Config) *hall_server {
	return &hall_server{
		app:    app,
		facade: facade,
		proto:  proto,
		config: config,
	}
}

func (hall *hall_server) OnAwake() error {
	if handler, ok := hall.facade.(GateHandler); !ok {
		return gira.ErrGateHandlerNotImplement
	} else {
		hall.handler = handler
	}
	hall.ctx, hall.cancelFunc = context.WithCancel(facade.Context())
	return nil
}

func (hall *hall_server) SelectPeer() *upstream_peer {
	return hall.peer
}

func (hall *hall_server) loginErrResponse(r gira.GateRequest, req gira.ProtoRequest, err error) error {
	log.Info(err)
	resp, err := hall.proto.NewResponse(req)
	if err == nil {
		return err
	}
	resp.SetErrorCode(gira.ErrCode(err))
	resp.SetErrorMsg(gira.ErrMsg(err))
	if data, err := hall.proto.ResponseEncode("Login", int32(r.ReqId()), resp); err != nil {
		return err
	} else {
		r.Response(data)
	}
	return nil
}

func (hall *hall_server) OnGateStream(client gira.GateConn) {
	sessionId := client.ID()
	log.Infow("stream open", "session_id", sessionId)
	defer func() {
		log.Infow("stream close", "session_id", sessionId)
	}()
	if hall.config.FrameWork.MaxSessionCount == -1 {
		log.Infow("stream close", "session_id", sessionId)
	} else if hall.config.FrameWork.MaxSessionCount == 0 || hall.sessionCount >= hall.config.FrameWork.MaxSessionCount {
		log.Warnw("reach max session count", "session_id", sessionId, "session_count", hall.sessionCount, "max_session_count", hall.config.FrameWork.MaxSessionCount)
		return
	}
	var req gira.GateRequest
	var err error
	var memberId string
	// 需要在一定时间内发登录协议
	timeoutCtx, timeoutFunc := context.WithTimeout(hall.ctx, time.Duration(hall.config.FrameWork.WaitLoginTimeout)*time.Second)
	defer timeoutFunc()
	req, err = client.Recv(timeoutCtx)
	if err != nil {
		log.Infow("recv fail", "session_id", sessionId, "error", err)
		return
	}
	// log.Infow("client=>gate request", "session_id", sessionId, "req_id", req.ReqId)
	if name, _, dataReq, err := hall.proto.RequestDecode(req.Payload()); err != nil {
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
		if claims, err := jwt.ParseJwtToken(token, hall.config.Jwt.Secret); err != nil {
			log.Errorw("token无效", "session_id", sessionId, "token", token, "error", err)
			hall.loginErrResponse(req, dataReq, err)
			return
		} else if claims.MemberId != memberId {
			log.Errorw("memberid无效", "token", token, "member_id", memberId, "expected_member_id", claims.MemberId)
			hall.loginErrResponse(req, dataReq, gira.ErrInvalidJwt)
			return
		} else {
			memberId = claims.MemberId
		}
		server := hall.SelectPeer()
		if server == nil {
			hall.loginErrResponse(req, dataReq, gira.ErrUpstreamUnavailable)
			return
		}

		session := newSession(hall, sessionId, memberId)
		session.serve(hall.ctx, server, client, req, dataReq)
	}
}

func (self *hall_server) OnPeerAdd(peer *gira.Peer) {
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

func (hall *hall_server) OnPeerDelete(peer *gira.Peer) {
	if peer.Name != hall.config.Upstream.Name {
		return
	}
	if peer.Id != hall.config.Upstream.Id {
		return
	}
	log.Infow("remove upstream", "id", peer.Id, "fullname", peer.FullName, "grpc_addr", peer.GrpcAddr)
	hall.peer = nil
}

func (hall *hall_server) OnPeerUpdate(peer *gira.Peer) {
	if peer.Name != hall.config.Upstream.Name {
		return
	}
	if peer.Id != hall.config.Upstream.Id {
		return
	}
	server := &upstream_peer{
		Id:       peer.Id,
		FullName: peer.FullName,
		Address:  peer.GrpcAddr,
	}
	log.Infow("update upstream", "id", peer.Id, "fullname", peer.FullName, "grpc_addr", peer.GrpcAddr)
	hall.peer = server
}
