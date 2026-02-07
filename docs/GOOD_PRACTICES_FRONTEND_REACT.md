# Boas Práticas Frontend React

## Princípios de Design


### Comentários

Minimize o máximo possível. Comentários são exceções como explicar uma lógica complexa se o nome da variável ou função não for claro.

**Evite comentários que:**
- Explicam o que o código faz (o código deve ser autoexplicativo)
- Estão desatualizados ou não refletem o código atual
- Duplicam informações já presentes no código
- Explicam o que um componente renderiza (o JSX deve ser autoexplicativo)

**Use comentários para:**
- Explicar o "porquê" de uma decisão técnica complexa
- Documentar workarounds ou limitações de bibliotecas
- Avisar sobre comportamentos inesperados ou side effects
- Documentar props complexas ou tipos TypeScript não triviais

#### Exemplo

```tsx
// ❌ Ruim - comentário desnecessário
// Renderiza o formulário de cadastro
const RegisterForm = () => {
    return <form>...</form>;
};

// ✅ Bom - código autoexplicativo
const RegisterForm = () => {
    return <form>...</form>;
};

// ✅ Bom - explica o porquê de uma decisão técnica
// Usa useMemo para evitar recálculo da regex em cada render
// Referência: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Regular_Expressions
const emailRegex = useMemo(() => /^[^\s@]+@[^\s@]+\.[^\s@]+$/, []);
```

### Component-Based Architecture

- **Componentes**: Crie componentes reutilizáveis que representem unidades da interface do usuário (ex: `Button`, `Input`, `UserCard`).
- **Separação de Responsabilidades**: Separe componentes de apresentação (dumb components) de componentes de lógica (smart components).
- **Composição**: Prefira composição sobre herança. Use children, render props ou compound components.
- **Single Source of Truth**: Mantenha o estado em um único lugar. Use Context API ou gerenciadores de estado para estado global.

#### Exemplo

```tsx
// Componente de apresentação (dumb component)
interface ButtonProps {
    children: React.ReactNode;
    onClick: () => void;
    variant?: 'primary' | 'secondary';
    disabled?: boolean;
}

const Button = ({ children, onClick, variant = 'primary', disabled = false }: ButtonProps) => {
    return (
        <button 
            onClick={onClick} 
            className={`btn btn-${variant}`}
            disabled={disabled}
        >
            {children}
        </button>
    );
};

// Componente de lógica (smart component)
const RegisterForm = () => {
    const [email, setEmail] = useState('');
    const { register, isLoading } = useAuth();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        await register({ email });
    };

    return (
        <form onSubmit={handleSubmit}>
            <Input 
                type="email" 
                value={email} 
                onChange={(e) => setEmail(e.target.value)} 
            />
            <Button type="submit" disabled={isLoading}>
                Cadastrar
            </Button>
        </form>
    );
};
```

### Princípios SOLID

- **S - Single Responsibility Principle**: Cada componente deve ter uma única responsabilidade. Se um componente faz muitas coisas, quebre-o em componentes menores.
- **O - Open/Closed Principle**: Componentes devem ser abertos para extensão, mas fechados para modificação. Use props, render props ou compound components.
- **L - Liskov Substitution Principle**: Componentes filhos devem ser substituíveis por componentes pais sem quebrar a aplicação.
- **I - Interface Segregation Principle**: Evite componentes que forcem a implementação de props desnecessárias. Prefira props opcionais e interfaces pequenas.
- **D - Dependency Inversion Principle**: Dependa de abstrações (props, interfaces) e não de implementações concretas. Injete dependências via props.

#### Exemplo

```tsx
// Interface segregada - props opcionais e específicas
interface UserCardProps {
    user: User;
    showEmail?: boolean;
    showPhone?: boolean;
    onEdit?: (user: User) => void;
}

const UserCard = ({ user, showEmail = false, showPhone = false, onEdit }: UserCardProps) => {
    return (
        <div>
            <h3>{user.name}</h3>
            {showEmail && <p>{user.email}</p>}
            {showPhone && <p>{user.phone}</p>}
            {onEdit && <Button onClick={() => onEdit(user)}>Editar</Button>}
        </div>
    );
};
```

