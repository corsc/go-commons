package cache

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockStorage is an autogenerated mock type for the Storage type
type MockStorage struct {
	mock.Mock
}

// Get provides a mock function with given fields: ctx, key
func (_m *MockStorage) Get(ctx context.Context, key string) ([]byte, error) {
	ret := _m.Called(ctx, key)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, string) []byte); ok {
		r0 = rf(ctx, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Invalidate provides a mock function with given fields: ctx, key
func (_m *MockStorage) Invalidate(ctx context.Context, key string) error {
	ret := _m.Called(ctx, key)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Set provides a mock function with given fields: ctx, key, bytes
func (_m *MockStorage) Set(ctx context.Context, key string, bytes []byte) error {
	ret := _m.Called(ctx, key, bytes)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []byte) error); ok {
		r0 = rf(ctx, key, bytes)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}