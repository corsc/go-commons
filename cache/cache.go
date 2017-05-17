package cache

import (
	"context"
	"encoding"
	"sync/atomic"
	"time"
)

// Client defines a cache instance.
//
// This can represent the cache for the entire system or for a particular use-case/type.
//
// If a cache is used for multiple purposes, then care must be taken to ensure uniqueness of cache keys.
//
// It is not recommended to change this struct's member data after creation as a data race will likely ensue.
type Client struct {
	// Storage is the cache storage scheme. (Required)
	Storage Storage

	// Logger defines a logger to used for errors during async cache writes (optional)
	Logger Logger

	// Metrics allow for tracking cache events (hit/miss/etc) (optional)
	Metrics Metrics

	// WriteTimeout is the max time spent waiting for cache writes to complete (optional - default 3 seconds)
	WriteTimeout time.Duration

	// track pending cache writes
	pendingWrites int64
}

// Get attempts to retrieve the value from cache and when it misses will run the builder func to create the value.
//
// It will asynchronously update/save the value in the cache on after a successful builder run
func (c *Client) Get(ctx context.Context, key string, dest BinaryEncoder, builder Builder) error {
	bytes, err := c.Storage.Get(ctx, key)
	if err != nil {
		if err == errCacheMiss {
			return c.onCacheMiss(ctx, key, dest, builder)
		}

		c.getMetrics().Track(CacheError)
		return err
	}

	return c.onCacheHit(dest, bytes)
}

func (c *Client) onCacheMiss(ctx context.Context, key string, dest BinaryEncoder, builder Builder) error {
	c.getMetrics().Track(CacheMiss)

	err := builder.Build(ctx, key, dest)
	if err != nil {
		c.getMetrics().Track(CacheLambdaError)
		return err
	}

	atomic.AddInt64(&c.pendingWrites, 1)
	go c.updateCache(key, dest)

	return err
}

func (c *Client) onCacheHit(dest encoding.BinaryUnmarshaler, bytes []byte) error {
	err := dest.UnmarshalBinary(bytes)
	if err != nil {
		c.getMetrics().Track(CacheUnmarshalError)
		return err
	}

	c.getMetrics().Track(CacheHit)
	return nil
}

// update the cache with the supplied key/value pair
func (c *Client) updateCache(key string, val encoding.BinaryMarshaler) {
	defer func() {
		// update tracking
		atomic.AddInt64(&c.pendingWrites, -1)
	}()

	// use independent context so we don't miss cache updated
	ctx, cancelFn := context.WithTimeout(context.Background(), c.getWriteTimeout())
	defer cancelFn()

	bytes, err := val.MarshalBinary()
	if err != nil {
		c.getLogger().Log("failed marshal '%s' from cache with err: %s", key, err)
		return
	}

	err = c.Storage.Set(ctx, key, bytes)
	if err != nil {
		c.getLogger().Log("failed to update item '%s' in cache with err: %s", key, err)
	}
}

// return the supplied logger or a no-op implementation
func (c *Client) getLogger() Logger {
	if c.Logger != nil {
		return c.Logger
	}

	return noopLogger
}

// return the supplied metric tracker or a no-op implementation
func (c *Client) getMetrics() Metrics {
	if c.Metrics != nil {
		return c.Metrics
	}

	return noopMetrics
}

// return the timeout on cache writes
func (c *Client) getWriteTimeout() time.Duration {
	if int64(c.WriteTimeout) > 0 {
		return c.WriteTimeout
	}

	return 3 * time.Second
}

// Builder builds the data for a key
type Builder interface {
	// Build returns the data for the supplied key by populating dest
	Build(ctx context.Context, key string, dest BinaryEncoder) error
}

// BuilderFunc implements Builder as a function
type BuilderFunc func(ctx context.Context, key string, dest BinaryEncoder) error

// Build implements Builder
func (b BuilderFunc) Build(ctx context.Context, key string, dest BinaryEncoder) error {
	return b(ctx, key, dest)
}

// BinaryEncoder encodes/decodes the receiver to and from binary form
type BinaryEncoder interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}
