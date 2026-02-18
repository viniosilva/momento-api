# Débitos Técnicos - Pinnado Backend

> **Última atualização**: 2026-02-18
> **Arquitetura**: Clean Architecture + DDD + Go 1.25.0

---

## 📊 Resumo Executivo

**Score da Arquitetura**: 8.5/10

**Fundação**:
- ✅ Clean Architecture bem implementada
- ✅ DDD aplicado corretamente
- ✅ Testabilidade excelente
- ✅ Migração de índices desacoplada do startup
- ✅ Middleware centralizado em `pkg/nethttp` (Chain, Request ID, Timeout, Recovery, Auth, Logging)
- ✅ `pkg/` agnóstico de domínio (sem importar `internal/`)
- ⚠️ Orquestração precária (DI manual)
- ⚠️ Escalabilidade limitada

---

## 🎯 Priorização

| Prioridade | Débito | Impacto | Esforço | Status |
|------------|--------|---------|---------|--------|
| 🔴 P0 | [#1] Dependency Injection Container | Alto | Médio | 🔜 Pendente |
| 🟠 P1 | [#2] Refatorar `shared/config` | Alto | Baixo | 🔜 Pendente |
| 🟡 P2 | [#13] Falta validação de sort fields | Baixo | Baixo | 🔜 Pendente |
| 🟢 P3 | [#7] Health Check Simplista | Baixo | Baixo | 🔜 Pendente |
| 🟢 P3 | [#17] Falta validação de env vars | Médio | Baixo | 🔜 Pendente |

---

## 🔴 P0 - Crítico

### #1: Dependency Injection Container

**Problema**: Todas as dependências são instanciadas manualmente no `cmd/api/main.go`

**Localização**: `cmd/api/main.go:72-78`

```go
// CÓDIGO ATUAL (problemático)
userRepository := authinfra.NewUserRepository(userCollection)
jwtService := authinfra.NewJWTService(config.JWT.Secret, config.JWT.Expiration)
authService := authapp.NewAuthService(userRepository, jwtService)
```

**Impactos**:
- ❌ Crescimento linear com novos módulos (exemplo: adicionar 10 módulos = +100 linhas no main)
- ❌ Viola Single Responsibility (main conhece implementação de TODOS os módulos)
- ❌ Dificulta testes de integração (impossível mockar setup)
- ❌ Impossível trocar implementações em runtime
- ❌ Manutenção cara (cada novo módulo = alterar main.go)

**Solução Proposta**:

**Opção 1: Wire (Google) - Recomendado**
```go
// internal/di/wire.go
//go:build wireinject
// +build wireinject

package di

import (
    "github.com/google/wire"
    "go.mongodb.org/mongo-driver/mongo"
    
    authapp "pinnado/internal/auth/application"
    authinfra "pinnado/internal/auth/infrastructure"
    "pinnado/internal/shared/infrastructure"
)

func InitializeAuthModule(db *mongo.Database, config infrastructure.Config) authapp.AuthService {
    wire.Build(
        authinfra.NewUserRepository,
        authinfra.NewJWTService,
        authapp.NewAuthService,
        wire.Bind(new(authapp.UserRepository), new(*authinfra.UserRepository)),
    )
    return nil
}

// cmd/api/main.go
authService := di.InitializeAuthModule(db, config)
```

**Opção 2: Dig (Uber)**
```go
// internal/di/container.go
package di

import (
    "go.uber.org/dig"
    "go.mongodb.org/mongo-driver/mongo"
)

type Container struct {
    *dig.Container
}

func NewContainer() *Container {
    return &Container{dig.New()}
}

func (c *Container) RegisterAuth(db *mongo.Database, config Config) error {
    c.Provide(authinfra.NewUserRepository)
    c.Provide(authinfra.NewJWTService)
    c.Provide(authapp.NewAuthService)
    return nil
}

// cmd/api/main.go
container := di.NewContainer()
container.RegisterAuth(db, config)
var authService authapp.AuthService
container.Invoke(func(s authapp.AuthService) { authService = s })
```

**Opção 3: Factory Pattern (Manual)**
```go
// internal/di/factory.go
package di

type Factory struct {
    db     *mongo.Database
    config Config
}

func NewFactory(db *mongo.Database, config Config) *Factory {
    return &Factory{db: db, config: config}
}

func (f *Factory) CreateAuthModule() authapp.AuthService {
    userCollection := f.db.Collection(authdomain.UsersCollectionName)
    userRepo := authinfra.NewUserRepository(userCollection)
    jwtService := authinfra.NewJWTService(f.config.JWT.Secret, f.config.JWT.Expiration)
    return authapp.NewAuthService(userRepo, jwtService)
}

// cmd/api/main.go
factory := di.NewFactory(db, config)
authService := factory.CreateAuthModule()
healthService := factory.CreateHealthModule()
```

**Recomendação**: **Wire** (compile-time safety + performance)

**Estimativa**: 2-3 dias (setup inicial + refactor de 2 módulos)

**Benefícios**:
- ✅ Reduz `main.go` de 127 para ~50 linhas
- ✅ Facilita adicionar novos módulos (sem mexer no main)
- ✅ Testabilidade (mockar container inteiro)
- ✅ Troca de implementações sem tocar no main

**Referências**:
- [Wire Tutorial](https://github.com/google/wire/blob/main/docs/guide.md)
- [Dig Tutorial](https://github.com/uber-go/dig)

---

## 🟠 P1 - Alto Impacto

### #2: Módulo `shared` com Config Monolítico

**Problema**: `internal/shared/infrastructure/config.go` conhece configuração de TODOS os módulos

**Localização**: `internal/shared/infrastructure/config.go:12-37`

```go
// CÓDIGO ATUAL (problemático)
type Config struct {
    Api   ApiConfig
    Mongo MongoConfig
    JWT   JWTConfig  // ← JWT é específico do módulo Auth!
}
```

**Impactos**:
- ❌ `shared` vira "God Module" (acoplamento alto)
- ❌ Adicionar novo módulo = quebrar interface de `shared`
- ❌ Impossível extrair `auth` para microsserviço (JWT acoplado)
- ❌ Viola Single Responsibility (Config conhece API + DB + Auth + ...)
- ❌ Testes de módulos ficam interdependentes

**Solução Proposta**:

**Estrutura Nova**:
```
/internal
  /shared
    /infrastructure
      /config.go  ← APENAS config comum (API + Mongo connection)
  /auth
    /infrastructure
      /config.go  ← Config específico do Auth (JWT + User collection)
  /products  # (futuro)
    /infrastructure
      /config.go  ← Config específico de Products
```

**Implementação**:

```go
// internal/shared/infrastructure/config.go (MINIMALISTA)
package infrastructure

type Config struct {
    Api   ApiConfig
    Mongo MongoConfig
}

type ApiConfig struct {
    Host string
    Port string
}

type MongoConfig struct {
    Host           string
    Port           string
    DBName         string
    User           string
    Pass           string
    MaxRetries     int
    RetryDelay     time.Duration
    ConnectTimeout time.Duration
}

func LoadConfig(opts ...LoadConfigOption) Config {
    // apenas config comum
}
```

```go
// internal/auth/infrastructure/config.go (NOVO)
package infrastructure

type AuthConfig struct {
    JWT JWTConfig
}

type JWTConfig struct {
    Secret     string
    Expiration time.Duration
}

func LoadAuthConfig() AuthConfig {
    return AuthConfig{
        JWT: JWTConfig{
            Secret:     getEnv("JWT_SECRET", "your-secret-key"),
            Expiration: getEnvAsDuration("JWT_EXPIRATION_MS", 24*time.Hour),
        },
    }
}
```

```go
// cmd/api/main.go (REFATORADO)
baseConfig := sharedinfra.LoadConfig()
authConfig := authinfra.LoadAuthConfig()

// uso
jwtService := authinfra.NewJWTService(authConfig.JWT.Secret, authConfig.JWT.Expiration)
```

**Estimativa**: 1 dia (refactor + testes)

**Benefícios**:
- ✅ `shared` volta a ser agnóstico
- ✅ Módulos desacoplados (podem virar microsserviços)
- ✅ Testes isolados por módulo
- ✅ Facilita remover módulos (sem quebrar shared)

**Checklist**:
- [ ] Mover `JWTConfig` para `auth/infrastructure/config.go`
- [ ] Criar `LoadAuthConfig()` no módulo auth
- [ ] Atualizar `cmd/api/main.go` para usar 2 configs
- [ ] Atualizar testes de `shared/infrastructure/config_test.go`
- [ ] Atualizar `.cursorrules` com novo padrão

---

## 🟡 P2 - Médio Impacto

_(Seção #13 abaixo)_

---

## 🟢 P3 - Baixo Impacto

### #7: Health Check Simplista

**Problema**: Health apenas pinga MongoDB

**Localização**: `internal/shared/application/health_service.go`

```go
// CÓDIGO ATUAL (simplista)
func (s *healthService) Check(ctx context.Context) (Health, error) {
    if err := s.mongoClient.Ping(ctx); err != nil {
        return Health{Status: "DOWN"}, err
    }
    return Health{Status: "UP"}, nil
}
```

**Impactos**:
- ❌ Não verifica: Redis (futuro), APIs externas, filas (RabbitMQ, Kafka)
- ❌ Kubernetes precisa de `/health/live` e `/health/ready` (semântica diferente)
- ❌ Impossível deploy zero-downtime (não sabe se app está pronto)
- ❌ Sem detalhes de qual dependência falhou

**Solução Proposta**:

```go
// internal/shared/domain/health.go (ATUALIZAR)
package domain

type Health struct {
    Status string           `json:"status"`  // UP, DOWN, DEGRADED
    Checks map[string]Check `json:"checks"`
}

type Check struct {
    Status  string `json:"status"`  // UP, DOWN
    Message string `json:"message,omitempty"`
}

const (
    StatusUp       = "UP"
    StatusDown     = "DOWN"
    StatusDegraded = "DEGRADED"  // algumas dependências falharam
)
```

```go
// internal/shared/application/health_service.go (REFATORAR)
package application

type healthService struct {
    mongoClient MongoClient
    // adicionar outras dependências (Redis, HTTP clients, etc)
}

func (s *healthService) Liveness(ctx context.Context) (domain.Health, error) {
    // SEMPRE retorna UP (processo está vivo)
    return domain.Health{
        Status: domain.StatusUp,
        Checks: map[string]domain.Check{
            "app": {Status: domain.StatusUp, Message: "application is running"},
        },
    }, nil
}

func (s *healthService) Readiness(ctx context.Context) (domain.Health, error) {
    checks := make(map[string]domain.Check)
    overallStatus := domain.StatusUp
    
    // Check MongoDB
    if err := s.mongoClient.Ping(ctx); err != nil {
        checks["mongodb"] = domain.Check{
            Status:  domain.StatusDown,
            Message: err.Error(),
        }
        overallStatus = domain.StatusDown
    } else {
        checks["mongodb"] = domain.Check{Status: domain.StatusUp}
    }
    
    // Adicionar outros checks (Redis, APIs, etc)
    // if err := s.redisClient.Ping(ctx); err != nil { ... }
    
    return domain.Health{
        Status: overallStatus,
        Checks: checks,
    }, nil
}
```

```go
// internal/shared/presentation/handler.go (ATUALIZAR)
func (h *HealthHandler) Liveness(w http.ResponseWriter, r *http.Request) {
    health, _ := h.healthService.Liveness(r.Context())
    nethttp.JSON(w, http.StatusOK, health)
}

func (h *HealthHandler) Readiness(w http.ResponseWriter, r *http.Request) {
    health, _ := h.healthService.Readiness(r.Context())
    
    statusCode := http.StatusOK
    if health.Status == domain.StatusDown {
        statusCode = http.StatusServiceUnavailable
    }
    
    nethttp.JSON(w, statusCode, health)
}
```

```go
// internal/shared/presentation/router.go (ATUALIZAR)
func SetupRouter(opts SetupRouterOptions) {
    handler := NewHealthHandler(opts.HealthService)
    
    opts.Mux.HandleFunc("GET /api/health/live", handler.Liveness)
    opts.Mux.HandleFunc("GET /api/health/ready", handler.Readiness)
}
```

**Kubernetes Deployment**:
```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
      - name: pinnado-api
        livenessProbe:
          httpGet:
            path: /api/health/live
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /api/health/ready
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 5
```

**Estimativa**: 1 dia (implementar readiness/liveness + testes)

**Benefícios**:
- ✅ Deploy zero-downtime (Kubernetes espera readiness)
- ✅ Observabilidade (qual dependência falhou)
- ✅ Status degradado (app funciona parcialmente)
- ✅ Integração com Kubernetes/Docker

**Exemplo de Response**:
```json
// GET /api/health/live (sempre 200)
{
  "status": "UP",
  "checks": {
    "app": {"status": "UP", "message": "application is running"}
  }
}

// GET /api/health/ready (200 se OK, 503 se falha)
{
  "status": "DEGRADED",
  "checks": {
    "mongodb": {"status": "UP"},
    "redis": {"status": "DOWN", "message": "connection refused"}
  }
}
```

**Checklist**:
- [ ] Adicionar constantes de status no domain
- [ ] Implementar `Liveness()` e `Readiness()` no service
- [ ] Adicionar rotas `/health/live` e `/health/ready`
- [ ] Atualizar Swagger
- [ ] Atualizar deployment.yaml (Kubernetes)
- [ ] Documentar diferença entre liveness e readiness no README

---

### #13: Falta de Validação de Sort Fields

**Problema**: `SortDTO` aceita qualquer campo sem validação

**Localização**: `internal/shared/application/dto/sort_dto.go`

```go
// CÓDIGO ATUAL (sem validação)
type SortDTO struct {
    Field string `json:"field"`
    Order string `json:"order"` // "asc" ou "desc"
}
```

**Impactos**:
- ❌ **Possível NoSQL injection**: Campo malicioso pode acessar dados sensíveis
- ❌ **Erros em runtime**: MongoDB retorna erro se campo não existe
- ❌ **Experiência ruim**: Cliente não sabe quais campos são válidos
- ❌ **Sem documentação**: Swagger não lista campos permitidos

**Solução Proposta**:

```go
// internal/shared/application/dto/sort_dto.go (REFATORADO)
package dto

import "fmt"

type SortDTO struct {
    Field string `json:"field" example:"created_at"`
    Order string `json:"order" example:"desc" enums:"asc,desc"`
}

func (s SortDTO) Validate(allowedFields []string) error {
    // Valida order
    if s.Order != "asc" && s.Order != "desc" {
        return fmt.Errorf("invalid sort order: %s (allowed: asc, desc)", s.Order)
    }
    
    // Valida field
    if !contains(allowedFields, s.Field) {
        return fmt.Errorf("invalid sort field: %s (allowed: %v)", s.Field, allowedFields)
    }
    
    return nil
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}
```

```go
// internal/notes/presentation/handler.go (USAR)
package presentation

var allowedSortFields = []string{"created_at", "updated_at", "title"}

func (h *NoteHandler) ListNotes(w http.ResponseWriter, r *http.Request) {
    // Parse query params
    sortField := r.URL.Query().Get("sort_field")
    sortOrder := r.URL.Query().Get("sort_order")
    
    sort := dto.SortDTO{Field: sortField, Order: sortOrder}
    
    // Validar campos permitidos
    if err := sort.Validate(allowedSortFields); err != nil {
        nethttp.JSON(w, http.StatusBadRequest, ErrorResponse{Message: err.Error()})
        return
    }
    
    // Continuar com listagem
    // ...
}
```

**Estimativa**: 1 hora

**Benefícios**:
- ✅ Segurança contra NoSQL injection
- ✅ Validação em tempo de request
- ✅ Mensagens claras de erro
- ✅ Documentação implícita (campos permitidos)

**Checklist**:
- [ ] Adicionar `Validate()` em `SortDTO`
- [ ] Definir `allowedSortFields` em cada handler que usa sort
- [ ] Adicionar validação antes de chamar service
- [ ] Atualizar Swagger com campos permitidos
- [ ] Adicionar testes de validação

---

### #17: Falta de Validação de Environment Variables

**Problema**: `LoadConfig()` usa defaults sem validar obrigatoriedade

**Localização**: `internal/shared/infrastructure/config.go`

**Impactos**:
- ❌ **App sobe com config inválida**: Ex: JWT secret vazio
- ❌ **Bugs em runtime**: Conexão MongoDB falha depois de 30s
- ❌ **Difícil debug**: Erro aparece longe da causa raiz
- ❌ **Insegurança**: Pode usar JWT secret padrão em produção

**Solução Proposta**:

```go
// internal/shared/infrastructure/config.go (ADICIONAR VALIDAÇÃO)
package infrastructure

import (
    "fmt"
    "log"
    "os"
)

func LoadConfig(opts ...LoadConfigOption) Config {
    config := Config{
        Api: ApiConfig{
            Host: getEnv("API_HOST", "0.0.0.0"),
            Port: getEnv("API_PORT", "8080"),
        },
        Mongo: MongoConfig{
            Host:   getEnv("MONGO_HOST", "localhost"),
            Port:   getEnv("MONGO_PORT", "27017"),
            DBName: getEnv("MONGO_DB_NAME", "pinnado"),
            User:   getEnv("MONGO_USER", ""),
            Pass:   getEnv("MONGO_PASS", ""),
            // ...
        },
    }
    
    // VALIDAR configurações obrigatórias
    if err := config.Validate(); err != nil {
        log.Fatalf("invalid configuration: %v", err)
    }
    
    return config
}

// Validate verifica se configuração é válida
func (c Config) Validate() error {
    // Validar MongoDB
    if c.Mongo.DBName == "" {
        return fmt.Errorf("MONGO_DB_NAME is required")
    }
    
    // Validar API
    if c.Api.Port == "" {
        return fmt.Errorf("API_PORT is required")
    }
    
    // Validar em produção
    if os.Getenv("ENV") == "production" {
        if c.Mongo.User == "" || c.Mongo.Pass == "" {
            return fmt.Errorf("MONGO_USER and MONGO_PASS are required in production")
        }
    }
    
    return nil
}
```

```go
// internal/auth/infrastructure/config.go (VALIDAR JWT)
func LoadAuthConfig() AuthConfig {
    config := AuthConfig{
        JWT: JWTConfig{
            Secret:     getEnv("JWT_SECRET", ""),
            Expiration: getEnvAsDuration("JWT_EXPIRATION_MS", 24*time.Hour),
        },
    }
    
    if err := config.Validate(); err != nil {
        log.Fatalf("invalid auth configuration: %v", err)
    }
    
    return config
}

func (c AuthConfig) Validate() error {
    if c.JWT.Secret == "" {
        return fmt.Errorf("JWT_SECRET is required")
    }
    
    if len(c.JWT.Secret) < 32 {
        return fmt.Errorf("JWT_SECRET must be at least 32 characters")
    }
    
    if c.JWT.Expiration < time.Minute {
        return fmt.Errorf("JWT_EXPIRATION must be at least 1 minute")
    }
    
    return nil
}
```

**Estimativa**: 2 horas

**Benefícios**:
- ✅ Fail-fast no startup
- ✅ Erros claros sobre config faltando
- ✅ Segurança (força JWT secret forte)
- ✅ Documentação implícita (quais configs são obrigatórias)

**Checklist**:
- [ ] Adicionar `Validate()` em `shared/infrastructure/config.go`
- [ ] Adicionar `Validate()` em `auth/infrastructure/config.go`
- [ ] Validar JWT secret tem > 32 caracteres
- [ ] Validar MongoDB credentials em produção
- [ ] Documentar env vars obrigatórias no README
- [ ] Adicionar exemplo de `.env.example`

---

## 📈 Roadmap de Implementação

### **Fase 1: Correções Críticas** (Semana 1)
- [x] Análise de débitos técnicos

**Resultado**: Arquitetura consistente + código limpo

---

### **Fase 2: Qualidade e Padronização** (Semana 2)
- [ ] #13: Validação de sort fields (1h)
- [ ] #2: Refatorar `shared/config` (1 dia)
- [ ] #17: Validar env vars obrigatórias (2h)

**Resultado**: Código padronizado + melhor DX

---

### **Fase 3: Observabilidade** (Semana 3)
- [ ] #7: Health check robusto (1 dia)

**Resultado**: Logs rastreáveis + observabilidade

---

### **Fase 4: Escalabilidade** (Semana 4-5)
- [ ] #1: Dependency Injection (Wire) (3 dias)
- [ ] Integração com APM (DataDog/New Relic) (2 dias)
- [ ] Métricas Prometheus (1 dia)

**Resultado**: Arquitetura escalável + monitoramento

---

## 🎯 Critérios de Sucesso

Quando este documento estiver **100% implementado**:

**Código**:
- ✅ `main.go` < 50 linhas (vs 150 atual)
- ✅ Startup < 100ms (vs ~2s atual)
- ✅ Logs rastreáveis com request_id
- ✅ Módulos 100% desacoplados
- ✅ Zero violações de arquitetura (`pkg/` agnóstico)
- ✅ Zero duplicação de código (middleware centralizado)
- ✅ Logging 100% estruturado (slog)

**Deploy**:
- ✅ Zero-downtime (Kubernetes readiness)
- ✅ Rollback de migrations (make migrate-down)
- ✅ CI/CD automatizado (migrations + deploy)
- ✅ Validação de env vars no startup

**Testes**:
- ✅ Cobertura unitária > 80%

**Monitoramento**:
- ✅ APM integrado (traces correlacionados)
- ✅ Métricas Prometheus (latência, erros, throughput)
- ✅ Health check detalhado (status por dependência)
- ✅ Request timeout configurável

---

## 📚 Referências

**Dependency Injection**:
- [Wire - Google](https://github.com/google/wire)
- [Dig - Uber](https://github.com/uber-go/dig)

**Observabilidade**:
- [Request ID Pattern](https://www.oreilly.com/library/view/cloud-native-go/9781492076322/)
- [Health Check API](https://microservices.io/patterns/observability/health-check-api.html)

**Migrations**:
- [Database Migrations Best Practices](https://www.prisma.io/dataguide/types/relational/migration-strategies)

**Middleware**:
- [Go Middleware Patterns](https://drstearns.github.io/tutorials/gomiddleware/)

---

## 📝 Notas

- **Data de criação**: 2026-02-08
- **Última atualização**: 2026-02-18
- **Responsável**: Equipe de Arquitetura
- **Revisão**: Mensal (ou quando adicionar novo módulo)
- **Arquivo vivo**: Sempre atualizar quando resolver débito

---

## 📋 Changelog

### 2026-02-18 - Remoção dos débitos #15, #16 e #18
- **#15** (Falta de testes de integração), **#16** (Makefile limitado) e **#18** (Estrutura de erros simples) removidos do documento por decisão de escopo.
- Tabela de priorização, seções e roadmap reescritos; Fase 4 (Infraestrutura) eliminada; Fase 5 renumerada para Fase 4.
- Critérios de sucesso atualizados (removida menção a testes de integração/E2E).

### 2026-02-18 - Débitos #11 e #14 resolvidos e removidos
- **#11** (User.Update não utilizado): Removido do documento (código já não continha o método; débito obsoleto).
- **#14** (Swagger Host desnecessária): Resolvido — removida a linha `docs.SwaggerInfo.Host = addr` em `cmd/api/main.go`. Swagger UI passa a inferir o host automaticamente.
- Tabela de priorização, seções e roadmap atualizados.

### 2026-02-18 - Débitos de log resolvidos e removidos
- Resolvidos e removidos do documento os débitos de logging: #5 (Logger duplicado), #12 (Logging inconsistente no main.go), #20 (Graceful shutdown sem logging)
- Implementado em `cmd/api/main.go`: um único logger, logging 100% estruturado (appLogger), shutdown com sinal, timeout e elapsed
- Tabela de priorização, seções P2/P3 e roadmap atualizados

### 2026-02-18 - Remoção de débitos resolvidos
- Removidos 7 débitos já implementados: #3 (migração de índices), #4 (Request ID), #6 (Middleware Chain), #8 (pkg → internal), #9 (middleware duplicado), #10 (inconsistência middleware), #19 (context timeout)
- Soluções aplicadas via `pkg/nethttp` (Chain, requestid, timeout, recovery, auth, logging) e `cmd/migrate`
- Tabela de priorização e roadmap atualizados apenas com itens pendentes

### 2026-02-15 - Auditoria de Estrutura
- Adicionados 13 novos débitos técnicos (#8 a #20)
- Identificada violação crítica: `pkg/` importando `internal/` (#8)
- Identificada duplicação de código em middlewares (#9, #10)
- Adicionados débitos de qualidade de código (#11, #12, #13, #14)
- Adicionados débitos de infraestrutura (#15, #16, #17)
- Adicionados débitos de melhorias opcionais (#18, #19, #20)
- Reorganizado roadmap em 5 fases
- Atualizado score da arquitetura (mantém 8/10)

### 2026-02-08 - Criação Inicial
- Documento criado com 7 débitos técnicos originais
- Definidas prioridades P0 a P3
- Estruturado roadmap de 3 fases

----
