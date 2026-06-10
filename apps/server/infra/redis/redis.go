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

func Get(ctx context.Context, key string) (string, error) {
	cli := GetClient()
	res, err := cli.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return res, err
}

func Set(ctx context.Context, key string, value string, ttlSeconds int) error {
	cli := GetClient()
	return cli.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
}

func Del(ctx context.Context, key string) error {
	cli := GetClient()
	return cli.Del(ctx, key).Err()
}

func DelByPattern(ctx context.Context, pattern string) error {
	cli := GetClient()
	var cursor uint64
	for {
		keys, next, err := cli.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err := cli.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}

		cursor = next
		if cursor == 0 {
			break
		}
	}

	return nil
}
