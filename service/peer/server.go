package peer

import (
	"context"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/facade"
	"github.com/lujingwei002/gira/options/service_options"
	"github.com/lujingwei002/gira/service/peer/peer_grpc"
	"google.golang.org/grpc"
)

type PeerService struct {
	ctx        context.Context
	peerServer *peer_server
}

type peer_server struct {
	peer_grpc.UnimplementedPeerServer
}

func (self *peer_server) HealthCheck(context.Context, *peer_grpc.HealthCheckRequest) (*peer_grpc.HealthCheckResponse, error) {
	resp := &peer_grpc.HealthCheckResponse{}
	return resp, nil
}

func NewService(facade gira.Application) *PeerService {
	return &PeerService{
		peerServer: &peer_server{},
	}
}

func (self *PeerService) Serve() error {
	<-self.ctx.Done()
	return nil
}

func (self *PeerService) OnStop() error {
	return nil
}

func (self *PeerService) OnStart(ctx context.Context) error {
	self.ctx = ctx
	facade.RegisterGrpc(func(server *grpc.Server) error {
		peer_grpc.RegisterPeerServer(server, self.peerServer)
		return nil
	})
	return nil
}

func GetServiceName() string {
	return facade.NewServiceName(peer_grpc.PeerServiceName, service_options.WithAsAppServiceOption(true))
}
