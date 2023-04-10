package gate

import (
	"context"
	"sync"
	"time"

	"github.com/lujingwei002/gira/facade"
	"github.com/lujingwei002/gira/log"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/framework/smallgame/account/jwt"
)

type hall_server struct {
	application     *Framework
	proto           gira.Proto
	peers           sync.Map
	ctx             context.Context
	cancelFunc      context.CancelFunc
	sessionCount    int64 // 会话数量
	connectionCount int64 // 连接数量
	handler         GateHandler
	config          *Config
}

func newHall(application *Framework, proto gira.Proto, config *Config) *hall_server {
	return &hall_server{
		application: application,
		proto:       proto,
		config:      config,
	}
}

func (hall *hall_server) OnAwake() error {
	var a interface{} = hall.application
	if handler, ok := a.(GateHandler); !ok {
		return gira.ErrGateHandlerNotImplement
	} else {
		hall.handler = handler
	}
	hall.ctx, hall.cancelFunc = context.WithCancel(facade.Context())
	return nil
}

// 选择一个节点
func (hall *hall_server) SelectPeer() *upstream_peer {
	var minPlayerCount int64 = 0xffffffff
	var selected *upstream_peer
	hall.peers.Range(func(key any, value any) bool {
		server := value.(*upstream_peer)
		if server.client != nil && server.playerCount < minPlayerCount {
			minPlayerCount = server.playerCount
			selected = server
		}
		return true
	})
	return selected
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

func (hall *hall_server) OnClientStream(client gira.GateConn) {
	sessionId := client.ID()
	log.Infow("client stream open", "session_id", sessionId)
	defer func() {
		log.Infow("client stream exit", "session_id", sessionId)
	}()
	if hall.config.FrameWork.MaxSessionCount == -1 {
		log.Infow("client stream close", "session_id", sessionId)
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
		session.serve(server, client, req, dataReq)
	}
}

func (hall *hall_server) OnPeerAdd(peer *gira.Peer) {
	if peer.Name != hall.config.Upstream.Name {
		return
	}
	if peer.Id < hall.config.Upstream.MinId || peer.Id > hall.config.Upstream.MaxId {
		return
	}
	log.Infow("add upstream", "id", peer.Id, "fullname", peer.FullName, "grpc_addr", peer.GrpcAddr)
	ctx, cancelFUnc := context.WithCancel(hall.ctx)
	server := &upstream_peer{
		Id:         peer.Id,
		FullName:   peer.FullName,
		Address:    peer.GrpcAddr,
		ctx:        ctx,
		cancelFunc: cancelFUnc,
		isAlive:    true,
	}
	if v, ok := hall.peers.Load(peer.Id); ok {
		lastServer := v.(*upstream_peer)
		if lastServer.isSuspend {
			lastServer.Resume()
		} else {
			lastServer.close()
		}
	} else {
		hall.peers.Store(peer.Id, server)
		go server.serve()
	}
}

func (hall *hall_server) OnPeerDelete(peer *gira.Peer) {
	if peer.Name != hall.config.Upstream.Name {
		return
	}
	if peer.Id < hall.config.Upstream.MinId || peer.Id > hall.config.Upstream.MaxId {
		return
	}
	log.Infow("remove upstream", "id", peer.Id, "fullname", peer.FullName, "grpc_addr", peer.GrpcAddr)
	if v, ok := hall.peers.Load(peer.Id); ok {
		lastServer := v.(*upstream_peer)
		if lastServer.isSuspend {
			lastServer.Suspend()
		} else {
			hall.peers.Delete(peer.Id)
			lastServer.close()
		}
	}
}

func (hall *hall_server) OnPeerUpdate(peer *gira.Peer) {
	if peer.Name != hall.config.Upstream.Name {
		return
	}
	if peer.Id < hall.config.Upstream.MinId || peer.Id > hall.config.Upstream.MaxId {
		return
	}
	log.Infow("update upstream", "id", peer.Id, "fullname", peer.FullName, "grpc_addr", peer.GrpcAddr)
	if v, ok := hall.peers.Load(peer.Id); ok {
		lastServer := v.(*upstream_peer)
		lastServer.Address = peer.GrpcAddr
	}
}
