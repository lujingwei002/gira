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

type AccountDbClient struct {
	cancelFunc context.CancelFunc
	cancelCtx  context.Context
	client     *mongo.Client
	config     gira.AccountDbConfig
}

func NewAccountDbClient() *AccountDbClient {
	return &AccountDbClient{}
}

func (self *AccountDbClient) GetMongoDatabase() *mongo.Database {
	return self.client.Database(self.config.Db)
}

func (self *AccountDbClient) GetMongoClient() *mongo.Client {
	return self.client
}

// 初始化并启动
func (self *AccountDbClient) Start(ctx context.Context, config gira.AccountDbConfig) error {
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
		log.Fatal("connect accountdb error", err)
	}
	ctx2, cancelFunc2 := context.WithTimeout(self.cancelCtx, 3*time.Second)
	defer cancelFunc2()
	if err = client.Ping(ctx2, readpref.Primary()); err != nil {
		log.Fatal("connect accountdb error", err)
		return err
	}
	self.client = client
	log.Info("connect accountdb success")
	return nil
}
