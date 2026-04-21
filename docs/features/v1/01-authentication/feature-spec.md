# Feature Spec — Authentication

## 1. Feature Name
Authentication System

---

## 2. Context / Problem
O sistema precisa garantir que apenas usuários autenticados possam acessar eventos e fotos privadas.

Sem autenticação:
- não há privacidade de eventos
- não há controle de acesso a fotos
- não há identidade de usuário no sistema

Esta é a base de segurança de todo o produto.

---

## 3. Goal (Value)
Permitir que usuários criem contas, façam login seguro e mantenham sessões autenticadas para acessar o sistema de forma contínua.

---

## 4. User Value
- Acesso seguro ao sistema
- Privacidade de dados pessoais e eventos
- Continuidade de sessão sem login constante
- Recuperação de conta em caso de perda de senha

---

## 5. Scope

### Included
- Sign up com email e senha
- Verificação de email
- Login
- Logout
- Refresh token
- Reset password (request + confirm)

### Excluded
- OAuth (Google/Apple)
- Magic link login
- MFA / 2FA
- Login social

---

## 6. Functional Requirements
- Usuário deve poder criar conta com email único
- Senha deve ser armazenada de forma segura (hash)
- Usuário só pode logar após verificar email
- Sistema deve emitir access token após login
- Sistema deve permitir renovação de sessão via refresh token
- Usuário deve poder solicitar reset de senha via email
- Usuário deve poder redefinir senha com token válido
- Tokens devem ter expiração e serem invalidados quando necessário

---

## 7. Data Model

### User
- id
- email
- password_hash
- email_verified
- created_at
- updated_at

### Auth Tokens
- refresh_token (hashado no banco)
- email_verification_token
- reset_password_token
- expiration timestamps

---

## 8. Acceptance Criteria (Definition of Done)

- Usuário consegue se registrar com email válido
- Usuário recebe email de verificação após signup
- Usuário não consegue logar sem email verificado
- Login retorna access token e refresh token válidos
- Refresh token gera novo access token sem novo login
- Logout invalida sessão corretamente
- Reset de senha funciona fim a fim (request + confirm)
- Usuário não autenticado não acessa rotas protegidas

---

## 9. Edge Cases
- Email já cadastrado
- Token de verificação expirado
- Token de reset expirado
- Tentativas de login com senha incorreta
- Múltiplos refresh tokens ativos
- Reuso de token inválido
- Usuário não verificado tentando login

---

## 10. Technical Notes
- JWT para access token
- Refresh token com armazenamento seguro (hash no DB)
- Cookies HTTP-only para sessão
- Middleware de autenticação obrigatório em rotas protegidas
- Serviço de email necessário para verificação e reset
- Expiração curta para access token e longa para refresh token

---

## 11. Metrics (optional)
- Taxa de conversão signup → email verified
- Taxa de login bem-sucedido
- Taxa de reset de senha concluído
- Tempo médio até primeira autenticação

---

## 12. Dependencies (optional)
- Serviço de envio de email
- MongoDB (User + tokens)
- Sistema de logging estruturado