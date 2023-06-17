package app

/*

实现gira.Application接口的客户程序

*/

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/lujingwei002/gira/log"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/db"
	"github.com/lujingwei002/gira/proj"
	"github.com/lujingwei002/gira/registry"
	"github.com/lujingwei002/gira/sdk"
	"github.com/lujingwei002/gira/service/admin/admin_grpc"

	_ "net/http/pprof"

	"github.com/robfig/cron"
	"golang.org/x/sync/errgroup"
)

type DefaultClientApplicationFacade struct {
}

func (s *DefaultClientApplicationFacade) OnConfigLoad(c *gira.Config) error {
	return nil
}
func (s *DefaultClientApplicationFacade) OnCreate() error {
	return nil
}
func (s *DefaultClientApplicationFacade) OnStart() error {
	return nil
}
func (s *DefaultClientApplicationFacade) OnStop() error {
	return nil
}

type Client struct {
	zone               string // 区名 wc|qq|hw|quick
	env                string // dev|local|qa|prd
	appId              int32
	appType            string /// 服务类型
	appName            string /// 服务名
	appFullName        string /// 完整的服务名 Name_Id
	cancelFunc         context.CancelFunc
	ctx                context.Context
	errCtx             context.Context
	errGroup           *errgroup.Group
	resourceLoader     gira.ResourceLoader
	config             *gira.Config
	stopChan           chan struct{}
	status             int64
	respositoryVersion string
	buildTime          int64
	upTime             int64
	projectFilePath    string /// 配置文件绝对路径, gira.yaml
	configDir          string /// config目录
	envDir             string /// env目录
	runConfigFilePath  string /// 运行时的配置文件
	workDir            string /// 工作目录
	logDir             string /// 日志目录
	runDir             string /// 运行目录
	applicationFacade  gira.ApplicationFacade
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
	sdk                *sdk.SdkComponent
	cron               *cron.Cron
	configFilePath     string
	dotEnvFilePath     string
}

func newClientApplication(args ApplicationArgs, applicationFacade gira.ApplicationFacade) *Client {
	ctx, cancelFunc := context.WithCancel(context.Background())
	errGroup, errCtx := errgroup.WithContext(ctx)
	application := &Client{
		respositoryVersion: args.RespositoryVersion,
		buildTime:          args.BuildTime,
		appId:              args.AppId,
		configFilePath:     args.ConfigFilePath,
		dotEnvFilePath:     args.DotEnvFilePath,
		applicationFacade:  applicationFacade,
		appType:            args.AppType,
		appName:            fmt.Sprintf("%s_%d", args.AppType, args.AppId),
		ctx:                ctx,
		cancelFunc:         cancelFunc,
		errCtx:             errCtx,
		errGroup:           errGroup,
		stopChan:           make(chan struct{}, 1),
	}
	return application
}

func (application *Client) init() error {
	// 初始化
	rand.Seed(time.Now().UnixNano())
	application.upTime = time.Now().Unix()
	// 项目配置初始化
	application.workDir = proj.Config.ProjectDir
	if err := os.Chdir(application.workDir); err != nil {
		return err
	}
	application.projectFilePath = proj.Config.ProjectConfFilePath
	if _, err := os.Stat(application.projectFilePath); err != nil {
		return err
	}
	application.envDir = proj.Config.EnvDir
	application.configDir = proj.Config.ConfigDir
	if _, err := os.Stat(application.configDir); err != nil {
		return err
	}
	application.runDir = proj.Config.RunDir
	if _, err := os.Stat(application.runDir); err != nil {
		if err := os.Mkdir(application.runDir, 0755); err != nil {
			return err
		}
	}
	application.runConfigFilePath = filepath.Join(application.runDir, fmt.Sprintf("%s", application.appFullName))
	application.logDir = proj.Config.LogDir

	// 读应用配置文件
	if c, err := gira.LoadApplicationConfig(application.configFilePath, application.dotEnvFilePath, application.appType, application.appId); err != nil {
		return err
	} else {
		application.env = c.Env
		application.zone = c.Zone
		application.appFullName = fmt.Sprintf("%s_%s_%s_%d", application.appType, application.zone, application.env, application.appId)
		application.config = c
	}
	if err := application.applicationFacade.OnConfigLoad(application.config); err != nil {
		return err
	}
	// 初始化日志
	if application.config.Log != nil {
		if err := log.ConfigLog(*application.config.Log); err != nil {
			return err
		}
	}
	runtime.GOMAXPROCS(application.config.Thread)
	return nil
}

