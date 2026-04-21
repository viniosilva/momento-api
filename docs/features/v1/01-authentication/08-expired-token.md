# Expired Token - Token Expirado

### Status: DONE

## Cenário: Acesso com access token expirado

**DADO** que o access token expirou  
**QUANDO** o usuário tenta acessar a API  
**ENTÃO** o sistema deve exigir refresh token  
**E** negar acesso direto com token expirado

## Descrição do Comportamento Esperado

Quando um usuário tenta acessar a API com um access token que expirou:

1. O sistema valida o token
2. Identifica que o token está expirado
3. Retorna erro HTTP 401 (Unauthorized)
4. A resposta indica que o token expirou
5. O cliente deve usar o refresh token para obter um novo access token

## Notas Técnicas

- O erro pode incluir código específico (ex: `token_expired`) para que o cliente possa tratar apropriadamente
- O cliente deve interceptar este erro e utilizar o refresh token para renovar a sessão
- O sistema não deve aceitar access tokens expirados diretamente na API
- O refresh token endpoint deve verificar a validade do refresh token antes deemitir novo access token