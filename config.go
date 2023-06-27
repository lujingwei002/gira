package gira

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/joho/godotenv"
	"github.com/lujingwei002/gira/proj"

	yaml "gopkg.in/yaml.v3"
)

const (
	config_instruction_type_none = iota
	config_instruction_type_include
	config_instruction_type_template
)

// 读取cli工具配置
// func application_field(reader *config_reader, env map[string]interface{}, key string) interface{} {
// 	if v, ok := env["application"]; ok {
// 		if applications, ok := v.(map[string]interface{}); ok {
// 			if v, ok := applications[reader.appName]; ok {
// 				if application, ok := v.(map[string]interface{}); ok {
// 					if value, ok := application[key]; ok {
// 						return value
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return ""
// }

// func other_application_field(reader *config_reader, otherAppType string, otherAppId int32, env map[string]interface{}, key string) interface{} {
// 	otherAppName := fmt.Sprintf("%s_%d", otherAppType, reader.appId+otherAppId)
// 	if v, ok := env["application"]; ok {
// 		if applications, ok := v.(map[string]interface{}); ok {
// 			if v, ok := applications[otherAppName]; ok {
// 				if application, ok := v.(map[string]interface{}); ok {
// 					if value, ok := application[key]; ok {
// 						return value
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return ""
// }

// func host_field(reader *config_reader, env map[string]interface{}, key string) interface{} {
// 	var hostName string
// 	var hostFound bool = false
// 	if v, ok := env["application"]; ok {
// 		if applications, ok := v.(map[string]interface{}); ok {
// 			if v, ok := applications[reader.appName]; ok {
// 				if application, ok := v.(map[string]interface{}); ok {
// 					hostName = application["host"].(string)
// 					hostFound = true
// 				}
// 			}
// 		}
// 	}
// 	if !hostFound {
// 		return ""
// 	}
// 	if v, ok := env["host"]; ok {
// 		if hosts, ok := v.(map[string]interface{}); ok {
// 			if v, ok := hosts[hostName]; ok {
// 				if host, ok := v.(map[string]interface{}); ok {
// 					if value, ok := host[key]; ok {
// 						return value
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return ""
// }

// func other_host_field(reader *config_reader, otherAppType string, otherAppId int32, env map[string]interface{}, key string) interface{} {
// 	appName := fmt.Sprintf("%s_%d", otherAppType, reader.appId+otherAppId)
// 	var hostName string
// 	var hostFound bool = false
// 	if v, ok := env["application"]; ok {
// 		if applications, ok := v.(map[string]interface{}); ok {
// 			if v, ok := applications[appName]; ok {
// 				if application, ok := v.(map[string]interface{}); ok {
// 					hostName = application["host"].(string)
// 					hostFound = true
// 				}
// 			}
// 		}
// 	}
// 	if !hostFound {
// 		return ""
// 	}
// 	if v, ok := env["host"]; ok {
// 		if hosts, ok := v.(map[string]interface{}); ok {
// 			if v, ok := hosts[hostName]; ok {
// 				if host, ok := v.(map[string]interface{}); ok {
// 					if value, ok := host[key]; ok {
// 						return value
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return ""

// }

// func printLines(str string) {
// 	lines := strings.Split(str, "\n")
// 	for num, line := range lines {
// 		fmt.Println(num+1, ":", line)
// 	}
// }

// 读取应该配置
func LoadCliConfig(configFilePath string, dotEnvFilePath string) (*Config, error) {
	if len(configFilePath) <= 0 {
		configFilePath = path.Join(proj.Config.ConfigDir, "cli.yaml")
	}
	if len(dotEnvFilePath) <= 0 {
		dotEnvFilePath = path.Join(proj.Config.EnvDir, ".env")
	}
	return LoadConfig(configFilePath, dotEnvFilePath, "cli", 0)
}

