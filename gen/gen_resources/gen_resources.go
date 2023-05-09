package gen_resources

/// 参考 https://www.cnblogs.com/f-ck-need-u/p/10053124.html

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/log"

	"github.com/lujingwei002/gira/proj"
	excelize "github.com/xuri/excelize/v2"
	yaml "gopkg.in/yaml.v3"
)

var cli_code = `
// Code generated by github.com/lujingwei002/gira. DO NOT EDIT.
// Code generated by github.com/lujingwei002/gira. DO NOT EDIT.
// Code generated by github.com/lujingwei002/gira. DO NOT EDIT.

package main

import (
	"context"
	"log"
	"os"
	"github.com/urfave/cli/v2"
	"<<.Module>>/gen/resource"
)

var uri string

func main() {
	app := &cli.App{
		Name: "gira-resource",
		Authors: []*cli.Author{{
			Name:  "lujingwei",
			Email: "lujingwei@xx.org",
		}},
		Description: "gira-cli",
		Flags:       []cli.Flag{},
		Commands: []*cli.Command{
			{
				Name:   "compress",
				Usage:  "compress yaml to binary",
				Action: compressAction,
			},
			{
				Name:      "push",
				Usage:     "push to database",
				Action:    pushAction,
				Flags:       []cli.Flag{
					&cli.StringFlag{
						Name: "uri",
						Required: true,
						Usage: "database uri",
						Destination: &uri,
					},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Println(err)
	}
}

func compressAction(args *cli.Context) error {
	if err := resource.Compress("resource"); err != nil {
		return err
	} else {
		log.Println("success")
		return nil
	}
}

func pushAction(args *cli.Context) error {
	return resource.Push(context.Background(), uri, "resource")
}
`

