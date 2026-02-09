# História de Usuário

Como um usuário, eu quero listar todos os meus registros para ter uma visão panorâmica da minha produção.

## Necessidade de Negócio

O objetivo é permitir que usuários autenticados visualizem todas as suas notas de forma organizada e eficiente, facilitando a localização e gestão do conteúdo criado. A listagem deve considerar apenas as notas do usuário autenticado, oferecer ordenação adequada, suportar paginação para grandes volumes de dados e fornecer informações resumidas de cada nota. Além disso, deve tratar casos de ausência de notas e falhas de carregamento.

## Cenários

### Cenário 01: Listagem com Sucesso

**DADO** que estou autenticado no sistema  
**E** possuo notas cadastradas  
**QUANDO** eu acessar a página de listagem de notas  
**ENTÃO** o sistema deve exibir todas as minhas notas em formato de lista ou grid  
**E** cada item deve mostrar informações resumidas: título (se houver), preview do conteúdo, data de criação e data de última modificação  
**E** as notas devem estar ordenadas por data de criação (mais recentes primeiro) por padrão  
**E** devo poder visualizar o número total de notas cadastradas

### Cenário 02: Listagem com Muitas Notas (Paginação)

**DADO** que estou autenticado no sistema  
**E** possuo mais notas do que o limite de itens por página  
**QUANDO** eu acessar a página de listagem de notas  
**ENTÃO** o sistema deve exibir apenas as primeiras notas (ex: 20 por página)  
**E** deve apresentar controles de paginação (próxima página, página anterior, número da página)  
**E** devo poder navegar entre as páginas para visualizar todas as minhas notas

### Cenário 03: Listagem sem Notas Cadastradas

**DADO** que estou autenticado no sistema  
**E** não possuo nenhuma nota cadastrada  
**QUANDO** eu acessar a página de listagem de notas  
**ENTÃO** o sistema deve exibir uma mensagem amigável: "Você ainda não possui notas. Comece criando sua primeira nota!"  
**E** deve oferecer um botão ou link para criar a primeira nota  
**E** não deve exibir uma lista vazia sem contexto

### Cenário 04: Tentativa de Acesso sem Autenticação

**DADO** que não estou autenticado no sistema  
**QUANDO** eu tentar acessar a página de listagem de notas  
**OU** tentar acessar via API  
**ENTÃO** o sistema deve redirecionar para a página de login  
**OU** retornar erro de autenticação (401 Unauthorized)  
**E** não deve exibir nenhuma nota

### Cenário 05: Ordenação Personalizada

**DADO** que estou autenticado no sistema  
**E** estou na página de listagem de notas  
**QUANDO** eu selecionar uma opção de ordenação diferente (ex: mais antigas primeiro, alfabética, última modificação)  
**ENTÃO** o sistema deve reordenar a lista conforme a opção selecionada  
**E** manter a ordenação selecionada durante a navegação na página  
**E** aplicar a ordenação em todas as páginas da paginação

### Cenário 06: Falha no Carregamento

**DADO** que estou autenticado no sistema  
**QUANDO** ocorrer uma falha na comunicação com o servidor ao carregar a lista  
**OU** ocorrer um erro interno do servidor  
**ENTÃO** o sistema deve exibir uma mensagem de erro amigável: "Não foi possível carregar suas notas. Tente novamente."  
**E** deve oferecer um botão para tentar recarregar a lista  
**E** não deve exibir uma tela em branco ou erro técnico

### Cenário 07: Filtro por Status (Arquivadas/Ativas)

**DADO** que estou autenticado no sistema  
**E** possuo notas ativas e arquivadas  
**QUANDO** eu selecionar o filtro "Apenas ativas" ou "Apenas arquivadas"  
**ENTÃO** o sistema deve exibir apenas as notas que correspondem ao filtro selecionado  
**E** deve atualizar o contador total de notas exibidas  
**E** deve manter o filtro selecionado durante a navegação

## Definition of Done

- Apenas usuários autenticados podem visualizar a lista de notas
- A listagem deve exibir apenas as notas pertencentes ao usuário autenticado
- Cada item da lista deve exibir: preview do conteúdo (primeiros caracteres), data de criação, data de última modificação
- A listagem deve suportar paginação com limite configurável de itens por página (sugestão: 20 itens)
- A ordenação padrão deve ser por data de criação (mais recentes primeiro)
- Deve ser possível ordenar por: data de criação (ascendente/descendente), data de modificação (ascendente/descendente), alfabética (se houver título)
- O sistema deve exibir uma mensagem apropriada quando não houver notas cadastradas
- O sistema deve oferecer filtros para exibir apenas notas ativas ou arquivadas
- Mensagens de erro devem ser claras e não expor detalhes técnicos
- O sistema deve exibir um indicador de carregamento durante a busca das notas
- A listagem deve ser performática mesmo com grande volume de notas (usar paginação)
- O sistema deve garantir que apenas o proprietário da nota possa visualizá-la na lista
- Em caso de falha de conexão, o sistema deve oferecer feedback adequado e permitir nova tentativa
- A listagem deve ser responsiva e funcionar em diferentes tamanhos de tela
