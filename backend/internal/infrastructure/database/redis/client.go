package redis

import (
	"context"
	"time"

	redisv8 "github.com/go-redis/redis/v8"
)

// Client defines the interface for Redis operations
type Client interface {
	// Basic operations
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	MGet(ctx context.Context, keys ...string) ([]string, error)
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, keys ...string) (int64, error)

	// Set operations
	SAdd(ctx context.Context, key string, members ...interface{}) error
	SRem(ctx context.Context, key string, members ...interface{}) error
	SMembers(ctx context.Context, key string) ([]string, error)
	SCard(ctx context.Context, key string) (int64, error)
	SIsMember(ctx context.Context, key, member string) (bool, error)

	// Connection management
	Close() error
	Ping(ctx context.Context) error
}

// RealClient wraps go-redis/v8 client
type RealClient struct {
	cli *redisv8.Client
}

// NewRealClient creates a real Redis client from addr/password/db
func NewRealClient(addr, password string, db int) *RealClient {
	return &RealClient{
		cli: redisv8.NewClient(&redisv8.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}
}

func (r *RealClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.cli.Set(ctx, key, value, expiration).Err()
}

func (r *RealClient) Get(ctx context.Context, key string) (string, error) {
	return r.cli.Get(ctx, key).Result()
}

func (r *RealClient) MGet(ctx context.Context, keys ...string) ([]string, error) {
	vals, err := r.cli.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	res := make([]string, len(vals))
	for i, v := range vals {
		if v == nil {
			res[i] = ""
			continue
		}
		if s, ok := v.(string); ok {
			res[i] = s
			continue
		}
		res[i] = ""
	}
	return res, nil
}

func (r *RealClient) Del(ctx context.Context, keys ...string) error {
	return r.cli.Del(ctx, keys...).Err()
}

func (r *RealClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.cli.Exists(ctx, keys...).Result()
}

func (r *RealClient) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return r.cli.SAdd(ctx, key, members...).Err()
}

func (r *RealClient) SRem(ctx context.Context, key string, members ...interface{}) error {
	return r.cli.SRem(ctx, key, members...).Err()
}

func (r *RealClient) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.cli.SMembers(ctx, key).Result()
}

func (r *RealClient) SCard(ctx context.Context, key string) (int64, error) {
	return r.cli.SCard(ctx, key).Result()
}

func (r *RealClient) SIsMember(ctx context.Context, key, member string) (bool, error) {
	return r.cli.SIsMember(ctx, key, member).Result()
}

func (r *RealClient) Close() error {
	return r.cli.Close()
}

func (r *RealClient) Ping(ctx context.Context) error {
	return r.cli.Ping(ctx).Err()
}

// MockClient is a temporary mock implementation for development
type MockClient struct{}

func (m *MockClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return nil
}
func (m *MockClient) Get(ctx context.Context, key string) (string, error) { return "", nil }
func (m *MockClient) MGet(ctx context.Context, keys ...string) ([]string, error) {
	return []string{}, nil
}
func (m *MockClient) Del(ctx context.Context, keys ...string) error                      { return nil }
func (m *MockClient) Exists(ctx context.Context, keys ...string) (int64, error)          { return 0, nil }
func (m *MockClient) SAdd(ctx context.Context, key string, members ...interface{}) error { return nil }
func (m *MockClient) SRem(ctx context.Context, key string, members ...interface{}) error { return nil }
func (m *MockClient) SMembers(ctx context.Context, key string) ([]string, error) {
	return []string{}, nil
}
func (m *MockClient) SCard(ctx context.Context, key string) (int64, error) { return 0, nil }
func (m *MockClient) SIsMember(ctx context.Context, key, member string) (bool, error) {
	return false, nil
}
func (m *MockClient) Close() error                   { return nil }
func (m *MockClient) Ping(ctx context.Context) error { return nil }

// NewMockClient creates a new mock Redis client
func NewMockClient() Client { return &MockClient{} }
