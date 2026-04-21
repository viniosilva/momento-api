# Logout - Encerramento de Sessão

### Status: PENDING

## Cenário: Logout do usuário

**DADO** que um usuário está autenticado  
**QUANDO** ele realiza logout  
**ENTÃO** o refresh token deve ser invalidado  
**E** futuras requisições devem exigir novo login

## Descrição do Comportamento Esperado

Quando um usuário autenticado decide fazer logout:

1. O sistema recebe a requisição de logout
2. Invalida o refresh token no banco de dados
3. Remove a sessão ativa
4. O access token continua válido até sua expiração natural

## Notas Técnicas

- O logout pode ser implementado como "single device" (invalida apenas tokens do dispositivo atual) ou "all devices" (invalida todos os tokens do usuário)
- O refresh token deve ser removido ou marcado como inválido no banco de dados
- O sistema deve responder à requisição de logout mesmo que o token já esteja expirado