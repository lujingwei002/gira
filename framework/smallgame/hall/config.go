package hall

import (
	"github.com/lujingwei002/gira"
	"gopkg.in/yaml.v2"
)

type Config struct {
	gira.Config
	Framework struct {
		Hall struct {
			GatewayReportInterval int64 `yaml:"gateway-report-interval"`
			// 玩家数据保存间隔
			BgSaveInterval       int64 `yaml:"bgsave-interval"`
			RequestBufferSize    int   `yaml:"request-buffer-size"`
			ResponseBufferSize   int   `yaml:"response-buffer-size"`
			PushBufferSize       int   `yaml:"push-buffer-size"`
			SessionActorBuffSize int   `yaml:"session-actor-buffer-size"`
		} `yaml:"hall"`
	} `yaml:"framework"`
}

func (c *Config) OnConfigLoad(config *gira.Config) error {
	c.Config = *config
	if err := yaml.Unmarshal(config.Raw, c); err != nil {
		return err
	}
	return nil
}
