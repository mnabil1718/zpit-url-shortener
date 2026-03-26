package cache

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

// Cache mock
type MockCache struct{ mock.Mock }

func (m *MockCache) Set(ctx context.Context, k string, v any, ttl time.Duration) error {
	return m.Called(ctx, k, v, ttl).Error(0)
}

func (m *MockCache) Get(ctx context.Context, k string) (string, error) {
	args := m.Called(ctx, k)
	return args.String(0), args.Error(1)
}

func (m *MockCache) Delete(ctx context.Context, k string) error {
	return m.Called(ctx, k).Error(0)
}

func (m *MockCache) GetDel(ctx context.Context, k string) (string, error) {
	args := m.Called(ctx, k)
	return args.String(0), args.Error(1)
}

func (m *MockCache) Inc(ctx context.Context, k string) error {
	return m.Called(ctx, k).Error(0)
}

func (m *MockCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	args := m.Called(ctx, pattern)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockCache) Close() error {
	return m.Called().Error(0)
}
