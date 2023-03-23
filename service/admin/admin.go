package admin

import (
	"context"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/service/admin/admin_grpc"
	"google.golang.org/grpc"
)

type admin_server struct {
	admin_grpc.UnimplementedAdminServer
	facade gira.ApplicationFacade
}

func (self admin_server) ReloadResource(context.Context, *admin_grpc.ReloadResourceRequest) (*admin_grpc.ReloadResourceResponse, error) {
	return nil, nil
}

func Register(facade gira.ApplicationFacade, server *grpc.Server) error {
	admin_grpc.RegisterAdminServer(server, admin_server{
		facade: facade,
	})
	return nil
}
