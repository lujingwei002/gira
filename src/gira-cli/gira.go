package main

import (
	"bytes"
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"os/exec"

	//https://cli.urfave.org/v2/examples/full-api-example/
	"github.com/urfave/cli/v2"

	"github.com/lujingwei/gira-cli/gen/gen_configs"
	"github.com/lujingwei/gira-cli/gen/gen_protocols"
	"github.com/lujingwei/gira-cli/gen/gen_services"
	"github.com/lujingwei/gira-cli/proj"
)

var state *proj.State

// TODO: afaf
func main() {
	app := &cli.App{
		Name: "Minigame",
		//app.Author = "lujingwei"
		///app.Email = "lujingwei@xx.org"
		Description: "Minigame",
		// flags
		Flags:  []cli.Flag{},
		Action: runAction,

		//log.SetFlags(log.LstdFlags | log.Lshortfile)
		Before: beforeAction,
		Commands: []*cli.Command{
			{
				Name:   "gen",
				Usage:  "gen code",
				Action: genAction,
				Before: beforeAction1,
			},
			{
				Name: "start",
				//Aliases: []string{"start"},
				Usage:  "Start minigame",
				Action: startAction,
			},
			{
				Name:   "status",
				Usage:  "Display status",
				Action: statusAction,
			},
			{
				Name:   "stop",
				Usage:  "Stop minigame",
				Action: stopAction,
			},
			{
				Name:   "restart",
				Usage:  "Restart minigame",
				Action: restartAction,
			},
			{
				Name:   "reload",
				Usage:  "Reload config",
				Action: reloadAction,
			},
		},
	}

	fmt.Println(os.Args)
	fmt.Println("start")
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("[main] startup server error %+v", err)
	}
}

func execCommand(name string, arg []string) (string, error) {
	cmd := exec.Command(name, arg...)
	//cmd.Stdin = strings.NewReader("")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return "", err
	}
	fmt.Println(out.String())
	return out.String(), nil
}

func beforeAction1(args *cli.Context) error {
	if err := state.LoadProjectConf(); err != nil {
		return err
	}
	return nil
}

func beforeAction(args *cli.Context) error {
	fmt.Println("beforeAction")
	var projectDir string
	if args.Args().Len() < 2 {
		projectDir = "."
	} else {
		projectDir = args.Args().Get(1)
	}
	fmt.Println(projectDir)

	// 当前目录
	err, s := proj.NewState(projectDir)
	if err != nil {
		return err
	}
	state = s
	return nil
}

// https://studygolang.com/articles/35019
func genAction(args *cli.Context) error {
	fmt.Println("1genAction")
	if err := gen_protocols.Gen(state); err != nil {
		return err
	}
	fmt.Println("gggggggggggggggg111")
	if err := gen_services.Gen(state); err != nil {
		return err
	}
	fmt.Println("gggggggggggggggg222")
	if err := gen_configs.Gen(state); err != nil {
		return err
	}
	return nil
}

func startAction(args *cli.Context) error {
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
	fmt.Println("fffffffffff1")
	fmt.Println("fffffffffff2")
	return nil
}
