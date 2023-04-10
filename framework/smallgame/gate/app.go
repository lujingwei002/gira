package gate

import (
	"github.com/lujingwei002/gira"
)

type Framework struct {
	hall *hall_server
	// 使用的协议，当中必须包括名为Login的协议
	Proto  gira.Proto
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
func (framework *Framework) SessionCount() int64 {
	return framework.hall.sessionCount
}

// 当前连接的数量
func (framework *Framework) ConnectionCount() int64 {
	return framework.hall.connectionCount
}

func (framework *Framework) OnFrameworkAwake(application gira.Application) error {
	framework.hall = newHall(framework, framework.Proto, framework.Config)
	if err := framework.hall.OnAwake(); err != nil {
		return err
	}
	return nil
}

func (framework *Framework) OnFrameworkStart() error {
	return nil
}

func (framework *Framework) OnFrameworkConfigLoad(c *gira.Config) error {
	framework.Config = &Config{}
	return framework.Config.OnConfigLoad(c)
}

func (framework *Framework) OnClientStream(conn gira.GateConn) {
	framework.hall.OnClientStream(conn)
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
