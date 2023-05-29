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
	zone              string // 区名 wc|qq|hw|quick
	env               string // dev|local|qa|prd
	appId             int32
	appType           string /// 服务类型
	appName           string /// 服务名
	appFullName       string /// 完整的服务名 Name_Id
	cancelFunc        context.CancelFunc
	ctx               context.Context
	errCtx            context.Context
	errGroup          *errgroup.Group
	resourceLoader    gira.ResourceLoader
	config            *gira.Config
	stopChan          chan struct{}
	status            int64
	buildVersion      string
	buildTime         int64
	projectFilePath   string /// 配置文件绝对路径, gira.yaml
	configDir         string /// config目录
	envDir            string /// env目录
	runConfigFilePath string /// 运行时的配置文件
	workDir           string /// 工作目录
	logDir            string /// 日志目录
	runDir            string /// 运行目录
	application       gira.Application
	frameworks        []gira.Framework
	// mainScene          *gira.Scene
	httpServer         *gins.HttpServer
	registry           *registry.Registry
	dbClients          map[string]gira.DbClient
	gameDbClient       gira.DbClient
	logDbClient        gira.DbClient
	behaviorDbClient   gira.DbClient
	accountDbClient    gira.DbClient
	statDbClient       gira.DbClient
	resourceDbClient   gira.DbClient
	accountCacheClient gira.DbClient
	adminCacheClient   gira.DbClient
	adminDbClient      gira.DbClient
	sdkComponent       *sdk.SdkComponent
	gate               *gate.Server
	grpcServer         *grpc.Server
	serviceComponent   *service.ServiceComponent
	moduleContainer    *module.ModuleContainer
}

func newRuntime(args ApplicationArgs, application gira.Application) *Runtime {
	ctx, cancelFunc := context.WithCancel(context.Background())
	errGroup, errCtx := errgroup.WithContext(ctx)
	runtime := &Runtime{
		buildVersion:     args.BuildVersion,
		buildTime:        args.BuildTime,
		appId:            args.AppId,
		application:      application,
		frameworks:       make([]gira.Framework, 0),
		appType:          args.AppType,
		appName:          fmt.Sprintf("%s_%d", args.AppType, args.AppId),
		ctx:              ctx,
		cancelFunc:       cancelFunc,
		errCtx:           errCtx,
		errGroup:         errGroup,
		stopChan:         make(chan struct{}, 1),
		serviceComponent: service.New(),
		moduleContainer:  module.New(),
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
	runtime.workDir = proj.Config.ProjectDir
	if err := os.Chdir(runtime.workDir); err != nil {
		return err
	}
	runtime.projectFilePath = proj.Config.ProjectConfFilePath
	if _, err := os.Stat(runtime.projectFilePath); err != nil {
		return err
	}
	runtime.envDir = proj.Config.EnvDir
	runtime.configDir = proj.Config.ConfigDir
	if _, err := os.Stat(runtime.configDir); err != nil {
		return err
	}
	runtime.runDir = proj.Config.RunDir
	if _, err := os.Stat(runtime.runDir); err != nil {
		if err := os.Mkdir(runtime.runDir, 0755); err != nil {
			return err
		}
	}
	runtime.runConfigFilePath = filepath.Join(runtime.runDir, fmt.Sprintf("%s", runtime.appFullName))
	runtime.logDir = proj.Config.LogDir

	/*
		app.ConfigFilePath = filepath.Join(app.ConfigDir, fmt.Sprintf("%sconf.yaml", app.Name))
		if _, err := os.Stat(app.ConfigFilePath); err != nil {
			return err
		}*/
	// 读应用配置文件
	if c, err := gira.LoadConfig(runtime.configDir, runtime.envDir, runtime.appType, runtime.appId); err != nil {
		return err
	} else {
		runtime.env = c.Env
		runtime.zone = c.Zone
		runtime.appFullName = fmt.Sprintf("%s_%s_%s_%d", runtime.appType, runtime.zone, runtime.env, runtime.appId)
		runtime.config = c
	}
	// 初始化框架
	runtime.frameworks = runtime.application.OnFrameworkInit()
	// 加载配置回调
	for _, fw := range runtime.frameworks {
		if err := fw.OnFrameworkConfigLoad(runtime.config); err != nil {
			return err
		}
	}
	if err := runtime.application.OnConfigLoad(runtime.config); err != nil {
		return err
	}
	// 初始化日志
	if runtime.config.Log != nil {
		if err := log.ConfigLog(runtime.application, *runtime.config.Log); err != nil {
			return err
		}
	}
	go_runtime.GOMAXPROCS(runtime.config.Thread)
	return nil
}

func (runtime *Runtime) start() (err error) {
	if r, ok := runtime.application.(RunnableApplication); !ok {
		return gira.ErrTodo
	} else {
		r.OnApplicationRun(runtime)
	}
	// 初始化
	if err = runtime.init(); err != nil {
		return
	}
	// 设置全局对象
	gira.OnApplicationCreate(runtime.application)
	if err = runtime.onCreate(); err != nil {
		return
	}
	if err = runtime.onStart(); err != nil {
		return
	}
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
	return nil
}

func (runtime *Runtime) onStop() {
	log.Infow("runtime on stop")
	runtime.application.OnStop()
	for _, fw := range runtime.frameworks {
		log.Info("framework on stop", "name")
		if err := fw.OnFrameworkStop(); err != nil {
			log.Warnw("framework on stop fail", "error", err)
		}
	}
	runtime.serviceComponent.Stop()
	runtime.moduleContainer.Stop()
}

func (runtime *Runtime) onStart() (err error) {
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
			runtime.errGroup.Wait()
		}
	}()
	runtime.errGroup.Go(func() error {
		return runtime.moduleContainer.Serve(runtime.errCtx)
	})
	runtime.errGroup.Go(func() error {
		return runtime.serviceComponent.Serve(runtime.errCtx)
	})
	// ==== registry ================
	if runtime.registry != nil {
		if err = runtime.moduleContainer.InstallModule("registry", runtime.registry); err != nil {
			return
		}
	}
	// ==== http ================
	if runtime.httpServer != nil {
		if err = runtime.moduleContainer.InstallModule("http server", runtime.httpServer); err != nil {
			return
		}
	}
	// ==== service ================
	if runtime.grpcServer != nil {
		// ==== admin service ================
		{
			service := admin_service.NewService(runtime.application)
			if err = runtime.serviceComponent.StartService("admin", service); err != nil {
				return
			}
		}
		// ==== peer service ================
		{
			service := peer_service.NewService(runtime.application)
			if err = runtime.serviceComponent.StartService("peer", service); err != nil {
				return
			}
		}
	}
	// ==== registry ================
	if runtime.registry != nil {
		runtime.registry.Notify()
	}
	// ==== framework start ================
	for _, fw := range runtime.frameworks {
		if err = fw.OnFrameworkStart(); err != nil {
			return
		}
	}
	// ==== application start ================
	if err = runtime.application.OnStart(); err != nil {
		return
	}
	// ==== grpc ================
	if runtime.grpcServer != nil {
		if err = runtime.moduleContainer.InstallModule("grpc server", runtime.grpcServer); err != nil {
			return
		}
	}
	return nil
}

