package app

import (
	"fmt"
	"os"
	"strconv"

	"github.com/lujingwei002/gira/log"
	"github.com/lujingwei002/gira/proj"

	"github.com/lujingwei002/gira"
	"github.com/urfave/cli/v2"
)

// 需要两个系参数 xx -id 1 start|stop|restart
func Cli(name string, buildVersion string, buildTime string, application gira.Application) error {

	app := &cli.App{
		Name: "gira service",
		// Authors:     []*cli.Author{&cli.Author{Name: "lujingwei", Email: "lujingwei002@qq.com"}},
		Description: "gira service",
		Flags:       []cli.Flag{},
		Action:      runAction,
		Metadata: map[string]interface{}{
			"application":  application,
			"name":         name,
			"buildTime":    buildTime,
			"buildVersion": buildVersion,
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
				},
			},
			{
				Name:   "status",
				Usage:  "Display status",
				Action: statusAction,
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
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Errorw("app run fail", "error", err)
	}
	return nil
}

func Start(application gira.Application, appId int32, appType string) error {
	runtime := newRuntime(ApplicationArgs{
		AppType: appType,
		AppId:   appId,
	}, application)
	return runtime.start()
}

// 启动应用
func startAction(args *cli.Context) error {
	appId := int32(args.Int("id"))
	application, _ := args.App.Metadata["application"].(gira.Application)
	appType, _ := args.App.Metadata["name"].(string)
	buildVersion, _ := args.App.Metadata["buildVersion"].(string)
	var buildTime int64
	if v, ok := args.App.Metadata["buildTime"].(string); ok {
		if t, err := strconv.Atoi(v); err != nil {
			buildTime = 0
		} else {
			buildTime = int64(t)
		}
	}
	log.Println("build version:", buildVersion)
	log.Println("build time:", buildTime)
	log.Infof("%s %d starting...", appType, appId)
	runtime := newRuntime(ApplicationArgs{
		AppType:      appType,
		AppId:        appId,
		BuildVersion: buildVersion,
		BuildTime:    buildTime,
	}, application)
	if err := runtime.start(); err != nil {
		return err
	}
	if err := runtime.wait(); err != nil {
		return err
	}
	return nil
}

// 打印应用构建版本
func versionAction(args *cli.Context) error {
	buildVersion, _ := args.App.Metadata["buildVersion"].(string)
	fmt.Println(buildVersion)
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
	if _, err := gira.LoadConfig(proj.Config.ConfigDir, proj.Config.EnvDir, appType, appId); err != nil {
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
	appType, _ := args.App.Metadata["name"].(string)
	if c, err := gira.LoadConfig(proj.Config.ConfigDir, proj.Config.EnvDir, appType, appId); err != nil {
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
	/*projectDir, err := getProjectDir()
	if err != nil {
		return err
	}
	configPath := fmt.Sprintf("%s/supervisord/supervisord.conf", projectDir)
	execCommand("supervisorctl", []string{"-c", configPath, "status"})*/
	return nil
}

func reloadAction(args *cli.Context) error {
	return nil
}

func runAction(args *cli.Context) error {
	return nil
}
