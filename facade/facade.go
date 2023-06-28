package facade

import (
	"context"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/errors"
	"github.com/lujingwei002/gira/options/service_options"
)

// 返回app配置
func GetConfig() *gira.Config {
	return gira.App().GetConfig()
}

// 返回app context
func Context() context.Context {
	return gira.App().Context()
}

// 返回app构建版本
func GetAppVersion() string {
	return gira.App().GetAppVersion()
}

// 返回app构建时间
func GetBuildTime() int64 {
	return gira.App().GetBuildTime()
}

// 返回resource版本
func GetResVersion() string {
	application := gira.App()
	if c, ok := application.(gira.ResourceSource); !ok {
		return ""
	} else if r := c.GetResourceLoader(); r == nil {
		return ""
	} else {
		return r.GetResVersion()
	}
}

// 返回resource loader版本
func GetLoaderVersion() string {
	application := gira.App()
	if c, ok := application.(gira.ResourceSource); !ok {
		return ""
	} else if r := c.GetResourceLoader(); r == nil {
		return ""
	} else {
		return r.GetLoaderVersion()
	}
}

// 返回启动时间
func GetUpTime() int64 {
	return gira.App().GetUpTime()
}

// 返回app id
func GetAppId() int32 {
	return gira.App().GetAppId()
}

// 返回app全名
func GetAppFullName() string {
	return gira.App().GetAppFullName()
}

// 返回当前所在的区
func GetZone() string {
	return gira.App().GetZone()
}

// 返回当前环境
func GetEnv() string {
	return gira.App().GetEnv()
}

// 是否开发环境
func IsDevEnv() bool {
	return gira.App().GetEnv() == gira.Env_DEV
}

// 是否测试环境
func IsQaEnv() bool {
	return gira.App().GetEnv() == gira.Env_QA
}

// 是否生产环境
func IsPrdEnv() bool {
	return gira.App().GetEnv() == gira.Env_prd
}

// 返回app类型
func GetAppType() string {
	return gira.App().GetAppType()
}

func Go(f func() error) {
	gira.App().Go(f)
}

func Done() <-chan struct{} {
	return gira.App().Done()
}

func GetLogDir() string {
	return gira.App().GetLogDir()
}

func GetWorkDir() string {
	return gira.App().GetWorkDir()
}

// ================= resource component =============================
// 重载配置
func ReloadResource() error {
	application := gira.App()
	if c, ok := application.(gira.ResourceSource); !ok {
		return errors.ErrResourceManagerNotImplement
	} else if r := c.GetResourceLoader(); r == nil {
		return errors.ErrResourceManagerNotImplement
	} else {
		return r.ReloadResource("resource")
	}
}

// 广播重载配置
func BroadcastReloadResource(ctx context.Context, name string) (result gira.BroadcastReloadResourceResult, err error) {
	application := gira.App()
	if h, ok := application.(gira.AdminClient); !ok {
		err = errors.ErrAdminClientNotImplement
		return
	} else {
		result, err = h.BroadcastReloadResource(ctx, name)
		return
	}
}

// ================= registry =============================
// 解锁user
func UnlockLocalUser(userId string) (*gira.Peer, error) {
	application := gira.App()
	if r := application.GetRegistry(); r == nil {
		return nil, errors.ErrRegistryNOtImplement
	} else {
		return r.UnlockLocalUser(userId)
	}
}

// 锁定user
func LockLocalUser(userId string) (*gira.Peer, error) {
	application := gira.App()
	if r := application.GetRegistry(); r == nil {
		return nil, errors.ErrRegistryNOtImplement
	} else {
		return r.LockLocalUser(userId)
	}
}

// 查找user所在的节点
func WhereIsUser(userId string) (*gira.Peer, error) {
	application := gira.App()
	if r := application.GetRegistry(); r != nil {
		return r.WhereIsUser(userId)
	} else if r := application.GetRegistryClient(); r != nil {
		return r.WhereIsUser(userId)
	} else {
		return nil, errors.ErrRegistryNOtImplement
	}
}

func ListPeerKvs() (peers map[string]string, err error) {
	application := gira.App()
	if r := application.GetRegistryClient(); r == nil {
		err = errors.ErrRegistryNOtImplement
		return
	} else {
		peers, err = r.ListPeerKvs()
		return
	}
}

// 遍历节点
func RangePeers(f func(k any, v any) bool) {
	application := gira.App()
	if r := application.GetRegistry(); r == nil {
		return
	} else {
		r.RangePeers(f)
	}
}

func ListServiceKvs() (peers map[string][]string, err error) {
	application := gira.App()
	if r := application.GetRegistryClient(); r == nil {
		err = errors.ErrRegistryNOtImplement
		return
	} else {
		peers, err = r.ListServiceKvs()
		return
	}
}

// 构造服务名
func NewServiceName(serviceName string, opt ...service_options.RegisterOption) string {
	application := gira.App()
	if r := application.GetRegistry(); r != nil {
		return r.NewServiceName(serviceName, opt...)
	} else if r := application.GetRegistryClient(); r != nil {
		return r.NewServiceName(serviceName, opt...)
	} else {
		return ""
	}
}

// 注册服务名
func RegisterServiceName(serviceName string, opt ...service_options.RegisterOption) (*gira.Peer, error) {
	application := gira.App()
	if r := application.GetRegistry(); r == nil {
		return nil, errors.ErrRegistryNOtImplement
	} else {
		return r.RegisterService(serviceName, opt...)
	}
}

