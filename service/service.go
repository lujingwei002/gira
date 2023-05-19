package service

import (
	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/facade"
)

type ServiceContainer struct {
}

func NewServiceContainer() *ServiceContainer {
	return &ServiceContainer{}
}

// 启动服务
func (c *ServiceContainer) StartService(name string, service gira.Service) error {
	if _, err := facade.RegisterServiceName(name); err != nil {
		return err
	}
	if err := service.OnStart(); err != nil {
		return err
	}
	return nil
}

// 停止服务
func (c *ServiceContainer) StopService(service gira.Service) error {
	service.OnStop()
	return nil
}
