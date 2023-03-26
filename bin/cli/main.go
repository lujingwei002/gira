package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	_ "net/http/pprof"
	"os"
	"os/exec"

	//https://cli.urfave.org/v2/examples/full-api-example/
	"github.com/lujingwei002/gira/log"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"

	"github.com/lujingwei002/gira/gen/gen_const"
	"github.com/lujingwei002/gira/gen/gen_model"
	"github.com/lujingwei002/gira/gen/gen_protocols"
	"github.com/lujingwei002/gira/gen/gen_resources"
	"github.com/lujingwei002/gira/gen/gen_services"
	"github.com/lujingwei002/gira/macro"
	"github.com/lujingwei002/gira/proj"
)

func main() {
	app := &cli.App{
		Name: "gira-cli",
		Authors: []*cli.Author{{
			Name:  "lujingwei",
			Email: "lujingwei@xx.org",
		}},
		Description: "gira-cli",
		Flags:       []cli.Flag{},
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
				Name:   "resource",
				Usage:  "resource [compress|migrate]",
				Before: beforeAction1,
				Subcommands: []*cli.Command{
					{
						Name:   "compress",
						Usage:  "compress yaml to binary",
						Action: resourceCompressAction,
					},
					{
						Name:   "push",
						Usage:  "push table to database",
						Action: resourcePushAction,
					},
					{
						Name:   "reload",
						Usage:  "reload",
						Action: resourceReloadAction,
					},
				},
			},
			{
				Name:   "env",
				Usage:  "env",
				Before: beforeAction1,
				Subcommands: []*cli.Command{
					{
						Name:   "list",
						Usage:  "list env name",
						Action: envListAction,
					},
					{
						Name:   "switch",
						Usage:  "switch env name",
						Action: envSwitchAction,
					},
				},
			},
			{
				Name:   "macro",
				Usage:  "macro code",
				Action: macroAction,
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
	if err := app.Run(os.Args); err != nil {
		log.Fatal("error", err)
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
	return nil
}

func command(bin string, arg ...string) error {
	if _, err := os.Stat(bin); err != nil && os.IsNotExist(err) {
		return err
	}
	log.Infow("run command", "command", fmt.Sprint(arg))
	cmd := exec.Command(bin, arg...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	log.Info("output:\n", string(output))
	return nil
}

func resourceCompressAction(args *cli.Context) error {
	bin := "bin/resource"
	return command(bin, "compress")
}

func resourcePushAction(args *cli.Context) error {
	return nil
	// config := proj.Config.ResourceDb
	// uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", config.User, config.Password, config.Host, config.Port, config.Db)
	// bin := "bin/resource"
	// return command(bin, "push", uri)
}

func resourceReloadAction(args *cli.Context) error {
	return nil
}

// 切换环境
func envSwitchAction(args *cli.Context) error {
	if args.NArg() < 1 {
		cli.ShowAppHelp(args)
		return nil
	}
	expectName := args.Args().Get(0)
	envDir := proj.Config.EnvDir
	found := false
	if d, err := os.Open(envDir); err != nil {
		return err
	} else {
		defer d.Close()
		if files, err := d.ReadDir(0); err != nil {
			return err
		} else {
			for _, file := range files {
				if file.IsDir() {
					envName := file.Name()
					if envName == expectName {
						found = true
					}
				}
			}
		}
	}
	if !found {
		log.Println("env not found, you can select one of below env")
		return envListAction(args)
	}
	if data, err := ioutil.ReadFile(proj.Config.DotEnvFilePath); err != nil {
		return err
	} else {
		var result map[string]interface{}
		if err := yaml.Unmarshal(data, &result); err != nil {
			return err
		}
		result["env"] = expectName
		if data, err := yaml.Marshal(result); err != nil {
			return err
		} else {
			if err := ioutil.WriteFile(proj.Config.DotEnvFilePath, data, 0644); err != nil {
				return err
			}
		}
	}
	log.Printf("switch %s success", expectName)
	return nil
}

// 环境列表
func envListAction(args *cli.Context) error {
	envDir := proj.Config.EnvDir
	if d, err := os.Open(envDir); err != nil {
		return err
	} else {
		defer d.Close()
		if files, err := d.ReadDir(0); err != nil {
			return err
		} else {
			for _, file := range files {
				if file.IsDir() {
					envName := file.Name()
					if envName == proj.Config.Env {
						log.Info("*", file.Name())
					} else {
						log.Info(file.Name())
					}
				}
			}
		}
	}
	return nil
}
func beforeAction(args *cli.Context) error {
	return nil
}

func genAction(args *cli.Context) error {
	log.Println("genAction")
	if err := gen_protocols.Gen(); err != nil {
		return err
	}
	if err := gen_services.Gen(); err != nil {
		return err
	}
	if err := gen_resources.Gen(); err != nil {
		return err
	}
	if err := gen_const.Gen(); err != nil {
		return err
	}
	if err := gen_model.Gen(); err != nil {
		return err
	}
	return nil
}

func resourceAction(args *cli.Context) error {
	log.Println("resourceAction")
	return nil
}

func macroAction(args *cli.Context) error {
	log.Println("macroAction")
	if err := macro.Gen(); err != nil {
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