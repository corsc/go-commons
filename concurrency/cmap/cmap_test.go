package cmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap_GetSet(t *testing.T) {
	scenarios := []struct {
		desc        string
		setup       func() *Map
		key         string
		expected    interface{}
		expectedErr error
	}{
		{
			desc: "no values",
			setup: func() *Map {
				return New()
			},
			key:         "foo",
			expected:    nil,
			expectedErr: ErrNoSuchItem,
		},
		{
			desc: "happy path",
			setup: func() *Map {
				myMap := New()
				_ = myMap.Set("foo", "bar")
				return myMap
			},
			key:         "foo",
			expected:    "bar",
			expectedErr: nil,
		},
		{
			desc: "unknown value",
			setup: func() *Map {
				myMap := New()
				_ = myMap.Set("foo", "bar")
				return myMap
			},
			key:         "apples",
			expected:    nil,
			expectedErr: ErrNoSuchItem,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			myMap := scenario.setup()

			result, resultErr := myMap.Get(scenario.key)
			assert.Equal(t, scenario.expected, result)
			assert.Equal(t, scenario.expectedErr, resultErr)
		})
	}

}

func TestMap_GetElseSet(t *testing.T) {
	scenarios := []struct {
		desc     string
		setup    func() *Map
		key      string
		newValue interface{}
		expected interface{}
	}{
		{
			desc: "no values",
			setup: func() *Map {
				return New()
			},
			key:      "foo",
			newValue: "bar",
			expected: "bar",
		},
		{
			desc: "known value",
			setup: func() *Map {
				myMap := New()
				_ = myMap.Set("foo", "bar")
				return myMap
			},
			key:      "foo",
			newValue: "apples",
			expected: "bar",
		},
		{
			desc: "unknown value",
			setup: func() *Map {
				myMap := New()
				_ = myMap.Set("foo", "bar")
				return myMap
			},
			key:      "apples",
			newValue: "are yummy",
			expected: "are yummy",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			myMap := scenario.setup()

			result, resultErr := myMap.GetElseSet(scenario.key, scenario.newValue)
			assert.Equal(t, scenario.expected, result)
			assert.Nil(t, resultErr)
		})
	}
}

func TestMap_Count(t *testing.T) {
	scenarios := []struct {
		desc     string
		setup    func() *Map
		expected int64
	}{
		{
			desc: "no values",
			setup: func() *Map {
				return New()
			},
			expected: 0,
		},
		{
			desc: "one value",
			setup: func() *Map {
				myMap := New()
				_ = myMap.Set("foo", "bar")
				return myMap
			},
			expected: 1,
		},
		{
			desc: "values",
			setup: func() *Map {
				myMap := New()
				_ = myMap.Set("foo", "bar")
				_ = myMap.Set("apples", "bar")
				_ = myMap.Set("oranges", "bar")
				return myMap
			},
			expected: 3,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			myMap := scenario.setup()

			result := myMap.Count()
			assert.Equal(t, scenario.expected, result)
		})
	}
}

func TestMap_Has(t *testing.T) {
	scenarios := []struct {
		desc     string
		setup    func() *Map
		key      string
		expected bool
	}{
		{
			desc: "no values",
			setup: func() *Map {
				return New()
			},
			key:      "foo",
			expected: false,
		},
		{
			desc: "exists",
			setup: func() *Map {
				myMap := New()
				_ = myMap.Set("foo", "bar")
				return myMap
			},
			key:      "foo",
			expected: true,
		},
		{
			desc: "values but doesn't exist",
			setup: func() *Map {
				myMap := New()
				_ = myMap.Set("foo", "bar")
				_ = myMap.Set("apples", "bar")
				_ = myMap.Set("oranges", "bar")
				return myMap
			},
			key:      "pairs",
			expected: false,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			myMap := scenario.setup()

			result := myMap.Has(scenario.key)
			assert.Equal(t, scenario.expected, result)
		})
	}
}

func TestMap_Remove(t *testing.T) {
	myMap := New()
	_ = myMap.Set("foo", "bar")
	_ = myMap.Set("apples", "bar")
	_ = myMap.Set("oranges", "bar")
	assert.Equal(t, int64(3), myMap.Count())

	myMap.Remove("apples")
	assert.Equal(t, int64(2), myMap.Count())
}

func TestMap_Iterator(t *testing.T) {
	expected := []Tuple{
		{Key: "foo", Value: "bar"},
		{Key: "apples", Value: "bar"},
		{Key: "oranges", Value: "bar"},
	}

	myMap := New()
	_ = myMap.Set("foo", "bar")
	_ = myMap.Set("apples", "bar")
	_ = myMap.Set("oranges", "bar")

	dataCh := myMap.Iterator()
	total := 0
	for tuple := range dataCh {
		assert.Contains(t, expected, tuple)
		total++
	}

	assert.Equal(t, 3, total)
}
