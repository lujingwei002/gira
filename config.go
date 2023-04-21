package gira

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/joho/godotenv"

	yaml "gopkg.in/yaml.v3"
)

const (
	config_instruction_type_none = iota
	config_instruction_type_include
)

// 日志配置
type LogConfig struct {
	MaxSize    int    `yaml:"max-size"`
	MaxBackups int    `yaml:"max-backups"`
	MaxAge     int    `yaml:"max-age"`
	Compress   bool   `yaml:"compress"`
	Level      string `yaml:"level"`
}

// jwt配置
type JwtConfig struct {
	Secret            string `yaml:"secret"`
	Expiretime        int64  `yaml:"expiretime"`
	RefreshExpiretime int64  `yaml:"refresh-expiretime"`
}

// 游戏数据库配置
type GameDbConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
}

// 账号数据库配置
type AccountCacheConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}

// 账号数据库配置
type AccountDbConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
}

// 资源数据库配置
type ResourceDbConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
}

func (self ResourceDbConfig) Uri() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", self.User, self.Password, self.Host, self.Port, self.Db)
}

func (self GameDbConfig) Uri() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", self.User, self.Password, self.Host, self.Port, self.Db)
}

type AdminDbConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
}

// 状态数据库配置
type StatDbConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
}

// etcd配置
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

// http模块配置
type HttpConfig struct {
	Addr         string `yaml:"addr"`
	ReadTimeout  int64  `yaml:"read-timeout"`
	WriteTimeout int64  `yaml:"write-timeout"`
}

// 网关模块配置
type GatewayConfig struct {
	Bind    string `yaml:"bind"`
	Address string `yaml:"address"`
	Debug   bool   `yaml:"debug"`
}

type TestSdkConfig struct {
	Secret string `yaml:"secret"`
}
type PwdSdkConfig struct {
	Secret string `yaml:"secret"`
}
type SdkConfig struct {
	Test *TestSdkConfig `yaml:"test"`
	Pwd  *PwdSdkConfig  `yaml:"pwd"`
}

type GrpcConfig struct {
	Address string `yaml:"address"`
}

type PprofConfig struct {
	Port int    `yaml:"port"`
	Bind string `yaml:"bind"`
}

type Config struct {
	Raw    []byte
	Thread int `yaml:"thread"`
	Env    string
	Zone   string
	Log    *LogConfig  `yaml:"log"`
	Pprof  PprofConfig `yaml:"pprof"`
	Module struct {
		ResourceDb   *ResourceDbConfig   `yaml:"resourcedb"`
		GameDb       *GameDbConfig       `yaml:"gamedb"`
		AccountDb    *AccountDbConfig    `yaml:"accountdb"`
		StatDb       *StatDbConfig       `yaml:"statdb"`
		AdminDb      *AdminDbConfig      `yaml:"admindb"`
		AccountCache *AccountCacheConfig `yaml:"account-cache"`
		Http         *HttpConfig         `yaml:"http,omitempty"`
		Etcd         *EtcdConfig         `yaml:"etcd"`
		Grpc         *GrpcConfig         `yaml:"grpc"`
		Sdk          *SdkConfig          `yaml:"sdk"`
		Jwt          *JwtConfig          `yaml:"jwt"`
		Gateway      *GatewayConfig      `yaml:"gateway"`
		Admin        *AdminConfig        `yaml:"admin"`
	} `yaml:"module"`
}

type AdminConfig struct {
	None string `yaml:"none"`
}

type dot_env_config struct {
	Env  string `yaml:"env"`
	Zone string `yaml:"zone"`
}

type config_reader struct {
	appId   int32
	appType string
	appName string
	zone    string
	env     string
}

// 读取应该配置
func LoadConfig(configDir string, envDir string, appType string, appId int32) (*Config, error) {
	c := &Config{}
	appName := fmt.Sprintf("%s_%d", appType, appId)
	reader := config_reader{
		appType: appType,
		appId:   appId,
		appName: appName,
	}
	if data, err := reader.read(configDir, envDir); err != nil {
		return nil, err
	} else {
		if err := c.unmarshal(data); err != nil {
			return nil, err
		}
	}
	c.Env = reader.env
	c.Zone = reader.zone
	return c, nil
}

func (c *Config) unmarshal(data []byte) error {
	// 解析yaml
	if err := yaml.Unmarshal(data, c); err != nil {
		return err
	}
	c.Raw = data
	// log.Infof("配置: %+v\n", c)
	// log.Infof("GameDb配置: %+v\n", c.GameDb)
	// log.Infof("Etcd配置: %+v\n", c.Etcd)
	// log.Infof("Grpc配置: %+v\n", c.Grpc)
	return nil
}

// 读取应该配置
func LoadCliConfig(configDir string, envDir string) (*Config, error) {
	return LoadConfig(configDir, envDir, "cli", 0)
}

// 读取cli工具配置
func application_field(reader *config_reader, env map[string]interface{}, key string) interface{} {
	if v, ok := env["application"]; ok {
		if applications, ok := v.(map[string]interface{}); ok {
			if v, ok := applications[reader.appName]; ok {
				if application, ok := v.(map[string]interface{}); ok {
					if value, ok := application[key]; ok {
						return value
					}
				}
			}
		}
	}
	return ""
}

func other_application_field(reader *config_reader, otherAppType string, otherAppId int32, env map[string]interface{}, key string) interface{} {
	otherAppName := fmt.Sprintf("%s_%d", otherAppType, reader.appId+otherAppId)
	if v, ok := env["application"]; ok {
		if applications, ok := v.(map[string]interface{}); ok {
			if v, ok := applications[otherAppName]; ok {
				if application, ok := v.(map[string]interface{}); ok {
					if value, ok := application[key]; ok {
						return value
					}
				}
			}
		}
	}
	return ""
}

