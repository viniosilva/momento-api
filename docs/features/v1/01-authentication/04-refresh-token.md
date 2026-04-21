# Refresh Token - Renovação de Sessão

### Status: DONE

## Cenário: Renovação de sessão com refresh token válido

**DADO** que um usuário possui um refresh token válido  
**QUANDO** ele solicita renovação de sessão  
**ENTÃO** o sistema deve emitir um novo access token  
**E** manter a sessão ativa

## Descrição do Comportamento Esperado

Quando o access token expira e o usuário utiliza um refresh token válido:

1. O sistema valida o refresh token
2. Verifica que o token não está revogado ou expirado
3. Gera um novo access token
4. Mantém o refresh token ou rotaciona (dependendo da política de segurança)
5. Retorna o novo access token

## Notas Técnicas

- O refresh token deve ser armazenado no banco de dados para permitir revogação
- O sistema deve validar se o refresh token está ativo antes de emitir novo access token
- Recomenda-se implementar rotação de refresh tokens por segurança
- tokens de refresh devem ter validade mais longa que access tokens