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

package cache

// Metrics allows for instrumenting the cache for hit/miss and errors
//go:generate mockery -name Metrics -inpkg -testonly -case underscore
type Metrics interface {
	// Tracks
	Track(event Event)
}

// MetricsFunc implements Metrics as a function
type MetricsFunc func(event Event)

// Track implements Metrics
func (m MetricsFunc) Track(event Event) {
	m(event)
}

// No op implementation of metrics
var noopMetrics = MetricsFunc(func(event Event) {
	// intentionally do nothing
})
