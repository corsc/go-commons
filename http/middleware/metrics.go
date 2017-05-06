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
	"net/http"
	"strconv"
	"time"

	"github.com/corsc/go-commons/http/middleware/internal"
)

// MetricsClient allows for tracking the endpoint via StatsD or similar
type MetricsClient interface {
	Duration(key string, start time.Time, tags ...string)
}

// Duration will track the duration (and usage) of the method.
//
// It is based on statsD but it could be used with any capable metrics client.
func Duration(handler http.Handler, metrics MetricsClient, extraTags ...string) http.HandlerFunc {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		defer metrics.Duration("api", time.Now(), extraTags...)

		handler.ServeHTTP(resp, req)
	})
}

// DurationStatus is similar to Duration but also tracks the HTTP response code (via tags)
//
// It is based on statsD but it could be used with any capable metrics client.
func DurationStatus(handler http.Handler, metrics MetricsClient, extraTags ...string) http.HandlerFunc {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		start := time.Now()

		crw := internal.NewCustomResponseWriter(resp)

		handler.ServeHTTP(resp, req)

		metrics.Duration("api", start, append(extraTags, "status:"+strconv.Itoa(crw.Status()))...)
	})
}
