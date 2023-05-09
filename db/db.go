package db

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDbClient struct {
	cancelFunc context.CancelFunc
	ctx        context.Context
	client     *mongo.Client
	config     gira.DbConfig
}

type RedisClient struct {
	cancelFunc context.CancelFunc
	ctx        context.Context
	client     *redis.Client
	config     gira.DbConfig
}

type MysqlClient struct {
	cancelFunc context.CancelFunc
	ctx        context.Context
	client     *sql.DB
	config     gira.DbConfig
}

func (self *MysqlClient) GetMysqlClient() *sql.DB {
	return self.client
}

func (self *RedisClient) GetRedisClient() *redis.Client {
	return self.client
}

func (self *MongoDbClient) GetMongoDatabase() *mongo.Database {
	return self.client.Database(self.config.Db)
}

func (self *MongoDbClient) GetMongoClient() *mongo.Client {
	return self.client
}

func ConfigDbClient(ctx context.Context, name string, config gira.DbConfig) (gira.DbClient, error) {
	switch config.Driver {
	case "mongo":
		return configMongoDbClient(ctx, name, config)
	case "redis":
		return configRedisClient(ctx, name, config)
	case "mysql":
		return configMysqlClient(ctx, name, config)
	default:
		return nil, gira.ErrTodo
	}
}

func configMongoDbClient(ctx context.Context, name string, config gira.DbConfig) (gira.DbClient, error) {
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	client := &MongoDbClient{
		config:     config,
		cancelFunc: cancelFunc,
		ctx:        cancelCtx,
	}
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d", config.User, config.Password, config.Host, config.Port)
	log.Info(uri)
	clientOpts := options.Client().ApplyURI(uri)
	ctx1, cancelFunc1 := context.WithTimeout(client.ctx, 3*time.Second)
	defer cancelFunc1()
	conn, err := mongo.Connect(ctx1, clientOpts)
	if err != nil {
		log.Errorw("connect database fail", "name", name, "error", err)
		return nil, err
	}
	ctx2, cancelFunc2 := context.WithTimeout(client.ctx, 3*time.Second)
	defer cancelFunc2()
	if err = conn.Ping(ctx2, readpref.Primary()); err != nil {
		log.Errorw("connect database fail", "name", name, "error", err)
		return nil, err
	}
	client.client = conn
	log.Infow("connect database success", "name", name)
	return client, nil
}

func configRedisClient(ctx context.Context, name string, config gira.DbConfig) (gira.DbClient, error) {
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	client := &RedisClient{
		config:     config,
		cancelFunc: cancelFunc,
		ctx:        cancelCtx,
	}
	var err error
	var db int
	db, err = strconv.Atoi(config.Db)
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       db,
	})
	ctx1, cancelFunc1 := context.WithTimeout(client.ctx, 3*time.Second)
	defer cancelFunc1()
	if _, err := rdb.Ping(ctx1).Result(); err != nil {
		log.Errorw("connect database fail", "name", name, "error", err)
		return nil, err
	}
	client.client = rdb
	log.Infow("connect database success", "name", name)
	return client, nil
}

func configMysqlClient(ctx context.Context, name string, config gira.DbConfig) (gira.DbClient, error) {
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	client := &MysqlClient{
		config:     config,
		cancelFunc: cancelFunc,
		ctx:        cancelCtx,
	}
	uri := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", config.User, config.Password, config.Host, config.Port, config.Db)
	db, err := sql.Open("mysql", uri)
	if err != nil {
		log.Errorw("connect database fail", "name", name, "error", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Errorw("connect database fail", "name", name, "error", err)
		return nil, err
	}
	client.client = db
	log.Infow("connect database success", "name", name)
	return client, nil
}
