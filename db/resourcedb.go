package db

import (
	"context"
	"fmt"
	"time"

	"github.com/lujingwei002/gira/log"

	"github.com/lujingwei002/gira"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type ResourceDbClient struct {
	cancelFunc context.CancelFunc
	cancelCtx  context.Context
	client     *mongo.Client
	config     gira.ResourceDbConfig
}

func NewResourceDbClient() *ResourceDbClient {
	return &ResourceDbClient{}
}

func (self *ResourceDbClient) GetMongoDatabase() *mongo.Database {
	return self.client.Database(self.config.Db)
}

func (self *ResourceDbClient) GetMongoClient() *mongo.Client {
	return self.client
}

// 初始化并启动
func (self *ResourceDbClient) OnAwake(ctx context.Context, config gira.ResourceDbConfig) error {
	self.config = config
	self.cancelCtx, self.cancelFunc = context.WithCancel(ctx)
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d", config.User, config.Password, config.Host, config.Port)
	log.Info(uri)
	clientOpts := options.Client().
		ApplyURI(uri)
	ctx, cancelFunc := context.WithTimeout(self.cancelCtx, 3*time.Second)
	defer cancelFunc()
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal("connect resourcedb error", err)
	}
	ctx2, cancelFunc2 := context.WithTimeout(self.cancelCtx, 3*time.Second)
	defer cancelFunc2()
	if err = client.Ping(ctx2, readpref.Primary()); err != nil {
		log.Fatal("connect resourcedb error", err)
		return err
	}
	self.client = client
	log.Info("connect resourcedb success")
	return nil
}
