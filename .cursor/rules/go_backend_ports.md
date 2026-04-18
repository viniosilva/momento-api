# Ports Layer - Handlers and Routes

## Location
`internal/{module}/ports/`

## Required Files
- `request_response.go` - Request (Swagger tags) and Response (primitives)
- `port.go` - Service interface (decoupling)
- `handler.go` - Handler implementation + `MapErrorToHTTPStatus`
- `handler_test.go` - Tests (one per possible status code)
- `router.go` - `SetupRouterOptions` + route registration
- `router_test.go` - Route tests

## Critical Rules
- ✅ COMPLETE Swagger (@Summary, @Description, @Tags, @Param, @Success, @Failure, @Router)
- ✅ Use `nethttp.JSON` (from `pkg`)
- ✅ Clean input: `strings.TrimSpace`
- ✅ Separate Request/Response from DTOs (ports doesn't expose domain models)
- ✅ `MapErrorToHTTPStatus` maps errors → status codes
- ✅ Use `r.Context()` when calling service
- ❌ DO NOT put business logic in handler
- ❌ DO NOT expose internal errors (500 = "internal server error")
- ❌ DO NOT import domain directly (use app.DTOs)

## Middleware Chain (Required Order)

```
Request
    ↓
Recovery     (outermost - catch panics)
    ↓
Logging      (record timing, request/response)
    ↓
Sanitization (validate structure: JSON, size, depth)
    ↓
Auth         (verify JWT token - skip for public routes)
    ↓
Handler      (business logic)
```

Implementation in `pkg/nethttp`:
```go
chain := nethttp.NewDefaultChain(logger)
chain.AddMiddleware(recovery.RecoveryMiddleware())
chain.AddMiddleware(logging.LoggingMiddleware(logger))
chain.AddMiddleware(sanitization.SanitizationMiddleware())
chain.AddMiddleware(auth.AuthMiddleware(jwtService))

mux.Handle("POST /api/notes", chain.ThenFunc(handler.CreateNote))
```

## Input Validation Strategy

**Separate concerns**:
- **Middleware**: Structural validation (JSON valid, size ≤ 1MB, nesting depth ≤ 10)
- **Domain**: Content sanitization + type-safe validation (XSS prevention via Value Objects)

## Patterns

### request_response.go
```go
type RegisterRequest struct {
    Email string `json:"email" example:"user@example.com"`
}
type RegisterResponse struct { ID string `json:"id"` }
type ErrorResponse struct { Message string `json:"message"` }
```

### port.go
```go
type AuthService interface {
    Register(ctx context.Context, input app.UserInput) (app.UserOutput, error)
}
```

### handler.go
```go
// @Summary Register a new user
// @Tags auth
// @Param request body RegisterRequest true "Registration request"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    // 1. Parse, 2. Clean, 3. Call service, 4. Return
}
func MapErrorToHTTPStatus(err error) (int, string) { /* maps */ }
```

### router.go
```go
type SetupRouterOptions struct { Mux *http.ServeMux; Prefix string; AuthService AuthService; Logger *slog.Logger }
func SetupRouter(opts SetupRouterOptions) { /* registers routes */ }
```

## Swagger (Required on All Endpoints)

### In main.go (header)
```go
// @title Momento API
// @version 1.0
// @description API documentation for Momento application
// @BasePath /api
```

### In structs
```go
type RegisterRequest struct {
    Email string `json:"email" example:"user@example.com"`
}
```