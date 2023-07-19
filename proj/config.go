package proj

import (
	"fmt"
	"path"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/config"
)

// 读取应该配置
func LoadCliConfig() (*gira.Config, error) {
	configFilePath := path.Join(Config.ConfigDir, "cli.yaml")
	dotEnvFilePath := path.Join(Config.EnvDir, ".env")
	return config.Load(Config.ProjectDir, configFilePath, dotEnvFilePath, "cli", 0)
}

// 读取测试配置
func LoadTestConfig() (*gira.Config, error) {
	configFilePath := path.Join(Config.ConfigDir, "test.yaml")
	dotEnvFilePath := path.Join(Config.EnvDir, ".env")
	return config.Load(Config.ProjectDir, configFilePath, dotEnvFilePath, "test", 0)
}

// 读取应用配置
func LoadApplicationConfig(appType string, appId int32) (*gira.Config, error) {
	configFilePath := path.Join(Config.ConfigDir, fmt.Sprintf("%s.yaml", appType))
	dotEnvFilePath := path.Join(Config.EnvDir, ".env")
	if c, err := config.Load(Config.ProjectDir, configFilePath, dotEnvFilePath, appType, appId); err != nil {
		return nil, err
	} else {

		return c, nil
	}
}
