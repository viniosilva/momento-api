# Plano de Implementação: Login com E-mail e Senha

## Resumo do Fluxo
Implementação do caso de uso de login, permitindo que usuários autentiquem-se com e-mail e senha. O sistema deve validar credenciais, comparar senha hasheada e retornar erro genérico para credenciais inválidas (por segurança).

## Sequência de Implementação

### 1. Domain Layer
**Status:** ✅ Reutilizar entidades e Value Objects existentes
- `Email` (Value Object) - já existe em `internal/auth/domain/email.go`
- `Password` (Value Object) - já existe em `internal/auth/domain/password.go` com método `Compare()`
- `User` (Entidade) - já existe em `internal/auth/domain/user.go`

**Tarefa:** Adicionar erro de domínio para credenciais inválidas
**Definição técnica:** Variável de erro `ErrInvalidCredentials` para uso quando e-mail não existe ou senha não confere.
**Arquivo:** `internal/auth/domain/user.go`

```go
var (
    ErrUserAlreadyExists = errors.New("user already exists")
    ErrUserNotFound      = errors.New("user not found")
    ErrInvalidCredentials = errors.New("invalid credentials")
)
```

---

### 2. Infrastructure Layer

**Tarefa:** Método FindByEmail no UserRepository
**Definição técnica:** Implementar método `FindByEmail(ctx context.Context, email domain.Email) (domain.User, error)` que busca usuário por e-mail (case-insensitive). Deve retornar `domain.ErrUserNotFound` quando não encontrar.
**Arquivo:** `internal/auth/infrastructure/user_repository.go`

```go
func (r *userRepository) FindByEmail(ctx context.Context, email domain.Email) (domain.User, error) {
    filter := bson.M{"email": string(email)}
    
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

---

### 3. Application Layer

**Tarefa:** Atualizar interface UserRepository (Port)
**Definição técnica:** Adicionar método `FindByEmail(ctx context.Context, email domain.Email) (domain.User, error)` na interface `UserRepository`.
**Arquivo:** `internal/auth/application/port.go`

```go
type UserRepository interface {
    Create(ctx context.Context, user domain.User) error
    ExistsByEmail(ctx context.Context, email domain.Email) (bool, error)
    FindByEmail(ctx context.Context, email domain.Email) (domain.User, error)
}
```

**Tarefa:** DTOs de Login
**Definição técnica:** Estruturas `LoginInput` (entrada com Email e Password como string) e `LoginOutput` (saída com ID, Email e timestamps - sem expor senha).
**Arquivo:** `internal/auth/application/dto.go`

```go
type LoginInput struct {
    Email    string
    Password string
}

