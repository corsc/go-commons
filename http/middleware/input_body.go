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
	"context"
	"encoding/json"
	"net/http"
	"reflect"
)

// key used to store the DTO in the context
type inputBodyKey int

var inputBodyDTO inputBodyKey

// LoggingClient allows for logging
type LoggingClient interface {
	Warn(msg string, args ...interface{})
}

// InputBody will attempt to populate a copy of the supplied struct and store it in request context.
func InputBody(handler http.Handler, dto interface{}, client ...LoggingClient) http.HandlerFunc {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		dtoType := reflect.TypeOf(dto)
		dtoCopy := reflect.New(dtoType.Elem()).Interface()

		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(dtoCopy)
		if err != nil {
			logWarn(client, "error during JSON decode of request body. err: %s", err)
			http.Error(resp, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// store the DTO into the context
		req = req.WithContext(context.WithValue(req.Context(), inputBodyDTO, dtoCopy))

		handler.ServeHTTP(resp, req)
	})
}

// InputBodyDTO returns the populated DTO for this request (or nil).
// NOTE: this method should be used in conjunction with `InputBody()`
func InputBodyDTO(req *http.Request) interface{} {
	return req.Context().Value(inputBodyDTO)
}

// convenience method for logging when the logging client was supplied (via optional/variadic arg)
func logWarn(client []LoggingClient, msg string, args ...interface{}) {
	if len(client) == 0 {
		return
	}

	client[0].Warn(msg, args...)
}
