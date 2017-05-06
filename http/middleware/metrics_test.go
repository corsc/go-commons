package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDuration(t *testing.T) {
	handler := getTestHandler()
	metricsClient := &mockMetricsClient{}

	handlerFunc := Duration(handler, metricsClient)
	assert.IsType(t, (http.HandlerFunc)(nil), handlerFunc)

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/who-cares", nil)
	handlerFunc(resp, req)

	assert.True(t, handler.wasCalled())
	assert.True(t, metricsClient.wasCalled())
}

func TestDurationStatus(t *testing.T) {
	handler := getTestHandler()
	metricsClient := &mockMetricsClient{}

	handlerFunc := DurationStatus(handler, metricsClient)
	assert.IsType(t, (http.HandlerFunc)(nil), handlerFunc)

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/who-cares", nil)
	handlerFunc(resp, req)

	assert.True(t, handler.wasCalled())

	assert.True(t, metricsClient.wasCalled())
	assert.True(t, metricsClient.callHasTag(0, "status"))
}

// mock implementation of MetricsClient
type mockMetricsClient struct {
	calls []map[string]interface{}
}

// implements MetricsClient
func (mock *mockMetricsClient) Duration(key string, start time.Time, tags ...string) {
	call := map[string]interface{}{
		"key":   key,
		"start": start,
		"tags":  tags,
	}

	mock.calls = append(mock.calls, call)
}

func (mock *mockMetricsClient) wasCalled() bool {
	return len(mock.calls) > 0
}

func (mock *mockMetricsClient) callHasTag(callNo int, tag string) bool {
	if len(mock.calls) <= callNo {
		return false
	}

	call := mock.calls[callNo]

	tags := call["tags"].([]string)
	for _, thisTag := range tags {
		if strings.HasPrefix(thisTag, tag) {
			return true
		}
	}
	return false
}
