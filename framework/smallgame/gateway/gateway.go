package gateway

import (
	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/framework/smallgame/gateway/app"
	"github.com/lujingwei002/gira/framework/smallgame/gateway/config"
)

// 需要实现的接口
type GatewayHandler interface {
}

// 登录后，第一个登录消息

type Gateway interface {
	gira.Framework
	// 配置
	GetConfig() *config.GatewayConfig
	// 当前会话的数量
	SessionCount() int64
	// 当前连接的数量
	ConnectionCount() int64
}

func NewFramework(proto gira.Proto) Gateway {
	return &app.Framework{
		Proto: proto,
	}
}
