package db

import (
	"context"
	"fmt"

	"github.com/lujingwei002/gira/log"

	"github.com/go-redis/redis/v8"
	"github.com/lujingwei002/gira"
)

type AdminCacheClient struct {
	cancelFunc context.CancelFunc
	cancelCtx  context.Context
	client     *redis.Client
	config     gira.AdminCacheConfig
}

func NewAdminCacheClient() *AdminCacheClient {
	return &AdminCacheClient{}
}

func (self *AdminCacheClient) GetRedisClient() *redis.Client {
	return self.client
}

// 初始化并启动
func (self *AdminCacheClient) OnAwake(ctx context.Context, config gira.AdminCacheConfig) error {
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
	log.Info("connect admin-cache success")
	return nil
}