var code = `
// Code generated by github.com/lujingwei002/gira. DO NOT EDIT.
// Code generated by github.com/lujingwei002/gira. DO NOT EDIT.
// Code generated by github.com/lujingwei002/gira. DO NOT EDIT.

package resource
import (
	"github.com/lujingwei002/gira"
	yaml "gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
	"encoding/gob"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/bson"
	"time"
	"github.com/lujingwei002/gira/log"
	"os"
	"fmt"
	"context"
	"net/url"
	"strings"
)

// 将Db类型的bundle推送覆盖到db上
//
// Parameters:
// uri - mongodb://root:123456@192.168.1.200:3331/resourcedb
// dir - bundle所在的目录
func Push(ctx context.Context, uri string, dir string) error {
	u, err := url.Parse(uri)
    if err != nil {
        log.Infow("parsing mongodb uri fail", "error", err)
        return err
    }
    path := strings.TrimPrefix(u.Path, "/")
	uri = strings.Replace(uri, u.Path, "", 1)
	log.Infow("connect database", "uri", uri, "path", path)

	clientOpts := options.Client().ApplyURI(uri)
	ctx1, cancelFunc1 := context.WithTimeout(ctx, 3*time.Second)
	defer cancelFunc1()
	client, err := mongo.Connect(ctx1, clientOpts)
	if err != nil {
		log.Infow("connect accountdb fail", "error", err)
		return err
	}
	ctx2, cancelFunc2 := context.WithTimeout(ctx, 3*time.Second)
	defer cancelFunc2()
	if err = client.Ping(ctx2, readpref.Primary()); err != nil {
		log.Errorw("connect database fail", "uri", uri, "error", err)
		return err
	}
	database := client.Database(path)

	<<- range .BundleArr>>
	<<- if  eq .BundleType "db">>
	<<.CapBundleStructName>> := &<<.BundleStructName>>{}
	// 推送<<.BundleName>>
	if err := <<.CapBundleStructName>>.SaveToDb(ctx, database, dir); err != nil {
        return err
	}
	<<- end>>
	<<- end>>
	return nil
}

// 将目录下的配置文件压缩成二进制bundle格式
//
// Parameters:
// dir - 配置文件yaml所在的目录
func Compress(dir string) error {
<<range .BundleArr>>
	// <<.BundleStructName>>
	<<.CapBundleStructName>> := &<<.BundleStructName>>{}
	if err := <<.CapBundleStructName>>.LoadFromYaml(dir); err != nil {
        return err
	}
	if err := <<.CapBundleStructName>>.SaveToBin(dir); err != nil {
        return err
	}
<<end>>
	return nil
}


<<range .LoaderArr>>
// <<.LoaderName>>的handler接口，必须实现这些方法对加载了的配置进行处理
type <<.HandlerStructName>> interface {
	<<- range .BundleArr>>
		<<- range .ResourceArr>>
	Convert<<.StructName>>(arr <<.WrapStructName>>) error
	Load<<.StructName>>(reload bool) error
		<<- end>>
	<<- end>>
}
<<end>>


<<range .LoaderArr>>
// <<.LoaderName>>加载器，负载加载拥有的bundle
type <<.LoaderStructName>> struct {
	<<- range .BundleArr>>
	<<.BundleStructName>>
	<<- end>>
}

// 从yaml文件加载资源
func (self *<<.LoaderStructName>>) LoadFromYaml(dir string) error {
	<<- range .BundleArr>>
	if err := self.<<.BundleStructName>>.LoadFromYaml(dir); err != nil {
		return err
	}
	<<- end>>
	return nil
}

// 从bin文件加载资源
func (self *<<.LoaderStructName>>) LoadFromBin(dir string) error {
	<<- range .BundleArr>>
	if err := self.<<.BundleStructName>>.LoadFromBin(dir); err != nil {
		return err
	}
	<<- end>>
	return nil
}

// 从db中加载资源
func (self *<<.LoaderStructName>>) LoadFromDb(ctx context.Context, database *mongo.Database) error {
	<<- range .BundleArr>>
	if err := self.<<.BundleStructName>>.LoadFromDb(ctx, database); err != nil {
		return err
	}
	<<- end>>
	return nil
}

// 根据bundle的类型，从相应的源中加载资源
func (self *<<.LoaderStructName>>) Load(ctx context.Context, database *mongo.Database, dir string) error {
	<<- range .BundleArr>>
	<<- if eq .BundleType "db">>
	if err := self.<<.BundleStructName>>.LoadFromDb(ctx, database); err != nil {
		return err
	}
	<<- end>>
	<<- if eq .BundleType "raw">>
	if err := self.<<.BundleStructName>>.LoadFromYaml(dir); err != nil {
		return err
	}
	<<- end>>
	<<- if eq .BundleType "bin">>
	if err := self.<<.BundleStructName>>.LoadFromBin(dir); err != nil {
		return err
	}
	<<- end>>
	<<- end>>
	return nil
}

// 加载成功后，对配置进行处理
func (self *<<.LoaderStructName>>) Convert(handler gira.ResourceHandler) error {
	handler.OnResourcePreLoad()
	h := handler.(<<.HandlerStructName>>)
	<<- range .BundleArr>>
		<<- $bundleStructName := .BundleStructName>>
		<<- range .ResourceArr>>
	if err := h.Convert<<.StructName>>(self.<<$bundleStructName>>.<<.WrapStructName>>); err != nil {
		return err
	}
		<<- end>>
	<<- end>>
	handler.OnResourcePostLoad()
	return nil
}
<<end>>

<<range .BundleArr>>
type <<.BundleStructName>> struct {
<<- range .ResourceArr>>
	<<.WrapStructName>> <<.WrapStructName>>
<<- end>>
}

func (self* <<.BundleStructName>>) Clear() {
	<<- range .ResourceArr>>
		<<- if .IsDeriveObject>>
		<<- else>>
	self.<<.WrapStructName>> = make(<<.WrapTypeName>>, 0)
		<<- end>>
	<<- end>>
}

// 从yaml文件加载资源
func (self* <<.BundleStructName>>) LoadFromYaml(dir string) error {
	self.Clear()
	<<- range .ResourceArr>>
	var <<.CapStructName>>filePath = filepath.Join(dir, "<<.YamlFileName>>")
	if err := self.<<.WrapStructName>>.LoadFromYaml(<<.CapStructName>>filePath); err != nil {
		return err
	}
	<<end>>
	return nil
}

// 从db中加载资源
func (self* <<.BundleStructName>>) LoadFromDb(ctx context.Context, database *mongo.Database) error {
	self.Clear()
	<<- range .ResourceArr>>
	if err := self.<<.WrapStructName>>.LoadFromDb(ctx, database); err != nil {
		return err
	}
	<<- end>>
	return nil
}

// 从bin文件中加载资源
func (self* <<.BundleStructName>>) LoadFromBin(dir string) error {
	var filePath = filepath.Join(dir, "<<.BundleName>>.bin")
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := gob.NewDecoder(f)
	if err := encoder.Decode(self); err != nil {
		return err
	}
	return nil
}

// 将资源保存到bin文件中
func (self *<<.BundleStructName>>) SaveToBin(dir string) error {
	var filePath = filepath.Join(dir, "<<.BundleName>>.bin")
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := gob.NewEncoder(f)
	if err := encoder.Encode(self); err != nil {
		return err
	}
	return nil
}

// 将资源保存到db上
func (self *<<.BundleStructName>>) SaveToDb(ctx context.Context, database *mongo.Database, dir string) error {
	<<- range .ResourceArr>>
	// 加载<<.ResourceName>>
	var <<.CapStructName>>filePath = filepath.Join(dir, "<<.YamlFileName>>")
	var <<.ArrTypeName>> <<.ArrTypeName>>
	if err := <<.ArrTypeName>>.LoadFromYaml(<<.CapStructName>>filePath); err != nil {
		return err
	}
	<<- end>>
	<<- range .ResourceArr>>
	// 保存 <<.ResourceName>>
	if coll := database.Collection("<<.TableName>>"); coll == nil {
		return fmt.Errorf("collection <<.TableName>> not found")
	} else {
		coll.Drop(ctx)
		models := make([]mongo.WriteModel, 0)
		for _, v := range <<.ArrTypeName>> {
			models = append(models, mongo.NewInsertOneModel().SetDocument(v))
		}
		if _, err := coll.BulkWrite(ctx, models); err != nil {
			log.Info("push <<.TableName>> fail", err)
			return err
		} else {
			log.Infof("push <<.TableName>>(%d) success", len(<<.ArrTypeName>>))
		}
	}
	<<end>>
	return nil
}
<<end>>

<<- range .ResourceArr>>
type <<.StructName>> struct {
	<<- range .FieldArr>>
	// <<.Comment>>
	<<.StructFieldName>> <<.GoTypeName>> <<quote>>bson:"<<.FieldName>>" json:"<<.FieldName>>" yaml:"<<.FieldName>>"<<quote>>
	<<- end>>
}

type <<.ArrTypeName>> []*<<.StructName>>

// 从yaml文件加载资源
func (self *<<.ArrTypeName>>) LoadFromYaml(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, self); err != nil {
		log.Warnw("load resource fail", "name", filePath)
		return err
	}
	return nil
}

// 从db中加载资源
func (self *<<.ArrTypeName>>) LoadFromDb(ctx context.Context, database *mongo.Database) error {
	coll := database.Collection("<<.TableName>>")
	if cursor, err := coll.Find(ctx, bson.D{}); err != nil {
		return err
	} else {
		defer cursor.Close(ctx)
		if err := cursor.All(ctx, self); err != nil {
			return err
		}
	}
	return nil
}

/// 字典类型的配置,转换成字典格式
<<- if .IsDeriveMap>>
type <<.MapTypeName>> <<.GoMapTypeName>>
func (self *<<.WrapStructName>>) LoadFromYaml(filePath string) error {
	var arr <<.ArrTypeName>>
	if err := arr.LoadFromYaml(filePath); err != nil {
		return err
	}
	return self.Make(arr)
}

func (self *<<.WrapStructName>>) LoadFromDb(ctx context.Context, database *mongo.Database) error {
	var arr <<.ArrTypeName>>
	if err := arr.LoadFromDb(ctx, database); err != nil {
		return err
	}
	return self.Make(arr)
}

func (self *<<.WrapStructName>>) Make(arr <<.ArrTypeName>>) error {
	if err := gira.Make<<len .MapKeyArr>>Key_<<join .MapKeyGoTypeNameArr "_">>(arr, *self<<- range $k, $v := .MapKeyArr>>, "<<camelString $v>>"<<- end>>); err != nil {
		return err
	}
	return nil
}
<<- end>>

/// 对像类型的配置,转换成对象格式
<<- if .IsDeriveObject>>
type <<.ObjectTypeName>> struct {
	<<$resource := .>>
	<<- range .ValueArr>>
	<<- $v := index . $resource.ObjectKeyIndex>>
	<<- if eq $v "">>
	<<- else>>
	<< camelString $v>> *<<$resource.StructName>>
	<<- end>>
	<<- end>>
}

func (self *<<.WrapStructName>>) LoadFromYaml(filePath string) error {
	var arr <<.ArrTypeName>>
	if err := arr.LoadFromYaml(filePath); err != nil {
		return err
	}
	return self.Make(arr)
}


func (self *<<.WrapStructName>>) LoadFromDb(ctx context.Context, database *mongo.Database) error {
	var arr <<.ArrTypeName>>
	if err := arr.LoadFromDb(ctx, database); err != nil {
		return err
	}
	return self.Make(arr)
}

func (self *<<.WrapStructName>>) Make(arr <<.ArrTypeName>>) error {
	dict := make(<<.GoObjectTypeName>> ,0)
	if err := gira.Make<<len .ObjectKeyArr>>Key_<<join .ObjectKeyGoTypeNameArr "_">>(arr, dict<<- range $k, $v := .ObjectKeyArr>>, "<<camelString $v>>"<<- end>>); err != nil {
		return err
	}
	var ok bool
	<<- range .ValueArr>>
		<<- $v := index . $resource.ObjectKeyIndex>>
		<<- if eq $v "">>
		<<- else>>
	if self.<<- camelString $v>>, ok = dict["<<$v>>"]; !ok {
		return fmt.Errorf("<<$resource.StructName>> <<$v>> key not found")
	}
		<<- end>>
	<<- end>>
	return nil
}
	<<- end>>


<<- end>>

`

