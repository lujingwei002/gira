package gira

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	yaml "gopkg.in/yaml.v3"
)

type JwtConfig struct {
	Secret            string `yaml:"secret"`
	Expiretime        int64  `yaml:"expiretime"`
	RefreshExpiretime int64  `yaml:"refresh-expiretime"`
}

type GameDbConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
}

type AccountCacheConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}

type AccountDbConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
}
type StatDbConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
}

type EtcdConfig struct {
	Endpoints []struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"endpoints"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	DialTimeout  int    `yaml:"dial-timeout"`
	LeaseTimeout int64  `yaml:"lease-timeout"`
	Address      string `yaml:"address"`
	Advertise    []struct {
		Name  string `yaml:"name"`
		Value string `yaml:"value"`
	} `yaml:"advertise"`
}

type HttpConfig struct {
	Addr         string `yaml:"addr"`
	ReadTimeout  int64  `yaml:"read-timeout"`
	WriteTimeout int64  `yaml:"write-timeout"`
}

type GateConfig struct {
	Bind    string `yaml:"bind"`
	Address string `yaml:"address"`
}

type TestSdkConfig struct {
	Secret string `yaml:"secret"`
}

type SdkConfig struct {
	Test *TestSdkConfig `yaml:"test"`
}

type GrpcConfig struct {
	Address string `yaml:"address"`
}

type Config struct {
	Raw          []byte
	Thread       int                 `yaml:"thread"`
	Name         string              `yaml:"name"`
	Http         *HttpConfig         `yaml:"http,omitempty"`
	GameDb       *GameDbConfig       `yaml:"gamedb"`
	AccountDb    *AccountDbConfig    `yaml:"accountdb"`
	StatDb       *StatDbConfig       `yaml:"statdb"`
	AccountCache *AccountCacheConfig `yaml:"account-cache"`
	Etcd         *EtcdConfig         `yaml:"etcd"`
	Grpc         *GrpcConfig         `yaml:"grpc"`
	Sdk          *SdkConfig          `yaml:"sdk"`
	Jwt          *JwtConfig          `yaml:"jwt"`
	Gate         *GateConfig         `yaml:"gate"`
}

func NewConfig() *Config {
	return &Config{}
}

type ConfigHandler interface {
	LoadConfig(c *Config) error
}

func (c *Config) Unmarshal(data []byte) error {
	// 解析yaml
	if err := yaml.Unmarshal(data, c); err != nil {
		log.Println(string(data))
		return err
	}
	c.Raw = data
	// log.Printf("配置: %+v\n", c)
	// log.Printf("GameDb配置: %+v\n", c.GameDb)
	// log.Printf("Etcd配置: %+v\n", c.Etcd)
	// log.Printf("Grpc配置: %+v\n", c.Grpc)
	return nil
}

func (c *Config) Parse(facade ApplicationFacade, dir string, zone string, env string, name string) error {
	if data, err := c.Read(facade, dir, zone, env, name); err != nil {
		return err
	} else {
		return c.Unmarshal(data)
	}
}

func applicationFieldMapFunc(facade ApplicationFacade, env map[string]interface{}, key string) interface{} {
	appName := fmt.Sprintf("%s-%d", facade.GetName(), facade.GetId())
	if v, ok := env["application"]; ok {
		if applications, ok := v.(map[string]interface{}); ok {
			if v, ok := applications[appName]; ok {
				if application, ok := v.(map[string]interface{}); ok {
					if value, ok := application[key]; ok {
						return value
					}
				}
			}
		}
	}
	return "none"
}

func otherApplicationFieldMapFunc(facade ApplicationFacade, name string, id int32, env map[string]interface{}, key string) interface{} {
	appName := fmt.Sprintf("%s-%d", name, facade.GetId()+id)
	if v, ok := env["application"]; ok {
		if applications, ok := v.(map[string]interface{}); ok {
			if v, ok := applications[appName]; ok {
				if application, ok := v.(map[string]interface{}); ok {
					if value, ok := application[key]; ok {
						return value
					}
				}
			}
		}
	}
	return "none"
}

func hostFieldMapFunc(facade ApplicationFacade, env map[string]interface{}, key string) interface{} {
	appName := fmt.Sprintf("%s-%d", facade.GetName(), facade.GetId())
	var hostName string
	var hostFound bool = false
	if v, ok := env["application"]; ok {
		if applications, ok := v.(map[string]interface{}); ok {
			if v, ok := applications[appName]; ok {
				if application, ok := v.(map[string]interface{}); ok {
					hostName = application["host"].(string)
					hostFound = true
				}
			}
		}
	}
	if !hostFound {
		return "none"
	}
	if v, ok := env["host"]; ok {
		if hosts, ok := v.(map[string]interface{}); ok {
			if v, ok := hosts[hostName]; ok {
				if host, ok := v.(map[string]interface{}); ok {
					if value, ok := host[key]; ok {
						return value
					}
				}
			}
		}
	}
	return "none"
}

