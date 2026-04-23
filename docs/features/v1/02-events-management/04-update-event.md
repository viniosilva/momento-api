---
issue: 04
feature: Events Management
group: Update Event
bootstrap: 00-initial-setup.md
---

### Status: DONE

# Feature: Update Event

## User Story

Como **owner de um evento**, quero **editar os detalhes do meu evento**, para que **eu possa corrigir informações incorretas ou atualizar detalhes**.

---

## Gherkin Scenarios

### Scenario: Owner edita evento com dados válidos

```gherkin
Feature: Update Event
  As an event owner
  I want to edit the details of my event
  So that I can correct incorrect information or update details

  Scenario: Update event with valid data
    Given I am authenticated as "user@example.com"
    And there is an event with id "event-123" owned by "user@example.com"
    And the event has:
      | field       | value              |
      | title      | Título Antigo      |
      | event_date | 2026-05-15T20:00Z |
    And I have updated event data:
      | field       | value              |
      | title      | Novo Título       |
      | event_date | 2026-06-01T20:00Z |
    When I send a PUT request to "/api/events/event-123" with the updated data
    Then the response status should be 200
    And the response body should contain:
      | field       | value              |
      | title      | Novo Título       |
      | event_date | 2026-06-01T20:00Z |
    And the event should be updated in the database
    And the updated_at timestamp should be updated
```

### Scenario: Owner atualiza apenas título

```gherkin
  Scenario: Update only event title
    Given I am authenticated as "user@example.com"
    And there is an event with id "event-123" owned by "user@example.com"
    And I have updated event data:
      | field | value          |
      | title | Título Editado |
    When I send a PUT request to "/api/events/event-123" with the updated data
    Then the response status should be 200
    And the response body should contain title "Título Editado"
    And other fields should remain unchanged
```

### Scenario: Non-owner tenta editar evento

```gherkin
  Scenario: Non-owner attempts to edit event
    Given I am authenticated as "user@example.com"
    And there is an event with id "event-456" owned by "other@example.com"
    And I have updated event data:
      | field | value     |
      | title | Hackeado |
    When I send a PUT request to "/api/events/event-456" with the updated data
    Then the response status should be 403
    And the error message should be "Access denied"
```

### Scenario: Owner tenta editar evento inexistente

```gherkin
  Scenario: Owner attempts to update non-existent event
    Given I am authenticated as "user@example.com"
    And there is no event with id "non-existent"
    And I have updated event data:
      | field | value     |
      | title | Não Existe |
    When I send a PUT request to "/api/events/non-existent" with the updated data
    Then the response status should be 404
    And the error message should be "Event not found"
```

### Scenario: Owner tenta editar com título vazio

```gherkin
  Scenario: Owner attempts to update with empty title
    Given I am authenticated as "user@example.com"
    And there is an event with id "event-123" owned by "user@example.com"
    And I have invalid event data:
      | field  | value |
      | title |       |
    When I send a PUT request to "/api/events/event-123" with the updated data
    Then the response status should be 400
    And the error message should be "Title is required"
```

### Scenario: Usuário tenta editar sem autenticação

```gherkin
  Scenario: Update event without authentication
    Given I am not authenticated
    And I have updated event data:
      | field | value     |
      | title | Novo Título |
    When I send a PUT request to "/api/events/event-123" with the updated data
    Then the response status should be 401
    And the error message should be "Unauthorized"
```

---

## Spec Hybrid — Detailed Scenarios

### Spec: Update Event — Success

**Given** o usuário é owner do evento  
**And** possui dados válidos para atualização  
**When** envia request PUT para `/api/events/{eventId}`  
**Then** o sistema atualiza o evento com os novos dados  
**And** atualiza o timestamp `updated_at`  

**Request:**
```json
{
  "title": "Novo Título",
  "description": "Nova descrição",
  "event_date": "2026-06-01T20:00:00Z"
}
```

**Response:**
```json
{
  "id": "event-123",
  "owner_id": "user_id",
  "title": "Novo Título",
  "description": "Nova descrição",
  "event_date": "2026-06-01T20:00:00Z",
  "status": "upcoming",
  "created_at": "2026-04-21T10:00:00Z",
  "updated_at": "2026-04-21T12:00:00Z"
}
```

### Spec: Update Event — Access Denied

**Given** o usuário está autenticado  
**And** o evento pertence a outro usuário  
**When** tenta atualizar o evento  
**Then** o sistema retorna 403 Forbidden  

**Error Response:**
```json
{
  "error": "Access denied",
  "message": "You do not have permission to update this event"
}
```

### Spec: Update Event — Not Found

**Given** o usuário está autenticado  
**And** o evento não existe  
**When** tenta atualizar o evento  
**Then** o sistema retorna 404 Not Found  

**Error Response:**
```json
{
  "error": "Event not found",
  "message": "No event found with id: non-existent"
}
```

### Spec: Update Event — Validation Error

**Given** o usuário é owner do evento  
**When** envia dados inválidos (título vazio)  
**Then** o sistema retorna 400 Bad Request  

**Error Response:**
```json
{
  "error": "Validation failed",
  "details": [
    { "field": "title", "message": "Title is required" }
  ]
}
```

### Spec: Update Event — Unauthorized

**Given** o usuário não está autenticado  
**When** tenta atualizar o evento  
**Then** o sistema retorna 401 Unauthorized  

---

## Acceptance Criteria

- [ ] Owner consegue atualizar título do evento
- [ ] Owner consegue atualizar descrição do evento
- [ ] Owner consegue atualizar event_date do evento
- [ ] Campos não informados permanecem inalterados
- [ ] Retorna 403 se usuário não é owner
- [ ] Retorna 404 se evento não existe
- [ ] Retorna 400 se título está vazio
- [ ] Retorna 401 se usuário não está autenticado
- [ ] Campo updated_at é atualizado automaticamente

---

## Dependencies

**Requires:**
- Authentication system (JWT validation)
- User model (to get owner_id)
- Create event feature (to have events to update)
- Get event details feature (to validate ownership)

**Provides:**
- Updated event data for other features