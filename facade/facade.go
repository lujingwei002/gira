package facade

import (
	"context"

	"github.com/lujingwei002/gira"
)

func Context() context.Context {
	return gira.Facade().Context()
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
