package iocloser

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClose(t *testing.T) {
	myCloser := &myCloser{}

	Close(myCloser)

	assert.True(t, myCloser.wasCalled)
}

func TestClose_withLogger(t *testing.T) {
	myCloser := &myCloser{
		err: errors.New("something failed"),
	}
	logger := &mockLogger{}

	Close(myCloser, logger.log)

	assert.True(t, myCloser.wasCalled)
	assert.True(t, logger.wasCalled)
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
