package gate

import (
	"github.com/lujingwei002/gira"
	"gopkg.in/yaml.v2"
)

type Config struct {
	gira.Config
	Upstream struct {
		MinId int32  `yaml:"min-id"`
		MaxId int32  `yaml:"max-id"`
		Name  string `yaml:"name"`
	} `yaml:"upstream"`

	FrameWork struct {
		// 登录超时
		WaitLoginTimeout int64 `yaml:"wait-login-timeout"`
		// 最大会话数量
		MaxSessionCount int64 `yaml:"max-session-count"`
	} `yaml:"framework"`
}

func (c *Config) OnConfigLoad(config *gira.Config) error {
	c.Config = *config
	if err := yaml.Unmarshal(config.Raw, c); err != nil {
		return err
	}
	return nil
}