func QuoteChar() interface{} {
	return "`"
}

// 字段类型
type field_type int

const (
	field_type_int field_type = iota
	field_type_int32
	field_type_int64
	field_type_string
	field_type_json
	field_type_bool
	field_type_string_arr
	field_type_int_arr
	field_type_float_arr
	field_type_struct
)

type resource_type int

const (
	resource_type_array = iota
	resource_type_map
	resource_type_object
)

// excel类型字符串和类型的对应关系
var type_name_dict = map[string]field_type{
	"int":      field_type_int,
	"int64":    field_type_int64,
	"int32":    field_type_int32,
	"string":   field_type_string,
	"json":     field_type_json,
	"bool":     field_type_bool,
	"string[]": field_type_string_arr,
	"int[]":    field_type_int_arr,
	"float[]":  field_type_float_arr,
}

// 和go类型的对应关系
var go_type_name_dict = map[field_type]string{
	field_type_int:        "int64",
	field_type_int64:      "int64",
	field_type_int32:      "int32",
	field_type_string:     "string",
	field_type_json:       "interface{}",
	field_type_bool:       "bool",
	field_type_string_arr: "[]string",
	field_type_int_arr:    "[]int64",
	field_type_float_arr:  "[]float64",
}

var resource_type_name_dict = map[string]resource_type{
	"map":    resource_type_map,
	"object": resource_type_object,
	"array":  resource_type_array,
}

