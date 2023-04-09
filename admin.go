package gira

import "context"

type AdminService interface {
	ResourceReload() error
}

type AdminClient interface {
	// 广播所有节点重载资源
	BroadcastReloadResource(ctx context.Context, name string) error
}