// 读取测试配置
func LoadTestConfig(configFilePath string, dotEnvFilePath string) (*Config, error) {
	if len(configFilePath) <= 0 {
		configFilePath = path.Join(proj.Config.ConfigDir, "test.yaml")
	}
	if len(dotEnvFilePath) <= 0 {
		dotEnvFilePath = path.Join(proj.Config.EnvDir, ".env")
	}
	return LoadConfig(configFilePath, dotEnvFilePath, "test", 0)
}

// 读取应用配置
func LoadApplicationConfig(configFilePath string, dotEnvFilePath string, appType string, appId int32) (*Config, error) {
	if len(configFilePath) <= 0 {
		configFilePath = path.Join(proj.Config.ConfigDir, fmt.Sprintf("%s.yaml", appType))
	}
	if len(dotEnvFilePath) <= 0 {
		dotEnvFilePath = path.Join(proj.Config.EnvDir, ".env")
	}
	if c, err := LoadConfig(configFilePath, dotEnvFilePath, appType, appId); err != nil {
		return nil, err
	} else {

		return c, nil
	}
}

// 读取配置
func LoadConfig(configFilePath string, dotEnvFilePath string, appType string, appId int32) (*Config, error) {
	c := newDefaultConfig()
	reader := config_reader{
		appType: appType,
		appId:   appId,
		appName: fmt.Sprintf("%s_%d", appType, appId),
	}
	if data, err := reader.read(configFilePath, dotEnvFilePath); err != nil {
		return nil, err
	} else if err := c.unmarshal(data); err != nil {
		return nil, NewSyntaxError(err.Error(), configFilePath, string(data))
	}
	// for _, v := range c.Db {
	// 	v.MaxOpenConns = 32
	// }
	return c, nil
}

func newDefaultConfig() *Config {
	return &Config{
		Thread:  1,
		Sandbox: 0,
	}
}

// 日志配置
type LogConfig struct {
	Console    bool   `yaml:"console"`
	File       bool   `yaml:"file"`
	MaxSize    int    `yaml:"max-size"`
	MaxBackups int    `yaml:"max-backups"`
	MaxAge     int    `yaml:"max-age"`
	Compress   bool   `yaml:"compress"`
	Level      string `yaml:"level"`
	DbLevel    string `yaml:"db-level"`
	Db         bool   `yaml:"db"`
	Dir        string `yaml:"dir"`
	Name       string `yaml:"name"`
}

// jwt配置
type JwtConfig struct {
	Secret            string `yaml:"secret"`
	Expiretime        int64  `yaml:"expiretime"`
	RefreshExpiretime int64  `yaml:"refresh-expiretime"`
}

type DbConfig struct {
	Driver          string        `yaml:"driver"`
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	Db              string        `yaml:"db"`
	Query           string        `yaml:"query"`
	MaxOpenConns    int           `yaml:"max-open-conns"`
	MaxIdleConns    int           `yaml:"max-idle-conns"`
	ConnMaxIdleTime time.Duration `yaml:"conn-max-idle-time"`
	ConnMaxLifetime time.Duration `yaml:"conn-max-lifetime"`
	ConnnectTimeout time.Duration `yaml:"connect-timeout"`
}

type BehaviorConfig struct {
	SyncInterval int64 `yaml:"sync-interval"`
	BatchInsert  int   `yaml:"batch-insert"`
}

// 完整的地址，包括path部分
func (self DbConfig) Uri() string {
	switch self.Driver {
	case MONGODB_NAME:
		return fmt.Sprintf("mongodb://%s:%s@%s:%d/?db=%s&%s", self.User, self.Password, self.Host, self.Port, self.Db, self.Query)
	case REDIS_NAME:
		return fmt.Sprintf("redis://%s:%s@%s:%d?%s", self.User, self.Password, self.Host, self.Port, self.Query)
	case MYSQL_NAME:
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&%s", self.User, self.Password, self.Host, self.Port, self.Db, self.Query)
	default:
		return fmt.Sprintf("%s not support", self.Driver)
	}
	//return fmt.Sprintf("%s://%s:%s@%s:%d/%s", self.Driver, self.User, self.Password, self.Host, self.Port, self.Db)
}

func (self DbConfig) GormUri() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", self.User, self.Password, self.Host, self.Port, self.Db)
}

