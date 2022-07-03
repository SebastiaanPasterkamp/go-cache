package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Memory is an in-memory implementation of the cache.Repository.
type Memory struct {
	mu sync.Mutex
	db map[string]memoryItem
}

// New returns a new cache.Memory instance based on the provided configuration.
func (*InMemorySettings) New(ctx context.Context) (*Memory, error) {
	return &Memory{
		mu: sync.Mutex{},
		db: map[string]memoryItem{},
	}, nil
}

// Set adds a new item to the in-memory cache.
func (m *Memory) Set(_ context.Context, key string, value interface{}, TTL time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrNotStored, err)
	}

	item := memoryItem{
		data: data,
	}

	if TTL > 0 {
		ttl := time.Now().Add(TTL)
		item.ttl = &ttl
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.db[key] = item

	return nil
}

// Get returns an item from the in-memory cache if it is available and not
// expired. Finding an expired item triggers an async pruning event.
func (m *Memory) Get(_ context.Context, key string, value interface{}) error {
	m.mu.Lock()
	item, ok := m.db[key]
	m.mu.Unlock()

	if !ok {
		return ErrNotFound
	}

	if item.expired() {
		defer m.prune()
		return ErrNotFound
	}

	if err := json.Unmarshal(item.data, value); err != nil {
		return fmt.Errorf("%w: %v", ErrNotRecovered, err)
	}

	return nil
}

// Del removes an item from the in-memory cache.
func (m *Memory) Del(_ context.Context, keys ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, key := range keys {
		delete(m.db, key)
	}

	return nil
}

func (m *Memory) prune() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for key := range m.db {
		if m.db[key].expired() {
			delete(m.db, key)
		}
	}
}

type memoryItem struct {
	data []byte
	ttl  *time.Time
}

func (item memoryItem) expired() bool {
	if item.ttl == nil {
		return false
	}

	return item.ttl.Before(time.Now())
}
