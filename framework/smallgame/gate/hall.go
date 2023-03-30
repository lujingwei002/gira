package gate

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/lujingwei002/gira/log"
	"github.com/lujingwei002/gira/sproto"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/framework/smallgame/account/jwt"
	"github.com/lujingwei002/gira/framework/smallgame/gen/grpc/hall_grpc"
	"google.golang.org/grpc"
)

type Hall struct {
	proto        *sproto.Sproto
	peer         *upstream_peer
	ctx          context.Context
	cancelFunc   context.CancelFunc
	sessionCount int64 // 会话数量
	connCount    int64 // 连接数量
	handler      GateHandler
}

var hall *Hall

type upstream_peer struct {
	Id       int32
	FullName string
	Address  string
}

func init() {
	hall = &Hall{}
}

func (self *Hall) Awake(facade gira.ApplicationFacade, proto *sproto.Sproto) error {
	if handler, ok := facade.(GateHandler); !ok {
		return gira.ErrGateHandlerNotImplement
	} else {
		hall.handler = handler
	}
	hall.ctx, hall.cancelFunc = context.WithCancel(facade.Context())
	hall.proto = proto
	return nil
	// go stat()
}

func (self *Hall) SelectPeer() *upstream_peer {
	return self.peer
}

func (self *Hall) loginErrResponse(r gira.GateRequest, req sproto.SprotoRequest, err error) error {
	log.Info(err)
	resp, err := self.proto.NewResponse(req)
	if err == nil {
		return err
	}
	resp.SetErrorCode(gira.ErrCode(err))
	resp.SetErrorMsg(gira.ErrMsg(err))
	if data, err := self.proto.ResponseEncode("Login", int32(r.ReqId()), resp); err != nil {
		return err
	} else {
		r.Response(data)
	}
	return nil
}

func (self *Hall) OnGateSessionClose(s gira.GateConn) {
	atomic.AddInt64(&self.connCount, -1)
}

func (self *Hall) OnGateSessionOpen(s gira.GateConn) {
	atomic.AddInt64(&self.connCount, 1)
}

func (self *Hall) OnGateStream(client gira.GateConn) {
	sessionId := client.ID()
	log.Infow("stream open", "session_id", sessionId)
	defer func() {
		log.Infow("stream close", "session_id", sessionId)
	}()

	var req gira.GateRequest
	var err error
	var memberId string
	timeoutCtx, timeoutFunc := context.WithTimeout(self.ctx, 10*time.Second)
	defer timeoutFunc()
	req, err = client.Recv(timeoutCtx)
	if err != nil {
		log.Infow("recv fail", "session_id", sessionId, "error", err)
		return
	}
	// log.Infow("client=>gate request", "session_id", sessionId, "req_id", req.ReqId)
	if _, name, _, dataReq, err := self.proto.RequestDecode(req.Payload()); err != nil {
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
		if claims, err := jwt.ParseJwtToken(token, application.Config.Jwt.Secret); err != nil {
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
		session := &Session{
			sessionId: sessionId,
			memberId:  memberId,
		}
		session.serve(self.ctx, stream, client, req)
	}
}

func (self *Hall) OnPeerAdd(peer *gira.Peer) {
	if peer.Name != application.Config.Upstream.Name {
		return
	}
	if peer.Id != application.Config.Upstream.Id {
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

func (self *Hall) OnPeerDelete(peer *gira.Peer) {
	log.Info("OnPeerDelete")
	if peer.Name != application.Config.Upstream.Name {
		return
	}
	if peer.Id != application.Config.Upstream.Id {
		return
	}
	self.peer = nil
}
func (self *Hall) OnPeerUpdate(peer *gira.Peer) {
	log.Info("OnPeerUpdate")
	if peer.Name != application.Config.Upstream.Name {
		return
	}
	if peer.Id != application.Config.Upstream.Id {
		return
	}
	server := &upstream_peer{
		Id:       peer.Id,
		FullName: peer.FullName,
		Address:  peer.GrpcAddr,
	}
	self.peer = server
}
