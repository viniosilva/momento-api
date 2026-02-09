# História de Usuário
Como um usuário, eu quero recuperar minha senha via e-mail para retomar o acesso em caso de esquecimento.

## Necessidade de Negócio
O objetivo é fornecer uma forma segura e autônoma para que o usuário recupere o acesso à sua conta. O processo deve garantir que apenas o proprietário do e-mail possa redefinir a senha, protegendo a plataforma contra tentativas de invasão e garantindo que o usuário não fique bloqueado fora do sistema.

### Cenário 01: Solicitação de Recuperação com Sucesso
DADO que estou na página de "Esqueci minha senha"
QUANDO eu inserir meu e-mail cadastrado e clicar em "Enviar link de recuperação"
ENTÃO o sistema deve enviar um e-mail com um token único e temporário para redefinição
E exibir uma mensagem informando que as instruções foram enviadas.

### Cenário 02: E-mail Não Cadastrado ou Inválido (Segurança/Exceção)
DADO que insiro um e-mail que não existe na base ou está em formato inválido
QUANDO eu clicar em "Enviar link de recuperação"
ENTÃO o sistema deve exibir a mesma mensagem de sucesso do Cenário 01
E não deve enviar e-mail algum, evitando que atacantes saibam quais e-mails estão na base (Prevenção de Enumeração de Usuários).

### Cenário 03: Redefinição com Token Expirado (Exceção)
DADO que eu clico em um link de recuperação que já expirou (após 24h) ou já foi utilizado
QUANDO a página de nova senha carregar
ENTÃO o sistema deve exibir a mensagem: "Este link de recuperação expirou. Por favor, solicite um novo."

## Definition of Done
- O token de recuperação deve ter validade máxima de 24 horas.
- O link de recuperação deve ser de uso único (invalidar após o primeiro sucesso).
- A nova senha deve seguir os mesmos critérios de força da tela de cadastro (mínimo 8 caracteres, letras, números e símbolos).
- Após a redefinição com sucesso, todos os outros tokens ativos para este usuário devem ser revogados.
- O sistema deve enviar um e-mail de confirmação avisando que a senha foi alterada com sucesso.

## Cenários de testes

### Cenários Positivos
- Deve enviar o e-mail de recuperação quando o e-mail inserido for válido e estiver cadastrado na base
- Deve exibir mensagem de sucesso confirmando o envio das instruções para qualquer e-mail informado (prevenção de enumeração)
- Deve garantir que o token gerado seja único e possua validade de 24 horas
- Deve permitir a redefinição da senha quando o token for válido e a nova senha atender aos requisitos de força
- Deve enviar um e-mail de confirmação ao usuário informando que a alteração de senha foi realizada com sucesso
- Deve invalidar o token imediatamente após o primeiro uso bem-sucedido
- Deve revogar todos os outros tokens de recuperação pendentes após a conclusão de uma troca de senha

### Cenários Negativos
- Deve retornar erro de validação quando o campo de e-mail estiver vazio ou apenas com espaços
- Deve retornar erro quando o formato do e-mail for inválido (ausência de @ ou domínio)
- Deve exibir mensagem de erro específica quando o usuário tentar usar um link/token expirado (após 24h)
- Deve exibir mensagem de erro específica ao tentar usar um link/token que já foi utilizado anteriormente
- Deve impedir a redefinição se a nova senha tiver menos de 8 caracteres
- Deve impedir a redefinição se a nova senha não contiver letras maiúsculas, minúsculas, números e símbolos
- Deve retornar erro ao tentar injetar scripts (XSS) ou comandos SQL no campo de e-mail ou nos campos de nova senha
- Deve impedir o envio excessivo de e-mails de recuperação para o mesmo endereço em um curto período (Rate Limit)
- Deve retornar erro se houver tentativa de submeter a nova senha com o servidor offline ou sem conexão à internet
- Deve impedir que o usuário utilize a mesma senha atual como "nova senha" (se houver essa regra de negócio)
- Deve garantir que o botão de "Enviar" apresente estado de "loading" para evitar múltiplos disparos de e-mail por cliques repetidos
- Deve invalidar o token caso o usuário altere o e-mail da conta antes de utilizar o link de recuperação
