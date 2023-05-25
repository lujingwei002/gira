package app

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	go_runtime "runtime"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/lujingwei002/gira/gins"
	"github.com/lujingwei002/gira/log"
	"github.com/lujingwei002/gira/module"
	"github.com/lujingwei002/gira/service"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/db"
	"github.com/lujingwei002/gira/gate"
	"github.com/lujingwei002/gira/grpc"
	"github.com/lujingwei002/gira/proj"
	"github.com/lujingwei002/gira/registry"
	"github.com/lujingwei002/gira/sdk"
	admin_service "github.com/lujingwei002/gira/service/admin"
	peer_service "github.com/lujingwei002/gira/service/peer"

	_ "net/http/pprof"

	"golang.org/x/sync/errgroup"
)

type ApplicationArgs struct {
	AppType      string /// 服务名
	AppId        int32  /// 服务id
	BuildTime    int64
	BuildVersion string
}

type RunnableApplication interface {
	OnApplicationRun(runtime *Runtime)
}

const (
	runtime_status_started = 1
	runtime_status_stopped = 2
)

// / @Component
type Runtime struct {
	gira.BaseComponent

	zone           string // 区名 wc|qq|hw|quick
	env            string // dev|local|qa|prd
	appId          int32
	appType        string /// 服务类型
	appName        string /// 服务名
	appFullName    string /// 完整的服务名 Name_Id
	cancelFunc     context.CancelFunc
	ctx            context.Context
	errCtx         context.Context
	errGroup       *errgroup.Group
	resourceLoader gira.ResourceLoader
	config         *gira.Config
	stopChan       chan struct{}
	status         int64

	BuildVersion       string
	BuildTime          int64
	ProjectFilePath    string /// 配置文件绝对路径, gira.yaml
	ConfigDir          string /// config目录
	EnvDir             string /// env目录
	ConfigFilePath     string /// 内置配置文件
	RunConfigFilePath  string /// 运行时的配置文件
	WorkDir            string /// 工作目录
	LogDir             string /// 日志目录
	RunDir             string /// 运行目录
	Application        gira.Application
	Frameworks         []gira.Framework
	MainScene          *gira.Scene
	HttpServer         *gins.HttpServer
	Registry           *registry.Registry
	DbClients          map[string]gira.DbClient
	GameDbClient       gira.DbClient
	LogDbClient        gira.DbClient
	BehaviorDbClient   gira.DbClient
	AccountDbClient    gira.DbClient
	StatDbClient       gira.DbClient
	ResourceDbClient   gira.DbClient
	AccountCacheClient gira.DbClient
	AdminCacheClient   gira.DbClient
	AdminDbClient      gira.DbClient
	SdkComponent       *sdk.SdkComponent
	Gate               *gate.Server
	GrpcServer         *grpc.Server
	ServiceComponent   *service.ServiceComponent
	ModuleContainer    *module.ModuleContainer
}

func newRuntime(args ApplicationArgs, application gira.Application) *Runtime {
	ctx, cancelFunc := context.WithCancel(context.Background())
	errGroup, errCtx := errgroup.WithContext(ctx)
	runtime := &Runtime{
		BuildVersion:     args.BuildVersion,
		BuildTime:        args.BuildTime,
		appId:            args.AppId,
		Application:      application,
		Frameworks:       make([]gira.Framework, 0),
		appType:          args.AppType,
		appName:          fmt.Sprintf("%s_%d", args.AppType, args.AppId),
		ctx:              ctx,
		cancelFunc:       cancelFunc,
		errCtx:           errCtx,
		errGroup:         errGroup,
		stopChan:         make(chan struct{}, 1),
		ServiceComponent: service.New(ctx),
		ModuleContainer:  module.New(ctx),
	}
	return runtime
}

func (runtime *Runtime) Err() error {
	return runtime.errCtx.Err()
}

func (runtime *Runtime) Go(f func() error) {
	runtime.errGroup.Go(f)
}

