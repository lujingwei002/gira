package facade

import (
	"context"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/options/registry_options"
	"google.golang.org/grpc"
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
func GetBuildVersion() string {
	return gira.App().GetBuildVersion()
}

// 返回app构建时间
func GetBuildTime() int64 {
	return gira.App().GetBuildTime()
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

// 重载配置
func ReloadResource() error {
	application := gira.App()
	if s, ok := application.(gira.ResourceLoader); !ok {
		return gira.ErrResourceLoaderNotImplement
	} else {
		return s.ReloadResource("resource")
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

// 解锁user
func UnlockLocalUser(userId string) (*gira.Peer, error) {
	application := gira.App()
	if h, ok := application.(gira.Registry); ok {
		return h.UnlockLocalUser(userId)
	} else {
		return nil, gira.ErrRegistryNOtImplement
	}
}

// 锁定user
func LockLocalUser(userId string) (*gira.Peer, error) {
	application := gira.App()
	if h, ok := application.(gira.Registry); ok {
		return h.LockLocalUser(userId)
	} else {
		return nil, gira.ErrRegistryNOtImplement
	}
}

// 查找user所在的节点
func WhereIsUser(userId string) (*gira.Peer, error) {
	application := gira.App()
	if h, ok := application.(gira.Registry); ok {
		return h.WhereIsUser(userId)
	} else {
		return nil, gira.ErrRegistryNOtImplement
	}
}

func ListPeerKvs() (peers map[string]string, err error) {
	application := gira.App()
	if h, ok := application.(gira.Registry); ok {
		peers, err = h.ListPeerKvs()
		return
	} else {
		err = gira.ErrRegistryNOtImplement
		return
	}
}

// 遍历节点
func RangePeers(f func(k any, v any) bool) {
	application := gira.App()
	if h, ok := application.(gira.Registry); ok {
		h.RangePeers(f)
	}
}

func ListServiceKvs() (peers map[string][]string, err error) {
	application := gira.App()
	if h, ok := application.(gira.Registry); ok {
		peers, err = h.ListServiceKvs()
		return
	} else {
		err = gira.ErrRegistryNOtImplement
		return
	}
}

// 构造服务名
func NewServiceName(serviceName string, opt ...registry_options.RegisterOption) string {
	application := gira.App()
	if h, ok := application.(gira.Registry); ok {
		return h.NewServiceName(serviceName, opt...)
	} else {
		return ""
	}
}

// 注册服务名
func RegisterServiceName(serviceName string, opt ...registry_options.RegisterOption) (*gira.Peer, error) {
	application := gira.App()
	if h, ok := application.(gira.Registry); ok {
		return h.RegisterService(serviceName, opt...)
	} else {
		return nil, gira.ErrRegistryNOtImplement
	}
}

// 反注册服务名
func UnregisterServiceName(serviceName string) (*gira.Peer, error) {
	application := gira.App()
	if h, ok := application.(gira.Registry); ok {
		return h.UnregisterService(serviceName)
	} else {
		return nil, gira.ErrRegistryNOtImplement
	}
}

// 查找服务
func WhereIsService(serviceName string, opt ...registry_options.WhereOption) ([]*gira.Peer, error) {
	application := gira.App()
	if h, ok := application.(gira.Registry); ok {
		return h.WhereIsService(serviceName, opt...)
	} else {
		return nil, gira.ErrRegistryNOtImplement
	}
}

// 返回预定义的admindb client
func GetAdminDbClient() gira.DbClient {
	application := gira.App()
	if h, ok := application.(gira.AdminDbClient); ok {
		return h.GetAdminDbClient()
	} else {
		return nil
	}
}

// 返回预定义的resourcedb client
func GetResourceDbClient() gira.DbClient {
	application := gira.App()
	if h, ok := application.(gira.ResourceDbClient); ok {
		return h.GetResourceDbClient()
	} else {
		return nil
	}
}

// 返回预定义的statdb client
func GetStatDbClient() gira.DbClient {
	application := gira.App()
	if h, ok := application.(gira.StatDbClient); ok {
		return h.GetStatDbClient()
	} else {
		return nil
	}
}

// 返回预定义的accountdb client
func GetAccountDbClient() gira.DbClient {
	application := gira.App()
	if h, ok := application.(gira.AccountDbClient); ok {
		return h.GetAccountDbClient()
	} else {
		return nil
	}
}

// 返回预定义的logdb client
func GetLogDbClient() gira.DbClient {
	application := gira.App()
	if h, ok := application.(gira.LogDbClient); ok {
		return h.GetLogDbClient()
	} else {
		return nil
	}
}

// 返回预定义的behaviordb client
func GetBehaviorDbClient() gira.DbClient {
	application := gira.App()
	if h, ok := application.(gira.BehaviorDbClient); ok {
		return h.GetBehaviorDbClient()
	} else {
		return nil
	}
}

// 返回预定义的admincache client
func GetAdminCacheClient() gira.DbClient {
	application := gira.App()
	if h, ok := application.(gira.AdminCacheClient); ok {
		return h.GetAdminCacheClient()
	} else {
		return nil
	}
}

// 返回预定义的accountcache client
func GetAccountCacheClient() gira.DbClient {
	application := gira.App()
	if h, ok := application.(gira.AccountCacheClient); ok {
		return h.GetAccountCacheClient()
	} else {
		return nil
	}
}

// 返回预定义的gamedb client
func GetGameDbClient() gira.DbClient {
	application := gira.App()
	if h, ok := application.(gira.GameDbClient); ok {
		return h.GetGameDbClient()
	} else {
		return nil
	}
}

func Go(f func() error) {
	gira.App().Go(f)
}

func Done() <-chan struct{} {
	return gira.App().Done()
}

// 注册grpc服务
func RegisterGrpc(f func(server *grpc.Server) error) error {
	application := gira.App()
	if s, ok := application.(gira.GrpcServer); !ok {
		return gira.ErrGrpcServerNotImplement
	} else {
		return s.RegisterGrpc(f)
	}
}

// 登录sdk
func SdkLogin(accountPlat string, openId string, token string) (*gira.SdkAccount, error) {
	application := gira.App()
	if s, ok := application.(gira.Sdk); !ok {
		return nil, gira.ErrSdkNotImplement
	} else {
		return s.SdkLogin(accountPlat, openId, token)
	}
}

// 停止服务
func StopService(service gira.Service) error {
	application := gira.App()
	if s, ok := application.(gira.ServiceContainer); !ok {
		return gira.ErrResourceLoaderNotImplement
	} else {
		return s.StopService(service)
	}
}

// 启动服务
func StartService(name string, service gira.Service) error {
	application := gira.App()
	if s, ok := application.(gira.ServiceContainer); !ok {
		return gira.ErrResourceLoaderNotImplement
	} else {
		return s.StartService(name, service)
	}
}
