package gira

import (
	"context"
)

type Framework interface {
	OnFrameworkConfigLoad(c *Config) error
	OnFrameworkCreate(application Application) error
	OnFrameworkStart() error
	OnFrameworkStop() error
}

type Application interface {
	// ======= 生命周期回调 ===========
	// 配置加载完成后接收通知
	OnConfigLoad(c *Config) error
	OnCreate() error
	OnStart() error
	OnStop() error

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

	// ======= 组件 ===========
	GetServiceComponent() ServiceComponent
	GetSdkComponent() SdkComponent
}

var application Application

func App() Application {
	return application
}

func OnApplicationCreate(app Application) {
	application = app
}

type ApplicationComponent interface {
}
