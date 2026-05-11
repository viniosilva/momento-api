package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"momento/internal/auth/domain"
)

func NewMockUserRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUserRepository {
	mock := &MockUserRepository{}
	mock.Mock.Test(t)
	t.Cleanup(func() { mock.AssertExpectations(t) })
	return mock
}

type MockUserRepository struct {
	mock.Mock
}

func (_m *MockUserRepository) Create(ctx context.Context, user domain.User) error {
	ret := _m.Called(ctx, user)
	return ret.Error(0)
}

func (_m *MockUserRepository) ExistsByEmail(ctx context.Context, email domain.Email) (bool, error) {
	ret := _m.Called(ctx, email)
	return ret.Bool(0), ret.Error(1)
}

func (_m *MockUserRepository) FindByEmail(ctx context.Context, email domain.Email) (domain.User, error) {
	ret := _m.Called(ctx, email)
	return ret.Get(0).(domain.User), ret.Error(1)
}

func (_m *MockUserRepository) FindVerifiedByEmail(ctx context.Context, email domain.Email) (domain.User, error) {
	ret := _m.Called(ctx, email)
	return ret.Get(0).(domain.User), ret.Error(1)
}

func (_m *MockUserRepository) FindByID(ctx context.Context, id string) (domain.User, error) {
	ret := _m.Called(ctx, id)
	return ret.Get(0).(domain.User), ret.Error(1)
}

func (_m *MockUserRepository) Update(ctx context.Context, user domain.User) error {
	ret := _m.Called(ctx, user)
	return ret.Error(0)
}

func (_m *MockUserRepository) EXPECT() *MockUserRepository_Expecter {
	return &MockUserRepository_Expecter{mock: &_m.Mock}
}

type MockUserRepository_Expecter struct {
	mock *mock.Mock
}

type MockUserRepositoryCreateCall struct {
	*mock.Call
}

func (_e *MockUserRepository_Expecter) Create(ctx interface{}, user interface{}) *MockUserRepositoryCreateCall {
	return &MockUserRepositoryCreateCall{Call: _e.mock.On("Create", ctx, user)}
}

func (_c *MockUserRepositoryCreateCall) Return(_a0 error) *MockUserRepositoryCreateCall {
	_c.Call.Return(_a0)
	return _c
}

type MockUserRepositoryExistsByEmailCall struct {
	*mock.Call
}

func (_e *MockUserRepository_Expecter) ExistsByEmail(ctx interface{}, email interface{}) *MockUserRepositoryExistsByEmailCall {
	return &MockUserRepositoryExistsByEmailCall{Call: _e.mock.On("ExistsByEmail", ctx, email)}
}

func (_c *MockUserRepositoryExistsByEmailCall) Return(_a0 bool, _a1 error) *MockUserRepositoryExistsByEmailCall {
	_c.Call.Return(_a0, _a1)
	return _c
}

type MockUserRepositoryFindByEmailCall struct {
	*mock.Call
}

func (_e *MockUserRepository_Expecter) FindByEmail(ctx interface{}, email interface{}) *MockUserRepositoryFindByEmailCall {
	return &MockUserRepositoryFindByEmailCall{Call: _e.mock.On("FindByEmail", ctx, email)}
}

func (_c *MockUserRepositoryFindByEmailCall) Return(_a0 domain.User, _a1 error) *MockUserRepositoryFindByEmailCall {
	_c.Call.Return(_a0, _a1)
	return _c
}

type MockUserRepositoryFindVerifiedByEmailCall struct {
	*mock.Call
}

func (_e *MockUserRepository_Expecter) FindVerifiedByEmail(ctx interface{}, email interface{}) *MockUserRepositoryFindVerifiedByEmailCall {
	return &MockUserRepositoryFindVerifiedByEmailCall{Call: _e.mock.On("FindVerifiedByEmail", ctx, email)}
}

func (_c *MockUserRepositoryFindVerifiedByEmailCall) Return(_a0 domain.User, _a1 error) *MockUserRepositoryFindVerifiedByEmailCall {
	_c.Call.Return(_a0, _a1)
	return _c
}

type MockUserRepositoryFindByIDCall struct {
	*mock.Call
}

func (_e *MockUserRepository_Expecter) FindByID(ctx interface{}, id interface{}) *MockUserRepositoryFindByIDCall {
	return &MockUserRepositoryFindByIDCall{Call: _e.mock.On("FindByID", ctx, id)}
}

func (_c *MockUserRepositoryFindByIDCall) Return(_a0 domain.User, _a1 error) *MockUserRepositoryFindByIDCall {
	_c.Call.Return(_a0, _a1)
	return _c
}

type MockUserRepositoryUpdateCall struct {
	*mock.Call
}

func (_e *MockUserRepository_Expecter) Update(ctx interface{}, user interface{}) *MockUserRepositoryUpdateCall {
	return &MockUserRepositoryUpdateCall{Call: _e.mock.On("Update", ctx, user)}
}

func (_c *MockUserRepositoryUpdateCall) Return(_a0 error) *MockUserRepositoryUpdateCall {
	_c.Call.Return(_a0)
	return _c
}
