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

type GameDbClient struct {
	cancelFunc context.CancelFunc
	ctx        context.Context
	client     *mongo.Client
	config     gira.GameDbConfig
}

func NewGameDbClient() *GameDbClient {
	return &GameDbClient{}
}

func ConfigGameDbClient(ctx context.Context, config gira.GameDbConfig) (client *GameDbClient, err error) {
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	client = &GameDbClient{
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
		log.Fatal("connect gamedb error", err)
		return
	}
	ctx2, cancelFunc2 := context.WithTimeout(client.ctx, 3*time.Second)
	defer cancelFunc2()
	if err = conn.Ping(ctx2, readpref.Primary()); err != nil {
		log.Fatal("connect gamedb error", err)
		return
	}
	client.client = conn
	log.Info("connect gamedb success")
	return
}

func (self *GameDbClient) GetMongoDatabase() *mongo.Database {
	return self.client.Database(self.config.Db)
}

func (self *GameDbClient) GetMongoClient() *mongo.Client {
	return self.client
}

// 初始化并启动
func (self *GameDbClient) OnAwake(ctx context.Context, config gira.GameDbConfig) error {
	return nil
}
