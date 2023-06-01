package gen_application

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/lujingwei002/gira/log"

	"github.com/lujingwei002/gira/proj"
	"gopkg.in/yaml.v2"
)

var code = `
// Code generated by github.com/lujingwei002/gira. DO NOT EDIT.
// Code generated by github.com/lujingwei002/gira. DO NOT EDIT.
// Code generated by github.com/lujingwei002/gira. DO NOT EDIT.

package main 

import (
	"github.com/lujingwei002/gira/log"
	"github.com/lujingwei002/gira"
	gira_app "github.com/lujingwei002/gira/app"
	{{.ApplicationName}}_app "{{.ModuleName}}/{{.ApplicationName}}/app"
)

var buildVersion string
var buildTime string

// 检查是否满足接口
var _ = (gira.ApplicationFacade)(&{{.ApplicationName}}_app.Application{})
var _ = (gira.ResourceSource)(&{{.ApplicationName}}_app.Application{})

func main() {
	app := {{.ApplicationName}}_app.NewApplication()
	err := gira_app.Cli("{{.ApplicationName}}", buildVersion, buildTime, app)
	if err != nil {
		log.Info(err)
	}
}

`

type applications_file struct {
	Application []struct {
		Name string `yaml:"name"`
	} `yaml:"application"`
}

type gen_state struct {
	applicationFilePath string
	applicationsFile    applications_file
}

func capUpperString(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}

func parse(state *gen_state) error {
	if data, err := ioutil.ReadFile(state.applicationFilePath); err != nil {
		return err
	} else {
		if err := yaml.Unmarshal(data, &state.applicationsFile); err != nil {
			return err
		}
	}
	return nil
}

func genApplications(state *gen_state) error {
	log.Info("生成go文件")
	for _, v := range state.applicationsFile.Application {
		sb := strings.Builder{}
		if _, err := os.Stat(proj.Config.SrcGenApplicationDir); err != nil && os.IsNotExist(err) {
			if err := os.Mkdir(proj.Config.SrcGenApplicationDir, 0755); err != nil {
				return err
			}
		}
		srcGenApplicationDir := path.Join(proj.Config.SrcGenApplicationDir, v.Name)
		if _, err := os.Stat(srcGenApplicationDir); err != nil && os.IsNotExist(err) {
			if err := os.Mkdir(srcGenApplicationDir, 0755); err != nil {
				return err
			}
		}
		applicationPath := path.Join(srcGenApplicationDir, fmt.Sprintf("%s.gen.go", v.Name))
		file, err := os.OpenFile(applicationPath, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		file.Truncate(0)
		defer file.Close()

		tmpl, err := template.New("app").Parse(code)
		if err != nil {
			return err
		}
		params := map[string]string{
			"ModuleName":      proj.Config.Module,
			"ApplicationName": v.Name,
		}
		if err := tmpl.Execute(&sb, params); err != nil {
			return err
		}
		file.WriteString(sb.String())
	}
	return nil
}

func Gen() error {
	log.Info("===============gen app start.===============")
	// 初始化
	state := &gen_state{
		applicationFilePath: path.Join(proj.Config.DocDir, "application.yaml"),
	}
	if err := parse(state); err != nil {
		log.Info(err)
		return err
	}
	if err := genApplications(state); err != nil {
		log.Info(err)
		return err
	}
	log.Info("===============gen app finished===============")
	return nil
}
