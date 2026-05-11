package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

func NewMockEmailSender(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockEmailSender {
	mock := &MockEmailSender{}
	mock.Mock.Test(t)
	t.Cleanup(func() { mock.AssertExpectations(t) })
	return mock
}

type MockEmailSender struct {
	mock.Mock
}

func (_m *MockEmailSender) SendResetPasswordEmail(ctx context.Context, to string, token string) error {
	ret := _m.Called(ctx, to, token)
	return ret.Error(0)
}

func (_m *MockEmailSender) SendVerificationEmail(ctx context.Context, to string, token string) error {
	ret := _m.Called(ctx, to, token)
	return ret.Error(0)
}

func (_m *MockEmailSender) EXPECT() *MockEmailSender_Expecter {
	return &MockEmailSender_Expecter{mock: &_m.Mock}
}

type MockEmailSender_Expecter struct {
	mock *mock.Mock
}

type MockEmailSenderSendResetPasswordEmailCall struct {
	*mock.Call
}

func (_e *MockEmailSender_Expecter) SendResetPasswordEmail(ctx interface{}, to interface{}, token interface{}) *MockEmailSenderSendResetPasswordEmailCall {
	return &MockEmailSenderSendResetPasswordEmailCall{Call: _e.mock.On("SendResetPasswordEmail", ctx, to, token)}
}

func (_c *MockEmailSenderSendResetPasswordEmailCall) Return(_a0 error) *MockEmailSenderSendResetPasswordEmailCall {
	_c.Call.Return(_a0)
	return _c
}

type MockEmailSenderSendVerificationEmailCall struct {
	*mock.Call
}

func (_e *MockEmailSender_Expecter) SendVerificationEmail(ctx interface{}, to interface{}, token interface{}) *MockEmailSenderSendVerificationEmailCall {
	return &MockEmailSenderSendVerificationEmailCall{Call: _e.mock.On("SendVerificationEmail", ctx, to, token)}
}

func (_c *MockEmailSenderSendVerificationEmailCall) Return(_a0 error) *MockEmailSenderSendVerificationEmailCall {
	_c.Call.Return(_a0)
	return _c
}
