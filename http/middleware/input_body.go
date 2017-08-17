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
	"io/ioutil"
	"net/http"
	"reflect"
)

// key used to store the DTO in the context
type inputBodyKey int

var inputBodyDTO inputBodyKey

// InputBody will attempt to populate a copy of the supplied struct and store it in request context.
func InputBody(handler http.Handler, dto interface{}, client ...InputBodyLogger) http.HandlerFunc {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		dtoType := reflect.TypeOf(dto)
		dtoCopy := reflect.New(dtoType.Elem()).Interface()

		contents, err := ioutil.ReadAll(req.Body)
		if err != nil {
			logWarn(req, client, "error reading request body. err: %s", err)
			http.Error(resp, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		logBody(req, client, contents)

		err = json.Unmarshal(contents, dtoCopy)
		if err != nil {
			logWarn(req, client, "error during JSON decode of request body. err: %s", err)
			http.Error(resp, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// store the DTO into the context
		req = InputBodySetDTO(req, dtoCopy)

		handler.ServeHTTP(resp, req)
	})
}

// InputBodyDTO returns the populated DTO for this request (or nil).
// NOTE: this method should be used in conjunction with `InputBody()`
func InputBodyDTO(req *http.Request) interface{} {
	return req.Context().Value(inputBodyDTO)
}

// InputBodySetDTO sets the supplied DTO into the context.
// This method will generally be used only during testing.
func InputBodySetDTO(req *http.Request, dto interface{}) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), inputBodyDTO, dto))
}

// convenience method for logging when the logging client was supplied (via optional/variadic arg)
func logWarn(req *http.Request, client []InputBodyLogger, msg string, args ...interface{}) {
	if len(client) == 0 {
		return
	}

	client[0].BadRequest(req, msg, args...)
}

func logBody(req *http.Request, client []InputBodyLogger, body []byte) {
	if len(client) == 0 {
		return
	}

	client[0].Request(req, body)
}

// InputBodyLogger allows for logging
type InputBodyLogger interface {
	// Log bad requests
	BadRequest(req *http.Request, msg string, args ...interface{})

	// Log request bodies
	Request(req *http.Request, body []byte)
}
