package gira

import "go.mongodb.org/mongo-driver/mongo"

type MongoClient interface {
	GetMongoClient() *mongo.Client
	GetMongoDatabase() *mongo.Database
}
