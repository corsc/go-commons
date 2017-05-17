package cache

import (
	"context"
	"encoding"
	"errors"
)

// ErrCacheMiss is returned when the cache does not contain the requested key
var ErrCacheMiss = errors.New("cache miss")

type Client struct {
	Storage Storage
	Logger  LoggerFn
}

func (c *Client) Get(ctx context.Context, key string, dest Binary, builder Builder) error {
	bytes, err := c.Storage.Get(ctx, key)
	if err != nil {
		if err != ErrCacheMiss {
			return err
		}

		err = builder.Build(ctx, key, dest)
		if err != nil {
			return err
		}

		// async update cache
		go c.updateCache(ctx, key, dest)

		return nil
	}

	return dest.UnmarshalBinary(bytes)
}

func (c *Client) updateCache(ctx context.Context, key string, val Binary) {
	bytes, err := val.MarshalBinary()
	if err != nil {
		c.getLogger()("failed marshal '%s' from cache with err: %s", key, err)
		return
	}

	err = c.Storage.Set(ctx, key, bytes)
	if err != nil {
		c.getLogger()("failed to update item '%s' in cache with err: %s", key, err)
	}
}

func (c *Client) getLogger() LoggerFn {
	if c.Logger != nil {
		return c.Logger
	}

	return func(message string, args ...interface{}) {
		// skip logging
	}
}

type LoggerFn func(message string, args ...interface{})

type Builder interface {
	Build(ctx context.Context, key string, dest Binary) error
}

// BuilderFunc implements Builder as a function
type BuilderFunc func(ctx context.Context, key string, dest Binary)

// Build implements Builder
func (b BuilderFunc) Build(ctx context.Context, key string, dest Binary) error {
	return b.Build(ctx, key, dest)
}

type Binary interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}
