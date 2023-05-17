package facade

import (
	"context"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/options/registry_options"
	"google.golang.org/grpc"
)

func GetConfig() *gira.Config {
	return gira.App().GetConfig()
}

func Context() context.Context {
	return gira.App().Context()
}

func GetBuildVersion() string {
	return gira.App().GetBuildVersion()
}

func GetBuildTime() int64 {
	return gira.App().GetBuildTime()
}

func GetAppId() int32 {
	return gira.App().GetAppId()
}

func GetAppFullName() string {
	return gira.App().GetAppFullName()
}

func GetAppType() string {
	return gira.App().GetAppType()
}

func ReloadResource() error {
	application := gira.App()
	if s, ok := application.(gira.ResourceLoader); !ok {
		return gira.ErrResourceLoaderNotImplement
	} else {
		return s.ReloadResource("resource")
	}
}

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

func GetAdminDbClient() gira.DbClient {
	application := gira.App()
	if h, ok := application.(gira.AdminDbClient); ok {
		return h.GetAdminDbClient()
	} else {
		return nil
	}
}

func UnlockLocalUser(userId string) (*gira.Peer, error) {
	application := gira.App()
	if h, ok := application.(gira.Registry); ok {
		return h.UnlockLocalUser(userId)
	} else {
		return nil, gira.ErrRegistryNOtImplement
	}
}
func LockLocalUser(userId string) (*gira.Peer, error) {
	application := gira.App()
	if h, ok := application.(gira.Registry); ok {
		return h.LockLocalUser(userId)
	} else {
		return nil, gira.ErrRegistryNOtImplement
	}
}

func WhereIsUser(userId string) (*gira.Peer, error) {
	application := gira.App()
	if h, ok := application.(gira.Registry); ok {
		return h.WhereIsUser(userId)
	} else {
		return nil, gira.ErrRegistryNOtImplement
	}
}

func RangePeers(f func(k any, v any) bool) {
	application := gira.App()
	if h, ok := application.(gira.Registry); ok {
		h.RangePeers(f)
	}
}

func NewServiceName(serviceName string, opt ...registry_options.RegisterOption) string {
	application := gira.App()
	if h, ok := application.(gira.Registry); ok {
		return h.NewServiceName(serviceName, opt...)
	} else {
		return ""
	}
}

func RegisterServiceName(serviceName string, opt ...registry_options.RegisterOption) (*gira.Peer, error) {
	application := gira.App()
	if h, ok := application.(gira.Registry); ok {
		return h.RegisterService(serviceName, opt...)
	} else {
		return nil, gira.ErrRegistryNOtImplement
	}
}

func UnregisterServiceName(serviceName string) (*gira.Peer, error) {
	application := gira.App()
	if h, ok := application.(gira.Registry); ok {
		return h.UnregisterService(serviceName)
	} else {
		return nil, gira.ErrRegistryNOtImplement
	}
}

func WhereIsService(serviceName string, opt ...registry_options.WhereOption) ([]*gira.Peer, error) {
	application := gira.App()
	if h, ok := application.(gira.Registry); ok {
		return h.WhereIsService(serviceName, opt...)
	} else {
		return nil, gira.ErrRegistryNOtImplement
	}
}

func GetResourceDbClient() gira.DbClient {
	application := gira.App()
	if h, ok := application.(gira.ResourceDbClient); ok {
		return h.GetResourceDbClient()
	} else {
		return nil
	}
}

func GetStatDbClient() gira.DbClient {
	application := gira.App()
	if h, ok := application.(gira.StatDbClient); ok {
		return h.GetStatDbClient()
	} else {
		return nil
	}
}

func GetAccountDbClient() gira.DbClient {
	application := gira.App()
	if h, ok := application.(gira.AccountDbClient); ok {
		return h.GetAccountDbClient()
	} else {
		return nil
	}
}

func GetLogDbClient() gira.DbClient {
	application := gira.App()
	if h, ok := application.(gira.LogDbClient); ok {
		return h.GetLogDbClient()
	} else {
		return nil
	}
}

func GetBehaviorDbClient() gira.DbClient {
	application := gira.App()
	if h, ok := application.(gira.BehaviorDbClient); ok {
		return h.GetBehaviorDbClient()
	} else {
		return nil
	}
}

func GetAdminCacheClient() gira.DbClient {
	application := gira.App()
	if h, ok := application.(gira.AdminCacheClient); ok {
		return h.GetAdminCacheClient()
	} else {
		return nil
	}
}

func GetAccountCacheClient() gira.DbClient {
	application := gira.App()
	if h, ok := application.(gira.AccountCacheClient); ok {
		return h.GetAccountCacheClient()
	} else {
		return nil
	}
}

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

func RegisterGrpc(f func(server *grpc.Server) error) error {
	application := gira.App()
	if s, ok := application.(gira.GrpcServer); !ok {
		return gira.ErrGrpcServerNotImplement
	} else {
		return s.RegisterGrpc(f)
	}
}

func SdkLogin(accountPlat string, openId string, token string) (*gira.SdkAccount, error) {
	application := gira.App()
	if s, ok := application.(gira.Sdk); !ok {
		return nil, gira.ErrSdkNotImplement
	} else {
		return s.SdkLogin(accountPlat, openId, token)
	}
}

func StopService(service gira.Service) error {
	application := gira.App()
	if s, ok := application.(gira.ServiceContainer); !ok {
		return gira.ErrResourceLoaderNotImplement
	} else {
		return s.StopService(service)
	}
}

func StartService(name string, service gira.Service) error {
	application := gira.App()
	if s, ok := application.(gira.ServiceContainer); !ok {
		return gira.ErrResourceLoaderNotImplement
	} else {
		return s.StartService(name, service)
	}
}
