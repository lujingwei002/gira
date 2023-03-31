package gate

import (
	"github.com/lujingwei002/gira"
	"gopkg.in/yaml.v2"
)

type Config struct {
	gira.Config
	Upstream struct {
		Id   int32  `yaml:"id"`
		Name string `yaml:"name"`
	} `yaml:"upstream"`
}

func (c *Config) OnConfigLoad(config *gira.Config) error {
	c.Config = *config
	if err := yaml.Unmarshal(config.Raw, c); err != nil {
		return err
	}
	return nil
}
