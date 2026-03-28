# ADR-002: HTTP Routing, Input Validation & Middleware Chain

**Status**: ✅ Accepted

**Date**: 2026-03-28

## Context

The project needed to choose:
- How to route HTTP requests (framework vs standard library)
- How to validate/sanitize inputs
- How to apply cross-cutting concerns (logging, auth, etc)

Constraints:
- Monolith Go project, minimize external dependencies
- Need middleware for recovery, logging, authentication
- Want clear request validation strategy
- Build-in `net/http` available in Go 1.22+

## Decision

### 1. Use `http.ServeMux` (No Framework)

Go 1.22+ supports path parameters natively:

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /notes/{id}", handler.GetUserNoteByID)
mux.HandleFunc("POST /notes", handler.CreateNote)
```

**Benefits**:
- Zero external dependencies for routing
- Better performance (no reflection)
- Explicit code, easy to debug
- Covers 95% of use cases

**Trade-off**: Manual middleware chain (solution: see below)

### 2. Implement Middleware Chain Pattern

Custom composition of `func(http.Handler) http.Handler`:

```go
chain := nethttp.NewDefaultChain(logger)
chain.AddMiddleware(recovery.RecoveryMiddleware())      // Catch panics
chain.AddMiddleware(logging.LoggingMiddleware())        // Log requests
chain.AddMiddleware(sanitization.SanitizationMiddleware()) // Validate structure
chain.AddMiddleware(auth.AuthMiddleware(jwtService))    // Verify token

mux.Handle("POST /api/notes", chain.ThenFunc(handler.CreateNote))
```

**Execution order**: Recovery → Logging → Sanitization → Auth → Handler

### 3. Input Validation Strategy

**Separate concerns**:
- **Middleware**: Validate structure (JSON valid, size ≤ 1MB, nesting depth ≤ 10)
- **Domain Value Objects**: Sanitize content (remove XSS) + type-safe validation

```go
// Middleware: Structural validation
func SanitizationMiddleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB limit
            
            body, err := io.ReadAll(r.Body)
            if err != nil {
                nethttp_utils.JSON(w, http.StatusBadRequest,
                    map[string]string{"message": "request body too large"})
                return
            }
            
            var data any
            if err := json.Unmarshal(body, &data); err != nil {
                nethttp_utils.JSON(w, http.StatusBadRequest,
                    map[string]string{"message": "invalid JSON"})
                return
            }
            
            // Validate structure (size, depth, key length)
            if err := validateStructure(data, 0); err != nil {
                nethttp_utils.JSON(w, http.StatusBadRequest,
                    map[string]string{"message": err.Error()})
                return
            }
            
            // Restore body without re-marshaling
            r.Body = io.NopCloser(bytes.NewReader(body))
            next.ServeHTTP(w, r)
        })
    }
}

// Domain: Content validation & sanitization
func NewNoteContent(value string) (NoteContent, error) {
    normalized := strings.TrimSpace(value)
    
    if normalized == "" {
        return "", ErrContentEmpty
    }
    
    // Remove XSS using bluemonday
    sanitized := bluemonday.StrictPolicy().Sanitize(normalized)
    return NoteContent(sanitized), nil
}
```

## Consequences

### Positive ✅

**Routing**:
- Zero dependencies for routing
- Better debuggability
- Explicit request handling
- ~0.1ms faster per request vs framework

**Middleware Chain**:
- Flexible composition
- Clear execution order
- Easy to test individually
- Reusable across multiple routes

**Input Validation**:
- Defense in depth (middleware + domain)
- XSS prevention at domain level (guaranteed)
- Type-safe via Value Objects
- Structure validation prevents DoS

### Negative ❌

**Routing**:
- Manual middleware needed
- No automatic route versioning
- Limited to ~100-200 routes before complexity

**Middleware Chain**:
- Backward iteration for execution order
- More boilerplate than framework

**Input Validation**:
- Overhead of re-parsing JSON (mitigated by not re-marshaling)
- Developers must use Value Objects (needs discipline)

## Alternatives Considered

### 1. Framework Router (Gin, Echo, Chi)

```go
router := gin.Default()
router.POST("/notes", handler.CreateNote)
```

- ✅ Middleware chain built-in
- ✅ Automatic request binding
- ❌ ~30+ transitive dependencies
- ❌ Framework lock-in
- ❌ Harder to migrate away

### 2. Gorilla Mux (Older Alternative)

- ❌ Less maintained
- ❌ Slower performance
- ❌ Complex API

### 3. Global Middleware Only

```go
handler := logging.LoggingMiddleware(
    auth.AuthMiddleware(
        mux,
    ),
)
http.ListenAndServe(":8080", handler)
```

- ✅ Simpler
- ❌ Applies all middleware to all routes
- ❌ No flexibility (can't skip auth on public routes)
- ❌ Order hard to understand

### 4. Handler Validation (No Middleware)

```go
func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
    // Validate here
    if r.ContentLength > 1<<20 {
        http.Error(w, "too large", 400)
        return
    }
}
```

- ✅ Simpler
- ❌ Repeated in every handler
- ❌ Inconsistent validation
- ❌ Security gaps

## Implementation Details

### Middleware Execution Order

```
Request Flow:
┌──────────────────┐
│     Request      │
└────────┬─────────┘
         ↓
    Recovery      (outermost - catch panics)
         ↓
    Logging       (record timing)
         ↓
    Sanitization  (validate structure: JSON, size, depth)
         ↓
    Auth          (verify JWT token)
         ↓
    Handler       (business logic)
         ↓
    Auth (return) (middleware ordering maintained)
         ↓
    Sanitization (return)
         ↓
    Logging      (return - record response)
         ↓
    Recovery     (return - innermost)
         ↓
    Response
