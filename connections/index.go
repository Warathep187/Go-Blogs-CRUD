package connections

import (
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

var MongoClient *mongo.Client
var RedisClient *redis.Client

func InitDatabaseConnection() {
	MongoClient = newMongoClient()
	RedisClient = newRedisClient()
}
