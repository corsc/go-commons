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

package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/corsc/go-commons/cache"
	"github.com/corsc/go-commons/http/middleware/internal"
)

var errCacheHandlerFailed = errors.New("response code from handler was not HTTP200 skipping cache")

// Cache will store the response bytes into cache based on the key generated supplied.
//
// Note: only responses with `http.StatusOK` response code will be cached
func Cache(handler http.Handler, cacheClient *cache.Client, generator func(*http.Request) string) http.HandlerFunc {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		key := generator(req)
		item := &cacheItem{}

		err := cacheClient.Get(req.Context(), key, item, cache.BuilderFunc(func(ctx context.Context, key string, dest cache.BinaryEncoder) error {
			// create a new response writer to catch the response body
			crw := internal.NewCustomResponseWriter(resp, true)

			handler.ServeHTTP(crw, req)

			if crw.Status() != http.StatusOK {
				// Skip non 200's
				return errCacheHandlerFailed
			}

			item.data = crw.Body()

			return nil
		}))

		if err != cache.ErrCacheMiss {
			// on cache hit, send the cached bytes and HTTP200
			_, err = resp.Write(item.data)
			if err != nil {
				http.Error(resp, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			resp.WriteHeader(http.StatusOK)
		}
	})
}

type cacheItem struct {
	data []byte
}

// MarshalBinary implements cache.BinaryEncoder
func (c *cacheItem) MarshalBinary() (data []byte, err error) {
	return c.data, nil
}

// UnmarshalBinary implements cache.BinaryEncoder
func (c *cacheItem) UnmarshalBinary(data []byte) error {
	c.data = data

	return nil
}
