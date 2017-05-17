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

package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCustomResponseWriter(t *testing.T) {
	writer := httptest.NewRecorder()
	myWriter := NewCustomResponseWriter(writer)
	assert.Implements(t, (*http.ResponseWriter)(nil), myWriter)
}

func TestCustomResponseWriter_WriteHeader(t *testing.T) {
	writer := httptest.NewRecorder()
	myWriter := NewCustomResponseWriter(writer)

	myWriter.WriteHeader(http.StatusBadRequest)
	assert.Equal(t, http.StatusBadRequest, myWriter.Status())
	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestCustomResponseWriter_Write(t *testing.T) {
	writer := httptest.NewRecorder()
	myWriter := NewCustomResponseWriter(writer)

	myWriter.Write([]byte("foo"))
	assert.Equal(t, http.StatusOK, myWriter.Status())
}
