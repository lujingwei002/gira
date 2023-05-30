package hall

import (
	"context"
	"time"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/framework/smallgame/gateway/config"
	"github.com/lujingwei002/gira/log"
	"github.com/lujingwei002/gira/service/hall/hall_grpc"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type upstream_peer struct {
	Id       int32
	FullName string
	Address  string

	client       hall_grpc.HallClient
	ctx          context.Context
	cancelFunc   context.CancelFunc
	playerCount  int64
	buildTime    int64
	buildVersion string
	hall         *HallServer
}

func (server *upstream_peer) close() {
	server.cancelFunc()
}

func (server *upstream_peer) NewClientStream(ctx context.Context) (stream hall_grpc.Hall_ClientStreamClient, err error) {
	if client := server.client; client == nil {
		err = gira.ErrNullPonter
		return
	} else {
		stream, err = client.ClientStream(ctx)
		return
	}
}

func (server *upstream_peer) serve() error {
	var conn *grpc.ClientConn
	var client hall_grpc.HallClient
	var stream hall_grpc.Hall_GateStreamClient
	var err error
	address := server.Address
	ticker := time.NewTicker(1 * time.Second)
	streamCtx, streamCancelFunc := context.WithCancel(server.ctx)
	defer func() {
		streamCancelFunc()
		ticker.Stop()
		// 关闭链接
		if conn != nil {
			conn.Close()
			conn = nil
		}
		server.client = nil
		log.Infow("server upstream exit", "full_name", server.FullName, "address", server.Address)
	}()
	for {
		// TODO: 有什么可选参数
		conn, err = grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Errorw("server dail fail", "error", err, "full_name", server.FullName, "address", server.Address)
			select {
			case <-server.ctx.Done():
				return server.ctx.Err()
			case <-ticker.C:
				// 重连
				continue
			}
		} else {
			log.Infow("server dial success", "full_name", server.FullName)
			break
		}
	}
	client = hall_grpc.NewHallClient(conn)
	for {
		stream, err = client.GateStream(streamCtx)
		if err != nil {
			log.Errorw("server create gate stream fail", "error", err, "full_name", server.FullName, "address", server.Address)
			select {
			case <-server.ctx.Done():
				return server.ctx.Err()
			case <-ticker.C:
				// 重连
				continue
			}
		} else {
			log.Infow("server create gate stream success", "full_name", server.FullName)
			break
		}
	}
	ticker.Stop()
	heartbeatTicker := time.NewTicker(time.Duration(config.Gateway.Framework.Gateway.Upstream.HeartbeatInvertal) * time.Second)
	defer func() {
		heartbeatTicker.Stop()
	}()
	// init
	{
		req := &hall_grpc.InfoRequest{}
		resp, err := client.Info(server.ctx, req)
		if err != nil {
			log.Warnw("hall info fail", "error", err)
			return err
		}
		server.buildTime = resp.BuildTime
		server.buildVersion = resp.BuildVersion
		log.Infow("server init", "full_name", server.FullName, "build_time", resp.BuildTime, "build_versioin", resp.BuildVersion)
	}
	server.client = client
	errGroup, errCtx := errgroup.WithContext(server.ctx)
	errGroup.Go(func() error {
		req := &hall_grpc.HeartbeatRequest{}
		for {
			select {
			case <-errCtx.Done():
				return errCtx.Err()
			case <-heartbeatTicker.C:
				resp, err := client.Heartbeat(server.ctx, req)
				if err != nil {
					log.Warnw("hall heartbeat fail", "error", err)
				} else {
					log.Infow("在线人数", "full_name", server.FullName, "session_count", resp.PlayerCount)
					server.playerCount = resp.PlayerCount
				}
			}
		}
	})
	errGroup.Go(func() error {
		for {
			var resp *hall_grpc.GateStreamResponse
			if resp, err = stream.Recv(); err != nil {
				log.Warnw("gate recv fail", "error", err)
				return err
			} else {
				log.Infow("gate recv", "resp", resp)
			}
		}
	})
	return errGroup.Wait()
}
