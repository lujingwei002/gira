package grpc

import (
	"context"
	"net"

	"github.com/lujingwei002/gira/log"

	"github.com/lujingwei002/gira"
	"google.golang.org/grpc"
)

// http://wzmmmmj.com/2020/09/06/grpc-stream/

type Server struct {
	config   gira.GrpcConfig
	server   *grpc.Server
	ctx      context.Context
	listener net.Listener
	services map[string]interface{}
}

func NewConfigGrpcServer(config gira.GrpcConfig) *Server {
	self := &Server{
		config:   config,
		server:   grpc.NewServer(),
		services: make(map[string]interface{}),
	}
	return self
}

// 实现接口 grpc.ServiceRegistrar
// 将impl保存起来，转发的时候要用到
func (self *Server) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	self.server.RegisterService(desc, impl)
	self.services[desc.ServiceName] = impl
}

func (self *Server) GetService(name string) (svr interface{}, ok bool) {
	svr, ok = self.services[name]
	return
}

func (self *Server) OnStart(ctx context.Context) error {
	self.ctx = ctx
	log.Debugw("grpc server started", "addr", self.config.Address)
	if listener, err := net.Listen("tcp", self.config.Address); err != nil {
		return err
	} else {
		self.listener = listener
		return nil
	}
}

func (self *Server) Serve() error {
	go func() {
		<-self.ctx.Done()
		self.server.GracefulStop()
	}()
	err := self.server.Serve(self.listener)
	return err
}

func (self *Server) OnStop() error {
	log.Debugw("gpc server on stop")
	return nil
}
