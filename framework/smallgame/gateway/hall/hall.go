package hall

import (
	"context"
	"sync"
	"time"

	"github.com/lujingwei002/gira/facade"
	"github.com/lujingwei002/gira/log"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/framework/smallgame/gateway"
	"github.com/lujingwei002/gira/framework/smallgame/gateway/config"
	jwt "github.com/lujingwei002/gira/jwt/account"
)

type HallServer struct {
	proto      gira.Proto
	peers      sync.Map
	ctx        context.Context
	cancelFunc context.CancelFunc

	SessionCount    int64 // 会话数量
	ConnectionCount int64 // 连接数量
}

func NewHall(proto gira.Proto) *HallServer {
	return &HallServer{
		proto: proto,
	}
}

func (hall *HallServer) OnCreate() error {
	hall.ctx, hall.cancelFunc = context.WithCancel(facade.Context())
	return nil
}

// 选择一个节点
func (hall *HallServer) SelectPeer() *upstream_peer {
	var selected *upstream_peer
	// 最大版本
	var maxBuildTime int64 = 0
	hall.peers.Range(func(key any, value any) bool {
		server := value.(*upstream_peer)
		if server.client != nil && server.buildTime > maxBuildTime {
			maxBuildTime = server.buildTime
		}
		return true
	})
	// 选择人数最小的
	var minPlayerCount int64 = 0xffffffff
	hall.peers.Range(func(key any, value any) bool {
		server := value.(*upstream_peer)
		if server.client != nil && server.buildTime == maxBuildTime && server.playerCount < minPlayerCount {
			minPlayerCount = server.playerCount
			selected = server
		}
		return true
	})
	return selected
}

func (hall *HallServer) loginErrResponse(message gira.GatewayMessage, req gira.ProtoRequest, err error) error {
	log.Info(err)
	resp, err := hall.proto.NewResponse(req)
	if err == nil {
		return err
	}
	resp.SetErrorCode(gira.ErrCode(err))
	resp.SetErrorMsg(gira.ErrMsg(err))
	if data, err := hall.proto.ResponseEncode("Login", int32(message.ReqId()), resp); err != nil {
		return err
	} else {
		message.Response(data)
	}
	return nil
}

// 客户端流数据
func (hall *HallServer) ServeClientStream(client gira.GatewayConn) {
	sessionId := client.Id()
	// 最大人数判断
	if config.Gateway.Framework.Gateway.MaxSessionCount == -1 {
		log.Infow("client stream close", "session_id", sessionId)
	} else if config.Gateway.Framework.Gateway.MaxSessionCount == 0 || hall.SessionCount >= config.Gateway.Framework.Gateway.MaxSessionCount {
		log.Warnw("reach max session count", "session_id", sessionId, "session_count", hall.SessionCount, "max_session_count", config.Gateway.Framework.Gateway.MaxSessionCount)
		return
	}
	log.Infow("client stream open", "session_id", sessionId)
	defer func() {
		log.Infow("client stream exit", "session_id", sessionId)
	}()
	var req gira.GatewayMessage
	var err error
	var memberId string
	// 需要在一定时间内发登录协议
	timeoutCtx, timeoutFunc := context.WithTimeout(hall.ctx, time.Duration(config.Gateway.Framework.Gateway.WaitLoginTimeout)*time.Second)
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
		// 检验token
		loginReq, ok := dataReq.(gateway.LoginRequest)
		if !ok {
			log.Infow("login request cast fail", "session_id", sessionId)
			return
		}
		memberId = loginReq.GetMemberId()
		token := loginReq.GetToken()
		log.Infow("client login", "session_id", sessionId, "member_id", memberId)
		log.Infow("token", "session_id", sessionId, "token", token)
		// 验证memberid和token
		if claims, err := jwt.ParseJwtToken(token, config.Gateway.Module.Jwt.Secret); err != nil {
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
		session := newSession(hall, sessionId, memberId)
		session.serve(client, req, dataReq)
	}
}

func (hall *HallServer) OnPeerAdd(peer *gira.Peer) {
	if peer.Name != config.Gateway.Framework.Gateway.Upstream.Name {
		return
	}
	if peer.Id < config.Gateway.Framework.Gateway.Upstream.MinId || peer.Id > config.Gateway.Framework.Gateway.Upstream.MaxId {
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
		hall:       hall,
	}
	if v, loaded := hall.peers.LoadOrStore(peer.Id, server); loaded {
		lastServer := v.(*upstream_peer)
		lastServer.close()
	}
	go server.serve()
}

func (hall *HallServer) OnPeerDelete(peer *gira.Peer) {
	if peer.Name != config.Gateway.Framework.Gateway.Upstream.Name {
		return
	}
	if peer.Id < config.Gateway.Framework.Gateway.Upstream.MinId || peer.Id > config.Gateway.Framework.Gateway.Upstream.MaxId {
		return
	}
	log.Infow("remove upstream", "id", peer.Id, "fullname", peer.FullName, "grpc_addr", peer.GrpcAddr)
	if v, loaded := hall.peers.LoadAndDelete(peer.Id); loaded {
		lastServer := v.(*upstream_peer)
		lastServer.close()
	}
}

func (hall *HallServer) OnPeerUpdate(peer *gira.Peer) {
	if peer.Name != config.Gateway.Framework.Gateway.Upstream.Name {
		return
	}
	if peer.Id < config.Gateway.Framework.Gateway.Upstream.MinId || peer.Id > config.Gateway.Framework.Gateway.Upstream.MaxId {
		return
	}
	log.Infow("update upstream", "id", peer.Id, "fullname", peer.FullName, "grpc_addr", peer.GrpcAddr)
	if v, ok := hall.peers.Load(peer.Id); ok {
		lastServer := v.(*upstream_peer)
		lastServer.Address = peer.GrpcAddr
	}
}