func (application *Client) start() (err error) {
	// 初始化
	if err = application.init(); err != nil {
		return
	}
	// 设置全局对象
	gira.OnApplicationCreate(application)
	if err = application.onCreate(); err != nil {
		return
	}
	if err = application.onStart(); err != nil {
		return
	}
	application.Go(func() error {
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
					application.onStop()
					application.cancelFunc()
					log.Info("application interrupt end")
					return gira.ErrInterrupt
				case syscall.SIGUSR1:
				case syscall.SIGUSR2:
				default:
				}
			// 主动停止
			case <-application.stopChan:
				log.Info("application stop begin")
				application.onStop()
				application.cancelFunc()
				log.Info("application stop end")
				return nil
			}
		}
	})
	application.status = application_status_started
	return nil
}

// 主动关闭, 启动成功后才可以调用
func (application *Client) Stop() error {
	if !atomic.CompareAndSwapInt64(&application.status, application_status_started, application_status_stopped) {
		return nil
	}
	application.stopChan <- struct{}{}
	return nil
}

func (application *Client) onStop() {
	log.Infow("runtime on stop")
	// application stop
	application.applicationFacade.OnStop()
	// service stop
	if application.registry != nil {
		application.registry.Stop()
	}
}

func (application *Client) onStart() (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
		if err != nil {
			log.Errorw("+++++++++++++++++++++++++++++")
			log.Errorw("++         启动异常       +++", "error", err)
			log.Errorw("+++++++++++++++++++++++++++++")
			application.onStop()
			application.cancelFunc()
			application.errGroup.Wait()
		}
	}()
	// ==== cron ================
	application.cron.Start()

	// ==== registry ================
	if application.registry != nil {
		if err = application.registry.StartAsClient(); err != nil {
			return
		}
	}
	// ==== registry ================
	if application.registry != nil {
		application.registry.Notify()
	}
	// ==== application start ================
	if err = application.applicationFacade.OnStart(); err != nil {
		return
	}
	// ==== registry ================
	// if application.registry != nil {
	// application.errGroup.Go(func() error {
	// return application.registry.Watch()
	// })
	// }
	return nil
}

func (application *Client) onCreate() error {
	// log.Info("applicationFacade", app.FullName, "start")
	applicationFacade := application.applicationFacade

	// ==== pprof ================
	if application.config.Pprof.Port != 0 {
		go func() {
			log.Infof("pprof start at http://%s:%d/debug", application.config.Pprof.Bind, application.config.Pprof.Port)
			http.ListenAndServe(fmt.Sprintf("%s:%d", application.config.Pprof.Bind, application.config.Pprof.Port), nil)
		}()
	}
	// ==== cron ================
	application.cron = cron.New()
	// ==== registry ================
	log.Println("ffffffffff")
	if application.config.Module.Etcd != nil {
		log.Println("ffffffffff222")
		if r, err := registry.NewConfigRegistry(application.ctx, application.config.Module.Etcd); err != nil {
			return err
		} else {
			application.registry = r
		}
	}
	// ==== db ================
	application.dbClients = make(map[string]gira.DbClient)
	for name, c := range application.config.Db {
		if client, err := db.NewConfigDbClient(application.ctx, name, *c); err != nil {
			return err
		} else {
			application.dbClients[name] = client
			if name == gira.GAMEDB_NAME {
				application.gameDbClient = client
			} else if name == gira.RESOURCEDB_NAME {
				application.resourceDbClient = client
			} else if name == gira.STATDB_NAME {
				application.statDbClient = client
			} else if name == gira.ACCOUNTDB_NAME {
				application.accountDbClient = client
			} else if name == gira.LOGDB_NAME {
				application.logDbClient = client
			} else if name == gira.BEHAVIORDB_NAME {
				application.behaviorDbClient = client
			} else if name == gira.ACCOUNTCACHE_NAME {
				application.accountCacheClient = client
			} else if name == gira.ADMINCACHE_NAME {
				application.adminCacheClient = client
			} else if name == gira.ADMINDB_NAME {
				application.adminDbClient = client
			}
		}
	}

	// ==== 加载resource ================
	if resourceComponent, ok := application.applicationFacade.(gira.ResourceSource); ok {
		resourceLoader := resourceComponent.GetResourceLoader()
		if resourceLoader != nil {
			application.resourceLoader = resourceLoader
			if err := application.resourceLoader.LoadResource("resource"); err != nil {
				return err
			}
		} else {
			// return gira.ErrResourceLoaderNotImplement
		}
	} else {
		// return gira.ErrResourceManagerNotImplement
	}

	// ==== sdk================
	if application.config.Module.Sdk != nil {
		application.sdk = sdk.NewConfigSdk(*application.config.Module.Sdk)
	}

	// ==== application create ================
	if err := applicationFacade.OnCreate(); err != nil {
		return err
	}
	return nil
}

