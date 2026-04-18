# Momento: Notas Colaborativas

Um aplicativo de notas para escrever e compartilhar conteúdos de forma colaborativa. Organizando projetos e estudos com eficiência e segurança.


## Tecnologias Utilizadas

- **Go 1.26**: Linguagem de programação
- **MongoDB 8**: Banco de dados NoSQL
- **Swagger/OpenAPI**: Documentação da API
- **Testify**: Framework de testes
- **Mockery**: Geração de mocks para testes
- **Make**


## Estrutura do Projeto

O projeto está organizado em módulos independentes dentro de `/internal`, seguindo Arquitetura Hexagonal (Ports & Adapters):
- **auth**: Autenticação e autorização (registro, login, JWT)
- **notes**: Gerenciamento de notas (criação, listagem, arquivo/restauração)
- **shared**: Código compartilhado entre módulos (health check)

```
/cmd                          // Pontos de entrada da aplicação
  /api                        // Ponto de entrada do servidor HTTP
    main.go                   // Composição da aplicação: injeção de dependências e inicialização
  /migrate                    // Ponto de entrada das migrações
    main.go                   // Criação de índices no MongoDB

/internal
  /{module}                   // Módulo independente (auth, notes, shared)
    /domain                   // Entidades, Value Objects e erros de domínio (sem dependências externas)
      {name}.go               // Entidade ou Value Object com validação
      {name}_test.go          // Testes unitários de domínio
    /app                      // Serviços de aplicação e interfaces (portas de entrada/saída)
      port.go                 // Interfaces que o app expõe e consome (UserRepository, JWTService, etc.)
      dto.go                  // Data Transfer Objects entre camadas (Input/Output)
      {name}_service.go       // Orquestração de casos de uso
      {name}_service_test.go  // Testes unitários do serviço com mocks
    /adapters                 // Implementações de infraestrutura (MongoDB, JWT, etc.)
      {name}_repository.go    // Implementação do repositório (satisfaz interface de app/port.go)
      {name}_model.go         // Documento de persistência com bson tags (desacoplado do domínio)
      mongo_indexes.go        // Orquestrador de criação de índices
      /indexes
        001_{name}.go         // Definição individual de índice MongoDB
    /mocks                    // Mocks gerados pelo Mockery (via make mock)
      mock_{name}.go          // Mock gerado a partir das interfaces de app/port.go
    /ports                    // Handlers HTTP, rotas e DTOs de request/response
      port.go                 // Interfaces do serviço consumidas pelo handler
      handler.go              // Handlers HTTP que delegam para o serviço de app
      handler_test.go         // Testes unitários dos handlers com mocks
      request_response.go     // Structs de request e response HTTP
      router.go               // Registro de rotas e middlewares
      router_test.go          // Testes das rotas via httptest
```


## Como Executar

### 1. Configuração do Ambiente

Crie um arquivo `.env` na raiz do projeto conforme o `.env.example`

### 2. Instalação de Dependências

```bash
make
```

### 3. Executar a Aplicação

```bash
make run
```

A API estará disponível em `http://localhost:8080/api`

### 4. Documentação da API (Swagger)

Após iniciar a aplicação, acesse a documentação Swagger em:

```
http://localhost:8080/docs/swagger/index.html
```


## Testes

### Executar Testes

```bash
make test
```


### Gerar Mocks

Gere os mocks com:

```bash
make mock
```


## Contribuindo

1. Siga os princípios de design descritos neste README
2. Mantenha a cobertura de testes acima de 80%
3. Use commits semânticos
4. Documente mudanças significativas


## Licença

Este projeto está sob a licença Apache 2.0.
