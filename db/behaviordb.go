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

type BehaviorDbClient struct {
	cancelFunc context.CancelFunc
	ctx        context.Context
	client     *mongo.Client
	config     gira.BehaviorDbConfig
}

func NewBehaviorDbClient() *BehaviorDbClient {
	return &BehaviorDbClient{}
}

func ConfigBehaviorDbClient(ctx context.Context, config gira.BehaviorDbConfig) (client *BehaviorDbClient, err error) {
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	client = &BehaviorDbClient{
		config:     config,
		cancelFunc: cancelFunc,
		ctx:        cancelCtx,
	}
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d", config.User, config.Password, config.Host, config.Port)
	log.Info(uri)
	clientOpts := options.Client().ApplyURI(uri)
	ctx, cancelFunc1 := context.WithTimeout(client.ctx, 3*time.Second)
	defer cancelFunc1()
	conn, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal("connect behaviordb error", err)
		return
	}
	ctx2, cancelFunc2 := context.WithTimeout(client.ctx, 3*time.Second)
	defer cancelFunc2()
	if err = conn.Ping(ctx2, readpref.Primary()); err != nil {
		log.Fatal("connect behaviordb error", err)
		return
	}
	client.client = conn
	log.Info("connect behaviordb success")
	return
}

func (self *BehaviorDbClient) GetMongoDatabase() *mongo.Database {
	return self.client.Database(self.config.Db)
}

func (self *BehaviorDbClient) GetMongoClient() *mongo.Client {
	return self.client
}

// 初始化并启动
func (self *BehaviorDbClient) OnAwake(ctx context.Context, config gira.BehaviorDbConfig) error {
	return nil
}
