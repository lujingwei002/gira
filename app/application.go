package app

import (
	"context"
	"time"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/service/admin/admin_grpc"
)

// 负载实现gira声明的接口和启动的模块

// ================== implement gira.Application ==================
// 返回配置
func (application *Application) GetConfig() *gira.Config {
	return application.config
}

// 返回构建版本
func (application *Application) GetRespositoryVersion() string {
	return application.respositoryVersion

}

// 返回构建时间
func (application *Application) GetBuildTime() int64 {
	return application.buildTime
}
func (application *Application) GetUpTime() int64 {
	return time.Now().Unix() - application.buildTime
}

// 返回应用id
func (application *Application) GetAppId() int32 {
	return application.appId
}

// 返回应用类型
func (application *Application) GetAppType() string {
	return application.appType
}

// 返回应用名
func (application *Application) GetAppName() string {
	return application.appName
}

// 返回应用全名
func (application *Application) GetAppFullName() string {
	return application.appFullName
}

func (application *Application) GetWorkDir() string {
	return application.workDir
}

func (application *Application) GetLogDir() string {
	return application.logDir
}

// ================== context =========================

func (application *Application) Context() context.Context {
	return application.ctx
}

func (application *Application) Quit() {
	application.cancelFunc()
}

func (application *Application) Done() <-chan struct{} {
	return application.ctx.Done()
}

func (application *Application) Go(f func() error) {
	application.errGroup.Go(f)
}

func (application *Application) Err() error {
	return application.errCtx.Err()
}

// ================== framework =========================
// 返回框架列表
func (application *Application) Frameworks() []gira.Framework {
	return application.frameworks
}

// ================== gira.SdkComponent ==================
func (application *Application) GetSdk() gira.Sdk {
	if application.sdk == nil {
		return nil
	} else {
		return application.sdk
	}
}

// ================== gira.RegistryComponent ==================
func (application *Application) GetRegistry() gira.Registry {
	if application.registry == nil {
		return nil
	} else {
		return application.registry
	}
}

// ================== gira.ServiceContainer ==================
func (application *Application) GetServiceContainer() gira.ServiceContainer {
	if application.serviceContainer == nil {
		return nil
	} else {
		return application.serviceContainer
	}
}

// ================== gira.GrpcServerComponent ==================
func (application *Application) GetGrpcServer() gira.GrpcServer {
	if application.grpcServer == nil {
		return nil
	} else {
		return application.grpcServer
	}
}

// ================== implement gira.ResourceComponent ==================
func (application *Application) GetResourceLoader() gira.ResourceLoader {
	return application.resourceLoader
}

// ================== implement gira.AdminClient ==================
// 重载配置
func (application *Application) BroadcastReloadResource(ctx context.Context, name string) (result gira.BroadcastReloadResourceResult, err error) {
	req := &admin_grpc.ReloadResourceRequest{
		Name: name,
	}
	result, err = admin_grpc.DefaultAdminClients.Broadcast().ReloadResource(ctx, req)
	return
}

// ================== implement gira.DbClientComponent ==================
func (application *Application) GetAccountDbClient() gira.DbClient {
	if application.accountDbClient == nil {
		return nil
	} else {
		return application.accountDbClient
	}
}

func (application *Application) GetGameDbClient() gira.DbClient {
	if application.gameDbClient == nil {
		return nil
	} else {
		return application.gameDbClient
	}
}

func (application *Application) GetLogDbClient() gira.DbClient {
	if application.logDbClient == nil {
		return nil
	} else {
		return application.logDbClient
	}
}

func (application *Application) GetBehaviorDbClient() gira.DbClient {
	if application.behaviorDbClient == nil {
		return nil
	} else {
		return application.behaviorDbClient
	}
}

func (application *Application) GetStatDbClient() gira.DbClient {
	if application.statDbClient == nil {
		return nil
	} else {
		return application.statDbClient
	}
}

func (application *Application) GetAccountCacheClient() gira.DbClient {
	if application.accountCacheClient == nil {
		return nil
	} else {
		return application.accountCacheClient
	}
}

func (application *Application) GetAdminCacheClient() gira.DbClient {
	if application.adminCacheClient == nil {
		return nil
	} else {
		return application.adminCacheClient
	}
}

func (application *Application) GetResourceDbClient() gira.DbClient {
	if application.resourceDbClient == nil {
		return nil
	} else {
		return application.resourceDbClient
	}
}

func (application *Application) GetAdminDbClient() gira.DbClient {
	if application.adminDbClient == nil {
		return nil
	} else {
		return application.adminDbClient
	}
}
