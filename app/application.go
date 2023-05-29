package app

import (
	"context"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/service/admin/admin_grpc"
)

// 负载实现gira声明的接口和启动的模块
type BaseApplication struct {
	runtime *Runtime
}

// 返回配置
func (application BaseApplication) GetConfig() *gira.Config {
	return application.runtime.config
}

// 返回构建版本
func (application *BaseApplication) GetBuildVersion() string {
	return application.runtime.buildVersion

}

// 返回构建时间
func (application *BaseApplication) GetBuildTime() int64 {
	return application.runtime.buildTime
}

// 返回应用id
func (application *BaseApplication) GetAppId() int32 {
	return application.runtime.appId
}

// 返回应用类型
func (application *BaseApplication) GetAppType() string {
	return application.runtime.appType
}

// 返回应用名
func (application *BaseApplication) GetAppName() string {
	return application.runtime.appName
}

// 返回应用全名
func (application *BaseApplication) GetAppFullName() string {
	return application.runtime.appFullName
}

func (application *BaseApplication) GetWorkDir() string {
	return application.runtime.workDir
}

func (application *BaseApplication) GetLogDir() string {
	return application.runtime.logDir
}

// ================== context =========================

func (application *BaseApplication) Context() context.Context {
	return application.runtime.ctx
}

func (application *BaseApplication) Quit() {
	application.runtime.cancelFunc()
}

func (application *BaseApplication) Done() <-chan struct{} {
	return application.runtime.ctx.Done()
}

func (application *BaseApplication) Go(f func() error) {
	application.runtime.errGroup.Go(f)
}

func (application *BaseApplication) Wait() error {
	return application.runtime.Wait()
}

// ================== 回调 =========================
// 框架启动回调
func (application *BaseApplication) OnFrameworkStart() error {
	return nil
}

// 设置runtime
func (application *BaseApplication) OnApplicationRun(runtime *Runtime) {
	application.runtime = runtime
}

// 框架销毁回调
func (application *BaseApplication) OnFrameworkStop() error {
	return nil
}

// 框架唤醒回调
func (self *BaseApplication) OnFrameworkCreate(application gira.Application) error {
	return nil
}

// 创建框架
func (application *BaseApplication) OnFrameworkInit() []gira.Framework {
	return nil
}

// 框架框架加载回调
func (application *BaseApplication) OnFrameworkConfigLoad(c *gira.Config) error {
	return nil
}

func (application *BaseApplication) OnLocalPlayerAdd(player *gira.LocalPlayer) {
}

func (application *BaseApplication) OnLocalPlayerDelete(player *gira.LocalPlayer) {
}

func (application *BaseApplication) OnLocalPlayerUpdate(player *gira.LocalPlayer) {
}

// ================== implement gira.DbClientComponent ==================
func (application *BaseApplication) GetAccountDbClient() gira.DbClient {
	if application.runtime.accountDbClient == nil {
		return nil
	} else {
		return application.runtime.accountDbClient
	}
}

func (application *BaseApplication) GetGameDbClient() gira.DbClient {
	if application.runtime.gameDbClient == nil {
		return nil
	} else {
		return application.runtime.gameDbClient
	}
}

func (application *BaseApplication) GetLogDbClient() gira.DbClient {
	if application.runtime.logDbClient == nil {
		return nil
	} else {
		return application.runtime.logDbClient
	}
}

func (application *BaseApplication) GetBehaviorDbClient() gira.DbClient {
	if application.runtime.behaviorDbClient == nil {
		return nil
	} else {
		return application.runtime.behaviorDbClient
	}
}

func (application *BaseApplication) GetStatDbClient() gira.DbClient {
	if application.runtime.statDbClient == nil {
		return nil
	} else {
		return application.runtime.statDbClient
	}
}

func (application *BaseApplication) GetAccountCacheClient() gira.DbClient {
	if application.runtime.accountCacheClient == nil {
		return nil
	} else {
		return application.runtime.accountCacheClient
	}
}

func (application *BaseApplication) GetAdminCacheClient() gira.DbClient {
	if application.runtime.adminCacheClient == nil {
		return nil
	} else {
		return application.runtime.adminCacheClient
	}
}

func (application *BaseApplication) GetResourceDbClient() gira.DbClient {
	if application.runtime.resourceDbClient == nil {
		return nil
	} else {
		return application.runtime.resourceDbClient
	}
}

func (application *BaseApplication) GetAdminDbClient() gira.DbClient {
	if application.runtime.adminDbClient == nil {
		return nil
	} else {
		return application.runtime.adminDbClient
	}
}

