package redisconfig

import (
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	once        sync.Once
)

func ConnectRedis() *redis.Client {
	once.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:         os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
			Password:     os.Getenv("REDIS_PASSWORD"),
			DB:           0,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		})
	})
	return redisClient
}
