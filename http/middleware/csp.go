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

	"github.com/corsc/go-commons/http/middleware/csp"
)

const (
	cspHeader = "Content-Security-Policy"

	cspPrefixDefault         = "default-src"
	cspPrefixBase            = "base-uri"
	cspPrefixChild           = "child-src"
	cspPrefixConnect         = "connect-src"
	cspPrefixFont            = "font-src"
	cspPrefixForm            = "form-action"
	cspPrefixFrame           = "frame-ancestors"
	cspPrefixImage           = "img-src"
	cspPrefixMedia           = "media-src"
	cspPrefixObject          = "object-src"
	cspPrefixPlugin          = "plugin-types"
	cspPrefixReport          = "report-uri"
	cspPrefixScript          = "script-src"
	cspPrefixStyle           = "style-src"
	cspPrefixUpgradeInsecure = "upgrade-insecure-requests"

	cspSeparator = ";"
)

// CSP sends Content Security Policy header.
//
// Reference:
// https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP
// https://www.html5rocks.com/en/tutorials/security/content-security-policy/
func CSP(handler http.Handler, config *csp.Config) http.HandlerFunc {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		value := ""
		cspAdd(&value, cspPrefixDefault, config.Default)
		cspAdd(&value, cspPrefixBase, config.Base)
		cspAdd(&value, cspPrefixChild, config.Child)
		cspAdd(&value, cspPrefixConnect, config.Connect)
		cspAdd(&value, cspPrefixFont, config.Font)
		cspAdd(&value, cspPrefixForm, config.Form)
		cspAdd(&value, cspPrefixFrame, config.Frame)
		cspAdd(&value, cspPrefixImage, config.Image)
		cspAdd(&value, cspPrefixMedia, config.Media)
		cspAdd(&value, cspPrefixObject, config.Object)
		cspAdd(&value, cspPrefixPlugin, config.Plugin)
		cspAdd(&value, cspPrefixScript, config.Script)
		cspAdd(&value, cspPrefixStyle, config.Style)

		cspAddString(&value, cspPrefixReport, config.Report)
		cspAddBool(&value, cspPrefixUpgradeInsecure, config.UpgradeInsecure)

		resp.Header().Set(cspHeader, value)

		handler.ServeHTTP(resp, req)
	})
}

func cspAdd(dest *string, prefix string, locations []string) {
	if len(locations) == 0 {
		return
	}

	if len(*dest) > 0 {
		*dest += " "
	}

	*dest += prefix

	for _, location := range locations {
		*dest += " " + location
	}

	*dest += cspSeparator
}

func cspAddString(dest *string, prefix string, value string) {
	if len(value) == 0 {
		return
	}

	if len(*dest) > 0 {
		*dest += " "
	}

	*dest += prefix + " " + value + cspSeparator
}

func cspAddBool(dest *string, output string, value bool) {
	if !value {
		return
	}

	if len(*dest) > 0 {
		*dest += " "
	}

	*dest += output + cspSeparator
}
