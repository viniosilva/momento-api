#### Criar Componente de Login

- **Arquivo:** `src/components/Login.js`
- **Descrição:** Componente responsável pela interface de login, incluindo campos para e-mail e senha, e um botão de "Entrar". Deve incluir validação de e-mail e exibir mensagens de erro ao usuário em caso de falha.

#### Criar Hook de Autenticação

- **Arquivo:** `src/hooks/useAuth.js`
- **Descrição:** Hook personalizado que utiliza o contexto de autenticação, permitindo que os componentes acessem e manipulem o estado de autenticação do usuário (login, logout, informações do usuário).

#### Criar Contexto de Autenticação

- **Arquivo:** `src/contexts/AuthContext.js`
- **Descrição:** Contexto que fornece acesso ao estado de autenticação e métodos, como login e logout, para toda a aplicação.

#### Criar Serviço de API para Login

- **Arquivo:** `src/services/api.js`
- **Descrição:** Implementar uma função para realizar chamadas de login à API, enviando as credenciais do usuário (e-mail e senha) e lidando com a resposta, incluindo o gerenciamento do token de autenticação.

#### Estruturas de Entrada e Saída para Login

- **Arquivo:** `src/types/userTypes.js`
- **Descrição:** Definir estruturas de entrada e saída para o login. A entrada deve incluir e-mail e senha, e a saída deve