func (self *DbConfig) Parse(uri string) error {
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}
	switch u.Scheme {
	case MONGODB_NAME:
		self.Driver = u.Scheme
	case REDIS_NAME:
		self.Driver = u.Scheme
	case MYSQL_NAME:
		self.Driver = u.Scheme
	default:
		return ErrDbNotSupport
	}
	host2 := strings.Split(u.Host, ":")
	if len(host2) == 2 {
		self.Host = host2[0]
		if v, err := strconv.Atoi(host2[1]); err != nil {
			return err
		} else {
			self.Port = v
		}
	} else if len(host2) == 1 {
		switch u.Scheme {
		case MONGODB_NAME:
			self.Port = 27017
		case REDIS_NAME:
			self.Port = 6379
		case MYSQL_NAME:
			self.Port = 3306
		}
	}
	self.User = u.User.Username()
	if v, set := u.User.Password(); set {
		self.Password = v
	} else {
		self.Password = ""
	}
	switch u.Scheme {
	case MONGODB_NAME:
		self.Db = u.Query().Get("db")
		query := u.Query()
		query.Del("db")
		self.Query = query.Encode()
	default:
		path := strings.TrimPrefix(u.Path, "/")
		self.Db = path
		self.Query = u.Query().Encode()
	}
	return nil
}

// registry配置
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

// registry配置
type EtcdClientConfig struct {
	Endpoints []struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"endpoints"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	DialTimeout  int    `yaml:"dial-timeout"`
	LeaseTimeout int64  `yaml:"lease-timeout"`
	Address      string `yaml:"address"`
}

// http模块配置
type HttpConfig struct {
	Addr         string `yaml:"addr"`
	ReadTimeout  int64  `yaml:"read-timeout"`
	WriteTimeout int64  `yaml:"write-timeout"`
	Ssl          bool   `yaml:"ssl"`
	CertFile     string `yaml:"cert-file"`
	KeyFile      string `yaml:"key-file"`
}

// 网关模块配置
type GatewayConfig struct {
	IsWebsocket      bool          `yaml:"is-websocket"`
	Bind             string        `yaml:"bind"`
	Address          string        `yaml:"address"`
	Debug            bool          `yaml:"debug"`
	Ssl              bool          `yaml:"ssl"`
	CertFile         string        `yaml:"cert-file"`
	KeyFile          string        `yaml:"key-file"`
	WsPath           string        `yaml:"ws-path"`
	RecvBuffSize     int           `yaml:"recv-buff-size"`
	RecvBacklog      int           `yaml:"recv-backlog"`
	SendBacklog      int           `yaml:"send-backlog"`
	HandshakeTimeout time.Duration `yaml:"handshake-timeout"`
	Heartbeat        time.Duration `yaml:"heartbeat"`
}

type TestSdkConfig struct {
	Secret string `yaml:"secret"`
}
type PwdSdkConfig struct {
	Secret string `yaml:"secret"`
}

type UltraSdkConfig struct {
	Secret string `yaml:"secret"`
}
type SdkConfig struct {
	Test  *TestSdkConfig  `yaml:"test"`
	Pwd   *PwdSdkConfig   `yaml:"pwd"`
	Ultra *UltraSdkConfig `yaml:"ultra"`
}

type GrpcConfig struct {
	Address  string `yaml:"address"`
	Workers  uint32 `yaml:"workers"`
	Resolver bool   `yaml:"resolver"` // 是否开启resolver
}

type PprofConfig struct {
	Port int    `yaml:"port"`
	Bind string `yaml:"bind"`
}

type ResourceConfig struct {
	Compress bool `yaml:"compress"`
}

