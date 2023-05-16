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
	"syscall"
	"time"

	"github.com/lujingwei002/gira/gins"
	"github.com/lujingwei002/gira/log"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/db"
	"github.com/lujingwei002/gira/gate"
	"github.com/lujingwei002/gira/grpc"
	"github.com/lujingwei002/gira/proj"
	"github.com/lujingwei002/gira/registry"
	"github.com/lujingwei002/gira/sdk"
	admin_service "github.com/lujingwei002/gira/service/admin"

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
	resourceLoader gira.ResourceLoader
	errGroup       *errgroup.Group
	config         *gira.Config

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
	Sdk                *sdk.Sdk
	Gate               *gate.Server
	GrpcServer         *grpc.GrpcServer
}

func newRuntime(args ApplicationArgs, application gira.Application) *Runtime {
	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	errGroup, errCtx := errgroup.WithContext(cancelCtx)
	runtime := &Runtime{
		BuildVersion: args.BuildVersion,
		BuildTime:    args.BuildTime,
		appId:        args.AppId,
		Application:  application,
		Frameworks:   make([]gira.Framework, 0),
		appType:      args.AppType,
		appName:      fmt.Sprintf("%s_%d", args.AppType, args.AppId),
		ctx:          cancelCtx,
		cancelFunc:   cancelFunc,
		errCtx:       errCtx,
		errGroup:     errGroup,
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

func (runtime *Runtime) serve() error {
	if err := runtime.start(); err != nil {
		return err
	}
	if err := runtime.wait(); err != nil {
		return err
	}
	return nil
}

func (runtime *Runtime) start() (err error) {
	if r, ok := runtime.Application.(RunnableApplication); !ok {
		return gira.ErrTodo
	} else {
		r.OnApplicationRun(runtime)
	}
	gira.OnApplicationNew(runtime.Application)
	runtime.Frameworks = runtime.Application.OnFrameworkInit()
	if err = runtime.init(); err != nil {
		return
	}
	defer func() {
		if err != nil {
			runtime.onDestory()
		}
	}()
	if err = runtime.onAwake(); err != nil {
		return
	}
	if err = runtime.onStart(); err != nil {
		return
	}
	return nil
}

func (runtime *Runtime) onDestory() {
	log.Infow("runtime onDestory")
	for _, fw := range runtime.Frameworks {
		log.Info("framework onDestory", "name")
		if err := fw.OnFrameworkDestory(); err != nil {
			log.Warnw("framework on destory fail", "error", err)
		}
	}
	if runtime.Registry != nil {
		log.Infow("registry onDestory")
		if err := runtime.Registry.OnDestory(); err != nil {
			log.Warnw("registry on destory fail", "error", err)
		}
	}
}

func (runtime *Runtime) onStart() error {
	application := runtime.Application

	if runtime.Registry != nil {
		if err := runtime.Registry.OnStart(); err != nil {
			return err
		}
	}
	// 注册grpc服务
	if runtime.GrpcServer != nil {
		if _, ok := application.(gira.GrpcServer); !ok {
			return gira.ErrGrpcServerNotImplement
		}
	}
	if runtime.config.Module.Admin != nil {
		if runtime.GrpcServer == nil {
			return gira.ErrGrpcServerNotOpen
		}
		service := admin_service.NewService(runtime.Application)
		if err := service.Register(runtime.GrpcServer.Server()); err != nil {
			return err
		}
	}
	for _, fw := range runtime.Frameworks {
		if err := fw.OnFrameworkStart(); err != nil {
			return err
		}
	}
	if err := runtime.Application.OnStart(); err != nil {
		return err
	}
	// 开始侦听
	if runtime.config.Module.Grpc != nil {
		if err := runtime.GrpcServer.OnStart(runtime.Application, runtime.errGroup, runtime.errCtx); err != nil {
			return err
		}
	}
	return nil
}

func (runtime *Runtime) onAwake() error {
	// log.Info("application", app.FullName, "start")
	application := runtime.Application

	// ====内置的服务=============
	if runtime.config.Pprof.Port != 0 {
		go func() {
			log.Infof("pprof start at http://%s:%d/debug", runtime.config.Pprof.Bind, runtime.config.Pprof.Port)
			http.ListenAndServe(fmt.Sprintf("%s:%d", runtime.config.Pprof.Bind, runtime.config.Pprof.Port), nil)
		}()
	}
	if runtime.config.Module.Sdk != nil {
		runtime.Sdk = sdk.NewConfigSdk(*runtime.config.Module.Sdk)
	}
	if runtime.config.Module.Etcd != nil {
		if r, err := registry.NewConfigRegistry(runtime.config.Module.Etcd, application); err != nil {
			return err
		} else {
			runtime.Registry = r
		}
	}
	runtime.DbClients = make(map[string]gira.DbClient)
	for name, c := range runtime.config.Db {
		if client, err := db.ConfigDbClient(runtime.ctx, name, *c); err != nil {
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
	// if runtime.config.Module.AdminCache != nil {
	// 	runtime.AdminCacheClient = db.NewAdminCacheClient()
	// 	if err := runtime.AdminCacheClient.OnAwake(runtime.ctx, *runtime.config.Module.AdminCache); err != nil {
	// 		return err
	// 	}
	// }
	// if runtime.config.Module.AccountCache != nil {
	// 	runtime.AccountCacheClient = db.NewAccountCacheClient()
	// 	if err := runtime.AccountCacheClient.OnAwake(runtime.ctx, *runtime.config.Module.AccountCache); err != nil {
	// 		return err
	// 	}
	// }
	if runtime.config.Module.Grpc != nil {
		runtime.GrpcServer = grpc.NewConfigGrpcServer(*runtime.config.Module.Grpc)
	}
	// if runtime.config.Module.GameDb != nil {
	// 	if client, err := db.ConfigGameDbClient(runtime.ctx, *runtime.config.Module.GameDb); err != nil {
	// 		return err
	// 	} else {
	// 		runtime.GameDbClient = client
	// 	}
	// }
	// if runtime.config.Module.BehaviorDb != nil {
	// 	if client, err := db.ConfigBehaviorDbClient(runtime.ctx, *runtime.config.Module.BehaviorDb); err != nil {
	// 		return err
	// 	} else {
	// 		runtime.BehaviorDbClient = client
	// 	}
	// }
	// if runtime.config.Module.AccountDb != nil {
	// 	if client, err := db.ConfigAccountDbClient(runtime.ctx, *runtime.config.Module.AccountDb); err != nil {
	// 		return err
	// 	} else {
	// 		runtime.AccountDbClient = client
	// 	}
	// }
	// if runtime.config.Module.StatDb != nil {
	// 	runtime.StatDbClient = db.NewStatDbClient()
	// 	if err := runtime.StatDbClient.OnAwake(runtime.ctx, *runtime.config.Module.StatDb); err != nil {
	// 		return err
	// 	}
	// }
	// if runtime.config.Module.AdminDb != nil {
	// 	runtime.adminDbClient = db.NewAdminDbClient()
	// 	if err := runtime.adminDbClient.OnAwake(runtime.ctx, *runtime.config.Module.AdminDb); err != nil {
	// 		return err
	// 	}
	// }
	// if runtime.config.Module.ResourceDb != nil {
	// 	runtime.ResourceDbClient = db.NewResourceDbClient()
	// 	if err := runtime.ResourceDbClient.OnAwake(runtime.ctx, *runtime.config.Module.ResourceDb); err != nil {
	// 		return err
	// 	}
	// }
	if runtime.config.Module.Http != nil {
		if handler, ok := application.(gira.HttpHandler); !ok {
			return gira.ErrHttpHandlerNotImplement
		} else {
			router := handler.HttpHandler()
			if httpServer, err := gins.NewConfigHttpServer(application, *runtime.config.Module.Http, router); err != nil {
				return err
			} else {
				runtime.HttpServer = httpServer

			}
		}
	}
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
	// res加载
	if resourceManager, ok := runtime.Application.(gira.ResourceManager); ok {
		resourceLoader := resourceManager.ResourceLoader()
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
	for _, fw := range runtime.Frameworks {
		if err := fw.OnFrameworkAwake(application); err != nil {
			return err
		}
	}
	if err := application.OnAwake(); err != nil {
		return err
	}
	// 创建场景
	// scene := CreateScene()
	// app.MainScene = scene
	// 等待关闭
	return nil
}

func (runtime *Runtime) wait() error {
	ctrlFunc := func() error {
		quit := make(chan os.Signal)
		defer close(quit)
		signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
		for {
			select {
			case s := <-quit:
				switch s {
				case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
					log.Info("ctrl shutdown begin.")
					runtime.onDestory()
					runtime.cancelFunc()
					log.Info("ctrl shutdown end.")
					return nil
				case syscall.SIGUSR1:
					log.Info("sigusr1.")
				case syscall.SIGUSR2:
					log.Info("sigusr2.")
				default:
					log.Info("single x")
				}
			case <-runtime.ctx.Done():
				log.Info("recv ctx:", runtime.Err().Error())
				return nil
			}
		}
	}
	runtime.Go(ctrlFunc)
	if err := runtime.errGroup.Wait(); err != nil {
		log.Infow("application down", "full_name", runtime.appFullName, "error", err)
		return err
	} else {
		log.Infow("application down", "full_name", runtime.appFullName)
		return nil
	}
}
