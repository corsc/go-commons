package skip

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIfNotSet(t *testing.T) {
	// test settings
	keyThatExists := fmt.Sprintf("exists-%d", time.Now().UnixNano())
	keyMissing := fmt.Sprintf("missing-%d", time.Now().UnixNano())
	os.Setenv(keyThatExists, "foo")

	// test for key that exists
	skipper := &mockSkipper{}
	IfNotSet(skipper, keyThatExists)
	assert.False(t, skipper.wasCalled)

	// test for key that does not exists
	skipper = &mockSkipper{}
	IfNotSet(skipper, keyMissing)
	assert.True(t, skipper.wasCalled)
}

type mockSkipper struct {
	wasCalled bool
}

// Skip implements Skipper
func (m *mockSkipper) Skip(args ...interface{}) {
	m.wasCalled = true
}
