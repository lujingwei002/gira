package gateway

// 需要实现的接口
type GatewayHandler interface {
}

// 登录后，第一个登录消息
type LoginRequest interface {
	GetMemberId() string
	GetToken() string
}
