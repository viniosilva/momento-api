# Guia de Planejamento - Cursor Plan Mode

Use este guia para **decompor User Stories em tarefas técnicas** antes de implementar.

> **Objetivo**: Criar um plano estruturado que possa ser executado passo a passo seguindo `.cursorrules`.

---

## 📋 Checklist de Decomposição

Para cada User Story, decomponha em **7 fases** seguindo Clean Architecture:

### **1. Domain Layer** (Regras de Negócio)
- [ ] Identificar Value Objects necessários (Email, Password, etc)
- [ ] Identificar Entidades (User, Product, etc)
- [ ] Definir erros de domínio (`var Err...`)
- [ ] Definir constantes (collection names, status, etc)

### **2. Application Layer - Contratos** (DTOs e Interfaces)
- [ ] Definir Input DTOs (tipos primitivos)
- [ ] Definir Output DTOs (pode ter Value Objects)
- [ ] Definir interfaces (Ports) para dependências

### **3. Infrastructure Layer** (Persistência e Integrações)
- [ ] Criar índices MongoDB (se necessário)
- [ ] Implementar Repository (se necessário)
- [ ] Implementar Services externos (JWT, APIs, etc)

### **4. Application Layer - Lógica** (Casos de Uso)
- [ ] Implementar Service com orquestração
- [ ] Validações via Value Objects
- [ ] Lógica de negócio

### **5. Presentation Layer** (HTTP Interface)
- [ ] Definir Request structs (com tags Swagger)
- [ ] Definir Response structs (primitivos)
- [ ] Definir interface do service (Port)
- [ ] Implementar Handler com Swagger completo
- [ ] Mapear erros para HTTP status codes
- [ ] Registrar rotas no Router

### **6. Orquestração** (Injeção de Dependências)
- [ ] Atualizar `cmd/api/main.go` com novo módulo
- [ ] Injetar dependências (Infrastructure → Application → Presentation)
- [ ] Registrar rotas

### **7. Finalização**
- [ ] Gerar Swagger: `make swag`
- [ ] Gerar mocks: `make mock`
- [ ] Executar testes: `make test`

---

## 📝 Formato de Output

Para cada fase, liste as tarefas assim:

```
### Fase [N]: [Nome da Fase]

**Tarefa**: [Nome do componente]
**Arquivo**: `caminho/do/arquivo.go`
**Definição**:
- [O que criar: struct, função, interface]
- [Campos/parâmetros principais]
- [Validações ou regras principais]
```

---

## 🎯 Exemplo Prático

**Input**: "Como usuário, quero me cadastrar via e-mail e senha"

**Output do Plan Mode**:

