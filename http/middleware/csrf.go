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
	"strings"
)

// CSRF sends custom header that in turn causes the request to considered "complex" and therefore CORS will apply.
//
// Reference: https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)_Prevention_Cheat_Sheet#Protecting_REST_Services:_Use_of_Custom_Request_Headers
func CSRF(headerKey string, handler http.Handler, logger CSRFLogger, ignoredPaths ...string) http.HandlerFunc {
	ignoredPathsLower := convertPathsToLower(ignoredPaths)

	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		currentPath := strings.ToLower(req.URL.Path)

		for _, ignoredPath := range ignoredPathsLower {
			if currentPath == ignoredPath {
				handler.ServeHTTP(resp, req)
				return
			}
		}

		if req.Header.Get(headerKey) == "" {
			if logger != nil {
				logger.BadRequest(req, "[CSRF] request to '%s' failed due to missing header '%s'", req.URL.Path, headerKey)
			}
			http.Error(resp, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		handler.ServeHTTP(resp, req)
	})
}

func convertPathsToLower(in []string) []string {
	out := make([]string, len(in))
	for index, path := range in {
		out[index] = strings.ToLower(path)
	}
	return out
}

// CSRFLogger allows for logging
type CSRFLogger interface {
	// log requests that fail CSRF check
	BadRequest(req *http.Request, msg string, args ...interface{})
}
