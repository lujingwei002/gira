package gira

import (
	"context"
)

type SdkAccount struct {
	UnionId         string
	AccessToken     string
	RefreshToken    string
	TokenExpireTime int64
	Nickname        string
	Gender          int32
	SmallPic        string
	LargePic        string
}

type ApplicationFacade interface {
	// ======= 生命周期回调 ===========
	OnApplicationLoad() error
	// 配置加载完成后接收通知
	OnConfigLoad(c *Config) error
	Awake() error
	Start() error

	// ======= 状态数据 ===========
	GetAppType() string
	GetAppName() string
	GetAppFullName() string
	GetAppId() int32
	GetLogDir() string

	// ======= 同步接口 ===========
	Wait() error
	Context() context.Context
	Go(f func() error)
	Done() <-chan struct{}

	// ======= 数据库连接 ===========
	GetGameDbClient() MongoClient
	GetStatDbClient() MongoClient
	GetAccountDbClient() MongoClient
	GetResourceDbClient() MongoClient
	GetAccountCacheClient() RedisClient

	// ======= sdk登录接口 ===========
	SdkLogin(accountPlat string, open_id string, token string) (*SdkAccount, error)

	// ======= registry接口 ===========
	// 如果失败，则返回当前所在的节点
	LockLocalMember(memberId string) (*Peer, error)
	UnlockLocalMember(memberId string) (*Peer, error)
	RangePeers(f func(k any, v any) bool)

	ReloadResource() error
	BroadcastReloadResource(ctx context.Context, name string) error
}

var app ApplicationFacade

func Facade() ApplicationFacade {
	return app
}

func SetApp(v ApplicationFacade) {
	app = v
}
