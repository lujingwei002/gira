package gateway

import (
	"context"
	"time"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/framework/smallgame/gen/grpc/hall_grpc"
	"github.com/lujingwei002/gira/log"
	"google.golang.org/grpc"
)

type upstream_peer struct {
	Id           int32
	FullName     string
	Address      string
	client       hall_grpc.HallClient
	clientStream hall_grpc.Hall_ClientStreamClient
	ctx          context.Context
	cancelFunc   context.CancelFunc
	playerCount  int64
	buildTime    int64
	buildVersion string
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
		log.Infow("server upstream exit")
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
				continue
			}
		} else {
			log.Infow("server create gate stream success", "full_name", server.FullName)
			break
		}
	}
	server.client = client
	for {
		var resp *hall_grpc.HallDataPush
		if resp, err = stream.Recv(); err != nil {
			log.Warnw("gate recv fail", "error", err)
			return err
		} else {
			log.Infow("gate recv", "resp", resp.Type)
			switch resp.Type {
			case hall_grpc.HallDataType_SERVER_REPORT:
				log.Infow("在线人数", "full_name", server.FullName, "session_count", resp.Report.PlayerCount)
				server.playerCount = resp.Report.PlayerCount
			case hall_grpc.HallDataType_SERVER_INIT:
				log.Infow("server init", "full_name", server.FullName, "build_time", resp.Init.BuildTime, "build_versioin", resp.Init.BuildVersion)
				server.buildTime = resp.Init.BuildTime
				server.buildVersion = resp.Init.BuildVersion
			}
		}
	}
}
