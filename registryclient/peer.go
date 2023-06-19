package registryclient

import (
	"context"
	"fmt"
	"strings"

	"github.com/lujingwei002/gira"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type peer_registry struct {
	prefix     string // /peer
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func newConfigPeerRegistry(r *RegistryClient) (*peer_registry, error) {
	ctx, cancelFunc := context.WithCancel(r.ctx)
	self := &peer_registry{
		prefix:     "/peer/",
		ctx:        ctx,
		cancelFunc: cancelFunc,
	}
	return self, nil
}

// 协程安全
func (self *peer_registry) getPeer(r *RegistryClient, fullName string) *gira.Peer {
	name, serverId, err := explodeServerFullName(fullName)
	if err != nil {
		return nil
	}
	client := r.client
	key := fmt.Sprintf("%s%s", self.prefix, fullName)
	kv := clientv3.NewKV(client)
	getResp, err := kv.Get(self.ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil
	}
	peer := &gira.Peer{
		Id:       serverId,
		Name:     name,
		FullName: r.fullName,
		Kvs:      make(map[string]string),
	}
	for _, kv := range getResp.Kvs {
		words := strings.Split(string(kv.Key), "/")
		if len(words) > 0 && words[len(words)-1] == GRPC_KEY {
			peer.GrpcAddr = string(kv.Value)
		} else {
			peer.Kvs[string(kv.Key)] = string(kv.Value)
		}
	}
	return peer
}
