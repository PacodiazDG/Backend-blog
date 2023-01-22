package database

import (
	"context"
	"os"
	"time"

	"github.com/PacodiazDG/Backend-blog/modules/logs"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Database *mongo.Database
var RedisCon *redis.Client
var Ctx = context.Background()

// Initializes the connection to MongoDB and transfers it to a global variable
func Initdb() {
	initvardb, err := newConnection()
	if err != nil {
		panic(err)
	}
	Database = initvardb.Database("Blog")
}

// Creates the connection to mongoDB and returns it in a pointer (*mongo.Client)
func newConnection() (*mongo.Client, error) {
	dbConfig := os.Getenv("DB_CONFIG")
	client, err := mongo.NewClient(options.Client().ApplyURI(dbConfig))
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		logs.WriteLogs(err, logs.CriticalError)
		panic(err)
	}
	client.Database("Blog")
	return client, nil
}

// Initializes the connection to redisDB and transfers it to a global variable
func InitRedis() {
	RedisCon = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("RedisAddr"),
		Password: os.Getenv("RedisPassword"),
		DB:       0,
		PoolSize: 30,
	})
	err := RedisCon.Set("key", "TestOk", 0).Err()
	if err != nil {
		logs.WriteLogs(err, logs.CriticalError)
		panic(err)
	}
}
