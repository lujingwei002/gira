package app

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/lujingwei002/gira/log"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/db"
	"github.com/lujingwei002/gira/gat"
	"github.com/lujingwei002/gira/grpc"
	"github.com/lujingwei002/gira/http"
	"github.com/lujingwei002/gira/registry"
	"github.com/lujingwei002/gira/sdk"
	admin_service "github.com/lujingwei002/gira/service/admin"

	"golang.org/x/sync/errgroup"
)

type ApplicationArgs struct {
	AppType string /// 服务名
	AppId   int32  /// 服务id
}

type ApplicationContext interface {
	Go(f func() error)
	Cancel() error
}

// / @Component
type Application struct {
	gira.BaseComponent
	zone               string // 区名 wc|qq|hw|quick
	env                string // dev|local|qa|prd
	appId              int32
	appType            string           /// 服务类型
	appName            string           /// 服务名
	appFullName        string           /// 完整的服务名 Name_Id
	ProjectConf        gira.ProjectConf /// gira.yaml配置
	ProjectFilePath    string           /// 配置文件绝对路径, gira.yaml
	ConfigDir          string           /// config目录
	EnvDir             string           /// env目录
	ConfigFilePath     string           /// 内置配置文件
	RunConfigFilePath  string           /// 运行时的配置文件
	WorkDir            string           /// 工作目录
	LogDir             string           /// 日志目录
	RunDir             string           /// 运行目录
	Facade             gira.ApplicationFacade
	cancelFunc         context.CancelFunc
	cancelCtx          context.Context
	errCtx             context.Context
	errGroup           *errgroup.Group
	MainScene          *gira.Scene
	Config             *gira.Config
	HttpServer         *http.HttpServer
	Registry           *registry.Registry
	GameDbClient       *db.GameDbClient
	AccountDbClient    *db.AccountDbClient
	StatDbClient       *db.StatDbClient
	ResourceDbClient   *db.ResourceDbClient
	AccountCacheClient *db.AccountCacheClient
	Sdk                *sdk.Sdk
	Gate               *gat.Gate
	GrpcServer         *grpc.GrpcServer
	resourceLoader     gira.ResourceLoader
	adminClient        AdminClient
	adminDbClient      *db.AdminDbClient
}

type FacadeSetApplication interface {
	SetApplication(application *Application)
}

func newApplication(args ApplicationArgs, facade gira.ApplicationFacade) *Application {
	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	errGroup, errCtx := errgroup.WithContext(cancelCtx)
	application := &Application{
		appId:      args.AppId,
		Facade:     facade,
		appType:    args.AppType,
		appName:    fmt.Sprintf("%s_%d", args.AppType, args.AppId),
		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,
		errCtx:     errCtx,
		errGroup:   errGroup,
	}
	if v, ok := facade.(FacadeSetApplication); ok {
		v.SetApplication(application)
	}
	return application
}

func (app *Application) Err() error {
	return app.errCtx.Err()
}

func (app *Application) Go(f func() error) {
	app.errGroup.Go(f)
}

func (app *Application) init() error {
	// 初始化
	rand.Seed(time.Now().UnixNano())
	// 目录初始化
	if workDir, err := os.Getwd(); err != nil {
		return err
	} else {
		dir := workDir
		for {
			projectFilePath := filepath.Join(dir, "gira.yaml")
			if _, err := os.Stat(projectFilePath); err == nil {
				app.WorkDir = dir
				break
			}
			dir = filepath.Dir(dir)
			if dir == "/" || dir == "" {
				break
			}
		}
	}
	if app.WorkDir == "" {
		return gira.ErrProjectFileNotFound
	}
	os.Chdir(app.WorkDir)
	app.ProjectFilePath = filepath.Join(app.WorkDir, "gira.yaml")
	if _, err := os.Stat(app.ProjectFilePath); err != nil {
		return err
	}
	app.EnvDir = filepath.Join(app.WorkDir, "env")
	app.ConfigDir = filepath.Join(app.WorkDir, "config")
	if _, err := os.Stat(app.ConfigDir); err != nil {
		return err
	}
	app.RunDir = filepath.Join(app.WorkDir, "run")
	if _, err := os.Stat(app.RunDir); err != nil {
		if err := os.Mkdir(app.RunDir, 0755); err != nil {
			return err
		}
	}
	app.RunConfigFilePath = filepath.Join(app.RunDir, fmt.Sprintf("%s", app.appFullName))
	app.LogDir = filepath.Join(app.RunDir, "log")

	/*
		app.ConfigFilePath = filepath.Join(app.ConfigDir, fmt.Sprintf("%sconf.yaml", app.Name))
		if _, err := os.Stat(app.ConfigFilePath); err != nil {
			return err
		}*/
	// 读项目配置文件
	if err := app.ProjectConf.Read(app.ProjectFilePath); err != nil {
		return err
	}
	// 读应用配置文件
	if c, err := gira.LoadConfig(app.ConfigDir, app.EnvDir, app.appType, app.appId); err != nil {
		return err
	} else {
		app.env = c.Env
		app.zone = c.Zone
		app.appFullName = fmt.Sprintf("%s_%s_%s_%d", app.appType, app.zone, app.env, app.appId)
		app.Config = c
	}
	if err := app.Facade.OnConfigLoad(app.Config); err != nil {
		return err
	}
	if app.Config.Log != nil {
		if err := log.ConfigLog(app.Facade, *app.Config.Log); err != nil {
			return err
		}
	}
	runtime.GOMAXPROCS(app.Config.Thread)
	//var serviceConf ServiceConf
	// for _, conf := range app.ProjectConf.Services {
	// 	if conf.Name == app.Name {
	// 		log.Info("=======", conf.Name)
	// 	}
	// }
	// log.Infof("%+v\n", app.ProjectConf)
	return nil
}

