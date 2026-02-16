# Débitos Técnicos - Pinnado Backend

> **Última atualização**: 2026-02-15
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
- ⚠️ Violação de arquitetura (pkg → internal)

---

## 🎯 Priorização

| Prioridade | Débito | Impacto | Esforço | Status |
|------------|--------|---------|---------|--------|
| 🔴 P0 | [#1] Dependency Injection Container | Alto | Médio | 🔜 Pendente |
| 🔴 P0 | [#8] `pkg/` importando `internal/` | Alto | Médio | 🔜 Pendente |
| 🟠 P1 | [#2] Refatorar `shared/config` | Alto | Baixo | 🔜 Pendente |
| 🟠 P1 | [#9] Middleware duplicado em routers | Médio | Baixo | 🔜 Pendente |
| 🟠 P1 | [#10] Inconsistência em middleware | Baixo | Muito Baixo | 🔜 Pendente |
| 🟡 P2 | [#3] Índices MongoDB no Startup | Médio | Baixo | 🔜 Pendente |
| 🟡 P2 | [#4] Contexto sem Request ID | Médio | Baixo | 🔜 Pendente |
| 🟡 P2 | [#5] Logger Duplicado | Baixo | Muito Baixo | 🔜 Pendente |
| 🟡 P2 | [#11] Método `User.Update()` não utilizado | Baixo | Muito Baixo | 🔜 Pendente |
| 🟡 P2 | [#12] Logging inconsistente no `main.go` | Baixo | Muito Baixo | 🔜 Pendente |
| 🟡 P2 | [#13] Falta validação de sort fields | Baixo | Baixo | 🔜 Pendente |
| 🟢 P3 | [#6] Ausência de Middleware Chain | Médio | Médio | 🔜 Pendente |
| 🟢 P3 | [#7] Health Check Simplista | Baixo | Baixo | 🔜 Pendente |
| 🟢 P3 | [#14] Configuração Swagger Host desnecessária | Baixo | Muito Baixo | 🔜 Pendente |
| 🟢 P3 | [#15] Falta de testes de integração | Médio | Alto | 🔜 Pendente |
| 🟢 P3 | [#16] Makefile limitado | Baixo | Baixo | 🔜 Pendente |
| 🟢 P3 | [#17] Falta validação de env vars | Médio | Baixo | 🔜 Pendente |
| 🟢 P3 | [#18] Estrutura de erros simples | Baixo | Médio | 🔜 Pendente |
| 🟢 P3 | [#19] Falta context timeout em handlers | Médio | Baixo | 🔜 Pendente |
| 🟢 P3 | [#20] Graceful shutdown sem logging | Baixo | Muito Baixo | 🔜 Pendente |

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

### #8: `pkg/` importando `internal/` (Violação de Arquitetura)

**Problema**: `pkg/nethttp/port.go` importa `pinnado/internal/auth/infrastructure`

**Localização**: 
- `pkg/nethttp/port.go:3`
- `internal/notes/presentation/router.go:8,16-18`

```go
// pkg/nethttp/port.go (VIOLAÇÃO CRÍTICA)
package nethttp

import "pinnado/internal/auth/infrastructure"

type JWTService interface {
    Validate(tokenString string) (*infrastructure.Claims, error)
}
```

**Impactos**:
- ❌ **Quebra Clean Architecture**: `pkg/` deve ser agnóstico de domínio
- ❌ **`pkg/` não é reutilizável**: Acoplado ao módulo auth
- ❌ **Impossível extrair auth**: `pkg/` depende de `internal/auth`
- ❌ **Viola regra fundamental**: `pkg/` NÃO pode conhecer `internal/`
- ❌ **Testabilidade comprometida**: `pkg/` precisa importar mocks de `internal/`

**Solução Proposta**:

**Opção 1: Mover `Claims` para Domain (Recomendado)**
```go
// internal/auth/domain/claims.go (NOVO)
package domain

import "time"

type Claims struct {
    UserID string    `json:"user_id"`
    Email  Email     `json:"email"`
    Exp    time.Time `json:"exp"`
}
```

```go
// pkg/nethttp/port.go (REFATORADO)
package nethttp

// Remover import de internal/!
// Interface genérica sem tipos concretos
type JWTService interface {
    Validate(tokenString string) (map[string]any, error)
}
```

```go
// pkg/nethttp/auth_middleware.go (REFATORADO)
package nethttp

import (
    "context"
    "net/http"
    "strings"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func AuthMiddleware(jwtService JWTService) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            token := extractToken(r)
            if token == "" {
                JSON(w, http.StatusUnauthorized, ErrorResponse{Message: "missing token"})
                return
            }
            
            claims, err := jwtService.Validate(token)
            if err != nil {
                JSON(w, http.StatusUnauthorized, ErrorResponse{Message: "invalid token"})
                return
            }
            
            // Extrai apenas o user_id (não depende de Claims struct)
            userID, ok := claims["user_id"].(string)
            if !ok {
                JSON(w, http.StatusUnauthorized, ErrorResponse{Message: "invalid claims"})
                return
            }
            
            ctx := context.WithValue(r.Context(), UserIDKey, userID)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

```go
// internal/auth/infrastructure/jwt_service.go (ATUALIZAR)
package infrastructure

func (s *jwtService) Validate(tokenString string) (map[string]any, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
        return []byte(s.secret), nil
    })
    
    if err != nil || !token.Valid {
        return nil, domain.ErrInvalidToken
    }
    
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return nil, domain.ErrInvalidClaims
    }
    
    // Retorna map genérico
    return map[string]any(claims), nil
}
```

**Opção 2: Criar `pkg/jwt` (Alternativa)**
```go
// pkg/jwt/claims.go (NOVO - agnóstico)
package jwt

import "time"

type Claims struct {
    UserID string                 `json:"user_id"`
    Email  string                 `json:"email"`
    Exp    time.Time              `json:"exp"`
    Extra  map[string]interface{} `json:"extra,omitempty"`
}
```

**Opção 3: Remover interface de `pkg/` (Mais Simples)**
```go
// Deletar: pkg/nethttp/port.go

// internal/notes/presentation/port.go (NOVO)
package presentation

import "pinnado/internal/auth/infrastructure"

type JWTService interface {
    Validate(tokenString string) (*infrastructure.Claims, error)
}
```

**Recomendação**: **Opção 3** (mais simples e mantém separação)

**Estimativa**: 3-4 horas

**Benefícios**:
- ✅ `pkg/` volta a ser agnóstico de domínio
- ✅ Reutilizável entre projetos
- ✅ Testável sem mocks de `internal/`
- ✅ Respeita Clean Architecture

**Checklist**:
- [ ] Deletar `pkg/nethttp/port.go`
- [ ] Mover interface `JWTService` para `internal/notes/presentation/port.go`
- [ ] Atualizar `pkg/nethttp/auth_middleware.go` para usar interface genérica ou injetada
- [ ] Atualizar `.cursorrules` (reforçar: `pkg/` NÃO importa `internal/`)
- [ ] Rodar `make test` para validar

---

### #9: Middleware Duplicado em 3 Routers

**Problema**: Funções `addMiddleware`, `addMiddlewares` e `makeLoggingMiddleware` duplicadas

**Localização**:
- `internal/auth/presentation/router.go:29-43`
- `internal/notes/presentation/router.go:51-70`
- `internal/shared/presentation/router.go:44-58`

```go
// CÓDIGO DUPLICADO (3x no projeto!)
type middlewareFunc func(http.Handler) http.Handler

func addMiddleware(handler http.HandlerFunc, middleware middlewareFunc) http.Handler {
    return middleware(handler)
}

func makeLoggingMiddleware(logger *slog.Logger) middlewareFunc {
    return func(handler http.Handler) http.Handler {
        if logger != nil {
            return nethttp.LoggingMiddleware(logger)(handler)
        }
        return handler
    }
}
```

**Impactos**:
- ❌ **Manutenção cara**: Alterar middleware = mexer em 3 arquivos
- ❌ **Inconsistência**: `auth` tem `addMiddleware`, `notes` tem `addMiddlewares` (plural)
- ❌ **Viola DRY**: 45 linhas duplicadas
- ❌ **Testes triplicados**: Cada router precisa testar as mesmas funções

**Solução Proposta**:

```go
// pkg/nethttp/middleware.go (NOVO)
package nethttp

import (
    "log/slog"
    "net/http"
)

type MiddlewareFunc func(http.Handler) http.Handler

// Chain aplica múltiplos middlewares em ordem
func Chain(middlewares ...MiddlewareFunc) MiddlewareFunc {
    return func(next http.Handler) http.Handler {
        for i := len(middlewares) - 1; i >= 0; i-- {
            next = middlewares[i](next)
        }
        return next
    }
}

// MakeLoggingMiddleware cria middleware de logging
func MakeLoggingMiddleware(logger *slog.Logger) MiddlewareFunc {
    return func(next http.Handler) http.Handler {
        if logger != nil {
            return LoggingMiddleware(logger)(next)
        }
        return next
    }
}
```

```go
// internal/auth/presentation/router.go (REFATORADO)
package presentation

import (
    "fmt"
    "log/slog"
    "net/http"
    "pinnado/pkg/nethttp"
)

func SetupRouter(options SetupRouterOptions) {
    handler := NewAuthHandler(options.AuthService)
    
    // Usar funções de pkg/nethttp
    middleware := nethttp.MakeLoggingMiddleware(options.Logger)
    
    registerHandler := middleware(handler.Register)
    loginHandler := middleware(handler.Login)
    
    options.Mux.Handle(fmt.Sprintf("POST %s/auth/register", options.Prefix), registerHandler)
    options.Mux.Handle(fmt.Sprintf("POST %s/auth/login", options.Prefix), loginHandler)
}

// DELETAR: addMiddleware, makeLoggingMiddleware, middlewareFunc
```

**Estimativa**: 2 horas

**Benefícios**:
- ✅ Reduz de 45 para 0 linhas duplicadas
- ✅ Manutenção centralizada
- ✅ Consistência entre módulos
- ✅ Facilita adicionar novos middlewares

**Checklist**:
- [ ] Criar `pkg/nethttp/middleware.go`
- [ ] Refatorar `internal/auth/presentation/router.go`
- [ ] Refatorar `internal/notes/presentation/router.go`
- [ ] Refatorar `internal/shared/presentation/router.go`
- [ ] Deletar funções duplicadas
- [ ] Rodar `make test`

---

### #10: Inconsistência em Funções de Middleware

**Problema**: `auth` e `shared` usam `addMiddleware` (singular), `notes` usa `addMiddlewares` (plural)

**Localização**:
- `internal/auth/presentation/router.go:31` - `addMiddleware`
- `internal/notes/presentation/router.go:53` - `addMiddlewares`
- `internal/shared/presentation/router.go:46` - `addMiddleware`

```go
// auth/router.go (singular)
func addMiddleware(handler http.HandlerFunc, middleware middlewareFunc) http.Handler

// notes/router.go (plural - aceita varargs)
func addMiddlewares(handler http.HandlerFunc, middlewares ...middlewareFunc) http.Handler
```

**Impactos**:
- ❌ **Confusão**: Desenvolvedores não sabem qual usar
- ❌ **Inconsistência**: Mesmo padrão, nomes diferentes
- ❌ **Code review difícil**: Precisa verificar qual versão está sendo usada

**Solução Proposta**:

Será resolvido automaticamente com **#9** (extrair para `pkg/nethttp`).

**Estimativa**: Incluído no #9

**Benefícios**:
- ✅ Padronização automática
- ✅ Uma única API (plural com varargs)

**Checklist**:
- [ ] Resolver #9 primeiro
- [ ] Validar que todos os módulos usam `nethttp.Chain()`

---

### #11: Método `User.Update()` Não Utilizado

**Problema**: Método `Update()` definido mas nunca chamado

**Localização**: `internal/auth/domain/user.go` (aproximadamente linha 30-35)

```go
// CÓDIGO ATUAL (não utilizado)
func (u *User) Update(email Email, password Password) {
    u.Email = email
    u.Password = password
    u.UpdatedAt = time.Now()
}
```

**Impactos**:
- ❌ **Dead code**: Aumenta complexidade sem valor
- ❌ **Confusão**: Desenvolvedores podem pensar que existe feature de update
- ❌ **Cobertura de teste**: Método não testado ou testado sem uso real

**Solução Proposta**:

**Opção 1: Remover (Recomendado)**
```go
// Deletar método Update() completamente
```

**Opção 2: Implementar funcionalidade completa**
```go
// internal/auth/application/auth_service.go (NOVO)
func (s *authService) UpdateUser(ctx context.Context, userID string, input UpdateUserInput) error {
    // 1. Buscar usuário
    user, err := s.repo.FindByID(ctx, userID)
    if err != nil {
        return fmt.Errorf("s.repo.FindByID: %w", err)
    }
    
    // 2. Validar novos dados
    email, err := domain.NewEmail(input.Email)
    if err != nil {
        return err
    }
    
    password, err := domain.NewPassword(input.Password)
    if err != nil {
        return err
    }
    
    // 3. Atualizar
    user.Update(email, password)
    
    // 4. Persistir
    if err := s.repo.Update(ctx, user); err != nil {
        return fmt.Errorf("s.repo.Update: %w", err)
    }
    
    return nil
}
```

**Recomendação**: **Opção 1** (remover)

**Estimativa**: 15 minutos

**Benefícios**:
- ✅ Reduz dead code
- ✅ Clareza sobre features existentes
- ✅ Menos manutenção

**Checklist**:
- [ ] Verificar se há algum uso oculto do método
- [ ] Deletar `User.Update()` de `domain/user.go`
- [ ] Deletar testes relacionados (se houver)
- [ ] Rodar `make test`

---

### #12: Logging Inconsistente no `main.go`

**Problema**: Mix de `log.Println` (stdlib) e `slog` (structured logging)

**Localização**: `cmd/api/main.go`

```go
// CÓDIGO ATUAL (inconsistente)
func main() {
    slog.SetDefault(logger.NewLogger("info"))  // Linha 48
    
    log.Println("loading configuration...")    // Linha 50 (stdlib!)
    log.Println("connecting to MongoDB...")    // Linha 53
    // ...
    
    appLogger := logger.NewLogger("info")      // Linha 98 (duplicado!)
}
```

**Impactos**:
- ❌ **Logs não estruturados**: `log.Println` não tem níveis/campos
- ❌ **Dificulta parsing**: Logs em formato texto puro
- ❌ **Inconsistência**: Alguns lugares usam `slog`, outros `log`
- ❌ **Impossível filtrar**: Não há `log.level=info` no stdlib

**Solução Proposta**:

```go
// cmd/api/main.go (REFATORADO)
func main() {
    appLogger := logger.NewLogger("info")
    slog.SetDefault(appLogger)
    
    appLogger.Info("loading configuration")
    config := infrastructure.LoadConfig()
    
    appLogger.Info("connecting to MongoDB")
    ctx := context.Background()
    mongoClient, err := mongodb.NewMongoClient(ctx, /* ... */)
    if err != nil {
        appLogger.Error("failed to connect to MongoDB", "error", err)
        os.Exit(1)
    }
    defer func() {
        appLogger.Info("disconnecting from MongoDB")
        if err := mongoClient.Disconnect(context.Background()); err != nil {
            appLogger.Error("error disconnecting from MongoDB", "error", err)
        }
    }()
    
    appLogger.Info("creating MongoDB indexes")
    if err := authinfra.CreateIndexes(ctx, mongoClient.Database(config.Mongo.DBName)); err != nil {
        appLogger.Error("failed to create MongoDB indexes", "error", err)
        os.Exit(1)
    }
    
    appLogger.Info("server starting", "address", addr)
    // ...
}
```

**Estimativa**: 30 minutos

**Benefícios**:
- ✅ Logs estruturados (JSON)
- ✅ Facilita parsing por ferramentas (DataDog, ELK)
- ✅ Consistência em todo o projeto
- ✅ Níveis de log configuráveis

**Checklist**:
- [ ] Substituir todos `log.Println` por `appLogger.Info`
- [ ] Substituir todos `log.Fatalf` por `appLogger.Error + os.Exit(1)`
- [ ] Remover import `"log"` de `main.go`
- [ ] Testar startup da aplicação

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

### #14: Configuração Swagger Host Desnecessária

**Problema**: `docs.SwaggerInfo.Host` configurado dinamicamente no `main.go`

**Localização**: `cmd/api/main.go:95-96`

```go
// CÓDIGO ATUAL (desnecessário)
addr := fmt.Sprintf("%s:%s", config.Api.Host, config.Api.Port)
docs.SwaggerInfo.Host = addr
```

**Impactos**:
- ⚠️ **Swagger quebra em produção**: Se host for `0.0.0.0:8080`, Swagger UI tenta acessar `0.0.0.0`
- ⚠️ **Não funciona com proxy**: Swagger não sabe o domínio real (ex: `api.pinnado.com`)
- ⚠️ **Desnecessário**: Swagger UI consegue inferir host automaticamente

**Solução Proposta**:

```go
// cmd/api/main.go (REMOVER)
// DELETAR: docs.SwaggerInfo.Host = addr

// Swagger UI infere host automaticamente do browser
```

```go
// cmd/api/main.go (ALTERNATIVA - apenas se necessário)
// Configurar apenas em produção via env var
if os.Getenv("SWAGGER_HOST") != "" {
    docs.SwaggerInfo.Host = os.Getenv("SWAGGER_HOST")
}
```

**Estimativa**: 10 minutos

**Benefícios**:
- ✅ Swagger funciona em qualquer ambiente
- ✅ Menos configuração manual
- ✅ Funciona com proxies/load balancers

**Checklist**:
- [ ] Remover linha `docs.SwaggerInfo.Host = addr`
- [ ] Testar Swagger UI localmente
- [ ] Testar Swagger UI em produção (se houver)

---

### #15: Falta de Testes de Integração

**Problema**: Apenas testes unitários, sem testes de integração com MongoDB

**Localização**: Ausência de `*_integration_test.go` ou `tests/integration/`

**Impactos**:
- ❌ **Bugs em produção**: Índices, queries, aggregations podem falhar
- ❌ **Refactor arriscado**: Sem garantia de que repository funciona com MongoDB real
- ❌ **Sem validação E2E**: Fluxo completo não é testado
- ❌ **Confiança baixa**: Deploy sem garantia de funcionamento

**Solução Proposta**:

**Opção 1: Testcontainers (Recomendado)**
```go
// internal/auth/infrastructure/user_repository_integration_test.go (NOVO)
//go:build integration
// +build integration

package infrastructure_test

import (
    "context"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/mongodb"
    
    "pinnado/internal/auth/domain"
    "pinnado/internal/auth/infrastructure"
)

func setupMongoContainer(t *testing.T) (*mongo.Client, func()) {
    ctx := context.Background()
    
    mongoC, err := mongodb.RunContainer(ctx,
        testcontainers.WithImage("mongo:8"),
    )
    require.NoError(t, err)
    
    endpoint, err := mongoC.Endpoint(ctx, "")
    require.NoError(t, err)
    
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(endpoint))
    require.NoError(t, err)
    
    cleanup := func() {
        client.Disconnect(ctx)
        mongoC.Terminate(ctx)
    }
    
    return client, cleanup
}

func TestUserRepository_Create_Integration(t *testing.T) {
    client, cleanup := setupMongoContainer(t)
    defer cleanup()
    
    db := client.Database("test")
    collection := db.Collection(domain.UsersCollectionName)
    repo := infrastructure.NewUserRepository(collection)
    
    email, _ := domain.NewEmail("test@example.com")
    password, _ := domain.NewPassword("Pass123!")
    user := domain.NewUser(email, password)
    
    err := repo.Create(context.Background(), user)
    
    require.NoError(t, err)
    
    // Verifica no banco
    exists, err := repo.ExistsByEmail(context.Background(), email)
    require.NoError(t, err)
    assert.True(t, exists)
}
```

```makefile
# Makefile (ADICIONAR)
test-integration:
	go test -tags=integration ./... -v

test-all:
	make test && make test-integration
```

**Opção 2: MongoDB Memory Server (Alternativa)**
```bash
# Menos recomendado, mas mais rápido
go get github.com/tryvium-travels/memongo
```

**Recomendação**: **Testcontainers** (ambiente real)

**Estimativa**: 2 dias (setup + testes de repositories)

**Benefícios**:
- ✅ Valida queries MongoDB
- ✅ Testa índices únicos
- ✅ Detecta bugs antes de produção
- ✅ Aumenta confiança em deploys

**Checklist**:
- [ ] Adicionar `testcontainers-go` no `go.mod`
- [ ] Criar helper `setupMongoContainer()`
- [ ] Adicionar testes de integração para `UserRepository`
- [ ] Adicionar testes de integração para `NoteRepository`
- [ ] Adicionar `make test-integration`
- [ ] Configurar CI/CD para rodar testes de integração

---

### #16: Makefile Limitado

**Problema**: Makefile tem apenas comandos básicos

**Localização**: `Makefile` (16 linhas apenas)

```makefile
# CÓDIGO ATUAL (limitado)
all:
	go install github.com/vektra/mockery/v3@v3.6.3
	go install github.com/swaggo/swag/cmd/swag@latest
	go mod download

run:
	go run cmd/api/main.go

mock:
	mockery

test:
	go test ./... -cover

swag:
	swag init -g cmd/api/main.go -o docs
```

**Impactos**:
- ❌ **Produtividade baixa**: Desenvolvedores digitam comandos longos
- ❌ **Falta comandos úteis**: Coverage HTML, race detector, lint, docker
- ❌ **Sem padronização**: Cada dev roda testes de forma diferente

**Solução Proposta**:

```makefile
# Makefile (EXPANDIDO)
.PHONY: all install run test test-race test-coverage test-integration mock swag lint format clean docker-up docker-down docker-logs help

# Instalação de dependências
all: install

install:
	@echo "Installing dependencies..."
	go install github.com/vektra/mockery/v3@v3.6.3
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go mod download

# Desenvolvimento
run:
	go run cmd/api/main.go

# Testes
test:
	go test ./... -cover -v

test-race:
	go test ./... -race

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

test-integration:
	go test -tags=integration ./... -v

test-all: test test-integration

# Geração de código
mock:
	mockery

swag:
	swag init -g cmd/api/main.go -o docs

# Qualidade de código
lint:
	golangci-lint run

format:
	go fmt ./...
	goimports -w .

# Docker
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Limpeza
clean:
	rm -f coverage.out coverage.html
	go clean -testcache

# Ajuda
help:
	@echo "Available commands:"
	@echo "  make install          - Install dependencies"
	@echo "  make run              - Run application"
	@echo "  make test             - Run unit tests"
	@echo "  make test-race        - Run tests with race detector"
	@echo "  make test-coverage    - Generate coverage report"
	@echo "  make test-integration - Run integration tests"
	@echo "  make test-all         - Run all tests"
	@echo "  make mock             - Generate mocks"
	@echo "  make swag             - Generate Swagger docs"
	@echo "  make lint             - Run linter"
	@echo "  make format           - Format code"
	@echo "  make docker-up        - Start Docker services"
	@echo "  make docker-down      - Stop Docker services"
	@echo "  make clean            - Clean build artifacts"
	@echo "  make help             - Show this help"
```

**Estimativa**: 1 hora

**Benefícios**:
- ✅ Comandos padronizados
- ✅ Produtividade aumentada
- ✅ Fácil onboarding de novos devs
- ✅ Integração fácil com CI/CD

**Checklist**:
- [ ] Expandir Makefile com novos comandos
- [ ] Adicionar `make help`
- [ ] Documentar comandos no README
- [ ] Configurar CI/CD para usar `make test-all`

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

### #18: Estrutura de Erros Simples

**Problema**: Erros são apenas `errors.New()` sem contexto adicional

**Localização**: `internal/auth/domain/errors.go`, etc

```go
// CÓDIGO ATUAL (simples)
var (
    ErrUserNotFound      = errors.New("user not found")
    ErrUserAlreadyExists = errors.New("user already exists")
    ErrInvalidEmail      = errors.New("invalid email format")
)
```

**Impactos**:
- ⚠️ **Sem contexto**: Não sabe qual email/user causou erro
- ⚠️ **Difícil debug**: Erro genérico sem detalhes
- ⚠️ **Logging pobre**: Logs não têm informações úteis
- ⚠️ **API responses genéricos**: Cliente não recebe detalhes

**Solução Proposta**:

```go
// internal/auth/domain/errors.go (ENRIQUECIDO)
package domain

import "fmt"

// DomainError representa um erro de domínio com contexto
type DomainError struct {
    Code    string         `json:"code"`
    Message string         `json:"message"`
    Details map[string]any `json:"details,omitempty"`
}

func (e *DomainError) Error() string {
    if len(e.Details) == 0 {
        return fmt.Sprintf("[%s] %s", e.Code, e.Message)
    }
    return fmt.Sprintf("[%s] %s (details: %v)", e.Code, e.Message, e.Details)
}

// Construtores de erros
func ErrUserNotFound(userID string) *DomainError {
    return &DomainError{
        Code:    "USER_NOT_FOUND",
        Message: "user not found",
        Details: map[string]any{"user_id": userID},
    }
}

func ErrUserAlreadyExists(email Email) *DomainError {
    return &DomainError{
        Code:    "USER_ALREADY_EXISTS",
        Message: "user with this email already exists",
        Details: map[string]any{"email": string(email)},
    }
}

func ErrInvalidEmail(value string) *DomainError {
    return &DomainError{
        Code:    "INVALID_EMAIL",
        Message: "invalid email format",
        Details: map[string]any{"email": value},
    }
}
```

```go
// internal/auth/presentation/handler.go (USAR)
func MapErrorToHTTPStatus(err error) (int, ErrorResponse) {
    var domainErr *domain.DomainError
    if errors.As(err, &domainErr) {
        switch domainErr.Code {
        case "USER_NOT_FOUND":
            return http.StatusNotFound, ErrorResponse{
                Message: domainErr.Message,
                Code:    domainErr.Code,
                Details: domainErr.Details,
            }
        case "USER_ALREADY_EXISTS":
            return http.StatusConflict, ErrorResponse{
                Message: domainErr.Message,
                Code:    domainErr.Code,
            }
        default:
            return http.StatusBadRequest, ErrorResponse{
                Message: domainErr.Message,
                Code:    domainErr.Code,
            }
        }
    }
    
    return http.StatusInternalServerError, ErrorResponse{
        Message: "internal server error",
    }
}
```

**Estimativa**: 1 dia (refactor de todos os erros)

**Benefícios**:
- ✅ Erros com contexto rico
- ✅ Facilita debugging
- ✅ API responses mais úteis
- ✅ Logs estruturados com detalhes

**Nota**: Implementação opcional. Pode adicionar complexidade desnecessária se erros simples são suficientes.

**Checklist**:
- [ ] Criar `DomainError` struct
- [ ] Refatorar erros de `auth/domain`
- [ ] Refatorar erros de `notes/domain`
- [ ] Atualizar `MapErrorToHTTPStatus`
- [ ] Atualizar testes

---

### #19: Falta Context Timeout em Handlers

**Problema**: Handlers não definem timeout para operações longas

**Localização**: Todos os handlers (ex: `internal/auth/presentation/handler.go`)

```go
// CÓDIGO ATUAL (sem timeout)
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    // r.Context() não tem timeout
    output, err := h.authService.Register(r.Context(), input)
    // se MongoDB travar, request fica pendurada eternamente
}
```

**Impactos**:
- ⚠️ **Requests penduradas**: Se MongoDB travar, requests não retornam
- ⚠️ **Sem controle**: Cliente não sabe quanto tempo esperar
- ⚠️ **Recursos vazando**: Goroutines esperando eternamente
- ⚠️ **Degradação**: Sistema fica lento sem timeout

**Solução Proposta**:

```go
// internal/auth/presentation/handler.go (COM TIMEOUT)
package presentation

import (
    "context"
    "net/http"
    "time"
)

const (
    requestTimeout = 30 * time.Second
)

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    // Criar context com timeout
    ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
    defer cancel()
    
    // Parse request
    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        nethttp.JSON(w, http.StatusBadRequest, ErrorResponse{Message: "invalid request"})
        return
    }
    
    // Chamar service com context com timeout
    output, err := h.authService.Register(ctx, application.UserInput{
        Email:    strings.TrimSpace(req.Email),
        Password: strings.TrimSpace(req.Password),
    })
    
    // Verificar se timeout ocorreu
    if ctx.Err() == context.DeadlineExceeded {
        nethttp.JSON(w, http.StatusGatewayTimeout, ErrorResponse{
            Message: "request timeout",
        })
        return
    }
    
    // Continuar normalmente
    // ...
}
```

**Alternativa: Middleware Global**
```go
// pkg/nethttp/timeout_middleware.go (NOVO)
package nethttp

import (
    "context"
    "net/http"
    "time"
)

func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx, cancel := context.WithTimeout(r.Context(), timeout)
            defer cancel()
            
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

```go
// internal/auth/presentation/router.go (USAR)
middleware := nethttp.Chain(
    nethttp.TimeoutMiddleware(30 * time.Second),
    nethttp.MakeLoggingMiddleware(options.Logger),
)
```

**Recomendação**: **Middleware Global** (mais simples)

**Estimativa**: 2 horas

**Benefícios**:
- ✅ Proteção contra requests penduradas
- ✅ Controle de tempo de resposta
- ✅ Melhor UX (timeout rápido vs espera infinita)
- ✅ Previne vazamento de recursos

**Checklist**:
- [ ] Criar `TimeoutMiddleware` em `pkg/nethttp`
- [ ] Adicionar em chain de middlewares de todos os routers
- [ ] Testar com MongoDB lento (simular timeout)
- [ ] Documentar timeout no Swagger (@Param timeout)

---

### #20: Graceful Shutdown Sem Logging Detalhado

**Problema**: Graceful shutdown não loga quantos requests foram finalizados/cancelados

**Localização**: `cmd/api/main.go:135-146`

```go
// CÓDIGO ATUAL (sem detalhes)
log.Println("shutting down server...")

shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
defer cancel()

if err := server.Shutdown(shutdownCtx); err != nil {
    log.Fatalf("server forced to shutdown: %v", err)
}

log.Println("server exited gracefully")
```

**Impactos**:
- ⚠️ **Sem visibilidade**: Não sabe se requests foram finalizados ou cancelados
- ⚠️ **Debug difícil**: Em caso de problema, não sabe o que aconteceu
- ⚠️ **Sem métricas**: Impossível medir tempo de shutdown

**Solução Proposta**:

```go
// cmd/api/main.go (COM LOGGING DETALHADO)
func main() {
    // ...
    
    appLogger.Info("server starting", "address", addr)
    
    go func() {
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            appLogger.Error("server failed to start", "error", err)
            os.Exit(1)
        }
    }()
    
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    sig := <-quit
    
    appLogger.Info("shutdown signal received", "signal", sig)
    
    shutdownStart := time.Now()
    shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
    defer cancel()
    
    appLogger.Info("shutting down server gracefully", 
        "timeout", shutdownTimeout.String())
    
    if err := server.Shutdown(shutdownCtx); err != nil {
        appLogger.Error("server forced to shutdown", 
            "error", err,
            "elapsed", time.Since(shutdownStart))
        os.Exit(1)
    }
    
    appLogger.Info("server exited gracefully", 
        "elapsed", time.Since(shutdownStart))
}
```

**Estimativa**: 30 minutos

**Benefícios**:
- ✅ Visibilidade de shutdown
- ✅ Métricas de tempo de shutdown
- ✅ Facilita debug de problemas
- ✅ Logs estruturados

**Checklist**:
- [ ] Adicionar logging detalhado no shutdown
- [ ] Logar sinal recebido (SIGINT, SIGTERM)
- [ ] Logar tempo de shutdown
- [ ] Testar com `kill -SIGTERM <pid>`

---

## 📈 Roadmap de Implementação

### **Fase 1: Correções Críticas** (Semana 1)
- [x] Análise de débitos técnicos
- [ ] #8: Resolver violação `pkg/` → `internal/` (4h)
- [ ] #9: Extrair middleware duplicado (2h)
- [ ] #10: Padronizar middleware (incluído no #9)
- [ ] #11: Remover `User.Update()` (15min)
- [ ] #5: Unificar logger (2h)

**Resultado**: Arquitetura consistente + código limpo

---

### **Fase 2: Qualidade e Padronização** (Semana 2)
- [ ] #12: Padronizar logging no `main.go` (30min)
- [ ] #13: Validação de sort fields (1h)
- [ ] #14: Remover Swagger Host desnecessário (10min)
- [ ] #2: Refatorar `shared/config` (1 dia)
- [ ] #16: Expandir Makefile (1h)
- [ ] #17: Validar env vars obrigatórias (2h)

**Resultado**: Código padronizado + melhor DX

---

### **Fase 3: Observabilidade** (Semana 3)
- [ ] #4: Adicionar Request ID (1 dia)
- [ ] #19: Context timeout em handlers (2h)
- [ ] #20: Logging detalhado no shutdown (30min)
- [ ] #7: Health check robusto (1 dia)

**Resultado**: Logs rastreáveis + observabilidade

---

### **Fase 4: Infraestrutura** (Semana 4)
- [ ] #3: Migração de índices (1 dia)
- [ ] #6: Middleware chain (1 dia)
- [ ] #15: Testes de integração (2 dias)

**Resultado**: Deploy mais seguro + cobertura de testes

---

### **Fase 5: Escalabilidade** (Semana 5-6)
- [ ] #1: Dependency Injection (Wire) (3 dias)
- [ ] #18: Estrutura de erros rica (1 dia) - Opcional
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
- ✅ Testes de integração (repositories)
- ✅ Testes E2E básicos

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
- **Última atualização**: 2026-02-15
- **Responsável**: Equipe de Arquitetura
- **Revisão**: Mensal (ou quando adicionar novo módulo)
- **Arquivo vivo**: Sempre atualizar quando resolver débito

---

## 📋 Changelog

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