type Config struct {
	Raw      []byte
	Thread   int         `yaml:"thread"`
	Env      string      `yaml:"env"`
	Zone     string      `yaml:"zone"`
	Log      *LogConfig  `yaml:"log"`
	CoreLog  *LogConfig  `yaml:"core-log"`
	Pprof    PprofConfig `yaml:"pprof"`
	Sandbox  int         `yaml:"sandbox"`
	Db       map[string]*DbConfig
	Resource ResourceConfig `yaml:"resource"`
	Module   struct {
		Behavior   *BehaviorConfig   `yaml:"behavior"`
		Http       *HttpConfig       `yaml:"http,omitempty"`
		Etcd       *EtcdConfig       `yaml:"etcd"`
		EtcdClient *EtcdClientConfig `yaml:"etcd-client"`
		Grpc       *GrpcConfig       `yaml:"grpc"`
		Sdk        *SdkConfig        `yaml:"sdk"`
		Jwt        *JwtConfig        `yaml:"jwt"`
		Gateway    *GatewayConfig    `yaml:"gateway"`
		Admin      *AdminConfig      `yaml:"admin"`
	} `yaml:"module"`
}

type AdminConfig struct {
	None string `yaml:"none"`
}

type config_reader struct {
	appId   int32
	appType string
	appName string
}

type config_file struct {
	builder   *strings.Builder
	templates []string // 模板
}

func (f *config_file) WriteString(s string) (int, error) {
	return f.builder.WriteString(s)
}

func (f *config_file) String() string {
	return f.builder.String()
}

func (c *Config) unmarshal(data []byte) error {
	// 解析yaml
	if err := yaml.Unmarshal(data, c); err != nil {
		return err
	}
	c.Raw = data
	return nil
}

// dotEnvFilePath 环境变量文件路径 .env
// 优先级 命令行 > 文件中appName指定变量 > 文件中变量
func (c *config_reader) readDotEnv(dotEnvFilePath string) error {
	fileEnv := make(map[string]string)
	appPrefixEnv := make(map[string]string)
	appNamePrefix := c.appName + "."
	if _, err := os.Stat(dotEnvFilePath); err == nil {
		if dict, err := godotenv.Read(dotEnvFilePath); err != nil {
			return NewSyntaxError(err.Error(), dotEnvFilePath, "")
		} else {
			for k, v := range dict {
				if strings.HasPrefix(k, appNamePrefix) {
					appPrefixEnv[strings.Replace(k, appNamePrefix, "", 1)] = v
				} else {
					fileEnv[k] = v
				}
			}
		}
	} else {
		return NewSyntaxError(err.Error(), dotEnvFilePath, "")
	}
	// 读文件中的环境变量到os.env中
	// if _, err := os.Stat(dotEnvFilePath); err == nil {
	// 	if err := godotenv.Load(dotEnvFilePath); err != nil && err != os.ErrNotExist {
	// 		return err
	// 	}
	// }
	for k, v := range appPrefixEnv {
		if _, ok := os.LookupEnv(k); !ok {
			os.Setenv(k, v)
		}
	}
	for k, v := range fileEnv {
		if _, ok := os.LookupEnv(k); !ok {
			os.Setenv(k, v)
		}
	}
	return nil
}

func (c *config_reader) templateFuncs() template.FuncMap {
	return template.FuncMap{
		"app_id": func() interface{} {
			return c.appId
		},
		"other_ap_id": func(id int32) interface{} {
			return c.appId + id
		},
		"app_type": func() interface{} {
			return c.appType
		},
		"app_name": func() interface{} {
			return c.appName
		},
		"work_dir": func() interface{} {
			return proj.Config.ProjectDir
		},
		"e": func(key string) interface{} {
			return os.Getenv(key)
		},
	}
}

