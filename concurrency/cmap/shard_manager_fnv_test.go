package cmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShardManagerFNV_implements(t *testing.T) {
	assert.Implements(t, (*ShardManager)(nil), &ShardManagerFNV{})
}

func TestShardManagerFNV_GetTotalShards(t *testing.T) {
	scenarios := []struct {
		desc     string
		setup    func() *ShardManagerFNV
		expected int64
	}{
		{
			desc: "default",
			setup: func() *ShardManagerFNV {
				return &ShardManagerFNV{}
			},
			expected: 32,
		},
		{
			desc: "custom value",
			setup: func() *ShardManagerFNV {
				return &ShardManagerFNV{
					TotalShards: 666,
				}
			},
			expected: 666,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			manager := scenario.setup()
			result := manager.GetTotalShards()

			assert.Equal(t, scenario.expected, result)
		})
	}
}

func TestShardManagerFNV_GetShardNo(t *testing.T) {
	scenarios := []struct {
		key      string
		expected int64
	}{
		{
			key:      "foo",
			expected: 19,
		},
		{
			key:      "bar",
			expected: 0,
		},
		{
			key:      "apples",
			expected: 6,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.key, func(t *testing.T) {
			manager := &ShardManagerFNV{}
			result, resultErr := manager.GetShardNo(scenario.key)

			assert.Equal(t, scenario.expected, result)
			assert.Nil(t, resultErr)
		})
	}
}