// ================== implement gira.SdkComponent ==================
// func (application *BaseApplication) SdkLogin(accountPlat string, openId string, token string) (*gira.SdkAccount, error) {
// 	return application.runtime.SdkComponent.Login(accountPlat, openId, token)
// }

func (application *BaseApplication) GetSdkComponent() gira.SdkComponent {
	if application.runtime.sdkComponent == nil {
		return nil
	} else {
		return application.runtime.sdkComponent
	}
}

// ================== implement gira.RegistryComponent ==================
func (application *BaseApplication) GetRegistry() gira.Registry {
	if application.runtime.registry == nil {
		return nil
	} else {
		return application.runtime.registry
	}
}

// ================== implement gira.ServiceComponent ==================
func (application *BaseApplication) GetServiceComponent() gira.ServiceComponent {
	if application.runtime.serviceComponent == nil {
		return nil
	} else {
		return application.runtime.serviceComponent
	}
}

// ================== implement gira.GrpcServerComponent ==================
func (application *BaseApplication) GetGrpcServer() gira.GrpcServer {
	if application.runtime.grpcServer == nil {
		return nil
	} else {
		return application.runtime.grpcServer
	}
}

// ================== implement gira.AdminClient ==================
// 重载配置
func (application *BaseApplication) BroadcastReloadResource(ctx context.Context, name string) (result gira.BroadcastReloadResourceResult, err error) {
	req := &admin_grpc.ReloadResourceRequest{
		Name: name,
	}
	result, err = admin_grpc.DefaultAdminClients.Broadcast().ReloadResource(ctx, req)
	return
}

// 返回框架列表
func (application *BaseApplication) Frameworks() []gira.Framework {
	return application.runtime.frameworks
}

/*
func (application *BaseApplication) NewServiceName(serviceName string, opt ...service_options.RegisterOption) string {
	return application.runtime.Registry.NewServiceName(serviceName, opt...)
}

func (application *BaseApplication) LockLocalUser(userId string) (*gira.Peer, error) {
	return application.runtime.Registry.LockLocalUser(userId)
}

func (application *BaseApplication) WhereIsUser(userId string) (*gira.Peer, error) {
	return application.runtime.Registry.WhereIsUser(userId)
}

func (application *BaseApplication) UnlockLocalUser(userId string) (*gira.Peer, error) {
	return application.runtime.Registry.UnlockLocalUser(userId)
}

// 注册服务
func (application *BaseApplication) RegisterService(serviceName string, opt ...service_options.RegisterOption) (*gira.Peer, error) {
	return application.runtime.Registry.RegisterService(serviceName, opt...)
}

// 查找服务
func (application *BaseApplication) WhereIsServiceName(serviceName string, opt ...service_options.WhereOption) ([]*gira.Peer, error) {
	return application.runtime.Registry.WhereIsService(serviceName, opt...)
}

// 反注册服务
func (application *BaseApplication) UnregisterService(serviceName string) (*gira.Peer, error) {
	return application.runtime.Registry.UnregisterService(serviceName)
}

// 返回全部节点
func (application *BaseApplication) RangePeers(f func(k any, v any) bool) {
	application.runtime.Registry.RangePeers(f)
}

func (application *BaseApplication) ListPeerKvs() (peers map[string]string, err error) {
	peers, err = application.runtime.Registry.ListPeerKvs()
	return
}

func (application *BaseApplication) ListServiceKvs() (services map[string][]string, err error) {
	services, err = application.runtime.Registry.ListServiceKvs()
	return
}
*/
// 重载配置
// func (application *BaseApplication) ReloadResource(dir string) error {
// 	if application.runtime.resourceLoader == nil {
// 		return gira.ErrResourceLoaderNotImplement
// 	}
// 	return application.runtime.resourceLoader.ReloadResource(dir)
// }

// 加载配置
// func (application *BaseApplication) LoadResource(dir string) error {
// 	if application.runtime.resourceLoader == nil {
// 		return gira.ErrResourceLoaderNotImplement
// 	}
// 	return application.runtime.resourceLoader.LoadResource(dir)
// }

// func (application *BaseApplication) StopService(service gira.Service) error {
// 	if s := application.runtime.ServiceContainer; s == nil {
// 		return gira.ErrServiceContainerNotImplement.Trace()
// 	} else {
// 		return s.StopService(service)
// 	}
// }

// func (application *BaseApplication) StartService(name string, service gira.Service) error {
// 	if s := application.runtime.ServiceContainer; s == nil {
// 		return gira.ErrServiceContainerNotImplement.Trace()
// 	} else {
// 		return s.StartService(name, service)
// 	}
// }
