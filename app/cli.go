package app

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/lujingwei002/gira/facade"
	"github.com/lujingwei002/gira/log"
	"github.com/lujingwei002/gira/proj"
	"github.com/lujingwei002/gira/service/admin/adminpb"
	peer_service "github.com/lujingwei002/gira/service/peer"
	"github.com/lujingwei002/gira/service/peer/peerpb"
	"gopkg.in/yaml.v3"

	"github.com/lujingwei002/gira"
	"github.com/urfave/cli/v2"
)

type ClientApplicationFacade struct {
}

func (s *ClientApplicationFacade) OnConfigLoad(c *gira.Config) error {
	return nil
}
func (s *ClientApplicationFacade) OnCreate() error {
	return nil
}
func (s *ClientApplicationFacade) OnStart() error {
	return nil
}
func (s *ClientApplicationFacade) OnStop() error {
	return nil
}

// 需要两个系参数 xx -id 1 start|stop|restart
func Cli(name string, appVersion string, buildTime string, applicationFacade gira.ApplicationFacade) error {

	app := &cli.App{
		Name: "gira service",
		// Authors:     []*cli.Author{&cli.Author{Name: "lujingwei", Email: "lujingwei002@qq.com"}},
		Description: "gira service",
		Flags:       []cli.Flag{},
		Action:      runAction,
		Metadata: map[string]interface{}{
			"application": applicationFacade,
			"name":        name,
			"buildTime":   buildTime,
			"appVersion":  appVersion,
		},

		Commands: []*cli.Command{
			{
				Name: "start",
				//Aliases: []string{"start"},
				Usage:  "Start service",
				Action: startAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "id",
						Usage:    "service id",
						Required: true,
					},
					&cli.StringFlag{
						Aliases:  []string{"f"},
						Name:     "file",
						Usage:    "config file",
						Required: false,
					},
					&cli.StringFlag{
						Aliases:  []string{"c"},
						Name:     "env",
						Usage:    "env config file",
						Required: false,
					},
				},
			},
			{
				Name:   "status",
				Usage:  "Display status",
				Action: statusAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "id",
						Usage:    "service id",
						Required: true,
					},
					&cli.StringFlag{
						Aliases:  []string{"f"},
						Name:     "file",
						Usage:    "config file",
						Required: false,
					},
					&cli.StringFlag{
						Aliases:  []string{"c"},
						Name:     "env",
						Usage:    "env config file",
						Required: false,
					},
				},
			},
			{
				Name:   "unregister",
				Usage:  "unregister service",
				Action: unregisterAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "id",
						Usage:    "service id",
						Required: true,
					},
					&cli.StringFlag{
						Aliases:  []string{"f"},
						Name:     "file",
						Usage:    "config file",
						Required: false,
					},
					&cli.StringFlag{
						Aliases:  []string{"c"},
						Name:     "env",
						Usage:    "env config file",
						Required: false,
					},
				},
			},
			{
				Name:   "stop",
				Usage:  "Stop service",
				Action: stopAction,
			},
			{
				Name:   "restart",
				Usage:  "Restart service",
				Action: restartAction,
			},
			{
				Name:   "reload",
				Usage:  "Reload config",
				Action: reloadAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "id",
						Usage:    "service id",
						Required: true,
					},
					&cli.StringFlag{
						Aliases:  []string{"f"},
						Name:     "file",
						Usage:    "config file",
						Required: false,
					},
					&cli.StringFlag{
						Aliases:  []string{"c"},
						Name:     "env",
						Usage:    "env config file",
						Required: false,
					},
				},
			},
			{
				Name:   "version",
				Usage:  "Build version",
				Action: versionAction,
			},
			{
				Name:   "time",
				Usage:  "Build time",
				Action: timeAction,
			},
			{
				Name:   "env",
				Usage:  "Print env",
				Action: envAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "id",
						Usage:    "service id",
						Required: true,
					},
					&cli.StringFlag{
						Aliases:  []string{"f"},
						Name:     "file",
						Usage:    "config file",
						Required: false,
					},
					&cli.StringFlag{
						Aliases:  []string{"c"},
						Name:     "env",
						Usage:    "env config file",
						Required: false,
					},
				},
			},
			{
				Name:   "config",
				Usage:  "Print config",
				Action: configAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "id",
						Usage:    "service id",
						Required: true,
					},
					&cli.StringFlag{
						Aliases:  []string{"f"},
						Name:     "file",
						Usage:    "config file",
						Required: false,
					},
					&cli.StringFlag{
						Aliases:  []string{"c"},
						Name:     "env",
						Usage:    "env config file",
						Required: false,
					},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		return err
	}
	return nil
}

func StartAsClient(applicationFacade gira.ApplicationFacade, appId int32, appType string, configFilePath string, dotEnvFilePath string) error {
	application := newApplication(gira.ApplicationArgs{
		AppType:        appType,
		AppId:          appId,
		ConfigFilePath: configFilePath,
		DotEnvFilePath: dotEnvFilePath,
		Facade:         applicationFacade,
	})
	return application.start()
}

func StartAsServer(applicationFacade gira.ApplicationFacade, appId int32, appType string, configFilePath string, dotEnvFilePath string) error {
	application := newApplication(gira.ApplicationArgs{
		AppType:        appType,
		AppId:          appId,
		ConfigFilePath: configFilePath,
		DotEnvFilePath: dotEnvFilePath,
		Facade:         applicationFacade,
	})
	return application.start()
}

