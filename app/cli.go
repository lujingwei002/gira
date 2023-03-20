package app

import (
	"log"
	"os"

	"github.com/lujingwei/gira"
	"github.com/urfave/cli/v2"
)

// 需要两个系参数 xx -id 1 start|stop|restart
func Cli(name string, facade gira.ApplicationFacade) error {
	app := &cli.App{
		Name: "gira service",
		//app.Author = "lujingwei"
		///app.Email = "lujingwei@xx.org"
		Description: "gira service",
		// flags
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "zone",
				Usage:    "zone, gf|wc|qq|quick",
				Required: true,
				Aliases:  []string{"z"},
			},
			&cli.StringFlag{
				Name:     "env",
				Usage:    "environment, local|dev|qa|prd",
				Required: true,
				Aliases:  []string{"e"},
			},
			&cli.StringFlag{
				Name:     "id",
				Usage:    "service id",
				Required: true,
			},
		},
		Action: runAction,
		Metadata: map[string]interface{}{
			"facade": facade,
			"name":   name,
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
	log.Println(os.Args)
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("app run error %+v", err)
	}
	return nil
}

func Start(facade gira.ApplicationFacade, zone string, env string, id int32, name string) error {
	// log.Println(name, zone, env, id, "start")
	application := newApplication(ApplicationArgs{
		Name: name,
		Id:   id,
		Env:  env,
		Zone: zone,
	}, facade)
	return application.start()
}

func startAction(args *cli.Context) error {
	zone := args.String("zone")
	env := args.String("env")
	id := int32(args.Int("id"))
	facade, _ := args.App.Metadata["facade"].(gira.ApplicationFacade)
	name, _ := args.App.Metadata["name"].(string)

	log.Println(name, zone, env, id, "start")
	application := newApplication(ApplicationArgs{
		Name: name,
		Id:   id,
		Env:  env,
		Zone: zone,
	}, facade)
	err := application.forver()
	return err
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
