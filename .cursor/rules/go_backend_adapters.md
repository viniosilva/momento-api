# Adapters Layer - Repositories and Integrations

## Location
`internal/{module}/adapters/`

## Common Files
- `{name}_repository.go` - Implements `app/port.go` (DB access)
- `{name}_service.go` - External services (JWT, APIs, etc)
- `{name}_service_test.go` - Tests for external services
- `mongo_indexes.go` - Index orchestrator
- `indexes/00N_{name}.go` - Specific indexes

## Critical Rules
- ✅ Implement interfaces from `app/port.go`
- ✅ Translate adapter errors to domain errors
- ✅ Context as first parameter
- ✅ Test external services (JWT, APIs) with unit tests
- ❌ DO NOT test repositories (leave for integration tests)
- ❌ DO NOT expose implementation details (MongoDB, etc)
- ❌ DO NOT add logging (logging is done via HTTP middleware)

## Patterns

### user_repository.go
```go
type userRepository struct { collection *mongo.Collection }
func NewUserRepository(col *mongo.Collection) *userRepository { }
func (r *userRepository) Create(ctx context.Context, user domain.User) error {
    _, err := r.collection.InsertOne(ctx, user)
    if mongo.IsDuplicateKeyError(err) {
        return domain.ErrUserAlreadyExists  // translate error
    }
    return err
}
```

### mongo_indexes.go
```go
func CreateIndexes(ctx context.Context, db *mongo.Database) error {
    return indexes.CreateUserEmailIndex(ctx, db)
}
```

### indexes/00N_users_unique_email.go
```go
func CreateUserEmailIndex(ctx context.Context, db *mongo.Database) error {
    // mongo.IndexModel with unique
}
```

See `internal/auth/adapters/` as reference.