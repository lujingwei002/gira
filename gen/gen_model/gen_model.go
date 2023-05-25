package gen_model

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
	"github.com/lujingwei002/gira/db"
	"github.com/urfave/cli/v2"
	"<<.Module>>/gen/model/<<.DbName>>"
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
	opts := make([]db.MigrateOption, 0)
	opts = append(opts, db.WithMigrateDropIndex(enabledDropIndex), db.WithMigrateConnectTimeout(connectTimeout))
	if client, err := db.NewDbClientFromUri(context.Background(), "<<.DbName>>", uri); err != nil {
		return err
	} else {
		return <<.DbName>>.Migrate(context.Background(), client, opts...)
	}
}
`

var protobuf_template = `
syntax = "proto3";  
package <<.DbName>>;
option go_package="src/gen/model/<<.DbName>>";

<<- range .CollectionArr>> 

message <<.StructName>>Pb {
	<<- range .FieldDict>> 
	/// <<.Comment>>
	<<.ProtobufTypeName>> <<.CamelName>> = <<.Tag>>;
	<<- end>>
}

<<- end>>
`

var model_template = `
// Code generated by github.com/lujingwei002/gira. DO NOT EDIT.
// Code generated by github.com/lujingwei002/gira. DO NOT EDIT.
// Code generated by github.com/lujingwei002/gira. DO NOT EDIT.

package <<.DbName>> 

// mongo模型
import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/go-redis/redis/v8"
	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/db"
	"github.com/lujingwei002/gira/log"
	"context"
	"encoding/json"
	"fmt"
	"time"
)



type <<.DriverInterfaceName>> interface {
	Migrate(ctx context.Context, opts ...db.MigrateOption) error
}
var globalDriver <<.DriverInterfaceName>>

func Use(ctx context.Context, client gira.DbClient) (<<.DriverInterfaceName>>, error) {
	switch c := client.(type) {
	case gira.MongoClient:
		return UseMongo(ctx, c)
	default:
		return nil, gira.ErrDbNotSupport
	}
}

func NewMongo() *<<.MongoDriverStructName>> {
	self := &<<.MongoDriverStructName>>{}
	<<- range .CollectionArr>> 
	self.<<.StructName>> = &<<.MongoDaoStructName>>{
		db: self,
	}
	<<- end>> 
	return self
}

func UseMongo(ctx context.Context, client gira.MongoClient) (*<<.MongoDriverStructName>>, error) {
	driver := NewMongo()
	if err := driver.Use(client); err != nil {
		return nil, err
	}
	return driver, nil
}

func NewRedis() *<<.RedisDriverStructName>> {
	self := &<<.RedisDriverStructName>>{}
	<<- range .CollectionArr>> 
	self.<<.StructName>> = &<<.RedisDaoStructName>>{
		db: self,
	}
	<<- end>> 
	return self
}

func UseRedis(ctx context.Context, client gira.RedisClient) (*<<.RedisDriverStructName>>, error) {
	driver := NewRedis()
	if err := driver.Use(client); err != nil {
		return nil, err
	}
	return driver, nil
}

func Migrate(ctx context.Context, client  gira.DbClient, opts ...db.MigrateOption) error {
	migrateOptions := &db.MigrateOptions {
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
		return gira.ErrDbNotSupport
	}
}

<<- range .CollectionArr>> 
// <<.CollName>>模型字段 
var <<.StructName>>Field *<<.StructName>>_Field = &<<.StructName>>_Field{
	<<- range .FieldArr>> 
	/// <<.Comment>>
	<<.CamelName>>: "<<.Name>>",
	<<- end>>
}
<<- end>> 

<<- range .CollectionArr>> 

<</* 模型字段 */>>
type <<.StructName>>_Field struct {
	<<- range .FieldArr>> 
	/// <<.Comment>>
	<<.CamelName>> string;
	<<- end>>
}

<</* 模型Data */>>
type <<.DataStructName>> struct {
	<<- range .FieldArr>> 
	/// <<.Comment>>
	<<.CamelName>> <<.GoTypeName>> <<quote>>bson:"<<.Name>>" json:"<<.Name>>"<<quote>>
	<<- end>>
}