func (runtime *Runtime) onCreate() error {
	// log.Info("application", app.FullName, "start")
	application := runtime.application

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
			runtime.registry = r
		}
	}
	// ==== db ================
	runtime.dbClients = make(map[string]gira.DbClient)
	for name, c := range runtime.config.Db {
		if client, err := db.NewConfigDbClient(runtime.ctx, name, *c); err != nil {
			return err
		} else {
			runtime.dbClients[name] = client
			if name == gira.GAMEDB_NAME {
				runtime.gameDbClient = client
			} else if name == gira.RESOURCEDB_NAME {
				runtime.resourceDbClient = client
			} else if name == gira.STATDB_NAME {
				runtime.statDbClient = client
			} else if name == gira.ACCOUNTDB_NAME {
				runtime.accountDbClient = client
			} else if name == gira.LOGDB_NAME {
				runtime.logDbClient = client
			} else if name == gira.BEHAVIORDB_NAME {
				runtime.behaviorDbClient = client
			} else if name == gira.ACCOUNTCACHE_NAME {
				runtime.accountCacheClient = client
			} else if name == gira.ADMINCACHE_NAME {
				runtime.adminCacheClient = client
			} else if name == gira.ADMINDB_NAME {
				runtime.adminDbClient = client
			}
		}
	}

	// ==== 加载resource ================
	if resourceComponent, ok := runtime.application.(gira.ResourceComponent); ok {
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
		runtime.grpcServer = grpc.NewConfigGrpcServer(*runtime.config.Module.Grpc)
	}

	// ==== sdk================
	if runtime.config.Module.Sdk != nil {
		runtime.sdkComponent = sdk.NewConfigSdk(*runtime.config.Module.Sdk)
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
				runtime.httpServer = httpServer
			}
		}
	}

	// ==== gateway ================
	if runtime.config.Module.Gateway != nil {
		var handler gira.GatewayHandler
		if h, ok := application.(gira.GatewayHandler); ok {
			handler = h
		} else {
			for _, fw := range runtime.frameworks {
				if h, ok = fw.(gira.GatewayHandler); ok {
					handler = h
					break
				}
			}
		}
		if handler == nil {
			return gira.ErrGateHandlerNotImplement
		}
		if gate, err := gate.NewConfigServer(runtime.application, handler, *runtime.config.Module.Gateway); err != nil {
			return err
		} else {
			runtime.gate = gate
		}
	}

	// ==== framework create ================
	for _, fw := range runtime.frameworks {
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
