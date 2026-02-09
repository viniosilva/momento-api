# História de Usuário

Como um usuário, eu quero registrar textos em novas notas para capturar ideias imediatamente.

## Necessidade de Negócio

O objetivo é permitir que usuários autenticados criem novas notas de forma rápida e intuitiva, capturando ideias e informações importantes no momento em que surgem. O sistema deve garantir que apenas usuários autenticados possam criar notas, validar a integridade dos dados inseridos e oferecer feedback claro sobre o sucesso ou falha da operação. Além disso, deve considerar limites de tamanho de texto e tratamento de erros de conexão ou persistência.

## Cenários

### Cenário 01: Criação de Nota com Sucesso

**DADO** que estou autenticado no sistema  
**E** estou na página de criação de nota  
**QUANDO** eu inserir um texto válido no campo de conteúdo  
**E** clicar em "Salvar"  
**ENTÃO** o sistema deve persistir a nota no banco de dados associada ao meu usuário  
**E** devo receber uma confirmação visual de que a nota foi criada com sucesso  
**E** ser redirecionado para a visualização da nota criada  
**E** a nota deve ter um identificador único gerado automaticamente  
**E** a nota deve ter data e hora de criação registradas automaticamente

### Cenário 02: Tentativa de Criar Nota sem Autenticação

**DADO** que não estou autenticado no sistema  
**QUANDO** eu tentar acessar a página de criação de nota  
**OU** tentar submeter uma nota via API  
**ENTÃO** o sistema deve redirecionar para a página de login  
**OU** retornar erro de autenticação (401 Unauthorized)  
**E** não deve criar nenhuma nota

### Cenário 03: Conteúdo Vazio ou Inválido

**DADO** que estou autenticado no sistema  
**E** estou na página de criação de nota  
**QUANDO** eu tentar salvar uma nota com conteúdo vazio  
**OU** contendo apenas espaços em branco  
**ENTÃO** o sistema deve destacar o campo de conteúdo com erro  
**E** exibir a mensagem: "O conteúdo da nota não pode estar vazio"  
**E** impedir a criação da nota até que um conteúdo válido seja inserido

### Cenário 04: Conteúdo Excede Limite Máximo

**DADO** que estou autenticado no sistema  
**E** estou na página de criação de nota  
**QUANDO** eu inserir um texto que exceda o limite máximo de 100.000 caracteres  
**E** tentar salvar a nota  
**ENTÃO** o sistema deve exibir uma mensagem informando o limite máximo: "O conteúdo da nota não pode exceder 100.000 caracteres"  
**E** mostrar um contador de caracteres indicando quantos foram utilizados (ex: "100.001 / 100.000")  
**E** impedir a criação da nota até que o conteúdo esteja dentro do limite

### Cenário 05: Falha na Persistência

**DADO** que estou autenticado no sistema  
**E** estou na página de criação de nota  
**E** inseri um conteúdo válido  
**QUANDO** ocorrer uma falha na comunicação com o banco de dados  
**OU** ocorrer um erro interno do servidor  
**ENTÃO** o sistema deve exibir uma mensagem de erro amigável: "Não foi possível salvar a nota. Tente novamente."  
**E** manter o conteúdo preenchido no formulário para evitar perda de dados  
**E** permitir nova tentativa de salvamento

## Definition of Done

- Apenas usuários autenticados podem criar notas
- O conteúdo da nota é obrigatório e não pode estar vazio ou conter apenas espaços em branco
- O conteúdo da nota deve ter um limite máximo de 100.000 caracteres (aproximadamente 100-200 KB em UTF-8, permitindo notas longas sem comprometer performance)
- Cada nota deve ter um identificador único (UUID)
- Cada nota deve estar associada ao usuário que a criou
- Cada nota deve ter data e hora de criação registradas automaticamente
- O sistema deve validar e sanitizar o conteúdo para prevenir ataques XSS
- Mensagens de erro devem ser claras e específicas, mas não devem expor detalhes técnicos do sistema
- O botão de salvar deve apresentar estado de "carregando" durante o processo de criação
- O sistema deve impedir múltiplas submissões simultâneas da mesma nota
- A nota criada deve ser persistida de forma atômica (transação)
- O sistema deve retornar o identificador da nota criada após o sucesso
- Em caso de falha de conexão, o sistema deve oferecer feedback adequado e permitir nova tentativa
