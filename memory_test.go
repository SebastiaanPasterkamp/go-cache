package cache_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/SebastiaanPasterkamp/go-cache"
)

func TestInMemoryGetSetDelete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	testCases := []struct {
		name          string
		inputKey      string
		inputValue    string
		inputTTL      time.Duration
		sleep         time.Duration
		outputKey     string
		expectedValue string
		expectedError error
	}{
		{"Working", "working", "good", 5 * time.Second,
			0, "working", "good", nil},
		{"Expired", "expired", "gone", 1 * time.Second,
			2 * time.Second, "expired", "", cache.ErrNotFound},
		{"Unknown", "working", "good", 5 * time.Second,
			0, "unknown", "", cache.ErrNotFound},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := cache.Configuration{
				InMemorySettings: &cache.InMemorySettings{},
			}

			c, err := cfg.Factory(ctx)
			if err != nil {
				t.Fatalf("Failed to initialize cache: %v", err)
			}

			c.Set(ctx, tt.inputKey, tt.inputValue, tt.inputTTL)

			time.Sleep(tt.sleep)

			value := ""
			err = c.Get(ctx, tt.outputKey, &value)
			if !errors.Is(err, tt.expectedError) {
				t.Fatalf("Unexpected error. Expected %v, got %v.",
					tt.expectedError, err)
			}

			if value != tt.expectedValue {
				t.Fatalf("Unexpected value. Expected %v, got %v.",
					tt.expectedValue, value)
			}

			c.Del(ctx, tt.inputKey)
			value = ""

			err = c.Get(ctx, tt.inputKey, &value)
			if !errors.Is(err, cache.ErrNotFound) {
				t.Fatalf("Unexpected error. Expected %v, got %v.",
					cache.ErrNotFound, err)
			}

			if value != "" {
				t.Fatalf("Unexpected value. Expected %v, got %v.",
					"", value)
			}
		})
	}
}
