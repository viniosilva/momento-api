# Style and Conventions - Code Patterns

## Object Calisthenics
Pragmatic guidelines. See `@.cursor/rules/go_backend_object_calisthenics.md` for full details. Simplicity > Rules

## Naming

- Public structs: `PascalCase` (e.g., `UserRepository`)
- Private structs: `camelCase` (e.g., `userRepository`)
- Constructors: `New{Name}` (e.g., `NewUserService`)
- Interfaces: behavior (e.g., `UserRepository`, not `IUserRepository`)
- Errors: `Err{Name}` (e.g., `ErrUserNotFound`)
- Files: `{name}.go`, `{name}_test.go`
- Test packages: `{name}_test` (e.g., `app_test`)

## Imports (order)
```go
import (
    "context"  // stdlib
    "github.com/stretchr/testify/assert"  // externals
    "momento/internal/auth/domain"  // internals
)
```

## Error Handling

**Flow**: Domain → Adapters (translate) → App (enrich) → Ports (map to HTTP)

- Domain errors: `var Err... = errors.New("...")`
- Compare: `errors.Is(err, domain.ErrExpected)`
- Wrap: `fmt.Errorf("s.repo.Find: %w", err)`
- DO NOT expose internal details (500 = "internal server error")

## SOLID Principles

- **S** (Single Responsibility): One responsibility per struct/file
- **O** (Open/Closed): Extensible via interfaces
- **L** (Liskov): Implementations replaceable by interfaces
- **I** (Interface Segregation): Small, focused interfaces
- **D** (Dependency Inversion): Depend on abstractions, inject dependencies

## Anti-Patterns (Avoid)

```go
// ❌ DON'T DO:
if err != nil { return err } else { return nil }    // use guard clause
import "momento/internal/..." // in pkg/              // pkg can't import internal
func NewService() userService { }                    // return *userService
return err                                           // wrap: fmt.Errorf("ctx: %w", err)
func (s *service) Create(input Input) error         // missing ctx context.Context
func NewClient() { host := os.Getenv("HOST") }      // receive as parameter
logger.Info("creating user") // in service          // DON'T log! Middleware does it
```

## Comments
Only when necessary. Explain "why", not "what".