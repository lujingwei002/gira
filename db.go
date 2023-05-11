package gira

import (
	"database/sql"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	GAMEDB_NAME       = "gamedb"
	RESOURCEDB_NAME   = "resourcedb"
	STATDB_NAME       = "statdb"
	ACCOUNTDB_NAME    = "accountdb"
	LOGDB_NAME        = "logdb"
	BEHAVIORDB_NAME   = "behaviordb"
	ACCOUNTCACHE_NAME = "accountcache"
	ADMINCACHE_NAME   = "admincache"
	ADMINDB_NAME      = "admindb"
)
const (
	MONGODB_NAME = "mongodb"
	REDIS_NAME   = "redis"
	MYSQL_NAME   = "mysql"
)

type DbClient interface {
	Uri() string
}

type RedisClient interface {
	DbClient
	GetRedisClient() *redis.Client
}

type MysqlClient interface {
	DbClient
	GetMysqlClient() *sql.DB
}

type MongoClient interface {
	DbClient
	GetMongoClient() *mongo.Client
	GetMongoDatabase() *mongo.Database
}

type GameDbClient interface {
	GetGameDbClient() DbClient
}

type StatDbClient interface {
	GetStatDbClient() DbClient
}

type AccountDbClient interface {
	GetAccountDbClient() DbClient
}

type ResourceDbClient interface {
	GetResourceDbClient() DbClient
}

type AdminDbClient interface {
	GetAdminDbClient() DbClient
}

type AccountCacheClient interface {
	GetAccountCacheClient() DbClient
}

type LogDbClient interface {
	GetLogDbClient() DbClient
}

type BehaviorDbClient interface {
	GetBehaviorDbClient() DbClient
}

type AdminCacheClient interface {
	GetAdminCacheClient() DbClient
}