func (self *<<.DataStructName>>) MarshalProtobuf(pb *<<.PbStructName>>) error {
	<<- range .FieldArr>> 
	<<- if eq .TypeName "id" >>
	pb.<<.CamelName>> = self.<<.CamelName>>.String()
	<<- else if .IsStruct >>
	if v, err := json.Marshal(self.<<.CamelName>>); err != nil {
		return err
	} else {
		pb.<<.CamelName>> = v
	}
	<<- else>>
	pb.<<.CamelName>> = self.<<.CamelName>> 
	<<- end>>
	<<- end>>
	return nil
}

func (self* <<.DataStructName>>)MarshalBinary() (data []byte, err error) {
	data, err = json.Marshal(self)
	return
}

<<- range .FieldArr>> 
func (self *<<.Coll.DataStructName>>) Get<<.CamelName>>() <<.GoTypeName>> {
	return self.<<.CamelName>> 
}
<<- end>> 

<<- end>> 




<</* ===========================derive user ================ */>>
<<- range .CollectionArr>> 
<<- if .IsDeriveUser>>
type <<.StructName>> struct {
	<<.DataStructName>>
	none  	bool // 是否为空数据
	dirty 	bool
}

func new<<.StructName>>() *<<.StructName>> {
	return &<<.StructName>> {
	}
}

/// 是否脏，如果为true,则需要保存到数据
func (self *<<.StructName>>) IsDirty() bool {
	return self.dirty
}

func (self *<<.StructName>>) SetDirty() {
	if self.dirty {
		return
	}
	self.dirty = true
}

/// 是否为空，即数据库还没有数据
func (self *<<.StructName>>) IsNone() bool {
	return self.none
}


<<- range .FieldArr>> 

func (self *<<.Coll.StructName>>) Set<<.CamelName>>(v <<.GoTypeName>>) {
	<<- if .IsComparable >>
	if self.<<.CamelName>> == v {
		return
	}
	<<- end>>
	self.<<.CamelName>> = v
	self.dirty = true
}

func (self *<<.Coll.StructName>>) Get<<.CamelName>>() <<.GoTypeName>> {
	return self.<<.CamelName>> 
}
<<- end>>

<<- end>><</* if .IsDeriveUser*/>>
<<- end>><</* range .CollectionArr*/>>



<</* ===========================derive userarr ================ */>>
<<- range .CollectionArr>> 
<<- if .IsDeriveUserArr>>
type <<.StructName>> struct {
	<<.DataStructName>>
	dirty   bool
	arr     *<<.ArrStructName>>
}

func (self *<<.StructName>>) SetDirty() {
	if self.dirty {
		return
	}
	self.dirty = true
	if self.arr != nil {
		self.arr.setDirty(self)
	}
}

<<- range .FieldArr>> 
<<- if .IsPrimaryKey>>
<<- else>>
<<- if .IsSecondaryKey>>
<<- else>>
func (self *<<.Coll.StructName>>) Set<<.CamelName>>(v <<.GoTypeName>>) {
	self.<<.CamelName>> = v
	if self.dirty {
		return
	}
	self.dirty = true
	if self.arr != nil {
		self.arr.setDirty(self)
	}
}
func (self *<<.Coll.StructName>>) Get<<.CamelName>>() <<.GoTypeName>> {
	return self.<<.CamelName>>
}
<<- end>>
<<- end>>
<<- end>>


type <<.ArrStructName>> struct {
	<<.CamelPrimaryKey>>	<<.PrimaryKeyField.GoTypeName>>	
	dict    map[<<.SecondaryKeyField.GoTypeName>>]*<<.StructName>>
	del		map[primitive.ObjectID]*<<.StructName>>
	dirty	map[primitive.ObjectID]*<<.StructName>>
	add     map[primitive.ObjectID]*<<.StructName>>
}

func new<<.StructName>>(arr *<<.ArrStructName>>) *<<.StructName>> {
	return &<<.StructName>>{
		arr: arr,
	}
}
func (self* <<.StructName>>)MarshalBinary() (data []byte, err error) {
	data, err = json.Marshal(self)
	return
}


