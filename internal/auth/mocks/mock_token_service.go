package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

func NewMockTokenService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTokenService {
	mock := &MockTokenService{}
	mock.Mock.Test(t)
	t.Cleanup(func() { mock.AssertExpectations(t) })
	return mock
}

type MockTokenService struct {
	mock.Mock
}

func (_m *MockTokenService) Store(ctx context.Context, token string, userID string, ttl time.Duration) error {
	ret := _m.Called(ctx, token, userID, ttl)
	return ret.Error(0)
}

func (_m *MockTokenService) Validate(ctx context.Context, token string) (string, error) {
	ret := _m.Called(ctx, token)
	return ret.Get(0).(string), ret.Error(1)
}

func (_m *MockTokenService) Invalidate(ctx context.Context, token string) error {
	ret := _m.Called(ctx, token)
	return ret.Error(0)
}

func (_m *MockTokenService) EXPECT() *MockTokenService_Expecter {
	return &MockTokenService_Expecter{mock: &_m.Mock}
}

type MockTokenService_Expecter struct {
	mock *mock.Mock
}

type MockTokenServiceStoreCall struct {
	*mock.Call
}

func (_e *MockTokenService_Expecter) Store(ctx interface{}, token interface{}, userID interface{}, ttl interface{}) *MockTokenServiceStoreCall {
	return &MockTokenServiceStoreCall{Call: _e.mock.On("Store", ctx, token, userID, ttl)}
}

func (_c *MockTokenServiceStoreCall) Return(_a0 error) *MockTokenServiceStoreCall {
	_c.Call.Return(_a0)
	return _c
}

type MockTokenServiceValidateCall struct {
	*mock.Call
}

func (_e *MockTokenService_Expecter) Validate(ctx interface{}, token interface{}) *MockTokenServiceValidateCall {
	return &MockTokenServiceValidateCall{Call: _e.mock.On("Validate", ctx, token)}
}

func (_c *MockTokenServiceValidateCall) Return(_a0 string, _a1 error) *MockTokenServiceValidateCall {
	_c.Call.Return(_a0, _a1)
	return _c
}

type MockTokenServiceInvalidateCall struct {
	*mock.Call
}

func (_e *MockTokenService_Expecter) Invalidate(ctx interface{}, token interface{}) *MockTokenServiceInvalidateCall {
	return &MockTokenServiceInvalidateCall{Call: _e.mock.On("Invalidate", ctx, token)}
}

func (_c *MockTokenServiceInvalidateCall) Return(_a0 error) *MockTokenServiceInvalidateCall {
	_c.Call.Return(_a0)
	return _c
}
