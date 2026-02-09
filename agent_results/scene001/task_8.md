# História de Usuário

Como um usuário, eu quero arquivar registros inativos para organizar meu espaço de trabalho sem excluir dados.

## Necessidade de Negócio

O objetivo é permitir que usuários autenticados arquivem notas que não estão mais em uso ativo, organizando o espaço de trabalho sem perder informações importantes. O arquivamento deve ser reversível, permitindo que o usuário desarquive notas quando necessário. O sistema deve garantir que apenas o proprietário da nota possa arquivar/desarquivar, atualizar o status adequadamente e oferecer filtros para visualizar apenas notas ativas ou arquivadas. Além disso, deve considerar que notas arquivadas não devem aparecer na listagem padrão, mas devem ser acessíveis quando necessário.

## Cenários

### Cenário 01: Arquivamento com Sucesso

**DADO** que estou autenticado no sistema  
**E** possuo uma nota ativa cadastrada  
**E** estou visualizando a nota  
**QUANDO** eu clicar no botão "Arquivar"  
**E** confirmar a ação  
**ENTÃO** o sistema deve atualizar o status da nota para "arquivada" no banco de dados  
**E** devo receber uma confirmação visual de que a nota foi arquivada com sucesso  
**E** a nota não deve mais aparecer na listagem padrão de notas ativas  
**E** a nota deve continuar acessível através do filtro "Apenas arquivadas"  
**E** a data de arquivamento deve ser registrada (opcional)

### Cenário 02: Desarquivamento com Sucesso

**DADO** que estou autenticado no sistema  
**E** possuo uma nota arquivada  
**E** estou visualizando a nota arquivada (via filtro)  
**QUANDO** eu clicar no botão "Desarquivar"  
**E** confirmar a ação  
**ENTÃO** o sistema deve atualizar o status da nota para "ativa" no banco de dados  
**E** devo receber uma confirmação visual de que a nota foi desarquivada com sucesso  
**E** a nota deve voltar a aparecer na listagem padrão de notas ativas  
**E** a nota não deve mais aparecer no filtro "Apenas arquivadas"

### Cenário 03: Tentativa de Arquivar Nota de Outro Usuário

**DADO** que estou autenticado no sistema  
**E** tento arquivar uma nota que pertence a outro usuário  
**QUANDO** eu tentar arquivar a nota via API  
**OU** tentar acessar a funcionalidade de arquivamento  
**ENTÃO** o sistema deve retornar erro de autorização (403 Forbidden)  
**OU** exibir mensagem: "Você não tem permissão para arquivar esta nota"  
**E** não deve alterar o status da nota

### Cenário 04: Arquivamento em Massa

**DADO** que estou autenticado no sistema  
**E** estou na página de listagem de notas  
**E** selecionei múltiplas notas ativas  
**QUANDO** eu clicar na ação "Arquivar selecionadas"  
**E** confirmar a ação  
**ENTÃO** o sistema deve arquivar todas as notas selecionadas  
**E** devo receber uma confirmação indicando quantas notas foram arquivadas  
**E** todas as notas selecionadas devem desaparecer da listagem padrão  
**E** todas devem estar disponíveis no filtro "Apenas arquivadas"

### Cenário 05: Nota Não Encontrada Durante Arquivamento

**DADO** que estou autenticado no sistema  
**E** estou tentando arquivar uma nota  
**QUANDO** a nota for deletada por outro processo antes do arquivamento  
**OU** a nota não existir mais no banco de dados  
**ENTÃO** o sistema deve exibir a mensagem: "Esta nota não existe mais ou foi removida"  
**E** deve redirecionar para a lista de notas  
**E** não deve tentar arquivar uma nota inexistente

### Cenário 06: Falha na Persistência do Arquivamento

**DADO** que estou autenticado no sistema  
**E** estou tentando arquivar uma nota  
**QUANDO** ocorrer uma falha na comunicação com o banco de dados  
**OU** ocorrer um erro interno do servidor  
**ENTÃO** o sistema deve exibir uma mensagem de erro amigável: "Não foi possível arquivar a nota. Tente novamente."  
**E** a nota deve manter seu status original (ativa)  
**E** deve permitir nova tentativa de arquivamento

### Cenário 07: Visualização de Notas Arquivadas

**DADO** que estou autenticado no sistema  
**E** possuo notas arquivadas  
**QUANDO** eu aplicar o filtro "Apenas arquivadas" na listagem  
**ENTÃO** o sistema deve exibir apenas as notas com status arquivado  
**E** deve indicar visualmente que as notas estão arquivadas  
**E** deve oferecer opção para desarquivar individualmente ou em massa

### Cenário 08: Tentativa de Acesso sem Autenticação

**DADO** que não estou autenticado no sistema  
**QUANDO** eu tentar arquivar ou desarquivar uma nota via API  
**ENTÃO** o sistema deve retornar erro de autenticação (401 Unauthorized)  
**E** não deve alterar o status de nenhuma nota

## Definition of Done

- Apenas o proprietário da nota pode arquivar ou desarquivar
- O sistema deve validar a existência da nota antes de alterar o status
- O arquivamento deve ser uma operação reversível (soft delete)
- Notas arquivadas não devem aparecer na listagem padrão (apenas ativas)
- Notas arquivadas devem ser acessíveis através de filtro específico
- O sistema deve oferecer filtro para visualizar "Apenas ativas" e "Apenas arquivadas"
- A data de arquivamento deve ser registrada (opcional, mas recomendado)
- O sistema deve retornar erro 404 quando a nota não existir durante o arquivamento
- O sistema deve retornar erro 403 quando o usuário não tiver permissão para arquivar
- O sistema deve retornar erro 401 quando o usuário não estiver autenticado
- Mensagens de erro devem ser claras e não expor detalhes técnicos do sistema
- O sistema deve oferecer confirmação antes de arquivar (especialmente em massa)
- A operação de arquivamento deve ser atômica (transação)
- O sistema deve suportar arquivamento em massa de múltiplas notas selecionadas
- O sistema deve suportar desarquivamento em massa de múltiplas notas selecionadas
- Notas arquivadas devem manter todos os seus dados (conteúdo, datas, metadados)
- A interface deve indicar visualmente o status arquivado das notas
- Em caso de falha de conexão, o sistema deve oferecer feedback adequado e permitir nova tentativa
- O sistema deve impedir múltiplas submissões simultâneas da mesma operação de arquivamento
