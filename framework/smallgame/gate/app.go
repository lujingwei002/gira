package gate

import (
	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/app"
	"github.com/lujingwei002/gira/log"
	"github.com/lujingwei002/gira/sproto"
	"google.golang.org/grpc"
)

var application *GateApplication

type GateApplication struct {
	app.BaseFacade
	Proto  *sproto.Sproto
	Config Config
}

type GateHandler interface {
	OnGateLogin(req sproto.SprotoRequest) (sproto.SprotoResponse, error)
	OnGateRequest(req sproto.SprotoRequest) (sproto.SprotoResponse, error)
}

type LoginRequest interface {
	GetMemberId() string
	GetToken() string
}

func (self *GateApplication) SessionCount() int64 {
	return hall.connCount
}

func (self *GateApplication) ConnectionCount() int64 {
	return hall.sessionCount
}

func (self *GateApplication) OnFrameworkAwake(facade gira.ApplicationFacade) error {
	application = self
	return hall.Awake(facade, self.Proto)
}

func (self *GateApplication) Start() error {
	log.Info("start")
	return nil
}
func (self *GateApplication) OnFrameworkConfigLoad(c *gira.Config) error {
	return self.Config.OnConfigLoad(c)
}

func (self *GateApplication) OnGateStream(conn gira.GateConn) {
	hall.OnGateStream(conn)
}

func (self *GateApplication) OnPeerAdd(peer *gira.Peer) {
	log.Info("OnPeerAdd")
	hall.OnPeerAdd(peer)
}

func (self *GateApplication) OnPeerDelete(peer *gira.Peer) {
	log.Info("OnPeerDelete")
	hall.OnPeerDelete(peer)
}

func (self *GateApplication) OnPeerUpdate(peer *gira.Peer) {
	log.Info("OnPeerUpdate")
	hall.OnPeerUpdate(peer)
}

func (self *GateApplication) OnGrpcServerStart(server *grpc.Server) error {
	return nil
}
func (self *GateApplication) OnFrameworkGrpcServerStart(server *grpc.Server) error {
	return nil
}
