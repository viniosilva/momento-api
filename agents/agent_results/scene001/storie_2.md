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
