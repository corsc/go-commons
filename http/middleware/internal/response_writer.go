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

package internal

import (
	"bytes"
	"net/http"
)

// NewCustomResponseWriter returns a CustomResponseWriter from the supplied http.ResponseWriter
func NewCustomResponseWriter(resp http.ResponseWriter, catchBody bool) *CustomResponseWriter {
	return &CustomResponseWriter{
		ResponseWriter: resp,
		catchBody:      catchBody,
	}
}

// CustomResponseWriter implements http.ResponseWriter and extends the default implementation by allowing us
// access to the http response code
type CustomResponseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool

	catchBody bool
	body      bytes.Buffer
}

// Write implements http.ResponseWriter
func (resp *CustomResponseWriter) Write(in []byte) (int, error) {
	if !resp.wroteHeader {
		resp.WriteHeader(http.StatusOK)
	}

	if resp.catchBody {
		_, err := resp.body.Write(in)
		if err != nil {
			return 0, err
		}
	}

	return resp.ResponseWriter.Write(in)
}

// WriteHeader implements http.ResponseWriter
func (resp *CustomResponseWriter) WriteHeader(code int) {
	resp.ResponseWriter.WriteHeader(code)

	// check after in case there's error handling in the wrapped http.ResponseWriter
	if resp.wroteHeader {
		return
	}

	resp.status = code
	resp.wroteHeader = true
}

// Status returns the HTTP status code (or 0 if not set)
func (resp *CustomResponseWriter) Status() int {
	return resp.status
}

// Body returns the HTTP body (if `catchBody == true` in the constructor)
func (resp *CustomResponseWriter) Body() []byte {
	return resp.body.Bytes()
}
