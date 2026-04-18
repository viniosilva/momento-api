# App Layer - Services, DTOs and Interfaces

## Location
`internal/{module}/app/`

## Required Files
- `dto.go` - Input (primitives) and Output (can have VOs)
- `port.go` - Interfaces for external dependencies
- `{name}_service.go` - Use case implementation

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

## See Also
- @.cursor/rules/go_backend_tests.md

See `internal/auth/app/` as reference.