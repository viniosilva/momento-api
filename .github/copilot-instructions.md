---
description: "Copilot instructions for code generation following the Clean Architecture with DDD patterns in the Momento project"
---

# Copilot Instructions - Momento Code Patterns

This document defines the architectural and code patterns that should be followed when generating new features. **Consult the ADRs (`/docs/adr/`) when you have questions about architecture decisions.**

---

## 📋 Architecture: Clean Architecture with DDD

The project follows **Clean Architecture in 4 layers** with **Domain-Driven Design**:

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

**Dependency Rule**: `Presentation → Application → Domain ← Infrastructure`

❌ **NEVER**: Domain → Application or Domain → Presentation

---

## 🎯 Guide by Layer

### 1. Domain Layer (`/internal/{domain}/domain/`)

**Responsibility**: Business rules, Entities, Value Objects (NO external dependencies)

**Patterns**:

#### Value Objects
Encapsulate validation and guarantee the data is ALWAYS valid:

```go
package domain

import (
	"errors"
	"regexp"
	"strings"
)

// Define domain error
var ErrInvalidEmail = errors.New("invalid email format")

// Value Object: guarantee of validity
type Email string

func NewEmail(value string) (Email, error) {
	normalized := strings.TrimSpace(strings.ToLower(value))
	
	if normalized == "" {
		return "", ErrInvalidEmail
	}
	
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if !regexp.MustCompile(pattern).MatchString(normalized) {
		return "", ErrInvalidEmail
	}
	
	return Email(normalized), nil
}

// String() for logs
func (e Email) String() string {
	return string(e)
}
```

#### Entities
Aggregates with unique identity:

```go
package domain

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const UsersCollectionName = "users"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	Email     Email              `bson:"email"`
	Password  Password           `bson:"password"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func NewUser(email Email, password Password) User {
	now := time.Now().UTC()
	return User{
		ID:        primitive.NewObjectID(),
		Email:     email,
		Password:  password,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
```

---

### 2. Application Layer (`/internal/{domain}/application/`)

**Responsibility**: Orchestrated business logic, DTOs, Interfaces (Ports)

#### DTOs (Input/Output)
Only primitive types + time.Time:

```go
package application

import "time"

// Input: primitive types only
type UserInput struct {
	Email    string
	Password string
}

// Output: can contain Value Objects, distinct from Request/Response
type UserOutput struct {
	ID        string    // string, not primitive.ObjectID
	Email     domain.Email
	CreatedAt time.Time
	UpdatedAt time.Time
}
```

#### Ports (Interfaces)
Define contracts for external dependencies:

```go
package application

import (
	"context"
	"momento/internal/auth/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	ExistsByEmail(ctx context.Context, email domain.Email) (bool, error)
	GetByEmail(ctx context.Context, email domain.Email) (domain.User, error)
}

type JWTService interface {
	Generate(userID string) (string, error)
	Validate(token string) (claims interface{ GetUserID() string }, error)
}
```

#### Services
Orchestrate logic using Repositories and Value Objects:

```go
package application

import (
	"context"
	"errors"
	"fmt"
	"momento/internal/auth/domain"
)

type AuthService struct {
	userRepository UserRepository
	jwtService     JWTService
}

func NewAuthService(userRepository UserRepository, jwtService JWTService) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		jwtService:     jwtService,
	}
}

