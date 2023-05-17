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

type GrpcServer struct {
	Config gira.GrpcConfig

	server   *grpc.Server
	errGroup *errgroup.Group
	errCtx   context.Context
}

func NewConfigGrpcServer(config gira.GrpcConfig, errGroup *errgroup.Group, errCtx context.Context) *GrpcServer {
	self := &GrpcServer{
		Config:   config,
		server:   grpc.NewServer(),
		errCtx:   errCtx,
		errGroup: errGroup,
	}
	return self
}

func (self *GrpcServer) Server() *grpc.Server {
	return self.server
}

func (self *GrpcServer) Stop() error {
	return nil
}

func (self *GrpcServer) Start() error {
	listen, err := net.Listen("tcp", self.Config.Address)
	if err != nil {
		panic(err)
	}
	self.errGroup.Go(func() error {
		self.server.Serve(listen)
		log.Info("gpc server shutdown")
		return nil
	})
	self.errGroup.Go(func() error {
		select {
		case <-self.errCtx.Done():
			self.server.GracefulStop()
		}
		return nil
	})
	return nil
}
