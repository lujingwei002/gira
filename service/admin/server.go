package admin

import (
	"context"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/facade"
	"github.com/lujingwei002/gira/options/service_options"
	"github.com/lujingwei002/gira/service/admin/admin_grpc"
	"google.golang.org/grpc"
)

type AdminService struct {
	ctx         context.Context
	adminServer *admin_server
}

type admin_server struct {
	admin_grpc.UnimplementedAdminServer
}

func (self *admin_server) ReloadResource(context.Context, *admin_grpc.ReloadResourceRequest) (*admin_grpc.ReloadResourceResponse, error) {
	resp := &admin_grpc.ReloadResourceResponse{}
	if err := facade.ReloadResource(); err != nil {
		return nil, err
	}
	return resp, nil
}

func NewService(facade gira.Application) *AdminService {
	return &AdminService{
		adminServer: &admin_server{},
	}
}

func (self *AdminService) OnStop() error {
	return nil
}

func (self *AdminService) Serve() error {
	<-self.ctx.Done()
	return nil
}

func (self *AdminService) OnStart(ctx context.Context) error {
	self.ctx = ctx
	facade.RegisterGrpc(func(server *grpc.Server) error {
		admin_grpc.RegisterAdminServer(server, self.adminServer)
		return nil
	})
	return nil
}

func GetServiceName() string {
	return facade.NewServiceName(admin_grpc.AdminServiceName, service_options.WithAsAppServiceOption(true))
}
