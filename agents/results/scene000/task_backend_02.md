# Tarefa: Infraestrutura Base, Persistência e Resiliência (Padrão /pkg)

1. Utilitários Técnicos (Agnósticos)
Tarefa: MongoDB Driver Setup
Definição técnica: Função NewMongoClient(host, port, dbName, user, pass string) que retorna (*mongo.Client, error). Lógica de Retry inclusa.
Arquivo: pkg/mongodb/client.go

2. Configuração de Ambiente (Orquestração)
Tarefa: Config Loader
Definição técnica: Struct Config mapeando PORT, MONGO_URI, MONGO_DB, APP_ENV. Uso de godotenv.
Arquivo: internal/shared/infrastructure/config.go

3. Aplicação e Monitoramento
Tarefa: HealthService com Integração de DB
Definição técnica: Interface HealthService com método HealthCheck(ctx). Implementação recebe *mongo.Client.
Arquivo: internal/shared/application/health_service.go

4. Camada de Apresentação e Documentação
Tarefa: Swagger e Rota de Healthcheck
Definição técnica: Handler de Healthcheck e registro de rota /docs/*any via http-swagger.
Arquivo: internal/shared/presentation/router.go

5. Ciclo de Vida e Bootstrap (Orquestrador)
Tarefa: Main Orchestrator
Definição técnica: Orquestração de startup (Config -> DB -> Health -> Server) e Graceful Shutdown.
Arquivo: cmd/api/main.go

6. Orquestração de Infra e Testes
Tarefa: Docker Compose e Makefile
Definição técnica:
- docker-compose.yaml: MongoDB 8.2
- Makefile: Adição de target make mocks para execução do mockery.
Arquivo: / (Root)

Tarefa: Configuração do Mockery
Definição técnica: Arquivo .mockery.yaml configurado para gerar mocks recursivamente a partir de internal/, enviando-os para internal/shared/mocks ou pastas mocks adjacentes, com a flag with-expecter: true.
Arquivo: .mockery.yaml