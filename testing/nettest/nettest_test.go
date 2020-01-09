package nettest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTCP(t *testing.T) {
	var results []int

	for x := 0; x < 20; x++ {
		port := GetTCP()

		assert.NotEqual(t, 0, port, "port should not be 0")
		assert.NotContains(t, results, port, "port should not repeat")

		results = append(results, port)
	}
}

func TestGetUDP(t *testing.T) {
	var results []int

	for x := 0; x < 20; x++ {
		port := GetUDP()

		assert.NotEqual(t, 0, port, "port should not be 0")
		assert.NotContains(t, results, port, "port should not repeat")

		results = append(results, port)
	}
}