// 等待中断
func (application *Client) Wait() error {
	err := application.errGroup.Wait()
	log.Infow("application down", "full_name", application.appFullName, "error", err)
	return err
}

// 负载实现gira声明的接口和启动的模块

// ================== implement gira.Application ==================
// 返回配置
func (application *Client) GetConfig() *gira.Config {
	return application.config
}

// 返回构建版本
func (application *Client) GetRespositoryVersion() string {
	return application.respositoryVersion

}

// 返回构建时间
func (application *Client) GetBuildTime() int64 {
	return application.buildTime
}
func (application *Client) GetUpTime() int64 {
	return time.Now().Unix() - application.buildTime
}

// 返回应用id
func (application *Client) GetAppId() int32 {
	return application.appId
}

// 返回应用类型
func (application *Client) GetAppType() string {
	return application.appType
}

// 返回应用名
func (application *Client) GetAppName() string {
	return application.appName
}

// 返回应用全名
func (application *Client) GetAppFullName() string {
	return application.appFullName
}

func (application *Client) GetWorkDir() string {
	return application.workDir
}

func (application *Client) GetLogDir() string {
	return application.logDir
}

// ================== context =========================

func (application *Client) Context() context.Context {
	return application.ctx
}

func (application *Client) Quit() {
	application.cancelFunc()
}

func (application *Client) Done() <-chan struct{} {
	return application.ctx.Done()
}

func (application *Client) Go(f func() error) {
	application.errGroup.Go(f)
}

func (application *Client) Err() error {
	return application.errCtx.Err()
}

// ================== framework =========================
// 返回框架列表
func (application *Client) Frameworks() []gira.Framework {
	return nil
}

// ================== gira.SdkComponent ==================
func (application *Client) GetSdk() gira.Sdk {
	if application.sdk == nil {
		return nil
	} else {
		return application.sdk
	}
}

// ================== gira.RegistryComponent ==================
func (application *Client) GetRegistry() gira.Registry {
	return nil
}
func (application *Client) GetRegistryClient() gira.RegistryClient {
	if application.registry == nil {
		return nil
	} else {
		return application.registry
	}
}

// ================== gira.ServiceContainer ==================
func (application *Client) GetServiceContainer() gira.ServiceContainer {
	return nil
}

// ================== gira.GrpcServerComponent ==================
func (application *Client) GetGrpcServer() gira.GrpcServer {
	return nil
}

// ================== implement gira.ResourceComponent ==================
func (application *Client) GetResourceLoader() gira.ResourceLoader {
	return application.resourceLoader
}

// ================== implement gira.AdminClient ==================
// 重载配置
func (application *Client) BroadcastReloadResource(ctx context.Context, name string) (result gira.BroadcastReloadResourceResult, err error) {
	req := &admin_grpc.ReloadResourceRequest{
		Name: name,
	}
	result, err = admin_grpc.DefaultAdminClients.Broadcast().ReloadResource(ctx, req)
	return
}

// ================== implement gira.DbClientComponent ==================
func (application *Client) GetAccountDbClient() gira.DbClient {
	if application.accountDbClient == nil {
		return nil
	} else {
		return application.accountDbClient
	}
}

func (application *Client) GetGameDbClient() gira.DbClient {
	if application.gameDbClient == nil {
		return nil
	} else {
		return application.gameDbClient
	}
}

func (application *Client) GetLogDbClient() gira.DbClient {
	if application.logDbClient == nil {
		return nil
	} else {
		return application.logDbClient
	}
}

func (application *Client) GetBehaviorDbClient() gira.DbClient {
	if application.behaviorDbClient == nil {
		return nil
	} else {
		return application.behaviorDbClient
	}
}

func (application *Client) GetStatDbClient() gira.DbClient {
	if application.statDbClient == nil {
		return nil
	} else {
		return application.statDbClient
	}
}

func (application *Client) GetAccountCacheClient() gira.DbClient {
	if application.accountCacheClient == nil {
		return nil
	} else {
		return application.accountCacheClient
	}
}

func (application *Client) GetAdminCacheClient() gira.DbClient {
	if application.adminCacheClient == nil {
		return nil
	} else {
		return application.adminCacheClient
	}
}

func (application *Client) GetResourceDbClient() gira.DbClient {
	if application.resourceDbClient == nil {
		return nil
	} else {
		return application.resourceDbClient
	}
}

func (application *Client) GetAdminDbClient() gira.DbClient {
	if application.adminDbClient == nil {
		return nil
	} else {
		return application.adminDbClient
	}
}

// ================== gira.Cron ==================
func (application *Client) GetCron() gira.Cron {
	if application.cron == nil {
		return nil
	} else {
		return application.cron
	}
}
