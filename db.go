package gira

import (
	"database/sql"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

type DbClient interface {
}

const (
	GAMEDB_NAME       = "gamedb"
	RESOURCEDB_NAME   = "resourcedb"
	STATDB_NAME       = "statdb"
	ACCOUNTDB_NAME    = "accountdb"
	LOGDB_NAME        = "logdb"
	ACCOUNTCACHE_NAME = "accountcache"
	ADMINCACHE_NAME   = "admincache"
	ADMINDB_NAME      = "admindb"
)

type RedisClient interface {
	GetRedisClient() *redis.Client
}

type MysqlClient interface {
	GetMysqlClient() *sql.DB
}

type MongoClient interface {
	GetMongoClient() *mongo.Client
	GetMongoDatabase() *mongo.Database
}

type GameDbClient interface {
	GetGameDbClient() MongoClient
}

type StatDbClient interface {
	GetStatDbClient() MongoClient
}

type AccountDbClient interface {
	GetAccountDbClient() MongoClient
}

type ResourceDbClient interface {
	GetResourceDbClient() MongoClient
}

type AdminDbClient interface {
	GetAdminDbClient() MysqlClient
}

type AccountCacheClient interface {
	GetAccountCacheClient() RedisClient
}

type LogDbClient interface {
	GetLogDbClient() MongoClient
}
type AdminCacheClient interface {
	GetAdminCacheClient() RedisClient
}
