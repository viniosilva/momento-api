# Feature Spec — Events Management

## 1. Feature Name
Events Management System (Core Domain)

## 2. Context / Problem
Usuários precisam de uma forma de organizar momentos em grupos.

Sem eventos:
- não existe agrupamento de fotos
- não existe contexto de memória
- não existe estrutura para convites e participantes

Eventos são a unidade central do produto.

## 3. Goal (Value)
Permitir que usuários criem e gerenciem eventos privados onde fotos e participantes são organizados em torno de um momento específico.

## 4. User Value
- Organização clara de momentos
- Separação de eventos distintos (viagens, festas, etc.)
- Controle sobre acesso às memórias
- Facilidade para revisitar eventos passados

## 5. Scope

### Included
- Criar evento
- Listar eventos do usuário
- Ver detalhes do evento
- Editar evento
- Arquivar evento
- Desarquivar evento
- Deletar evento
- Filtros:
  - upcoming events
  - past events (memories)
  - archived events

### Excluded
- Eventos públicos
- Busca pública de eventos
- Tags avançadas ou categorias complexas
- Sub-eventos ou hierarquias

## 6. Functional Requirements
- Usuário autenticado pode criar eventos ilimitados
- Evento deve pertencer a um único owner
- Apenas owner pode editar evento
- Apenas owner pode arquivar/desarquivar evento
- Apenas owner pode deletar evento
- Eventos devem suportar estados:
  - upcoming
  - past (baseado em data)
  - archived
- Sistema deve permitir listar eventos filtrados por estado
- Evento deve conter data e informações básicas (nome, descrição opcional)

## 7. Data Model (if needed)

### Event
- id
- owner_id
- title
- description (optional)
- event_date
- status (upcoming | past | archived)
- created_at
- updated_at
- deleted_at (optional hard delete alternative)

### Relationships
- User (1) → Events (N)
- Event (1) → Invitations (N)
- Event (1) → Photos (N)

## 8. Acceptance Criteria (Definition of Done)
- Usuário consegue criar um evento com sucesso
- Eventos aparecem listados corretamente por usuário
- Eventos podem ser filtrados por status
- Apenas owner consegue editar evento
- Apenas owner consegue arquivar e desarquivar evento
- Apenas owner consegue deletar evento
- Evento reflete corretamente estado (upcoming/past/archived)

## 9. Edge Cases
- Evento sem data definida (bloqueado ou inválido)
- Usuário tentando acessar evento de outro usuário
- Mudança de status automática de upcoming → past
- Deleção de evento com fotos associadas
- Arquivar evento já arquivado
- Desarquivar evento não arquivado

## 10. Technical Notes
- Status derivado parcialmente da event_date (past/upcoming)
- Archived é estado manual separado
- MongoDB collection: events
- Queries devem suportar filtros por owner + status
- Autorização obrigatória em todas as operações
- Soft delete opcional para auditoria futura

## 11. Metrics (optional)
- número de eventos criados por usuário
- taxa de eventos com fotos adicionadas
- distribuição entre upcoming vs past events
- taxa de arquivamento

## 12. Dependencies (optional)
- Authentication system
- User model
- Authorization layer
- Date/time handling utilities