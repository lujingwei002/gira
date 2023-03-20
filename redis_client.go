package gira

import "github.com/go-redis/redis/v8"

type RedisClient interface {
	GetRedisClient() *redis.Client
}
