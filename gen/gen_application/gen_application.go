package gen_application

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/lujingwei002/gira/log"
	"github.com/lujingwei002/gira/proj"
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

type application struct {
	ApplicationName string
	ModuleName      string
}
type gen_state struct {
	applications []application
}

type Parser interface {
	parse(state *gen_state) error
}

func gen(state *gen_state) error {
	if _, err := os.Stat(proj.Config.SrcGenApplicationDir); err != nil && os.IsNotExist(err) {
		if err := os.Mkdir(proj.Config.SrcGenApplicationDir, 0755); err != nil {
			return err
		}
	}
	for _, v := range state.applications {
		sb := strings.Builder{}
		srcGenApplicationDir := path.Join(proj.Config.SrcGenApplicationDir, v.ApplicationName)
		if _, err := os.Stat(srcGenApplicationDir); err != nil && os.IsNotExist(err) {
			if err := os.Mkdir(srcGenApplicationDir, 0755); err != nil {
				return err
			}
		}
		applicationPath := path.Join(srcGenApplicationDir, fmt.Sprintf("%s.gen.go", v.ApplicationName))
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
		if err := tmpl.Execute(&sb, v); err != nil {
			return err
		}
		log.Printf("gen application %v", v.ApplicationName)
		file.WriteString(sb.String())
	}
	return nil
}

func Gen() error {
	log.Info("===============gen app start===============")
	// 初始化
	state := &gen_state{}
	if true {

	}
	var p Parser
	if true {
		p = &golang_parser{}
	} else {
		p = &yaml_parser{}
	}
	if err := p.parse(state); err != nil {
		log.Info(err)
		return err
	}
	if err := gen(state); err != nil {
		log.Info(err)
		return err
	}
	log.Info("===============gen app finished===============")
	return nil
}