func new<<.ArrStructName>>(<<.CapCamelPrimaryKey>> <<.PrimaryKeyField.GoTypeName>>) *<<.ArrStructName>> {
	return &<<.ArrStructName>>{
		<<.CamelPrimaryKey>>: <<.CapCamelPrimaryKey>>,
		dict: make(map[<<.SecondaryKeyField.GoTypeName>>]*<<.StructName>>),
		del: make(map[primitive.ObjectID]*<<.StructName>>),
		add: make(map[primitive.ObjectID]*<<.StructName>>),
		dirty: make(map[primitive.ObjectID]*<<.StructName>>),
	}
}

func (self *<<.ArrStructName>>) setDirty(doc *<<.StructName>>) {
	if _, ok := self.add[doc.Id]; ok {
		return
	}
	if _, ok := self.dirty[doc.Id]; ok {
		return
	}
	if _, ok := self.del[doc.Id]; ok {
		return
	}
	self.dirty[doc.Id] = doc
}

func (self *<<.ArrStructName>>) Add(<<.CapCamelSecondaryKey>> <<.SecondaryKeyField.GoTypeName>>) (*<<.StructName>>, error) {
	doc := new<<.StructName>>(self)
	doc.Id = primitive.NewObjectID()
	doc.dirty = false
	doc.<<.CamelSecondaryKey>> = <<.CapCamelSecondaryKey>>
	doc.<<.CamelPrimaryKey>> = self.<<.CamelPrimaryKey>>
	// 已经有相同id的数据了
	if _, ok := self.dict[<<.CapCamelSecondaryKey>>]; ok {
		return nil, gira.ErrDataExist
	}
	self.dict[<<.CapCamelSecondaryKey>>] = doc
	self.add[doc.Id] = doc 
	return doc, nil
}

func (self *<<.ArrStructName>>) Delete(<<.CapCamelSecondaryKey>> <<.SecondaryKeyField.GoTypeName>>) error {
	if doc, ok := self.dict[<<.CapCamelSecondaryKey>>]; !ok {
		return gira.ErrDataNotFound
	} else {
		// 已经准备删除也，也当成成功返回
		if _, ok := self.del[doc.Id]; ok {
			return nil 
		}
		// 除了del， 从各个字典中删除
		delete(self.dict, <<.CapCamelSecondaryKey>>)
		delete(self.dirty, doc.Id)
		delete(self.add, doc.Id)
		self.del[doc.Id] = doc
		return nil
	}
}

func (self *<<.ArrStructName>>) Count() int {
	return len(self.dict)
}

func (self *<<.ArrStructName>>) Range(f func(<<.CapCamelSecondaryKey>> <<.SecondaryKeyField.GoTypeName>>, value *<<.StructName>>) bool)  {
	for k, v := range self.dict {
		if !f(k, v) {
			break
		}
	}
}

func (self *<<.ArrStructName>>) Clear() error {
	for _, v := range self.dict {
		self.del[v.Id] = v
	}
	for oi := range self.add {
		delete(self.add, oi)
	}
	for oi := range self.dirty {
		delete(self.dirty, oi)
	}
	for oi := range self.dict {
		delete(self.dict , oi)
	}
	return nil
}


func (self *<<.ArrStructName>>) Get(<<.CapCamelSecondaryKey>> <<.SecondaryKeyField.GoTypeName>>) (*<<.StructName>>, bool) {
	v, ok := self.dict[<<.CapCamelSecondaryKey>>]
	return v, ok
}



<<- end>><</* if .IsDeriveUserArr*/>>
<<- end>><</* range .CollectionArr*/>>



// mongo 
type <<.MongoDriverStructName>> struct {
	client		*mongo.Client
	database	*mongo.Database
	<<- range .CollectionArr>> 
	<<.StructName>>  *<<.MongoDaoStructName>>
	<<- end>>
}

func (self *<<.MongoDriverStructName>>) Use(client gira.MongoClient) error {
	if self.client != nil {
		return gira.ErrTodo
	}
	self.client = client.GetMongoClient()
	self.database = client.GetMongoDatabase()
	return nil
}

func (self *<<.MongoDriverStructName>>) Migrate(ctx context.Context, opts ...db.MigrateOption) error {
<<- range .CollectionArr>>
	if err := self.<<.StructName>>.Migrate(ctx, opts...); err != nil {
		return err
	}
<<- end>> 
	return nil
}

<<- range .CollectionArr>> 

