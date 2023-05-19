package module

import (
	"context"
	"sync/atomic"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/log"
	"golang.org/x/sync/errgroup"
)

const (
	module_status_started = 1
	module_status_stopped = 2
)

type Module struct {
	status     int32
	name       string
	handler    gira.Module
	err        error // first start error
	ctx        context.Context
	cancelFunc context.CancelFunc
}

type ModuleContainer struct {
	Modules    map[gira.Module]*Module
	ctx        context.Context
	cancelFunc context.CancelFunc
	errCtx     context.Context
	errGroup   *errgroup.Group
}

func New(ctx context.Context) *ModuleContainer {
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	errGroup, errCtx := errgroup.WithContext(cancelCtx)
	return &ModuleContainer{
		ctx:        cancelCtx,
		cancelFunc: cancelFunc,
		errCtx:     errCtx,
		errGroup:   errGroup,
		Modules:    make(map[gira.Module]*Module, 0),
	}
}

func (self *ModuleContainer) InstallModule(name string, m gira.Module) error {
	log.Debugw("module install", "name", name)
	cancelCtx, cancelFunc := context.WithCancel(self.ctx)
	module := &Module{
		name:       name,
		ctx:        cancelCtx,
		cancelFunc: cancelFunc,
		handler:    m,
	}
	self.Modules[m] = module
	if err := m.OnStart(cancelCtx); err != nil {
		module.err = err
		return err
	}
	module.status = module_status_started
	self.errGroup.Go(func() error {
		err := m.Serve()
		m.OnStop()
		return err
	})
	return nil
}

func (self *ModuleContainer) UninstallModule(m gira.Module) error {
	if module, ok := self.Modules[m]; !ok {
		return gira.ErrTodo.Trace()
	} else if !atomic.CompareAndSwapInt32(&module.status, module_status_started, module_status_stopped) {
		return gira.ErrTodo.Trace()
	} else {
		log.Debugw("module uninstall", "name", module.name)
		module.cancelFunc()
		return nil
	}
}

func (self *ModuleContainer) Wait() error {
	return self.errGroup.Wait()
}
