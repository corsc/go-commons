package middleware

import "net/http"

// returns a HTTP handler with properties that track calls and call values
func getTestHandler() *testHandler {
	return &testHandler{}
}

// handle foo requests
type testHandler struct {
	calls []*http.Request
}

// ServeHTTP implements the http.Handler
func (h *testHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	h.calls = append(h.calls, req)
}

func (h *testHandler) wasCalled() bool {
	return len(h.calls) > 0
}
