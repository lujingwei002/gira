package proj

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

const (
	projectConfFileName = "gira.yaml"
)

type ProjectConf struct {
	Version string `yaml:version`
}

type State struct {
	projectConf         ProjectConf
	projectDir          string
	projectConfFilePath string
	DocDir              string
	GenDir              string
	SrcDir              string
	GenSrcDir           string
}

var GlobalState *State

func NewState(projectDir string) (error, *State) {
	if GlobalState != nil {
		return nil, GlobalState
	}
	state := &State{}
	if dir, err := filepath.Abs(projectDir); err != nil {
		return err, nil
	} else {
		state.projectDir = dir
	}
	/* if exePath, err := exec.LookPath(os.Args[0]); err != nil {
		return err, nil
	} else {
		if absExePath, err := filepath.Abs(exePath); err != nil {
			return err, nil
		} else {
			state.projectDir = path.Dir(path.Dir(absExePath))
		}
	} */
	state.projectConfFilePath = path.Join(state.projectDir, projectConfFileName)
	state.DocDir = path.Join(state.projectDir, "doc")
	state.GenDir = path.Join(state.projectDir, "gen")
	state.SrcDir = path.Join(state.projectDir, "src")
	state.GenSrcDir = path.Join(state.projectDir, "src", "gen")
	fmt.Println(state.projectDir)
	return nil, state
}

func (p *ProjectConf) read(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	if err := yaml.Unmarshal(data, p); err != nil { //使用yaml.Unmarshal将yaml文件中的信息反序列化给Config结构体
		fmt.Printf("err: %v\n", err)
		return err
	}
	return nil
}

func (s *State) LoadProjectConf() error {
	if _, err := os.Stat(s.projectConfFilePath); os.IsNotExist(err) {
		return fmt.Errorf("%s not found", s.projectConfFilePath)
	}
	if err := s.projectConf.read(s.projectConfFilePath); err != nil {
		return err
	}
	return nil
}
