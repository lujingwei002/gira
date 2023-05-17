package rpc

import (
	"sync"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/facade"
	"github.com/lujingwei002/gira/service/hall/hall_grpc"
	"google.golang.org/grpc"
)

type HallClient struct {
	clientPool map[string]*sync.Pool
	streamPool map[string]*sync.Pool
	mu         sync.Mutex
}

var Hall *HallClient

func OnAwake() {
	Hall = &HallClient{
		clientPool: make(map[string]*sync.Pool),
	}
}

func (self *HallClient) getStream(address string) (hall_grpc.Hall_PushStreamClient, error) {
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
				if stream, err := client.PushStream(facade.Context()); err != nil {
					return err
				} else {
					return stream
				}
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
		return v.(hall_grpc.Hall_PushStreamClient), nil
	}
}

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

// func (self *HallClient) UserInstead(ctx context.Context, peer *gira.Peer, userId string) (*hall_grpc.UserInsteadResponse, error) {
// 	if client, err := self.getClient(peer.GrpcAddr); err != nil {
// 		return nil, err
// 	} else {
// 		req := &hall_grpc.UserInsteadRequest{
// 			UserId: userId,
// 		}
// 		if resp, err := client.UserInstead(ctx, req); err != nil {
// 			return nil, err
// 		} else {
// 			return resp, nil
// 		}
// 	}
// }

// func (self *HallClient) Push(ctx context.Context, peer *gira.Peer, userId string, req []byte) error {
// 	if stream, err := self.getStream(peer.GrpcAddr); err != nil {
// 		return err
// 	} else {
// 		req := &hall_grpc.PushStreamNotify{
// 			UserId: userId,
// 			Data:   req,
// 		}
// 		if err := stream.Send(req); err != nil {
// 			return err
// 		} else {
// 			return nil
// 		}
// 	}
// }

// func (self *HallClient) MustPush(ctx context.Context, peer *gira.Peer, userId string, req []byte) (*hall_grpc.MustPushResponse, error) {
// 	if client, err := self.getClient(peer.GrpcAddr); err != nil {
// 		return nil, err
// 	} else {
// 		req := &hall_grpc.MustPushRequest{
// 			UserId: userId,
// 			Data:   req,
// 		}
// 		if resp, err := client.MustPush(ctx, req); err != nil {
// 			return nil, err
// 		} else {
// 			return resp, nil
// 		}
// 	}
// }
