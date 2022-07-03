package cache_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/SebastiaanPasterkamp/go-cache"
)

func TestFactory(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	testCases := []struct {
		name          string
		cfg           cache.Configuration
		expectedType  cache.Repository
		expectedError error
	}{
		{"Fail without settings", cache.Configuration{},
			nil, cache.ErrMissingConfig},
		{"InMemory before Redis", cache.Configuration{
			InMemorySettings: &cache.InMemorySettings{},
			RedisSettings:    &cache.RedisSettings{Address: "mock"},
		},
			&cache.Memory{}, nil},
		{"Init InMemory", cache.Configuration{
			InMemorySettings: &cache.InMemorySettings{},
		},
			&cache.Memory{}, nil},
		{"Init Redis", cache.Configuration{
			RedisSettings: &cache.RedisSettings{Address: "mock"},
		},
			&cache.Redis{}, nil},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c, err := tt.cfg.Factory(ctx)
			if !errors.Is(err, tt.expectedError) {
				t.Fatalf("Unexpected error. Expected %v, got %v.",
					tt.expectedError, err)
			}

			if reflect.TypeOf(c) != reflect.TypeOf(tt.expectedType) {
				t.Fatalf("Unexpected implementation. Expected %T, got %T.",
					tt.expectedType, c)
			}
		})
	}
}
