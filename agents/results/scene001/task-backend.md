# Tarefa: Arquitetura Completa do Módulo de Autenticação (Registro e Login Automático)

## 1. Domain Layer

Tarefa: Value Object Email
Definição técnica: Tipo `Email` com construtor para validação RFC 5322 e retorno de valor primitivo.
Arquivo: internal/auth/domain/email.go

Tarefa: Value Object Password
Definição técnica: Tipo `Password` com construtor para validação de complexidade, métodos para geração de hash e comparação segura.
Arquivo: internal/auth/domain/password.go

Tarefa: Entidade User
Definição técnica: Entidade `User` composta pelos tipos `Email` e `Password`, garantindo integridade de campos obrigatórios e timestamps.
Arquivo: internal/auth/domain/user.go

## 2. Infrastructure Layer (Persistence & Adapters)

Tarefa: MongoDB Indexes
Definição técnica: Função para garantia de índice único no campo de e-mail na collection de usuários. Também adicionar func na cmd/api/main.go para rodar as indexes
Arquivo: internal/auth/infrastructure/mongo_indexes.go

Tarefa: User Repository (MongoDB)
Definição técnica: Implementação dos métodos de persistência e busca de usuários.
Arquivo: internal/auth/infrastructure/user_repository.go

## 3. Application Layer (Usecases & Ports)

Tarefa: DTOs de Fluxo
Definição técnica: Estruturas para transporte de dados de registro (entrada) e perfil de usuário.
Arquivo: internal/auth/application/dto.go

Tarefa: Interfaces (Ports)
Definição técnica: Definição dos contratos para `UserRepository` e `AuthService`.
Arquivo: internal/auth/application/port.go

Tarefa: Auth Service
Definição técnica: Implementação do caso de uso criar usuário. Deve orquestrar validação de domínio, unicidade, persistência e geração de credenciais.
Arquivo: internal/auth/application/user_service.go

## 4. Presentation Layer (API)

Tarefa: Contratos HTTP
Definição técnica: Estruturas de `Request` com tags de validação de campos e `Response` formatada para o frontend.
Arquivo: internal/auth/presentation/request_response.go

Tarefa: Auth Handler
Definição técnica: Handler para processar submissões de formulário, sanitização de entrada e mapeamento de erros de negócio para status codes adequados.
Arquivo: internal/auth/presentation/handler.go

Tarefa: Router Setup
Definição técnica: Configuração do endpoint de registro no roteador da aplicação.
Arquivo: internal/auth/presentation/router.go

## 5. Orquestração (Bootstrap)

Tarefa: Dependency Injection
Definição técnica: Inicialização e conexão de todas as camadas no ponto de entrada da aplicação, incluindo provedores de token e banco de dados.
Arquivo: cmd/api/main.go
