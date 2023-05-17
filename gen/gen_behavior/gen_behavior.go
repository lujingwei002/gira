package gen_behavior

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/lujingwei002/gira/log"

	"github.com/lujingwei002/gira/proj"

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
	"github.com/lujingwei002/gira/behavior"
	"github.com/lujingwei002/gira/db"
	"github.com/urfave/cli/v2"
	"<<.Module>>/gen/behavior/<<.DbName>>"
)

var enabledDropIndex bool
var uri string
var connectTimeout int64

func main() {
	app := &cli.App{
		Name: "migrate-<<.DbName>>",
		Authors: []*cli.Author{{
			Name:  "lujingwei",
			Email: "lujingwei@xx.org",
		}},
		Description: "migrate-<<.DbName>>",
		Flags:       []cli.Flag{},
		Commands: []*cli.Command{
			{
				Name:   "migrate",
				Usage:  "migrate scheme to database",
				Action: migrateAction,
				Flags:       []cli.Flag{
					&cli.BoolFlag{
						Name: "drop-index",
						Value: false,
						Usage: "enable drop index",
						Destination: &enabledDropIndex,
					},
					&cli.StringFlag{
						Name: "uri",
						Required: true,
						Usage: "database uri",
						Destination: &uri,
					},
					&cli.Int64Flag{
						Name: "connect-timeout",
						Value: 5,
						Usage: "connect database timeout",
						Destination: &connectTimeout,
					},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Println(err)
	}
}

func migrateAction(args *cli.Context) error {
	opts := make([]behavior.MigrateOption, 0)
	opts = append(opts, behavior.WithMigrateDropIndex(enabledDropIndex), behavior.WithMigrateConnectTimeout(connectTimeout))
	if client, err := db.NewDbClientFromUri(context.Background(), "<<.DbName>>", uri); err != nil {
		return err
	} else {
		return <<.DbName>>.Migrate(context.Background(), client, opts...)
	}
}
`

var model_template = `
// Code generated by github.com/lujingwei002/gira. DO NOT EDIT.
// Code generated by github.com/lujingwei002/gira. DO NOT EDIT.
// Code generated by github.com/lujingwei002/gira. DO NOT EDIT.

package <<.DbName>> 

// mongo模型
import (
	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/log"
	"github.com/lujingwei002/gira/facade"
	"github.com/lujingwei002/gira/behavior"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	yaml "gopkg.in/yaml.v3"
	"context"
	"bufio"
	"os"
	"sync"
	"io"
	"time"
)





<<- range .CollectionArr>> 
<</* 模型Data */>>
// <<.Comment>>
type <<.StructName>> struct {
	<<- range .FieldArr>> 
	/// <<.Comment>>
	<<.CamelName>> <<.GoTypeName>> <<quote>>bson:"<<.Name>>" json:"<<.Name>>"<<quote>>
	<<- end>>
}
<<- end>> 

type <<.DriverInterfaceName>> interface {
	Sync(ctx context.Context, opts ...behavior.SyncOption) (err error) 
	Migrate(ctx context.Context, opts ...behavior.MigrateOption) error
	<<- range .CollectionArr>> 
	Log<<.StructName>>(doc *<<.StructName>>) error
	<<- end>>
}

// mongo 
type <<.MongoDriverStructName>> struct {
	client		*mongo.Client
	database	*mongo.Database
	<<range .CollectionArr>> 
	// <<.Comment>>
	<<.StructName>>  *<<.MongoDaoStructName>>
	<<- end>>
}

var globalDriver <<.DriverInterfaceName>>

func NewMongo() *<<.MongoDriverStructName>> {
	self := &<<.MongoDriverStructName>>{}
	<<- range .CollectionArr>> 
	self.<<.StructName>> = &<<.MongoDaoStructName>>{
		db: self,
		models: make([]mongo.WriteModel, 0),
	}
	<<- end>> 
	return self
}

func Migrate(ctx context.Context, client  gira.DbClient, opts ...behavior.MigrateOption) error {
	migrateOptions := &behavior.MigrateOptions {
	}
	for _, v := range opts {
		v.ConfigMigrateOptions(migrateOptions)
	}
	switch client2 := client.(type) {
	case gira.MongoClient:
		driver := NewMongo()
		driver.Use(client2)
		return driver.Migrate(ctx, opts...)
	default:
		return gira.ErrDbNotSupport.Trace()
	}
}

func Use(ctx context.Context, client gira.DbClient, config gira.BehaviorConfig) error {
	switch client2 := client.(type) {
	case gira.MongoClient:
		return UseMongo(ctx, client2, config)
	default:
		return gira.ErrDbNotSupport.Trace()
	}
}

func UseMongo(ctx context.Context, client gira.MongoClient, config gira.BehaviorConfig) error {
	if globalDriver != nil {
		return gira.ErrTodo
	}
	driver := NewMongo()
	if err := driver.Use(client); err != nil {
		return err
	}
	globalDriver = driver
	facade.Go(func() error {
		return driver.Serve(ctx, config)
	})
	return nil
}

func Sync(ctx context.Context, opts ...behavior.SyncOption) error {
	if globalDriver == nil {
		return gira.ErrBehaviorNotInit
	}
	return globalDriver.Sync(ctx, opts...)
}

<<- range .CollectionArr>>
// <<.Comment>>
func Log<<.StructName>>(doc *<<.StructName>>) error {
	return globalDriver.Log<<.StructName>>(doc)
}
<<- end>> 

func (self *<<.MongoDriverStructName>>) Use(client gira.MongoClient) error {
	if self.client != nil {
		return gira.ErrTodo
	}
	self.client = client.GetMongoClient()
	self.database = client.GetMongoDatabase()
	return nil
}

func (self *<<.MongoDriverStructName>>) Migrate(ctx context.Context, opts ...behavior.MigrateOption) error {
<<- range .CollectionArr>>
	if err := self.<<.StructName>>.Migrate(ctx, opts...); err != nil {
		return err
	}
<<- end>> 
	return nil
}

func (self *<<.MongoDriverStructName>>) Serve(ctx context.Context, config gira.BehaviorConfig) (err error) {
	ticker := time.NewTicker(time.Duration(config.SyncInterval)*time.Second)
	defer func() {
		ticker.Stop()
	}()
	opts := make([]behavior.SyncOption, 0)
	if config.BatchInsert != 0 {
		opts = append(opts, behavior.WitchBatchInsertOption(config.BatchInsert))
	}
	for {
		select {
		case <-ctx.Done():
			self.Sync(context.TODO())
			return nil
		case <-ticker.C:
			self.Sync(ctx, opts...)
		}
	}
}

func (self *<<.MongoDriverStructName>>) Sync(ctx context.Context, opts ...behavior.SyncOption) (err error) {
<<- range .CollectionArr>>
	if _, err = self.<<.StructName>>.Sync(ctx, opts...); err != nil {
		log.Error(err)
	} else {

	}
<<- end>> 
	return 
}

<<- range .CollectionArr>> 
// <<.Comment>>
func (self *<<.MongoDriverStructName>>) Log<<.StructName>>(doc *<<.StructName>>) error {
	return self.<<.StructName>>.Log(doc)
}
<<- end>> 

<<- range .CollectionArr>> 

// <<.Comment>>
type <<.MongoDaoStructName>> struct {
	db 		*<<$.MongoDriverStructName>>
	models	[]mongo.WriteModel
	mu		sync.Mutex
}



// <<.Comment>>
func (self *<<.MongoDaoStructName>>) Log(doc *<<.StructName>>) error {
	doc.Id = primitive.NewObjectID()
	<<- if .HasLogTimeField>>
	doc.LogTime = time.Now().Unix()
	<<- end>>
	log.Infow("<<.CollName>>",
	<<- range .FieldArr>> 
		"<<.Name>>", doc.<<.CamelName>>, // <<.Comment>>
	<<- end>>
	)
	self.mu.Lock()
	self.models = append(self.models, mongo.NewInsertOneModel().SetDocument(doc))
	self.mu.Unlock()
	return nil
}

func (self *<<.MongoDaoStructName>>) Migrate(ctx context.Context, opts ...behavior.MigrateOption) error {
	migrateOptions := &behavior.MigrateOptions {
	}
	for _, v := range opts {
		v.ConfigMigrateOptions(migrateOptions)
	}
	database := self.db.database
	coll := database.Collection("<<.CollName>>")
	indexView := coll.Indexes()
	listOpts := options.ListIndexes().SetMaxTime(2 * time.Second)
	cursor, err := indexView.List(ctx, listOpts)
	if err != nil {
		return err
	}
	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return err
	}
	collName := "<<.CollName>>"
	log.Printf("%s migrate index", collName)
	// 已经有的
	own := make(map[string]bson.M)
	for _, v := range results {
		own[v["name"].(string)] = v
		log.Printf("[ ]%s.%s", collName, v["name"].(string))
	}
	// 配置的
	indexes := map[string]bool {
	<<- range .IndexArr>> 
		"<<.FullName>>": true,
	<<- end>>
	}

	// 新增索引
	<<- range .IndexArr>> 
	if _, ok := own["<<.FullName>>"]; !ok {
		keys := bson.D{
			<<- range .KeyArr>> 
			{Key: "<<.Key>>", Value: <<.Value>>},
			<<- end>>
		}
		log.Printf("[+]%s.<<.FullName>>", collName)
		if _, err := indexView.CreateOne(ctx, mongo.IndexModel{
            Keys: keys,
        }); err != nil {
			return err
		}
	}
	<<- end>>
	// 删除索引
	for name, _ := range own {
		if name == "_id_" {
			continue
		}
		if _, ok := indexes[name]; !ok {
			if migrateOptions.EnabledDropIndex {
				log.Printf("[-]%s.%s", collName, name)
				if _, err := indexView.DropOne(ctx, name); err != nil {
					return err
				}
			} else {
				log.Printf("[*]%s.%s", collName, name)
			}
		}
	}
	log.Println()
	// log.Println(own)
	return nil
}

func (self *<<.MongoDaoStructName>>) Sync(ctx context.Context, opts ...behavior.SyncOption) (n int, err error) {
	syncOptions := &behavior.SyncOptions {
	}
	for _, v := range opts {
		v.ConfigSyncOptions(syncOptions)
	}
	database := self.db.database
	coll := database.Collection("<<.CollName>>")
	if len(self.models) <= 0 {
		return
	}
	self.mu.Lock()
	var models []mongo.WriteModel
	if syncOptions.Step != 0 && syncOptions.Step <= len(self.models) {
		models = self.models[:syncOptions.Step]
		self.models = self.models[syncOptions.Step:]
	} else {
		models = self.models
		self.models = make([]mongo.WriteModel, 0)
	}
	self.mu.Unlock()
	<<- if .HasCreateTimeField>>
	for _, model := range models {
		v := model.(*mongo.InsertOneModel)
		doc := v.Document.(*<<.StructName>>)
		doc.CreateTime = time.Now().Unix()
	}
	<<- end>>
	writeOpts := options.BulkWrite().SetOrdered(false)
	_, err = coll.BulkWrite(ctx, models, writeOpts)
	if err != nil {
		self.mu.Lock()
		self.models = append(self.models, models...)
		self.mu.Unlock()
		log.Errorw("sync behavior fail", "name", "<<.CollName>>", "len", len(models), "error", err)
	    return
	} else {
		log.Infow("sync behavior", "name", "<<.CollName>>", "len", len(models))
		n = len(models)
		return 
	}
}

func (self *<<.MongoDaoStructName>>) BatchWrite(ctx context.Context, filePath string) error {
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
    if err != nil {
        return err
    }
	database := self.db.database
	coll := database.Collection("<<.CollName>>")
	models := make([]mongo.WriteModel, 0)
    defer f.Close()
    reader := bufio.NewReader(f)
    // 按行处理txt
    for  {
        line, _, err := reader.ReadLine()
        if err == io.EOF {
            break
        }
		doc := <<.StructName>> {
		}
		if err := yaml.Unmarshal(line, &doc); err != nil {
			return err
		}
		models = append(models, mongo.NewInsertOneModel().SetDocument(&doc))
    }
	if len(models) <= 0 {
		return nil
	}
	opts := options.BulkWrite().SetOrdered(false)
	_, err = coll.BulkWrite(ctx, models, opts)
	if err != nil {
		log.Error(err)
	    return err
	}
	return nil
}
<<- end>>





`

type field_type int

const (
	field_type_int field_type = iota
	field_type_int32
	field_type_int64
	field_type_string
	field_type_message
	field_type_objectid
	field_type_bool
	field_type_bytes
	field_type_int_arr
	field_type_int64_arr
	field_type_struct
)

const field_id_name string = "id"

var type_name_dict = map[string]field_type{
	"int":     field_type_int,
	"int32":   field_type_int32,
	"int64":   field_type_int64,
	"string":  field_type_string,
	"id":      field_type_objectid,
	"bool":    field_type_bool,
	"bytes":   field_type_bytes,
	"[]int":   field_type_int_arr,
	"[]int64": field_type_int64_arr,
}

var go_type_name_dict = map[field_type]string{
	field_type_int:       "int",
	field_type_int32:     "int32",
	field_type_int64:     "int64",
	field_type_string:    "string",
	field_type_objectid:  "primitive.ObjectID",
	field_type_bool:      "bool",
	field_type_bytes:     "[]byte",
	field_type_int_arr:   "[]int64",
	field_type_int64_arr: "[]int64",
}

type Field struct {
	Tag        int
	Name       string
	CamelName  string
	Type       field_type
	TypeName   string
	GoTypeName string
	Default    interface{}
	Comment    string
	Coll       *Collection
}
type SortFieldByName []*Field

func (self SortFieldByName) Len() int           { return len(self) }
func (self SortFieldByName) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }
func (self SortFieldByName) Less(i, j int) bool { return self[i].Tag < self[j].Tag }

func (f *Field) IsComparable() bool {
	switch f.GoTypeName {
	case "int32":
	case "int":
	case "int64":
	case "string":
	case "bool":
	case "primitive.ObjectID":
		return true
	default:
		return false
	}
	return false
}

func capLowerString(s string) string {
	if len(s) <= 0 {
		return s
	}
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

type Key struct {
	Key   string
	Value interface{}
}
type Index struct {
	Name     string
	KeyDict  map[string]*Key
	KeyArr   []*Key
	Tag      int
	FullName string
}
type SortIndexByName []*Index

func (self SortIndexByName) Len() int           { return len(self) }
func (self SortIndexByName) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }
func (self SortIndexByName) Less(i, j int) bool { return self[i].Tag < self[j].Tag }

type Collection struct {
	CollName              string // 表名
	StructName            string // 表名的驼峰格式
	MongoDaoStructName    string // mongo dao 结构的名称
	Derive                string
	Comment               string
	FieldDict             map[string]*Field
	FieldArr              []*Field
	IndexDict             map[string]*Index
	IndexArr              []*Index
	HasLogTimeField       bool
	HasCreateTimeField    bool
	MongoDriverStructName string
}
type SortCollectionByName []*Collection

func (self SortCollectionByName) Len() int           { return len(self) }
func (self SortCollectionByName) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }
func (self SortCollectionByName) Less(i, j int) bool { return self[i].CollName < self[j].CollName }

type Database struct {
	Module                string
	Driver                string
	DbStructName          string // 数据库名的驼峰格式
	MongoDriverStructName string // mongo 的 dao 结构名字
	DriverInterfaceName   string
	DbName                string
	GenModelDir           string        // 生成的文件路径，在 gen/{{DbName}}
	GenModelFilePath      string        // 生成的文件路径，在 gen/{{DbName}}/{{DbName}}.gen.go
	GenBinFilePath        string        // 生成的文件路径，在 gen/{{DbName}}/bin/{{DbName}}.gen.go
	GenBinDir             string        // 生成的文件路径，在 gen/{{DbName}}/bin
	CollectionArr         []*Collection // 所有的模型
}

// 生成协议的状态
type gen_state struct {
	databaseArr []*Database
}

func QuoteChar() interface{} {
	return "`"
}

func (coll *Collection) IsDeriveUser() bool {
	return coll.Derive == "user"
}

func (coll *Collection) IsDeriveUserArr() bool {
	return coll.Derive == "userarr"
}

func (coll *Collection) parseIndex(attrs map[string]interface{}) error {
	coll.IndexDict = make(map[string]*Index)
	coll.IndexArr = make([]*Index, 0)
	equalRegexp := regexp.MustCompile("[^=]+")

	for valueStr, v := range attrs {
		var indexName string
		var tagStr string
		var tag int
		var err error
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
		indexName = args[0]
		index := &Index{
			Name: indexName,
			Tag:  tag,
		}
		for _, option := range optionArr {
			optionDict := option.(map[string]interface{})
			if keyArr, ok := optionDict["key"]; ok {
				var fullName string
				index.KeyDict = make(map[string]*Key)
				if keyArr2, ok := keyArr.([]interface{}); ok {
					for _, keyObject := range keyArr2 {
						if keyObject2, ok := keyObject.(map[string]interface{}); ok {
							for k, v := range keyObject2 {
								indexKey := &Key{Key: k, Value: v}
								index.KeyDict[k] = indexKey
								index.KeyArr = append(index.KeyArr, indexKey)
								if len(index.KeyDict) == 1 {
									fullName = fmt.Sprintf("%s_%v", k, v)
								} else {
									fullName = fmt.Sprintf("%s_%s_%v", fullName, k, v)
								}
							}
						}
					}
				}
				index.FullName = fullName
			}
		}
		// log.Println(index.KeyDict)
		// log.Println(index.FullName)
		coll.IndexDict[indexName] = index
		coll.IndexArr = append(coll.IndexArr, index)
	}
	sort.Sort(SortIndexByName(coll.IndexArr))
	return nil
}

func (coll *Collection) parseStruct(attrs map[string]interface{}) error {
	coll.FieldDict = make(map[string]*Field)
	coll.FieldArr = make([]*Field, 0)
	spaceRegexp := regexp.MustCompile("[^\\s]+")
	equalRegexp := regexp.MustCompile("[^=]+")

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
		typeStr = args[1]
		fieldName = args[0]
		field := &Field{
			Coll:      coll,
			Name:      fieldName,
			CamelName: camelString(fieldName),
			Tag:       tag,
		}
		if fieldName == "id" {
			field.Name = "_id"
			field.CamelName = "Id"
		} else {
			field.Name = fieldName
			field.CamelName = camelString(fieldName)
		}
		field.TypeName = typeStr
		if typeValue, ok := type_name_dict[typeStr]; ok {
			field.Type = typeValue
			field.GoTypeName = go_type_name_dict[field.Type]
		} else {
			field.Type = field_type_struct
			field.GoTypeName = typeStr
		}
		// fmt.Println(optionArr)
		for _, option := range optionArr {
			optionDict := option.(map[string]interface{})
			if defaultVal, ok := optionDict["default"]; ok {
				field.Default = defaultVal
			}
			if comment, ok := optionDict["comment"]; ok {
				field.Comment = comment.(string)
			}
		}
		coll.FieldDict[fieldName] = field
		coll.FieldArr = append(coll.FieldArr, field)
	}
	for _, f := range coll.FieldArr {
		if f.Name == "log_time" && f.Type == field_type_int64 {
			f.Coll.HasLogTimeField = true
		}
		if f.Name == "create_time" && f.Type == field_type_int64 {
			f.Coll.HasCreateTimeField = true
		}
	}
	sort.Sort(SortFieldByName(coll.FieldArr))
	return nil
}

