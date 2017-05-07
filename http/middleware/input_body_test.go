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
// limitations un

package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInputBody(t *testing.T) {
	// define a special handler to inspect the request context
	wasCalled := false
	handler := func(_ http.ResponseWriter, req *http.Request) {
		requestDTO := InputBodyDTO(req).(*fooRequest)
		assert.NotNil(t, requestDTO)
		assert.Equal(t, "bob", requestDTO.Name)
		assert.Equal(t, "bob@home.com", requestDTO.Email)
		assert.Equal(t, 23, requestDTO.Age)
		wasCalled = true
	}

	handlerFunc := InputBody(http.HandlerFunc(handler), &fooRequest{}, &myLogger{})
	assert.IsType(t, (http.HandlerFunc)(nil), handlerFunc)

	requestBody := `
{
	"name" : "bob",
	"email" : "bob@home.com",
	"age": 23
}`

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/who-cares", bytes.NewBufferString(requestBody))
	handlerFunc(resp, req)

	assert.True(t, wasCalled)
}

type fooRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}
