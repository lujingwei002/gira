package gira

import "context"

// 模块接口
type Module interface {
	// 模块启动
	OnStart(ctx context.Context) error
	Serve() error
	// 模块停止
	OnStop() error
}
