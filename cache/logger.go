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

// Logger allows for logging errors in the asynchronous calls
type Logger interface {
	// Build returns the data for the supplied key by populating dest
	Log(message string, args ...interface{})
}

// LoggerFunc implements Logger as a function
type LoggerFunc func(message string, args ...interface{})

// Log implements Logger
func (l LoggerFunc) Log(message string, args ...interface{}) {
	l(message, args...)
}

// No op implementation of Logger
var noopLogger = LoggerFunc(func(message string, args ...interface{}) {
	// intentionally do nothing
})
