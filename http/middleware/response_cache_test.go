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

	"github.com/corsc/go-commons/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCache_cacheMiss(t *testing.T) {
	handler := getTestHandler()

	cacheKey := "ID"

	mockStorage := &cache.MockStorage{}
	mockStorage.On("Get", mock.Anything, cacheKey).Return(nil, cache.ErrCacheMiss).Once()
	mockStorage.On("Set", mock.Anything, cacheKey, mock.Anything).Return(nil).Once()

	cacheClient := &cache.Client{
		Storage: mockStorage,
	}

	keyGenerator := func(req *http.Request) string {
		return cacheKey
	}

	handlerFunc := Cache(handler, cacheClient, keyGenerator)
	assert.IsType(t, (http.HandlerFunc)(nil), handlerFunc)

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/who-cares", nil)
	handlerFunc(resp, req)

	assert.True(t, handler.wasCalled())
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestCache_cacheHit(t *testing.T) {
	handler := getTestHandler()

	cacheKey := "ID"

	mockStorage := &cache.MockStorage{}
	mockStorage.On("Get", mock.Anything, cacheKey).Return([]byte("foo"), nil).Once()

	cacheClient := &cache.Client{
		Storage: mockStorage,
	}

	keyGenerator := func(req *http.Request) string {
		return cacheKey
	}

	handlerFunc := Cache(handler, cacheClient, keyGenerator)
	assert.IsType(t, (http.HandlerFunc)(nil), handlerFunc)

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/who-cares", nil)
	handlerFunc(resp, req)

	assert.False(t, handler.wasCalled())
	assert.Equal(t, http.StatusOK, resp.Code)
}
