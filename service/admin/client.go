package admin

import (
	"context"
	"sync"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/facade"
	"github.com/lujingwei002/gira/log"
	"github.com/lujingwei002/gira/service/admin/admin_grpc"
	"google.golang.org/grpc"
)

type AdminClient struct {
	clientPool map[string]*sync.Pool
	mu         sync.Mutex
}

func NewAdminClient() *AdminClient {
	return &AdminClient{
		clientPool: make(map[string]*sync.Pool),
	}
}

func (self *AdminClient) getClient(address string) (admin_grpc.AdminClient, error) {
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
				client := admin_grpc.NewAdminClient(conn)
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
		return v.(admin_grpc.AdminClient), nil
	}
}

func (self *AdminClient) BroadcastReloadResource(ctx context.Context, name string) (err error) {
	facade.RangePeers(func(k any, v any) bool {
		peer := v.(*gira.Peer)
		var client admin_grpc.AdminClient
		if client, err = self.getClient(peer.GrpcAddr); err != nil {
			return true
		} else {
			req := &admin_grpc.ReloadResourceRequest{
				Name: name,
			}
			log.Infow("reload resource", "peer", peer.Name, "address", peer.GrpcAddr)
			var resp *admin_grpc.ReloadResourceResponse
			if resp, err = client.ReloadResource(ctx, req); err != nil {
				log.Infow("reload resource fail", "error", err)
				return true
			} else {
				log.Infow("reload resource success", "resp", resp)
				return true
			}
		}
	})
	return nil
}

func (self *AdminClient) ReloadResource(ctx context.Context, peer *gira.Peer, name string) (*admin_grpc.ReloadResourceResponse, error) {
	if client, err := self.getClient(peer.GrpcAddr); err != nil {
		return nil, err
	} else {
		req := &admin_grpc.ReloadResourceRequest{
			Name: name,
		}
		if resp, err := client.ReloadResource(ctx, req); err != nil {
			return nil, err
		} else {
			return resp, nil
		}
	}
}
