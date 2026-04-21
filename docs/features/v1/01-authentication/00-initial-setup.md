# Initial Setup - Configuração Inicial

### Status: DONE

## Visão Geral

Este documento estabelece a infraestrutura fundamental necessária para implementar todas as funcionalidades de autenticação do sistema. Os componentes descritos abaixo são compartilhados por múltiplas features e devem ser implementados primeiro para desbloquear o desenvolvimento das funcionalidades dependentes.

---

## Componentes de Domínio

### User Model

Estrutura que representa o usuário no sistema, contendo informações essenciais para autenticação.

**Localização:** `internal/auth/domain/user.go`

**Campos:**
- `id` - Identificador único do usuário (UUID)
- `email` - Email do usuário (Value Object)
- `password` - Senha hashada (Value Object)
- `created_at` - Data de criação
- `updated_at` - Data de atualização
- `email_verified_at` - Status de verificação do email (datetime)

**Features desbloqueadas:**
- Sign Up (01-sign-up.md)
- Login (03-login.md)
- Email Verification (02-email-verification.md)
- Reset Password (06-reset-password.md)

### Email Value Object

Value Object que encapsula a validação e manipulação de emails.

**Localização:** `internal/auth/domain/email.go`

**Funcionalidades:**
- Validação de formato RFC 5322
- Normalização (trim, lowercase)
- Criação com validação via construtor

**Erros:**
- `ErrEmailIsEmpty` - Email vazio
- `ErrInvalidEmail` - Formato inválido

**Features desbloqueadas:**
- Sign Up (01-sign-up.md)
- Login (03-login.md)
- Email Verification (02-email-verification.md)
- Reset Password (06-reset-password.md)

### Password Value Object

Value Object que encapsula a segurança da senha com hashing bcrypt.

**Localização:** `internal/auth/domain/password.go`

**Funcionalidades:**
- Validação de requisitos mínimos (6-64 caracteres, maiúscula, minúscula, número, símbolo)
- Geração de hash bcrypt com custo computacional adequado (cost=12)
- Comparação segura de senhas

**Erros:**
- `ErrInvalidPassword` - Senha inválida
- `ErrPasswordTooShort` - Senha muito curta
- `ErrPasswordTooLong` - Senha muito longa
- `ErrPasswordMissingUpper` - Falta maiúscula
- `ErrPasswordMissingLower` - Falta minúscula
- `ErrPasswordMissingNumber` - Falta número
- `ErrPasswordMissingSymbol` - Falta símbolo

**Features desbloqueadas:**
- Sign Up (01-sign-up.md)
- Login (03-login.md)
- Reset Password (06-reset-password.md)

### Domain Errors

Definição de erros de domínio para operações de autenticação.

**Localização:** `internal/auth/domain/user.go`

**Erros definidos:**
- `ErrUserAlreadyExists` - Usuário já existe
- `ErrUserNotFound` - Usuário não encontrado
- `ErrInvalidCredentials` - Credenciais inválidas
- `ErrRefreshTokenInvalid` - Refresh token inválido
- `ErrRefreshTokenExpired` - Refresh token expirado
- `ErrRefreshTokenNotFound` - Refresh token não encontrado

**Features desbloqueadas:**
- Sign Up (01-sign-up.md) - ErrUserAlreadyExists
- Login (03-login.md) - ErrUserNotFound, ErrInvalidCredentials
- Refresh Token (04-refresh-token.md) - ErrRefreshTokenInvalid, ErrRefreshTokenExpired, ErrRefreshTokenNotFound

---

## Serviços de Token

### JWT Service

Serviço para manipulação de JSON Web Tokens (access tokens).

**Localização:** `internal/auth/adapters/jwt_service.go`

**Interface:**
```go
type JWTService interface {
    Generate(userID string, email domain.Email) (string, error)
    Validate(tokenString string) (UserClaims, error)
}
```

