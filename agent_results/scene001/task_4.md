# História de Usuário: Registro de Notas Rápidas

**História:** Como um usuário, eu quero registrar textos em novas notas para capturar ideias imediatamente.

## 1. Necessidade de Negócio
O objetivo é fornecer uma ferramenta de captura rápida de informações para usuários autenticados. A prioridade é a **agilidade no registro** e a **segurança na persistência**. O sistema deve evitar a perda de dados em casos de falha de rede e garantir que o conteúdo seja armazenado de forma íntegra e associado exclusivamente ao proprietário da conta, impedindo vulnerabilidades como XSS (Cross-Site Scripting).

## 2. Cenários BDD (Dado/Quando/Então)

### Cenário 01: Criação de Nota com Sucesso (Caminho Feliz)
**DADO** que estou autenticado no sistema
**E** estou na página de criação de nota
**QUANDO** eu inserir um texto válido no campo de conteúdo
**E** clicar em "Salvar"
**ENTÃO** o sistema deve persistir a nota associada ao meu `user_id` via transação atômica
**E** deve retornar um status `201 Created`
**E** devo visualizar uma notificação de sucesso (Toast)
**E** ser redirecionado para a tela de edição ou visualização da nota recém-criada.

### Cenário 02: Validação de Conteúdo e Feedback em Tempo Real
**DADO** que estou na página de criação de nota
**QUANDO** o campo de texto estiver vazio ou apenas com espaços
**ENTÃO** o botão "Salvar" deve permanecer desabilitado
**E** ao começar a digitar, o sistema deve exibir um contador de caracteres progressivo (Ex: "150 / 100.000").

### Cenário 03: Bloqueio por Excesso de Caracteres
**DADO** que inseri um texto próximo ao limite de 100.000 caracteres
**QUANDO** eu exceder esse limite
**ENTÃO** o contador deve ficar vermelho
**E** o botão "Salvar" deve ser desabilitado
**E** deve exibir a mensagem: "O conteúdo excede o limite máximo permitido".

### Cenário 04: Proteção contra Injeção (Segurança)
**DADO** que sou um usuário mal-intencionado
**QUANDO** eu inserir scripts ou tags HTML (ex: `<script>`, `<iframe>`) no corpo da nota
**E** tentar salvar
**ENTÃO** o sistema deve sanitizar o conteúdo no backend antes de persistir
**E** ao carregar a nota, o texto deve ser renderizado como string literal, sem executar comandos no navegador.

### Cenário 05: Resiliência em Falhas de Conexão
**DADO** que inseri um conteúdo válido
**QUANDO** clicar em "Salvar" e houver uma queda de conexão ou erro `500` do servidor
**ENTÃO** o sistema deve exibir a mensagem: "Erro ao salvar. Verifique sua conexão e tente novamente"
**E** o conteúdo digitado deve permanecer no campo para que eu não perca o que escrevi.

### Cenário 06: Tentativa de Acesso Não Autenticado
**DADO** que não realizei login
**QUANDO** eu tentar acessar a URL direta de criação de nota ou enviar um POST para o endpoint de notas
**ENTÃO** o sistema deve redirecionar para `/login` (Web)
**OU** retornar erro `401 Unauthorized` (API).

## 3. Definition of Done (DoD)

### Funcional e UX
- [ ] O botão "Salvar" deve exibir estado de *loading* para impedir cliques duplos e duplicidade de IDs.
- [ ] Interface responsiva: o campo de texto deve ser adaptável a diferentes tamanhos de tela.
- [ ] O contador de caracteres deve ser atualizado a cada *keystroke*.

### Técnico e Segurança
- [ ] **ID da Nota:** Gerar obrigatoriamente um UUID v4 no backend.
- [ ] **Data/Hora:** O campo `created_at` deve ser preenchido pelo servidor (UTC).
- [ ] **Autorização:** O backend deve validar se o `user_id` enviado no corpo da requisição condiz com o ID presente no Token JWT.
- [ ] **Sanitização:** Implementar biblioteca de sanitização (ex: DOMPurify ou similar) para prevenir XSS.
- [ ] **Performance:** O endpoint de criação deve responder em menos de 300ms em condições normais.
- [ ] **API:** Retornar o objeto criado no body da resposta e o header `Location` com a URL do novo recurso.
