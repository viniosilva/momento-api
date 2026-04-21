# Access Protected Route - Acesso a Rota Protegida

### Status: DONE

## Cenário: Acesso a rota protegida sem autenticação

**DADO** que um usuário não está autenticado  
**QUANDO** ele tenta acessar uma rota protegida  
**ENTÃO** o sistema deve negar acesso com erro 401

## Descrição do Comportamento Esperado

Quando um usuário não autenticado tenta acessar uma rota que requer autenticação:

1. O sistema verifica os dados de autenticação na requisição
2. Identifica que não há token ou o token é inválido
3. Retorna erro HTTP 401 (Unauthorized)
4. Inclui mensagem indicando que autenticação é necessária

## Notas Técnicas

- O sistema deve verificar o token antes de processar a lógica da rota
- O erro 401 indica que autenticação é necessária, não necessariamente credenciais inválidas
- O cabeçalho de resposta pode incluir `WWW-Authenticate` com o esquema de autenticação