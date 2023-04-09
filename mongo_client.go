package gira

import (
	"go.mongodb.org/mongo-driver/mongo"
)

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
