package app

import (
	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/framework/smallgame/gateway"
	"github.com/lujingwei002/gira/framework/smallgame/gateway/config"
	"github.com/lujingwei002/gira/framework/smallgame/gateway/hall"
)

type Framework struct {
	hall *hall.HallServer
	// 使用的协议，当中必须包括名为Login的协议
	Proto  gira.Proto
	Config *config.Config
}

func NewFramework(proto gira.Proto) gateway.GatewayFramework {
	return &Framework{
		Proto: proto,
	}
}

// 当前会话的数量
func (framework *Framework) SessionCount() int64 {
	return framework.hall.SessionCount
}

// 当前连接的数量
func (framework *Framework) ConnectionCount() int64 {
	return framework.hall.ConnectionCount
}

func (framework *Framework) GetConfig() *config.GatewayConfig {
	return &framework.Config.Framework.Gateway
}

func (framework *Framework) OnFrameworkCreate(application gira.Application) error {
	framework.hall = hall.NewHall(framework.Proto)
	if err := framework.hall.OnCreate(); err != nil {
		return err
	}
	return nil
}

func (framework *Framework) OnFrameworkStart() error {
	return nil
}

func (framework *Framework) OnFrameworkStop() error {
	return nil
}

func (framework *Framework) OnFrameworkConfigLoad(c *gira.Config) error {
	framework.Config = &config.Config{}
	return framework.Config.OnConfigLoad(c)
}

func (framework *Framework) ServeClientStream(conn gira.GatewayConn) {
	framework.hall.ServeClientStream(conn)
}

func (framework *Framework) OnPeerAdd(peer *gira.Peer) {
	framework.hall.OnPeerAdd(peer)
}

func (framework *Framework) OnPeerDelete(peer *gira.Peer) {
	framework.hall.OnPeerDelete(peer)
}

func (framework *Framework) OnPeerUpdate(peer *gira.Peer) {
	framework.hall.OnPeerUpdate(peer)
}
