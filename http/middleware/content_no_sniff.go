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
	xctoHeader = "X-Content-Type-Options"
	xctoValue  = "nosniff"
)

// ContentNoSniff sends X-Content-Type-Options header.
//
// Reference: https://msdn.microsoft.com/en-us/library/gg622941%28v=vs.85%29.aspx
func ContentNoSniff(handler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set(xctoHeader, xctoValue)

		handler.ServeHTTP(resp, req)
	})
}