type <<.MongoDaoStructName>> struct {
	db *<<$.MongoDriverStructName>>
}

func (self *<<.MongoDaoStructName>>) New() *<<.DataStructName>> {
	doc := &<<.DataStructName>>{}
	doc.Id = primitive.NewObjectID()
	return doc
}

func (self *<<.MongoDaoStructName>>) Migrate(ctx context.Context, opts ...db.MigrateOption) error {
	migrateOptions := &db.MigrateOptions {
	}
	for _, v := range opts {
		v.ConfigMigrateOptions(migrateOptions)
	}
	database := self.db.database
	<<- if .IsCapped >>
	createOpts := options.CreateCollection()
	createOpts.SetCapped(true)
	createOpts.SetMaxDocuments(<<.Capped>>)
	createOpts.SetSizeInBytes(<<.Capped>>)
	if err := database.CreateCollection(ctx, "<<.CollName>>", createOpts); err != nil {
		log.Warn(err)
	}
	<<- end>>
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

func (self *<<.MongoDaoStructName>>) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (*<<.DataStructName>>, error) {
	database := self.db.database
	coll := database.Collection("<<.CollName>>")
	doc := &<<.DataStructName>>{}
	err := coll.FindOne(ctx, filter, opts...).Decode(doc)
	if err != nil {
		return nil, err
	}
	return doc, err
}


func (self *<<.MongoDaoStructName>>) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]*<<.DataStructName>>, error) {
	database := self.db.database
	coll := database.Collection("<<.CollName>>")
	cursor, err := coll.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	var results []*<<.DataStructName>> = make([]*<<.DataStructName>>, 0)
	for {
		if !cursor.Next(ctx) {
			break
		}
		v := &<<.DataStructName>>{}
		if err := cursor.Decode(v); err != nil {
			return results, err
		}
		results = append(results, v)
	}
	return results, nil
}

func (self *<<.MongoDaoStructName>>) UpdateOne(ctx context.Context, filter interface{}, update bson.D, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	database := self.db.database
	coll := database.Collection("<<.CollName>>")
	result, err := coll.UpdateOne(ctx, filter, update, opts...)
	return result, err
}

func (self *<<.MongoDaoStructName>>) ReplaceOne(ctx context.Context, filter interface{}, replacement *<<.DataStructName>>, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	database := self.db.database
	coll := database.Collection("<<.CollName>>")
	result, err := coll.ReplaceOne(ctx, filter, &replacement, opts...)
	return result, err
}

func (self *<<.MongoDaoStructName>>) UpsertOne(ctx context.Context, filter interface{}, replacement *<<.DataStructName>>, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	database := self.db.database
	coll := database.Collection("<<.CollName>>")
	opts = append(opts, options.Replace().SetUpsert(true))
	result, err := coll.ReplaceOne(ctx, filter, &replacement, opts...)
	return result, err
}

func (self *<<.MongoDaoStructName>>) InsertOne(ctx context.Context, doc *<<.DataStructName>>, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	database := self.db.database
	coll := database.Collection("<<.CollName>>")
	result, err := coll.InsertOne(ctx, doc, opts...)
	return result, err
}
<<- if .IsDeriveUser>>

func (self *<<.MongoDaoStructName>>) Save(ctx context.Context, doc *<<.StructName>>) error {
	log.Infow("<<.CollName>> save", "dirty", doc.dirty, "id", doc.Id)
	if !doc.dirty {
		return nil
	}
	database := self.db.database
	coll := database.Collection("<<.CollName>>")
	opts := options.Replace().SetUpsert(true)
	if result, err := coll.ReplaceOne(ctx, bson.D{{"_id", doc.Id}}, (*doc).<<.DataStructName>>, opts); err != nil {
		return err
	} else {
		if doc.none && result.MatchedCount != 0 {
			log.Infof("<<.CollName >> matched count is not zero, %+v", doc)
		} else if !doc.none && result.MatchedCount != 1 {
			log.Infof("<<.CollName >> matched count is zero, %+v", doc)
		}
		if result.ModifiedCount + result.UpsertedCount != 1 {
			log.Infof("<<.CollName >> unchanged, %+v", doc)
		}
	}
	doc.dirty = false
	return nil
}

