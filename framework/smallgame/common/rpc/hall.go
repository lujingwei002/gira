package rpc

import (
	"context"
	"sync"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/framework/smallgame/gen/grpc/hall_grpc"
	"google.golang.org/grpc"
)

func init() {
	Hall = &HallClient{
		clientPool: make(map[string]*sync.Pool),
	}
}

type HallClient struct {
	clientPool map[string]*sync.Pool
	mu         sync.Mutex
}

var Hall *HallClient

func (self *HallClient) getClient(address string) (hall_grpc.HallClient, error) {
	self.mu.Lock()
	var pool *sync.Pool
	var ok bool
	if pool, ok = self.clientPool[address]; !ok {
		pool = &sync.Pool{
			New: func() any {
				conn, err := grpc.Dial(address, grpc.WithInsecure())
				if err != nil {
					return err
				}
				client := hall_grpc.NewHallClient(conn)
				return client
			},
		}
		self.clientPool[address] = pool
		self.mu.Unlock()
	} else {
		self.mu.Unlock()
	}
	if v := pool.Get(); v == nil {
		return nil, gira.ErrGrpcClientPoolNil
	} else if err, ok := v.(error); ok {
		return nil, err
	} else {
		return v.(hall_grpc.HallClient), nil
	}
}

func (self *HallClient) UserInstead(ctx context.Context, peer *gira.Peer, userId string) (*hall_grpc.UserInsteadResponse, error) {
	if client, err := self.getClient(peer.GrpcAddr); err != nil {
		return nil, err
	} else {
		req := &hall_grpc.UserInsteadRequest{
			UserId: userId,
		}
		if resp, err := client.UserInstead(ctx, req); err != nil {
			return nil, err
		} else {
			return resp, nil
		}
	}
}
