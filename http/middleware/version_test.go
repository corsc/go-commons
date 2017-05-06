package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	handler := getTestHandler()
	headerTag := "version"
	headerValue := "1.2.3"

	handlerFunc := Version(handler, headerTag, headerValue)
	assert.IsType(t, (http.HandlerFunc)(nil), handlerFunc)

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/who-cares", nil)
	handlerFunc(resp, req)

	assert.True(t, handler.wasCalled())

	result := resp.Header().Get(headerTag)
	assert.Equal(t, headerValue, result)
}
