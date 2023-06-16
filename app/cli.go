package app

import (
	"fmt"
	"os"
	"strconv"

	"github.com/lujingwei002/gira/log"

	"github.com/lujingwei002/gira"
	"github.com/urfave/cli/v2"
)

// 需要两个系参数 xx -id 1 start|stop|restart
func Cli(name string, respositoryVersion string, buildTime string, applicationFacade gira.ApplicationFacade) error {

	app := &cli.App{
		Name: "gira service",
		// Authors:     []*cli.Author{&cli.Author{Name: "lujingwei", Email: "lujingwei002@qq.com"}},
		Description: "gira service",
		Flags:       []cli.Flag{},
		Action:      runAction,
		Metadata: map[string]interface{}{
			"application":        applicationFacade,
			"name":               name,
			"buildTime":          buildTime,
			"respositoryVersion": respositoryVersion,
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
				Usage:  "打印环境变量",
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
				Usage:  "打印配置",
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

func Start(applicationFacade gira.ApplicationFacade, appId int32, appType string, configFilePath string, dotEnvFilePath string) error {
	application := newApplication(ApplicationArgs{
		AppType:        appType,
		AppId:          appId,
		ConfigFilePath: configFilePath,
		DotEnvFilePath: dotEnvFilePath,
	}, applicationFacade)
	return application.start()
}

// 启动应用
func startAction(args *cli.Context) error {
	appId := int32(args.Int("id"))
	configFilePath := args.String("file")
	dotEnvFilePath := args.String("env")
	applicationFacade, _ := args.App.Metadata["application"].(gira.ApplicationFacade)
	appType, _ := args.App.Metadata["name"].(string)
	respositoryVersion, _ := args.App.Metadata["respositoryVersion"].(string)
	var buildTime int64
	if v, ok := args.App.Metadata["buildTime"].(string); ok {
		if t, err := strconv.Atoi(v); err != nil {
			buildTime = 0
		} else {
			buildTime = int64(t)
		}
	}
	log.Println("build version:", respositoryVersion)
	log.Println("build time:", buildTime)
	log.Println("config file path:", configFilePath)
	log.Println("env file path:", dotEnvFilePath)
	log.Infof("%s %d starting...", appType, appId)
	runtime := newApplication(ApplicationArgs{
		AppType:            appType,
		AppId:              appId,
		RespositoryVersion: respositoryVersion,
		BuildTime:          buildTime,
		DotEnvFilePath:     dotEnvFilePath,
		ConfigFilePath:     configFilePath,
	}, applicationFacade)

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
	respositoryVersion, _ := args.App.Metadata["respositoryVersion"].(string)
	fmt.Println(respositoryVersion)
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
		fmt.Println(string(c.Raw))
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
	appType, _ := args.App.Metadata["name"].(string)
	if _, err := gira.LoadApplicationConfig(configFilePath, dotEnvFilePath, appType, appId); err != nil {
		return err
	} else {
		// ctx := context.Background()
		// r, err := registry.NewConfigRegistry(ctx, c.Module.Etcd)
		// if err != nil {
		// 	return err
		// }
		// err = r.InitAsClient()
		// if err != nil {
		// 	return err
		// }
		return nil
	}
}

func reloadAction(args *cli.Context) error {
	return nil
}

func runAction(args *cli.Context) error {
	return nil
}
