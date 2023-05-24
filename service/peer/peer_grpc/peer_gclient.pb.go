// Code generated by protoc-gen-go-gclient. DO NOT EDIT.
// versions:
// - protoc-gen-go-gclient v1.3.0
// - protoc             v3.12.4
// source: service/peer/peer.proto

package peer_grpc

import (
	context "context"
	fmt "fmt"
	gira "github.com/lujingwei002/gira"
	facade "github.com/lujingwei002/gira/facade"
	service_options "github.com/lujingwei002/gira/options/service_options"
	grpc "google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
	sync "sync"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

type HealthCheckResponse_MulticastResult struct {
	errors       []error
	peerCount    int
	successPeers []*gira.Peer
	errorPeers   []*gira.Peer
	responses    []*HealthCheckResponse
}

func (r *HealthCheckResponse_MulticastResult) Error() error {
	if len(r.errors) <= 0 {
		return nil
	}
	return r.errors[0]
}
func (r *HealthCheckResponse_MulticastResult) Response(index int) *HealthCheckResponse {
	if index < 0 || index >= len(r.responses) {
		return nil
	}
	return r.responses[index]
}
func (r *HealthCheckResponse_MulticastResult) SuccessPeer(index int) *gira.Peer {
	if index < 0 || index >= len(r.successPeers) {
		return nil
	}
	return r.successPeers[index]
}
func (r *HealthCheckResponse_MulticastResult) ErrorPeer(index int) *gira.Peer {
	if index < 0 || index >= len(r.errorPeers) {
		return nil
	}
	return r.errorPeers[index]
}
func (r *HealthCheckResponse_MulticastResult) PeerCount() int {
	return r.peerCount
}
func (r *HealthCheckResponse_MulticastResult) SuccessCount() int {
	return len(r.successPeers)
}
func (r *HealthCheckResponse_MulticastResult) ErrorCount() int {
	return len(r.errorPeers)
}
func (r *HealthCheckResponse_MulticastResult) Errors(index int) error {
	if index < 0 || index >= len(r.errors) {
		return nil
	}
	return r.errors[index]
}

const (
	PeerServerName = "peer_grpc.Peer"
)

// PeerClient is the client API for Peer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PeerClients interface {
	WithServiceName(serviceName string) PeerClients
	Unicast() PeerClientsUnicast
	Multicast(count int) PeerClientsMulticast
	Broadcast() PeerClientsMulticast

	HealthCheck(ctx context.Context, address string, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse, error)
}

type PeerClientsMulticast interface {
	WhereRegex(regex string) PeerClientsMulticast
	WherePrefix(prefix string) PeerClientsMulticast
	HealthCheck(ctx context.Context, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse_MulticastResult, error)
}

type PeerClientsUnicast interface {
	Where(serviceName string) PeerClientsUnicast
	WherePeer(peer *gira.Peer) PeerClientsUnicast
	WhereAddress(address string) PeerClientsUnicast
	WhereUserId(userId string) PeerClientsUnicast
	WhereUserCatalog(userId string) PeerClientsUnicast

	HealthCheck(ctx context.Context, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse, error)
}

type peerClients struct {
	mu          sync.Mutex
	clientPool  map[string]*sync.Pool
	serviceName string
}

func NewPeerClients() PeerClients {
	return &peerClients{
		serviceName: PeerServerName,
		clientPool:  make(map[string]*sync.Pool, 0),
	}
}

var DefaultPeerClients = NewPeerClients()

func (c *peerClients) getClient(address string) (PeerClient, error) {
	c.mu.Lock()
	var pool *sync.Pool
	var ok bool
	if pool, ok = c.clientPool[address]; !ok {
		pool = &sync.Pool{
			New: func() any {
				conn, err := grpc.Dial(address, grpc.WithInsecure())
				if err != nil {
					return err
				}
				client := NewPeerClient(conn)
				return client
			},
		}
		c.clientPool[address] = pool
		c.mu.Unlock()
	} else {
		c.mu.Unlock()
	}
	if v := pool.Get(); v == nil {
		return nil, gira.ErrGrpcClientPoolNil
	} else if err, ok := v.(error); ok {
		return nil, err
	} else {
		return v.(PeerClient), nil
	}
}

func (c *peerClients) putClient(address string, client PeerClient) {
	c.mu.Lock()
	var pool *sync.Pool
	var ok bool
	if pool, ok = c.clientPool[address]; ok {
		pool.Put(client)
	}
	c.mu.Unlock()
}

func (c *peerClients) WithServiceName(serviceName string) PeerClients {
	c.serviceName = serviceName
	return c
}

func (c *peerClients) Unicast() PeerClientsUnicast {
	headers := make(map[string]string)
	u := &peerClientsUnicast{
		headers: metadata.New(headers),
		client:  c,
	}
	return u
}

