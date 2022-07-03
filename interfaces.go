package cache

import (
	"context"
	"time"
)

// Repository is a simple interface to any cache offering Set, Get, and Del
// operations.
type Repository interface {
	Set(ctx context.Context, key string, value interface{}, TTL time.Duration) error
	Get(ctx context.Context, key string, value interface{}) error
	Del(ctx context.Context, keys ...string) error
}
