# Tarefa: Telemetria e Logs Estruturados (Padrão /pkg)

1. Utilitários Técnicos (Agnósticos no /pkg)
Tarefa: Logger Provider
Definição técnica: 
- Função `NewLogger(level string)` que retorna uma instância de `*slog.Logger`.
- Implementar lógica para extrair `trace_id` e `request_id` de um `context.Context` e anexá-los como atributos dinâmicos.
- **Regra:** Não lê variáveis de ambiente. Recebe `level` e formato via parâmetro.
Arquivo: pkg/logger/logger.go

Tarefa: Telemetry SDK Setup
Definição técnica:
- Função `InitTracer(serviceName, exporterEndpoint string)` que configura o TracerProvider e o OTLP Exporter.
- Retorno de uma função `func(context.Context) error` para o Shutdown limpo.
- **Regra:** Recebe o endpoint e o nome do serviço via parâmetro.
Arquivo: pkg/telemetry/otel.go

2. Camada de Apresentação (Transporte)
Tarefa: Middlewares de Rastreabilidade e Log
Definição técnica:
- `RequestIDMiddleware`: Middleware que garante a presença de `X-Request-ID` no contexto e no header de resposta.
- `TelemetryMiddleware`: Inicia spans do OpenTelemetry correlacionados ao contexto da requisição.
- `LoggingMiddleware`: Utiliza o logger do `/pkg` para registrar o ciclo de vida da requisição (Início/Fim, Latência, Status Code).
Arquivo: internal/shared/presentation/middleware.go

3. Orquestração e Ciclo de Vida (Main)
Tarefa: Main Bootstrap
Definição técnica:
- No `main.go`:
  1. Chamar `telemetry.InitTracer(cfg.ServiceName, cfg.OtelEndpoint)`.
  2. Inicializar logger global via `logger.NewLogger(cfg.LogLevel, cfg.LogJSON)`.
  3. Registrar middlewares na ordem: RequestID -> Telemetry -> Logging.
  4. Adicionar o `Shutdown` da telemetria no fluxo de Graceful Shutdown.
Arquivo: cmd/api/main.go

Regras de Design
1. **Injeção de Parâmetros**: Todo componente no `/pkg` deve ser testável isoladamente, recebendo suas dependências no construtor.
2. **Context Propagation**: O `context.Context` deve ser passado de ponta a ponta (Handler -> Service -> Repo/Lib) para que o Logger e o Tracer funcionem corretamente.
3. **Isolamento**: Pacotes em `/pkg` não importam nada de `internal/`.
