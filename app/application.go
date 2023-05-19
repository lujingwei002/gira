package app

import (
	"context"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/options/registry_options"
	"github.com/lujingwei002/gira/service/admin/admin_grpc"
	"google.golang.org/grpc"
)

// 负载实现gira声明的接口和启动的模块
type BaseApplication struct {
	runtime *Runtime
}

// 返回配置
func (application BaseApplication) GetConfig() *gira.Config {
	return application.runtime.config
}

// 框架启动回调
func (application *BaseApplication) OnFrameworkStart() error {
	return nil
}

// 框架销毁回调
func (application *BaseApplication) OnFrameworkStop() error {
	return nil
}

// 框架唤醒回调
func (self *BaseApplication) OnFrameworkCreate(application gira.Application) error {
	return nil
}

// 框架框架加载回调
func (application *BaseApplication) OnFrameworkConfigLoad(c *gira.Config) error {
	return nil
}

// 返回构建版本
func (application *BaseApplication) GetBuildVersion() string {
	return application.runtime.BuildVersion

}

// 返回构建时间
func (application *BaseApplication) GetBuildTime() int64 {
	return application.runtime.BuildTime
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

func (application *BaseApplication) Go(f func() error) {
	application.runtime.errGroup.Go(f)
}

func (application *BaseApplication) Done() <-chan struct{} {
	return application.runtime.ctx.Done()
}

func (application *BaseApplication) Quit() {
	application.runtime.cancelFunc()
}

func (application *BaseApplication) Context() context.Context {
	return application.runtime.ctx
}

// 设置runtime
func (application *BaseApplication) OnApplicationRun(runtime *Runtime) {
	application.runtime = runtime
}

func (application *BaseApplication) GetWorkDir() string {
	return application.runtime.WorkDir
}

func (application *BaseApplication) GetLogDir() string {
	return application.runtime.LogDir
}

func (application *BaseApplication) Wait() error {
	return application.runtime.wait()
}

func (application *BaseApplication) GetAccountDbClient() gira.DbClient {
	return application.runtime.AccountDbClient
}

func (application *BaseApplication) GetGameDbClient() gira.DbClient {
	return application.runtime.GameDbClient
}

func (application *BaseApplication) GetLogDbClient() gira.DbClient {
	return application.runtime.LogDbClient
}

func (application *BaseApplication) GetBehaviorDbClient() gira.DbClient {
	return application.runtime.BehaviorDbClient
}

func (application *BaseApplication) GetStatDbClient() gira.DbClient {
	return application.runtime.StatDbClient
}

func (application *BaseApplication) GetAccountCacheClient() gira.DbClient {
	return application.runtime.AccountCacheClient
}

func (application *BaseApplication) GetAdminCacheClient() gira.DbClient {
	return application.runtime.AdminCacheClient
}

func (application *BaseApplication) GetResourceDbClient() gira.DbClient {
	return application.runtime.ResourceDbClient
}

func (application *BaseApplication) GetAdminDbClient() gira.DbClient {
	return application.runtime.AdminDbClient
}

// implement gira.Sdk
func (application *BaseApplication) SdkLogin(accountPlat string, openId string, token string) (*gira.SdkAccount, error) {
	return application.runtime.Sdk.Login(accountPlat, openId, token)
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

func (application *BaseApplication) OnLocalPlayerAdd(player *gira.LocalPlayer) {
}

func (application *BaseApplication) OnLocalPlayerDelete(player *gira.LocalPlayer) {
}

func (application *BaseApplication) OnLocalPlayerUpdate(player *gira.LocalPlayer) {

}

func (application *BaseApplication) NewServiceName(serviceName string, opt ...registry_options.RegisterOption) string {
	return application.runtime.Registry.NewServiceName(serviceName, opt...)
}

// 注册服务
func (application *BaseApplication) RegisterService(serviceName string, opt ...registry_options.RegisterOption) (*gira.Peer, error) {
	return application.runtime.Registry.RegisterService(serviceName, opt...)
}

// 查找服务
func (application *BaseApplication) WhereIsService(serviceName string, opt ...registry_options.WhereOption) ([]*gira.Peer, error) {
	return application.runtime.Registry.WhereIsService(serviceName, opt...)
}

// 反注册服务
func (application *BaseApplication) UnregisterService(serviceName string) (*gira.Peer, error) {
	return application.runtime.Registry.UnregisterService(serviceName)
}

// 重载配置
func (application *BaseApplication) ReloadResource(dir string) error {
	if application.runtime.resourceLoader == nil {
		return gira.ErrResourceLoaderNotImplement
	}
	return application.runtime.resourceLoader.ReloadResource(dir)
}

// 加载配置
func (application *BaseApplication) LoadResource(dir string) error {
	if application.runtime.resourceLoader == nil {
		return gira.ErrResourceLoaderNotImplement
	}
	return application.runtime.resourceLoader.LoadResource(dir)
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

// 重载配置
func (application *BaseApplication) BroadcastReloadResource(ctx context.Context, name string) (result gira.BroadcastReloadResourceResult, err error) {
	req := &admin_grpc.ReloadResourceRequest{
		Name: name,
	}
	result, err = admin_grpc.DefaultAdminClients.WithBroadcast().ReloadResource(ctx, req)
	return
}

// 创建框架
func (application *BaseApplication) OnFrameworkInit() []gira.Framework {
	return nil
}

// 返回框架列表
func (application *BaseApplication) Frameworks() []gira.Framework {
	return application.runtime.Frameworks
}

// 注册grpc服务
func (application *BaseApplication) RegisterGrpc(f func(server *grpc.Server) error) error {
	if s := application.runtime.GrpcServer; s == nil {
		return gira.ErrGrpcServerNotOpen
	} else if err := s.Register(f); err != nil {
		return err
	} else {
		return nil
	}
}

func (application *BaseApplication) StopService(service gira.Service) error {
	if s := application.runtime.ServiceContainer; s == nil {
		return gira.ErrServiceContainerNotImplement.Trace()
	} else {
		return s.StopService(service)
	}
}

func (application *BaseApplication) StartService(name string, service gira.Service) error {
	if s := application.runtime.ServiceContainer; s == nil {
		return gira.ErrServiceContainerNotImplement.Trace()
	} else {
		return s.StartService(name, service)
	}
}
