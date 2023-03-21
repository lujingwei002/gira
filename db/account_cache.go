package db

import (
	"context"
	"fmt"

	"github.com/lujingwei/gira/log"

	"github.com/go-redis/redis/v8"
	"github.com/lujingwei/gira"
)

type AccountCacheClient struct {
	cancelFunc context.CancelFunc
	cancelCtx  context.Context
	client     *redis.Client
	config     gira.AccountCacheConfig
}

func NewAccountCacheClient() *AccountCacheClient {
	return &AccountCacheClient{}
}

func (self *AccountCacheClient) GetRedisClient() *redis.Client {
	return self.client
}

// 初始化并启动
func (self *AccountCacheClient) Start(ctx context.Context, config gira.AccountCacheConfig) error {
	self.config = config
	self.cancelCtx, self.cancelFunc = context.WithCancel(ctx)

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password, // no password set
		DB:       config.Db,       // use default DB
	})
	if _, err := rdb.Ping(self.cancelCtx).Result(); err != nil {
		return err
	}
	self.client = rdb
	log.Info("connect account-cache success")
	return nil
}
