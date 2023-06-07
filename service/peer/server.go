package peer

import (
	"context"

	"github.com/lujingwei002/gira/facade"
	"github.com/lujingwei002/gira/options/service_options"
	"github.com/lujingwei002/gira/service/peer/peer_grpc"
)

type PeerService struct {
	ctx        context.Context
	peerServer *peer_server
}

type peer_server struct {
	peer_grpc.UnimplementedPeerServer
}

func (self *peer_server) HealthCheck(context.Context, *peer_grpc.HealthCheckRequest) (*peer_grpc.HealthCheckResponse, error) {
	resp := &peer_grpc.HealthCheckResponse{
		BuildTime:    facade.GetBuildTime(),
		BuildVersion: facade.GetBuildVersion(),
		UpTime:       facade.GetUpTime(),
	}
	return resp, nil
}

func NewService() *PeerService {
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
	if _, err := facade.RegisterServiceName(GetServiceName()); err != nil {
		return err
	}
	peer_grpc.RegisterPeerServer(facade.GrpcServer(), self.peerServer)
	return nil
}

func GetServiceName() string {
	return facade.NewServiceName(peer_grpc.PeerServerName, service_options.WithAsAppServiceOption())
}
