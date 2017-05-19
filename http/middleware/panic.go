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
	"runtime"
)

// Panic will catch all panics, optionally log the stack trace and set the response to `http.StatusInternalServerError`
func Panic(handler http.Handler, logger ...func(format string, args ...interface{})) http.HandlerFunc {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				if len(logger) > 0 {
					stackDump(r, logger[0])
				}

				http.Error(resp, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		handler.ServeHTTP(resp, req)
	})
}

func stackDump(r interface{}, logger func(format string, args ...interface{})) {
	_, file, line, _ := runtime.Caller(4)

	buf := make([]byte, 1<<10)
	runtime.Stack(buf, false)

	logger("PANIC: '%v': %s:%d\n", r, file, line)
	logger("%s\n", buf)
}
