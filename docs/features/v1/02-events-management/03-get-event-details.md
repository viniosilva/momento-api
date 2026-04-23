---
issue: 03
feature: Events Management
group: Get Event Details
bootstrap: 00-initial-setup.md
---

### Status: DONE

# Feature: Get Event Details

## User Story

Como **usuário autenticado**, quiero **ver os detalles de um evento específico**, para que **eu possa obter todas as informações sobre aquele evento**.

---

## Gherkin Scenarios

### Scenario: Usuário vê detalhes de um evento existente

```gherkin
Feature: Get Event Details
  As an authenticated user
  I want to see the details of a specific event
  So that I can get all information about that event

  Scenario: Get event details successfully
    Given I am authenticated as "user@example.com"
    And there is an event with id "event-123" owned by "user@example.com"
    And the event has:
      | field         | value                    |
      | title        | Aniversário da Maria      |
      | description  | Festa de 30 anos         |
      | event_date   | 2026-05-15T20:00:00Z   |
      | status       | upcoming                 |
    When I send a GET request to "/api/events/event-123"
    Then the response status should be 200
    And the response body should contain:
      | field         | value                    |
      | id          | event-123                |
      | title       | Aniversário da Maria      |
      | description | Festa de 30 anos         |
      | event_date  | 2026-05-15T20:00:00Z   |
      | status      | upcoming                 |
      | owner_id    | <current_user_id>        |
```

### Scenario: Usuário tenta ver detalhes de evento inexistente

```gherkin
  Scenario: Get details of non-existent event
    Given I am authenticated as "user@example.com"
    And there is no event with id "non-existent"
    When I send a GET request to "/api/events/non-existent"
    Then the response status should be 404
    And the error message should be "Event not found"
```

### Scenario: Usuário tenta ver detalhes de evento de outro usuário

```gherkin
  Scenario: Get details of another user's event
    Given I am authenticated as "user@example.com"
    And there is an event with id "event-456" owned by "other@example.com"
    When I send a GET request to "/api/events/event-456"
    Then the response status should be 403
    And the error message should be "Access denied"
```

### Scenario: Usuário tenta ver detalhes sem autenticação

```gherkin
  Scenario: Get event details without authentication
    Given I am not authenticated
    When I send a GET request to "/api/events/event-123"
    Then the response status should be 401
    And the error message should be "Unauthorized"
```

---

## Spec Hybrid — Detailed Scenarios

### Spec: Get Event Details — Success

**Given** o usuário está autenticado  
**And** o evento com o ID especificado existe  
**And** o evento pertence ao usuário autenticado (owner)  
**When** envia request GET para `/api/events/{eventId}`  
**Then** o sistema retorna os detalhes completos do evento  

**Response:**
```json
{
  "id": "event-123",
  "owner_id": "user_id",
  "title": "Aniversário da Maria",
  "description": "Festa de 30 anos",
  "event_date": "2026-05-15T20:00:00Z",
  "status": "upcoming",
  "created_at": "2026-04-21T10:00:00Z",
  "updated_at": "2026-04-21T10:00:00Z",
  "deleted_at": null
}
```

### Spec: Get Event Details — Not Found

**Given** o usuário está autenticado  
**And** o evento com o ID especificado não existe  
**When** tenta obter detalhes do evento  
**Then** o sistema retorna 404 Not Found  

**Error Response:**
```json
{
  "error": "Event not found",
  "message": "No event found with id: non-existent"
}
```

### Spec: Get Event Details — Access Denied

**Given** o usuário está autenticado  
**And** o evento existe mas pertence a outro usuário  
**When** tenta obter detalhes do evento  
**Then** o sistema retorna 403 Forbidden  

**Error Response:**
```json
{
  "error": "Access denied",
  "message": "You do not have permission to view this event"
}
```

### Spec: Get Event Details — Unauthorized

**Given** o usuário não está autenticado (token inválido ou ausente)  
**When** tenta obter detalhes do evento  
**Then** o sistema retorna 401 Unauthorized  

---

## Acceptance Criteria

- [ ] Usuário owner consegue ver detalhes do evento
- [ ] Retorna todos os campos do evento (id, title, description, event_date, status, owner_id, timestamps)
- [ ] Retorna 404 se evento não existe
- [ ] Retorna 403 se evento pertence a outro usuário
- [ ] Retorna 401 se usuário não está autenticado
- [ ] owner_id no response deve corresponder ao ID interno do usuário

---

## Dependencies

**Requires:**
- Authentication system (JWT validation)
- User model (to get owner_id)
- Create event feature (to have events to fetch)

**Provides:**
- Event details for update-event feature
- Event details for archive-event feature
- Event details for delete-event feature