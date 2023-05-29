package module

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/log"
	"golang.org/x/sync/errgroup"
)

const (
	module_status_started = 1
	module_status_stopped = 2
)

func New() *ModuleContainer {
	return &ModuleContainer{}
}

type Module struct {
	status     int32
	name       string
	handler    gira.Module
	err        error
	ctx        context.Context
	cancelFunc context.CancelFunc
}

type ModuleContainer struct {
	Modules    sync.Map
	ctx        context.Context
	cancelFunc context.CancelFunc
	errCtx     context.Context
	errGroup   *errgroup.Group
}

func (self *ModuleContainer) Serve(ctx context.Context) error {
	self.ctx, self.cancelFunc = context.WithCancel(ctx)
	self.errGroup, self.errCtx = errgroup.WithContext(self.ctx)
	return self.errGroup.Wait()
}

func (self *ModuleContainer) InstallModule(name string, module gira.Module) error {
	log.Debugw("module install", "name", name)
	mod := &Module{
		name:    name,
		handler: module,
	}
	if _, loaded := self.Modules.LoadOrStore(module, mod); loaded {
		return gira.ErrTodo.Trace()
	}
	mod.ctx, mod.cancelFunc = context.WithCancel(self.ctx)
	if err := module.OnStart(mod.ctx); err != nil {
		mod.err = err
		return err
	}
	mod.status = module_status_started
	self.errGroup.Go(func() error {
		err := module.Serve()
		module.OnStop()
		return err
	})
	return nil
}

func (self *ModuleContainer) UninstallModule(module gira.Module) error {
	if v, ok := self.Modules.Load(module); !ok {
		return gira.ErrTodo.Trace()
	} else {
		mod := v.(*Module)
		if !atomic.CompareAndSwapInt32(&mod.status, module_status_started, module_status_stopped) {
			return gira.ErrTodo.Trace()
		} else {
			log.Debugw("module uninstall", "name", mod.name)
			mod.cancelFunc()
			return nil
		}
	}
}

func (self *ModuleContainer) Stop() error {
	self.Modules.Range(func(key, value any) bool {
		mod := value.(*Module)
		mod.cancelFunc()
		return true
	})
	return self.errGroup.Wait()
}
