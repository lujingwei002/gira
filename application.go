package gira

import (
	"context"
)

type Framework interface {
	OnFrameworkConfigLoad(c *Config) error
	OnFrameworkAwake(application Application) error
	OnFrameworkStart() error
	OnFrameworkDestory() error
}

type Application interface {
	// ======= 生命周期回调 ===========
	// 配置加载完成后接收通知
	OnConfigLoad(c *Config) error
	OnAwake() error
	OnStart() error

	// ======= 状态数据 ===========
	GetConfig() *Config
	GetAppType() string
	GetAppName() string
	GetAppFullName() string
	GetAppId() int32
	GetLogDir() string
	GetBuildVersion() string
	GetBuildTime() int64

	// ======= 同步接口 ===========
	Wait() error
	Context() context.Context
	Go(f func() error)
	Done() <-chan struct{}

	OnFrameworkInit() []Framework
	Frameworks() []Framework
}

var app Application

func App() Application {
	return app
}

func OnApplicationNew(v Application) {
	app = v
}
