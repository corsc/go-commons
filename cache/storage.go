package cache

import (
	"context"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

// Storage is an abstract definition of the underlying cache storage
type Storage interface {
	// Set will attempt to get a value from storage
	Get(ctx context.Context, key string) ([]byte, error)

	// Set will save a value into storage
	Set(ctx context.Context, key string, bytes []byte) error

	// return the TTL duration in seconds
	GetTTL() int64
}

// RedigoStorage implements Storage
type RedigoStorage struct {
	// Pool is the redis connection pool (required)
	Pool *redis.Pool

	// TTL is the max TTL for cache items (optional but recommended)
	TTL time.Duration

	// calculated version of TTL
	ttlInSeconds int64
	ttlOnce      sync.Once
}

// Get implements Storage
func (r *RedigoStorage) Get(ctx context.Context, key string) ([]byte, error) {
	select {
	case resp := <-r.asyncDo("GET", key):
		return resp.reply, resp.err

	case <-ctx.Done():
		// context cancel/timeout
		return nil, ctx.Err()
	}
}

// Set implements Storage
func (r *RedigoStorage) Set(ctx context.Context, key string, bytes []byte) error {
	select {
	case resp := <-r.asyncDo("SETEX", key, r.GetTTL(), bytes):
		return resp.err

	case <-ctx.Done():
		// context cancel/timeout
		return ctx.Err()
	}
}

// GetTTL implements Storage
func (r *RedigoStorage) GetTTL() int64 {
	r.ttlOnce.Do(func() {
		r.ttlInSeconds = int64(r.TTL / time.Second)
	})

	return r.ttlInSeconds
}

func (r *RedigoStorage) asyncDo(command string, args ...interface{}) chan *redisResponse {
	redisRespCh := make(chan *redisResponse, 1)

	go func() {
		con := r.Pool.Get()

		output := &redisResponse{}
		output.reply, output.err = redis.Bytes(con.Do(command, args...))
		redisRespCh <- output
	}()

	return redisRespCh
}

type redisResponse struct {
	reply []byte
	err   error
}
