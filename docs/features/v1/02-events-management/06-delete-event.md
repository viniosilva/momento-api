---
issue: 06
feature: Events Management
group: Delete Event
bootstrap: 00-initial-setup.md
---

### Status: DONE

# Feature: Delete Event

## User Story

Como **owner de um evento**, quero **deletar permanentemente o meu evento**, para que **eu possa remover o evento completamente do sistema**.

---

## Gherkin Scenarios

### Scenario: Owner deleta evento com sucesso

```gherkin
Feature: Delete Event
  As an event owner
  I want to permanently delete my event
  So that I can completely remove the event from the system

  Scenario: Delete event successfully
    Given I am authenticated as "user@example.com"
    And there is an event with id "event-123" owned by "user@example.com"
    When I send a DELETE request to "/api/events/event-123"
    Then the response status should be 204
    And the event should be permanently removed from the database
```

### Scenario: Non-owner tenta deletar evento

```gherkin
  Scenario: Non-owner attempts to delete event
    Given I am authenticated as "user@example.com"
    And there is an event with id "event-456" owned by "other@example.com"
    When I send a DELETE request to "/api/events/event-456"
    Then the response status should be 403
    And the error message should be "Access denied"
```

### Scenario: Owner tenta deletar evento inexistente

```gherkin
  Scenario: Owner attempts to delete non-existent event
    Given I am authenticated as "user@example.com"
    And there is no event with id "non-existent"
    When I send a DELETE request to "/api/events/non-existent"
    Then the response status should be 404
    And the error message should be "Event not found"
```

### Scenario: Usuário tenta deletar sem autenticação

```gherkin
  Scenario: Delete event without authentication
    Given I am not authenticated
    When I send a DELETE request to "/api/events/event-123"
    Then the response status should be 401
    And the error message should be "Unauthorized"
```

---

## Spec Hybrid — Detailed Scenarios

### Spec: Delete Event — Success

**Given** o usuário é owner do evento  
**When** envia request DELETE para `/api/events/{eventId}`  
**Then** o sistema remove o evento permanentemente do banco de dados  
**And** retorna 204 No Content  

**Response Headers:**
```
HTTP/1.1 204 No Content
```

### Spec: Delete Event — Access Denied

**Given** o usuário está autenticado  
**And** o evento pertence a outro usuário  
**When** tenta deletar o evento  
**Then** o sistema retorna 403 Forbidden  

**Error Response:**
```json
{
  "error": "Access denied",
  "message": "You do not have permission to delete this event"
}
```

### Spec: Delete Event — Not Found

**Given** o usuário está autenticado  
**And** o evento não existe  
**When** tenta deletar o evento  
**Then** o sistema retorna 404 Not Found  

**Error Response:**
```json
{
  "error": "Event not found",
  "message": "No event found with id: non-existent"
}
```

### Spec: Delete Event — Unauthorized

**Given** o usuário não está autenticado (token inválido ou ausente)  
**When** tenta deletar o evento  
**Then** o sistema retorna 401 Unauthorized  

---

## Acceptance Criteria

- [ ] Owner consegue deletar evento permanentemente
- [ ] Evento é removido do banco de dados (hard delete)
- [ ] Retorna 204 No Content em caso de sucesso
- [ ] Retorna 403 se usuário não é owner
- [ ] Retorna 404 se evento não existe
- [ ] Retorna 401 se usuário não está autenticado

---

## Dependencies

**Requires:**
- Authentication system (JWT validation)
- User model (to get owner_id)
- Create event feature (to have events to delete)
- Get event details feature (to validate ownership)

**Provides:**
- Removal of event data from the system