func (runtime *Runtime) init() error {
	// 初始化
	rand.Seed(time.Now().UnixNano())
	// 项目配置初始化
	runtime.WorkDir = proj.Config.ProjectDir
	os.Chdir(runtime.WorkDir)
	runtime.ProjectFilePath = proj.Config.ProjectConfFilePath
	if _, err := os.Stat(runtime.ProjectFilePath); err != nil {
		return err
	}
	runtime.EnvDir = proj.Config.EnvDir
	runtime.ConfigDir = proj.Config.ConfigDir
	if _, err := os.Stat(runtime.ConfigDir); err != nil {
		return err
	}
	runtime.RunDir = proj.Config.RunDir
	if _, err := os.Stat(runtime.RunDir); err != nil {
		if err := os.Mkdir(runtime.RunDir, 0755); err != nil {
			return err
		}
	}
	runtime.RunConfigFilePath = filepath.Join(runtime.RunDir, fmt.Sprintf("%s", runtime.appFullName))
	runtime.LogDir = proj.Config.LogDir

	/*
		app.ConfigFilePath = filepath.Join(app.ConfigDir, fmt.Sprintf("%sconf.yaml", app.Name))
		if _, err := os.Stat(app.ConfigFilePath); err != nil {
			return err
		}*/
	// 读应用配置文件
	if c, err := gira.LoadConfig(runtime.ConfigDir, runtime.EnvDir, runtime.appType, runtime.appId); err != nil {
		return err
	} else {
		runtime.env = c.Env
		runtime.zone = c.Zone
		runtime.appFullName = fmt.Sprintf("%s_%s_%s_%d", runtime.appType, runtime.zone, runtime.env, runtime.appId)
		runtime.config = c
	}
	// 初始化框架
	runtime.Frameworks = runtime.Application.OnFrameworkInit()
	for _, fw := range runtime.Frameworks {
		if err := fw.OnFrameworkConfigLoad(runtime.config); err != nil {
			return err
		}
	}
	if err := runtime.Application.OnConfigLoad(runtime.config); err != nil {
		return err
	}
	// 初始化日志
	if runtime.config.Log != nil {
		if err := log.ConfigLog(runtime.Application, *runtime.config.Log); err != nil {
			return err
		}
	}
	go_runtime.GOMAXPROCS(runtime.config.Thread)
	return nil
}

func (runtime *Runtime) start() (err error) {
	if r, ok := runtime.Application.(RunnableApplication); !ok {
		return gira.ErrTodo
	} else {
		r.OnApplicationRun(runtime)
	}
	// 初始化
	if err = runtime.init(); err != nil {
		return
	}
	// 设置全局对象
	gira.OnApplicationCreate(runtime.Application)
	if err = runtime.onCreate(); err != nil {
		return
	}
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
		if err != nil {
			log.Errorw("+++++++++++++++++++++++++++++")
			log.Errorw("++         启动异常       +++", "error", err)
			log.Errorw("+++++++++++++++++++++++++++++")
			runtime.onStop()
			runtime.cancelFunc()
			runtime.ModuleContainer.Wait()
			runtime.ServiceComponent.Wait()
			runtime.errGroup.Wait()
		}
	}()
	if err = runtime.onStart(); err != nil {
		return
	}
	runtime.Go(func() error {
		runtime.ServiceComponent.Wait()
		runtime.ModuleContainer.Wait()
		return nil
	})
	runtime.Go(func() error {
		quit := make(chan os.Signal, 1)
		defer close(quit)
		signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
		for {
			select {
			// 中断
			case s := <-quit:
				log.Infow("application recv signal", "signal", s)
				switch s {
				case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
					log.Infow("+++++++++++++++++++++++++++++")
					log.Infow("++    signal interrupt    +++", "signal", s)
					log.Infow("+++++++++++++++++++++++++++++")
					runtime.onStop()
					runtime.cancelFunc()
					log.Info("application interrupt end")
					return gira.ErrInterrupt
				case syscall.SIGUSR1:
				case syscall.SIGUSR2:
				default:
				}
			// 主动停止
			case <-runtime.stopChan:
				log.Info("application stop begin")
				runtime.onStop()
				runtime.cancelFunc()
				log.Info("application stop end")
				return nil
			}
		}
	})
	runtime.status = runtime_status_started
	return nil
}

// 主动关闭, 启动成功后才可以调用
func (runtime *Runtime) Stop() error {
	if !atomic.CompareAndSwapInt64(&runtime.status, runtime_status_started, runtime_status_stopped) {
		return nil
	}
	runtime.stopChan <- struct{}{}
	err := runtime.errGroup.Wait()
	return err
}

func (runtime *Runtime) onStop() {
	log.Infow("runtime on stop")
	runtime.Application.OnStop()
	for _, fw := range runtime.Frameworks {
		log.Info("framework on stop", "name")
		if err := fw.OnFrameworkStop(); err != nil {
			log.Warnw("framework on stop fail", "error", err)
		}
	}
	runtime.ServiceComponent.Stop()
	runtime.ModuleContainer.Stop()
	// runtime.ModuleContainer.UninstallModule(runtime.Registry)
	// runtime.ModuleContainer.UninstallModule(runtime.GrpcServer)
	// if runtime.HttpServer != nil {
	// 	runtime.ModuleContainer.UninstallModule(runtime.HttpServer)
	// }
}

func (runtime *Runtime) onStart() error {

	// ==== registry ================
	if runtime.Registry != nil {
		if err := runtime.ModuleContainer.InstallModule("registry", runtime.Registry); err != nil {
			return err
		}
	}

	// ==== service ================
	if runtime.GrpcServer != nil {
		// ==== admin service ================
		{
			service := admin_service.NewService(runtime.Application)
			if err := runtime.ServiceComponent.StartService("admin", service); err != nil {
				return err
			}
		}
		// ==== peer service ================
		{
			service := peer_service.NewService(runtime.Application)
			if err := runtime.ServiceComponent.StartService("peer", service); err != nil {
				return err
			}
		}
	}
	// ==== registry ================
	if runtime.Registry != nil {
		runtime.Registry.Notify()
	}

	// ==== http ================
	if runtime.HttpServer != nil {
		if err := runtime.ModuleContainer.InstallModule("http server", runtime.HttpServer); err != nil {
			return err
		}
	}
	// ==== framework start ================
	for _, fw := range runtime.Frameworks {
		if err := fw.OnFrameworkStart(); err != nil {
			return err
		}
	}
	// ==== application start ================
	if err := runtime.Application.OnStart(); err != nil {
		return err
	}
	// ==== grpc ================
	if runtime.GrpcServer != nil {
		if err := runtime.ModuleContainer.InstallModule("grpc server", runtime.GrpcServer); err != nil {
			return err
		}
	}
	return nil
}

