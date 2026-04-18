# App Layer - Services, DTOs and Interfaces

## Location
`internal/{module}/app/`

## Required Files
- `dto.go` - Input (primitives) and Output (can have VOs)
- `port.go` - Interfaces for external dependencies
- `{name}_service.go` - Use case implementation
- `{name}_service_test.go` - Unit tests

## Critical Rules
- ✅ Context as **first parameter** always
- ✅ Wrap errors: `fmt.Errorf("s.repo.Method: %w", err)`
- ✅ Guard clauses (NO `else`)
- ✅ Validate input via Value Objects (domain)
- ❌ DO NOT import `ports` or `adapters` (interfaces only)
- ❌ DO NOT add logging (logging is done via HTTP middleware)

## Patterns

### dto.go
```go
type UserInput struct { Email string; Password string }      // primitives
type UserOutput struct { ID string; Email domain.Email }     // can have VO
```

### port.go
```go
type UserRepository interface {
    Create(ctx context.Context, user domain.User) error
    ExistsByEmail(ctx context.Context, email domain.Email) (bool, error)
}
```

### auth_service.go
```go
type authService struct { repo UserRepository }
func NewAuthService(repo UserRepository) *authService { }
func (s *authService) Register(ctx context.Context, input UserInput) (UserOutput, error) {
    email, err := domain.NewEmail(input.Email)
    if err != nil { return UserOutput{}, err }
    // business rules
    // persist
    // return DTO
}
```

## Tests (100% coverage when possible)

### Structure (Arrange-Act-Assert)
```go
func TestAuthService_Register(t *testing.T) {
    t.Run("should create user successfully", func(t *testing.T) {
        mock := mocks.NewMockUserRepository(t)
        mock.EXPECT().ExistsByEmail(mock.Anything, mock.Anything).Return(false, nil).Once()
        mock.EXPECT().Create(mock.Anything, mock.Anything).Return(nil).Once()
        svc := app.NewAuthService(mock)
        
        got, err := svc.Register(t.Context(), defaultInput)
        
        require.NoError(t, err)
        assert.NotEmpty(t, got.ID)
    })
}
```

### Naming
- Function: `TestStructName_MethodName`
- Subtest: `t.Run("should {behavior} when {condition}")`

See `internal/auth/app/` as reference.