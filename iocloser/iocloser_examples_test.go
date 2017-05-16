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
