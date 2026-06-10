package redis

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	once   sync.Once
)

func GetClient() *redis.Client {
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_URL"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		})
	})
	return client
}

func SetIfNotExists(ctx context.Context, key string, value string, ttlSeconds int) (bool, error) {
	cli := GetClient()
	res, err := cli.SetNX(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Result()
	return res, err
}

func Exists(ctx context.Context, key string) (bool, error) {
	cli := GetClient()
	res, err := cli.Exists(ctx, key).Result()
	return res > 0, err
}
