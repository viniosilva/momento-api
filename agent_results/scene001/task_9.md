# História de Usuário

Como um usuário, eu quero gerar links de compartilhamento para permitir o acesso de terceiros à minha nota.

## Necessidade de Negócio

O objetivo é permitir que usuários autenticados compartilhem suas notas com terceiros através de links únicos e seguros, facilitando a colaboração e o compartilhamento de informações. O sistema deve gerar links com tokens únicos, permitir configuração de permissões (leitura/edição), definir prazos de expiração opcionais e garantir que apenas o proprietário da nota possa gerar links de compartilhamento. Além disso, deve considerar segurança contra acesso não autorizado, rastreamento de compartilhamentos ativos e revogação de links quando necessário.

## Cenários

### Cenário 01: Geração de Link com Sucesso

**DADO** que estou autenticado no sistema  
**E** possuo uma nota cadastrada  
**E** estou visualizando a nota  
**QUANDO** eu clicar no botão "Compartilhar" ou "Gerar Link"  
**E** configurar as permissões (leitura ou edição)  
**E** definir um prazo de expiração (opcional)  
**E** confirmar a geração do link  
**ENTÃO** o sistema deve gerar um link único com token seguro  
**E** devo receber o link completo para compartilhamento  
**E** o link deve ser copiável para a área de transferência com um clique  
**E** o sistema deve registrar o link compartilhado associado à nota  
**E** as permissões e data de expiração devem ser salvas

### Cenário 02: Acesso via Link Compartilhado com Permissão de Leitura

**DADO** que possuo um link de compartilhamento válido com permissão de leitura  
**E** o link não expirou  
**QUANDO** eu acessar o link compartilhado  
**ENTÃO** o sistema deve exibir o conteúdo completo da nota  
**E** deve indicar que estou visualizando via link compartilhado  
**E** não deve permitir edição do conteúdo  
**E** não deve exibir opções de ação que requeiram autenticação (editar, excluir, arquivar)

### Cenário 03: Acesso via Link Compartilhado com Permissão de Edição

**DADO** que possuo um link de compartilhamento válido com permissão de edição  
**E** o link não expirou  
**QUANDO** eu acessar o link compartilhado  
**ENTÃO** o sistema deve exibir o conteúdo completo da nota  
**E** deve permitir edição do conteúdo  
**E** deve indicar que estou editando via link compartilhado  
**E** deve salvar as alterações normalmente

### Cenário 04: Link Compartilhado Expirado

**DADO** que possuo um link de compartilhamento que já expirou  
**QUANDO** eu tentar acessar o link  
**ENTÃO** o sistema deve exibir a mensagem: "Este link de compartilhamento expirou"  
**E** não deve exibir o conteúdo da nota  
**E** deve retornar erro 410 (Gone) ou 404 (Not Found) via API

### Cenário 05: Tentativa de Gerar Link para Nota de Outro Usuário

**DADO** que estou autenticado no sistema  
**E** tento gerar um link para uma nota que pertence a outro usuário  
**QUANDO** eu tentar gerar o link via API  
**OU** tentar acessar a funcionalidade de compartilhamento  
**ENTÃO** o sistema deve retornar erro de autorização (403 Forbidden)  
**OU** exibir mensagem: "Você não tem permissão para compartilhar esta nota"  
**E** não deve gerar nenhum link

### Cenário 06: Revogação de Link Compartilhado

**DADO** que estou autenticado no sistema  
**E** gerei um link de compartilhamento para uma nota  
**E** estou na página de gerenciamento de compartilhamentos  
**QUANDO** eu clicar em "Revogar Link" ou "Desativar Link"  
**E** confirmar a ação  
**ENTÃO** o sistema deve invalidar o link compartilhado  
**E** o link não deve mais permitir acesso à nota  
**E** devo receber confirmação de que o link foi revogado  
**E** tentativas de acesso ao link revogado devem retornar erro apropriado

### Cenário 07: Listagem de Links Compartilhados Ativos

**DADO** que estou autenticado no sistema  
**E** gerei múltiplos links de compartilhamento para minhas notas  
**QUANDO** eu acessar a página de gerenciamento de compartilhamentos  
**OU** visualizar os links ativos de uma nota específica  
**ENTÃO** o sistema deve exibir todos os links compartilhados ativos  
**E** deve mostrar informações: permissões, data de criação, data de expiração (se houver), status  
**E** deve oferecer opção para revogar cada link individualmente

### Cenário 08: Geração de Link sem Prazo de Expiração

**DADO** que estou autenticado no sistema  
**E** possuo uma nota cadastrada  
**QUANDO** eu gerar um link de compartilhamento sem definir prazo de expiração  
**ENTÃO** o sistema deve gerar o link normalmente  
**E** o link deve permanecer válido indefinidamente (até ser revogado manualmente)  
**E** deve indicar na interface que o link não expira

### Cenário 09: Tentativa de Acesso sem Autenticação (Link Público)

**DADO** que possuo um link de compartilhamento válido  
**E** não estou autenticado no sistema  
**QUANDO** eu acessar o link compartilhado  
**ENTÃO** o sistema deve permitir o acesso conforme as permissões do link  
**E** não deve exigir autenticação para visualizar/editar via link compartilhado  
**E** deve indicar que o acesso é via link compartilhado

### Cenário 10: Falha na Geração do Link

**DADO** que estou autenticado no sistema  
**E** estou tentando gerar um link de compartilhamento  
**QUANDO** ocorrer uma falha na comunicação com o banco de dados  
**OU** ocorrer um erro interno do servidor  
**ENTÃO** o sistema deve exibir uma mensagem de erro amigável: "Não foi possível gerar o link. Tente novamente."  
**E** não deve gerar um link parcial ou inválido  
**E** deve permitir nova tentativa

## Definition of Done

- Apenas o proprietário da nota pode gerar links de compartilhamento
- Cada link compartilhado deve ter um token único e seguro (UUID ou hash criptográfico)
- O sistema deve permitir configurar permissões: "Apenas leitura" ou "Leitura e edição"
- O sistema deve permitir definir prazo de expiração opcional para o link
- Links sem prazo de expiração devem permanecer válidos até serem revogados manualmente
- O link gerado deve ser copiável para a área de transferência com facilidade
- O sistema deve registrar todos os links compartilhados ativos associados a cada nota
- O sistema deve permitir revogar links compartilhados individualmente
- Links revogados não devem mais permitir acesso à nota
- Links expirados devem retornar erro apropriado (410 Gone ou 404 Not Found)
- O sistema deve validar a existência da nota antes de gerar o link
- O sistema deve retornar erro 404 quando a nota não existir
- O sistema deve retornar erro 403 quando o usuário não tiver permissão para compartilhar
- O sistema deve retornar erro 401 quando necessário (exceto para acesso via link público válido)
- Mensagens de erro devem ser claras e não expor detalhes técnicos do sistema
- O sistema deve oferecer interface para listar e gerenciar todos os links compartilhados ativos
- A geração do link deve ser atômica (transação)
- O token do link deve ser suficientemente aleatório para prevenir adivinhação
- Links compartilhados devem funcionar sem exigir autenticação do usuário que acessa
- O sistema deve indicar visualmente quando o acesso é via link compartilhado
- Em caso de falha de conexão, o sistema deve oferecer feedback adequado e permitir nova tentativa
- O sistema deve impedir múltiplas gerações simultâneas do mesmo link
- Idealmente, o sistema deve rastrear quantas vezes o link foi acessado (opcional)
