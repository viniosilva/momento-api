# Pinnado: Notas Colaborativas

Um aplicativo de notas para escrever e compartilhar conteúdos de forma colaborativa. Organizando projetos e estudos com eficiência e segurança.


## Tecnologias Utilizadas

- **Go 1.25.0**: Linguagem de programação
- **MongoDB 8**: Banco de dados NoSQL
- **Swagger/OpenAPI**: Documentação da API
- **Testify**: Framework de testes
- **Mockery**: Geração de mocks para testes
- **Make**


## Estrutura do Projeto

O projeto está organizado em módulos independentes dentro de `/internal`:
- **auth**: Módulo de autenticação e autorização
- **shared**: Código compartilhado entre módulos (health checks, configuração)

```
/cmd                    	  // Diretório raiz para pontos de entrada da aplicação
  /api                  	  // Módulo específico para API HTTP (pontos de entrada da aplicação)
    main.go             	  // Arquivo principal que inicializa e executa o servidor
/internal               	  // Código interno da aplicação (não exportado para outros projetos)
  /{module_name}        	  // Placeholder para nome do módulo (ex: user, product, note)
    /application        	  // Camada de aplicação: Serviços e Casos de Uso (Orquestração)
	  dto.go            	  // Data Transfer Objects para comunicação entre camadas
	  {name}_service_test.go  // Testes unitários do serviço de aplicação
	  {name}_service.go       // Implementação do serviço que orquestra casos de uso
	  port.go           	  // Interfaces (portas) que definem contratos de serviços
    /domain             	  // Camada de domínio: Entidades, Value Objects e Interfaces do Domínio
	  {name}_test.go    	  // Testes unitários das entidades e value objects
	  {name}.go         	  // Implementação de entidades e value objects do domínio
    /infrastructure     	  // Camada de infraestrutura: Implementações de Repositórios, DB e Clientes API
    /presentation       	  // Camada de apresentação: Handlers HTTP/gRPC, DTOs e Definição de Rotas
	  handler_test.go    	  // Testes unitários dos handlers HTTP
	  handler.go        	  // Implementação dos handlers que processam requisições HTTP
	  request_response.go  	  // DTOs específicos para requisições e respostas HTTP
	  router_test.go   		  // Testes unitários das definições de rotas
	  router.go        		  // Definição e configuração das rotas da API
/pkg                   		  // Código compartilhável entre projetos (bibliotecas reutilizáveis)
  /{name}              		  // Placeholder para nome do pacote compartilhado
	{name}_test.go     		  // Testes unitários do pacote compartilhado
	{name}.go          		  // Implementação do pacote compartilhado
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
