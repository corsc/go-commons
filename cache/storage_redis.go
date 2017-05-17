package cache

import (
	"context"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

// RedisStorage implements Storage
type RedisStorage struct {
	// Pool is the redis connection pool (required)
	Pool *redis.Pool

	// TTL is the max TTL for cache items (required)
	TTL time.Duration

	// calculated version of TTL
	ttlInSeconds int64
	ttlOnce      sync.Once
}

// Get implements Storage
func (r *RedisStorage) Get(ctx context.Context, key string) ([]byte, error) {
	select {
	case resp := <-r.asyncDo("GET", key):
		bytes, err := redis.Bytes(resp.reply, resp.err)
		if err == redis.ErrNil {
			return nil, errCacheMiss
		}
		return bytes, err

	case <-ctx.Done():
		// context cancel/timeout
		return nil, ctx.Err()
	}
}

// Set implements Storage
func (r *RedisStorage) Set(ctx context.Context, key string, bytes []byte) error {
	select {
	case resp := <-r.asyncDo("SETEX", key, r.GetTTL(), bytes):
		return resp.err

	case <-ctx.Done():
		// context cancel/timeout
		return ctx.Err()
	}
}

// GetTTL implements Storage
func (r *RedisStorage) GetTTL() int64 {
	r.ttlOnce.Do(func() {
		r.ttlInSeconds = int64(r.TTL / time.Second)
	})

	return r.ttlInSeconds
}

// Asynchronously perform a redis command
//
// This method is used to add "context" support to redigo calls
func (r *RedisStorage) asyncDo(command string, args ...interface{}) chan *redisResponse {
	redisRespCh := make(chan *redisResponse, 1)

	go func() {
		con := r.Pool.Get()

		output := &redisResponse{}
		output.reply, output.err = con.Do(command, args...)
		redisRespCh <- output
	}()

	return redisRespCh
}

// dto for passing responses out of RedisStorage.asyncDo
type redisResponse struct {
	reply interface{}
	err   error
}
