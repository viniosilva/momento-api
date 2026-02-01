# Arquitetura de Tarefa: Implementação de Healthcheck

1. Domain Layer
Tarefa: Value Object Health
Definição técnica: Struct `HealthStatus` com campo `Status` (enum).
Arquivo: internal/shared/domain/health.go

2. Application Layer
Tarefa: Interfaces e DTOs de Health
Definição técnica: 
- Struct `HealthOutput` (Status string).
- Interface `HealthService` com assinatura `HealthCheck(ctx context.Context) HealthOutput`.
Arquivos: internal/shared/application/port.go
          internal/shared/application/dto.go


Tarefa: Application Service
Definição técnica: Struct `healthService` que implementa `HealthService`. O método `HealthCheck` deve instanciar o domínio com status "ok" e retornar o DTO de saída.
Arquivo: internal/shared/application/health_service.go

3. Presentation Layer
Tarefa: Contratos de API
Definição técnica: Struct `HealthResponse` com tags JSON para `status`.
Arquivo: internal/shared/presentation/request_response.go

Tarefa: Health Handler
Definição técnica: Método `HealthCheck` que invoca `HealthService.HealthCheck` e renderiza o `HealthResponse` com código HTTP 200.
Arquivo: internal/shared/presentation/handler.go

Tarefa: Router
Definição técnica: Função `SetupHealthRouter` para registrar a rota `GET /health` apontando para o handler correspondente.
Arquivo: internal/shared/presentation/router.go

4. Infrastructure Layer
Tarefa: Entrypoint / Server Setup
Definição técnica: Inicialização do router e injeção de dependência do `HealthService` no `Handler`.
Arquivo: cmd/api/main.go