func (self *<<.MongoDaoStructName>>) Load(ctx context.Context, id primitive.ObjectID) (*<<.StructName>>, error) {
    doc := new<<.StructName>>()
	database := self.db.database
	log.Infow("<<.CollName>> load", "id", id)
	coll := database.Collection("<<.CollName>>")
	err := coll.FindOne(ctx, bson.D{{"_id", id}}).Decode(&doc.<<.DataStructName>>)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			doc.Id = id
			doc.none = true
		} else {
			return nil, err
		}
	} else {
		doc.Id = id
	}
	return doc, nil
}

func (self *<<.MongoDaoStructName>>) Delete(ctx context.Context, doc *<<.StructName>>) error {
	if doc == nil {
		return gira.ErrNullPonter
	}
	if doc.none {
		return gira.ErrDataNotExist
	}
	if doc.Id.IsZero() {
		return gira.ErrDataNotExist
	}
	database := self.db.database
	coll := database.Collection("<<.CollName>>")
	if result, err := coll.DeleteOne(ctx, bson.D{{"_id", doc.Id}}); err != nil {
		return err
	} else {
		if result.DeletedCount != 1 {
			return gira.ErrDataDeleteFail
		}
	}
	doc.dirty = false
	doc.none = true
	return nil
}


<<- end>>


<<- if .IsDeriveUserArr>>


func (self *<<.MongoDaoStructName>>) Save(ctx context.Context, doc *<<.ArrStructName>>) error {
	database := self.db.database
	coll := database.Collection("<<.CollName>>")
	// opts := options.Replace().SetUpsert(true)
	log.Infow("<<.CollName>> save", "<<.CapCamelPrimaryKey>>", doc.<<.CamelPrimaryKey>>)
	models := make([]mongo.WriteModel, 0)
	if len(doc.del) > 0 {
		for _, v := range(doc.del) {
			log.Debugw("<<.CollName>> del", "id", v.Id, "<<.CapCamelSecondaryKey>>", v.<<.CamelSecondaryKey>>)
			models = append(models, mongo.NewDeleteOneModel().SetFilter(bson.D{{"_id", v.Id}}))
		}
	}
	if len(doc.add) > 0 {
		for _, v := range(doc.add) {
			log.Debugw("<<.CollName>> add", "id", v.Id, "<<.CapCamelSecondaryKey>>", v.<<.CamelSecondaryKey>>)
			models = append(models, mongo.NewInsertOneModel().SetDocument(&v.<<.DataStructName>>))
		}
	}
	if len(doc.dirty) > 0 {
		for _, v := range(doc.dirty) {
			log.Debugw("<<.CollName>> update", "id", v.Id, "<<.CapCamelSecondaryKey>>", v.<<.CamelSecondaryKey>>)
			models = append(models, mongo.NewReplaceOneModel().
				SetFilter(bson.D{{"_id", v.Id}}).
				SetReplacement(&v.<<.DataStructName>>))
		}
	}

	for k, v := range doc.add {
		v.dirty = false
		delete(doc.add, k)
	}
	for k, v := range doc.dirty {
		v.dirty = false
		delete(doc.dirty, k)
	}
	for oi := range doc.del {
		delete(doc.dirty, oi)
	}
	if len(models) <= 0 {
		return nil
	}
	opts := options.BulkWrite().SetOrdered(false)
	_, err := coll.BulkWrite(ctx, models, opts)
	if err != nil {
		log.Error(err)
	    return err
	}
	return nil
}

func (self *<<.MongoDaoStructName>>) Load(ctx context.Context, <<.CapCamelPrimaryKey>> <<.PrimaryKeyField.GoTypeName>>) (*<<.ArrStructName>>, error) {
	database := self.db.database
	coll := database.Collection("<<.CollName>>")
	log.Infow("<<.CollName>> load", "<<.PrimaryKey>>", <<.CapCamelPrimaryKey>>)
	cursor, err := coll.Find(ctx, bson.D{{"<<.PrimaryKey>>", <<.CapCamelPrimaryKey>>}})
	if err != nil {
		return nil, err
	}
	var results []*<<.StructName>> = make([]*<<.StructName>>, 0)
	for {
		if !cursor.Next(ctx) {
			break
		}
		v := &<<.StructName>>{}
		if err := cursor.Decode(&v.<<.DataStructName>>); err != nil {
			return nil, err
		}
		results = append(results, v)
	}
    arr := new<<.ArrStructName>>(<<.CapCamelPrimaryKey>>)
	for _, v := range results {
		v.arr = arr
		arr.dict[v.<<.CamelSecondaryKey>>] = v
	}
	return arr, nil
}
<<- end>>
<<- end>>








