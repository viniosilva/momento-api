# Sign Up - Cadastro de Usuário

### Status: DONE

## Cenário: Cadastro com email válido

**DADO** que um usuário não possui conta no sistema  
**QUANDO** ele envia um email e senha válidos para cadastro  
**ENTÃO** o sistema deve criar o usuário como não verificado  
**E** enviar um email de verificação

## Descrição do Comportamento Esperado

Quando um novo usuário realiza o cadastro com dados válidos (email no formato correto e senha que atende aos requisitos mínimos), o sistema:

1. Cria um registro do usuário no banco de dados
2. Define o status do usuário como "não verificado" (email não confirmado)
3. Gera um token de verificação de email
4. Envia um email contendo o link/token de verificação

## Notas Técnicas

- O email deve ser único no sistema (não permitir duplicatas)
- A senha deve atender aos requisitos mínimos de segurança (mínimo 8 caracteres)
- O token de verificação possui validade limitada (geralmente 24 horas)
- O usuário não pode fazer login até verificar o email
