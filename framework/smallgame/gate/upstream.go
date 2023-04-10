package gate

import (
	"context"
	"time"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/framework/smallgame/gen/grpc/hall_grpc"
	"github.com/lujingwei002/gira/log"
	"golang.org/x/sync/errgroup"
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
	isSuspend    bool
	resumeChan   chan struct{}
}

func (peer *upstream_peer) close() {
	peer.cancelFunc()
}

func (server *upstream_peer) Suspend() {
	log.Infow("挂起", "full_name", server.FullName)
}

func (server *upstream_peer) Resume() {
	log.Infow("恢复", "full_name", server.FullName)
	server.resumeChan <- struct{}{}
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

	var stateFunc func(address string) error
	var stateDial func(address string) error
	var stateStream func(address string) error
	var stateSuspend func(address string) error
	server.resumeChan = make(chan struct{})
	defer func() {
		close(server.resumeChan)
		log.Infow("server upstream exit")
	}()
	// 连接状态
	stateDial = func(address string) error {
		// 重连
		ticker := time.NewTicker(1 * time.Second)
		defer func() {
			ticker.Stop()
		}()
	label_dial:
		// TODO: 有什么可选参数
		conn, err = grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Errorw("server dail fail", "error", err, "full_name", server.FullName, "address", server.Address)
			select {
			case <-server.ctx.Done():
				return server.ctx.Err()
			case <-ticker.C:
				goto label_dial
			}
		} else {
			log.Infow("server dial success", "full_name", server.FullName)
		}
		stateFunc = stateStream
		return nil
	}
	// stream状态
	stateStream = func(address string) error {
		ticker := time.NewTicker(1 * time.Second)
		// 创建stream
		streamCtx, streamCancelFunc := context.WithCancel(server.ctx)
		defer func() {
			ticker.Stop()
			streamCancelFunc()
			conn.Close()
			conn = nil
			server.client = nil
		}()
		client = hall_grpc.NewHallClient(conn)
	label_newstream:
		stream, err = client.GateStream(streamCtx)
		if err != nil {
			log.Errorw("server create gate stream fail", "error", err, "full_name", server.FullName, "address", server.Address)
			select {
			case <-server.ctx.Done():
				return server.ctx.Err()
			case <-ticker.C:
				goto label_newstream
			}
		}
		log.Infow("server create gate stream success", "full_name", server.FullName)
		server.client = client
		server.isSuspend = false
		for {
			var resp *hall_grpc.GateDataResponse
			if resp, err = stream.Recv(); err != nil {
				log.Warnw("gate recv fail", "error", err)
				if server.isSuspend {
					return nil
				} else {
					return err
				}
			} else {
				log.Infow("gate recv", "resp", resp.Type)
				switch resp.Type {
				case hall_grpc.GateDataType_SERVER_ONLINE:
					log.Infow("在线人数", "full_name", server.FullName, "session_count", resp.OnlineResponse.SessionCount)
				case hall_grpc.GateDataType_SERVER_SUSPEND:
					server.isSuspend = true
					stateFunc = stateSuspend
					return nil
				}
			}
		}
	}
	// 等待恢复状态
	stateSuspend = func(address string) error {
		// 最多等待时间
		log.Infow("等待恢复", "full_name", server.FullName)
		timer := time.NewTimer(5 * time.Second)
		defer func() {
			timer.Stop()
		}()
		select {
		case <-server.ctx.Done():
			return server.ctx.Err()
		case <-server.resumeChan:
			stateFunc = stateDial
			return nil
		case <-timer.C:
			log.Warnw("等待超时，取消挂起", "full_name", server.FullName)
			return gira.ErrTodo
		}
	}
	errGroup, _ := errgroup.WithContext(server.ctx)
	stateFunc = stateDial
	errGroup.Go(func() error {
		for {
			if err := stateFunc(server.Address); err != nil {
				return err
			}
		}
	})
	err = errGroup.Wait()
	return err
}
