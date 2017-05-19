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

	"github.com/corsc/go-commons/http/middleware/csp"
	"github.com/stretchr/testify/assert"
)

func TestCSP(t *testing.T) {
	scenarios := []struct {
		desc     string
		config   *csp.Config
		expected string
	}{
		{
			desc: "default",
			config: &csp.Config{
				Default: []string{"https:"},
			},
			expected: "default-src https:;",
		},
		{
			desc: "base",
			config: &csp.Config{
				Base: []string{csp.Self, "https://*.example.com/"},
			},
			expected: "base-uri 'self' https://*.example.com/;",
		},
		{
			desc: "child",
			config: &csp.Config{
				Child: []string{csp.Self, "https://youtube.com/"},
			},
			expected: "child-src 'self' https://youtube.com/;",
		},
		{
			desc: "connect",
			config: &csp.Config{
				Connect: []string{csp.Self, "https://example.com/"},
			},
			expected: "connect-src 'self' https://example.com/;",
		},
		{
			desc: "font",
			config: &csp.Config{
				Font: []string{csp.Self, "https://themes.googleusercontent.com"},
			},
			expected: "font-src 'self' https://themes.googleusercontent.com;",
		},
		{
			desc: "form",
			config: &csp.Config{
				Form: []string{csp.Self, "https://example.com/"},
			},
			expected: "form-action 'self' https://example.com/;",
		},
		{
			desc: "frame",
			config: &csp.Config{
				Frame: []string{csp.Self, "https://example.com/"},
			},
			expected: "frame-ancestors 'self' https://example.com/;",
		},
		{
			desc: "image",
			config: &csp.Config{
				Image: []string{csp.Self, "https://youtube.com/"},
			},
			expected: "img-src 'self' https://youtube.com/;",
		},
		{
			desc: "media",
			config: &csp.Config{
				Media: []string{csp.Self, "https://youtube.com/"},
			},
			expected: "media-src 'self' https://youtube.com/;",
		},
		{
			desc: "object",
			config: &csp.Config{
				Object: []string{csp.Self, "https://example.com/"},
			},
			expected: "object-src 'self' https://example.com/;",
		},
		{
			desc: "plugin",
			config: &csp.Config{
				Plugin: []string{"application/x-shockwave-flash", "application/x-java-applet"},
			},
			expected: "plugin-types application/x-shockwave-flash application/x-java-applet;",
		},
		{
			desc: "script",
			config: &csp.Config{
				Script: []string{csp.Self, "https://apis.google.com"},
			},
			expected: "script-src 'self' https://apis.google.com;",
		},
		{
			desc: "style",
			config: &csp.Config{
				Style: []string{csp.Self, "http://*.example.com", "'unsafe-inline'"},
			},
			expected: "style-src 'self' http://*.example.com 'unsafe-inline';",
		},
		{
			desc: "report",
			config: &csp.Config{
				Report: "/csp-report-endpoint",
			},
			expected: "report-uri /csp-report-endpoint;",
		},
		{
			desc: "upgrade insecure",
			config: &csp.Config{
				UpgradeInsecure: true,
			},
			expected: "upgrade-insecure-requests;",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			handler := getTestHandler()

			handlerFunc := CSP(handler, scenario.config)
			assert.IsType(t, (http.HandlerFunc)(nil), handlerFunc)

			resp := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/who-cares", nil)
			handlerFunc(resp, req)

			assert.True(t, handler.wasCalled())

			result := resp.Header().Get(cspHeader)
			assert.Equal(t, scenario.expected, result)
		})
	}
}

func TestCSP_extended(t *testing.T) {
	config := &csp.Config{
		Script:          []string{csp.Self, "https://apis.google.com"},
		UpgradeInsecure: true,
	}
	expected := "script-src 'self' https://apis.google.com; upgrade-insecure-requests;"

	handler := getTestHandler()

	handlerFunc := CSP(handler, config)
	assert.IsType(t, (http.HandlerFunc)(nil), handlerFunc)

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/who-cares", nil)
	handlerFunc(resp, req)

	assert.True(t, handler.wasCalled())

	result := resp.Header().Get(cspHeader)
	assert.Equal(t, expected, result)
}

func TestCSP_cdn(t *testing.T) {
	config := &csp.Config{
		Default: []string{"https://cdn.example.net"},
		Child:   []string{csp.None},
		Object:  []string{csp.None},
	}
	expected := "default-src https://cdn.example.net; child-src 'none'; object-src 'none';"

	handler := getTestHandler()

	handlerFunc := CSP(handler, config)
	assert.IsType(t, (http.HandlerFunc)(nil), handlerFunc)

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/who-cares", nil)
	handlerFunc(resp, req)

	assert.True(t, handler.wasCalled())

	result := resp.Header().Get(cspHeader)
	assert.Equal(t, expected, result)
}
