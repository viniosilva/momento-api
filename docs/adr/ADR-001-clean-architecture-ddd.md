# ADR-001: Clean Architecture with Domain-Driven Design

**Status**: ✅ Accepted

**Date**: 2026-03-28

## Context

The project needed a scalable and testable structure that enables:
- Clear separation of concerns
- Easy isolated testing (unit tests, integration tests)
- Independence from frameworks and libraries
- Maintainability and efficient onboarding

A Go monolith with standard `net/http`, no web framework or ORM.

## Decision

Implement **Clean Architecture** with **4 layers** and **Domain-Driven Design** patterns:

```
┌──────────────────────────────────────────────┐
│    Presentation (HTTP Handlers & Routes)     │
│      /internal/{domain}/presentation         │
├──────────────────────────────────────────────┤
│   Application (Business Logic & Services)    │
│      /internal/{domain}/application          │
├──────────────────────────────────────────────┤
│  Domain (Entities, Value Objects, Rules)     │
│         /internal/{domain}/domain            │
├──────────────────────────────────────────────┤
│ Infrastructure (DB, External Services, HTTP) │
│      /internal/{domain}/infrastructure       │
└──────────────────────────────────────────────┘
```

### Dependency Rules

```
📌 Allowed flow: Presentation → Application → Domain ← Infrastructure
❌ Forbidden: Domain → Application or Domain → Presentation
```

- **Domain**: Zero dependencies on other layers (only std lib)
- **Application**: Knows Domain, defines Ports (interfaces)
- **Infrastructure**: Implements Ports, provides external integrations
- **Presentation**: Orchestrates Application/Infrastructure, exposes HTTP API

### Value Objects Pattern

Implement validation at domain layer using **Value Objects**:

```go
// Domain: Guaranteed NoteContent is ALWAYS valid
type NoteContent string

func NewNoteContent(value string) (NoteContent, error) {
    normalized := strings.TrimSpace(value)
    
    if normalized == "" {
        return "", ErrContentEmpty
    }
    if len(normalized) > 100000 {
        return "", ErrContentTooLong
    }
    
    // Sanitization happens here (see ADR-002)
    sanitized := bluemonday.StrictPolicy().Sanitize(normalized)
    return NoteContent(sanitized), nil
}

// Usage: Guaranteed valid content
content, err := NewNoteContent(userInput)
if err != nil {
    return err // Validation error
}
// If here: content is GUARANTEED valid
note := domain.NewNote(id, content) // No need to re-validate
```

### Repository Pattern

Define interfaces in Application layer, implement in Infrastructure:

```go
// Application Layer: Define contract
package application

type NoteRepository interface {
    Save(ctx context.Context, note domain.Note) error
    FindByID(ctx context.Context, id string) (domain.Note, error)
}

// Infrastructure Layer: Implement with MongoDB
package infrastructure

type mongoNoteRepository struct {
    client *mongo.Client
}

func (r *mongoNoteRepository) Save(ctx context.Context, note domain.Note) error {
    collection := r.client.Database("pinnado").Collection("notes")
    _, err := collection.InsertOne(ctx, bson.M{
        "_id":     note.ID,
        "content": string(note.Content),
    })
    return err
}
```

## Consequences

### Positive ✅

- **Testability**: Mock dependencies via Ports (interfaces)
- **Independence**: Domain layer free from external dependencies
- **Maintainability**: Database changes don't affect business logic
- **Scalability**: Easy to add new domains/features
- **Reusability**: Application logic works for CLI, gRPC, WebSocket, etc
- **Type Safety**: Value Objects guarantee valid data at compile time

### Negative ❌

- **Overhead**: More layers = more files and boilerplate
- **Learning curve**: Pattern requires architectural understanding
- **Latency**: Each handler traverses all 4 layers

## Alternatives Considered

### 1. Simple 2-3 Layer Architecture

Handler → Service → Repository

- ❌ Domain logic mixed with business logic
- ❌ Hard to test in isolation
- ✅ Simpler initially

### 2. Anemic Domain Model

Services handle everything, entities are just data containers

- ❌ Business rules scattered in services
- ❌ Hard to understand domain rules
- ❌ Not DDD

### 3. Framework with ORM (Gin + GORM)

- ✅ Faster initial development
- ❌ Strong framework lock-in
- ❌ Difficult to migrate later
- ❌ Validation/sanitization unclear

### 4. Hexagonal Architecture (Ports & Adapters)

- ✅ Similar to chosen, focuses on inversion of control
- 🤔 More complex to implement in Go

## Implementation Details

### Project Structure

```
internal/notes/
├── domain/
│   ├── note.go              (Entities)
│   ├── note_content.go      (Value Object with validation)
│   └── note_content_test.go
├── application/
│   ├── port.go              (Interfaces: NoteRepository, NoteService)
│   ├── dto.go               (Input/Output structures)
│   └── note_service.go      (Business logic)
├── infrastructure/
│   ├── note_repository.go   (MongoDB implementation)
│   └── mongo_indexes.go     (Database setup)
├── presentation/
│   ├── handler.go           (HTTP handlers)
│   ├── router.go            (Route registration)
│   └── request_response.go  (DTOs)
└── mocks/
    └── mock_note_repository.go (For testing)
```

### Example: Value Object

```go
// internal/notes/domain/note_content.go
type NoteContent string

func NewNoteContent(value string) (NoteContent, error) {
    normalized := strings.TrimSpace(value)
    
    if normalized == "" {
        return "", ErrContentEmpty
    }
    
    if len(normalized) > 100000 {
        return "", ErrContentTooLong
    }
    
    // Sanitize XSS
    policy := bluemonday.StrictPolicy()
    sanitized := policy.Sanitize(normalized)
    
    return NoteContent(sanitized), nil
}
```

### Example: Repository Testing

```go
func TestCreateNote(t *testing.T) {
    // Mock repository
    repo := &mocks.MockNoteRepository{
        SaveFunc: func(ctx context.Context, note domain.Note) error {
            return nil // Mock: always success
        },
    }
    
    service := application.NewNoteService(repo)
    
    note, err := service.CreateNote(context.Background(), NoteInput{
        UserID:  "user123",
        Content: "Test note",
    })
    
    assert.NoError(t, err)
    assert.NotEmpty(t, note.ID)
}
```

## Data Flow

```
HTTP Request
    ↓
[Presentation] ParseRequest → CreateNoteRequest {Content: string}
    ↓
[Application] NoteInput {Content: string}
    ↓
[Service] service.CreateNote(input)
    ↓
[Domain] NewNoteContent(input.Content) ← ✅ Validation & Sanitization
    ↓
[Domain] NewNote(id, validContent) ← ✅ Create entity safely
    ↓
[Infrastructure] repo.Save(note) ← ✅ Persist
```

## Related Decisions

- ADR-002: HTTP Routing & Middleware (implements Presentation layer)
- [Value Objects Testing](/docs/adr/ADR-001-clean-architecture.md#example-value-object)

## References

- [Clean Architecture - Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design - Eric Evans](https://www.domainlanguage.com/ddd/)
- [Ports & Adapters Pattern](https://alistair.cockburn.us/hexagonal-architecture/)

---

**Key Takeaway**: Value Objects guarantee domain integrity, Repository Pattern ensures independence from persistence, layers ensure clear separation of concerns.
