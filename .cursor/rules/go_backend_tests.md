# Tests - Unit and Mocks

## Philosophy
- **Pragmatic coverage**: 80% is good. 100% is often negative ROI.
- **Priority**: domain tests > app > adapters > ports
- **Maintenance**: Tests that break on refactors indicate bad design

## Strategy by Layer

| Layer | Type | How to test |
|-------|------|-------------|
| `domain` | Unit | VOs validation, entity constructors, domain errors |
| `app` | Unit | Service with mocks |
| `adapters` | Unit | External services (JWT, APIs) |
| `ports` | Unit | Handlers with mocks, Routes with httptest |

## Unit Tests

### Structure (Arrange-Act-Assert)
```go
func TestAuthService_Register(t *testing.T) {
    t.Run("should create user successfully", func(t *testing.T) {
        mock := mocks.NewMockUserRepository(t)
        mock.EXPECT().ExistsByEmail(mock.Anything, mock.Anything).Return(false, nil).Once()
        mock.EXPECT().Create(mock.Anything, mock.Anything).Return(nil).Once()
        svc := app.NewAuthService(mock)

        got, err := svc.Register(t.Context(), defaultInput)
        require.NoError(t, err)

        assert.NotEmpty(t, got.ID)
    })
}
```

### Naming
- Function: `Test{StructName}_{MethodName}`
- Subtest: `t.Run("should {behavior} when {condition}")`

### Test Package
```go
package domain_test  // ALWAYS with _test suffix
```

## Mocks

### Generate with Mockery
```bash
make mock  # defined in Makefile
```

### Rule: only for interfaces in `app/port.go`
```go
// mocks auto-generated in /mocks/
type MockUserRepository struct { ... }
```

### In test code
```go
mock := mocks.NewMockUserRepository(t)
mock.EXPECT().Method(args).Return(result, nil).Once()
```

## Coverage

### Target
- Ideal: 80%
- Exceptional: 90%+ (only for critical domains)

### Run
```bash
make test        # run tests
make coverage    # view coverage report
```

## Handler Tests

```go
func TestAuthHandler_Register(t *testing.T) {
    userID := primitive.NewObjectID()

    t.Run("should return status created when note is created successfully", func(t *testing.T) {
        mockRepo := mocks.NewMockUserRepository(t)
        svc := application.NewUserService(mockRepo)
        handler := presentation.NewAuthHandler(svc)

        mockRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(nil).Once()

        mux := http.NewServeMux()
        mux.HandleFunc("POST /auth/register", func(w http.ResponseWriter, r *http.Request) {
            ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID.Hex())
            handler.Register(w, r.WithContext(ctx))
        })

        reqBody := map[string]any{
            "email":    "test@example.com",
            "password": "ValidPass123©",
        }
        body, _ := json.Marshal(reqBody)
        req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/auth/register", bytes.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        rec := httptest.NewRecorder()
        mux.ServeHTTP(rec, req)

        assert.Equal(t, http.StatusCreated, rec.Code)

        var got presentation.RegisterResponse
        err := json.Unmarshal(rec.Body.Bytes(), &got)
        require.NoError(t, err)

        assert.NotEmpty(t, got.ID)
    })

    t.Run("should return status bad request when email is invalid", func(t *testing.T) {
        mockRepo := mocks.NewMockUserRepository(t)
        svc := application.NewUserService(mockRepo)
        handler := presentation.NewAuthHandler(svc)

        mux := http.NewServeMux()
        mux.HandleFunc("POST /auth/register", func(w http.ResponseWriter, r *http.Request) {
            ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID.Hex())
            handler.Register(w, r.WithContext(ctx))
        })

        reqBody := map[string]any{
            "email":    "invalid-email",
            "password": "ValidPass123©",
        }
        body, _ := json.Marshal(reqBody)
        req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/auth/register", bytes.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        rec := httptest.NewRecorder()
        mux.ServeHTTP(rec, req)

        assert.Equal(t, http.StatusBadRequest, rec.Code)

        var got sharedresp.ErrorResponse
        err := json.NewDecoder(rec.Body).Decode(&got)
        require.NoError(t, err)

        assert.Equal(t, "invalid email format", got.Message)
    })

    t.Run("should return status internal server error when service returns error", func(t *testing.T) {
        mockRepo := mocks.NewMockUserRepository(t)
        svc := application.NewUserService(mockRepo)
        handler := presentation.NewAuthHandler(svc)

        mockRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(assert.AnError).Once()

        mux := http.NewServeMux()
        mux.HandleFunc("POST /auth/register", func(w http.ResponseWriter, r *http.Request) {
            ctx := context.WithValue(r.Context(), nethttp_auth.ContextKeyUserID, userID.Hex())
            handler.Register(w, r.WithContext(ctx))
        })

        reqBody := map[string]any{
            "email":    "test@example.com",
            "password": "ValidPass123©",
        }
        body, _ := json.Marshal(reqBody)
        req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/auth/register", bytes.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        rec := httptest.NewRecorder()
        mux.ServeHTTP(rec, req)

        assert.Equal(t, http.StatusInternalServerError, rec.Code)

        var got sharedresp.ErrorResponse
        err := json.NewDecoder(rec.Body).Decode(&got)
        require.NoError(t, err)

        assert.Equal(t, "internal server error", got.Message)
    })
}
```

## Anti-Patterns

```go
// ❌ DON'T:
// Test internal implementation instead of behavior
// Mock everything (loses test value)
// Tests that only test wrappers (useless coverage)
// Ignore failing tests (symptom of bad design)
// Leave failing tests in CI
```

## Quick Reference

| What | Where | How |
|------|-------|-----|
| Value Object validation | domain | unit |
| Entity constructor | domain | unit |
| Domain errors | domain | unit |
| Service logic | app | unit with mocks |
| External services (JWT) | adapters | unit |
| HTTP Handler | ports | unit with mocks |
| Router | ports | httptest |