func (runtime *Runtime) onCreate() error {
	// log.Info("application", app.FullName, "start")
	application := runtime.Application

	// ==== 模块 =============

	// ==== pprof ================
	if runtime.config.Pprof.Port != 0 {
		go func() {
			log.Infof("pprof start at http://%s:%d/debug", runtime.config.Pprof.Bind, runtime.config.Pprof.Port)
			http.ListenAndServe(fmt.Sprintf("%s:%d", runtime.config.Pprof.Bind, runtime.config.Pprof.Port), nil)
		}()
	}

	// ==== registry ================
	if runtime.config.Module.Etcd != nil {
		if r, err := registry.NewConfigRegistry(runtime.config.Module.Etcd, application); err != nil {
			return err
		} else {
			runtime.Registry = r
		}
	}
	// ==== db ================
	runtime.DbClients = make(map[string]gira.DbClient)
	for name, c := range runtime.config.Db {
		if client, err := db.NewConfigDbClient(runtime.ctx, name, *c); err != nil {
			return err
		} else {
			runtime.DbClients[name] = client
			if name == gira.GAMEDB_NAME {
				runtime.GameDbClient = client
			} else if name == gira.RESOURCEDB_NAME {
				runtime.ResourceDbClient = client
			} else if name == gira.STATDB_NAME {
				runtime.StatDbClient = client
			} else if name == gira.ACCOUNTDB_NAME {
				runtime.AccountDbClient = client
			} else if name == gira.LOGDB_NAME {
				runtime.LogDbClient = client
			} else if name == gira.BEHAVIORDB_NAME {
				runtime.BehaviorDbClient = client
			} else if name == gira.ACCOUNTCACHE_NAME {
				runtime.AccountCacheClient = client
			} else if name == gira.ADMINCACHE_NAME {
				runtime.AdminCacheClient = client
			} else if name == gira.ADMINDB_NAME {
				runtime.AdminDbClient = client
			}
		}
	}

	// ==== 加载resource ================
	if resourceComponent, ok := runtime.Application.(gira.ResourceComponent); ok {
		resourceLoader := resourceComponent.GetResourceLoader()
		if resourceLoader != nil {
			runtime.resourceLoader = resourceLoader
			if err := runtime.resourceLoader.LoadResource("resource"); err != nil {
				return err
			}
		} else {
			return gira.ErrResourceLoaderNotImplement
		}
	} else {
		return gira.ErrResourceManagerNotImplement
	}

	// ==== grpc ================
	if runtime.config.Module.Grpc != nil {
		if _, ok := application.(gira.GrpcServerComponent); !ok {
			return gira.ErrGrpcServerNotImplement
		}
		runtime.GrpcServer = grpc.NewConfigGrpcServer(*runtime.config.Module.Grpc)
	}

	// ==== sdk================
	if runtime.config.Module.Sdk != nil {
		runtime.SdkComponent = sdk.NewConfigSdk(*runtime.config.Module.Sdk)
	}

	// ==== http ================
	if runtime.config.Module.Http != nil {
		if handler, ok := application.(gira.HttpHandler); !ok {
			return gira.ErrHttpHandlerNotImplement
		} else {
			router := handler.HttpHandler()
			if httpServer, err := gins.NewConfigHttpServer(*runtime.config.Module.Http, router); err != nil {
				return err
			} else {
				runtime.HttpServer = httpServer
			}
		}
	}

	// ==== gateway ================
	if runtime.config.Module.Gateway != nil {
		var handler gira.GatewayHandler
		if h, ok := application.(gira.GatewayHandler); ok {
			handler = h
		} else {
			for _, fw := range runtime.Frameworks {
				if h, ok = fw.(gira.GatewayHandler); ok {
					handler = h
					break
				}
			}
		}
		if handler == nil {
			return gira.ErrGateHandlerNotImplement
		}
		if gate, err := gate.NewConfigServer(runtime.Application, handler, *runtime.config.Module.Gateway); err != nil {
			return err
		} else {
			runtime.Gate = gate
		}
	}

	// ==== framework create ================
	for _, fw := range runtime.Frameworks {
		if err := fw.OnFrameworkCreate(application); err != nil {
			return err
		}
	}

	// ==== application create ================
	if err := application.OnCreate(); err != nil {
		return err
	}
	return nil
}

// 等待中断
func (runtime *Runtime) Wait() error {
	err := runtime.errGroup.Wait()
	log.Infow("application down", "full_name", runtime.appFullName, "error", err)
	return err
}
