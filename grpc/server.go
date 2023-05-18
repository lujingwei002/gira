package grpc

import (
	"context"
	"net"

	"github.com/lujingwei002/gira/log"

	"github.com/lujingwei002/gira"
	"google.golang.org/grpc"
)

// http://wzmmmmj.com/2020/09/06/grpc-stream/

type GrpcServer struct {
	config   gira.GrpcConfig
	server   *grpc.Server
	ctx      context.Context
	listener net.Listener
}

func NewConfigGrpcServer(config gira.GrpcConfig) *GrpcServer {
	self := &GrpcServer{
		config: config,
		server: grpc.NewServer(),
	}
	return self
}

func (self *GrpcServer) Register(f func(server *grpc.Server) error) error {
	return f(self.server)
}

func (self *GrpcServer) OnStart(ctx context.Context) error {
	self.ctx = ctx
	log.Debugw("grpc server started", "addr", self.config.Address)
	if listener, err := net.Listen("tcp", self.config.Address); err != nil {
		return err
	} else {
		self.listener = listener
		return nil
	}
}

func (self *GrpcServer) Serve() error {
	go func() {
		<-self.ctx.Done()
		self.server.GracefulStop()
	}()
	err := self.server.Serve(self.listener)
	return err
}

func (self *GrpcServer) OnStop() error {
	log.Debugw("gpc server on stop")
	return nil
}