func (coll *Collection) Unmarshal(genState *gen_state, v interface{}) error {
	var derive string
	row := v.(map[string]interface{})
	if _, ok := row["struct"]; !ok {
		return fmt.Errorf("collection %s struct part not found", coll.CollName)
	}
	structPart := row["struct"]
	if _, ok := structPart.(map[string]interface{}); !ok {
		return fmt.Errorf("collection %s struct part not map", coll.CollName)
	}
	coll.Derive = derive
	if err := coll.parseStruct(row["struct"].(map[string]interface{})); err != nil {
		return err
	}
	// 解析index
	if v, ok := row["index"]; !ok {
	} else if v2, ok := v.(map[string]interface{}); !ok {
	} else if err := coll.parseIndex(v2); err != nil {
		return err
	}
	if v, ok := row["comment"]; ok {
		coll.Comment = v.(string)
	}
	return nil
}

func parse(state *gen_state, filePathArr []string) error {
	for _, fileName := range filePathArr {
		filePath := path.Join(proj.Config.DocBehaviorDir, fileName)
		log.Info("处理文件", filePath)
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}
		dbName := strings.Replace(fileName, ".yaml", "", 1)
		database := &Database{
			Module:                proj.Config.Module,
			CollectionArr:         make([]*Collection, 0),
			GenModelDir:           path.Join(proj.Config.SrcGenBehaviorDir, dbName),
			GenModelFilePath:      path.Join(proj.Config.SrcGenBehaviorDir, dbName, fmt.Sprintf("%s.gen.go", dbName)),
			GenBinFilePath:        path.Join(proj.Config.SrcGenBehaviorDir, dbName, "bin", fmt.Sprintf("%s.gen.go", dbName)),
			GenBinDir:             path.Join(proj.Config.SrcGenBehaviorDir, dbName, "bin"),
			DbName:                dbName,
			DbStructName:          camelString(dbName),
			MongoDriverStructName: fmt.Sprintf("%sMongoDriver", camelString(dbName)),
			DriverInterfaceName:   fmt.Sprintf("%sDriver", camelString(dbName)),
		}
		result := make(map[string]interface{})
		if err := yaml.Unmarshal(data, result); err != nil {
			return err
		}
		for k, v := range result {
			if k == "$driver" {
				database.Driver = v.(string)
			} else {
				collName := k
				coll := &Collection{
					MongoDriverStructName: database.MongoDriverStructName,
					CollName:              collName,
					StructName:            camelString(collName),
					MongoDaoStructName:    fmt.Sprintf("%sMongoDao", camelString(collName)),
				}
				if err := coll.Unmarshal(state, v); err != nil {
					return err
				}
				database.CollectionArr = append(database.CollectionArr, coll)
			}
		}
		sort.Sort(SortCollectionByName(database.CollectionArr))
		state.databaseArr = append(state.databaseArr, database)
	}
	return nil
}

