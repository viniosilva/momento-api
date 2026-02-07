# Boas Práticas Backend Go

## Princípios de Design


### Comentários

Minimize o máximo possível. Comentários são exceções como explicar uma regex se o nome da variável não for claro.

**Evite comentários que:**
- Explicam o que o código faz (o código deve ser autoexplicativo)
- Estão desatualizados ou não refletem o código atual
- Duplicam informações já presentes no código

**Use comentários para:**
- Explicar o "porquê" de uma decisão técnica complexa
- Documentar regex complexas ou algoritmos não triviais
- Avisar sobre limitações ou comportamentos inesperados

#### Exemplo

```go
// ❌ Ruim - comentário desnecessário
// Cria um novo usuário
func CreateUser(user User) error {
    // ...
}

// ✅ Bom - código autoexplicativo
func CreateUser(user User) error {
    // ...
}

// ✅ Bom - explica o porquê de uma decisão técnica
// Usa regex complexa para validar formato E.164 internacional
// Referência: https://en.wikipedia.org/wiki/E.164
phoneRegex := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
```

### Domain-Driven Design (DDD)

- Entidades: Crie modelos que representam objetos puros do domínio (ex: `User`, `Product`).
- Agregados: Use agregados para garantir a consistência dos dados e encapsular regras de negócio.
- Repositórios: Utilize interfaces de repositório para isolar a lógica de acesso a dados da camada de domínio.


#### Exemplo

```go
type User struct {
    ID    uuid.UUID
    Name  string
    Email string
}

type UserRepository interface {
    Save(ctx context.Context, user User) error
    FindByID(ctx context.Context, id uuid.UUID) (User, error)
}
```


### Princípios SOLID

- S - Single Responsibility Principle: Cada struct ou arquivo deve ter uma única responsabilidade. 
- O - Open/Closed Principle: Entidades abertas para extensão, mas fechadas para modificação.
- L - Liskov Substitution Principle: Subtipos devem ser substituíveis por seu tipo base.
- I - Interface Segregation Principle: Quanto maior a interface, mais fraca a abstração. Prefira interfaces pequenas e granulares.
- D - Dependency Inversion Principle: Dependa de abstrações (interfaces), não de implementações concretas. Injete dependências via construtores.

#### Exemplo

```go
type UserService interface {
    Register(user User) error
}

type userService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
    return &userService{repo: repo}
}
```


### Object Calisthenics
- Uma struct por responsabilidade: Agrupe comportamentos relacionados de forma coesa.
- Caminho Feliz à Esquerda: Proibido o uso de else. Use retornos antecipados (Guard Clauses) para tratar erros.
- Métodos Pequenos: Mantenha as funções curtas, focadas e com nomes descritivos.
- Baixa Indentação: Minimize níveis de indentação. Se o código estiver muito aninhado, extraia para funções privadas.
- Contexto: context.Context deve ser sempre o primeiro parâmetro em funções de I/O, DB e chamadas externas.
- Rule: "Tell, don’t ask": Delegue a lógica para o objeto que possui os dados.

#### Exemplo

```go
var ErrUserAlreadyExists = errors.New("user already exists")

func (s *userService) CreateUser(ctx context.Context, input CreateUserInput) (UserOutput, error) {
	phone, err := domain.NewPhone(input.Phone)
	if err != nil {
		return UserOutput{}, fmt.Errorf("domain.NewPhone: %w", err)
	}

	password, err := domain.NewPassword(input.Password)
	if err != nil {
		return UserOutput{}, fmt.Errorf("domain.NewPassword: %w", err)
	}

	existing, err := s.userRepository.HasByPhone(ctx, phone)
    if err != nil {
		return UserOutput{}, err
	}
    if existing {
		return UserOutput{}, ErrUserAlreadyExists
	}

	user := domain.NewUser(phone, password)
	if err := s.userRepository.CreateUser(ctx, user); err != nil {
		return UserOutput{}, fmt.Errorf("s.userRepository.CreateUser: %w", err)
	}

	return UserOutput{
		ID:    user.ID.Hex(),
		Phone: user.Phone,
	}, nil
}
```


## Testes Unitários e Mocking com mockery

- Ferramentas: Use testing, testify/assert e testify/require.
- Mockery: Gere mocks via make mocks. Use a interface .EXPECT() para segurança de tipos.
- Subtestes: Utilize t.Run para organizar casos de sucesso e falha.
- Helpers: Use funções auxiliares com t.Helper() para reduzir boilerplate nos testes.

#### Exemplo

**Nota:** Este exemplo mostra a estrutura completa de testes. Na prática, seguindo TDD, você testaria apenas o caminho feliz e os cenários de erro por vez.

