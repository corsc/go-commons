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
)

const (
	xxssHeader = "X-XSS-Protection"
	xxssValue  = "1"
	xxssBlock  = "; mode=block"
)

// XXSS sends X-XSS-Protection response header
//
// By default the "mode=block" is not sent.
//
// Reference: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-XSS-Protection
func XXSS(handler http.Handler, block ...bool) http.HandlerFunc {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		value := xxssValue
		if len(block) > 0 && block[0] {
			value += xxssBlock
		}

		resp.Header().Set(xxssHeader, value)

		handler.ServeHTTP(resp, req)
	})
}
