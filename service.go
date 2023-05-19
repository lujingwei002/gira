package gira

import "context"

type Service interface {
	OnStart(ctx context.Context) error
	Serve() error
	OnStop() error
}

type ServiceContainer interface {
	StartService(name string, service Service) error
	StopService(service Service) error
}
