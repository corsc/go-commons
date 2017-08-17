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
	"fmt"
	"net/http"
)

// returns a HTTP handler with properties that track calls and call values
func getTestHandler() *testHandler {
	return &testHandler{}
}

// handle foo requests
type testHandler struct {
	calls []*http.Request
}

// ServeHTTP implements the http.Handler
func (h *testHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	h.calls = append(h.calls, req)
}

func (h *testHandler) wasCalled() bool {
	return len(h.calls) > 0
}

// simple implementation of the LoggingClient interface
type myLogger struct{}

// BadRequest implements CSRFLogger and InputBodyLogger
func (l *myLogger) BadRequest(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
}

// Request implements InputBodyLogger
func (l *myLogger) Request(body []byte) {
	fmt.Printf("Body: %s\n", string(body))
}

// returns a HTTP handler that always panics
func getPanicingHandler() http.HandlerFunc {
	return func(_ http.ResponseWriter, _ *http.Request) {
		panic("foo")
	}
}
