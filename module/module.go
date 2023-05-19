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

func New(ctx context.Context) *ModuleContainer {
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	errGroup, errCtx := errgroup.WithContext(cancelCtx)
	return &ModuleContainer{
		ctx:        cancelCtx,
		cancelFunc: cancelFunc,
		errCtx:     errCtx,
		errGroup:   errGroup,
	}
}

func (self *ModuleContainer) InstallModule(name string, module gira.Module) error {
	log.Debugw("module install", "name", name)
	wrap := &Module{
		name:    name,
		handler: module,
	}
	if _, loaded := self.Modules.LoadOrStore(module, wrap); loaded {
		return gira.ErrTodo.Trace()
	}
	wrap.ctx, wrap.cancelFunc = context.WithCancel(self.ctx)
	if err := module.OnStart(wrap.ctx); err != nil {
		wrap.err = err
		return err
	}
	wrap.status = module_status_started
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
		wrap := v.(*Module)
		if !atomic.CompareAndSwapInt32(&wrap.status, module_status_started, module_status_stopped) {
			return gira.ErrTodo.Trace()
		} else {
			log.Debugw("module uninstall", "name", wrap.name)
			wrap.cancelFunc()
			return nil
		}
	}
}

func (self *ModuleContainer) Wait() error {
	return self.errGroup.Wait()
}

func (self *ModuleContainer) Stop() error {
	self.Modules.Range(func(key, value any) bool {
		wrap := value.(*Module)
		wrap.cancelFunc()
		return true
	})
	return self.errGroup.Wait()
}
