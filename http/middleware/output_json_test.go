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
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOutputJSON(t *testing.T) {
	dto := &fooResponse{
		Name:  "bob",
		Email: "bob@home.com",
		Age:   12,
	}

	expected := `{"name":"bob","email":"bob@home.com","age":12}
`

	resp := httptest.NewRecorder()
	_ = OutputJSON(resp, dto)

	assert.Equal(t, expected, resp.Body.String())
}

type fooResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}
