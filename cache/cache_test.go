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
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_cacheHit(t *testing.T) {
	// inputs
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	key := getTestKey()
	dest := &myDTO{}
	data := []byte(`{"name": "bob", "email":"bob@home.com"}`)

	// build a client and mock storage
	storage := &MockStorage{}
	storage.On("Get", mock.Anything, key).Return(data, nil)

	metrics := &MockMetrics{}
	metrics.On("Track", CacheHit)

	client := &Client{
		Storage: storage,
		Metrics: metrics,
	}

	// make the call
	resultErr := client.Get(ctx, key, dest, BuilderFunc(func(ctx context.Context, key string, dest BinaryEncoder) error {
		return errors.New("not implemented")
	}))

	assert.Nil(t, resultErr)

	assert.True(t, storage.AssertExpectations(t))
	assert.True(t, metrics.AssertExpectations(t))
}

func TestClient_cacheMiss(t *testing.T) {
	// inputs
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	key := getTestKey()
	dest := &myDTO{}

	// build a client and mock storage
	storage := &MockStorage{}
	storage.On("Get", mock.Anything, key).Return(nil, errCacheMiss)
	storage.On("Set", mock.Anything, key, mock.Anything).Return(nil, nil)

	metrics := &MockMetrics{}
	metrics.On("Track", CacheMiss)

	client := &Client{
		Storage: storage,
		Metrics: metrics,
	}

	// make the call
	resultErr := client.Get(ctx, key, dest, BuilderFunc(func(ctx context.Context, key string, dest BinaryEncoder) error {
		concrete := dest.(*myDTO)
		concrete.Name = "bob"
		concrete.Email = "bob@home.com"

		return nil
	}))

	assert.Nil(t, resultErr)

	err := client.waitForPending(1 * time.Second)
	assert.Nil(t, err)

	assert.True(t, storage.AssertExpectations(t))
	assert.True(t, metrics.AssertExpectations(t))
}

func TestClient_Invalidate(t *testing.T) {
	// inputs
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	key := getTestKey()

	// build a client and mock storage
	storage := &MockStorage{}
	storage.On("Invalidate", mock.Anything, key).Return(nil)

	client := &Client{
		Storage: storage,
	}

	// make the call
	resultErr := client.Invalidate(ctx, key)

	assert.Nil(t, resultErr)

	assert.True(t, storage.AssertExpectations(t))
}

func TestClient_cacheLambdaError(t *testing.T) {
	// inputs
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	key := getTestKey()
	dest := &myDTO{}

	// build a client and mock storage
	storage := &MockStorage{}
	storage.On("Get", mock.Anything, key).Return(nil, errCacheMiss)

	metrics := &MockMetrics{}
	metrics.On("Track", CacheMiss)
	metrics.On("Track", CacheLambdaError)

	client := &Client{
		Storage: storage,
		Metrics: metrics,
	}

	// make the call
	resultErr := client.Get(ctx, key, dest, BuilderFunc(func(ctx context.Context, key string, dest BinaryEncoder) error {
		// simulate user lambda error
		return errors.New("something failed")
	}))

	assert.NotNil(t, resultErr)

	assert.True(t, storage.AssertExpectations(t))
	assert.True(t, metrics.AssertExpectations(t))
}

func TestClient_cacheCacheError(t *testing.T) {
	// inputs
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	key := getTestKey()
	dest := &myDTO{}

	// build a client and mock storage
	storage := &MockStorage{}
	storage.On("Get", mock.Anything, key).Return(nil, errors.New("something failed"))

	metrics := &MockMetrics{}
	metrics.On("Track", CacheError)

	client := &Client{
		Storage: storage,
		Metrics: metrics,
	}

	// make the call
	resultErr := client.Get(ctx, key, dest, BuilderFunc(func(ctx context.Context, key string, dest BinaryEncoder) error {
		return errors.New("not implemented")
	}))

	assert.NotNil(t, resultErr)

	assert.True(t, storage.AssertExpectations(t))
	assert.True(t, metrics.AssertExpectations(t))
}

func TestClient_cacheUnmarshalError(t *testing.T) {
	// inputs
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	key := getTestKey()
	dest := &myDTO{}
	data := []byte(`bad data format`)

	// build a client and mock storage
	storage := &MockStorage{}
	storage.On("Get", mock.Anything, key).Return(data, nil)

	metrics := &MockMetrics{}
	metrics.On("Track", CacheUnmarshalError)

	client := &Client{
		Storage: storage,
		Metrics: metrics,
	}

	// make the call
	resultErr := client.Get(ctx, key, dest, BuilderFunc(func(ctx context.Context, key string, dest BinaryEncoder) error {
		return errors.New("not implemented")
	}))

	assert.NotNil(t, resultErr)

	assert.True(t, storage.AssertExpectations(t))
	assert.True(t, metrics.AssertExpectations(t))
}

type myDTO struct {
	Name  string
	Email string
}

// MarshalBinary implements cache.BinaryEncoder
func (m *myDTO) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

// UnmarshalBinary implements cache.BinaryEncoder
func (m *myDTO) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

// extra method to make tests predictable - wait for pending writes
func (c *Client) waitForPending(maxWait time.Duration) error {
	ticker := time.NewTicker(10 * time.Millisecond)
	for range ticker.C {
		select {
		case <-time.After(maxWait):
			return errors.New("timeout waiting for pending writes")

		default:
			if atomic.LoadInt64(&c.pendingWrites) == 0 {
				return nil
			}
		}
	}

	return errors.New("cannot happen")
}

func getTestKey() string {
	return fmt.Sprintf("redis.test.key.%d", time.Now().UnixNano())
}
