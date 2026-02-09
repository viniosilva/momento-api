# História de Usuário

Como um usuário, eu quero visualizar o conteúdo de notas existentes para consulta de informações.

## Necessidade de Negócio

O objetivo é permitir que usuários autenticados visualizem o conteúdo completo de suas notas de forma clara e legível, facilitando a consulta de informações previamente registradas. O sistema deve garantir que apenas o proprietário da nota (ou usuários com permissão de leitura via compartilhamento) possa visualizar o conteúdo, tratar casos de notas inexistentes ou deletadas, e oferecer uma experiência de leitura adequada. Além disso, deve considerar segurança contra acesso não autorizado e tratamento de erros.

## Cenários

### Cenário 01: Visualização com Sucesso

**DADO** que estou autenticado no sistema  
**E** possuo uma nota cadastrada  
**QUANDO** eu clicar em uma nota na lista  
**OU** acessar diretamente a URL da nota  
**ENTÃO** o sistema deve exibir o conteúdo completo da nota  
**E** deve mostrar informações adicionais: data de criação, data de última modificação, status (ativa/arquivada)  
**E** deve oferecer opções de ação: editar, arquivar, compartilhar, excluir  
**E** o conteúdo deve ser exibido de forma legível e formatada adequadamente

### Cenário 02: Tentativa de Visualizar Nota de Outro Usuário

**DADO** que estou autenticado no sistema  
**E** tento acessar uma nota que pertence a outro usuário  
**QUANDO** eu tentar acessar a nota  
**ENTÃO** o sistema deve retornar erro de autorização (403 Forbidden)  
**OU** exibir mensagem: "Você não tem permissão para visualizar esta nota"  
**E** não deve exibir o conteúdo da nota

### Cenário 03: Nota Não Encontrada

**DADO** que estou autenticado no sistema  
**QUANDO** eu tentar acessar uma nota que não existe  
**OU** que foi deletada permanentemente  
**OU** usar um identificador inválido  
**ENTÃO** o sistema deve exibir uma mensagem: "Nota não encontrada"  
**E** deve oferecer um link para retornar à lista de notas  
**E** deve retornar erro 404 (Not Found) via API

### Cenário 04: Tentativa de Acesso sem Autenticação

**DADO** que não estou autenticado no sistema  
**QUANDO** eu tentar acessar a visualização de uma nota  
**OU** tentar acessar via API  
**ENTÃO** o sistema deve redirecionar para a página de login  
**OU** retornar erro de autenticação (401 Unauthorized)  
**E** não deve exibir o conteúdo da nota

### Cenário 05: Falha no Carregamento

**DADO** que estou autenticado no sistema  
**QUANDO** ocorrer uma falha na comunicação com o servidor ao carregar a nota  
**OU** ocorrer um erro interno do servidor  
**ENTÃO** o sistema deve exibir uma mensagem de erro amigável: "Não foi possível carregar a nota. Tente novamente."  
**E** deve oferecer um botão para tentar recarregar  
**E** não deve exibir uma tela em branco ou erro técnico

### Cenário 06: Visualização de Nota Arquivada

**DADO** que estou autenticado no sistema  
**E** possuo uma nota arquivada  
**QUANDO** eu acessar a visualização da nota arquivada  
**ENTÃO** o sistema deve exibir o conteúdo completo normalmente  
**E** deve indicar visualmente que a nota está arquivada  
**E** deve oferecer opção para desarquivar a nota

## Definition of Done

- Apenas o proprietário da nota ou usuários pode visualizar o conteúdo
- O sistema deve validar a existência da nota antes de exibir o conteúdo
- O conteúdo da nota deve ser exibido de forma segura, prevenindo XSS (sanitização/escape adequado)
- A visualização deve incluir metadados: data de criação, data de última modificação, status
- O sistema deve retornar erro 404 quando a nota não existir
- O sistema deve retornar erro 403 quando o usuário não tiver permissão de visualização
- O sistema deve retornar erro 401 quando o usuário não estiver autenticado
- Mensagens de erro devem ser claras e não expor detalhes técnicos do sistema
- O sistema deve exibir um indicador de carregamento durante a busca da nota
- A visualização deve ser responsiva e funcionar em diferentes tamanhos de tela
- O conteúdo deve ser formatado adequadamente (preservar quebras de linha, espaçamento)
- O sistema deve oferecer navegação de volta para a lista de notas
- Em caso de falha de conexão, o sistema deve oferecer feedback adequado e permitir nova tentativa
- Notas arquivadas devem ser visualizáveis pelo proprietário, com indicação visual do status
