# Débitos Técnicos - Pinnado Backend

> **Última atualização**: 2026-02-08
> **Arquitetura**: Clean Architecture + DDD + Go 1.25.0

---

## 📊 Resumo Executivo

**Score da Arquitetura**: 8/10

**Fundação**:
- ✅ Clean Architecture bem implementada
- ✅ DDD aplicado corretamente
- ✅ Testabilidade excelente
- ⚠️ Orquestração precária (DI manual)
- ⚠️ Escalabilidade limitada

---

## 🎯 Priorização

| Prioridade | Débito | Impacto | Esforço | Status |
|------------|--------|---------|---------|--------|
| 🔴 P0 | [#1] Dependency Injection Container | Alto | Médio | 🔜 Pendente |
| 🟠 P1 | [#2] Refatorar `shared/config` | Alto | Baixo | 🔜 Pendente |
| 🟡 P2 | [#3] Índices MongoDB no Startup | Médio | Baixo | 🔜 Pendente |
| 🟡 P2 | [#4] Contexto sem Request ID | Médio | Baixo | 🔜 Pendente |
| 🟡 P2 | [#5] Logger Duplicado | Baixo | Muito Baixo | 🔜 Pendente |
| 🟢 P3 | [#6] Ausência de Middleware Chain | Médio | Médio | 🔜 Pendente |
| 🟢 P3 | [#7] Health Check Simplista | Baixo | Baixo | 🔜 Pendente |

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

### #3: Índices MongoDB Síncronos no Startup

**Problema**: `authinfra.CreateIndexes()` é bloqueante no startup da aplicação

**Localização**: `cmd/api/main.go:67-69`

```go
// CÓDIGO ATUAL (bloqueante)
log.Println("creating MongoDB indexes...")
if err := authinfra.CreateIndexes(ctx, mongoClient.Database(config.Mongo.DBName)); err != nil {
    log.Fatalf("failed to create MongoDB indexes: %v", err)
}
```

**Impactos**:
- ❌ Startup lento (com N módulos = N * tempo de índice)
- ❌ Falha de índice = aplicação não sobe (muito restritivo)
- ❌ Dificulta rollback de migração (índice já criado)
- ❌ Impossível versionar índices (sem controle de migração)
- ❌ Containers Kubernetes ficam em `CrashLoopBackOff` se índice falha

**Solução Proposta**:

**Opção 1: Índices Assíncronos (Simples)**
```go
// cmd/api/main.go
go func() {
    log.Println("creating MongoDB indexes in background...")
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := authinfra.CreateIndexes(ctx, db); err != nil {
        log.Printf("WARNING: failed to create indexes: %v", err)
        // notificar monitoramento (Sentry, DataDog, etc)
    } else {
        log.Println("MongoDB indexes created successfully")
    }
}()
```

**Opção 2: Ferramenta de Migração (Recomendado)**
```bash
# Makefile
migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down

# Deployment
# 1. make migrate-up   (antes de subir app)
# 2. make run          (app sobe sem criar índices)
```

```go
// cmd/migrate/main.go
package main

import (
    "context"
    "log"
    "os"
    
    authinfra "pinnado/internal/auth/infrastructure"
    "pinnado/internal/shared/infrastructure"
    "pinnado/pkg/mongodb"
)

func main() {
    if len(os.Args) < 2 {
        log.Fatal("usage: migrate [up|down]")
    }
    
    config := infrastructure.LoadConfig()
    ctx := context.Background()
    
    client, err := mongodb.NewMongoClient(ctx, config.Mongo...)
    if err != nil {
        log.Fatalf("failed to connect: %v", err)
    }
    defer client.Disconnect(ctx)
    
    db := client.Database(config.Mongo.DBName)
    
    switch os.Args[1] {
    case "up":
        log.Println("Running migrations UP...")
        if err := authinfra.CreateIndexes(ctx, db); err != nil {
            log.Fatalf("migration failed: %v", err)
        }
        log.Println("Migrations completed successfully")
    case "down":
        log.Println("Running migrations DOWN...")
        if err := authinfra.DropIndexes(ctx, db); err != nil {
            log.Fatalf("rollback failed: %v", err)
        }
        log.Println("Rollback completed successfully")
    default:
        log.Fatalf("unknown command: %s", os.Args[1])
    }
}
```

```go
// internal/auth/infrastructure/mongo_indexes.go (adicionar)
func DropIndexes(ctx context.Context, db *mongo.Database) error {
    collection := db.Collection(authdomain.UsersCollectionName)
    _, err := collection.Indexes().DropOne(ctx, "unique_email")
    return err
}
```

**Recomendação**: **Opção 2** (migração separada)

**Estimativa**: 1 dia (criar cmd/migrate + testes)

**Benefícios**:
- ✅ Startup rápido (~100ms vs ~2s)
- ✅ Rollback de índices (make migrate-down)
- ✅ Versionamento de migrations (git)
- ✅ Deploy mais seguro (migrar antes de atualizar app)

**Checklist**:
- [ ] Criar `cmd/migrate/main.go`
- [ ] Implementar `DropIndexes()` em cada módulo
- [ ] Adicionar `make migrate-up` e `make migrate-down`
- [ ] Documentar no README.md
- [ ] Atualizar CI/CD para rodar migrations antes do deploy

---

### #4: Contexto sem Request ID (Tracing)

**Problema**: `context.Background()` usado sem enriquecimento de request ID

**Localização**: Handlers HTTP (ex: `internal/auth/presentation/handler.go`)

```go
// CÓDIGO ATUAL (sem rastreamento)
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    // r.Context() não tem request_id
    output, err := h.authService.Register(r.Context(), input)
    // impossível rastrear esta requisição nos logs
}
```

**Impactos**:
- ❌ Impossível rastrear requisições específicas nos logs
- ❌ Dificulta debugging em produção (logs desconectados)
- ❌ Sem correlação entre: handler → service → repository
- ❌ Impossível integrar com APM (DataDog, New Relic)
- ❌ Logs parecem "sopa" (sem agrupamento por request)

**Solução Proposta**:

```go
// pkg/middleware/request_id.go (NOVO)
package middleware

import (
    "context"
    "net/http"
    
    "github.com/google/uuid"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

func RequestID(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := r.Header.Get("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        
        ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
        w.Header().Set("X-Request-ID", requestID)
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func GetRequestID(ctx context.Context) string {
    if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
        return requestID
    }
    return "unknown"
}
```

```go
// pkg/logger/logger.go (ATUALIZAR)
func InfoContext(ctx context.Context, msg string, args ...any) {
    requestID := middleware.GetRequestID(ctx)
    slog.Info(msg, append([]any{"request_id", requestID}, args...)...)
}
```

```go
// internal/auth/presentation/router.go (USAR)
func SetupRouter(opts SetupRouterOptions) {
    handler := authpres.NewAuthHandler(opts.AuthService)
    
    opts.Mux.Handle("/api/auth/register", 
        middleware.RequestID(http.HandlerFunc(handler.Register)))
}
```

```go
// internal/auth/application/auth_service.go (USAR)
func (s *authService) Register(ctx context.Context, input UserInput) (UserOutput, error) {
    logger.InfoContext(ctx, "registering user", "email", input.Email)
    // logs agora tem request_id automaticamente
}
```

**Estimativa**: 1 dia (implementar + adicionar em todos os handlers)

**Benefícios**:
- ✅ Logs rastreáveis (agrupa por request_id)
- ✅ Debugging fácil (buscar por request_id)
- ✅ Integração com APM (correlação de traces)
- ✅ Facilita troubleshooting em produção

**Exemplo de Logs**:
```json
// ANTES (impossível rastrear)
{"level":"info","msg":"registering user","email":"user@example.com"}
{"level":"info","msg":"user created","user_id":"abc123"}

// DEPOIS (rastreável)
{"level":"info","request_id":"550e8400-e29b","msg":"registering user","email":"user@example.com"}
{"level":"info","request_id":"550e8400-e29b","msg":"user created","user_id":"abc123"}
```

**Checklist**:
- [ ] Criar `pkg/middleware/request_id.go`
- [ ] Atualizar `pkg/logger/logger.go` com `InfoContext()`
- [ ] Adicionar middleware em todos os routers
- [ ] Atualizar services para usar `logger.InfoContext()`
- [ ] Documentar uso no `.cursorrules`

---

### #5: Logger Duplicado

**Problema**: Dois loggers criados no `main.go`

**Localização**: `cmd/api/main.go:40,83`

```go
// CÓDIGO ATUAL (duplicado)
slog.SetDefault(logger.NewLogger("info"))  // Linha 40 (global)
// ...
appLogger := logger.NewLogger("info")      // Linha 83 (mesmo logger!)
```

**Impactos**:
- ❌ Confusão sobre qual logger usar (global vs injetado)
- ❌ `slog.SetDefault()` global dificulta testes unitários
- ❌ Inconsistência (alguns lugares usam `slog`, outros `appLogger`)
- ❌ Impossível trocar nível de log em runtime por módulo

**Solução Proposta**:

```go
// cmd/api/main.go (REFATORADO)
func main() {
    // REMOVER: slog.SetDefault(logger.NewLogger("info"))
    
    appLogger := logger.NewLogger("info")  // ← UM ÚNICO LOGGER
    
    // Injetar em TODOS os lugares
    presentation.SetupRouter(presentation.SetupRouterOptions{
        // ...
        Logger: appLogger,
    })
    
    authpres.SetupRouter(authpres.SetupRouterOptions{
        // ...
        Logger: appLogger,
    })
}
```

```go
// internal/auth/application/auth_service.go (INJETAR)
type authService struct {
    repo   UserRepository
    jwt    JWTService
    logger *slog.Logger  // ← adicionar
}

func NewAuthService(repo UserRepository, jwt JWTService, logger *slog.Logger) *authService {
    return &authService{repo: repo, jwt: jwt, logger: logger}
}

func (s *authService) Register(ctx context.Context, input UserInput) (UserOutput, error) {
    s.logger.InfoContext(ctx, "registering user", "email", input.Email)
    // ...
}
```

**Estimativa**: 2 horas (remover SetDefault + injetar em services)

**Benefícios**:
- ✅ Um único ponto de configuração de logger
- ✅ Testabilidade (mockar logger nos testes)
- ✅ Possível ter loggers diferentes por módulo (debug apenas em auth)
- ✅ Consistência (sempre injetado)

**Checklist**:
- [ ] Remover `slog.SetDefault()` do main.go
- [ ] Adicionar `logger *slog.Logger` nos construtores de services
- [ ] Atualizar todos os `NewXXXService()` para receber logger
- [ ] Atualizar testes para passar logger mock
- [ ] Documentar padrão no `.cursorrules`

---

## 🟢 P3 - Baixo Impacto

### #6: Ausência de Middleware Chain

**Problema**: Routers registram rotas sem middleware organizado

**Localização**: `internal/auth/presentation/router.go`, `internal/shared/presentation/router.go`

```go
// CÓDIGO ATUAL (sem middleware chain)
func SetupRouter(opts SetupRouterOptions) {
    handler := NewAuthHandler(opts.AuthService)
    
    // cada handler é registrado "cru"
    opts.Mux.HandleFunc("POST /api/auth/register", handler.Register)
    opts.Mux.HandleFunc("POST /api/auth/login", handler.Login)
}
```

**Impactos**:
- ❌ Difícil adicionar: CORS, Rate Limiting, Auth, Recovery
- ❌ Código duplicado (cada handler precisa logar manualmente)
- ❌ Impossível aplicar middleware globalmente
- ❌ Ordem de middleware não clara

**Solução Proposta**:

```go
// pkg/middleware/chain.go (NOVO)
package middleware

import "net/http"

type Chain struct {
    middlewares []func(http.Handler) http.Handler
}

func NewChain(middlewares ...func(http.Handler) http.Handler) *Chain {
    return &Chain{middlewares: middlewares}
}

func (c *Chain) Then(h http.Handler) http.Handler {
    for i := len(c.middlewares) - 1; i >= 0; i-- {
        h = c.middlewares[i](h)
    }
    return h
}

func (c *Chain) ThenFunc(fn http.HandlerFunc) http.Handler {
    return c.Then(fn)
}
```

```go
// pkg/middleware/logging.go (NOVO)
package middleware

import (
    "log/slog"
    "net/http"
    "time"
)

func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            next.ServeHTTP(w, r)
            
            logger.InfoContext(r.Context(), "request completed",
                "method", r.Method,
                "path", r.URL.Path,
                "duration_ms", time.Since(start).Milliseconds())
        })
    }
}
```

```go
// pkg/middleware/recovery.go (NOVO)
package middleware

import (
    "log/slog"
    "net/http"
)

func Recovery(logger *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            defer func() {
                if err := recover(); err != nil {
                    logger.ErrorContext(r.Context(), "panic recovered",
                        "error", err,
                        "path", r.URL.Path)
                    
                    w.WriteHeader(http.StatusInternalServerError)
                    w.Write([]byte(`{"message":"internal server error"}`))
                }
            }()
            
            next.ServeHTTP(w, r)
        })
    }
}
```

```go
// internal/auth/presentation/router.go (USAR)
func SetupRouter(opts SetupRouterOptions) {
    handler := NewAuthHandler(opts.AuthService)
    
    // chain padrão para todas as rotas
    chain := middleware.NewChain(
        middleware.Recovery(opts.Logger),
        middleware.RequestID,
        middleware.Logging(opts.Logger),
    )
    
    // rotas públicas
    opts.Mux.Handle("POST /api/auth/register", chain.ThenFunc(handler.Register))
    opts.Mux.Handle("POST /api/auth/login", chain.ThenFunc(handler.Login))
    
    // rotas autenticadas (futuro)
    authChain := middleware.NewChain(
        middleware.Recovery(opts.Logger),
        middleware.RequestID,
        middleware.Logging(opts.Logger),
        middleware.Auth(opts.JWTService),  // ← adiciona auth
    )
    
    opts.Mux.Handle("GET /api/auth/me", authChain.ThenFunc(handler.Me))
}
```

**Estimativa**: 1 dia (implementar chain + middlewares básicos)

**Benefícios**:
- ✅ Middleware reutilizável e componível
- ✅ Ordem clara (Recovery → RequestID → Logging → Auth)
- ✅ Fácil adicionar CORS, Rate Limiting, etc
- ✅ Código limpo (sem repetição)

**Middlewares Sugeridos**:
- ✅ Recovery (panic handling)
- ✅ RequestID (tracing)
- ✅ Logging (request/response)
- 🔜 CORS (cross-origin)
- 🔜 RateLimit (proteção contra abuse)
- 🔜 Auth (JWT validation)
- 🔜 Metrics (Prometheus)

**Checklist**:
- [ ] Criar `pkg/middleware/chain.go`
- [ ] Criar `pkg/middleware/recovery.go`
- [ ] Criar `pkg/middleware/logging.go`
- [ ] Refatorar routers para usar chain
- [ ] Adicionar testes de middleware
- [ ] Documentar no `.cursorrules`

---

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

## 📈 Roadmap de Implementação

### **Fase 1: Fundação** (Semana 1-2)
- [x] Análise de débitos técnicos
- [ ] #5: Unificar logger (2h)
- [ ] #2: Refatorar shared/config (1 dia)
- [ ] #4: Adicionar Request ID (1 dia)

**Resultado**: Logs rastreáveis + módulos desacoplados

---

### **Fase 2: Infraestrutura** (Semana 3-4)
- [ ] #3: Migração de índices (1 dia)
- [ ] #6: Middleware chain (1 dia)
- [ ] #7: Health check robusto (1 dia)

**Resultado**: Deploy mais seguro + observabilidade

---

### **Fase 3: Escalabilidade** (Semana 5-6)
- [ ] #1: Dependency Injection (Wire) (3 dias)
- [ ] Integração com APM (DataDog/New Relic) (2 dias)
- [ ] Métricas Prometheus (1 dia)

**Resultado**: Arquitetura escalável + monitoramento

---

## 🎯 Critérios de Sucesso

Quando este documento estiver **100% implementado**:

**Código**:
- ✅ `main.go` < 50 linhas (vs 127 atual)
- ✅ Startup < 100ms (vs ~2s atual)
- ✅ Logs rastreáveis com request_id
- ✅ Módulos 100% desacoplados

**Deploy**:
- ✅ Zero-downtime (Kubernetes readiness)
- ✅ Rollback de migrations (make migrate-down)
- ✅ CI/CD automatizado (migrations + deploy)

**Monitoramento**:
- ✅ APM integrado (traces correlacionados)
- ✅ Métricas Prometheus (latência, erros, throughput)
- ✅ Health check detalhado (status por dependência)

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
- **Responsável**: Equipe de Arquitetura
- **Revisão**: Mensal (ou quando adicionar novo módulo)
- **Arquivo vivo**: Sempre atualizar quando resolver débito

----
