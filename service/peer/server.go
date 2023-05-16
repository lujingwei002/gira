package peer

import (
	"context"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/facade"
	"github.com/lujingwei002/gira/options/registry_options"
	"github.com/lujingwei002/gira/service/peer/peer_service"
	"google.golang.org/grpc"
)

type PeerService struct {
	facade gira.Application
}

type peer_server struct {
	peer_service.UnimplementedPeerServer
	facade gira.Application
}

func (self *peer_server) HealthCheck(context.Context, *peer_service.HealthCheckRequest) (*peer_service.HealthCheckResponse, error) {
	resp := &peer_service.HealthCheckResponse{}
	return resp, nil
}

func NewService(facade gira.Application) *PeerService {
	return &PeerService{
		facade: facade,
	}
}

func (self *PeerService) Register(server *grpc.Server) error {
	peer_service.RegisterPeerServer(server, &peer_server{
		facade: self.facade,
	})
	if _, err := facade.RegisterService(peer_service.PeerServiceName, registry_options.WithRegisterAsGroupOption(true), registry_options.WithRegisterCatAppIdOption(true)); err != nil {
		return err
	}
	return nil
}
