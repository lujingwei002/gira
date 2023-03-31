package gate

import (
	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/app"
	"github.com/lujingwei002/gira/log"
	"github.com/lujingwei002/gira/sproto"
	"google.golang.org/grpc"
)

type GateApplication struct {
	app.BaseFacade
	hall *hall
	// 使用的协议，当中必须包括名为Login的协议
	Proto  *sproto.Sproto
	Config *Config
}

// 需要实现的接口
type GateHandler interface {
}

// 登录后，第一个登录消息
type LoginRequest interface {
	GetMemberId() string
	GetToken() string
}

// 当前会话的数量
func (self *GateApplication) SessionCount() int64 {
	return self.hall.sessionCount
}

// 当前连接的数量
func (self *GateApplication) ConnectionCount() int64 {
	return self.hall.connectionCount
}

func (self *GateApplication) OnFrameworkAwake(facade gira.ApplicationFacade) error {
	self.hall = newHall(facade, self.Proto, self.Config)
	if err := self.hall.OnAwake(); err != nil {
		return err
	}
	return nil
}

func (self *GateApplication) OnFrameworkStart() error {
	log.Info("framework start")
	return nil
}
func (self *GateApplication) OnFrameworkConfigLoad(c *gira.Config) error {
	self.Config = &Config{}
	return self.Config.OnConfigLoad(c)
}

func (self *GateApplication) OnGateStream(conn gira.GateConn) {
	self.hall.OnGateStream(conn)
}

func (self *GateApplication) OnPeerAdd(peer *gira.Peer) {
	log.Info("OnPeerAdd")
	self.hall.OnPeerAdd(peer)
}

func (self *GateApplication) OnPeerDelete(peer *gira.Peer) {
	log.Info("OnPeerDelete")
	self.hall.OnPeerDelete(peer)
}

func (self *GateApplication) OnPeerUpdate(peer *gira.Peer) {
	log.Info("OnPeerUpdate")
	self.hall.OnPeerUpdate(peer)
}

func (self *GateApplication) OnGrpcServerStart(server *grpc.Server) error {
	return nil
}

func (self *GateApplication) OnFrameworkGrpcServerStart(server *grpc.Server) error {
	return nil
}
