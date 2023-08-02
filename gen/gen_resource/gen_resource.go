package gen_resource

/// 参考 https://www.cnblogs.com/f-ck-need-u/p/10053124.html

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	log "github.com/lujingwei002/gira/corelog"
	"github.com/lujingwei002/gira/errors"
	"gopkg.in/yaml.v3"

	"github.com/lujingwei002/gira/proj"
	excelize "github.com/xuri/excelize/v2"
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
	"fmt"
	"github.com/lujingwei002/gira/db"
	"github.com/lujingwei002/gira/options/resource_options"
	"github.com/urfave/cli/v2"
	"<<.Module>>/gen/resource"
)

var resVersion string
var buildTime string
var uri string
var drop bool
var force bool

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
					&cli.BoolFlag{
						Name: "drop",
						Value: false,
						Usage: "drop collection",
						Destination: &drop,
					},
					&cli.BoolFlag{
						Name: "force",
						Value: false,
						Usage: "force overwrite collection",
						Destination: &force,
					},
				},
			},
			{
				Name:      "version",
				Usage:     "res version",
				Action:    versionAction,
			},
			{
				Name:      "time",
				Usage:     "Build time",
				Action:    timeAction,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Println(err)
	}
}

// 打印应用构建版本
func versionAction(args *cli.Context) error {
	fmt.Println(fmt.Sprintf("<<.ResVersion>>.%s", resVersion))
	return nil
}

// 打印应用构建时间
func timeAction(args *cli.Context) error {
	fmt.Println(buildTime)
	return nil
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
	opts := make([]resource_options.PushOption, 0)
	if drop {
		opts = append(opts, resource_options.WithPushDropOption())
	}
	if force {
		opts = append(opts, resource_options.WithPushForceOption())
	}
	if client, err := db.NewDbClientFromUri(context.Background(), "resource", uri); err != nil {
		return err
	} else {
		return resource.Push(context.Background(), client, "resource", opts...)
	}
}
`

var code = `
// Code generated by github.com/lujingwei002/gira. DO NOT EDIT.
// Code generated by github.com/lujingwei002/gira. DO NOT EDIT.
// Code generated by github.com/lujingwei002/gira. DO NOT EDIT.

package resource
import (
	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/errors"
	"github.com/lujingwei002/gira/options/resource_options"
	yaml "gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
	"encoding/gob"
	"path"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	log "github.com/lujingwei002/gira/corelog"
	"os"
	"fmt"
	"context"
<<range .ImportArr>>
	"<<.>>"
<<- end>>
)

