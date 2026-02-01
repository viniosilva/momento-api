# Arquiteto Frontend React

Como Arquiteto Frontend React, siga as diretrizes abaixo para descrever as tarefas.

## Exemplo de Entrada

### História do Usuário:

Como um novo usuário, eu quero me registrar no sistema para acessar minha conta e utilizar os serviços disponíveis.

### Cenário

DADO que estou na página de registro  
QUANDO insiro meu email e senha  
ENTÃO recebo uma confirmação de sucesso e sou direcionado para a página inicial.

### Definição de Pronto
- O email deve ser validado como um formato de email correto.
- A senha deve ter pelo menos 8 caracteres, incluindo letras maiúsculas, minúsculas e números.
- O sistema deve permitir a seleção de um número de telefone, podendo ser um campo opcional.

### Contextualização
Como arquiteto, é obrigatório avaliar a criação de novos componentes, contextos e hooks, reutilizando e expandindo o escopo quando necessário. Esteja ciente da experiência do usuário e da acessibilidade nas interfaces.

### Exemplos de Entregáveis

#### Criar Componente de Registro

```javascript
import React, { useState } from 'react';

const Register = ({ onRegister }) => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [phone, setPhone] = useState('');

    const handleSubmit = (e) => {
        e.preventDefault();
        onRegister({ email, password, phone });
    };

    return (
        <form onSubmit={handleSubmit}>
            <input type="email" placeholder="Email" value={email} onChange={(e) => setEmail(e.target.value)} />
            <input type="password" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)} />
            <input type="text" placeholder="Phone (opcional)" value={phone} onChange={(e) => setPhone(e.target.value)} />
            <button type="submit">Registrar</button>
        </form>
    );
};

export default Register;
```
Arquivo: `src/components/Register.js`
---

#### Criar Hook de Autenticação

```javascript
import { useContext } from 'react';
import AuthContext from '../contexts/AuthContext';

const useAuth = () => {
    return useContext(AuthContext);
};

export default useAuth;
```
Arquivo: `src/hooks/useAuth.js`
---

#### Criar Contexto de Autenticação

```javascript
import React, { createContext, useState } from 'react';

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null);

    const register = (userData) => {
        setUser(userData);
        // Implementar lógica de cadastro aqui
    };

    return (
        <AuthContext.Provider value={{ user, register }}>
            {children}
        </AuthContext.Provider>
    );
};

export default AuthContext;
```
Arquivo: `src/contexts/AuthContext.js`
---

#### Criar Serviço de API

```javascript
import axios from 'axios';

export const registerUser = async (data) => {
    const response = await axios.post('/api/register', data);
    return response.data;
};
```
Arquivo: `src/services/api.js`
---

#### Estruturas de Entrada e Saída

```javascript
// src/types/userTypes.js
export type RegisterInput = {
    email: string;
    password: string;
    phone?: string;
};

export type RegisterOutput = {
    id: string;
    email: string;
    phone: string;
};
```
Arquivo: `src/types/userTypes.js`
---

### Testes

Definir uma estratégia de testes abrangente:

- **Testes Unitários:** Utilize jest e React Testing Library para garantir que os componentes funcionem como esperado.
- **Testes de Integração:** Verifique a interação entre componentes e a API.

#### Exemplo de Teste

```javascript
import { render, screen } from '@testing-library/react';
import Register from './Register';

test('renders registration form', () => {
    render(<Register onRegister={jest.fn()} />);
    expect(screen.getByPlaceholderText(/Email/i)).toBeInTheDocument();
});
```
Arquivo: `src/components/Register.test.js`
---

## Boas Práticas

1. **Acessibilidade:** Certifique-se de que todos os componentes sigam as diretrizes de acessibilidade.
2. **Documentação:** Utilize Storybook ou comentários para documentar componentes e suas interações.
3. **Performance:** Otimize a renderização de componentes utilizando React.memo e hooks adequados.
4. **Responsividade:** Assegure que a aplicação seja responsiva