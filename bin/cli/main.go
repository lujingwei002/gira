package main

import (
	"bufio"
	"fmt"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	//https://cli.urfave.org/v2/examples/full-api-example/
	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/log"
	"github.com/urfave/cli/v2"

	"github.com/lujingwei002/gira/gen/gen_application"
	"github.com/lujingwei002/gira/gen/gen_const"
	"github.com/lujingwei002/gira/gen/gen_macro"
	"github.com/lujingwei002/gira/gen/gen_model"
	"github.com/lujingwei002/gira/gen/gen_protocols"
	"github.com/lujingwei002/gira/gen/gen_resources"
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
				Usage:  "gen [resource|model|protocol|const]",
				Before: beforeAction1,
				Subcommands: []*cli.Command{
					{
						Name:   "all",
						Usage:  "gen all",
						Action: genAllAction,
					},
					{
						Name:   "resource",
						Usage:  "gen resource",
						Action: genResourceAction,
					},
					{
						Name:   "model",
						Usage:  "gen model",
						Action: genModelAction,
					},
					{
						Name:   "protocol",
						Usage:  "gen protocol",
						Action: genProtocolAction,
					},
					{
						Name:   "const",
						Usage:  "gen const",
						Action: genConstAction,
					},
					{
						Name:   "app",
						Usage:  "gen app",
						Action: genAppAction,
					},
				},
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
			{
				Name:   "run",
				Usage:  "run command",
				Action: runAction,
			},
			{
				Name:   "build",
				Usage:  "build command",
				Action: buildAction,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Info(err)
	}
}

func execCommandLine(line string) error {
	lastWd, err := os.Getwd()
	if err != nil {
		return err
	}
	defer func() {
		os.Chdir(lastWd)
	}()
	arr := strings.Split(line, ";")
	for _, v := range arr {
		pats := strings.Fields(v)
		name := pats[0]
		args := pats[1:]
		switch name {
		case "cd":
			if len(args) > 0 {
				os.Chdir(args[0])
			} else {
				os.Chdir("")
			}
			log.Printf("[OK] %s", v)
		default:
			log.Printf("%s", v)
			if err := execCommand(name, args); err != nil {
				log.Printf("[FAIL] %s", v)
				return err
			} else {
				log.Printf("[OK] %s", v)
			}
		}
	}
	return nil
}

func execCommand(name string, arg []string) error {
	cmd := exec.Command(name, arg...)
	// 获取命令的标准输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	// 启动命令
	if err := cmd.Start(); err != nil {
		return err
	}
	// 创建一个channel，用于接收信号
	c := make(chan os.Signal, 1)
	// 监听SIGINT信号
	signal.Notify(c, os.Interrupt, syscall.SIGINT)
	defer func() {
		signal.Reset(os.Interrupt, syscall.SIGINT)
	}()
	// 创建一个 Scanner 对象，对命令的标准输出和标准错误输出进行扫描
	scanner1 := bufio.NewScanner(stdout)
	go func() {
		for scanner1.Scan() {
			// 输出命令的标准输出
			log.Println(scanner1.Text())
		}
	}()
	scanner2 := bufio.NewScanner(stderr)
	go func() {
		for scanner2.Scan() {
			// 输出命令的标准错误输出
			fmt.Fprintln(os.Stderr, scanner2.Text())
		}
	}()
	go func() {
		// 等待信号
		<-c
	}()
	// 等待命令执行完成
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
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
	if config, err := gira.LoadCliConfig(proj.Config.ConfigDir, proj.Config.EnvDir); err != nil {
		return err
	} else {
		dbConfig := config.ResourceDb
		uri := dbConfig.Uri()
		bin := "bin/resource"
		return command(bin, "push", uri)
	}
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
	os.Remove(filepath.Join(proj.Config.EnvDir, ".env"))
	os.Remove(filepath.Join(proj.Config.EnvDir, ".config.yaml"))
	if err := os.Symlink(filepath.Join(proj.Config.EnvDir, expectName, ".env"), filepath.Join(proj.Config.EnvDir, ".env")); err != nil {
		return err
	}
	if err := os.Symlink(filepath.Join(proj.Config.EnvDir, expectName, ".config.yaml"), filepath.Join(proj.Config.EnvDir, ".config.yaml")); err != nil {
		return err
	}
	log.Printf("switch %s success", expectName)
	return envListAction(args)
}

// 环境列表
func envListAction(args *cli.Context) error {
	if config, err := gira.LoadCliConfig(proj.Config.ConfigDir, proj.Config.EnvDir); err != nil {
		return err
	} else {
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
						if envName == config.Env {
							log.Info("*", file.Name())
						} else {
							log.Info(file.Name())
						}
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

func genResourceAction(args *cli.Context) error {
	if err := gen_resources.Gen(gen_resources.Config{
		Force: true,
	}); err != nil {
		return err
	}
	return nil
}

func genProtocolAction(args *cli.Context) error {
	if err := gen_protocols.Gen(); err != nil {
		return err
	}
	return nil
}

func genModelAction(args *cli.Context) error {
	if err := gen_model.Gen(); err != nil {
		return err
	}
	return nil
}

func genAppAction(args *cli.Context) error {
	if err := gen_application.Gen(); err != nil {
		return err
	}
	return nil
}

func genConstAction(args *cli.Context) error {
	if err := gen_const.Gen(); err != nil {
		return err
	}
	return nil
}

func genAllAction(args *cli.Context) error {
	if err := gen_protocols.Gen(); err != nil {
		return err
	}
	if err := gen_application.Gen(); err != nil {
		return err
	}
	if err := gen_resources.Gen(gen_resources.Config{
		Force: true,
	}); err != nil {
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
	if err := gen_macro.Gen(&gen_macro.Config{}); err != nil {
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

func runAction(c *cli.Context) error {
	if c.Args().Len() <= 0 {
		return nil
	}
	name := c.Args().First()
	args := c.Args().Tail()
	log.Println(args)
	if arr, ok := proj.Config.Run[name]; !ok {
		return nil
	} else {
		for _, line := range arr {
			// 替换命令中的变量
			for k, v := range args {
				line = strings.Replace(line, fmt.Sprintf("$(%d)", k+1), v, 1)
			}
			if err := execCommandLine(line); err != nil {
				return err
			}
		}
	}
	return nil
}

func buildAction(c *cli.Context) error {
	if c.Args().Len() <= 0 {
		return nil
	}
	name := c.Args().First()
	var buildFunc func(target string) error
	buildFunc = func(target string) error {
		if build, ok := proj.Config.Build[target]; !ok {
			return nil
		} else {
			if len(build.Dependency) > 0 {
				for _, v := range build.Dependency {
					if err := buildFunc(v); err != nil {
						log.Printf("[FAIL] build %s\n", v)
						return err
					} else {
						log.Printf("[OK] build %s\n", v)
					}
				}
			}
			log.Printf(build.Description)
			for _, v := range build.Run {
				if err := execCommandLine(v); err != nil {
					return err
				}
			}
			return nil
		}
	}
	return buildFunc(name)
}
