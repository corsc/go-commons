package cache

import (
	"context"
	"sync"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/garyburd/redigo/redis"
)

// RedisStorage implements Storage
//
// It is strongly recommended that users customize the circuit breaker settings with a call similar to:
//
//    hystrix.ConfigureCommand(cache.CbRedisStorage, hystrix.CommandConfig{
//        Timeout: 1 * 1000,
//        MaxConcurrentRequests: 1000,
//        ErrorPercentThreshold: 50,
//        })
//
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
	resp, err := r.do(ctx, redisGet, key)
	if err != nil {
		return nil, err
	}

	bytes, err := redis.Bytes(resp, err)
	if err == redis.ErrNil {
		return nil, errCacheMiss
	}
	return bytes, err
}

// Set implements Storage
func (r *RedisStorage) Set(ctx context.Context, key string, bytes []byte) error {
	_, err := r.do(ctx, redisSetex, key, r.GetTTL(), bytes)
	return err
}

// GetTTL implements Storage
func (r *RedisStorage) GetTTL() int64 {
	r.ttlOnce.Do(func() {
		r.ttlInSeconds = int64(r.TTL / time.Second)
	})

	return r.ttlInSeconds
}

// calls to redis protected by a circuit breaker
func (r *RedisStorage) do(ctx context.Context, command string, args ...interface{}) (interface{}, error) {
	resultCh := make(chan interface{}, 1)
	errorCh := hystrix.Go(CbRedisStorage, func() error {
		con := r.Pool.Get()

		reply, err := con.Do(command, args...)
		if err != nil {
			return err
		}

		resultCh <- reply

		return nil
	}, nil)

	select {
	case result := <-resultCh:
		// success
		return result, nil

	case <-ctx.Done():
		// timeout/context cancelled
		return nil, ctx.Err()

	case err := <-errorCh:
		// failure
		return nil, err
	}
}
