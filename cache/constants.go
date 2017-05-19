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

const (
	// CbRedisStorage is tag for redis storage circuit breaker.
	// This should be used for in calls to `hystrix.ConfigureCommand()`
	CbRedisStorage = "CbRedisStorage"

	// redis commands
	redisGet    = "GET"
	redisSetex  = "SETEX"
	redisExpire = "EXPIRE"

	// CbDynamoDbStorage is tag for DynamoDB storage circuit breaker.
	// This should be used for in calls to `hystrix.ConfigureCommand()`
	CbDynamoDbStorage = "CbDynamoDbStorage"

	// dynamo constants
	ddbKey  = "key"
	ddbData = "data"
	ddbTTL  = "ttl"
)
