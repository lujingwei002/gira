package app

import (
	"context"

	"github.com/lujingwei002/gira"
	"google.golang.org/grpc"
)

// 负载实现gira声明的接口和启动的模块
type BaseApplication struct {
	runtime *Runtime
}

func (application BaseApplication) GetConfig() *gira.Config {
	return application.runtime.config
}

func (application *BaseApplication) OnFrameworkStart() error {
	return nil
}

func (self *BaseApplication) OnFrameworkAwake(application gira.Application) error {
	return nil
}

func (application *BaseApplication) OnFrameworkConfigLoad(c *gira.Config) error {
	return nil
}

func (application *BaseApplication) GetBuildVersion() string {
	return application.runtime.BuildVersion

}

func (application *BaseApplication) GetBuildTime() int64 {
	return application.runtime.BuildTime
}

func (application *BaseApplication) GetAppId() int32 {
	return application.runtime.appId
}

func (application *BaseApplication) GetAppType() string {
	return application.runtime.appType
}

func (application *BaseApplication) GetAppName() string {
	return application.runtime.appName
}

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

func (application *BaseApplication) GetAccountDbClient() gira.MongoClient {
	return application.runtime.AccountDbClient
}

func (application *BaseApplication) GetGameDbClient() gira.MongoClient {
	return application.runtime.GameDbClient
}

func (application *BaseApplication) GetLogDbClient() gira.MongoClient {
	return application.runtime.LogDbClient
}

func (application *BaseApplication) GetStatDbClient() gira.MongoClient {
	return application.runtime.StatDbClient
}

func (application *BaseApplication) GetAccountCacheClient() gira.RedisClient {
	return application.runtime.AccountCacheClient
}

func (application *BaseApplication) GetAdminCacheClient() gira.RedisClient {
	return application.runtime.AdminCacheClient
}

func (application *BaseApplication) GetResourceDbClient() gira.MongoClient {
	return application.runtime.ResourceDbClient
}

func (application *BaseApplication) GetAdminDbClient() gira.MysqlClient {
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

func (application *BaseApplication) ReloadResource() error {
	if application.runtime.resourceLoader == nil {
		return gira.ErrResourceLoaderNotImplement
	}
	return application.runtime.resourceLoader.ReloadResource("resource")
}

func (application *BaseApplication) RangePeers(f func(k any, v any) bool) {
	application.runtime.Registry.RangePeers(f)
}

func (application *BaseApplication) BroadcastReloadResource(ctx context.Context, name string) error {
	return application.runtime.adminClient.BroadcastReloadResource(ctx, name)
}

func (application *BaseApplication) OnFrameworkInit() []gira.Framework {
	return nil
}

func (application *BaseApplication) Frameworks() []gira.Framework {
	return application.runtime.Frameworks
}

func (application *BaseApplication) RegisterGrpc(f func(server *grpc.Server) error) error {
	if s := application.runtime.GrpcServer; s == nil {
		return gira.ErrGrpcServerNotOpen
	} else if err := f(s.Server()); err != nil {
		return err
	} else {
		return nil
	}
}
