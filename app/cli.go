package app

import (
	"os"
	"strconv"

	"github.com/lujingwei002/gira/log"

	"github.com/lujingwei002/gira"
	"github.com/urfave/cli/v2"
)

// 需要两个系参数 xx -id 1 start|stop|restart
func Cli(name string, buildVersion string, buildTime string, application gira.Application) error {
	log.Println("build version:", buildVersion)
	log.Println("build time:", buildTime)
	app := &cli.App{
		Name: "gira service",
		// Authors:     []*cli.Author{&cli.Author{Name: "lujingwei", Email: "lujingwei002@qq.com"}},
		Description: "gira service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "id",
				Usage:    "service id",
				Required: true,
			},
		},
		Action: runAction,
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
		},
	}
	log.Info(os.Args)
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("app run error %+v", err)
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
	log.Infof("%s %d starting...", appType, appId)
	runtime := newRuntime(ApplicationArgs{
		AppType:      appType,
		AppId:        appId,
		BuildVersion: buildVersion,
		BuildTime:    buildTime,
	}, application)
	return runtime.serve()
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
