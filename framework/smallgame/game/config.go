package game

import (
	"github.com/lujingwei002/gira"
	"gopkg.in/yaml.v2"
)

type Config struct {
	gira.Config
	Framework struct {
		GateReportInterval int64 `yaml:"gate-report-interval"`
		// 玩家数据保存间隔
		BgSaveInterval int64 `yaml:"bgsave-interval"`
	} `yaml:"framework"`
}

func (c *Config) OnConfigLoad(config *gira.Config) error {
	c.Config = *config
	if err := yaml.Unmarshal(config.Raw, c); err != nil {
		return err
	}
	return nil
}