// 启动应用
func startAction(args *cli.Context) error {
	appId := int32(args.Int("id"))
	configFilePath := args.String("file")
	dotEnvFilePath := args.String("env")
	applicationFacade, _ := args.App.Metadata["application"].(gira.ApplicationFacade)
	appType, _ := args.App.Metadata["name"].(string)
	appVersion, _ := args.App.Metadata["appVersion"].(string)
	var buildTime int64
	if v, ok := args.App.Metadata["buildTime"].(string); ok {
		if t, err := strconv.Atoi(v); err != nil {
			buildTime = 0
		} else {
			buildTime = int64(t)
		}
	}
	log.Println("build version:", appVersion)
	log.Println("build time:", buildTime)
	log.Infof("%s %d starting...", appType, appId)
	runtime := newApplication(gira.ApplicationArgs{
		AppType:        appType,
		AppId:          appId,
		AppVersion:     appVersion,
		BuildTime:      buildTime,
		DotEnvFilePath: dotEnvFilePath,
		ConfigFilePath: configFilePath,
		Facade:         applicationFacade,
	})

	if err := runtime.start(); err != nil {
		return err
	}
	if err := runtime.Wait(); err != nil {
		return err
	}
	return nil
}

// 打印应用构建版本
func versionAction(args *cli.Context) error {
	appVersion, _ := args.App.Metadata["appVersion"].(string)
	fmt.Println(appVersion)
	return nil
}

// 打印应用构建时间
func timeAction(args *cli.Context) error {
	buildTime, _ := args.App.Metadata["buildTime"].(string)
	fmt.Println(buildTime)
	return nil
}

// 打印应用环境变量
func envAction(args *cli.Context) error {
	appId := int32(args.Int("id"))
	appType, _ := args.App.Metadata["name"].(string)
	configFilePath := args.String("file")
	dotEnvFilePath := args.String("env")
	if _, err := gira.LoadApplicationConfig(configFilePath, dotEnvFilePath, appType, appId); err != nil {
		return err
	} else {
		for _, k := range os.Environ() {
			fmt.Println(k, os.Getenv(k))
		}
	}
	return nil
}

// 打印应用配置
func configAction(args *cli.Context) error {
	appId := int32(args.Int("id"))
	configFilePath := args.String("file")
	dotEnvFilePath := args.String("env")
	appType, _ := args.App.Metadata["name"].(string)
	if c, err := gira.LoadApplicationConfig(configFilePath, dotEnvFilePath, appType, appId); err != nil {
		return err
	} else {
		c.Raw = nil
		if s, err := yaml.Marshal(c); err != nil {
			return err
		} else {
			fmt.Println(string(s))
		}
	}
	return nil
}

func stopAction(args *cli.Context) error {
	return nil
}

func restartAction(args *cli.Context) error {
	return nil
}

func statusAction(args *cli.Context) error {
	appId := int32(args.Int("id"))
	configFilePath := args.String("file")
	dotEnvFilePath := args.String("env")
	if len(configFilePath) <= 0 {
		configFilePath = path.Join(proj.Config.ConfigDir, "cli.yaml")
	}
	if len(dotEnvFilePath) <= 0 {
		dotEnvFilePath = path.Join(proj.Config.EnvDir, ".env")
	}
	if err := StartAsClient(&ClientApplicationFacade{}, appId, "cli", configFilePath, dotEnvFilePath); err != nil {
		return err
	}
	ctx := facade.Context()
	serviceName := peer_service.GetServiceName()
	if _, err := peerpb.DefaultPeerClients.Unicast().Where(serviceName).HealthCheck(ctx, &peerpb.HealthCheckRequest{}); err != nil {
		log.Println(err)
		log.Println("dead")
		return nil
	} else {
		log.Println("alive")
		return nil
	}
}

func unregisterAction(args *cli.Context) error {
	appId := int32(args.Int("id"))
	appType, _ := args.App.Metadata["name"].(string)
	configFilePath := args.String("file")
	dotEnvFilePath := args.String("env")
	if len(configFilePath) <= 0 {
		configFilePath = path.Join(proj.Config.ConfigDir, "cli.yaml")
	}

	if len(dotEnvFilePath) <= 0 {
		dotEnvFilePath = path.Join(proj.Config.EnvDir, ".env")
	}
	if err := StartAsClient(&ClientApplicationFacade{}, appId, "cli", configFilePath, dotEnvFilePath); err != nil {
		return err
	}
	appFullName := gira.FormatAppFullName(appType, appId, facade.GetZone(), facade.GetEnv())
	if err := facade.UnregisterPeer(appFullName); err != nil {
		log.Println(err)
		return nil
	} else {
		log.Println("OK")
		return nil
	}
}

func reloadAction(args *cli.Context) error {
	appId := int32(args.Int("id"))
	configFilePath := args.String("file")
	dotEnvFilePath := args.String("env")
	if len(configFilePath) <= 0 {
		configFilePath = path.Join(proj.Config.ConfigDir, "cli.yaml")
	}
	if len(dotEnvFilePath) <= 0 {
		dotEnvFilePath = path.Join(proj.Config.EnvDir, ".env")
	}
	if err := StartAsClient(&ClientApplicationFacade{}, appId, "cli", configFilePath, dotEnvFilePath); err != nil {
		return err
	}
	ctx := facade.Context()
	serviceName := peer_service.GetServiceName()
	if _, err := adminpb.DefaultAdminClients.Unicast().Where(serviceName).ReloadResource(ctx, &adminpb.ReloadResourceRequest{}); err != nil {
		log.Println(err)
		return nil
	} else {
		log.Println("OK")
		return nil
	}
}

func runAction(args *cli.Context) error {
	return nil
}
