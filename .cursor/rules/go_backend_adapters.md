# Adapters Layer - Repositories and Integrations

## Location
`internal/{module}/adapters/`

## Common Files
- `{name}_repository.go` - Implements `app/port.go` (DB access)
- `{name}_model.go` - Database document model + converters
- `{name}_service.go` - External services (JWT, APIs, etc)
- `mongo_indexes.go` - Index orchestrator
- `indexes/00N_{name}.go` - Specific indexes

## Critical Rules
- ✅ Implement interfaces from `app/port.go`
- ✅ Translate adapter errors to domain errors
- ✅ Context as first parameter
- ✅ Separate database model (`{name}_model.go`) from domain entities
- ✅ Use converter functions (`toUserDocument`, `toUserDomain`) to map between layers
- ❌ DO NOT expose implementation details (MongoDB types in domain)
- ❌ DO NOT add logging (logging is done via HTTP middleware)

## See Also
- @.cursor/rules/go_backend_tests.md

## Patterns

### user_repository.go
```go
type userRepository struct { collection *mongo.Collection }
func NewUserRepository(col *mongo.Collection) *userRepository { }
func (r *userRepository) Create(ctx context.Context, user domain.User) error {
    doc, err := toUserDocument(user)
    if err != nil { return err }

    _, err = r.collection.InsertOne(ctx, doc)
    if mongo.IsDuplicateKeyError(err) {
        return domain.ErrUserAlreadyExists  // translate error
    }
    return err
}
```

### user_model.go (Database Model)
```go
type userDocument struct {
    ID        primitive.ObjectID `bson:"_id"`
    Email     string             `bson:"email"`
    Password  string             `bson:"password"`
    CreatedAt time.Time          `bson:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at"`
}

func toUserDocument(u domain.User) (userDocument, error) {
    id, err := primitive.ObjectIDFromHex(u.ID)
    if err != nil { return userDocument{}, err }

    return userDocument{
        ID:        id,
        Email:     string(u.Email),
        Password:  string(u.Password),
        CreatedAt: u.CreatedAt,
        UpdatedAt: u.UpdatedAt,
    }, nil
}

func toUserDomain(d userDocument) domain.User {
    return domain.User{
        ID:        d.ID.Hex(),
        Email:     domain.Email(d.Email),
        Password:  domain.Password(d.Password),
        CreatedAt: d.CreatedAt,
        UpdatedAt: d.UpdatedAt,
    }
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