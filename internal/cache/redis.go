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

	opts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(opts)

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

func (r *RedisClient) GetDel(ctx context.Context, k string) (string, error) {
	v, err := r.client.GetDel(ctx, k).Result() // avoid race condition for goroutines
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

func (r *RedisClient) Inc(ctx context.Context, k string) error {
	return r.client.Incr(ctx, k).Err()
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}
