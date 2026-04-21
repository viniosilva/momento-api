# Logout - Encerramento de Sessão

### Status: DONE

## Cenário: Logout do usuário

**DADO** que um usuário está autenticado  
**QUANDO** ele realiza logout  
**ENTÃO** o refresh token deve ser invalidated  
**E** futuras requisições devem exigir novo login

## Descrição do Comportamento Esperado

Quando um usuário autenticado decide fazer logout:

1. O sistema recebe a requisição de logout com o refresh_token
2. Invalida o refresh token no Redis (single device logout)
3. Retorna sucesso
4. O access token continua válido até sua expiração natural

---

## Cenário 2: Logout com token já invalidado

**DADO** que um usuário tenta fazer logout com um token já invalidado  
**QUANDO** o sistema processa a requisição  
**ENTÃO** o logout deve ser bem-sucedido

## Descrição do Comportamento Esperado

Quando um usuário tenta fazer logout com um token que já foi invalidado (ou não existe):

1. O sistema recebe a requisição de logout
2. Tenta invalidar o token no Redis
3. Retorna sucesso mesmo se o token não existe
4. Isso previne enumeração de sessões ativas

## Endpoint

**POST** `/api/auth/logout`

### Request Body
```json
{
  "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2ggdG9rZW4..."
}
```

### Response (204 No Content)
```
(no body)
```

### Errors
- **400**: refresh_token é obrigatório
- **400**: request body inválido

## Notas Técnicas

- Implementado como "single device" - invalida apenas o token do dispositivo atual
- O refresh token é armazenado no Redis e removido na invalidação
- O sistema retorna sucesso mesmo se o token já está invalidado (previne enumeração)
- O access token continua válido até expiração natural (15 minutos)
- Segue o mesmo padrão de implementação dos outros endpoints de autenticação