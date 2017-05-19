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

func TestXXSS(t *testing.T) {
	scenarios := []struct {
		desc      string
		extraArgs []bool
		expected  string
	}{
		{
			desc:      "no block",
			extraArgs: nil,
			expected:  "1",
		},
		{
			desc:      "with block",
			extraArgs: []bool{true},
			expected:  "1; mode=block",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			handler := getTestHandler()

			handlerFunc := XXSS(handler, scenario.extraArgs...)
			assert.IsType(t, (http.HandlerFunc)(nil), handlerFunc)

			resp := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/who-cares", nil)
			handlerFunc(resp, req)

			assert.True(t, handler.wasCalled())

			result := resp.Header().Get(xxssHeader)
			assert.Equal(t, scenario.expected, result)
		})
	}
}
