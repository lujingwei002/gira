package gen_configs

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/lujingwei/gira-cli/proj"
	excelize "github.com/xuri/excelize/v2"
	yaml "gopkg.in/yaml.v3"
)

// 字段类型
type field_type int

const (
	field_type_int field_type = iota
	field_type_string
	field_type_json
)

var type_name_dict = map[string]field_type{
	"int":    field_type_int,
	"string": field_type_string,
	"json":   field_type_json,
}

var go_type_name_dict = map[field_type]string{
	field_type_int:    "int",
	field_type_string: "string",
	field_type_json:   "interface{}",
}

// 字段结构
type config_field struct {
	tag  int
	name string     // 字段名
	Type field_type // 字段类型
}

type config_data struct {
	fieldArr []config_field  // 字段信息
	valueArr [][]interface{} // 字段值
}

type config_descriptor struct {
	yamlFilePath string
	filePath     string
	name         string
}
type configs_file struct {
	descriptorDict map[string]*config_descriptor
	loaderDict     map[string][]string
}

// 生成协议的状态
type config_state struct {
	configFilePath string
	docDir         string
	configsDir     string
	configsFile    configs_file
	configDataDict map[string]*config_data
}

func camelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

func toGoString(descriptor *config_descriptor, data *config_data) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("type %s struct {%s", descriptor.name, fmt.Sprintln()))
	for _, field := range data.fieldArr {
		typeName := go_type_name_dict[field.Type]
		sb.WriteString(fmt.Sprintf("    %s %s `yaml:%s`", camelString(field.name), typeName, field.name))
		sb.WriteString(fmt.Sprintln())
	}
	sb.WriteString(fmt.Sprintf("}%s", fmt.Sprintln()))

	sb.WriteString(fmt.Sprintf(`
type %sArr []%s
`, descriptor.name, descriptor.name))

	sb.WriteString(fmt.Sprintf(`
func (self *%sArr) Load(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, self); err != nil {
		return err
	}
	return nil
}
`, descriptor.name))
	return sb.String()
}

func toYamlString(descriptor *config_descriptor, data *config_data) string {
	sb := strings.Builder{}
	for _, row := range data.valueArr {
		sb.WriteString(fmt.Sprintln("-"))
		for index, v := range row {
			field := data.fieldArr[index]
			if field.Type == field_type_json && v == "" {
				sb.WriteString(fmt.Sprintf("  %s: %s%s", field.name, v, fmt.Sprintln()))
			} else if field.Type == field_type_json {
				sb.WriteString(fmt.Sprintf("  %s: |-%s", field.name, fmt.Sprintln()))
				sb.WriteString(fmt.Sprintf("    %s%s", v, fmt.Sprintln()))
			} else if field.Type == field_type_string {
				sb.WriteString(fmt.Sprintf("  %s: %s%s", field.name, v, fmt.Sprintln()))
			} else {
				sb.WriteString(fmt.Sprintf("  %s: %v%s", field.name, v, fmt.Sprintln()))
			}
		}
	}
	sb.WriteString(fmt.Sprintln())
	return sb.String()
}

func genConfigs1(configState *config_state, filePathArr []string) error {
	fmt.Println(filePathArr)

	if data, err := ioutil.ReadFile(configState.configFilePath); err != nil {
		return err
	} else {
		result := make(map[string]interface{})
		if err := yaml.Unmarshal(data, result); err != nil {
			return err
		}
		descriptors, _ := result["descriptor"]
		if descriptors == nil {
			return nil
		}
		for _, row := range descriptors.([]interface{}) {
			arr := row.([]interface{})
			filePath := arr[0].(string)
			yamlFilePath := strings.Replace(filePath, filepath.Ext(filePath), ".yaml", -1)
			item := &config_descriptor{
				name:         arr[1].(string),
				filePath:     filePath,
				yamlFilePath: yamlFilePath,
			}
			configState.configsFile.descriptorDict[filePath] = item
		}
		loaders, _ := result["loader"]
		for name, rows := range loaders.(map[string]interface{}) {
			arr := make([]string, 0)
			for _, row := range rows.([]interface{}) {
				arr = append(arr, row.(string))
			}
			configState.configsFile.loaderDict[name] = arr
		}
	}

	for _, filePath := range filePathArr {
		f, err := excelize.OpenFile(filePath)
		if err != nil {
			log.Println(err)
			return err
		}
		config := &config_data{
			fieldArr: make([]config_field, 0),
			valueArr: make([][]interface{}, 0),
		}
		// 获取 Sheet1 上所有单元格
		rows, err := f.GetRows("Sheet1")
		typeRow := rows[4]
		nameRow := rows[3]
		// 字段名
		for index, v := range nameRow {
			if v != "" {
				typeName := typeRow[index]
				if realType, ok := type_name_dict[typeRow[index]]; ok {
					field := config_field{
						name: v,
						Type: realType,
						tag:  index,
					}
					config.fieldArr = append(config.fieldArr, field)
				} else {
					return fmt.Errorf("invalid type %s", typeName)
				}
			}
		}
		// 值
		for index, row := range rows {
			if index <= 4 {
				continue
			}
			valueArr := make([]interface{}, 0)
			for _, field := range config.fieldArr {
				var v interface{}
				if len(row) > field.tag {
					v = row[field.tag]
				} else {
					v = ""
				}
				if field.Type == field_type_string {
				} else if field.Type == field_type_json {
				} else {
					if v == "" {
						v = 0
					}
				}
				valueArr = append(valueArr, v)
			}
			config.valueArr = append(config.valueArr, valueArr)
		}
		configName := strings.Replace(filePath, configState.configsDir, "", -1)
		configName = strings.Replace(configName, string(filepath.Separator), "", -1)
		configState.configDataDict[configName] = config
	}
	return nil
}

