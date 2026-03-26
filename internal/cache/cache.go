package cache

import (
	"context"
	"errors"
	"time"
)

type ICache interface {
	Set(ctx context.Context, k string, v any, ttl time.Duration) error
	Get(ctx context.Context, k string) (string, error)
	Delete(ctx context.Context, k string) error
	GetDel(ctx context.Context, k string) (string, error)
	Inc(ctx context.Context, k string) error
	Keys(ctx context.Context, pattern string) ([]string, error)
	Close() error
}

var (
	ErrCacheMiss = errors.New("cache miss")
)
