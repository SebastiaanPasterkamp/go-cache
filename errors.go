package cache

import (
	"fmt"
)

var (
	// ErrMissingConfig is the error returned if the Configuration does not
	// contain settings for a supported cache.Repository type.
	ErrMissingConfig = fmt.Errorf("no configuration for any cache repository type")
	// ErrBadConfig is the error returned if the factory cannot instantiate a
	// cache repository based on the provided configuration.
	ErrBadConfig = fmt.Errorf("cannot create cache repository from configuration")
	// ErrNotFound is the error returned if the requested item does not exist in
	// the cache repository.
	ErrNotFound = fmt.Errorf("item not found")
	// ErrNotStored is the error returned if the stored item cannot be stored in
	// the repository.
	ErrNotStored = fmt.Errorf("item not stored")
	// ErrNotRecovered  is the error returned if the stored item cannot be
	// retrieved or the retrieved item is malformed.
	ErrNotRecovered = fmt.Errorf("failed to recover item")
)