type LoginOutput struct {
    ID        string
    Email     domain.Email
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

**Tarefa:** Atualizar interface AuthService (Port)
**Definição técnica:** Adicionar método `Login(ctx context.Context, input application.LoginInput) (application.LoginOutput, error)` na interface `AuthService`.
**Arquivo:** `internal/auth/presentation/port.go`

```go
type AuthService interface {
    Register(ctx context.Context, input application.UserInput) (application.UserOutput, error)
    Login(ctx context.Context, input application.LoginInput) (application.LoginOutput, error)
}
```

**Tarefa:** Método Login no AuthService
**Definição técnica:** Implementar caso de uso de login:
1. Validar e-mail (Fail-Fast via `domain.NewEmail`)
2. Buscar usuário por e-mail via repositório
3. Comparar senha fornecida com hash armazenado usando `password.Compare()`
4. Retornar erro genérico `domain.ErrInvalidCredentials` se usuário não existir OU senha não conferir (por segurança)
5. Retornar `LoginOutput` em caso de sucesso
**Arquivo:** `internal/auth/application/auth_service.go`

```go
func (s *authService) Login(ctx context.Context, input LoginInput) (LoginOutput, error) {
    email, err := domain.NewEmail(input.Email)
    if err != nil {
        return LoginOutput{}, err
    }

    user, err := s.userRepository.FindByEmail(ctx, email)
    if err != nil {
        if errors.Is(err, domain.ErrUserNotFound) {
            return LoginOutput{}, domain.ErrInvalidCredentials
        }
        return LoginOutput{}, fmt.Errorf("s.userRepository.FindByEmail: %w", err)
    }

    if err := user.Password.Compare(input.Password); err != nil {
        return LoginOutput{}, domain.ErrInvalidCredentials
    }

    return LoginOutput{
        ID:        user.ID.Hex(),
        Email:     user.Email,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }, nil
}
```

**Nota:** O método `Login` deve ser adicionado ao struct `authService` que implementa a interface `AuthService` definida em `internal/auth/presentation/port.go`.

---

### 4. Presentation Layer

**Tarefa:** Request e Response de Login
**Definição técnica:** Estruturas `LoginRequest` (com tags JSON e validação) e `LoginResponse` (com ID e Email). Reutilizar `ErrorResponse` existente.
**Arquivo:** `internal/auth/presentation/request_response.go`

```go
type LoginRequest struct {
    Email    string `json:"email" binding:"required" example:"user@example.com"`
    Password string `json:"password" binding:"required" example:"ValidPass123!"`
}

type LoginResponse struct {
    ID    string `json:"id" example:"507f1f77bcf86cd799439011"`
    Email string `json:"email" example:"user@example.com"`
}
```

**Tarefa:** Handler HTTP Login
**Definição técnica:** Implementar handler `Login` com:
- Decodificação do JSON do request body
- Tratamento de trim em e-mail e senha
- Chamada ao `authService.Login`
- Mapeamento de erros para HTTP status codes (400 para validação, 401 para credenciais inválidas, 500 para erros internos)
- Anotações Swaggo completas
**Arquivo:** `internal/auth/presentation/handler.go`

```go
// Login godoc
// @Summary Login with email and password
// @Description Authenticates a user with email and password credentials
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request"
// @Success 200 {object} LoginResponse "Login successful"
// @Failure 400 {object} ErrorResponse "Invalid request data"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        nethttp.JSON(w, http.StatusBadRequest, ErrorResponse{
            Message: "invalid request body",
        })
        return
    }

    input := application.LoginInput{
        Email:    strings.TrimSpace(req.Email),
        Password: strings.TrimSpace(req.Password),
    }

    output, err := h.authService.Login(r.Context(), input)
    if err != nil {
        statusCode, message := MapErrorToHTTPStatus(err)
        nethttp.JSON(w, statusCode, ErrorResponse{
            Message: message,
        })
        return
    }

    response := LoginResponse{
        ID:    output.ID,
        Email: string(output.Email),
    }

    nethttp.JSON(w, http.StatusOK, response)
}
```

**Tarefa:** Atualizar MapErrorToHTTPStatus
**Definição técnica:** Adicionar mapeamento de `domain.ErrInvalidCredentials` para HTTP 401 (Unauthorized) e mensagem genérica "E-mail ou senha incorretos".
**Arquivo:** `internal/auth/presentation/handler.go`

```go
func MapErrorToHTTPStatus(err error) (int, string) {
    if errors.Is(err, domain.ErrInvalidCredentials) {
        return http.StatusUnauthorized, "E-mail ou senha incorretos."
    }
    
    if errors.Is(err, domain.ErrUserNotFound) {
        return http.StatusUnauthorized, "E-mail ou senha incorretos."
    }
    
    // ... resto do código existente
}
```

**Tarefa:** Registrar rota de Login
**Definição técnica:** Adicionar rota `POST /api/auth/login` no router, aplicando o mesmo middleware de logging.
**Arquivo:** `internal/auth/presentation/router.go`

```go
func SetupRouter(options SetupRouterOptions) {
    handler := NewAuthHandler(options.AuthService)
    loggingMiddleware := makeLoggingMiddleware(options.Logger)

    registerHandler := addMiddleware(handler.Register, loggingMiddleware)
    options.Mux.Handle(fmt.Sprintf("POST %s/auth/register", options.Prefix), registerHandler)

    loginHandler := addMiddleware(handler.Login, loggingMiddleware)
    options.Mux.Handle(fmt.Sprintf("POST %s/auth/login", options.Prefix), loginHandler)
}
```

---

### 5. Orquestração
**Status:** ✅ Não requer alterações
- A injeção de dependências já está configurada em `cmd/api/main.go`
- O `AuthService` já é injetado no handler via `SetupRouter`

---

## Observações de Segurança

1. **Erro Genérico:** Sempre retornar "E-mail ou senha incorretos" para ambos os casos (usuário não existe ou senha incorreta) para evitar enumeração de usuários.

2. **Validação de E-mail:** O `domain.NewEmail` já faz trim e normaliza para lowercase, garantindo case-insensitive.

3. **Comparação de Senha:** Usar `password.Compare()` que utiliza `bcrypt.CompareHashAndPassword` internamente.

4. **Proteção contra Brute Force:** A história menciona bloqueio após 5 tentativas. Esta funcionalidade pode ser implementada em uma camada futura (middleware ou service de rate limiting), não faz parte do escopo inicial do login básico.

5. **SQL Injection:** Prevenção automática via driver MongoDB e uso de Value Objects.

6. **XSS:** Prevenção via serialização JSON adequada e não renderização de conteúdo do usuário diretamente em HTML (responsabilidade do frontend).

---

## Testes Unitários Recomendados

### Domain Layer
- Testar `ErrInvalidCredentials` quando necessário

### Application Layer
- Login com credenciais válidas
- Login com e-mail inválido (formato)
- Login com e-mail não cadastrado
- Login com senha incorreta
- Login com erro no repositório

### Infrastructure Layer
- `FindByEmail` retorna usuário quando existe
- `FindByEmail` retorna `ErrUserNotFound` quando não existe
- `FindByEmail` trata erros do MongoDB

### Presentation Layer
- Handler retorna 200 com credenciais válidas
- Handler retorna 400 com body inválido
- Handler retorna 400 com e-mail inválido
- Handler retorna 401 com credenciais inválidas (genérico)
- Handler retorna 500 com erro interno
