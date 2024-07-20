package connections

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func newRedisClient() *redis.Client {
	rds := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	if err := rds.Ping(context.TODO()).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Redis connected")
	return rds
}
