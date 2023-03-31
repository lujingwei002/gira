package facade

import (
	"context"

	"github.com/lujingwei002/gira"
)

func Context() context.Context {
	return gira.Facade().Context()
}
func GetAppId() int32 {
	return gira.Facade().GetAppId()
}
func GetAppFullName() string {
	return gira.Facade().GetAppFullName()
}
func GetResourceDbClient() gira.MongoClient {
	return gira.Facade().GetResourceDbClient()
}
func ReloadResource() error {
	return gira.Facade().ReloadResource()
}
func RangePeers(f func(k any, v any) bool) {
	gira.Facade().RangePeers(f)
}

func BroadcastReloadResource(ctx context.Context, name string) error {
	return gira.Facade().BroadcastReloadResource(ctx, name)
}
func GetAdminDbClient() gira.MysqlClient {
	return gira.Facade().GetAdminDbClient()
}

func UnlockLocalUser(userId string) (*gira.Peer, error) {
	return gira.Facade().UnlockLocalUser(userId)
}
func LockLocalUser(userId string) (*gira.Peer, error) {
	return gira.Facade().LockLocalUser(userId)
}

func WhereIsUser(userId string) (*gira.Peer, error) {
	return gira.Facade().WhereIsUser(userId)
}
