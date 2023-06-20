package gira

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

/*
	生命周期

OnFrameworkInit

	|

OnFrameworkConfigLoad

	|

OnConfigLoad

	|

OnFrameworkCreate

	|

OnCreate

	|

OnFrameworkStart

	|

OnStart

	|

OnStop

	|

OnFrameworkStop
*/

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
	GetWorkDir() string
	GetRespositoryVersion() string
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
	GetCron() Cron
	GetGrpcServer() GrpcServer
	GetRegistry() Registry
	GetRegistryClient() RegistryClient
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

func FormatAppFullName(appType string, appId int32, zone string, env string) string {
	return fmt.Sprintf("%s_%s_%s_%d", appType, zone, env, appId)
}

func ParseAppFullName(fullName string) (name string, id int32, err error) {
	pats := strings.Split(string(fullName), "_")
	if len(pats) != 4 {
		err = ErrInvalidPeer
		return
	}
	name = pats[0]
	var v int
	if v, err = strconv.Atoi(pats[3]); err == nil {
		id = int32(v)
	}
	return
}