// 字段结构
type Field struct {
	Tag             int
	FieldName       string     // 字段名
	StructFieldName string     // 字段名
	Type            field_type // 字段类型
	GoTypeName      string
	Comment         string
	Default         interface{}
}
type SortFieldByName []*Field

func (self SortFieldByName) Len() int           { return len(self) }
func (self SortFieldByName) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }
func (self SortFieldByName) Less(i, j int) bool { return self[i].Tag < self[j].Tag }

type Resource struct {
	ResourceName   string
	StructName     string
	TableName      string
	CapStructName  string
	WrapStructName string
	WrapTypeName   string
	YamlFileName   string
	ArrTypeName    string // ErrorCodeArr
	FilePath       string
	Type           resource_type

	// map类型
	MapTypeName            string // ErrorCodeMap
	GoMapTypeName          string // map[int] *ErrorCode
	MapKeyArr              []string
	ObjectKeyGoTypeNameArr []string

	// object类型
	ObjectTypeName      string
	GoObjectTypeName    string // map[int] *ErrorCode
	MapKeyGoTypeNameArr []string
	ObjectKeyArr        []string
	ObjectKeyIndex      int

	FieldDict map[string]*Field // 字段信息
	FieldArr  []*Field          // 字段信息
	ValueArr  [][]interface{}   // 字段值
}

type SortResourceByName []*Resource

func (self SortResourceByName) Len() int           { return len(self) }
func (self SortResourceByName) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }
func (self SortResourceByName) Less(i, j int) bool { return self[i].StructName < self[j].StructName }

func (self *Resource) IsDeriveMap() bool {
	return self.Type == resource_type_map
}

func (self *Resource) IsDeriveObject() bool {
	return self.Type == resource_type_object
}

type Bundle struct {
	BundleType          string // file binary db
	BundleName          string
	BundleStructName    string
	CapBundleStructName string
	ResourceNameArr     []string
	ResourceArr         []*Resource
}

