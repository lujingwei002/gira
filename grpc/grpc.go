package grpc

import (
	"context"
	"log"
	"net"

	"github.com/lujingwei/gira"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// http://wzmmmmj.com/2020/09/06/grpc-stream/

type GrpcHandler interface {
	OnRegisterGrpcServer(server *grpc.Server) error
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
	}
	return self
}

func (self *GrpcServer) Start(facade gira.ApplicationFacade, errGroup *errgroup.Group, errCtx context.Context) error {
	self.facade = facade
	self.errCtx = errCtx
	self.errGroup = errGroup
	listen, err := net.Listen("tcp", self.Config.Address)
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	self.server = server
	if handler, ok := facade.(GrpcHandler); !ok {
		return gira.ErrGrpcHandlerNotImplement
	} else {
		if err := handler.OnRegisterGrpcServer(self.server); err != nil {
			return err
		}
	}
	errGroup.Go(func() error {
		server.Serve(listen)
		log.Println("gpc server shutdown")
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
