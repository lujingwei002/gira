package gateway

import (
	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/framework/smallgame/gateway/config"
)

type LoginRequest interface {
	GetMemberId() string
	GetToken() string
}

// 需要实现的接口
type GatewayHandler interface {
}

type GatewayFramework interface {
	gira.Framework
	// 配置
	GetConfig() *config.GatewayConfig
	// 当前会话的数量
	SessionCount() int64
	// 当前连接的数量
	ConnectionCount() int64
}
