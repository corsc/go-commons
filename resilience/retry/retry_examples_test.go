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

package retry_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/corsc/go-commons/resilience/retry"
)

func ExampleRetry_Do() {
	// simplest usage; using the defaults
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	err := (&retry.Client{}).Do(ctx, "myretry", func() error {
		// do something amazing here
		return nil
	})

	// Output:
	// error was <nil>
	fmt.Printf("error was %v", err)
}

func ExampleRetry_Do_customErrorHandling() {
	// some custom errors
	var ErrUserError = errors.New("bad user input - retrying won't help")
	var ErrNetworkError = errors.New("bad connection - retrying might help")

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	client := &retry.Client{
		CanRetry: func(err error) bool {
			if err == ErrUserError {
				return false
			}
			if err == ErrNetworkError {
				return true
			}
			return true
		},
	}

	err := client.Do(ctx, "myretry", func() error {
		resp, err := http.Get("http://www.google.com/")
		if err != nil {
			// pass HTTP client errors out
			return err
		}

		// convert response codes into errors (where appropriate)
		switch resp.StatusCode {
		case http.StatusOK:
			return nil

		case http.StatusBadRequest:
			return ErrUserError

		case http.StatusBadGateway:
			return ErrNetworkError
		}

		return nil
	})

	// Output:
	// error was <nil>
	fmt.Printf("error was %v", err)
}
