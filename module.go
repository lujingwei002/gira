package gira

import "context"

// 模块接口
type Module interface {
	OnStart(ctx context.Context) error
	Serve() error
	OnStop() error
}
