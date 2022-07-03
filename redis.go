package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// redisClient is a subset of the redis.Cmdable interface limited to only the
// functions used in the cache.Redis implementation.
type redisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

// Redis is an redis based implementation of the cache.Repository.
// If Address is 'mock', then the redis client is replaced with the in-memory
// cache implementation.
type Redis struct {
	rdb redisClient
}

// New returns a new cache.Memory instance based on the provided configuration.
func (cfg *RedisSettings) New(ctx context.Context) (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.Database,
	})

	if cfg.Address == "mock" {
		mem := &InMemorySettings{}
		mock, err := mem.New(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create mock: %v", err)
		}

		return &Redis{rdb: &mockRedis{db: mock}}, err
	}

	err := rdb.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBadConfig, err)
	}

	return &Redis{rdb: rdb}, err
}

// Set adds a new item to the redis cache.
func (r *Redis) Set(ctx context.Context, key string, value interface{}, TTL time.Duration) error {
	obj, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrNotStored, err)
	}

	if err := r.rdb.Set(ctx, key, obj, TTL).Err(); err != nil {
		return fmt.Errorf("%w: %v", ErrNotStored, err)
	}

	return nil
}

// Get returns an item from the redis cache if it is available and not
// expired.
func (r *Redis) Get(ctx context.Context, key string, value interface{}) error {
	obj, err := r.rdb.Get(ctx, key).Result()
	switch err {
	case redis.Nil:
		return ErrNotFound
	case nil:
	default:
		return fmt.Errorf("%w: %v", ErrNotRecovered, err)
	}

	if err := json.Unmarshal([]byte(obj), value); err != nil {
		return fmt.Errorf("%w: %v", ErrNotRecovered, err)
	}

	return nil
}

// Del removes an item from the redis cache.
func (r *Redis) Del(ctx context.Context, keys ...string) error {
	if err := r.rdb.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("failed to delete from redis: %w", err)
	}

	return nil
}

// mockRedis uses the in-memory cache.Repository implementation to validate the
// interaction with the redis client.
type mockRedis struct {
	db *Memory
}

func (m *mockRedis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	bytes, ok := value.([]byte)
	if !ok {
		panic("mockRedis expects to be storing a byte array castable to string")
	}

	err := m.db.Set(ctx, key, string(bytes), expiration)
	cmd := redis.NewStatusCmd(ctx, "set", key, string(bytes), expiration)
	cmd.SetErr(err)
	return cmd
}

func (m *mockRedis) Get(ctx context.Context, key string) *redis.StringCmd {
	value := ""
	err := m.db.Get(ctx, key, &value)
	cmd := redis.NewStringCmd(ctx, "get", key)
	cmd.SetVal(string(value))
	if errors.Is(err, ErrNotFound) {
		cmd.SetErr(redis.Nil)
	} else {
		cmd.SetErr(err)
	}
	return cmd
}

func (m *mockRedis) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	err := m.db.Del(ctx, keys...)
	args := make([]interface{}, 1+len(keys))
	args[0] = "del"
	for i, key := range keys {
		args[1+i] = key
	}
	cmd := redis.NewIntCmd(ctx, args...)
	cmd.SetErr(err)
	return cmd
}