type SortBundleByName []*Bundle

func (self SortBundleByName) Len() int           { return len(self) }
func (self SortBundleByName) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }
func (self SortBundleByName) Less(i, j int) bool { return self[i].BundleName < self[j].BundleName }

type Loader struct {
	LoaderStructName  string
	LoaderName        string
	HandlerStructName string
	bundleNameArr     []string
	BundleArr         []*Bundle
}

type SortLoaderByName []*Loader

func (self SortLoaderByName) Len() int           { return len(self) }
func (self SortLoaderByName) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }
func (self SortLoaderByName) Less(i, j int) bool { return self[i].LoaderName < self[j].LoaderName }

// 生成协议的状态
type gen_state struct {
	Config       Config
	Module       string
	ResourceDict map[string]*Resource
	ResourceArr  []*Resource
	BundleDict   map[string]*Bundle
	BundleArr    []*Bundle
	LoaderArr    []*Loader
}

func capUpperString(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}

func capLowerString(s string) string {
	return strings.ToLower(s[0:1]) + s[1:]
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

func fotmatYamlString(resource *Resource) string {
	sb := strings.Builder{}
	for _, row := range resource.ValueArr {
		sb.WriteString(fmt.Sprintln("-"))
		for index, v := range row {
			field := resource.FieldArr[index]
			if field.Type == field_type_json && v == "" {
				sb.WriteString(fmt.Sprintf("  %s: %s%s", field.FieldName, v, fmt.Sprintln()))
			} else if field.Type == field_type_json {
				sb.WriteString(fmt.Sprintf("  %s: |-%s", field.FieldName, fmt.Sprintln()))
				sb.WriteString(fmt.Sprintf("    %s%s", v, fmt.Sprintln()))
			} else if field.Type == field_type_string {
				str := v.(string)
				str = strings.ReplaceAll(str, "\r\n", " ")
				str = strings.ReplaceAll(str, "\n", " ")
				sb.WriteString(fmt.Sprintf("  %s: %s%s", field.FieldName, str, fmt.Sprintln()))
			} else if field.Type == field_type_struct {
				sb.WriteString(fmt.Sprintf("  %s: %v%s", field.FieldName, v, fmt.Sprintln()))
			} else {
				sb.WriteString(fmt.Sprintf("  %s: %v%s", field.FieldName, v, fmt.Sprintln()))
			}
		}
	}
	sb.WriteString(fmt.Sprintln())
	return sb.String()
}

type struct_type struct {
	Format    string                 `yaml:"format"`
	Excel     string                 `yaml:"excel"`
	Bundles   []string               `yaml:"bundles"`
	Resources []string               `yaml:"resources"`
	Struct    map[string]interface{} `yaml:"struct"`
	Map       []string               `yaml:"map"`
	Object    []string               `yaml:"object"`
}

func (r *Resource) readExcel(filePath string) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Info(err)
		return err
	}
	// 获取 Sheet1 上所有单元格
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return err
	}
	commentRow := rows[0]
	nameRow := rows[3]
	typeRow := rows[4]
	// 字段名
	for index, v := range nameRow {
		if v != "" {
			typeName := typeRow[index]
			comment := ""
			if index < len(commentRow) {
				comment = commentRow[index]
			}
			// 字段类型
			if realType, ok := type_name_dict[typeRow[index]]; ok {
				field := &Field{
					FieldName:       v,
					StructFieldName: capUpperString(camelString(v)),
					Type:            realType,
					Tag:             index,
					GoTypeName:      go_type_name_dict[realType],
					Comment:         comment,
				}
				r.FieldArr = append(r.FieldArr, field)
				r.FieldDict[v] = field
			} else {
				log.Warnw("invalid type", "type", typeName, "field", v)
				return fmt.Errorf("invalid type %s", typeName)
			}
		}
	}
	//值
	for index, row := range rows {
		if index <= 4 {
			continue
		}
		valueArr := make([]interface{}, 0)
		for _, field := range r.FieldArr {
			var v interface{}
			if len(row) > field.Tag {
				v = row[field.Tag]
			} else {
				v = ""
			}
			if field.Type == field_type_string {
			} else if field.Type == field_type_bool {
				if v == "" {
					v = "false"
				}
			} else if field.Type == field_type_json {
			} else if field.Type == field_type_string_arr {
			} else if field.Type == field_type_int_arr {
				v = fmt.Sprintf("[%s]", v)
			} else if field.Type == field_type_float_arr {
				v = fmt.Sprintf("[%s]", v)
			} else {
				if v == "" {
					v = 0
				}
			}
			valueArr = append(valueArr, v)
		}
		r.ValueArr = append(r.ValueArr, valueArr)
	}
	// sort.Sort(SortFieldByName(r.FieldArr))

	return nil
}

