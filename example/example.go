package main

import (
	"context"
	"log"
	"time"

	"github.com/SebastiaanPasterkamp/go-cache"
)

func main() {
	ctx := context.Background()

	// This configuration is for demonstration purpose only. Only one of the
	// settings will be used to initialize the cache where in-memory takes
	// the priority.
	test := cache.Configuration{
		// Demonstrate the in-memory cache implementation, which should work
		// as an example
		InMemorySettings: &cache.InMemorySettings{},
		// For real-world use-cases use the Redis settings instead
		RedisSettings: &cache.RedisSettings{
			Address:  "redis-host:6379",
			Password: "S3(r37!", // Please configure this using an environment variable
			Database: 0,
		},
	}

	c, err := test.Factory(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize cache: %v", err)
	}

	err = c.Set(ctx, "key", "value", 5*time.Second)
	if err != nil {
		log.Fatalf("Failed to store cache item: %v", err)
	}

	value := ""
	err = c.Get(ctx, "key", &value)
	if err != nil {
		log.Fatalf("Failed to retrieve cache item: %v", err)
	}

	err = c.Del(ctx, "key")
	if err != nil {
		log.Fatalf("Failed to delete cache item: %v", err)
	}

	log.Printf("Value retrieved: %q\n", value)
}
