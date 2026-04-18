# Domain Layer - Entities and Business Rules

## Location
`internal/{module}/domain/`

## Responsibilities
- Entities and Value Objects
- Pure business rules
- Domain errors
- Constants (e.g., collection names)

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
const UsersCollectionName = "users"
type User struct {
    ID        primitive.ObjectID `bson:"_id"`
    Email     Email              `bson:"email"`
    CreatedAt time.Time          `bson:"created_at"`
}
func NewUser(email Email) User { /* constructor */ }
```

See `internal/auth/domain/` as reference.