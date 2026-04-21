# Email Verification - Verificação de Email

### Status: PENDING

## Cenário: Verificação de email com token válido

**DADO** que um usuário recebeu um token de verificação válido  
**QUANDO** ele acessa o link de verificação  
**ENTÃO** o sistema deve marcar o usuário como verificado  
**E** permitir login

## Descrição do Comportamento Esperado

Após receber o email de verificação, o usuário acessa o link containing the verification token. O sistema então:

1. Valida o token de verificação
2. Atualiza o status do usuário para "verificado"
3. Confirma o email como verificado no banco de dados
4. Permite que o usuário faça login no sistema

## Notas Técnicas

- O token de verificação é de uso único
- O token possui tempo de expiração (geralmente 24 horas)
- Após verificação bem-sucedida, o token é invalidado
- O sistema deve validar a integridade do token antes de confirmar