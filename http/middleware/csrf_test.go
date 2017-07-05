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
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSRF(t *testing.T) {
	headerKey := "x-my-site"
	testURL := "/who-cares"

	scenarios := []struct {
		desc         string
		setup        func(req *http.Request)
		ignoredPaths []string
		expected     int
	}{
		{
			desc: "happy path",
			setup: func(req *http.Request) {
				req.Header.Set(headerKey, "ok")
			},
			expected: http.StatusOK,
		},
		{
			desc:         "happy path - path ignored",
			setup:        func(req *http.Request) {},
			ignoredPaths: []string{testURL},
			expected:     http.StatusOK,
		},
		{
			desc:     "header missing / not set",
			setup:    func(req *http.Request) {},
			expected: http.StatusBadRequest,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			handler := getTestHandler()

			handlerFunc := CSRF(headerKey, handler, scenario.ignoredPaths...)
			assert.IsType(t, (http.HandlerFunc)(nil), handlerFunc)

			resp := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", testURL, nil)
			scenario.setup(req)

			handlerFunc(resp, req)
			assert.Equal(t, scenario.expected, resp.Code, scenario.desc)
		})
	}
}
