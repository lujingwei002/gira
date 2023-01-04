package gira

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v3"
)

type ProjectConf struct {
	Version string `yaml:version`
}

func (p *ProjectConf) Read(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	if err := yaml.Unmarshal(data, p); err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	return nil
}
