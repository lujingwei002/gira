package grpc

import (
	"context"
	"net"
	"sync"

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
	servers  map[string]interface{}
	mu       sync.Mutex
}

func NewConfigServer(config gira.GrpcConfig) *Server {
	opts := []grpc.ServerOption{grpc.NumStreamWorkers(config.Workers)}
	self := &Server{
		config:  config,
		server:  grpc.NewServer(opts...),
		servers: make(map[string]interface{}),
	}
	return self
}

// 实现接口 grpc.ServiceRegistrar
// 将impl保存起来，转发的时候要用到
func (self *Server) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	self.server.RegisterService(desc, impl)
	self.mu.Lock()
	self.servers[desc.ServiceName] = impl
	self.mu.Unlock()
}

func (self *Server) GetServer(name string) (svr interface{}, ok bool) {
	self.mu.Lock()
	svr, ok = self.servers[name]
	self.mu.Unlock()
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
