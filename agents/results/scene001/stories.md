# História de Usuário
Como um novo usuário, eu quero criar uma conta para iniciar o uso da plataforma.

## Necessidade de Negócio
O objetivo é permitir que novos usuários se registrem de forma segura. É fundamental validar a integridade dos dados (e-mail e senha) e garantir que não existam contas duplicadas, oferecendo uma experiência fluida no "caminho feliz" e orientações claras em caso de erros de preenchimento ou regras de negócio.

### Cenário 01: Cadastro com Sucesso
DADO que estou na página de cadastro e não possuo uma conta ativa
QUANDO eu preencher nome, e-mail válido e uma senha forte
E clicar em "Criar Conta"
ENTÃO o sistema deve persistir meus dados, realizar o login automático e me redirecionar para a tela de boas-vindas.

### Cenário 02: E-mail Já Cadastrado
DADO que insiro um e-mail que já consta na base de dados
QUANDO eu clicar em "Criar Conta"
ENTÃO o sistema deve exibir a mensagem: "Este e-mail já está em uso. Deseja realizar o login ou recuperar sua senha?" e oferecer os links correspondentes.

### Cenário 03: Validação de Dados Inválidos
DADO que preencho o formulário com um e-mail em formato incorreto ou uma senha que não atende aos requisitos de segurança
QUANDO eu tentar submeter o formulário
ENTÃO o sistema deve destacar os campos com erro e exibir mensagens instrutivas (ex: "Senha deve conter ao menos 1 número").

## Definition of Done
- O campo de e-mail deve seguir o padrão RFC 5322 (exemplo@dominio.com).
- A senha deve ter no mínimo 8 caracteres, incluindo letras maiúsculas, minúsculas, números e símbolos.
- O sistema deve impedir a criação de duplicatas (e-mail único).
- Todas as senhas devem ser armazenadas utilizando algoritmos de Hash seguros (ex: BCrypt).
- O botão de envio deve apresentar um estado de "carregando" (loading) para evitar múltiplos cliques acidentais.

## Cenários de testes

### Cenários Positivos
- Deve criar a conta com sucesso quando nome, e-mail e senha atenderem a todos os requisitos
- Deve realizar o login automático imediatamente após o cadastro bem-sucedido
- Deve redirecionar o usuário para a tela de boas-vindas após o registro
- Deve garantir que a senha seja armazenada como Hash (BCrypt) no banco de dados
- Deve aceitar e-mails com domínios longos ou internacionais válidos pela RFC 5322

### Cenários Negativos
- Deve retornar erro quando o e-mail já estiver cadastrado no sistema
- Deve retornar erro quando o campo de e-mail estiver vazio
- Deve retornar erro quando o e-mail não possuir um domínio válido (ex: usuario@dominio)
- Deve retornar erro quando a senha tiver menos de 8 caracteres
- Deve retornar erro quando a senha não contiver ao menos uma letra maiúscula
- Deve retornar erro quando a senha não contiver ao menos uma letra minúscula
- Deve retornar erro quando a senha não contiver ao menos um número
- Deve retornar erro quando a senha não contiver ao menos um símbolo (caractere especial)
- Deve retornar erro quando a senha for composta apenas por espaços em branco
- Deve retornar erro quando o campo de nome estiver vazio
- Deve retornar erro ao tentar submeter o formulário com caracteres de script (XSS) no campo nome
- Deve retornar erro ao tentar realizar SQL Injection no campo de e-mail
- Deve impedir a criação de contas duplicadas ao efetuar cliques rápidos e repetidos no botão de envio
- Deve exibir mensagens de erro específicas abaixo de cada campo com falha de validação
- Deve manter o estado de "loading" no botão de envio enquanto a requisição é processada

---

# História de Usuário
Como um usuário, eu quero realizar o login com e-mail e senha para acessar meu ambiente personalizado.

## Cenário 01: Login com Sucesso
DADO que estou na página de login
QUANDO eu inserir meu e-mail e senha cadastrados corretamente
ENTÃO o sistema deve me redirecionar para o dashboard.

## Cenário 02: E-mail com Formato Inválido
DADO que estou na página de login
QUANDO eu inserir um e-mail sem o formato padrão (ex: "usuario@dominio")
ENTÃO o sistema deve exibir uma mensagem: "Por favor, insira um e-mail válido."

## Cenário 03: Credenciais Incorretas
DADO que estou na página de login
QUANDO eu inserir um e-mail ou senha que não constam na base de dados
ENTÃO o sistema deve exibir a mensagem: "E-mail ou senha incorretos."
E não deve especificar qual dos dois campos está errado por razões de segurança.

## Definition of Done
- Validação de campos obrigatórios antes do envio (client-side)
- Feedback visual claro para erros de validação
- Proteção contra ataques de força bruta (ex: bloqueio temporário após 5 tentativas falhas)
- A senha deve ser mascarada com asteriscos ou pontos

## Cenários de testes

### Cenários Positivos
- Deve redirecionar o usuário para o dashboard ao inserir e-mail e senha corretos
- Deve garantir que a senha seja exibida de forma mascarada (asteriscos ou pontos) durante a digitação
- Deve permitir o login com e-mails que contenham letras maiúsculas, tratando-os como case-insensitive
- Deve aceitar e-mails com domínios longos ou internacionais válidos pela RFC 5322
- Deve ignorar espaços em branco acidentais inseridos antes ou depois do e-mail no formulário

### Cenários Negativos
- Deve retornar erro quando o campo de e-mail estiver vazio
- Deve retornar erro quando o campo de senha estiver vazio
- Deve retornar erro quando o e-mail não possuir um domínio válido (ex: usuario@dominio)
- Deve retornar erro genérico "E-mail ou senha incorretos" quando o e-mail não existir na base
- Deve retornar erro genérico "E-mail ou senha incorretos" quando a senha estiver incorreta para um e-mail válido
- Deve retornar erro ao tentar realizar SQL Injection no campo de e-mail (ex: ' OR 1=1 --)
- Deve retornar erro ao tentar submeter scripts (XSS) no campo de e-mail ou senha
- Deve bloquear o acesso temporariamente após a 5ª tentativa de login consecutiva sem sucesso
- Deve impedir múltiplas requisições ao efetuar cliques rápidos no botão de login (debouncing/loading state)
- Deve exibir mensagens de erro claras abaixo dos campos quando a validação client-side falhar
- Deve garantir que o botão de login permaneça no estado de "loading" enquanto a autenticação é processada
- Deve retornar erro de timeout caso o serviço de autenticação demore a responder além do limite definido

---

# História de Usuário
Como um usuário, eu quero recuperar minha senha via e-mail para retomar o acesso em caso de esquecimento.

## Necessidade de Negócio
O objetivo é fornecer uma forma segura e autônoma para que o usuário recupere o acesso à sua conta. O processo deve garantir que apenas o proprietário do e-mail possa redefinir a senha, protegendo a plataforma contra tentativas de invasão e garantindo que o usuário não fique bloqueado fora do sistema.

### Cenário 01: Solicitação de Recuperação com Sucesso (Caminho Feliz)
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
