package service

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/log"
	"golang.org/x/sync/errgroup"
)

const (
	service_status_started = 1
	service_status_stopped = 2
)

type Service struct {
	status     int32
	name       string
	handler    gira.Service
	ctx        context.Context
	cancelFunc context.CancelFunc
}

type ServiceContainer struct {
	Services   sync.Map
	ctx        context.Context
	cancelFunc context.CancelFunc
	errCtx     context.Context
	errGroup   *errgroup.Group
}

func New(ctx context.Context) *ServiceContainer {
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	errGroup, errCtx := errgroup.WithContext(cancelCtx)
	return &ServiceContainer{
		ctx:        cancelCtx,
		cancelFunc: cancelFunc,
		errCtx:     errCtx,
		errGroup:   errGroup,
	}
}

func (self *ServiceContainer) Serve() error {
	<-self.ctx.Done()
	return self.errGroup.Wait()
}

// 启动服务
func (self *ServiceContainer) StartService(name string, service gira.Service) error {
	log.Debugw("start service", "name", name)
	s := &Service{
		name:    name,
		handler: service,
	}
	if _, loaded := self.Services.LoadOrStore(service, s); loaded {
		return gira.ErrServiceAlreadyStarted.Trace()
	}
	s.ctx, s.cancelFunc = context.WithCancel(self.ctx)
	if err := service.OnStart(s.ctx); err != nil {
		return err
	}
	s.status = service_status_started
	self.errGroup.Go(func() error {
		err := service.Serve()
		service.OnStop()
		return err
	})
	return nil
}

// 停止服务
func (self *ServiceContainer) StopService(service gira.Service) error {
	if v, ok := self.Services.Load(service); !ok {
		return gira.ErrServiceNotFound.Trace()
	} else {
		s := v.(*Service)
		if !atomic.CompareAndSwapInt32(&s.status, service_status_started, service_status_stopped) {
			return gira.ErrServiceAlreadyStopped.Trace()
		} else {
			log.Debugw("stop service", "name", s.name)
			s.cancelFunc()
			return nil
		}
	}
}

// 停止服务并等待
func (self *ServiceContainer) Stop() error {
	self.Services.Range(func(key, value any) bool {
		s := value.(*Service)
		s.cancelFunc()
		return true
	})
	return self.errGroup.Wait()
}
