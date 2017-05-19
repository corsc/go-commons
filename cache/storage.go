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
)

// Storage is an abstract definition of the underlying cache storage
//go:generate mockery -name Storage -inpkg -case underscore
type Storage interface {
	// Get will attempt to get a value from storage
	Get(ctx context.Context, key string) ([]byte, error)

	// Set will save a value into storage
	Set(ctx context.Context, key string, bytes []byte) error

	// Invalidate will force invalidate/remove a key from storage
	Invalidate(ctx context.Context, key string) error
}