<</* redis操作 */>>
type <<.RedisDriverStructName>> struct {
	client		*redis.Client
	<<- range .CollectionArr>> 
	<<.StructName>>  *<<.RedisDaoStructName>>
	<<- end>>
}



func (self *<<.RedisDriverStructName>>) Migrate(ctx context.Context, opts ...db.MigrateOption) error {
	return nil
}

func (self *<<.RedisDriverStructName>>) Use(client gira.RedisClient) error {
	self.client = client.GetRedisClient()
	return nil
}
<<- range .CollectionArr>> 

type <<.RedisDaoStructName>> struct {
	db *<<$.RedisDriverStructName>>
}

func (self *<<.RedisDaoStructName>>) New() *<<.DataStructName>> {
	doc := &<<.DataStructName>>{}
	doc.Id = primitive.NewObjectID()
	return doc
}


func (self *<<.RedisDaoStructName>>) Set(ctx context.Context, key primitive.ObjectID, value *<<.DataStructName>>, expiration time.Duration) *redis.StatusCmd {
	rkey := fmt.Sprintf("%s@<<.CollName>>", key.Hex())
	result := self.db.client.Set(ctx, rkey, value, expiration)
	return result
}

func (self *<<.RedisDaoStructName>>) Get(ctx context.Context, key primitive.ObjectID) (*<<.DataStructName>>, error) {
	rkey := fmt.Sprintf("%s@<<.CollName>>", key.Hex())
	result := self.db.client.Get(ctx, rkey)
	if result.Err() != nil {
		return nil, result.Err()
	}
	if b, err := result.Bytes(); err != nil {
		return nil, err
	} else {
		doc := &<<.DataStructName>>{}
		if err := json.Unmarshal(b, doc); err != nil {
			return nil, err
		}
		return doc, nil
	}
}




<<- if .IsDeriveUserArr>>
func (self *<<.RedisDaoStructName>>) HSet(ctx context.Context, key primitive.ObjectID, values ...interface{}) *redis.IntCmd {
	rkey := fmt.Sprintf("%s@<<.CollName>>", key.Hex())
	result := self.db.client.HSet(ctx, rkey, values...)
	return result
}
<<- end>>



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
	field_type_int:       "int64",
	field_type_int32:     "int32",
	field_type_int64:     "int64",
	field_type_string:    "string",
	field_type_objectid:  "primitive.ObjectID",
	field_type_bool:      "bool",
	field_type_bytes:     "[]byte",
	field_type_int_arr:   "[]int64",
	field_type_int64_arr: "[]int64",
}

var protobuf_type_name_dict = map[field_type]string{
	field_type_int:       "int64",
	field_type_int32:     "int32",
	field_type_int64:     "int64",
	field_type_string:    "string",
	field_type_objectid:  "string",
	field_type_bool:      "bool",
	field_type_bytes:     "bytes",
	field_type_int64_arr: "repeated int64",
	field_type_int_arr:   "repeated int64",
}

type message_type int

const (
	message_type_struct message_type = iota
	message_type_request
	message_type_response
	message_type_notify
	message_type_push
)

