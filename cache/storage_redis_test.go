// +build redis

// Copyright 2017 Corey Scott http://www.sage42.org/
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	assert.Equal(t, errCacheMiss, resultErr)

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
	assert.Equal(t, context.Canceled, resultErr)
}

func TestRedisStorage_setWithCtxDone(t *testing.T) {
	// inputs
	ctx, cancelFn := context.WithCancel(context.Background())
	key := getTestKey()

	storage := getTestRedisStorage()

	// attempt to get with a cancelled context
	cancelFn()

	resultErr := storage.Set(ctx, key, []byte("this is foo"))
	assert.Equal(t, context.Canceled, resultErr)
}

func getTestRedisStorage() *RedisStorage {
	return &RedisStorage{
		Pool: &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", ":6379") },
		},
		TTL: 60 * time.Second,
	}
}