// 反注册服务名
func UnregisterServiceName(serviceName string) (*gira.Peer, error) {
	application := gira.App()
	if r := application.GetRegistry(); r == nil {
		return nil, errors.ErrRegistryNOtImplement
	} else {
		return r.UnregisterService(serviceName)
	}
}

// 查找服务
func WhereIsServiceName(serviceName string, opt ...service_options.WhereOption) ([]*gira.Peer, error) {
	application := gira.App()
	if r := application.GetRegistry(); r != nil {
		return r.WhereIsService(serviceName, opt...)
	} else if r := application.GetRegistryClient(); r != nil {
		return r.WhereIsService(serviceName, opt...)
	} else {
		return nil, errors.ErrRegistryNOtImplement
	}
}

func UnregisterPeer(appFullName string) error {
	application := gira.App()
	if r := application.GetRegistryClient(); r != nil {
		return r.UnregisterPeer(appFullName)
	} else {
		return errors.ErrRegistryNOtImplement
	}
}

// ================= db client component =============================
// 返回预定义的admindb client
func GetAdminDbClient() gira.DbClient {
	application := gira.App()
	if c, ok := application.(gira.DbClientComponent); ok {
		return c.GetAdminDbClient()
	} else {
		return nil
	}
}

// 返回预定义的resourcedb client
func GetResourceDbClient() gira.DbClient {
	application := gira.App()
	if c, ok := application.(gira.DbClientComponent); ok {
		return c.GetResourceDbClient()
	} else {
		return nil
	}
}

// 返回预定义的statdb client
func GetStatDbClient() gira.DbClient {
	application := gira.App()
	if c, ok := application.(gira.DbClientComponent); ok {
		return c.GetStatDbClient()
	} else {
		return nil
	}
}

// 返回预定义的accountdb client
func GetAccountDbClient() gira.DbClient {
	application := gira.App()
	if c, ok := application.(gira.DbClientComponent); ok {
		return c.GetAccountDbClient()
	} else {
		return nil
	}
}

// 返回预定义的logdb client
func GetLogDbClient() gira.DbClient {
	application := gira.App()
	if c, ok := application.(gira.DbClientComponent); ok {
		return c.GetLogDbClient()
	} else {
		return nil
	}
}

// 返回预定义的behaviordb client
func GetBehaviorDbClient() gira.DbClient {
	application := gira.App()
	if c, ok := application.(gira.DbClientComponent); ok {
		return c.GetBehaviorDbClient()
	} else {
		return nil
	}
}

// 返回预定义的admincache client
func GetAdminCacheClient() gira.DbClient {
	application := gira.App()
	if c, ok := application.(gira.DbClientComponent); ok {
		return c.GetAdminCacheClient()
	} else {
		return nil
	}
}

// 返回预定义的accountcache client
func GetAccountCacheClient() gira.DbClient {
	application := gira.App()
	if c, ok := application.(gira.DbClientComponent); ok {
		return c.GetAccountCacheClient()
	} else {
		return nil
	}
}

// 返回预定义的gamedb client
func GetGameDbClient() gira.DbClient {
	application := gira.App()
	if c, ok := application.(gira.DbClientComponent); ok {
		return c.GetGameDbClient()
	} else {
		return nil
	}
}

// ================= grpc server component =============================
func GrpcServer() gira.GrpcServer {
	application := gira.App()
	return application.GetGrpcServer()
}

func IsEnableResolver() bool {
	application := gira.App()
	if cfg := application.GetConfig().Module.Grpc; cfg != nil {
		return cfg.Resolver
	} else {
		return false
	}
}

// 查看grpc server
func WhereIsServer(name string) (svr interface{}, ok bool) {
	application := gira.App()
	if s := application.GetGrpcServer(); s == nil {
		return nil, false
	} else {
		return s.GetServer(name)
	}
}

// ================= sdk =============================
// 登录sdk
func SdkLogin(accountPlat string, openId string, token string, authUrl string, appId string, appSecret string) (*gira.SdkAccount, error) {
	application := gira.App()
	if s := application.GetSdk(); s == nil {
		return nil, errors.ErrSdkComponentNotImplement
	} else {
		return s.Login(accountPlat, openId, token, authUrl, appId, appSecret)
	}
}

// 验证sdk订单
func SdkPayOrderCheck(accountPlat string, args map[string]interface{}, paySecret string) (*gira.SdkPayOrder, error) {
	application := gira.App()
	if s := application.GetSdk(); s == nil {
		return nil, errors.ErrSdkComponentNotImplement
	} else {
		return s.PayOrderCheck(accountPlat, args, paySecret)
	}
}

// ================= service =============================
// 停止服务
func StopService(service gira.Service) error {
	application := gira.App()
	if s := application.GetServiceContainer(); s == nil {
		return errors.ErrServiceContainerNotImplement
	} else {
		return s.StopService(service)
	}
}

// 启动服务
func StartService(name string, service gira.Service) error {
	application := gira.App()
	if s := application.GetServiceContainer(); s == nil {
		return errors.ErrServiceContainerNotImplement
	} else {
		return s.StartService(name, service)
	}
}

// ================= cron =============================
// 设置定时调度函数
// 并发不安全, github.com/robfig/cron 中的AddFunc函数会对同一个切片进行修改
func Cron(spec string, cmd func()) error {
	application := gira.App()
	if s := application.GetCron(); s == nil {
		return errors.ErrCronNotImplement
	} else {
		return s.AddFunc(spec, cmd)
	}
}
