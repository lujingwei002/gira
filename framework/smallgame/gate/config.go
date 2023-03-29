package gate

import (
	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/log"
	"gopkg.in/yaml.v2"
)

var Configs *Config

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
	Configs = c
	log.Infof("配置: %+v\n", string(c.Raw))
	return nil
}
