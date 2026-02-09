# História de Usuário

Como um usuário, eu quero modificar o texto de notas salvas para manter a informação atualizada.

## Necessidade de Negócio

O objetivo é permitir que usuários autenticados editem o conteúdo de suas notas existentes, garantindo que possam atualizar informações, corrigir erros e manter o conteúdo relevante ao longo do tempo. O sistema deve garantir que apenas o proprietário da nota (ou usuários com permissão de edição via compartilhamento) possa modificar o conteúdo, validar a integridade dos dados, atualizar a data de modificação automaticamente e oferecer feedback claro sobre o sucesso ou falha da operação. Além disso, deve considerar tratamento de conflitos de edição simultânea e preservação de dados em caso de falhas.

## Cenários

### Cenário 01: Edição com Sucesso

**DADO** que estou autenticado no sistema  
**E** possuo uma nota cadastrada  
**E** estou visualizando o conteúdo da nota  
**QUANDO** eu clicar no botão "Editar"  
**E** modificar o texto do conteúdo  
**E** clicar em "Salvar" ou "Atualizar"  
**ENTÃO** o sistema deve persistir as alterações no banco de dados  
**E** devo receber uma confirmação visual de que a nota foi atualizada com sucesso  
**E** a data de última modificação deve ser atualizada automaticamente  
**E** devo ser redirecionado para a visualização da nota atualizada  
**OU** a interface deve sair do modo de edição e exibir o conteúdo atualizado

### Cenário 02: Tentativa de Editar Nota de Outro Usuário

**DADO** que estou autenticado no sistema  
**E** tento editar uma nota que pertence a outro usuário  
**E** não possuo permissão de edição via compartilhamento  
**QUANDO** eu tentar acessar a edição da nota  
**OU** tentar submeter alterações via API  
**ENTÃO** o sistema deve retornar erro de autorização (403 Forbidden)  
**OU** exibir mensagem: "Você não tem permissão para editar esta nota"  
**E** não deve permitir a modificação do conteúdo

### Cenário 03: Edição com Conteúdo Vazio

**DADO** que estou autenticado no sistema  
**E** estou editando uma nota  
**QUANDO** eu remover todo o conteúdo da nota  
**OU** deixar apenas espaços em branco  
**E** tentar salvar  
**ENTÃO** o sistema deve destacar o campo de conteúdo com erro  
**E** exibir a mensagem: "O conteúdo da nota não pode estar vazio"  
**E** impedir a atualização até que um conteúdo válido seja inserido  
**E** manter o conteúdo original da nota no banco de dados

### Cenário 04: Conteúdo Excede Limite Máximo

**DADO** que estou autenticado no sistema  
**E** estou editando uma nota  
**QUANDO** eu inserir um texto que exceda o limite máximo de caracteres permitidos  
**E** tentar salvar  
**ENTÃO** o sistema deve exibir uma mensagem informando o limite máximo de caracteres  
**E** mostrar um contador de caracteres indicando quantos foram utilizados  
**E** impedir a atualização até que o conteúdo esteja dentro do limite

### Cenário 05: Cancelamento de Edição

**DADO** que estou autenticado no sistema  
**E** estou editando uma nota  
**E** fiz alterações no conteúdo  
**QUANDO** eu clicar no botão "Cancelar"  
**OU** fechar a página sem salvar  
**ENTÃO** o sistema deve descartar as alterações não salvas  
**E** retornar para a visualização da nota com o conteúdo original  
**E** não deve atualizar a data de modificação  
**E** idealmente deve solicitar confirmação antes de descartar alterações não salvas

### Cenário 06: Falha na Persistência

**DADO** que estou autenticado no sistema  
**E** estou editando uma nota  
**E** fiz alterações válidas no conteúdo  
**QUANDO** ocorrer uma falha na comunicação com o banco de dados  
**OU** ocorrer um erro interno do servidor  
**ENTÃO** o sistema deve exibir uma mensagem de erro amigável: "Não foi possível salvar as alterações. Tente novamente."  
**E** manter as alterações preenchidas no formulário para evitar perda de dados  
**E** permitir nova tentativa de salvamento  
**E** não deve atualizar a data de modificação

### Cenário 07: Nota Não Encontrada Durante Edição

**DADO** que estou autenticado no sistema  
**E** estou editando uma nota  
**QUANDO** a nota for deletada por outro processo durante a edição  
**OU** a nota não existir mais no banco de dados  
**E** eu tentar salvar as alterações  
**ENTÃO** o sistema deve exibir a mensagem: "Esta nota não existe mais ou foi removida"  
**E** deve redirecionar para a lista de notas  
**E** não deve tentar persistir alterações em uma nota inexistente

### Cenário 08: Edição via Link Compartilhado (Permissão de Edição)

**DADO** que possuo um link de compartilhamento válido com permissão de edição  
**E** o link não expirou  
**QUANDO** eu acessar o link e clicar em "Editar"  
**E** modificar o conteúdo  
**E** salvar  
**ENTÃO** o sistema deve permitir a edição e persistir as alterações  
**E** deve atualizar a data de última modificação  
**E** deve indicar que a edição foi feita via link compartilhado (opcional: registrar quem editou)

## Definition of Done

- Apenas o proprietário da nota ou usuários com permissão de edição via link compartilhado podem modificar o conteúdo
- O sistema deve validar a existência da nota antes de permitir a edição
- O conteúdo editado é obrigatório e não pode estar vazio ou conter apenas espaços em branco
- O conteúdo editado deve respeitar o limite máximo de caracteres (mesmo limite da criação)
- A data de última modificação deve ser atualizada automaticamente após cada edição bem-sucedida
- O sistema deve validar e sanitizar o conteúdo para prevenir ataques XSS
- Mensagens de erro devem ser claras e específicas, mas não devem expor detalhes técnicos do sistema
- O botão de salvar deve apresentar estado de "carregando" durante o processo de atualização
- O sistema deve impedir múltiplas submissões simultâneas da mesma edição
- A atualização deve ser persistida de forma atômica (transação)
- O sistema deve retornar erro 404 quando a nota não existir durante a edição
- O sistema deve retornar erro 403 quando o usuário não tiver permissão de edição
- O sistema deve retornar erro 401 quando o usuário não estiver autenticado
- Idealmente, o sistema deve oferecer confirmação antes de descartar alterações não salvas
- Em caso de falha de conexão, o sistema deve oferecer feedback adequado e permitir nova tentativa
- O sistema deve preservar o conteúdo original em caso de falha na atualização
- A interface de edição deve ser intuitiva e clara, diferenciando visualmente o modo de edição do modo de visualização
