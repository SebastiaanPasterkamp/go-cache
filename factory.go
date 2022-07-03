package cache

import (
	"context"
)

// Factory returns a* cache.Repository instance based on the provided
// Configuration.
func (cfg Configuration) Factory(ctx context.Context) (repo Repository, err error) {
	switch {
	case cfg.InMemorySettings != nil:
		repo, err = cfg.InMemorySettings.New(ctx)
	case cfg.RedisSettings != nil:
		repo, err = cfg.RedisSettings.New(ctx)
	default:
		err = ErrMissingConfig
	}

	return
}
