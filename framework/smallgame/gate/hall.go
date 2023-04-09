package gate

import (
	"context"
	"sync"
	"time"

	"github.com/lujingwei002/gira/facade"
	"github.com/lujingwei002/gira/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/framework/smallgame/account/jwt"
	"github.com/lujingwei002/gira/framework/smallgame/gen/grpc/hall_grpc"
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

type upstream_peer struct {
	Id           int32
	FullName     string
	Address      string
	client       hall_grpc.HallClient
	clientStream hall_grpc.Hall_ClientStreamClient
	ctx          context.Context
	cancelFunc   context.CancelFunc
	playerCount  int64
	isSuspend    bool
}

func (peer *upstream_peer) close() {
	peer.cancelFunc()
}

func (server *upstream_peer) NewClientStream(ctx context.Context) (stream hall_grpc.Hall_ClientStreamClient, err error) {
	if client := server.client; client == nil {
		err = gira.ErrNullPonter
		return
	} else if stream, err = client.ClientStream(ctx); err != nil {
		return
	} else {
		return
	}
}

func (server *upstream_peer) serve() error {
	defer func() {
		log.Infow("server upstream exit")
	}()

	f := func(address string) error {
		var conn *grpc.ClientConn
		var client hall_grpc.HallClient
		var stream hall_grpc.Hall_GateStreamClient
		var err error
		// 连接hall
		conn, err = grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Errorw("server dail fail", "error", err, "full_name", server.FullName, "address", server.Address)
			return err
		}
		defer func() {
			conn.Close()
		}()
		client = hall_grpc.NewHallClient(conn)
		log.Infow("server dial success", "full_name", server.FullName)
		// 创建stream
		streamCtx, streamCancelFunc := context.WithCancel(server.ctx)
		defer func() {
			streamCancelFunc()
		}()
		stream, err = client.GateStream(streamCtx)
		if err != nil {
			log.Errorw("server create gate stream fail", "error", err, "full_name", server.FullName, "address", server.Address)
			return err
		}
		log.Infow("server create gate stream success", "full_name", server.FullName)
		server.client = client
		for {
			var resp *hall_grpc.GateDataResponse
			if resp, err = stream.Recv(); err != nil {
				log.Warnw("gate recv fail", "error", err)
				return err
			} else {
				log.Infow("gate recv", "resp", resp)
				switch resp.Type {
				case hall_grpc.GateDataType_SERVER_SUSPEND:
					server.isSuspend = true
				case hall_grpc.GateDataType_SERVER_RESUME:
					server.isSuspend = false
				}
			}
		}
	}
	errGroup, _ := errgroup.WithContext(server.ctx)
	errGroup.Go(func() error {
		for {
			if err := f(server.Address); err != nil {
			} else {
			}
			select {
			case <-server.ctx.Done():
				return server.ctx.Err()
			default:
				//expect
				time.Sleep(1 * time.Second)
			}
		}
	})
	err := errGroup.Wait()
	return err
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

func (hall *hall_server) SelectPeer() *upstream_peer {
	var minPlayerCount int64 = 0xffffffff
	var selected *upstream_peer
	hall.peers.Range(func(key any, value any) bool {
		server := value.(*upstream_peer)
		log.Println("ggggggggggggggggg", server, server.client)
		if server.client != nil && server.playerCount < minPlayerCount {
			log.Println("qqqq")
			minPlayerCount = server.playerCount
			selected = server
		}
		return true
	})
	log.Println("ffffffff", selected)
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
	// if peer.Id != hall.config.Upstream.Id {
	// 	return
	// }
	log.Infow("add upstream", "id", peer.Id, "fullname", peer.FullName, "grpc_addr", peer.GrpcAddr)
	// conn, err := grpc.Dial(peer.GrpcAddr, grpc.WithInsecure())
	// if err != nil {
	// 	log.Errorw("dial upstream fail", "error", err)
	// 	return
	// }
	ctx, cancelFUnc := context.WithCancel(hall.ctx)
	server := &upstream_peer{
		Id:         peer.Id,
		FullName:   peer.FullName,
		Address:    peer.GrpcAddr,
		ctx:        ctx,
		cancelFunc: cancelFUnc,
	}
	if v, loaded := hall.peers.LoadOrStore(peer.Id, server); loaded {
		lastServer := v.(*upstream_peer)
		lastServer.close()
	}
	go server.serve()
}

func (hall *hall_server) OnPeerDelete(peer *gira.Peer) {
	if peer.Name != hall.config.Upstream.Name {
		return
	}
	// if peer.Id != hall.config.Upstream.Id {
	// 	return
	// }
	log.Infow("remove upstream", "id", peer.Id, "fullname", peer.FullName, "grpc_addr", peer.GrpcAddr)
	if v, loaded := hall.peers.LoadAndDelete(peer.Id); loaded {
		lastServer := v.(*upstream_peer)
		lastServer.close()
	}
}

func (hall *hall_server) OnPeerUpdate(peer *gira.Peer) {
	if peer.Name != hall.config.Upstream.Name {
		return
	}
	// if peer.Id != hall.config.Upstream.Id {
	// 	return
	// }
	log.Infow("update upstream", "id", peer.Id, "fullname", peer.FullName, "grpc_addr", peer.GrpcAddr)
	if v, ok := hall.peers.Load(peer.Id); ok {
		lastServer := v.(*upstream_peer)
		lastServer.Address = peer.GrpcAddr
	}
}
