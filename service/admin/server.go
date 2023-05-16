package admin

import (
	"context"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/facade"
	"github.com/lujingwei002/gira/options/registry_options"
	"github.com/lujingwei002/gira/service/admin/admin_service"
	"google.golang.org/grpc"
)

type AdminService struct {
	facade gira.Application
}

type admin_server struct {
	admin_service.UnimplementedAdminServer
	facade gira.Application
}

func (self *admin_server) ReloadResource(context.Context, *admin_service.ReloadResourceRequest) (*admin_service.ReloadResourceResponse, error) {
	resp := &admin_service.ReloadResourceResponse{}
	if err := facade.ReloadResource(); err != nil {
		return nil, err
	}
	return resp, nil
}

func NewService(facade gira.Application) *AdminService {
	return &AdminService{
		facade: facade,
	}
}

func (self *AdminService) Register(server *grpc.Server) error {
	admin_service.RegisterAdminServer(server, &admin_server{
		facade: self.facade,
	})
	if _, err := facade.RegisterService(admin_service.AdminServiceName, registry_options.WithRegisterAsGroupOption(true), registry_options.WithRegisterCatAppIdOption(true)); err != nil {
		return err
	}
	return nil
}