func getSrcFileHash(arr []string) string {
	sort.Strings(arr)
	sb := strings.Builder{}
	for _, filePath := range arr {
		if v, err := os.Stat(filePath); err == nil {
			sb.WriteString(fmt.Sprintf("%s %v\n", filePath, v.ModTime()))
		}
	}
	hash := md5.New()
	hash.Write([]byte(sb.String()))
	md5Hash := hex.EncodeToString(hash.Sum(nil))
	return md5Hash
}

func (resource *Resource) parseStruct(attrs map[string]interface{}) error {
	spaceRegexp := regexp.MustCompile("[^\\s]+")
	equalRegexp := regexp.MustCompile("[^=]+")
	fileArr := make([]*Field, 0)

	for valueStr, v := range attrs {
		var tag int
		var err error
		var fieldName string
		var typeStr string
		var tagStr string
		var optionArr []interface{}
		switch v.(type) {
		case nil:
			break
		case []interface{}:
			optionArr = v.([]interface{})
		default:
			return fmt.Errorf("%+v invalid11", v)
		}
		args := equalRegexp.FindAllString(valueStr, -1)
		if len(args) != 2 {
			return fmt.Errorf("%s invalid", valueStr)
		}
		tagStr = strings.TrimSpace(args[1])
		if tag, err = strconv.Atoi(tagStr); err != nil {
			return err
		}
		args = spaceRegexp.FindAllString(args[0], -1)
		if len(args) != 2 {
			return fmt.Errorf("%s invalid", valueStr)
		}
		fieldName = args[0]
		typeStr = args[1]
		field := &Field{
			FieldName:       fieldName,
			StructFieldName: camelString(fieldName),
			Tag:             tag,
		}
		if typeValue, ok := type_name_dict[typeStr]; ok {
			field.Type = typeValue
			field.GoTypeName = go_type_name_dict[field.Type]
		} else {
			field.Type = field_type_struct
			field.GoTypeName = typeStr
		}
		for _, option := range optionArr {
			optionDict := option.(map[string]interface{})
			if defaultVal, ok := optionDict["default"]; ok {
				field.Default = defaultVal
			}
			if comment, ok := optionDict["comment"]; ok {
				field.Comment = comment.(string)
			}
		}
		if last, ok := resource.FieldDict[fieldName]; ok {
			resource.FieldDict[fieldName] = field
			// 替换之前的
			for k, v := range resource.FieldArr {
				if v == last {
					resource.FieldArr[k] = field
					break
				}
			}
		} else {
			resource.FieldDict[fieldName] = field
			fileArr = append(fileArr, field)

		}
	}
	sort.Sort(SortFieldByName(fileArr))
	resource.FieldArr = append(resource.FieldArr, fileArr...)
	return nil
}

