package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/mnabil1718/zp.it/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(cfg *config.Config) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL,
		Password: "",
		DB:       0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Sprintf("redis connection error: %v", err))
	}

	return &RedisClient{client: rdb}
}

func (r *RedisClient) Set(ctx context.Context, k string, v any, ttl time.Duration) error {
	return r.client.Set(ctx, k, v, ttl).Err()
}

func (r *RedisClient) Get(ctx context.Context, k string) (string, error) {
	v, err := r.client.Get(ctx, k).Result()
	if err == redis.Nil {
		return "", ErrCacheMiss
	}

	if err != nil {
		return "", err
	}

	return v, nil
}

func (r *RedisClient) Delete(ctx context.Context, k string) error {
	return r.client.Del(ctx, k).Err()
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}
