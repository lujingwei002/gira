package proj

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

const (
	project_config_file_name = "gira.yaml"
)

var (
	Config *ProjectConfig
)

type ProjectConfig struct {
	Version             string `yaml:"version"`
	Module              string `yaml:"module"`
	Env                 string //当前环境 local|dev|qa|prd
	Zone                string // 当前区 wc|qq|gf|review
	ProjectDir          string
	DotEnvFilePath      string // .env
	projectConfFilePath string // gira.yaml
	DocDir              string // doc
	ConfigDir           string // config
	EnvDir              string // env
	ConstDir            string //
	ResourceDir         string // resource
	GenDir              string // gen
	SrcTestDir          string // src/test
	GenModelDir         string // gen/model
	GenProtocolDir      string // gen/protocol
	SrcDir              string // src
	SrcGenDir           string // src/gen/
	SrcGenModelDir      string // src/gen/model/
	SrcGenProtocolDir   string // src/gen/protocol/
	SrcGenConstDir      string // src/gen/const/
	SrcGenResourceDir   string // src/gen/resource/
	ExcelDir            string // doc/resource/
	ConstDocFilePath    string // doc/const.yaml
	DocResourceFilePath string // doc/resource.yaml
	DocProtocolFilePath string // doc/protocol.yaml
	DocProtocolDir      string // doc/protocol/
	DocModelDir         string // doc/model/
}

func init() {
	Config = &ProjectConfig{}
	if err := Config.load(); err != nil {
		panic(err)
	}
}

func (p *ProjectConfig) load() error {
	// 初始化
	// 向上查找gira.yaml文件
	if workDir, err := os.Getwd(); err != nil {
		return err
	} else {
		dir := workDir
		for {
			projectFilePath := filepath.Join(dir, "gira.yaml")
			if _, err := os.Stat(projectFilePath); err == nil {
				p.ProjectDir = dir
				break
			}
			dir = filepath.Dir(dir)
			if dir == "/" || dir == "" {
				break
			}
		}
	}
	p.projectConfFilePath = path.Join(p.ProjectDir, project_config_file_name)
	p.EnvDir = path.Join(p.ProjectDir, "env")
	p.ConfigDir = path.Join(p.ProjectDir, "config")
	p.DotEnvFilePath = path.Join(p.ConfigDir, ".env")
	p.DocDir = path.Join(p.ProjectDir, "doc")
	p.ResourceDir = path.Join(p.ProjectDir, "resource")
	p.GenDir = path.Join(p.ProjectDir, "gen")
	p.SrcTestDir = path.Join(p.ProjectDir, "src", "test")
	p.GenModelDir = path.Join(p.ProjectDir, "gen", "model")
	p.SrcDir = path.Join(p.ProjectDir, "src")
	p.SrcGenDir = path.Join(p.ProjectDir, "src", "gen")
	p.SrcGenConstDir = path.Join(p.SrcGenDir, "const")
	p.SrcGenModelDir = path.Join(p.SrcGenDir, "model")
	p.SrcGenResourceDir = path.Join(p.SrcGenDir, "resource")
	p.SrcGenProtocolDir = path.Join(p.SrcGenDir, "protocol")
	p.ExcelDir = path.Join(p.DocDir, "resource")
	p.ConstDocFilePath = path.Join(p.DocDir, "const.yaml")
	p.DocResourceFilePath = path.Join(p.DocDir, "resource.yaml")
	p.DocProtocolFilePath = path.Join(p.DocDir, "protocol.yaml")
	p.DocProtocolDir = path.Join(p.DocDir, "protocol")
	p.DocModelDir = path.Join(p.DocDir, "model")
	p.GenProtocolDir = path.Join(p.GenDir, "protocol")
	if _, err := os.Stat(p.projectConfFilePath); err != nil && os.IsNotExist(err) {
		return err
	}
	data, err := ioutil.ReadFile(p.projectConfFilePath)
	if err != nil {
		return err
	}
	//使用yaml.Unmarshal将yaml文件中的信息反序列化给Config结构体
	if err := yaml.Unmarshal(data, p); err != nil {
		return err
	}
	return nil
}