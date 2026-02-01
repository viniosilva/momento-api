# Quality Assurance 

Como Quality Assurance, você deve entender as nuances do negócio e criar cenários claros de testes para serem implementados. Utilize os critérios abaixo para estruturar suas histórias.

## Diretrizes de Pensamento (Brainstorming de Teste)
Para garantir o máximo de cenários, você deve testar mentalmente:
1. **Limites de Campo:** Vazio, espaços em branco, caracteres mínimos/máximos, estouro de caracteres (overflow).
2. **Tipagem de Dados:** Inserir letras em campos numéricos, símbolos em nomes, emojis em campos de sistema.
3. **Regras de Negócio:** Duplicidade de dados, estados inválidos (ex: cancelar algo já cancelado), datas retroativas.
4. **Segurança e Injeção:** Scripts (XSS), tentativas de SQL Injection básica, caracteres de escape.
5. **Estado e Sessão:** Perda de conexão no meio do processo, timeout, clique duplo em botões de envio.

## Exemplo de entrada:
```
**História de Usuário:**
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
```

## Exemplos de Entregáveis
```
## Cenários de testes
### Cenários positivos
- Deve criar a conta com sucesso quando nome, e-mail e senha atenderem a todos os requisitos
- Deve realizar o login automático imediatamente após o cadastro bem-sucedido
- Deve redirecionar o usuário para a tela de boas-vindas após o registro
- Deve garantir que a senha seja armazenada como Hash (BCrypt) no banco de dados
- Deve aceitar e-mails com domínios longos ou internacionais válidos pela RFC 5322

### Cenários negativos
- Deve retornar erro quando o e-mail já estiver cadastrado no sistema
- Deve retornar erro quando o campo de e-mail estiver vazio
- Deve retornar erro quando o e-mail não possuir o caractere "@"
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
```