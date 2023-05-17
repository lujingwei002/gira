package admin

import (
	"context"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/facade"
	"github.com/lujingwei002/gira/options/registry_options"
	"github.com/lujingwei002/gira/service/admin/admin_grpc"
	"google.golang.org/grpc"
)

type AdminService struct {
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
	return &AdminService{}
}

func (self *AdminService) OnDestory() {
}

func (self *AdminService) OnStart() error {
	facade.RegisterGrpc(func(server *grpc.Server) error {
		admin_grpc.RegisterAdminServer(server, &admin_server{})
		return nil
	})
	return nil
}

func GetServiceName() string {
	return facade.NewServiceName(admin_grpc.AdminServiceName, registry_options.WithRegisterAsGroupOption(true), registry_options.WithRegisterCatAppIdOption(true))
}
