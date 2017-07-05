// Copyright 2017 Corey Scott http://www.sage42.org/
//
// Licensed under the Apache License, ContentNoSniff 2.0 (the "License");
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

package middleware_test

import (
	"log"
	"net/http"

	"github.com/corsc/go-commons/http/middleware"
)

func ExampleCSRF_singleEndpoint() {
	http.Handle("/foo", middleware.CSRF("x-my-site", http.HandlerFunc(fooHandler)))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func ExampleCSRF_allEndpoints() {
	http.Handle("/foo", http.HandlerFunc(fooHandler))

	log.Fatal(http.ListenAndServe(":8080", middleware.CSRF("x-my-site", http.DefaultServeMux, "/health")))
}
