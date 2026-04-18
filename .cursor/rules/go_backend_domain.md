# Domain Layer - Entities and Business Rules

## Location
`internal/{module}/domain/`

## Responsibilities
- Entities and Value Objects
- Pure business rules
- Domain errors
- ID generation (via `pkg/uid`)

## Critical Rules
- ❌ DO NOT import: `app`, `adapters`, `ports`
- ❌ DO NOT access: database, external APIs
- ✅ Only pure business logic
- ✅ Value Objects validate themselves (Fail-Fast)
- ✅ Business naming (ScheduleTraining, CancelTraining, not generic Create/Delete)
- ✅ Value Objects sanitize input (e.g., XSS prevention via bluemonday)

## Patterns

### Domain errors (global variables)
```go
var (
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidEmail = errors.New("invalid email format")
)
```

### Value Object with validation
```go
type Email string
func NewEmail(value string) (Email, error) { /* validate and return */ }
func ValidateEmail(value string) error { /* validation logic */ }
```

### Entity
```go
type User struct {
    ID        string    // string (not MongoDB ObjectID)
    Email     Email
    Password  Password
    CreatedAt time.Time
    UpdatedAt time.Time
}
func NewUser(email Email, password Password) User {
    return User{
        ID:        uid.New(),  // use pkg/uid for ID generation
        Email:     email,
        Password:  password,
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
    }
}
```

**Important**: Domain entities must NOT know about database implementation:
- Use primitive types (string, int) - NOT `primitive.ObjectID`
- NO database tags (`bson`, `json` from DB)
- Use UUID or string IDs from `pkg/uid`

**See Also**
- @.cursor/rules/go_backend_adapters.md (for database model with converters)

See `internal/auth/domain/` as reference.