func genConfigs2(state *proj.State, configState *config_state) error {
	fmt.Println("生成yaml文件")
	configsDir := path.Join(state.GenDir, "configs")
	if err := os.RemoveAll(configsDir); err != nil {
		return err
	}
	if err := os.Mkdir(configsDir, 0755); err != nil {
		return err
	}
	for name, v := range configState.configsFile.descriptorDict {
		fmt.Println(name, "==>", v.yamlFilePath)
		filePath := path.Join(state.GenDir, "configs", v.yamlFilePath)
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		file.Truncate(0)
		if configData, ok := configState.configDataDict[name]; ok {
			file.WriteString(toYamlString(v, configData))
			file.WriteString("\n")
			file.Close()
		} else {
			file.Close()
			return fmt.Errorf("%s not found", name)
		}
	}
	fmt.Println("生成go文件")

	srcConfigsDir := path.Join(state.GenSrcDir, "configs")
	if err := os.RemoveAll(srcConfigsDir); err != nil {
		return err
	}
	if err := os.Mkdir(srcConfigsDir, 0755); err != nil {
		return err
	}
	configsPath := path.Join(srcConfigsDir, "configs.go")
	file, err := os.OpenFile(configsPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	file.Truncate(0)
	defer file.Close()

	file.WriteString(fmt.Sprintf(
		`package configs
import (
	yaml "gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
)

`))

	file.WriteString(fmt.Sprintf("type Config struct {%s", fmt.Sprintln()))
	for name, v := range configState.configsFile.descriptorDict {
		if _, ok := configState.configDataDict[name]; ok {
			file.WriteString(fmt.Sprintf("    %sArr %sArr%s", v.name, v.name, fmt.Sprintln()))
		} else {
			return fmt.Errorf("%s not found", name)
		}
	}
	file.WriteString(fmt.Sprintf("}%s", fmt.Sprintln()))
	file.WriteString(fmt.Sprintf(
		`
func (self *Config) LoadFromYaml(dir string) error {
	var filePath string
`))
	for _, v := range configState.configsFile.descriptorDict {
		fmt.Println(v.yamlFilePath)
		file.WriteString(fmt.Sprintf(
			`
	filePath = filepath.Join(dir, "%s")
	if err := self.%sArr.Load(filePath); err != nil {
		return err
	}`, v.yamlFilePath, v.name))
	}
	file.WriteString(fmt.Sprintf(
		`
	return nil
}
`))

	for name, v := range configState.configsFile.descriptorDict {
		if configData, ok := configState.configDataDict[name]; ok {
			file.WriteString("\n")
			file.WriteString(toGoString(v, configData))
			file.WriteString("\n")
		} else {
			return fmt.Errorf("%s not found", name)
		}
	}
	file.Close()
	return nil
}

// 生成协议
func Gen(state *proj.State) error {
	// 初始化
	configFilePath := path.Join(state.DocDir, "configs.yaml")
	configState := &config_state{
		configDataDict: make(map[string]*config_data),
		configFilePath: configFilePath,
		docDir:         state.DocDir,
		configsDir:     filepath.Join(state.DocDir, "configs"),
		configsFile: configs_file{
			descriptorDict: make(map[string]*config_descriptor, 0),
			loaderDict:     make(map[string][]string),
		},
	}

	configFilePathArr := make([]string, 0)
	filepath.WalkDir(configState.configsDir, func(path string, d os.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		log.Println(path[0:1])
		if filepath.HasPrefix(filepath.Base(path), ".") {
			return nil
		}
		if filepath.Ext(path) == ".xlsx" {
			configFilePathArr = append(configFilePathArr, path)
		}
		return nil
	})

	if err := genConfigs1(configState, configFilePathArr); err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("ccccccccccccccc1")
	if err := genConfigs2(state, configState); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
