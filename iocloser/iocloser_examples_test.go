package iocloser_test

import (
	"fmt"

	"errors"

	"github.com/corsc/go-commons/iocloser"
)

func ExampleClose() {
	myCloser := &myCloser{}

	iocloser.Close(myCloser)

	// Output:
	// wasCalled: true
	fmt.Printf("wasCalled: %v", myCloser.wasCalled)
}

func ExampleClose_withLogger() {
	myCloser := &myCloser{
		err: errors.New("something failed"),
	}
	logger := &mockLogger{}

	iocloser.Close(myCloser, logger.log)

	// Output:
	// closer.wasCalled: true
	// logger.wasCalled: true
	fmt.Printf("closer.wasCalled: %v\n", myCloser.wasCalled)
	fmt.Printf("logger.wasCalled: %v", logger.wasCalled)
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

type mockLogger struct {
	wasCalled bool
}

func (m *mockLogger) log(_ string, _ ...interface{}) {
	m.wasCalled = true
}
