---
issue: 05
feature: Events Management
group: Archive Event
bootstrap: 00-initial-setup.md
---

### Status: DONE

# Feature: Archive Event

## User Story

Como **owner de um evento**, quiero ** arquivar o meu evento**, para que **eu possa remover o evento da lista ativa sem deletá-lo permanentemente**.

---

## Gherkin Scenarios

### Scenario: Owner arquiva evento com sucesso

```gherkin
Feature: Archive Event
  As an event owner
  I want to archive my event
  So that I can remove the event from the active list without deleting it permanently

  Scenario: Archive event successfully
    Given I am authenticated as "user@example.com"
    And there is an event with id "event-123" owned by "user@example.com"
    And the event has status "upcoming"
    When I send a POST request to "/api/events/event-123/archive"
    Then the response status should be 200
    And the response body should contain status "archived"
    And the event should be updated in the database
    And the deleted_at timestamp should be set
```

### Scenario: Owner desarquiva evento com sucesso

```gherkin
  Scenario: Unarchive event successfully
    Given I am authenticated as "user@example.com"
    And there is an event with id "event-123" owned by "user@example.com"
    And the event has status "archived"
    When I send a POST request to "/api/events/event-123/unarchive"
    Then the response status should be 200
    And the response body should contain status "upcoming"
    And the deleted_at timestamp should be set to null
```

### Scenario: Owner tenta arquivar evento já arquivado

```gherkin
  Scenario: Owner attempts to archive already archived event
    Given I am authenticated as "user@example.com"
    And there is an event with id "event-123" owned by "user@example.com"
    And the event has status "archived"
    When I send a POST request to "/api/events/event-123/archive"
    Then the response status should be 400
    And the error message should be "Event is already archived"
```

### Scenario: Owner tenta desarquivar evento não arquivado

```gherkin
  Scenario: Owner attempts to unarchive non-archived event
    Given I am authenticated as "user@example.com"
    And there is an event with id "event-123" owned by "user@example.com"
    And the event has status "upcoming"
    When I send a POST request to "/api/events/event-123/unarchive"
    Then the response status should be 400
    And the error message should be "Event is not archived"
```

### Scenario: Non-owner tenta arquivar evento

```gherkin
  Scenario: Non-owner attempts to archive event
    Given I am authenticated as "user@example.com"
    And there is an event with id "event-456" owned by "other@example.com"
    When I send a POST request to "/api/events/event-456/archive"
    Then the response status should be 403
    And the error message should be "Access denied"
```

### Scenario: Owner tenta arquivar evento inexistente

```gherkin
  Scenario: Owner attempts to archive non-existent event
    Given I am authenticated as "user@example.com"
    And there is no event with id "non-existent"
    When I send a POST request to "/api/events/non-existent/archive"
    Then the response status should be 404
    And the error message should be "Event not found"
```

### Scenario: Usuário tenta arquivar sem autenticação

```gherkin
  Scenario: Archive event without authentication
    Given I am not authenticated
    When I send a POST request to "/api/events/event-123/archive"
    Then the response status should be 401
    And the error message should be "Unauthorized"
```

---

## Spec Hybrid — Detailed Scenarios

### Spec: Archive Event — Success

**Given** o usuário é owner do evento  
**And** o evento tem status "upcoming" ou "past"  
**When** envia request POST para `/api/events/{eventId}/archive`  
**Then** o sistema altera status para "archived"  
**And** define `deleted_at` com timestamp atual  

**Response:**
```json
{
  "id": "event-123",
  "owner_id": "user_id",
  "title": "My Event",
  "status": "archived",
  "deleted_at": "2026-04-21T12:00:00Z"
}
```

### Spec: Unarchive Event — Success

**Given** o usuário é owner do evento  
**And** o evento tem status "archived"  
**When** envia request POST para `/api/events/{eventId}/unarchive`  
**Then** o sistema altera status para "upcoming"  
**And** define `deleted_at` como null  

**Response:**
```json
{
  "id": "event-123",
  "owner_id": "user_id",
  "title": "My Event",
  "status": "upcoming",
  "deleted_at": null
}
```

### Spec: Archive Event — Already Archived

**Given** o usuário é owner do evento  
**And** o evento já tem status "archived"  
**When** tenta arquivar novamente  
**Then** o sistema retorna 400 Bad Request  

**Error Response:**
```json
{
  "error": "Event is already archived",
  "message": "This event has already been archived"
}
```

### Spec: Unarchive Event — Not Archived

**Given** o usuário é owner do evento  
**And** o evento não tem status "archived"  
**When** tenta desarquivar  
**Then** o sistema retorna 400 Bad Request  

**Error Response:**
```json
{
  "error": "Event is not archived",
  "message": "This event is not archived"
}
```

### Spec: Archive Event — Access Denied

**Given** o usuário não é owner do evento  
**When** tenta arquivar o evento  
**Then** o sistema retorna 403 Forbidden  

### Spec: Archive Event — Not Found

**Given** o evento não existe  
**When** tenta arquivar o evento  
**Then** o sistema retorna 404 Not Found  

### Spec: Archive Event — Unauthorized

**Given** o usuário não está autenticado  
**When** tenta arquivar o evento  
**Then** o sistema retorna 401 Unauthorized  

---

## Acceptance Criteria

- [ ] Owner consegue arquivar evento
- [ ] Status muda para "archived"
- [ ] Campo deleted_at é definido com timestamp
- [ ] Owner consegue desarquivar evento
- [ ] Status muda para "upcoming"
- [ ] Campo deleted_at volta a ser null
- [ ] Retorna 400 se evento já está arquivado (archive)
- [ ] Retorna 400 se evento não está arquivado (unarchive)
- [ ] Retorna 403 se usuário não é owner
- [ ] Retorna 404 se evento não existe
- [ ] Retorna 401 se usuário não está autenticado

---

## Dependencies

**Requires:**
- Authentication system (JWT validation)
- User model (to get owner_id)
- Create event feature (to have events to archive)
- Get event details feature (to validate ownership)

**Provides:**
- Archived status for list-events filter
- Archived/unarchived events accessible by owner