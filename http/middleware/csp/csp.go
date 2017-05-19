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

package csp

const (
	// None matches nothing (i.e. disable this media type)
	None = "'none'"

	// Self matches the current origin, but not its subdomains.
	Self = "'self'"
)

// Config is the config for the CSP header
// Note: all config is optional
// Note: `None` and `Self` can be added to most of the settings in this struct
type Config struct {
	// Default serves as a fallback for the other CSP fetch directives
	Default []string

	// Base restricts the URLs that can appear in a page’s <base> element
	Base []string

	// Child lists the URLs for workers and embedded frame contents.
	Child []string

	// Connect limits the origins to which you can connect (via XHR, WebSockets, and EventSource).
	Connect []string

	// Font specifies the origins that can serve web fonts.
	Font []string

	// Form lists valid endpoints for submission from `<form>` tags
	Form []string

	// Frame specifies the sources that can embed the current page.
	Frame []string

	// Image defines the origins from which images can be loaded.
	Image []string

	// Media restricts the origins allowed to deliver video and audio.
	Media []string

	// Object allows control over Flash and other plugins.
	Object []string

	// Plugin limits the kinds of plugins a page may invoke.
	Plugin []string

	// Script defines the origins from which scripts can be loaded.
	Script []string

	// Style defines the origins from which stylesheets can be loaded.
	Style []string

	// Report specifies a URL where a browser will send reports when a content security policy is violated.
	Report string

	// UpgradeInsecure instructs user agents to rewrite URL schemes, changing HTTP to HTTPS.
	UpgradeInsecure bool
}
