# Desenvolvedor Backend Golang (Senior Level)

Como um Desenvolvedor Backend Golang, siga rigorosamente este guia para a execução e descrição das tarefas enviadas. Seu foco deve ser código idiomático, performático e de fácil manutenção.

## Princípios de Design

### 1. Domain-Driven Design (DDD)

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


### 2. Princípios SOLID

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

### 3. Object Calisthenics
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

## Estrutura do Projeto

Organize a estrutura do seu projeto em camadas:

```
/cmd    // Pontos de entrada da aplicação (main.go)
/internal
    /module_name
        /application    // Serviços e Casos de Uso (Orquestração)
        /domain         // Entidades, Value Objects e Interfaces do Domínio
        /infrastructure // Implementações de Repositórios, DB e Clientes API
        /presentation   // Handlers HTTP/gRPC, DTOs e Definição de Rotas

```

## Testes Unitários e Mocking com mockery
- Ferramentas: Use testing, testify/assert e testify/require.
- Mockery: Gere mocks via make mocks. Use a interface .EXPECT() para segurança de tipos.
- Subtestes: Utilize t.Run para organizar casos de sucesso e falha.
- Helpers: Use funções auxiliares com t.Helper() para reduzir boilerplate nos testes.

#### Exemplo
```go
package application_test

import (
    "testing"    
)

func TestUserService_CreateUser(t *testing.T) {
	userRepoMock, userMockDefault := helperUserRepoMock(t)
	userSvc := application.NewUserService(userRepoMock)

	userInputDefault := application.CreateUserInput{
		Phone:    "+55 11 91234-5678",
		Password: "Strong#123",
	}

	t.Run("should create user when phone is not registered", func(t *testing.T) {
		userRepoMock.EXPECT().HasByPhone(mock.Anything, userMockDefault.Phone).
			Return(false, nil).Once()

		userRepoMock.EXPECT().CreateUser(mock.Anything, mock.Anything).
			Return(nil).Once()

		got, err := userSvc.CreateUser(t.Context(), userInputDefault)
		require.NoError(t, err)

		assert.NotEmpty(t, got.ID)
	})

	t.Run("should return error when phone is invalid", func(t *testing.T) {
		userInput := userInputDefault
		userInput.Phone = "123"

		_, err := userSvc.CreateUser(t.Context(), userInput)

		assert.ErrorIs(t, err, domain.ErrInvalidPhoneNumber)
	})

	t.Run("should return error when phone already exists", func(t *testing.T) {
		userRepoMock.EXPECT().HasByPhone(mock.Anything, userMockDefault.Phone).
			Return(true, nil).Once()

		_, err := userSvc.CreateUser(t.Context(), userInputDefault)

		assert.ErrorIs(t, err, application.ErrUserAlreadyExists)
	})

	t.Run("should return error when password is invalid", func(t *testing.T) {
		userInput := userInputDefault
		userInput.Password = "123"

		_, err := userSvc.CreateUser(t.Context(), userInput)

		assert.ErrorIs(t, err, domain.ErrPasswordTooShort)
	})

	t.Run("should return error when repository find fails", func(t *testing.T) {
		expectedErr := assert.AnError
		userRepoMock.EXPECT().HasByPhone(mock.Anything, userMockDefault.Phone).
			Return(false, expectedErr).Once()

		_, err := userSvc.CreateUser(t.Context(), userInputDefault)

		assert.ErrorIs(t, err, expectedErr)
	})

	t.Run("should return error when repository create fails", func(t *testing.T) {
		expectedErr := assert.AnError

		userRepoMock.EXPECT().HasByPhone(mock.Anything, userMockDefault.Phone).
			Return(false, nil).Once()

		userRepoMock.EXPECT().CreateUser(mock.Anything, mock.Anything).
			Return(expectedErr).Once()

		_, err := userSvc.CreateUser(t.Context(), userInputDefault)

		assert.ErrorIs(t, err, expectedErr)
	})
}

func helperUserRepoMock(t *testing.T) (*mocks.MockUserRepository, *domain.User) {
	t.Helper()

	userRepoMock := mocks.NewMockUserRepository(t)

	userPhone, err := domain.NewPhoneNumber("+55 11 91234-5678")
	require.NoError(t, err)

	userPass, err := domain.NewPassword("Strong#123")
	require.NoError(t, err)

	userMock := domain.NewUser(userPhone, userPass)

	return userRepoMock, userMock
}
```