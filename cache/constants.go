package cache

import (
	"errors"
)

// errCacheMiss is returned when the cache does not contain the requested key
var errCacheMiss = errors.New("cache miss")

// Event denote the cache event type
type Event int

const (
	// CacheHit denotes the key was found in the cache and successfully returned
	CacheHit Event = iota

	// CacheMiss denotes the key was not found in the cache
	CacheMiss

	// CacheError denotes an error occurred with the cache or underlying storage
	CacheError

	// CacheLambdaError denotes an error occurred in the user code (e.g. the Builder caller)
	// Note: The builder is only called after a cache miss and therefore CacheLambdaError events should not be counted
	// towards "total cache usage" or used in "cache hit rate" calcuations
	CacheLambdaError

	// CacheUnmarshalError denotes an error occurred in while calling BinaryEncoder.UnmarshalBinary
	//
	// If the BinaryEncoder is implemented correctly, this event should never happen
	CacheUnmarshalError
)