```markdown
## Plano Técnico: Cadastro de Usuário

### Fase 1: Domain Layer

**Tarefa**: Value Object Email
**Arquivo**: `internal/auth/domain/email.go`
**Definição**:
- Type: `Email string`
- Construtor: `NewEmail(value string) (Email, error)`
- Validação: RFC 5322 (usar `mail.ParseAddress`)
- Normalização: lowercase e trim

**Tarefa**: Value Object Password  
**Arquivo**: `internal/auth/domain/password.go`
**Definição**:
- Type: `Password string`
- Construtor: `NewPassword(plain string) (Password, error)`
- Validações: min 6, max 64, uppercase, lowercase, número, símbolo
- Hash: bcrypt com cost 12

**Tarefa**: Entidade User
**Arquivo**: `internal/auth/domain/user.go`
**Definição**:
- Struct: `User { ID, Email, Password, CreatedAt, UpdatedAt }`
- Construtor: `NewUser(email Email, password Password) User`
- Constante: `UsersCollectionName = "users"`

**Tarefa**: Erros de Domínio
**Arquivo**: `internal/auth/domain/user.go`
**Definição**:
- `ErrUserAlreadyExists = errors.New("user already exists")`
- `ErrUserNotFound = errors.New("user not found")`

---

### Fase 2: Application Layer - Contratos

**Tarefa**: DTOs
**Arquivo**: `internal/auth/application/dto.go`
**Definição**:
- Input: `UserInput { Email string; Password string }`
- Output: `UserOutput { ID string; Email domain.Email; CreatedAt time.Time }`

**Tarefa**: Interfaces (Ports)
**Arquivo**: `internal/auth/application/port.go`
**Definição**:
```go
type UserRepository interface {
    Create(ctx context.Context, user domain.User) error
    ExistsByEmail(ctx context.Context, email domain.Email) (bool, error)
}
```

---

### Fase 3: Infrastructure Layer

**Tarefa**: MongoDB Index
**Arquivo**: `internal/auth/infrastructure/indexes/001_users_unique_email.go`
**Definição**:
- Função: `CreateUserEmailIndex(ctx, db) error`
- Index: `email` field com `unique: true`

**Tarefa**: Repository
**Arquivo**: `internal/auth/infrastructure/user_repository.go`
**Definição**:
- Struct privada: `userRepository { collection *mongo.Collection }`
- Implementar: `Create()`, `ExistsByEmail()`
- Traduzir: `mongo.IsDuplicateKeyError` → `domain.ErrUserAlreadyExists`

---

### Fase 4: Application Layer - Service

**Tarefa**: Auth Service
**Arquivo**: `internal/auth/application/auth_service.go`
**Definição**:
- Struct privada: `authService { userRepository UserRepository }`
- Método: `Register(ctx context.Context, input UserInput) (UserOutput, error)`
- Fluxo:
  1. Validar email via `domain.NewEmail()`
  2. Validar senha via `domain.NewPassword()`
  3. Verificar se email existe via `ExistsByEmail()`
  4. Criar user via `Create()`
  5. Retornar `UserOutput`

---

### Fase 5: Presentation Layer

**Tarefa**: Request/Response
**Arquivo**: `internal/auth/presentation/request_response.go`
**Definição**:
- `RegisterRequest { Email string; Password string }` com tags Swagger
- `RegisterResponse { ID string; Email string }`
- `ErrorResponse { Message string }`

**Tarefa**: Port
**Arquivo**: `internal/auth/presentation/port.go`
**Definição**:
```go
type AuthService interface {
    Register(ctx, input application.UserInput) (application.UserOutput, error)
}
```

**Tarefa**: Handler
**Arquivo**: `internal/auth/presentation/handler.go`
**Definição**:
- Struct: `AuthHandler { authService AuthService }`
- Método: `Register(w http.ResponseWriter, r *http.Request)`
- Swagger completo com @Summary, @Tags, @Param, @Success, @Failure, @Router
- Função: `MapErrorToHTTPStatus(err error) (int, string)`

**Tarefa**: Router
**Arquivo**: `internal/auth/presentation/router.go`
**Definição**:
- Função: `SetupRouter(opts SetupRouterOptions)`
- Registrar: `POST /api/auth/register`

---

### Fase 6: Orquestração

**Tarefa**: Main
**Arquivo**: `cmd/api/main.go`
**Definição**:
1. Criar índices: `authinfra.CreateIndexes(ctx, db)`
2. Instanciar repository: `authinfra.NewUserRepository(collection)`
3. Instanciar service: `authapp.NewAuthService(repo)`
4. Configurar rotas: `authpres.SetupRouter(opts)`

---

### Fase 7: Finalização

**Comandos**:
```bash
make swag  # Gerar Swagger
make mock  # Gerar mocks
make test  # Executar testes
```

---

## 💡 Dicas para Plan Mode

1. **Seja específico**: Liste arquivos exatos, não "criar arquivos necessários"
2. **Siga a ordem**: Domain → Application → Infrastructure → Application → Presentation → Main
3. **Reutilize**: Sempre verificar se Value Objects/Entidades já existem
4. **Foque no contrato**: Assinaturas de métodos, campos de structs, não implementação
5. **Pense em testes**: Cada componente deve ser testável

---

## 🔄 Fluxo de Trabalho

```
User Story
    ↓
Este guia (decomposição)
    ↓
Plan Mode output estruturado
    ↓
Implementação seguindo .cursorrules
    ↓
Feature completa
```
