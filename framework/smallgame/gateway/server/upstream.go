package server

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/errors"
	"github.com/lujingwei002/gira/framework/smallgame/gateway/config"
	"github.com/lujingwei002/gira/framework/smallgame/gen/service/hall_grpc"
	"github.com/lujingwei002/gira/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// hall.ctx
//
//	   |
//		  |----ctx---
//				|
//				|--- upstream ctx
type Upstream struct {
	Id       int32
	FullName string
	Address  string

	client      hall_grpc.HallClient
	ctx         context.Context
	cancelFunc  context.CancelFunc
	playerCount int64
	buildTime   int64
	status      hall_grpc.HallStatus
	appVersion  string
}

type upstream_map struct {
	sync.Map
}

// 选择一个节点
func (peers *upstream_map) SelectPeer() *Upstream {
	var selected *Upstream
	// 最大版本
	var maxBuildTime int64 = 0
	peers.Range(func(key any, value any) bool {
		server := value.(*Upstream)
		if server.client != nil && server.buildTime > maxBuildTime {
			maxBuildTime = server.buildTime
		}
		return true
	})
	// 选择人数最小的
	var minPlayerCount int64 = 0xffffffff
	peers.Range(func(key any, value any) bool {
		server := value.(*Upstream)
		if server.client != nil && server.buildTime == maxBuildTime && server.playerCount < minPlayerCount {
			select {
			case <-server.ctx.Done():
				return true
			default:
				minPlayerCount = server.playerCount
				selected = server
			}
		}
		return true
	})
	// return selected
	if selected != nil {
		if err := selected.HealthCheck(); err != nil {
			log.Warnw("select peer fail", "error", err, "full_name", selected.FullName)
			return nil
		} else {
			return selected
		}
	} else {
		return nil
	}
}

func newUpstream(ctx context.Context, peer *gira.Peer) *Upstream {
	ctx, cancelFUnc := context.WithCancel(ctx)
	return &Upstream{
		Id:         peer.Id,
		FullName:   peer.FullName,
		Address:    peer.Address,
		ctx:        ctx,
		cancelFunc: cancelFUnc,
		status:     hall_grpc.HallStatus_UnAvailable,
	}
}

func (servers *upstream_map) OnPeerAdd(ctx context.Context, peer *gira.Peer) {
	log.Infow("add upstream", "id", peer.Id, "fullname", peer.FullName, "grpc_addr", peer.Address)
	upstream := newUpstream(ctx, peer)
	if v, loaded := servers.LoadOrStore(peer.Id, upstream); loaded {
		lastHall := v.(*Upstream)
		lastHall.close()
	}
	go upstream.serve()
}

func (servers *upstream_map) OnPeerDelete(peer *gira.Peer) {
	log.Infow("remove upstream", "id", peer.Id, "fullname", peer.FullName, "grpc_addr", peer.Address)
	if v, loaded := servers.LoadAndDelete(peer.Id); loaded {
		lastHall := v.(*Upstream)
		lastHall.close()
	}
}

func (servers *upstream_map) OnPeerUpdate(peer *gira.Peer) {
	log.Infow("update upstream", "id", peer.Id, "fullname", peer.FullName, "grpc_addr", peer.Address)
	if v, ok := servers.Load(peer.Id); ok {
		lastHall := v.(*Upstream)
		lastHall.Address = peer.Address
	}
}

func (server *Upstream) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("full_name:%s address:%s", server.FullName, server.Address))
	return sb.String()
}

func (server *Upstream) close() {
	server.cancelFunc()
}

func (server *Upstream) NewClientStream(ctx context.Context) (stream hall_grpc.Hall_ClientStreamClient, err error) {
	if client := server.client; client == nil {
		err = errors.ErrNullPointer
		return
	} else {
		stream, err = client.ClientStream(ctx)
		return
	}
}

func (server *Upstream) HealthCheck() (err error) {
	if client := server.client; client == nil {
		err = errors.ErrNullPointer
		return
	} else {
		req := &hall_grpc.HealthCheckRequest{}
		var resp *hall_grpc.HealthCheckResponse
		if resp, err = client.HealthCheck(server.ctx, req); err != nil {
			return
		} else {
			server.status = resp.Status
			if resp.Status != hall_grpc.HallStatus_OK {
				err = errors.ErrUpstreamUnavailable
				return
			} else {
				return
			}
		}
	}
}

func (server *Upstream) serve() error {
	var conn *grpc.ClientConn
	var client hall_grpc.HallClient
	var stream hall_grpc.Hall_GateStreamClient
	var err error
	address := server.Address
	dialTicker := time.NewTicker(1 * time.Second)
	streamCtx, streamCancelFunc := context.WithCancel(server.ctx)
	defer func() {
		server.client = nil
		streamCancelFunc()
		server.cancelFunc()
		dialTicker.Stop()
		log.Infow("server upstream exit", "full_name", server.FullName, "address", server.Address)
	}()
	// 1.dial
	for {
		// TODO: 有什么可选参数
		conn, err = grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Errorw("server dail fail", "error", err, "full_name", server.FullName, "address", server.Address)
			select {
			case <-server.ctx.Done():
				return server.ctx.Err()
			case <-dialTicker.C:
				// 重连
				continue
			}
		} else {
			log.Infow("server dial success", "full_name", server.FullName)
			break
		}
	}
	defer conn.Close()
	// 2.new stream
	client = hall_grpc.NewHallClient(conn)
	for {
		stream, err = client.GateStream(streamCtx)
		if err != nil {
			log.Errorw("server create gate stream fail", "error", err, "full_name", server.FullName, "address", server.Address)
			select {
			case <-server.ctx.Done():
				return server.ctx.Err()
			case <-dialTicker.C:
				// 重连
				continue
			}
		} else {
			log.Infow("server create gate stream success", "full_name", server.FullName)
			break
		}
	}
	dialTicker.Stop()
	heartbeatTicker := time.NewTicker(time.Duration(config.Gateway.Framework.Gateway.Upstream.HeartbeatInvertal) * time.Second)
	defer heartbeatTicker.Stop()
	// 3.init
	{
		req := &hall_grpc.InfoRequest{}
		resp, err := client.Info(server.ctx, req)
		if err != nil {
			log.Warnw("hall info fail", "error", err)
			return err
		}
		server.buildTime = resp.BuildTime
		server.appVersion = resp.AppVersion
		log.Infow("server init success", "full_name", server.FullName, "build_time", resp.BuildTime, "respository_version", resp.AppVersion)
	}
	server.client = client
	errGroup, errCtx := errgroup.WithContext(server.ctx)
	// heartbeat
	errGroup.Go(func() error {
		req := &hall_grpc.HealthCheckRequest{}
		for {
			select {
			case <-errCtx.Done():
				return errCtx.Err()
			case <-heartbeatTicker.C:
				resp, err := client.HealthCheck(server.ctx, req)
				if err != nil {
					log.Warnw("hall heartbeat fail", "error", err)
				} else {
					log.Infow("在线人数", "full_name", server.FullName, "session_count", resp.PlayerCount)
					server.playerCount = resp.PlayerCount
					server.status = resp.Status
				}
			}
		}
	})
	// stream
	errGroup.Go(func() error {
		for {
			var resp *hall_grpc.GateStreamResponse
			if resp, err = stream.Recv(); err != nil {
				log.Warnw("gate recv fail", "error", err)
				// server.client = nil
				// return err
				server.status = hall_grpc.HallStatus_UnAvailable
				return nil
			} else {
				log.Infow("gate recv", "resp", resp)
			}
		}
	})
	return errGroup.Wait()
}
