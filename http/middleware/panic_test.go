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
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPanic(t *testing.T) {
	handler := getPanicingHandler()

	loggedMessage := ""
	logFn := func(format string, args ...interface{}) {
		loggedMessage += fmt.Sprintf(format, args...)
	}

	handlerFunc := Panic(handler, logFn)
	assert.IsType(t, (http.HandlerFunc)(nil), handlerFunc)

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/who-cares", nil)
	handlerFunc(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.NotEqual(t, "", loggedMessage)
}
