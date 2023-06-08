package facade

import (
	"context"

	"github.com/lujingwei002/gira"
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
func GetRespositoryVersion() string {
	return gira.App().GetRespositoryVersion()
}

// 返回app构建时间
func GetBuildTime() int64 {
	return gira.App().GetBuildTime()
}

func GetResourceRespositoryVersion() string {
	application := gira.App()
	if c, ok := application.(gira.ResourceSource); !ok {
		return ""
	} else if r := c.GetResourceLoader(); r == nil {
		return ""
	} else {
		return r.GetFileRespositoryVersion()
	}
}

func GetResourceBuildTime() int64 {
	application := gira.App()
	if c, ok := application.(gira.ResourceSource); !ok {
		return 0
	} else if r := c.GetResourceLoader(); r == nil {
		return 0
	} else {
		return r.GetFileBuildTime()
	}
}

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

// ================= resource component =============================
// 重载配置
func ReloadResource() error {
	application := gira.App()
	if c, ok := application.(gira.ResourceSource); !ok {
		return gira.ErrResourceManagerNotImplement
	} else if r := c.GetResourceLoader(); r == nil {
		return gira.ErrResourceManagerNotImplement
	} else {
		return r.ReloadResource("resource")
	}
}

// 广播重载配置
func BroadcastReloadResource(ctx context.Context, name string) (result gira.BroadcastReloadResourceResult, err error) {
	application := gira.App()
	if h, ok := application.(gira.AdminClient); !ok {
		err = gira.ErrAdminClientNotImplement
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
		return nil, gira.ErrRegistryNOtImplement
	} else {
		return r.UnlockLocalUser(userId)
	}
}

// 锁定user
func LockLocalUser(userId string) (*gira.Peer, error) {
	application := gira.App()
	if r := application.GetRegistry(); r == nil {
		return nil, gira.ErrRegistryNOtImplement
	} else {
		return r.LockLocalUser(userId)
	}
}

// 查找user所在的节点
func WhereIsUser(userId string) (*gira.Peer, error) {
	application := gira.App()
	if r := application.GetRegistry(); r == nil {
		return nil, gira.ErrRegistryNOtImplement
	} else {
		return r.WhereIsUser(userId)
	}
}

func ListPeerKvs() (peers map[string]string, err error) {
	application := gira.App()
	if r := application.GetRegistry(); r == nil {
		err = gira.ErrRegistryNOtImplement
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
	if r := application.GetRegistry(); r == nil {
		err = gira.ErrRegistryNOtImplement
		return
	} else {
		peers, err = r.ListServiceKvs()
		return
	}
}

// 构造服务名
func NewServiceName(serviceName string, opt ...service_options.RegisterOption) string {
	application := gira.App()
	if r := application.GetRegistry(); r == nil {
		return ""
	} else {
		return r.NewServiceName(serviceName, opt...)
	}
}

// 注册服务名
func RegisterServiceName(serviceName string, opt ...service_options.RegisterOption) (*gira.Peer, error) {
	application := gira.App()
	if r := application.GetRegistry(); r == nil {
		return nil, gira.ErrRegistryNOtImplement
	} else {
		return r.RegisterService(serviceName, opt...)
	}
}

// 反注册服务名
func UnregisterServiceName(serviceName string) (*gira.Peer, error) {
	application := gira.App()
	if r := application.GetRegistry(); r == nil {
		return nil, gira.ErrRegistryNOtImplement
	} else {
		return r.UnregisterService(serviceName)
	}
}

// 查找服务
func WhereIsServiceName(serviceName string, opt ...service_options.WhereOption) ([]*gira.Peer, error) {
	application := gira.App()
	if r := application.GetRegistry(); r == nil {
		return nil, gira.ErrRegistryNOtImplement
	} else {
		return r.WhereIsService(serviceName, opt...)
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
		return nil, gira.ErrSdkComponentNotImplement
	} else {
		return s.Login(accountPlat, openId, token, authUrl, appId, appSecret)
	}
}

func SdkPayOrderCheck(accountPlat string, args map[string]interface{}, paySecret string) (*gira.SdkPayOrder, error) {
	application := gira.App()
	if s := application.GetSdk(); s == nil {
		return nil, gira.ErrSdkComponentNotImplement
	} else {
		return s.PayOrderCheck(accountPlat, args, paySecret)
	}
}

// ================= service =============================
// 停止服务
func StopService(service gira.Service) error {
	application := gira.App()
	if s := application.GetServiceContainer(); s == nil {
		return gira.ErrServiceContainerNotImplement
	} else {
		return s.StopService(service)
	}
}

// 启动服务
func StartService(name string, service gira.Service) error {
	application := gira.App()
	if s := application.GetServiceContainer(); s == nil {
		return gira.ErrServiceContainerNotImplement
	} else {
		return s.StartService(name, service)
	}
}
