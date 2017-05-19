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

package middleware_test

import (
	"fmt"
	"net/http"
	"time"
)

// handle foo requests
func fooHandler(resp http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(resp, "hello foo")
}

// a handler that always panics
func panicHandler(_ http.ResponseWriter, _ *http.Request) {
	panic("foo")
}

// simple implementation of the MetricsClient interface
type myMetricsClient struct{}

// Duration implements MetricsClient
func (m *myMetricsClient) Duration(key string, start time.Time, tags ...string) {
	// send metrics to server
}