func (s *AuthService) Register(ctx context.Context, input UserInput) (UserOutput, error) {
	// 1. Create Value Objects (validation)
	email, err := domain.NewEmail(input.Email)
	if err != nil {
		return UserOutput{}, err // Domain validation error
	}

	password, err := domain.NewPassword(input.Password)
	if err != nil {
		return UserOutput{}, err
	}

	// 2. Business logic
	exists, err := s.userRepository.ExistsByEmail(ctx, email)
	if err != nil {
		return UserOutput{}, fmt.Errorf("s.userRepository.ExistsByEmail: %w", err)
	}
	if exists {
		return UserOutput{}, domain.ErrUserAlreadyExists
	}

	// 3. Create entity and persist
	user := domain.NewUser(email, password)
	if err := s.userRepository.Create(ctx, user); err != nil {
		return UserOutput{}, fmt.Errorf("s.userRepository.Create: %w", err)
	}

	return UserOutput{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
```

---

### 3. Infrastructure Layer (`/internal/{domain}/infrastructure/`)

**Responsibility**: Implementation of Ports, persistence, external integrations

#### Repository Pattern
Implement interfaces defined in Application:

```go
package infrastructure

import (
	"context"
	"errors"
	"momento/internal/auth/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) *userRepository {
	return &userRepository{
		collection: collection,
	}
}

func (r *userRepository) Create(ctx context.Context, user domain.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email domain.Email) (bool, error) {
	filter := bson.M{"email": email.String()}
	
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email domain.Email) (domain.User, error) {
	filter := bson.M{"email": email.String()}
	
	var user domain.User
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}
	
	return user, nil
}
```

#### MongoDB Indexes
Create indexes in separate files:

```go
// internal/auth/infrastructure/mongo_indexes.go
package infrastructure

import (
	"context"
	"momento/internal/auth/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateIndexes(ctx context.Context, db *mongo.Database) error {
	usersCollection := db.Collection(domain.UsersCollectionName)
	
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "email", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	
	_, err := usersCollection.Indexes().CreateOne(ctx, indexModel)
	return err
}
```

---

### 4. Presentation Layer (`/internal/{domain}/presentation/`)

**Responsibility**: HTTP Handlers, Routers, Request/Response structs

#### Request/Response Structs
Only primitive types with Swagger tags:

```go
package presentation

type RegisterRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"SecurePass123!"`
}

type UserResponse struct {
	ID        string `json:"id" example:"507f1f77bcf86cd799439011"`
	Email     string `json:"email" example:"user@example.com"`
	CreatedAt string `json:"created_at" example:"2026-02-08T10:30:00Z"`
	UpdatedAt string `json:"updated_at" example:"2026-02-08T10:30:00Z"`
}
```

#### Handler (Controller)
With complete Swagger documentation:

```go
package presentation

import (
	"encoding/json"
	"net/http"
	"time"

	"momento/internal/auth/application"
	"momento/internal/auth/domain"
	sharedresp "momento/internal/shared/presentation/response"
	nethttp_utils "momento/pkg/nethttp/utils"
)

type authHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *authHandler {
	return &authHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Create a new user account
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "User registration data"
// @Success 201 {object} UserResponse "User successfully created"
// @Failure 400 {object} response.ErrorResponse "Validation error"
// @Failure 409 {object} response.ErrorResponse "User already exists"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/auth/register [post]
func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		nethttp_utils.JSON(w, http.StatusBadRequest, sharedresp.ErrorResponse{
			Message: "invalid request body",
		})
		return
	}

	input := application.UserInput{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.authService.Register(r.Context(), input)
	if err != nil {
		statusCode, message := MapErrorToHTTPStatus(err)
		nethttp_utils.JSON(w, statusCode, sharedresp.ErrorResponse{
			Message: message,
		})
		return
	}

	res := UserResponse{
		ID:        output.ID,
		Email:     output.Email.String(),
		CreatedAt: output.CreatedAt.Format(time.RFC3339),
		UpdatedAt: output.UpdatedAt.Format(time.RFC3339),
	}

	nethttp_utils.JSON(w, http.StatusCreated, res)
}

// MapErrorToHTTPStatus maps domain errors to HTTP status codes
func MapErrorToHTTPStatus(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrInvalidEmail):
		return http.StatusBadRequest, "invalid email format"
	case errors.Is(err, domain.ErrUserAlreadyExists):
		return http.StatusConflict, "user already exists"
	case errors.Is(err, domain.ErrUserNotFound):
		return http.StatusNotFound, "user not found"
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
```

#### Router
With middleware chain and dependency injection:

```go
package presentation

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"momento/pkg/nethttp"
	auth "momento/pkg/nethttp/auth"
	logging "momento/pkg/nethttp/logging"
	sanitization "momento/pkg/nethttp/sanitization"
)

type SetupRouterOptions struct {
	Mux         *http.ServeMux
	Prefix      string
	AuthService AuthService
	JWTService  JWTService
	Logger      *slog.Logger
	Timeout     *time.Duration
}

func SetupRouter(options SetupRouterOptions) {
	handler := NewAuthHandler(options.AuthService)

	// Middleware chain (order matters)
	chain := nethttp.NewDefaultChain(options.Logger, nethttp.WithTimeout(options.Timeout))
	chain.AddMiddleware(sanitization.SanitizationMiddleware())
	chain.AddMiddleware(logging.LoggingMiddleware(options.Logger))

	// Public routes
	options.Mux.Handle(
		fmt.Sprintf("POST %s/auth/register", options.Prefix),
		chain.ThenFunc(handler.Register),
	)

	// Middleware with auth for authenticated routes
	authChain := nethttp.NewDefaultChain(options.Logger, nethttp.WithTimeout(options.Timeout))
	authChain.AddMiddleware(sanitization.SanitizationMiddleware())
	authChain.AddMiddleware(logging.LoggingMiddleware(options.Logger))
	authChain.AddMiddleware(auth.AuthMiddleware(options.JWTService))

	options.Mux.Handle(
		fmt.Sprintf("GET %s/auth/profile", options.Prefix),
		authChain.ThenFunc(handler.GetProfile),
	)
}
```

---

## 🛠️ Specific Patterns

### Error Handling
1. **Domain errors**: Define with `errors.New()` and return directly
2. **Wrapping**: Use `fmt.Errorf("context: %w", err)` to add context
3. **HTTP Mapping**: Translate domain errors to status codes in Presentation

```go
// ✅ GOOD
if err != nil {
	return UserOutput{}, fmt.Errorf("s.userRepository.Create: %w", err)
}

// ❌ AVOID
if err != nil {
	return UserOutput{}, err // Loses context
}
```

### Structured Logging
Use `*slog.Logger` (stdlib) with structured fields:

```go
logger.InfoContext(ctx, "user registered",
	"user_id", userID,
	"email", email,
)

logger.ErrorContext(ctx, "registration failed",
	"error", err,
	"email", email,
)
```

### HTTP Utilities
Use `nethttp_utils.JSON()` for standardized responses:

```go
nethttp_utils.JSON(w, http.StatusOK, UserResponse{
	ID:    "123",
	Email: "user@example.com",
})
```

### MongoDB
- **Indexes**: Define in `infrastructure/mongo_indexes.go`
- **Conversion**: Use `primitive.ObjectIDFromHex()` with error handling
- **BSON tags**: Always use `bson` tags (not `json`)

```go
// ✅ GOOD
type User struct {
	ID    primitive.ObjectID `bson:"_id"`
	Email Email              `bson:"email"`
}

// ❌ AVOID
type User struct {
	ID    primitive.ObjectID `json:"_id"` // Always use bson
	Email Email              `json:"email"`
}
```

---

## 📖 References

### ADRs (Architecture Decision Records)
- **[ADR-001: Clean Architecture with DDD](/docs/adr/ADR-001-clean-architecture-ddd.md)** - 4-layer architecture, Value Objects Pattern, Repository Pattern
- **[ADR-002: HTTP Routing and Middleware](/docs/adr/ADR-002-http-routing-middleware.md)** - Use of `http.ServeMux`, middleware chain, input validation

### Examples in the Project
- **Handlers**: `internal/notes/presentation/handler.go`, `internal/auth/presentation/handler.go`
- **Services**: `internal/notes/application/note_service.go`, `internal/auth/application/auth_service.go`
- **Repositories**: `internal/notes/infrastructure/note_repository.go`, `internal/auth/infrastructure/user_repository.go`
- **Routers**: `internal/notes/presentation/router.go`, `internal/auth/presentation/router.go`

---

## ✅ Checklist before generating code

- [ ] Did I consult the ADRs if the decision is architectural?
- [ ] Does the new layer follow the dependency hierarchy (Presentation → Application → Domain ← Infrastructure)?
- [ ] Are validations in the Domain (Value Objects)?
- [ ] Do DTOs use only primitive types?
- [ ] Do Handlers have complete Swagger documentation?
- [ ] Are domain errors mapped to HTTP status codes?
- [ ] Do Repositories implement Ports defined in Application?
- [ ] Does MongoDB use BSON tags and indexes?
- [ ] Is logging using structured format with `slog`?

