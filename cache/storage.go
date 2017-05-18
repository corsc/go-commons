package cache

import (
	"context"
)

// Storage is an abstract definition of the underlying cache storage
//go:generate mockery -name Storage -inpkg -testonly -case underscore
type Storage interface {
	// Get will attempt to get a value from storage
	Get(ctx context.Context, key string) ([]byte, error)

	// Set will save a value into storage
	Set(ctx context.Context, key string, bytes []byte) error
}