// 加载配置
// 1.加载环境变量
// 2.加载模板变量
// 3.处理模板
// 4.yaml decode
func (c *config_reader) read(configFilePath string, dotEnvFilePath string) ([]byte, error) {
	var err error
	// 加载.env环境变量
	if err := c.readDotEnv(dotEnvFilePath); err != nil {
		return nil, err
	}
	file := &config_file{
		builder: &strings.Builder{},
	}
	// 预处理
	if err := c.preprocess(file, "", configFilePath); err != nil {
		return nil, err
	}
	// 处理模板变量
	if len(file.templates) > 0 {
		envData := make(map[string]interface{})
		for _, filePath := range file.templates {
			if envData, err = c.readTemplateData(filePath); err != nil {
				return nil, err
			}
		}
		// 替换环境变量
		t := template.New("config").Delims("<<", ">>")
		t.Funcs(c.templateFuncs())
		t, err = t.Parse(file.String())
		if err != nil {
			return nil, err
		}
		out := strings.Builder{}
		if err := t.Execute(&out, envData); err != nil {
			return nil, err
		}
		// log.Infof("替换环境变量后\n%v\n", out.String())
		return []byte(out.String()), nil
	} else {
		// 替换环境变量
		t := template.New("config").Delims("<<", ">>")
		t.Funcs(c.templateFuncs())
		t, err = t.Parse(file.String())
		if err != nil {
			return nil, err
		}
		out := strings.Builder{}
		if err := t.Execute(&out, nil); err != nil {
			return nil, err
		}
		// log.Infof("替换环境变量后\n%v\n", out.String())
		return []byte(out.String()), nil
	}
}

// 加载模板变量并返回
func (c *config_reader) readTemplateData(filePath string) (map[string]interface{}, error) {
	if _, err := os.Stat(filePath); err != nil {
		return nil, NewSyntaxError(err.Error(), filePath, "")
	}
	file := &config_file{
		builder: &strings.Builder{},
	}
	if err := c.preprocess(file, "", filePath); err != nil {
		return nil, err
	}
	var err error
	// 替换环境变量
	t := template.New("config").Delims("<<", ">>")
	t.Funcs(c.templateFuncs())
	t, err = t.Parse(file.String())
	if err != nil {
		return nil, NewSyntaxError(err.Error(), filePath, file.String())
	}
	out := strings.Builder{}
	if err := t.Execute(&out, nil); err != nil {
		return nil, NewSyntaxError(err.Error(), filePath, file.String())
	}
	envData := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(out.String()), envData); err != nil {
		return nil, NewSyntaxError(err.Error(), filePath, file.String())
	}
	return envData, nil
}

// 预处理文件filePath,输出到file中
// include指令
// template指令
// 替换${var}成环境变量
func (c *config_reader) preprocess(file *config_file, indent string, filePath string) error {
	dir := path.Dir(filePath)
	lines, err := os.Open(filePath)
	if err != nil {
		return NewSyntaxError(err.Error(), filePath, "")
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
			} else if strings.HasPrefix(sline, `$template:`) {
				instruction = config_instruction_type_template
			} else {
				// 循环匹配所有环境变量
				for _, match := range re.FindAllString(line, -1) {
					// 获取环境变量的值
					envName := match[2 : len(match)-1]
					envValue := os.Getenv(envName)
					// 将环境变量替换为其值
					line = strings.Replace(line, match, envValue, 1)
				}
				// 循环匹配所有环境变量
				file.WriteString(indent)
				file.WriteString(line)
				file.WriteString("\n")
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
				if err := c.preprocess(file, indent2, path.Join(dir, includeFilePath)); err != nil {
					return err
				}
			}
			if !found {
				instruction = config_instruction_type_none
				file.WriteString(indent)
				file.WriteString(line)
				file.WriteString("\n")
			}
		}
		if instruction == config_instruction_type_template {
			var templateFilePath string
			var found bool = false
			if strings.HasPrefix(sline, `$template:`) {
				pats := strings.SplitN(line, ":", 2)
				if len(pats) == 2 {
					templateFilePath = strings.TrimSpace(pats[1])
				}
				found = true
			} else {
				str := strings.TrimSpace(line)
				if strings.HasPrefix(str, "-") {
					str = strings.Replace(line, `-`, "", 1)
					templateFilePath = strings.TrimSpace(str)
					found = true
				}
			}
			if len(templateFilePath) > 0 {
				file.templates = append(file.templates, path.Join(dir, templateFilePath))
			}
			if !found {
				instruction = config_instruction_type_none
				file.WriteString(indent)
				file.WriteString(line)
				file.WriteString("\n")
			}
		}
	}
	lines.Close()
	return nil
}
