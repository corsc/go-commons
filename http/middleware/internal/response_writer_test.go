package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCustomResponseWriter(t *testing.T) {
	writer := httptest.NewRecorder()
	myWriter := NewCustomResponseWriter(writer)
	assert.Implements(t, (*http.ResponseWriter)(nil), myWriter)
}

func TestCustomResponseWriter_WriteHeader(t *testing.T) {
	writer := httptest.NewRecorder()
	myWriter := NewCustomResponseWriter(writer)

	myWriter.WriteHeader(http.StatusBadRequest)
	assert.Equal(t, http.StatusBadRequest, myWriter.Status())
	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestCustomResponseWriter_Write(t *testing.T) {
	writer := httptest.NewRecorder()
	myWriter := NewCustomResponseWriter(writer)

	myWriter.Write([]byte("foo"))
	assert.Equal(t, http.StatusOK, myWriter.Status())
}