```go
package application_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	
	"pinnado/internal/auth/application"
	"pinnado/internal/auth/domain"
	"pinnado/internal/auth/application/mocks"
)

func TestUserService_CreateUser(t *testing.T) {
	userRepoMock, userMockDefault := helperUserRepoMock(t)
	userSvc := application.NewUserService(userRepoMock)

	userInputDefault := application.CreateUserInput{
		Phone:    "+55 11 91234-5678",
		Password: "Strong#123",
	}

	// Caminho feliz - sempre testar
	t.Run("should create user when phone is not registered", func(t *testing.T) {
		userRepoMock.EXPECT().HasByPhone(mock.Anything, userMockDefault.Phone).
			Return(false, nil).Once()

		userRepoMock.EXPECT().CreateUser(mock.Anything, mock.Anything).
			Return(nil).Once()

		got, err := userSvc.CreateUser(context.Background(), userInputDefault)
		require.NoError(t, err)

		assert.NotEmpty(t, got.ID)
	})

	// Apenas um cenário de erro representativo (seguindo TDD)
	t.Run("should return error when phone already exists", func(t *testing.T) {
		userRepoMock.EXPECT().HasByPhone(mock.Anything, userMockDefault.Phone).
			Return(true, nil).Once()

		_, err := userSvc.CreateUser(context.Background(), userInputDefault)

		assert.ErrorIs(t, err, application.ErrUserAlreadyExists)
	})
}

func helperUserRepoMock(t *testing.T) (*mocks.MockUserRepository, *domain.User) {
	t.Helper()

	userRepoMock := mocks.NewMockUserRepository(t)

	userPhone, err := domain.NewPhone("+55 11 91234-5678")
	require.NoError(t, err)

	userPass, err := domain.NewPassword("Strong#123")
	require.NoError(t, err)

	userMock := domain.NewUser(userPhone, userPass)

	return userRepoMock, userMock
}
```

### Test-Driven Development (TDD)

TDD (Test-Driven Development) é uma prática de desenvolvimento onde você escreve os testes **antes** de implementar a funcionalidade. O ciclo TDD segue três etapas: **Red → Green → Refactor**.

**Domain Layer (Value Objects e Entidades):**
- Teste primeiro as validações e regras de negócio a fim que o resultado seja de sucesso
- Garanta que construtores falhem rápido (Fail-Fast) com dados inválidos
- Teste comportamentos, não implementação
- Teste com 100% de cobertura se possível

**Application Layer (Services):**
- Teste casos de uso completos que o resultado seja de sucesso
- Use mocks para dependências externas com banco de dados e integrações
- Teste tanto o caminho feliz quanto os casos de erro
- Para os casos de erro, apenas um cenário de erro que possibilite testar é só suficiente.
  Ex: um teste que o usuário falhe ao tentar criar, não precisa testar todos os casos de falha, um somente.

**Infrastructure Layer (Repositories):**
- Não testar, iremos testar com testes de integração

**Presentation Layer (Handlers):**
- Teste parsing de requests e validações
- Teste formatação de responses
- Use httptest para testar handlers HTTP
- Para os casos de erro, apenas um cenário de erro que possibilite testar é só suficiente.
  Ex: um teste que o usuário falhe ao tentar criar, não precisa testar todos os casos de falha, um somente.
- Foco em um teste para cada status code que é possível retornar


#### Boas Práticas TDD em Go

1. **Um teste, uma responsabilidade**: Cada teste deve verificar um comportamento específico
2. **Nomes descritivos**: Use `should [comportamento] when [condição]` nos nomes dos testes
3. **Arrange-Act-Assert**: Estruture testes em três fases claras
4. **Testes independentes**: Cada teste deve poder rodar isoladamente
5. **Mocks apenas quando necessário**: Prefira dependências reais quando possível
6. **Teste comportamentos, não implementação**: Foque no "o quê", não no "como"


#### Ferramentas para TDD em Go

- **testing**: Pacote padrão do Go
- **testify/assert**: Asserções mais legíveis
- **testify/require**: Asserções que param a execução em caso de falha
- **testify/mock**: Criação de mocks
- **mockery**: Geração automática de mocks a partir de interfaces
- **httptest**: Testes para handlers HTTP
- **go test -race**: Detecção de race conditions

#### Comandos Úteis

```bash
# Rodar testes em modo watch (requer ferramentas externas)
go test ./... -watch

# Rodar testes com cobertura
go test ./... -cover

# Rodar testes com cobertura detalhada
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Rodar testes em paralelo
go test ./... -parallel 4

# Rodar testes específicos
go test ./internal/auth/domain -run TestNewPhone
```

#### Checklist TDD

Antes de começar a implementar:
- [ ] Entendi o requisito e o comportamento esperado?
- [ ] Escrevi testes que descrevem o comportamento?
- [ ] Os testes falham por causa certa (não por erro de sintaxe)?
- [ ] Implementei o mínimo necessário para passar?
- [ ] Refatorei mantendo os testes verdes?
- [ ] Todos os casos de borda estão cobertos?
- [ ] O código está limpo e legível?