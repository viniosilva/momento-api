---
issue: 01
feature: Events Management
group: Create Event
bootstrap: 00-initial-setup.md
---

### Status: DONE

# Feature: Create Event

## User Story

Como **usuário autenticado**, quero **criar um novo evento**, para que **eu possa organizar momentos importantes com photos e participantes**.

---

## Gherkin Scenarios

### Scenario: Usuário cria evento com dados válidos

```gherkin
Feature: Create Event
  As an authenticated user
  I want to create a new event
  So that I can organize important moments

  Scenario: Create event with valid data
    Given I am authenticated as "user@example.com"
    And I have valid event data:
      | field         | value                    |
      | title        | Aniversário da Maria      |
      | description  | Festa de 30 anos         |
      | event_date   | 2026-05-15T20:00:00Z   |
    When I send a POST request to "/api/events" with the event data
    Then the response status should be 201
    And the response body should contain:
      | field         | value                    |
      | title        | Aniversário da Maria      |
      | status       | upcoming                 |
      | owner_id     | <current_user_id>        |
    And the event should be persisted in the database
```

### Scenario: Usuário tenta criar evento sem autenticação

```gherkin
  Scenario: Create event without authentication
    Given I am not authenticated
    And I have valid event data:
      | field       | value              |
      | title       | Evento Teste       |
      | event_date  | 2026-06-01T10:00Z |
    When I send a POST request to "/api/events" with the event data
    Then the response status should be 401
    And the error message should be "Unauthorized"
```

### Scenario: Usuário tenta criar evento sem título

```gherkin
  Scenario: Create event without title
    Given I am authenticated as "user@example.com"
    And I have invalid event data:
      | field        | value              |
      | description  | Sem título        |
      | event_date   | 2026-06-01T10:00Z |
    When I send a POST request to "/api/events" with the event data
    Then the response status should be 400
    And the error message should be "Title is required"
```

### Scenario: Usuário tenta criar evento sem data

```gherkin
  Scenario: Create event without event_date
    Given I am authenticated as "user@example.com"
    And I have invalid event data:
      | field  | value            |
      | title  | Evento Sem Data |
    When I send a POST request to "/api/events" with the event data
    Then the response status should be 400
    And the error message should be "Event date is required"
```

### Scenario: Usuário tenta criar evento com data no passado

```gherkin
  Scenario: Create event with past date
    Given I am authenticated as "user@example.com"
    And I have event data with past date:
      | field       | value                    |
      | title       | Evento Passado          |
      | event_date  | 2020-01-01T00:00:00Z   |
    When I send a POST request to "/api/events" with the event data
    Then the response status should be 201
    And the event status should be "past"
```

---

## Spec Hybrid — Detailed Scenarios

### Spec: Create Event Successfully

**Given** o usuário está autenticado no sistema  
**And** possui dados de evento válidos (título, data, descrição opcional)  
**When** envia request POST para `/api/events`  
**Then** o sistema cria o evento com status padrão "upcoming" (se data futura) ou "past" (se data passada)  
**And** associa o evento ao owner_id do usuário autenticado  
**And** persiste o evento no MongoDB com timestamps de criação  

**Data Model Created:**
```json
{
  "id": "generated_uuid",
  "owner_id": "user_id_from_token",
  "title": "provided_title",
  "description": "provided_or_null",
  "event_date": "2026-05-15T20:00:00Z",
  "status": "upcoming",
  "created_at": "2026-04-21T10:00:00Z",
  "updated_at": "2026-04-21T10:00:00Z",
  "deleted_at": null
}
```

### Spec: Create Event — Validation Errors

**Given** o usuário está autenticado  
**When** envia dados inválidos (faltando título ou data)  
**Then** o sistema retorna 400 Bad Request  
**And** inclui mensagem de erro específica para cada campo inválido  

**Error Response:**
```json
{
  "error": "Validation failed",
  "details": [
    { "field": "title", "message": "Title is required" },
    { "field": "event_date", "message": "Event date is required" }
  ]
}
```

### Spec: Create Event — Unauthorized

**Given** o usuário não está autenticado (token inválido ou ausente)  
**When** tenta criar evento  
**Then** o sistema retorna 401 Unauthorized  

---

## Acceptance Criteria

- [ ] Usuário autenticado consegue criar evento com título e data válida
- [ ] Evento é persistido com owner_id correto
- [ ] Status é automaticamente definido como "upcoming" ou "past" baseado na data
- [ ] Descrição é opcional (campo pode ser null)
- [ ] Retorna 401 se usuário não está autenticado
- [ ] Retorna 400 se título está ausente
- [ ] Retorna 400 se event_date está ausente
- [ ] timestamps (created_at, updated_at) são definidos automaticamente
- [ ] ID único é gerado para o evento

---

## Dependencies

**Requires:**
- Authentication system (JWT validation)
- User model (to get owner_id)

**Provides:**
- Event creation for list-events feature
- Event creation for get-event-details feature