```

### Chain Implementation

```go
// pkg/nethttp/chain.go
type Chain struct {
    middlewares []func(http.Handler) http.Handler
    logger      *slog.Logger
    timeout     *time.Duration
}

func NewDefaultChain(logger *slog.Logger, opts ...Option) *Chain {
    return &Chain{
        middlewares: []func(http.Handler) http.Handler{},
        logger:      logger,
    }
}

func (c *Chain) AddMiddleware(m func(http.Handler) http.Handler) *Chain {
    c.middlewares = append(c.middlewares, m)
    return c
}

// Compose: Apply middlewares in REVERSE order (last added executes first)
func (c *Chain) Then(h http.Handler) http.Handler {
    for i := len(c.middlewares) - 1; i >= 0; i-- {
        h = c.middlewares[i](h)
    }
    return h
}

func (c *Chain) ThenFunc(h http.HandlerFunc) http.Handler {
    return c.Then(http.HandlerFunc(h))
}
```

### Using Chain in Router

```go
// internal/notes/presentation/router.go
func SetupRouter(options SetupRouterOptions) {
    handler := NewNoteHandler(options.NoteService)
    
    chain := nethttp.NewDefaultChain(options.Logger)
    chain.AddMiddleware(recovery.RecoveryMiddleware())
    chain.AddMiddleware(logging.LoggingMiddleware(options.Logger))
    chain.AddMiddleware(nethttp_sanitization.SanitizationMiddleware())
    chain.AddMiddleware(auth.AuthMiddleware(options.JWTService))
    
    options.Mux.Handle("POST /api/notes", chain.ThenFunc(handler.CreateNote))
    options.Mux.Handle("GET /api/notes", chain.ThenFunc(handler.ListNotes))
    options.Mux.Handle("GET /api/notes/{id}", chain.ThenFunc(handler.GetUserNoteByID))
}
```

### Example Middlewares

#### Recovery Middleware

```go
func RecoveryMiddleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            defer func() {
                if err := recover(); err != nil {
                    logger.Error("panic recovered", "error", err)
                    nethttp_utils.JSON(w, http.StatusInternalServerError,
                        map[string]string{"message": "internal server error"})
                }
            }()
            next.ServeHTTP(w, r)
        })
    }
}
```

#### Auth Middleware

```go
const ContextKeyUserID = "userID"

func AuthMiddleware(jwtService JWTService) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            token := extractBearerToken(r)
            if token == "" {
                nethttp_utils.JSON(w, http.StatusUnauthorized,
                    map[string]string{"message": "missing token"})
                return
            }
            
            userID, err := jwtService.VerifyToken(token)
            if err != nil {
                nethttp_utils.JSON(w, http.StatusUnauthorized,
                    map[string]string{"message": "invalid token"})
                return
            }
            
            ctx := context.WithValue(r.Context(), ContextKeyUserID, userID)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

### Sanitization Strategy

**Middleware level** (structural):
```
- Content-Length validation (max 1MB)
- JSON format validation
- Nesting depth check (max 10 levels)
- Key length check (max 100 chars)
```

**Domain level** (content):
```
- XSS removal via bluemonday.StrictPolicy()
- Type validation in Value Objects
- Business rule validation
```

Result: Input guaranteed valid AND safe at domain level.

## Testing Middleware Chain

```go
func TestMiddlewareChain(t *testing.T) {
    calls := []string{}
    
    mw1 := func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            calls = append(calls, "mw1-start")
            next.ServeHTTP(w, r)
            calls = append(calls, "mw1-end")
        })
    }
    
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        calls = append(calls, "handler")
        w.WriteHeader(http.StatusOK)
    })
    
    chain := NewDefaultChain(nil)
    chain.AddMiddleware(mw1)
    
    req := httptest.NewRequest("GET", "/", nil)
    rec := httptest.NewRecorder()
    
    chain.Then(handler).ServeHTTP(rec, req)
    
    assert.Equal(t, []string{"mw1-start", "handler", "mw1-end"}, calls)
}
```

## Security Considerations

1. **Size limits**: Prevent memory exhaustion
2. **Depth limits**: Prevent DoS via deeply nested JSON
3. **Key validation**: Block oversized keys
4. **XSS prevention**: bluemonday.StrictPolicy() removes dangerous tags
5. **JWT validation**: Auth middleware verifies all tokens

## Related Decisions

- ADR-001: Clean Architecture (Presentation layer uses this routing)

## References

- [Go 1.22 Release Notes: Path Parameters](https://go.dev/blog/go1.22)
- [Middleware Pattern in Go](https://www.alexedwards.net/blog/making-and-using-middleware)
- [bluemonday: XSS Prevention](https://github.com/microcosm-cc/bluemonday)
- [OWASP Input Validation](https://cheatsheetseries.owasp.org/cheatsheets/Input_Validation_Cheat_Sheet.html)

---

**Key Takeaway**: Separate structural validation (middleware) from content validation (domain). Use Value Objects for type safety. Chain middlewares for clean, composable request handling.
