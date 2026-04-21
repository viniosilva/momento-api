# Login - Autenticação de Usuário

### Status: DONE

## Cenário 1: Login com credenciais válidas

**DADO** que um usuário possui conta verificada  
**QUANDO** ele envia email e senha corretos  
**ENTÃO** o sistema deve autenticar o usuário  
**E** retornar access token e refresh token válidos

## Descrição do Comportamento Esperado

Quando um usuário com conta verificada fornece credenciais corretas:

1. O sistema valida as credenciais
2. Gera um access token (JWT) com tempo de expiração curto (ex: 15 minutos)
3. Gera um refresh token com tempo de expiração longo (ex: 7 dias)
4. Retorna ambos os tokens na resposta

---

## Cenário 2: Login com email não verificado

**DADO** que um usuário possui conta mas não verificou o email  
**QUANDO** ele tenta fazer login  
**ENTÃO** o sistema deve negar autenticação  
**E** não emitir tokens

## Descrição do Comportamento Esperado

Quando um usuário tenta fazer login antes de verificar seu email:

1. O sistema valida as credenciais
2. Identifica que o email não foi verificado
3. Retorna erro de autenticação
4. Não emite tokens de sessão
5. Instrui o usuário a verificar seu email

---

## Cenário 3: Login com credenciais incorretas

**DADO** que um usuário envia credenciais incorretas  
**QUANDO** o sistema processa o login  
**ENTÃO** a autenticação deve falhar  
**E** nenhum token deve ser emitido

## Descrição do Comportamento Esperado

Quando um usuário fornece email ou senha incorretos:

1. O sistema tenta validar as credenciais
2. Identifica que as credenciais são inválidas
3. Retorna erro de autenticação
4. Não emite tokens

## Notas Técnicas

- O sistema deve implementar rate limiting para evitar ataques de força bruta
- Não revelar se o email existe ou não no sistema (prevenir enumeração)
- Os tokens devem ser armazenados de forma segura