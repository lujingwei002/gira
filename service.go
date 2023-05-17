package gira

type Service interface {
	OnStart() error
	OnDestory()
}

type ServiceContainer interface {
	StartService(name string, service Service) error
	StopService(service Service) error
}
