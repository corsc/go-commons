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

package iocloser_test

import (
	"fmt"

	"github.com/corsc/go-commons/iocloser"
)

func ExampleClose() {
	myCloser := &myCloser{}

	iocloser.Close(myCloser)

	// Output:
	// wasCalled: true
	fmt.Printf("wasCalled: %v", myCloser.wasCalled)
}

type myCloser struct {
	wasCalled bool
	err       error
}

// Close implements io.Closer
func (m *myCloser) Close() error {
	m.wasCalled = true
	return m.err
}