func parse(state *gen_state) error {
	// srcFilePathArr := make([]string, 0)
	// srcFilePathArr = append(srcFilePathArr, proj.Config.DocResourceFilePath)

	if data, err := ioutil.ReadFile(proj.Config.DocResourceFilePath); err != nil {
		return err
	} else {
		structDict := make(map[string]*struct_type)
		if err := yaml.Unmarshal(data, &structDict); err != nil {
			return err
		}
		for name, v := range structDict {
			// bundle
			if len(v.Resources) > 0 {
				bundle := &Bundle{
					BundleType:          v.Format,
					BundleName:          name,
					BundleStructName:    name,
					CapBundleStructName: capLowerString(name),
					ResourceNameArr:     v.Resources,
					ResourceArr:         make([]*Resource, 0),
				}
				state.BundleDict[name] = bundle
				state.BundleArr = append(state.BundleArr, bundle)

			} else if len(v.Bundles) > 0 {
				// loader
				loader := &Loader{
					LoaderName:        name,
					LoaderStructName:  name,
					HandlerStructName: "I" + name + "Handler",
					bundleNameArr:     v.Bundles,
					BundleArr:         make([]*Bundle, 0),
				}
				state.LoaderArr = append(state.LoaderArr, loader)

			} else if v.Excel != "" {
				var filePath string
				filePath = v.Excel
				// 类型
				var typ resource_type = resource_type_array
				if len(v.Map) > 0 {
					typ = resource_type_map
				}
				if len(v.Object) > 0 {
					typ = resource_type_object
				}
				// 解析key arr
				// xlsx替换成yaml
				//yamlFileName := strings.Replace(filePath, filepath.Ext(filePath), ".yaml", -1)
				yamlFileName := fmt.Sprintf("%s.yaml", name)
				camelName := camelString(name)
				r := &Resource{
					FieldDict:      make(map[string]*Field, 0),
					FieldArr:       make([]*Field, 0),
					ValueArr:       make([][]interface{}, 0),
					Type:           typ,
					ResourceName:   name,
					StructName:     camelName,
					TableName:      name,
					CapStructName:  capLowerString(camelName),
					ArrTypeName:    camelName + "Arr",
					MapTypeName:    camelName + "Map",
					ObjectTypeName: camelName + "Object",
					FilePath:       filePath,
					YamlFileName:   yamlFileName,
					MapKeyArr:      v.Map,
					ObjectKeyArr:   v.Object,
				}
				if typ == resource_type_array {
					r.WrapStructName = r.ArrTypeName
				} else if typ == resource_type_map {
					r.WrapStructName = r.MapTypeName
				} else if typ == resource_type_object {
					r.WrapStructName = r.ObjectTypeName
				}
				// 解析excel
				// filePath := path.Join(proj.Config.ExcelDir, v.FilePath)
				if err := r.readExcel(path.Join(proj.Config.ExcelDir, filePath)); err != nil {
					return err
				}
				// 解析struct
				if len(v.Struct) > 0 {
					if err := r.parseStruct(v.Struct); err != nil {
						return err
					}
				}
				// 处理不同的转换类型， 转map, 转object
				//GoMapTypeName
				r.GoMapTypeName = ""
				for _, v := range r.MapKeyArr {
					if field, ok := r.FieldDict[v]; ok {
						r.GoMapTypeName = r.GoMapTypeName + fmt.Sprintf(`map[%s]`, field.GoTypeName)
					} else {
						return fmt.Errorf("resource %s key %s not found\n", r.StructName, v)
					}
				}
				r.GoMapTypeName = r.GoMapTypeName + fmt.Sprintf(" *%s", r.StructName)
				// GoObjectTypeName
				r.GoObjectTypeName = ""
				for _, v := range r.ObjectKeyArr {
					if field, ok := r.FieldDict[v]; ok {
						r.GoObjectTypeName = r.GoObjectTypeName + fmt.Sprintf(`map[%s]`, field.GoTypeName)
					} else {
						return fmt.Errorf("resource %s key %s not found\n", r.StructName, v)
					}
				}
				r.GoObjectTypeName = r.GoObjectTypeName + fmt.Sprintf(" *%s", r.StructName)
				// WrapTypeName
				if r.Type == resource_type_map {
					r.WrapTypeName = r.GoMapTypeName
				} else if r.Type == resource_type_array {
					r.WrapTypeName = fmt.Sprintf("[]* %s", r.StructName)
				} else if r.Type == resource_type_object {
					for k, v := range r.FieldArr {
						if v.FieldName == r.ObjectKeyArr[0] {
							r.ObjectKeyIndex = k
							break
						}
					}
				}
				// make key时用的
				mapKeyGoTypeNameArr := make([]string, 0)
				for _, k := range r.MapKeyArr {
					f, _ := r.FieldDict[k]
					mapKeyGoTypeNameArr = append(mapKeyGoTypeNameArr, f.GoTypeName)
				}
				r.MapKeyGoTypeNameArr = mapKeyGoTypeNameArr

				// make key时用的
				objectKeyGoTypeNameArr := make([]string, 0)
				for _, k := range r.ObjectKeyArr {
					f, _ := r.FieldDict[k]
					objectKeyGoTypeNameArr = append(objectKeyGoTypeNameArr, f.GoTypeName)
				}
				r.ObjectKeyGoTypeNameArr = objectKeyGoTypeNameArr
				state.ResourceDict[name] = r
				state.ResourceArr = append(state.ResourceArr, r)
			}
		}
		for _, v := range state.BundleArr {
			for _, name := range v.ResourceNameArr {
				if r, ok := state.ResourceDict[name]; ok {
					v.ResourceArr = append(v.ResourceArr, r)
				}
			}
		}
		for _, v := range state.LoaderArr {
			for _, name := range v.bundleNameArr {
				if r, ok := state.BundleDict[name]; ok {
					v.BundleArr = append(v.BundleArr, r)
				}
			}
		}
		// 排序
		sort.Sort(SortResourceByName(state.ResourceArr))
		sort.Sort(SortBundleByName(state.BundleArr))
		sort.Sort(SortLoaderByName(state.LoaderArr))
	}

	return nil
}

