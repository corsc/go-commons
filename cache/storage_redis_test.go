// +build redis

package cache

import (
	"context"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/assert"
)

func TestRedisStorage_implements(t *testing.T) {
	assert.Implements(t, (*Storage)(nil), &RedisStorage{})
}

func TestRedisStorage_happyPath(t *testing.T) {
	// inputs
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	key := getTestKey()

	storage := getTestRedisStorage()

	// get a value (should fail)
	result, resultErr := storage.Get(ctx, key)
	assert.Nil(t, result)
	assert.Equal(t, resultErr, errCacheMiss)

	// set a value
	data := []byte(`this is foo`)
	resultErr = storage.Set(ctx, key, data)
	assert.Nil(t, resultErr)

	// get a value
	result, resultErr = storage.Get(ctx, key)
	assert.Equal(t, data, result)
	assert.Nil(t, resultErr)
}

func TestRedisStorage_getWithCtxDone(t *testing.T) {
	// inputs
	ctx, cancelFn := context.WithCancel(context.Background())
	key := getTestKey()

	storage := getTestRedisStorage()

	// attempt to get with a cancelled context
	cancelFn()

	result, resultErr := storage.Get(ctx, key)
	assert.Nil(t, result)
	assert.Equal(t, resultErr, context.Canceled)
}

func TestRedisStorage_setWithCtxDone(t *testing.T) {
	// inputs
	ctx, cancelFn := context.WithCancel(context.Background())
	key := getTestKey()

	storage := getTestRedisStorage()

	// attempt to get with a cancelled context
	cancelFn()

	resultErr := storage.Set(ctx, key, []byte("this is foo"))
	assert.Equal(t, resultErr, context.Canceled)
}

func getTestRedisStorage() *RedisStorage {
	return &RedisStorage{
		Pool: &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", "localhost:6379") },
		},
		TTL: 60 * time.Second,
	}
}