func host_field(reader *config_reader, env map[string]interface{}, key string) interface{} {
	var hostName string
	var hostFound bool = false
	if v, ok := env["application"]; ok {
		if applications, ok := v.(map[string]interface{}); ok {
			if v, ok := applications[reader.appName]; ok {
				if application, ok := v.(map[string]interface{}); ok {
					hostName = application["host"].(string)
					hostFound = true
				}
			}
		}
	}
	if !hostFound {
		return ""
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
	return ""
}

func other_host_field(reader *config_reader, otherAppType string, otherAppId int32, env map[string]interface{}, key string) interface{} {
	appName := fmt.Sprintf("%s_%d", otherAppType, reader.appId+otherAppId)
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
		return ""
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
	return ""

}

// 加载配置
// 根据环境，区名，服务名， 服务组合配置文件路径，规则是config/app/<<name>>.yaml
func (c *config_reader) read(dir string, envDir string) ([]byte, error) {
	var configFilePath = filepath.Join(dir, fmt.Sprintf("%s.yaml", c.appType))
	sb := strings.Builder{}
	// 预处理
	if err := c.preprocess(&sb, "", configFilePath); err != nil {
		return nil, err
	}
	// log.Infof("配置预处理后\n%v\n", sb.String())
	// 读环境变量
	yamlEnvFilePath := path.Join(envDir, ".config.yaml")
	dotEnvFilePath := path.Join(envDir, ".env")
	envData, err := c.readEnv(yamlEnvFilePath, dotEnvFilePath)
	if err != nil {
		return nil, err
	}
	envData["app_type"] = c.appType
	envData["app_name"] = c.appName
	envData["app_id"] = c.appId
	funcMap := template.FuncMap{
		// The name "title" is what the function will be called in the template text.
		"application_field": func(env map[string]interface{}, key string) interface{} {
			return application_field(c, env, key)
		},
		"other_application_field": func(otherAppType string, otherAppId int32, env map[string]interface{}, key string) interface{} {
			return other_application_field(c, otherAppType, otherAppId, env, key)
		},
		"host_field": func(env map[string]interface{}, key string) interface{} {
			return host_field(c, env, key)
		},

		"other_host_field": func(otherAppType string, otherAppId int32, env map[string]interface{}, key string) interface{} {
			return other_host_field(c, otherAppType, otherAppId, env, key)
		},

		"application_id": func() interface{} {
			return c.appId
		},

		"other_application_id": func(id int32) interface{} {
			return c.appId + id
		},
		"application_name": func() interface{} {
			return c.appType
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
	// log.Infof("替换环境变量后\n%v\n", out.String())
	return []byte(out.String()), nil
}

// 执行include指令
func (c *config_reader) readEnv(filePath string, dotEnvFilePath string) (map[string]interface{}, error) {
	if _, err := os.Stat(dotEnvFilePath); err == nil {
		if err := godotenv.Load(dotEnvFilePath); err != nil && err != os.ErrNotExist {
			return nil, err
		}
	}
	envData := make(map[string]interface{})
	sb := strings.Builder{}
	if _, err := os.Stat(filePath); err != nil {
		return envData, err
	}
	if err := c.preprocess(&sb, "", filePath); err != nil {
		return envData, err
	}
	if err := yaml.Unmarshal([]byte(sb.String()), envData); err != nil {
		return envData, err
	}
	if _, err := os.Stat(dotEnvFilePath); err == nil {
		if dict, err := godotenv.Read(dotEnvFilePath); err != nil {
			return nil, err
		} else {
			for k, v := range dict {
				envData[k] = v
			}
		}
	}
	c.env = envData["env"].(string)
	c.zone = envData["zone"].(string)
	return envData, nil
}

// 执行include指令
func (c *config_reader) preprocess(sb *strings.Builder, indent string, filePath string) error {
	dir := path.Dir(filePath)
	lines, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// 正则表达式，用于匹配形如${VAR_NAME}的环境变量
	re := regexp.MustCompile(`\${\w+}`)
	scanner := bufio.NewScanner(lines)
	scanner.Split(bufio.ScanLines)
	instruction := config_instruction_type_none
	for scanner.Scan() {
		line := scanner.Text()
		sline := strings.TrimSpace(line)
		if instruction == config_instruction_type_none {
			if strings.HasPrefix(sline, `$include:`) {
				instruction = config_instruction_type_include
			} else {
				// 循环匹配所有环境变量
				for _, match := range re.FindAllString(line, -1) {
					// 获取环境变量的值
					envValue := os.Getenv(match[2 : len(match)-1])
					// 将环境变量替换为其值
					line = re.ReplaceAllString(line, envValue)
				}
				// 循环匹配所有环境变量
				sb.WriteString(indent)
				sb.WriteString(line)
				sb.WriteString("\n")
			}
		}
		if instruction == config_instruction_type_include {
			var includeFilePath string
			var found bool = false
			if strings.HasPrefix(sline, `$include:`) {
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
				indent2 := strings.Replace(line, sline, "", 1) + indent
				if err := c.preprocess(sb, indent2, path.Join(dir, includeFilePath)); err != nil {
					return err
				}
			}
			if !found {
				instruction = config_instruction_type_none
				sb.WriteString(indent)
				sb.WriteString(line)
				sb.WriteString("\n")
			}
		}
	}
	lines.Close()
	return nil
}