func (c *peerClients) Multicast(count int) PeerClientsMulticast {
	u := &peerClientsMulticast{
		count:       count,
		serviceName: fmt.Sprintf("%s/", c.serviceName),
		client:      c,
	}
	return u
}

func (c *peerClients) Broadcast() PeerClientsMulticast {
	u := &peerClientsMulticast{
		count:       -1,
		serviceName: fmt.Sprintf("%s/", c.serviceName),
		client:      c,
	}
	return u
}

func (c *peerClients) HealthCheck(ctx context.Context, address string, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse, error) {
	client, err := c.getClient(address)
	if err != nil {
		return nil, err
	}
	defer c.putClient(address, client)
	out, err := client.HealthCheck(ctx, in, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type peerClientsUnicast struct {
	peer        *gira.Peer
	serviceName string
	address     string
	userId      string
	client      *peerClients
	headers     metadata.MD
}

func (c *peerClientsUnicast) Where(serviceName string) PeerClientsUnicast {
	c.serviceName = fmt.Sprintf("%s/%s", c.client.serviceName, serviceName)
	return c
}

func (c *peerClientsUnicast) WherePeer(peer *gira.Peer) PeerClientsUnicast {
	c.peer = peer
	return c
}

func (c *peerClientsUnicast) WhereAddress(address string) PeerClientsUnicast {
	c.address = address
	return c
}

func (c *peerClientsUnicast) WhereUserId(userId string) PeerClientsUnicast {
	c.userId = userId
	return c
}

func (c *peerClientsUnicast) WhereUserCatalog(userId string) PeerClientsUnicast {
	c.userId = userId
	c.headers.Append("catalog-key", userId)
	return c
}

func (c *peerClientsUnicast) HealthCheck(ctx context.Context, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse, error) {
	var address string
	if len(c.address) > 0 {
		address = c.address
	} else if c.peer != nil {
		address = c.peer.GrpcAddr
	} else if len(c.serviceName) > 0 {
		if peers, err := facade.WhereIsServiceName(c.serviceName); err != nil {
			return nil, err
		} else if len(peers) < 1 {
			return nil, gira.ErrPeerNotFound.Trace()
		} else {
			address = peers[0].GrpcAddr
		}
	} else if len(c.userId) > 0 {
		if peer, err := facade.WhereIsUser(c.userId); err != nil {
			return nil, err
		} else {
			address = peer.GrpcAddr
		}
	}
	if len(address) <= 0 {
		return nil, gira.ErrInvalidArgs.Trace()
	}
	client, err := c.client.getClient(address)
	if err != nil {
		return nil, err
	}
	defer c.client.putClient(address, client)
	if c.headers.Len() > 0 {
		ctx = metadata.NewOutgoingContext(ctx, c.headers)
	}
	out, err := client.HealthCheck(ctx, in, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type peerClientsMulticast struct {
	count       int
	serviceName string
	regex       string
	prefix      string
	client      *peerClients
}

func (c *peerClientsMulticast) WhereRegex(regex string) PeerClientsMulticast {
	c.regex = regex
	return c
}

func (c *peerClientsMulticast) WherePrefix(prefix string) PeerClientsMulticast {
	c.prefix = prefix
	return c
}

func (c *peerClientsMulticast) HealthCheck(ctx context.Context, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse_MulticastResult, error) {
	var peers []*gira.Peer
	var whereOpts []service_options.WhereOption
	// 多播
	whereOpts = append(whereOpts, service_options.WithWhereCatalogOption())
	if c.count > 0 {
		whereOpts = append(whereOpts, service_options.WithWhereMaxCountOption(c.count))
	}
	serviceName := c.serviceName
	if len(c.regex) > 0 {
		serviceName = fmt.Sprintf("%s/%s", c.serviceName, c.regex)
		whereOpts = append(whereOpts, service_options.WithWhereRegexOption())
	}
	if len(c.prefix) > 0 {
		serviceName = fmt.Sprintf("%s/%s", c.serviceName, c.prefix)
		whereOpts = append(whereOpts, service_options.WithWherePrefixOption())
	}
	peers, err := facade.WhereIsServiceName(serviceName, whereOpts...)
	if err != nil {
		return nil, err
	}
	result := &HealthCheckResponse_MulticastResult{}
	result.peerCount = len(peers)
	for _, peer := range peers {
		client, err := c.client.getClient(peer.GrpcAddr)
		if err != nil {
			result.errors = append(result.errors, err)
			result.errorPeers = append(result.errorPeers, peer)
			continue
		}
		out, err := client.HealthCheck(ctx, in, opts...)
		if err != nil {
			result.errors = append(result.errors, err)
			result.errorPeers = append(result.errorPeers, peer)
			c.client.putClient(peer.GrpcAddr, client)
			continue
		}
		c.client.putClient(peer.GrpcAddr, client)
		result.responses = append(result.responses, out)
		result.successPeers = append(result.successPeers, peer)
	}
	return result, nil
}