func genResourcesYamlAndGo(state *gen_state) error {
	log.Info("生成yaml文件")
	if _, err := os.Stat(proj.Config.ResourceDir); os.IsNotExist(err) {
		if err := os.Mkdir(proj.Config.ResourceDir, 0755); err != nil {
			return err
		}
	}
	for name, v := range state.ResourceDict {
		log.Info(name, "==>", v.YamlFileName)
		filePath := path.Join(proj.Config.ResourceDir, v.YamlFileName)
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		file.Truncate(0)
		file.WriteString(fotmatYamlString(v))
		file.WriteString("\n")
		file.Close()
	}
	log.Info("生成go文件")
	if _, err := os.Stat(proj.Config.SrcGenResourceDir); err != nil && os.IsNotExist(err) {
		if err := os.Mkdir(proj.Config.SrcGenResourceDir, 0755); err != nil {
			return err
		}
	}
	funcMap := template.FuncMap{
		"join":        strings.Join,
		"quote":       QuoteChar,
		"capUpper":    capUpperString,
		"camelString": camelString,
	}
	resourcesPath := path.Join(proj.Config.SrcGenResourceDir, "resource.gen.go")
	file, err := os.OpenFile(resourcesPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	file.Truncate(0)
	defer file.Close()
	tmpl := template.New("resource").Delims("<<", ">>")
	tmpl.Funcs(funcMap)
	if tmpl, err := tmpl.Parse(code); err != nil {
		return err
	} else {
		if err := tmpl.Execute(file, state); err != nil {
			return err
		}
	}
	return nil
}

func genResourceCli(state *gen_state) error {
	if _, err := os.Stat(path.Join(proj.Config.SrcGenResourceDir, "bin")); err != nil && os.IsNotExist(err) {
		os.Mkdir(path.Join(proj.Config.SrcGenResourceDir, "bin"), 0755)
	}
	resourceFilePath := path.Join(proj.Config.SrcGenResourceDir, "bin", "resource.gen.go")
	file, err := os.OpenFile(resourceFilePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	file.Truncate(0)
	defer file.Close()
	tmpl := template.New("resource").Delims("<<", ">>")
	if tmpl, err := tmpl.Parse(cli_code); err != nil {
		return err
	} else {
		if err := tmpl.Execute(file, state); err != nil {
			return err
		}
	}
	return nil
}

type Config struct {
	Force bool
}

// 生成协议
func Gen(config Config) error {
	log.Info("===============gen resource start===============")
	// 初始化
	state := &gen_state{
		Config:       config,
		Module:       proj.Config.Module,
		ResourceDict: make(map[string]*Resource, 0),
		ResourceArr:  make([]*Resource, 0),
		BundleDict:   make(map[string]*Bundle, 0),
		BundleArr:    make([]*Bundle, 0),
		LoaderArr:    make([]*Loader, 0),
	}

	// for _, v := range state.ResourceArr {
	// 	filePath := path.Join(proj.Config.ExcelDir, v.FilePath)
	// 	srcFilePathArr = append(srcFilePathArr, filePath)
	// }
	// srcHash := getSrcFileHash(srcFilePathArr)
	// if srcHash == proj.Config.GenResourceHash {
	// 	return gira.ErrGenNotChange
	// }
	//for _, v := range state.ResourceArr {
	// 解析excel
	//filePath := path.Join(proj.Config.ExcelDir, v.FilePath)
	//if err := v.readExcel(filePath); err != nil {
	////	log.Println(filePath)
	// return err
	//}
	//}
	// proj.Update("gen_resource_hash", srcHash)

	if err := parse(state); err != nil && err == gira.ErrGenNotChange {
		log.Info("===============gen resource finished, not change===============")
		return nil
	} else if err != nil {
		return err
	}
	// 生成YAML和go
	if err := genResourcesYamlAndGo(state); err != nil {
		return err
	}
	// 生成cli程序
	if err := genResourceCli(state); err != nil {
		return err
	}
	log.Info("===============gen resource finished===============")
	return nil
}
