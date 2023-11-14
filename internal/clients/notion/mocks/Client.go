// Code generated by mockery v2.35.4. DO NOT EDIT.

package mocks

import (
	context "context"

	notion "github.com/nickysemenza/gourd/internal/clients/notion"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// GetAll provides a mock function with given fields: ctx, lookback, pageID
func (_m *Client) GetAll(ctx context.Context, lookback time.Duration, pageID string) ([]notion.Recipe, error) {
	ret := _m.Called(ctx, lookback, pageID)

	var r0 []notion.Recipe
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, time.Duration, string) ([]notion.Recipe, error)); ok {
		return rf(ctx, lookback, pageID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, time.Duration, string) []notion.Recipe); ok {
		r0 = rf(ctx, lookback, pageID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]notion.Recipe)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, time.Duration, string) error); ok {
		r1 = rf(ctx, lookback, pageID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewClient creates a new instance of Client. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *Client {
	mock := &Client{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