func (app *Application) createHttpServer() {

}

func (app *Application) forver() error {
	if err := app.start(); err != nil {
		return err
	}
	if err := app.wait(); err != nil {
		return err
	}
	return nil
}

func (app *Application) start() error {
	// log.Info("application", app.FullName, "start")
	if err := app.init(); err != nil {
		return err
	}
	app.Facade.Awake()
	// 内置的服务
	if app.Config.Sdk != nil {
		app.Sdk = sdk.NewConfigSdk(*app.Config.Sdk)
	}

	if app.Config.Etcd != nil {
		if r, err := registry.NewConfigRegistry(app.Config.Etcd, app.Facade); err != nil {
			return err
		} else {
			app.Registry = r
		}
	}
	if app.Config.AccountCache != nil {
		app.AccountCacheClient = db.NewAccountCacheClient()
		if err := app.AccountCacheClient.Start(app.cancelCtx, *app.Config.AccountCache); err != nil {
			return err
		}
	}
	if app.Config.Grpc != nil {
		app.GrpcServer = grpc.NewConfigGrpcServer(*app.Config.Grpc)
		// admin service
		app.adminClient = admin_service.NewAdminClient()
		if app.Config.Admin != nil {
			if err := admin_service.Register(app.Facade, app.GrpcServer.Server()); err != nil {
				return err
			}
		}
		if err := app.GrpcServer.Start(app.Facade, app.errGroup, app.errCtx); err != nil {
			return err
		}
	}
	if app.Config.GameDb != nil {
		app.GameDbClient = db.NewGameDbClient()
		if err := app.GameDbClient.Start(app.cancelCtx, *app.Config.GameDb); err != nil {
			return err
		}
	}
	if app.Config.AccountDb != nil {
		app.AccountDbClient = db.NewAccountDbClient()
		if err := app.AccountDbClient.Start(app.cancelCtx, *app.Config.AccountDb); err != nil {
			return err
		}
	}
	if app.Config.StatDb != nil {
		app.StatDbClient = db.NewStatDbClient()
		if err := app.StatDbClient.Start(app.cancelCtx, *app.Config.StatDb); err != nil {
			return err
		}
	}
	if app.Config.AdminDb != nil {
		app.adminDbClient = db.NewAdminDbClient()
		if err := app.adminDbClient.Start(app.cancelCtx, *app.Config.AdminDb); err != nil {
			return err
		}
	}

	if app.Config.ResourceDb != nil {
		app.ResourceDbClient = db.NewResourceDbClient()
		if err := app.ResourceDbClient.Start(app.cancelCtx, *app.Config.ResourceDb); err != nil {
			return err
		}
	}
	// res加载
	if resourceManager, ok := app.Facade.(gira.ResourceManager); ok {
		resourceLoader := resourceManager.ResourceLoader()
		if resourceLoader != nil {
			app.resourceLoader = resourceLoader
			if err := app.resourceLoader.LoadResource("resource"); err != nil {
				return err
			}
		} else {
			return gira.ErrResourceLoaderNotImplement
		}
	} else {
		return gira.ErrResourceManagerNotImplement
	}

	if app.Config.Http != nil {
		if handler, ok := app.Facade.(http.HttpHandler); !ok {
			return gira.ErrHttpHandlerNotImplement
		} else {
			router := handler.HttpHandler()
			if httpServer, err := http.NewConfigHttpServer(app.Facade, *app.Config.Http, router); err != nil {
				return err
			} else {
				app.HttpServer = httpServer

			}
		}
	}
	if app.Config.Gate != nil {
		if gate, err := gat.NewConfigGate(app.Facade, *app.Config.Gate); err != nil {
			return err
		} else {
			app.Gate = gate
		}
	}
	if app.Registry != nil {
		if err := app.Registry.Notify(); err != nil {
			return err
		}
	}
	if err := app.Facade.OnApplicationLoad(); err != nil {
		return err
	}
	// 创建场景
	// scene := CreateScene()
	// app.MainScene = scene
	// 等待关闭
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
					app.cancelFunc()
					log.Info("ctrl shutdown end.")
					return nil
				case syscall.SIGUSR1:
					log.Info("sigusr1.")
				case syscall.SIGUSR2:
					log.Info("sigusr2.")
				default:
					log.Info("single x")
				}
			case <-app.cancelCtx.Done():
				log.Info("recv ctx:", app.Err().Error())
				return nil
			}
		}
	}
	// app.Go(scene.forver)
	app.Go(ctrlFunc)
	return nil
}

func (app *Application) wait() error {
	if err := app.errGroup.Wait(); err != nil {
		log.Infow("application down", "full_name", app.appFullName, "error", err)
		return err
	} else {
		log.Infow("application down", "full_name", app.appFullName)
		return nil
	}
}
