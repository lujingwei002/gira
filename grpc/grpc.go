package grpc

import (
	"context"
	"net"

	"github.com/lujingwei002/gira/log"

	"github.com/lujingwei002/gira"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// http://wzmmmmj.com/2020/09/06/grpc-stream/

type GrpcHandler interface {
	OnGrpcServerStart(server *grpc.Server) error
	OnFrameworkGrpcServerStart(server *grpc.Server) error
}

type GrpcServer struct {
	Config   gira.GrpcConfig
	server   *grpc.Server
	errGroup *errgroup.Group
	errCtx   context.Context
	facade   gira.ApplicationFacade
}

func NewConfigGrpcServer(config gira.GrpcConfig) *GrpcServer {
	self := &GrpcServer{
		Config: config,
		server: grpc.NewServer(),
	}
	return self
}

func (self *GrpcServer) Server() *grpc.Server {
	return self.server
}

func (self *GrpcServer) Start(facade gira.ApplicationFacade, errGroup *errgroup.Group, errCtx context.Context) error {
	self.facade = facade
	self.errCtx = errCtx
	self.errGroup = errGroup
	listen, err := net.Listen("tcp", self.Config.Address)
	if err != nil {
		panic(err)
	}

	errGroup.Go(func() error {
		self.server.Serve(listen)
		log.Info("gpc server shutdown")
		return nil
	})
	errGroup.Go(func() error {
		select {
		case <-errCtx.Done():
			self.server.GracefulStop()
		}
		return nil
	})
	return nil
}