// 将Db类型的bundle推送覆盖到db上
//
// Parameters:
// uri - mongodb://root:123456@192.168.1.200:3331/resourcedb
// dir - bundle所在的目录
func Push(ctx context.Context, client gira.DbClient, dir string, opts ...resource_options.PushOption) error {
	dir = path.Join(dir, "conf")
	<<- range .BundleArr>>
	<<- if  eq .BundleType "db">>
	<<.CapBundleStructName>> := &<<.BundleStructName>>{}
	// 推送<<.BundleName>>
	if err := <<.CapBundleStructName>>.SaveToDb(ctx, client, dir, opts...); err != nil {
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
	dir = path.Join(dir, "conf")
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

type version_file struct {
	ResVersion string <<quote>>yaml:"res_version"<<quote>>
}

<<- $config := .Config>>
<<- $state := .>>

<<range .LoaderArr>>
// <<.LoaderName>>加载器，负载加载拥有的bundle
type <<.LoaderStructName>> struct {
	handler gira.ResourceHandler
	version version_file
	<<- range .BundleArr>>
	<<.BundleStructName>>
	<<- end>>
}

func New<<.LoaderStructName>>(handler gira.ResourceHandler) *<<.LoaderStructName>>{
	return &<<.LoaderStructName>>{
		handler: handler,
	}
}

func (self *<<.LoaderStructName>>) GetResVersion() string {
	return self.version.ResVersion
}

func (self *<<.LoaderStructName>>) GetLoaderVersion() string {
	return "<<$state.ResVersion>>.<<$config.RespositoryVersion>>"
}

// 加载版本文件
func (self *<<.LoaderStructName>>) LoadVersion(dir string) error {
	filePath := filepath.Join(dir, ".version.yaml")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, &self.version); err != nil {
		log.Warnw("load resource version fail", "name", filePath)
		return err
	}
	return nil
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
func (self *<<.LoaderStructName>>) LoadFromDb(ctx context.Context, client gira.DbClient) error {
	<<- range .BundleArr>>
	if err := self.<<.BundleStructName>>.LoadFromDb(ctx, client); err != nil {
		return err
	}
	<<- end>>
	return nil
}

// 根据bundle的类型，从相应的源中加载资源
func (self *<<.LoaderStructName>>) LoadResource(ctx context.Context, client gira.DbClient, dir string, compress bool) error {
	if err := self.LoadVersion(dir); err != nil {
		return err
	}
	<<- range .BundleArr>>
	<<- if eq .BundleType "db">>
	if err := self.<<.BundleStructName>>.LoadFromDb(ctx, client); err != nil {
		return err
	}
	<<- end>>
	<<- if eq .BundleType "raw">>
	if compress {
		if err := self.<<.BundleStructName>>.LoadFromBin(dir); err != nil {
			return err
		}
	} else {
		if err := self.<<.BundleStructName>>.LoadFromYaml(dir); err != nil {
			return err
		}
	}
	<<- end>>
	<<- if eq .BundleType "bin">>
	if err := self.<<.BundleStructName>>.LoadFromBin(dir); err != nil {
		return err
	}
	<<- end>>
	<<- end>>
	if self.handler != nil {
		self.handler.OnResourcePreLoad(false)
		if err := self.convert(self.handler); err != nil {
			return err
		}
		self.handler.OnResourcePostLoad(false)
	}
	return nil
}

// 根据bundle的类型，从相应的源中加载资源
func (self *<<.LoaderStructName>>) ReloadResource(ctx context.Context, client gira.DbClient, dir string, compress bool) error {
	if err := self.LoadVersion(dir); err != nil {
		return err
	}
	<<- range .BundleArr>>
	<<- if eq .BundleType "db">>
	if err := self.<<.BundleStructName>>.LoadFromDb(ctx, client); err != nil {
		return err
	}
	<<- end>>
	<<- if eq .BundleType "raw">>
	if compress {
		if err := self.<<.BundleStructName>>.LoadFromBin(dir); err != nil {
			return err
		}
	} else {
		if err := self.<<.BundleStructName>>.LoadFromYaml(dir); err != nil {
			return err
		}
	}
	<<- end>>
	<<- if eq .BundleType "bin">>
	if err := self.<<.BundleStructName>>.LoadFromBin(dir); err != nil {
		return err
	}
	<<- end>>
	<<- end>>
	if self.handler != nil {
		self.handler.OnResourcePreLoad(true)
		if err := self.convert(self.handler); err != nil {
			return err
		}
		self.handler.OnResourcePostLoad(true)
	}
	return nil
}

// 加载成功后，对配置进行处理
func (self *<<.LoaderStructName>>) convert(handler gira.ResourceHandler) error {
	h := handler.(<<.HandlerStructName>>)
	<<- range .BundleArr>>
		<<- $bundleStructName := .BundleStructName>>
		<<- range .ResourceArr>>
	if err := h.Convert<<.StructName>>(self.<<$bundleStructName>>.<<.WrapStructName>>); err != nil {
		return err
	}
		<<- end>>
	<<- end>>
	
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
	<<- $bundleName := .BundleName>>
	<<- range .ResourceArr>>
	var <<.CapStructName>>filePath = filepath.Join(dir, "<<$bundleName>>", "<<.YamlFileName>>")
	if err := self.<<.WrapStructName>>.LoadFromYaml(<<.CapStructName>>filePath); err != nil {
		return err
	}
	<<end>>
	return nil
}

// 从db中加载资源
func (self* <<.BundleStructName>>) LoadFromDb(ctx context.Context, client gira.DbClient) error {
	self.Clear()
	<<- range .ResourceArr>>
	if err := self.<<.WrapStructName>>.LoadFromDb(ctx, client); err != nil {
		return err
	}
	<<- end>>
	return nil
}

// 从bin文件中加载资源
func (self* <<.BundleStructName>>) LoadFromBin(dir string) error {
	var filePath = filepath.Join(dir, "<<.BundleName>>.dat")
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
	var filePath = filepath.Join(dir, "<<.BundleName>>.dat")
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

func (self *<<.BundleStructName>>) SaveToDb(ctx context.Context, client gira.DbClient, dir string, opts ...resource_options.PushOption) error {
	switch c := client.(type) {
	case gira.MongoClient:
		return self.SaveToMongo(ctx, c.GetMongoDatabase(), dir, opts...)
	default:
		return errors.ErrDbNotSupport
	}
}

// 将资源保存到db上
func (self *<<.BundleStructName>>) SaveToMongo(ctx context.Context, database *mongo.Database, dir string, opts ...resource_options.PushOption) error {
	pushOptions := &resource_options.PushOptions {
	}
	for _, v := range opts {
		v.ConfigPushOption(pushOptions)
	}
	<<- $bundleName := .BundleName >>
	<<- range .ResourceArr>>
	// 加载<<.ResourceName>>
	var <<.CapStructName>>filePath = filepath.Join(dir, "<<$bundleName>>", "<<.YamlFileName>>")
	var <<.ArrTypeName>> <<.ArrTypeName>>
	if err := <<.ArrTypeName>>.LoadFromYaml(<<.CapStructName>>filePath); err != nil {
		return err
	}
	<<- end>>
	<<- range .ResourceArr>>
	// 保存 <<.ResourceName>>
	if coll := database.Collection("<<.TableName>>"); coll == nil {
		return fmt.Errorf("collection <<.TableName>> not found")
	} else if pushOptions.Drop {
		coll.Drop(ctx)
	} else if pushOptions.Force {
		coll.Drop(ctx)
		models := make([]mongo.WriteModel, 0)
		for _, v := range <<.ArrTypeName>> {
			models = append(models, mongo.NewInsertOneModel().SetDocument(v))
		}
		if len(models) > 0 {
			if _, err := coll.BulkWrite(ctx, models); err != nil {
				log.Info("push <<.TableName>> force fail", err)
				return err
			} else {
				log.Infof("push <<.TableName>>(%d) force success", len(models))
			} 
		} else {
			log.Infof("push <<.TableName>>(%d) force success", len(models))
		}
	} else {
		models := make([]mongo.WriteModel, 0)
		<<- if gt .KeyLen 0 >>
		<<- $key0 := .Key0 >>
		for _, v := range <<.ArrTypeName>> {
			doc := &<<.StructName>>{}
			err := coll.FindOne(ctx, bson.D{{"<< $key0 >>", v.<<camelString $key0>>}}).Decode(doc)
			if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
				log.Infof("insert <<.TableName>> %v", v.<<camelString $key0>>)
				models = append(models, mongo.NewInsertOneModel().SetDocument(v))
			} else if err != nil {
				return err
			}
		}
		<<- end>>
		if len(models) > 0 {
			if _, err := coll.BulkWrite(ctx, models); err != nil {
				log.Info("push <<.TableName>> fail", err)
				return err
			} else {
				log.Infof("push <<.TableName>>(%d) success", len(models))
			}
		} else {
			log.Infof("push <<.TableName>>(%d) success", len(models))
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
func (self *<<.ArrTypeName>>) LoadFromMongo(ctx context.Context, database *mongo.Database) error {
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

func (self *<<.ArrTypeName>>) LoadFromDb(ctx context.Context, client gira.DbClient) error {
	switch c := client.(type) {
	case gira.MongoClient:
		return self.LoadFromMongo(ctx, c.GetMongoDatabase())
	default:
		panic(errors.ErrDbNotSupport)
	}
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

func (self *<<.WrapStructName>>) LoadFromDb(ctx context.Context, client gira.DbClient) error {
	switch c := client.(type) {
	case gira.MongoClient:
		return self.LoadFromMongo(ctx, c.GetMongoDatabase())
	default:
		panic(errors.ErrDbNotSupport)
	}
}

func (self *<<.WrapStructName>>) LoadFromMongo(ctx context.Context, database *mongo.Database) error {
	var arr <<.ArrTypeName>>
	if err := arr.LoadFromMongo(ctx, database); err != nil {
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


func (self *<<.WrapStructName>>) LoadFromDb(ctx context.Context, client gira.DbClient) error {
	switch c := client.(type) {
	case gira.MongoClient:
		return self.LoadFromMongo(ctx, c.GetMongoDatabase())
	default:
		panic(errors.ErrDbNotSupport)
	}
}

func (self *<<.WrapStructName>>) LoadFromMongo(ctx context.Context, database *mongo.Database) error {
	var arr <<.ArrTypeName>>
	if err := arr.LoadFromMongo(ctx, database); err != nil {
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

var xml_template = `
<!--Code generated by github.com/lujingwei002/gira. DO NOT EDIT.-->
<!--Code generated by github.com/lujingwei002/gira. DO NOT EDIT.-->
<!--Code generated by github.com/lujingwei002/gira. DO NOT EDIT.-->




<collections>

<<- range .ResourceArr>> 
<</* 模型Data */>>
	<collection name="<<.ResourceName>>">
	<<- range .FieldArr>> 
		<!--<<.Comment>>-->
		<field name="<<.FieldName>>" type="<<.GoTypeName>>"/>
	<<- end>>
	</collection>
<<- end>> 
</collections>

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
	field_type_float
	field_type_object_id
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
	"float":    field_type_float,
	"int64":    field_type_int64,
	"long":     field_type_int64,
	"int32":    field_type_int32,
	"string":   field_type_string,
	"json":     field_type_json,
	"bool":     field_type_bool,
	"string[]": field_type_string_arr,
	"int[]":    field_type_int_arr,
	"float[]":  field_type_float_arr,
	"id":       field_type_object_id,
}

// 和go类型的对应关系
var go_type_name_dict = map[field_type]string{
	field_type_int:        "int64",
	field_type_float:      "float64",
	field_type_int64:      "int64",
	field_type_int32:      "int32",
	field_type_string:     "string",
	field_type_json:       "interface{}",
	field_type_bool:       "bool",
	field_type_string_arr: "[]string",
	field_type_int_arr:    "[]int64",
	field_type_float_arr:  "[]float64",
	field_type_object_id:  "primitive.ObjectID",
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
}

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
	KeyArr         []string

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

func (self *Resource) Key0() string {
	return self.KeyArr[0]
}

func (self *Resource) KeyLen() int {
	return len(self.KeyArr)
}

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

type Loader struct {
	LoaderStructName  string
	LoaderName        string
	HandlerStructName string
	bundleNameArr     []string
	BundleArr         []*Bundle
}

// 生成协议的状态
type gen_state struct {
	Config       Config
	Module       string
	ResourceDict map[string]*Resource
	ResourceArr  []*Resource
	BundleDict   map[string]*Bundle
	BundleArr    []*Bundle
	LoaderArr    []*Loader
	ImportArr    []string
	ResVersion   string
}

type Parser interface {
	parse(constState *gen_state) error
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

type import_type []string

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
					FieldName:  v,
					Type:       realType,
					Tag:        index,
					GoTypeName: go_type_name_dict[realType],
					Comment:    comment,
				}
				if v == "_id" {
					field.StructFieldName = "Id"
				} else {
					field.StructFieldName = capUpperString(camelString(v))
				}
				r.FieldArr = append(r.FieldArr, field)
				r.FieldDict[v] = field
			} else {
				log.Warnw("invalid type", "file_path", filePath, "type", typeName, "field", v)
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

type version_file struct {
	ResVersion string `yaml:"res_version"`
}

func genResourcesVersion(state *gen_state) error {
	log.Info("生成version文件")
	filePath := filepath.Join(proj.Dir.ResourceDir, "conf", ".version.yaml")
	v := version_file{}
	v.ResVersion = fmt.Sprintf("%s.%s", state.ResVersion, state.Config.RespositoryVersion)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	file.Truncate(0)
	if data, err := yaml.Marshal(&v); err != nil {
		return err
	} else {
		if _, err := file.Write(data); err != nil {
			return err
		}
	}
	return nil
}

func genXml(state *gen_state) error {
	log.Info("生成xml文件")
	if _, err := os.Stat(proj.Dir.GenResourceDir); os.IsNotExist(err) {
		if err := os.Mkdir(proj.Dir.GenResourceDir, 0755); err != nil {
			return err
		}
	}
	for _, bundle := range state.BundleArr {
		if bundle.BundleType == "db" {
			fmt.Println(bundle.BundleName)
			var err error
			filePath := path.Join(proj.Dir.GenResourceDir, fmt.Sprintf("%s.xml", bundle.BundleName))
			file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				return err
			}
			file.Truncate(0)
			funcMap := template.FuncMap{
				"quote": QuoteChar,
			}
			tmpl := template.New("xml").Delims("<<", ">>")
			tmpl.Funcs(funcMap)
			tmpl, err = tmpl.Parse(xml_template)
			if err != nil {
				return err
			}
			if err = tmpl.Execute(file, bundle); err != nil {
				return err
			}
			file.Close()
		}
	}
	return nil
}

func genResourcesYaml(state *gen_state) error {
	log.Info("生成yaml文件")
	if _, err := os.Stat(proj.Dir.ResourceDir); os.IsNotExist(err) {
		if err := os.Mkdir(proj.Dir.ResourceDir, 0755); err != nil {
			return err
		}
	}
	for _, bundle := range state.BundleArr {
		for name, v := range bundle.ResourceArr {
			log.Info(name, "==>", path.Join(bundle.BundleName, v.YamlFileName))
			if _, err := os.Stat(path.Join(proj.Dir.ResourceDir, "conf")); os.IsNotExist(err) {
				if err := os.Mkdir(path.Join(proj.Dir.ResourceDir, "conf"), 0755); err != nil {
					return err
				}
			}
			if _, err := os.Stat(path.Join(proj.Dir.ResourceDir, "conf", bundle.BundleName)); os.IsNotExist(err) {
				if err := os.Mkdir(path.Join(proj.Dir.ResourceDir, "conf", bundle.BundleName), 0755); err != nil {
					return err
				}
			}
			filePath := path.Join(proj.Dir.ResourceDir, "conf", bundle.BundleName, v.YamlFileName)
			file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				return err
			}
			file.Truncate(0)
			file.WriteString(fotmatYamlString(v))
			file.WriteString("\n")
			file.Close()
		}
	}
	return nil

	// for name, v := range state.ResourceDict {
	// 	log.Info(name, "==>", v.YamlFileName)
	// 	filePath := path.Join(proj.Dir.ResourceDir, v.YamlFileName)
	// 	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	file.Truncate(0)
	// 	file.WriteString(fotmatYamlString(v))
	// 	file.WriteString("\n")
	// 	file.Close()
	// }
	return nil
}

func genResourcesLoader(state *gen_state) error {
	log.Info("生成go文件")
	if _, err := os.Stat(proj.Dir.SrcGenResourceDir); err != nil && os.IsNotExist(err) {
		if err := os.Mkdir(proj.Dir.SrcGenResourceDir, 0755); err != nil {
			return err
		}
	}
	funcMap := template.FuncMap{
		"join":        strings.Join,
		"quote":       QuoteChar,
		"capUpper":    capUpperString,
		"camelString": camelString,
	}
	resourcesPath := path.Join(proj.Dir.SrcGenResourceDir, "resource.gen.go")
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
	if _, err := os.Stat(path.Join(proj.Dir.SrcGenResourceDir, "bin")); err != nil && os.IsNotExist(err) {
		os.Mkdir(path.Join(proj.Dir.SrcGenResourceDir, "bin"), 0755)
	}
	resourceFilePath := path.Join(proj.Dir.SrcGenResourceDir, "bin", "resource.gen.go")
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
	RespositoryVersion string
	Force              bool
}

// 生成协议
func Gen(config Config) error {
	log.Info("===============gen resource start===============")
	// 初始化
	state := &gen_state{
		Config:       config,
		Module:       proj.Module,
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
	// 	return errors.ErrGenNotChange
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
	var p Parser
	if true {
		p = &golang_parser{}
	} else {
		p = &yaml_parser{}
	}
	if err := p.parse(state); err != nil && err == errors.ErrGenNotChange {
		log.Info("===============gen resource finished, not change===============")
		return nil
	} else if err != nil {
		return err
	}
	// 生成xml
	if err := genXml(state); err != nil {
		return err
	}
	// 生成YAML
	if err := genResourcesYaml(state); err != nil {
		return err
	}
	// 生成go
	if err := genResourcesLoader(state); err != nil {
		return err
	}
	// 生成cli程序
	if err := genResourceCli(state); err != nil {
		return err
	}
	// 生成.version
	if err := genResourcesVersion(state); err != nil {
		return err
	}
	fmt.Println("=================", config.RespositoryVersion)
	log.Info("===============gen resource finished===============")
	return nil
}
