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
	"encoding"
)

// API defines the main API for this package
type API interface {
	// Get attempts to retrieve the value from cache and when it misses will run the builder func to create the value.
	Get(ctx context.Context, key string, dest BinaryEncoder, builder Builder) error

	// Set will update the cache with the supplied key/value pair
	// NOTE: generally this need not be called is it is called implicitly by Get
	Set(ctx context.Context, key string, val encoding.BinaryMarshaler)

	// Invalidate will force invalidate any matching key in the cache
	Invalidate(ctx context.Context, key string) error
}
