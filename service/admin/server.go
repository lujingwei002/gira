package admin

import (
	"context"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/facade"
	"github.com/lujingwei002/gira/log"
	"github.com/lujingwei002/gira/service/admin/admin_grpc"
	"google.golang.org/grpc"
)

type AdminService struct {
	facade gira.ApplicationFacade
}

type admin_server struct {
	admin_grpc.UnimplementedAdminServer
	facade gira.ApplicationFacade
}

func (self admin_server) ReloadResource(context.Context, *admin_grpc.ReloadResourceRequest) (*admin_grpc.ReloadResourceResponse, error) {
	resp := &admin_grpc.ReloadResourceResponse{}
	if err := facade.ReloadResource(); err != nil {
		return nil, err
	}
	return resp, nil
}

func NewService(facade gira.ApplicationFacade) *AdminService {
	return &AdminService{
		facade: facade,
	}
}

func (self *AdminService) Register(server *grpc.Server) error {
	log.Infof("注册Admin服务 %p", server)
	admin_grpc.RegisterAdminServer(server, admin_server{
		facade: self.facade,
	})
	return nil

}
