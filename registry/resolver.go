package registry

import (
	"fmt"
	"sync"

	"github.com/lujingwei002/gira"
	"google.golang.org/grpc/resolver"
)

// grpc resolver

const (
	PeerScheme string = "peer"
)

func formatPeerUrl(fullName string) string {
	return fmt.Sprintf("%s:///%s", PeerScheme, fullName)
}

type peer_resolver_builder struct {
	r         *Registry
	resolvers sync.Map
}

type peer_resolver struct {
	cc       resolver.ClientConn
	builder  *peer_resolver_builder
	r        *Registry
	fullName string
}

func (builder *peer_resolver_builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &peer_resolver{
		cc:       cc,
		r:        builder.r,
		builder:  builder,
		fullName: target.Endpoint,
	}

	if last, loaded := builder.resolvers.LoadOrStore(r.fullName, r); loaded {
		return last.(*peer_resolver), nil
	} else {
		if peer := r.r.peerRegistry.getPeer(r.r, r.fullName); peer == nil {
			return r, nil
		} else {
			addr := resolver.Address{
				Addr: peer.Address,
			}
			cc.UpdateState(resolver.State{
				Addresses: []resolver.Address{addr},
			})
			return r, nil
		}
	}
}

func (builder *peer_resolver_builder) onPeerUpdate(r *Registry, peer *gira.Peer) error {
	if last, ok := builder.resolvers.Load(peer.FullName); !ok {
		return nil
	} else {
		peerResolver := last.(*peer_resolver)
		return peerResolver.onPeerUpdate(r, peer)
	}
}

func (builder *peer_resolver_builder) Scheme() string {
	return PeerScheme
}

func (self *peer_resolver) ResolveNow(o resolver.ResolveNowOptions) {
}

func (self *peer_resolver) Close() {
	self.builder.resolvers.Delete(self.fullName)
}

func (self *peer_resolver) onPeerUpdate(r *Registry, peer *gira.Peer) error {
	addr := resolver.Address{
		Addr: peer.Address,
	}
	self.cc.UpdateState(resolver.State{
		Addresses: []resolver.Address{addr},
	})
	return nil
}
