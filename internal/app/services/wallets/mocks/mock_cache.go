// Code generated by mockery v2.53.4. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockCache is an autogenerated mock type for the cache type
type MockCache struct {
	mock.Mock
}

// GetBalance provides a mock function with given fields: ctx, walletID
func (_m *MockCache) GetBalance(ctx context.Context, walletID string) (*int, error) {
	ret := _m.Called(ctx, walletID)

	if len(ret) == 0 {
		panic("no return value specified for GetBalance")
	}

	var r0 *int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*int, error)); ok {
		return rf(ctx, walletID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *int); ok {
		r0 = rf(ctx, walletID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*int)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, walletID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockCache creates a new instance of MockCache. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCache(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCache {
	mock := &MockCache{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