func genModel(state *gen_state) error {
	for _, db := range state.databaseArr {
		if _, err := os.Stat(db.GenModelDir); err != nil {
			if os.IsNotExist(err) {
				if err := os.Mkdir(db.GenModelDir, 0755); err != nil {
					return err
				}
			} else {
				return err
			}
		}
		var err error
		file, err := os.OpenFile(db.GenModelFilePath, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		file.Truncate(0)
		funcMap := template.FuncMap{
			"quote": QuoteChar,
		}
		tmpl := template.New("model").Delims("<<", ">>")
		tmpl.Funcs(funcMap)
		tmpl, err = tmpl.Parse(model_template)
		if err != nil {
			return err
		}
		if err = tmpl.Execute(file, db); err != nil {
			return err
		}
		file.Close()
	}
	return nil
}

func genCli(state *gen_state) error {
	for _, db := range state.databaseArr {
		if _, err := os.Stat(db.GenBinDir); err != nil && os.IsNotExist(err) {
			if os.IsNotExist(err) {
				if err := os.Mkdir(db.GenBinDir, 0755); err != nil {
					return err
				}
			} else {
				return err
			}
		}
		file, err := os.OpenFile(db.GenBinFilePath, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		file.Truncate(0)
		defer file.Close()
		tmpl := template.New("cli").Delims("<<", ">>")
		if tmpl, err := tmpl.Parse(cli_code); err != nil {
			return err
		} else {
			if err := tmpl.Execute(file, db); err != nil {
				return err
			}
		}
	}
	return nil
}

// 生成协议
func Gen() error {
	log.Info("===============gen behavior start===============")
	if _, err := os.Stat(proj.Config.SrcGenBehaviorDir); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(proj.Config.SrcGenBehaviorDir, 0755); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	fileNameArr := make([]string, 0)
	if err := filepath.WalkDir(proj.Config.DocBehaviorDir, func(path string, d os.DirEntry, err error) error {
		if d == nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(d.Name()) == ".yaml" {
			fileNameArr = append(fileNameArr, d.Name())
		}
		return nil
	}); err != nil {
		return err
	}
	sort.Strings(fileNameArr)
	state := &gen_state{
		databaseArr: make([]*Database, 0),
	}
	if err := parse(state, fileNameArr); err != nil {
		log.Info(err)
		return err
	}
	if err := genModel(state); err != nil {
		log.Info(err)
		return err
	}
	// 生成cli程序
	if err := genCli(state); err != nil {
		return err
	}
	log.Info("===============gen behavior finished===============")
	return nil
}
