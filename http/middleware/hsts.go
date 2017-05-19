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
)

const (
	hstsHeader     = "Strict-Transport-Security"
	hstsMaxAge     = "max-age="
	hstsIncludeSub = "; includeSubDomains"
)

// HSTS sends HTTP Strict Transport Security header.
//
// By default the "includeSubDomains" is added; which is recommended.
//
// Reference: https://www.owasp.org/index.php/HTTP_Strict_Transport_Security_Cheat_Sheet
func HSTS(handler http.Handler, maxAge time.Duration, excludeSubDomains ...bool) http.HandlerFunc {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		value := hstsMaxAge + strconv.FormatInt(int64(maxAge.Seconds()), 10)
		if len(excludeSubDomains) == 0 || !excludeSubDomains[0] {
			value += hstsIncludeSub
		}

		resp.Header().Set(hstsHeader, value)

		handler.ServeHTTP(resp, req)
	})
}