func otherHostFieldMapFunc(facade ApplicationFacade, name string, id int32, env map[string]interface{}, key string) interface{} {
	appName := fmt.Sprintf("%s-%d", name, facade.GetId()+id)
	log.Println("bb", appName, key)
	var hostName string
	var hostFound bool = false
	if v, ok := env["application"]; ok {
		if applications, ok := v.(map[string]interface{}); ok {
			if v, ok := applications[appName]; ok {
				if application, ok := v.(map[string]interface{}); ok {
					hostName = application["host"].(string)
					hostFound = true
				}
			}
		}
	}
	if !hostFound {
		return "none"
	}
	if v, ok := env["host"]; ok {
		if hosts, ok := v.(map[string]interface{}); ok {
			if v, ok := hosts[hostName]; ok {
				if host, ok := v.(map[string]interface{}); ok {
					if value, ok := host[key]; ok {
						return value
					}
				}
			}
		}
	}
	return "none"
}

// 加载配置
// 根据环境，区名，服务名， 服务组合配置文件路径，规则是config/app/<<name>>.yaml
func (c *Config) Read(facade ApplicationFacade, dir string, zone string, env string, name string) ([]byte, error) {
	var configFilePath = filepath.Join(dir, fmt.Sprintf("%s.yaml", name))
	sb := strings.Builder{}
	// 预处理
	if err := c.preprocess(&sb, configFilePath); err != nil {
		return nil, err
	}
	// log.Printf("配置预处理后\n%v\n", sb.String())
	// 读环境变量
	envFilePath := path.Join(dir, "env", env, "env.yaml")
	envData, err := c.readEnv(envFilePath)
	if err != nil {
		return nil, err
	}

	funcMap := template.FuncMap{
		// The name "title" is what the function will be called in the template text.
		"application_field": func(env map[string]interface{}, key string) interface{} {
			return applicationFieldMapFunc(facade, env, key)
		},
		"other_application_field": func(name string, id int32, env map[string]interface{}, key string) interface{} {
			return otherApplicationFieldMapFunc(facade, name, id, env, key)
		},
		"host_field": func(env map[string]interface{}, key string) interface{} {
			return hostFieldMapFunc(facade, env, key)
		},

		"other_host_field": func(name string, id int32, env map[string]interface{}, key string) interface{} {
			return otherHostFieldMapFunc(facade, name, id, env, key)
		},

		"application_id": func() interface{} {
			return facade.GetId()
		},

		"other_application_id": func(id int32) interface{} {
			return facade.GetId() + id
		},
		"application_name": func() interface{} {
			return facade.GetName()
		},
	}

	// 替换环境变量
	t := template.New("config").Delims("<<", ">>")
	t.Funcs(funcMap)
	t, err = t.Parse(sb.String())
	if err != nil {
		return nil, err
	}
	out := strings.Builder{}
	t.Execute(&out, envData)
	// log.Printf("替换环境变量后\n%v\n", out.String())
	return []byte(out.String()), nil
}

// 执行include指令
func (c *Config) readEnv(filePath string) (map[string]interface{}, error) {
	envData := make(map[string]interface{})
	sb := strings.Builder{}
	if _, err := os.Stat(filePath); err != nil {
		return envData, err
	}
	if err := c.preprocess(&sb, filePath); err != nil {
		return envData, err
	}
	if err := yaml.Unmarshal([]byte(sb.String()), envData); err != nil {
		return envData, err
	}
	return envData, nil
}

const (
	instruction_none    = iota
	instruction_include = iota
)

// 执行include指令
func (c *Config) preprocess(sb *strings.Builder, filePath string) error {
	dir := path.Dir(filePath)

	lines, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	scanner := bufio.NewScanner(lines)
	scanner.Split(bufio.ScanLines)
	instruction := instruction_none
	for scanner.Scan() {
		line := scanner.Text()
		if instruction == instruction_none {
			if strings.HasPrefix(line, `$include:`) {
				instruction = instruction_include
			} else {
				sb.WriteString(line)
				sb.WriteString("\n")
			}
		}
		if instruction == instruction_include {
			var includeFilePath string
			var found bool = false
			if strings.HasPrefix(line, `$include:`) {
				pats := strings.SplitN(line, ":", 2)
				if len(pats) == 2 {
					includeFilePath = strings.TrimSpace(pats[1])
				}
				found = true
			} else {
				str := strings.TrimSpace(line)
				if strings.HasPrefix(str, "-") {
					str = strings.Replace(line, `-`, "", 1)
					includeFilePath = strings.TrimSpace(str)
					found = true
				}
			}
			if len(includeFilePath) > 0 {
				if err := c.preprocess(sb, path.Join(dir, includeFilePath)); err != nil {
					return err
				}
			}
			if !found {
				instruction = instruction_none
				sb.WriteString(line)
				sb.WriteString("\n")
			}
		}
	}
	lines.Close()
	return nil
}