### React Best Practices

- **Um componente, uma responsabilidade**: Agrupe comportamentos relacionados de forma coesa.
- **Early Returns**: Use retornos antecipados para evitar aninhamento excessivo e melhorar legibilidade.
- **Componentes Pequenos**: Mantenha os componentes curtos, focados e com nomes descritivos. Se um componente tem mais de 200 linhas, considere quebrá-lo.
- **Baixa Indentação**: Minimize níveis de indentação. Se o JSX estiver muito aninhado, extraia para componentes ou funções auxiliares.
- **Hooks Customizados**: Extraia lógica reutilizável para hooks customizados. Use `use` como prefixo (ex: `useAuth`, `useForm`).
- **Rule: "Composition over Configuration"**: Prefira composição de componentes sobre configuração complexa via props.

#### Exemplo

```tsx
// ❌ Ruim - componente grande com muita lógica
const UserProfile = ({ user }) => {
    const [isEditing, setIsEditing] = useState(false);
    const [name, setName] = useState(user.name);
    const [email, setEmail] = useState(user.email);
    // ... mais 50 linhas de código

    return (
        <div>
            {/* JSX muito aninhado */}
        </div>
    );
};

// ✅ Bom - componentes pequenos e focados
const UserProfile = ({ user }: { user: User }) => {
    const { isEditing, toggleEdit } = useEditMode();
    const { formData, handleChange, handleSubmit } = useUserForm(user);

    if (!user) {
        return <EmptyState message="Usuário não encontrado" />;
    }

    return (
        <Card>
            {isEditing ? (
                <UserEditForm 
                    user={user} 
                    onSubmit={handleSubmit}
                    onCancel={toggleEdit}
                />
            ) : (
                <UserView user={user} onEdit={toggleEdit} />
            )}
        </Card>
    );
};

// Hook customizado para lógica de edição
const useEditMode = () => {
    const [isEditing, setIsEditing] = useState(false);
    
    const toggleEdit = () => setIsEditing(prev => !prev);
    
    return { isEditing, toggleEdit };
};
```

## Testes Unitários e Mocking

- Ferramentas: Use Jest, React Testing Library e @testing-library/user-event.
- Mocks: Use `jest.mock()` para mockar módulos e `jest.fn()` para mockar funções.
- Subtestes: Utilize `describe` e `it` para organizar casos de sucesso e falha.
- Helpers: Use funções auxiliares e `render` customizado para reduzir boilerplate nos testes.

#### Exemplo

**Nota:** Este exemplo mostra a estrutura completa de testes. Na prática, seguindo TDD, você testaria apenas o caminho feliz e os cenários de erro por vez.

```tsx
// src/components/RegisterForm.test.tsx
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { RegisterForm } from './RegisterForm';
import * as authService from '../../services/auth.service';

jest.mock('../../services/auth.service');

describe('RegisterForm', () => {
    const mockRegister = jest.fn();
    
    beforeEach(() => {
        jest.clearAllMocks();
        (authService.registerUser as jest.Mock) = mockRegister;
    });

    // Caminho feliz - sempre testar
    it('should register user successfully when form is valid', async () => {
        const user = userEvent.setup();
        mockRegister.mockResolvedValue({ id: '1', email: 'test@example.com' });

        render(<RegisterForm />);

        await user.type(screen.getByLabelText(/email/i), 'test@example.com');
        await user.type(screen.getByLabelText(/senha/i), 'Strong#123');
        await user.click(screen.getByRole('button', { name: /cadastrar/i }));

        await waitFor(() => {
            expect(mockRegister).toHaveBeenCalledWith({
                email: 'test@example.com',
                password: 'Strong#123',
            });
        });
    });

    // Apenas um cenário de erro representativo (seguindo TDD)
    it('should display error message when registration fails', async () => {
        const user = userEvent.setup();
        mockRegister.mockRejectedValue(new Error('Email já cadastrado'));

        render(<RegisterForm />);

        await user.type(screen.getByLabelText(/email/i), 'existing@example.com');
        await user.type(screen.getByLabelText(/senha/i), 'Strong#123');
        await user.click(screen.getByRole('button', { name: /cadastrar/i }));

        await waitFor(() => {
            expect(screen.getByText(/email já cadastrado/i)).toBeInTheDocument();
        });
    });
});
```

