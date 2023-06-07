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

type ApplicationFacade interface {
	// ======= 生命周期回调 ===========
	// 配置加载完成后接收通知
	OnConfigLoad(c *Config) error
	OnCreate() error
	OnStart() error
	OnStop() error
}

type ApplicationFramework interface {
	OnFrameworkInit() []Framework
}

type Application interface {

	// ======= 状态数据 ===========
	GetConfig() *Config
	GetAppType() string
	GetAppName() string
	GetAppFullName() string
	GetAppId() int32
	GetLogDir() string
	GetBuildVersion() string
	GetBuildTime() int64
	GetUpTime() int64

	// ======= 同步接口 ===========
	Wait() error
	Context() context.Context
	Go(f func() error)
	Done() <-chan struct{}

	Frameworks() []Framework

	// ======= 组件 ===========
	GetServiceContainer() ServiceContainer
	GetSdk() Sdk
	GetGrpcServer() GrpcServer
	GetRegistry() Registry
}

// 当前正在运行的应用
var application Application

func App() Application {
	return application
}

// 创建完成时回调
func OnApplicationCreate(app Application) {
	application = app
}