type Field struct {
	Tag              int
	Name             string
	CamelName        string
	Type             field_type
	Array            bool
	TypeName         string
	GoTypeName       string
	ProtobufTypeName string
	Default          interface{}
	Comment          string
	Coll             *Collection
	IsPrimaryKey     bool /// 是否主键，目前只对userarr类型的表格有效果
	IsSecondaryKey   bool /// 是否次键，目前只对userarr类型的表格有效果
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

func (f *Field) IsStruct() bool {
	return f.Type == field_type_struct
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
	PbStructName          string
	ArrStructName         string
	MongoDaoStructName    string // mongo dao 结构的名称
	RedisDaoStructName    string // redis dao 结构的名称
	Derive                string
	KeyArr                []string
	DataStructName        string
	FieldDict             map[string]*Field
	FieldArr              []*Field
	SecondaryKey          string
	CamelSecondaryKey     string
	CapCamelSecondaryKey  string
	PrimaryKey            string
	CamelPrimaryKey       string
	CapCamelPrimaryKey    string
	PrimaryKeyField       *Field
	SecondaryKeyField     *Field
	IndexDict             map[string]*Index
	IndexArr              []*Index
	MongoDriverStructName string
	Capped                int64
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
	RedisDriverStructName string // redis 的 dao 结构名字
	DbName                string
	GenBinFilePath        string        // 生成的文件路径，在 gen/model/{{DbName}}/bin/{{DbName}}.gen.go
	GenBinDir             string        // 生成的文件路径，在 gen/model/{{DbName}}/bin
	GenModelDir           string        // 生成的文件路径，在 gen/model/{{DbName}}
	GenModelFilePath      string        // 生成的文件路径，在 gen/model/{{DbName}}/{{DbName}}.gen.go
	GenProtobufFilePath   string        // 生成的protobuf文件路径， 在gen/{{DbName}}/{{DbName}}.gen.proto
	CollectionArr         []*Collection // 所有的模型
	DriverInterfaceName   string
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

func (coll *Collection) IsCapped() bool {
	return coll.Capped != 0
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

func (descriptor *Collection) parseStruct(attrs map[string]interface{}) error {
	descriptor.FieldDict = make(map[string]*Field)
	descriptor.FieldArr = make([]*Field, 0)
	//commaRegexp := regexp.MustCompile("[^,]+")
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
			Coll:      descriptor,
			Name:      fieldName,
			CamelName: camelString(fieldName),
			Array:     false,
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
			field.ProtobufTypeName = protobuf_type_name_dict[field.Type]
		} else {
			field.Type = field_type_struct
			field.GoTypeName = typeStr
			field.ProtobufTypeName = "bytes"
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

		descriptor.FieldDict[fieldName] = field
		descriptor.FieldArr = append(descriptor.FieldArr, field)
	}
	sort.Sort(SortFieldByName(descriptor.FieldArr))
	return nil
}

func (coll *Collection) Unmarshal(genState *gen_state, v interface{}) error {
	var derive string
	keyArr := make([]string, 0)
	row := v.(map[string]interface{})
	if v, ok := row["derive"]; ok {
		derive = v.(string)
	}
	// } else {
	// 	return fmt.Errorf("collection %s derive part not found", coll.CollName)
	// }
	if _, ok := row["struct"]; !ok {
		return fmt.Errorf("collection %s struct part not found", coll.CollName)
	}
	structPart := row["struct"]
	if _, ok := structPart.(map[string]interface{}); !ok {
		return fmt.Errorf("collection %s struct part not map", coll.CollName)
	}
	if _, ok := row["key"]; ok {
		keyPart := row["key"]
		if _, ok := keyPart.([]interface{}); !ok {
			return fmt.Errorf("collection %s key part not array", coll.CollName)
		}
		for _, v := range keyPart.([]interface{}) {
			keyArr = append(keyArr, v.(string))
		}
	}
	// } else {
	// 	return fmt.Errorf("collection %s key part not found", coll.CollName)

	// }
	coll.Derive = derive
	coll.KeyArr = keyArr
	if err := coll.parseStruct(row["struct"].(map[string]interface{})); err != nil {
		return err
	}
	// 解析index
	if v, ok := row["index"]; !ok {
	} else if v2, ok := v.(map[string]interface{}); !ok {
	} else if err := coll.parseIndex(v2); err != nil {
		return err
	}
	if coll.Derive == "userarr" {
		if len(keyArr) != 2 {
			return fmt.Errorf("collection %s derive userarr need 2 key", coll.CollName)
		}
		primaryKey := keyArr[0]
		secondaryKey := keyArr[1]
		coll.PrimaryKey = primaryKey
		coll.CamelPrimaryKey = camelString(primaryKey)
		coll.CapCamelPrimaryKey = capLowerString(coll.CamelPrimaryKey)
		coll.SecondaryKey = secondaryKey
		coll.CamelSecondaryKey = camelString(secondaryKey)
		coll.CapCamelSecondaryKey = capLowerString(coll.CamelSecondaryKey)
		if field, ok := coll.FieldDict[primaryKey]; ok {
			coll.PrimaryKeyField = field
			field.IsPrimaryKey = true
		} else {
			return fmt.Errorf("collection %s derive userarr, but primary key %s not found", coll.CollName, primaryKey)
		}
		if field, ok := coll.FieldDict[secondaryKey]; ok {
			coll.SecondaryKeyField = field
			field.IsSecondaryKey = true
		} else {
			return fmt.Errorf("collection %s derive userarr, but secondary key %s not found", coll.CollName, secondaryKey)
		}
	}
	if v, ok := row["capped"]; ok {
		coll.Capped = int64(v.(int))
	}
	return nil
}

func parse(state *gen_state, filePathArr []string) error {
	for _, fileName := range filePathArr {
		filePath := path.Join(proj.Config.DocModelDir, fileName)
		log.Info("处理文件", filePath)
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}
		dbName := strings.Replace(fileName, ".yaml", "", 1)
		database := &Database{
			Module:                proj.Config.Module,
			CollectionArr:         make([]*Collection, 0),
			GenBinFilePath:        path.Join(proj.Config.SrcGenModelDir, dbName, "bin", fmt.Sprintf("%s.gen.go", dbName)),
			GenBinDir:             path.Join(proj.Config.SrcGenModelDir, dbName, "bin"),
			GenModelDir:           path.Join(proj.Config.SrcGenModelDir, dbName),
			GenModelFilePath:      path.Join(proj.Config.SrcGenModelDir, dbName, fmt.Sprintf("%s.gen.go", dbName)),
			GenProtobufFilePath:   path.Join(proj.Config.GenModelDir, dbName, fmt.Sprintf("%s.gen.proto", dbName)),
			DbName:                dbName,
			DbStructName:          camelString(dbName),
			MongoDriverStructName: fmt.Sprintf("%sMongoDriver", camelString(dbName)),
			RedisDriverStructName: fmt.Sprintf("%sRedisDriver", camelString(dbName)),
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
					PbStructName:          fmt.Sprintf("%sPb", camelString(collName)),
					ArrStructName:         fmt.Sprintf("%sArr", camelString(collName)),
					DataStructName:        fmt.Sprintf("%sData", camelString(collName)),
					MongoDaoStructName:    fmt.Sprintf("%sMongoDao", camelString(collName)),
					RedisDaoStructName:    fmt.Sprintf("%sRedisDao", camelString(collName)),
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

func genModel(protocolState *gen_state) error {
	for _, db := range protocolState.databaseArr {
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

func genProtobuf(protocolState *gen_state) error {
	for _, db := range protocolState.databaseArr {

		dir := path.Dir(db.GenProtobufFilePath)
		if _, err := os.Stat(dir); err != nil {
			if os.IsNotExist(err) {
				if err := os.Mkdir(dir, 0755); err != nil {
					return err
				}
			} else {
				return err
			}
		}

		var err error
		file, err := os.OpenFile(db.GenProtobufFilePath, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		file.Truncate(0)
		funcMap := template.FuncMap{
			"quote": QuoteChar,
		}
		tmpl := template.New("model").Delims("<<", ">>")
		tmpl.Funcs(funcMap)
		tmpl, err = tmpl.Parse(protobuf_template)
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
	log.Info("===============gen model start===============")
	if _, err := os.Stat(proj.Config.SrcGenModelDir); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(proj.Config.SrcGenModelDir, 0755); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	fileNameArr := make([]string, 0)
	filepath.WalkDir(proj.Config.DocModelDir, func(path string, d os.DirEntry, err error) error {
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
	})
	sort.Strings(fileNameArr)
	state := &gen_state{
		databaseArr: make([]*Database, 0),
	}
	if err := parse(state, fileNameArr); err != nil {
		log.Info(err)
		return err
	}
	if err := genProtobuf(state); err != nil {
		log.Info(err)
		return err
	}
	if err := genModel(state); err != nil {
		log.Info(err)
		return err
	}
	if err := genCli(state); err != nil {
		log.Info(err)
		return err
	}
	log.Info("===============gen model finished===============")
	return nil
}