### Test-Driven Development (TDD)

TDD (Test-Driven Development) é uma prática de desenvolvimento onde você escreve os testes **antes** de implementar a funcionalidade. O ciclo TDD segue três etapas: **Red → Green → Refactor**.

**Components Layer (Componentes de UI):**
- Teste primeiro a renderização e interações do usuário a fim que o resultado seja de sucesso
- Garanta que validações de formulário falhem rápido (Fail-Fast) com dados inválidos
- Teste comportamentos visíveis, não implementação interna
- Teste acessibilidade (labels, ARIA, navegação por teclado)

**Hooks Layer (Hooks Customizados):**
- Teste casos de uso completos que o resultado seja de sucesso
- Use mocks para dependências externas (APIs, contextos)
- Teste tanto o caminho feliz quanto os casos de erro
- Para os casos de erro, apenas um cenário de erro que possibilite testar é só suficiente.
  Ex: um teste que o hook falhe ao tentar fazer login, não precisa testar todos os casos de falha, um somente.

**Services Layer (Chamadas de API):**
- Não testar diretamente, iremos testar com testes de integração ou através dos hooks que os utilizam

**Contexts Layer (Context API):**
- Teste o comportamento do provider e consumo do contexto
- Teste atualizações de estado e side effects
- Use mocks para serviços externos
- Para os casos de erro, apenas um cenário de erro que possibilite testar é só suficiente.

**Pages Layer (Páginas):**
- Teste integração entre componentes e roteamento
- Teste fluxos completos do usuário
- Foco em um teste para cada fluxo principal da página


#### Boas Práticas TDD em React

1. **Um teste, uma responsabilidade**: Cada teste deve verificar um comportamento específico
2. **Nomes descritivos**: Use `should [comportamento] when [condição]` nos nomes dos testes
3. **Arrange-Act-Assert**: Estruture testes em três fases claras
4. **Testes independentes**: Cada teste deve poder rodar isoladamente
5. **Mocks apenas quando necessário**: Prefira dependências reais quando possível
6. **Teste comportamentos, não implementação**: Foque no "o quê" o usuário vê/faz, não no "como" está implementado
7. **Teste acessibilidade**: Use `getByRole`, `getByLabelText` ao invés de `getByTestId`

#### Ferramentas para TDD em React

- **Jest**: Framework de testes
- **React Testing Library**: Biblioteca para testar componentes focando no comportamento do usuário
- **@testing-library/user-event**: Simulação de interações do usuário
- **@testing-library/jest-dom**: Matchers adicionais para DOM
- **MSW (Mock Service Worker)**: Mock de requisições HTTP para testes de integração
- **Vitest**: Alternativa moderna ao Jest (mais rápido)

#### Comandos Úteis

```bash
# Rodar testes em modo watch
npm test -- --watch

# Rodar testes com cobertura
npm test -- --coverage

# Rodar testes com cobertura detalhada
npm test -- --coverage --coverageReporters=html

# Rodar testes específicos
npm test -- RegisterForm

# Rodar testes em modo verbose
npm test -- --verbose

# Rodar testes em paralelo (padrão do Jest)
npm test -- --maxWorkers=4
```

#### Checklist TDD

Antes de começar a implementar:
- [ ] Entendi o requisito e o comportamento esperado?
- [ ] Escrevi testes que descrevem o comportamento do usuário?
- [ ] Os testes falham por causa certa (não por erro de sintaxe)?
- [ ] Implementei o mínimo necessário para passar?
- [ ] Refatorei mantendo os testes verdes?
- [ ] Todos os casos de borda estão cobertos?
- [ ] O componente está acessível (labels, ARIA, navegação por teclado)?
- [ ] O código está limpo e legível?