**Funcionalidades:**
- Geração de tokens JWT com claims (userID, email, exp, iat)
- Validação de tokens (assinatura, expiração)
- Assinatura HMAC-SHA256

**Features desbloqueadas:**
- Login (03-login.md) - Geração de access token
- Refresh Token (04-refresh-token.md) - Geração de novo access token
- Protected Route (07-access-protected-route.md) - Validação de access token
- Expired Token (08-expired-token.md) - Tratamento de token expirado

### Secure Token Service

Serviço para manipulação de refresh tokens com armazenamento seguro.

**Localização:** `internal/auth/adapters/secure_token_service.go`

**Interface:**
```go
type SecureTokenService interface {
    Generate(ctx context.Context, userID, email string) (string, error)
    Refresh(ctx context.Context, token string) (userID, email, newToken string, err error)
}
```

**Funcionalidades:**
- Geração de tokens criptograficamente seguros (32 bytes, base64url)
- Armazenamento em Redis com TTL configurável
- Renovação de tokens (invalida antigo, cria novo)
- Refresh automático

**Features desbloqueadas:**
- Login (03-login.md) - Geração de refresh token
- Refresh Token (04-refresh-token.md) - Renovação de sessão
- Logout (05-logout.md) - Invalidação de token

---

## Repositórios

### User Repository

Repositório para persistência de usuários no MongoDB.

**Localização:** `internal/auth/adapters/user_repository.go`

**Interface:**
```go
type UserRepository interface {
    Create(ctx context.Context, user domain.User) error
    ExistsByEmail(ctx context.Context, email domain.Email) (bool, error)
    FindByEmail(ctx context.Context, email domain.Email) (domain.User, error)
    Update(ctx context.Context, user domain.User) error
}
```

**Funcionalidades:**
- Criação de usuário com validação de email único
- Verificação de existência de email
- Busca por email
- Atualização de usuário

**Índices:**
- Email único (restrição de duplicatas)

**Features desbloqueadas:**
- Sign Up (01-sign-up.md) - Create, ExistsByEmail
- Login (03-login.md) - FindByEmail
- Email Verification (02-email-verification.md) - Update
- Reset Password (06-reset-password.md) - Update

---

## Middleware

### Auth Middleware

Middleware para proteção de rotas via JWT.

**Localização:** `pkg/nethttp/auth/auth_middleware.go`

**Funcionalidades:**
- Extração de token do header Authorization (Bearer)
- Validação de token JWT
- Injeção de userID e email no contexto da requisição
- Resposta 401 para tokens inválidos/expirados

**Contexto injetado:**
- `user_id` - ID do usuário autenticado
- `email` - Email do usuário autenticado

**Features desbloqueadas:**
- Protected Route (07-access-protected-route.md)
- Expired Token (08-expired-token.md)

---

## Dependências Externas

### MongoDB

Banco de dados principal para persistência de usuários e tokens.

**Collections:**
- `users` - Armazenamento de usuários

### Redis

Armazenamento de refresh tokens com TTL.

**Funcionalidades:**
- Armazenamento temporário de refresh tokens
- Expiração automática de tokens
- Busca O(1) por token

---

## Fluxo de Implementação Recomendado

1. **Domain Layer** (primeiro)
   - User model
   - Email value object
   - Password value object
   - Domain errors

2. **Infraestrutura** (segundo)
   - MongoDB setup
   - Redis setup
   - User Repository

3. **Serviços** (terceiro)
   - JWT Service
   - Secure Token Service

4. **Application Layer** (quarto)
   - Auth Service (combina repositório + serviços)

5. **Presentation Layer** (quinto)
   - HTTP Handlers
   - Auth Middleware
   - Router

---

## Testes Unitários

Cada componente deve possui testes unitários cobrindo:

- Casos de sucesso
- Casos de erro/validação
- Edge cases

**Localização:** Arquivos com sufixo `_test.go` no mesmo diretório do componente.
