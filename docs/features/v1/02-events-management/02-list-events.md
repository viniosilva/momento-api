---
issue: 02
feature: Events Management
group: List Events
bootstrap: 00-initial-setup.md
---

### Status: DONE

# Feature: List Events

## User Story

Como **usuário autenticado**, quero **listar meus eventos com filtros**, para que **eu possa encontrar eventos específicos rapidamente**.

---

## Gherkin Scenarios

### Scenario: Usuário lista eventos com filtro padrão

```gherkin
Feature: List Events
  As an authenticated user
  I want to list my events with filters
  So that I can find specific events quickly

  Scenario: List events with default filter
    Given I am authenticated as "user@example.com"
    And I have the following events in the database:
      | title           | status    | event_date             |
      | Evento Futuro 1 | upcoming | 2026-12-01T20:00:00Z |
      | Evento Futuro 2 | upcoming | 2026-11-15T20:00:00Z |
      | Evento Passado  | past     | 2020-05-01T20:00:00Z |
    When I send a GET request to "/api/events"
    Then the response status should be 200
    And the response body should contain a list of events
    And the events should be sorted by event_date descending
```

### Scenario: Usuário lista eventos com filtro upcoming

```gherkin
  Scenario: List events with upcoming filter
    Given I am authenticated as "user@example.com"
    And I have "upcoming" events in the database
    When I send a GET request to "/api/events?status=upcoming"
    Then the response status should be 200
    And all returned events should have status "upcoming"
```

### Scenario: Usuário lista eventos com filtro past

```gherkin
  Scenario: List events with past filter
    Given I am authenticated as "user@example.com"
    And I have "past" events in the database
    When I send a GET request to "/api/events?status=past"
    Then the response status should be 200
    And all returned events should have status "past"
```

### Scenario: Usuário lista eventos com filtro archived

```gherkin
  Scenario: List events with archived filter
    Given I am authenticated as "user@example.com"
    And I have "archived" events in the database
    When I send a GET request to "/api/events?status=archived"
    Then the response status should be 200
    And all returned events should have status "archived"
```

### Scenario: Usuário tenta listar eventos sem autenticação

```gherkin
  Scenario: List events without authentication
    Given I am not authenticated
    When I send a GET request to "/api/events"
    Then the response status should be 401
    And the error message should be "Unauthorized"
```

### Scenario: Usuário tenta listar eventos com filtro inválido

```gherkin
  Scenario: List events with invalid filter
    Given I am authenticated as "user@example.com"
    When I send a GET request to "/api/events?status=invalid"
    Then the response status should be 400
    And the error message should be "Invalid status filter. Allowed values: upcoming, past, archived"
```

---

## Spec Hybrid — Detailed Scenarios

### Spec: List Events — Default Behavior

**Given** o usuário está autenticado  
**And** possui eventos no banco de dados  
**When** envia request GET para `/api/events` sem filtro  
**Then** o sistema retorna todos os eventos não-arquivados do usuário  
**And** ordena por event_date decrescente (mais próximos primeiro)  

**Response:**
```json
{
  "events": [
    {
      "id": "uuid-1",
      "title": "Evento Futuro 1",
      "description": "Descrição 1",
      "event_date": "2026-12-01T20:00:00Z",
      "status": "upcoming",
      "created_at": "2026-04-21T10:00:00Z"
    },
    {
      "id": "uuid-2",
      "title": "Evento Futuro 2",
      "description": "Descrição 2",
      "event_date": "2026-11-15T20:00:00Z",
      "status": "upcoming",
      "created_at": "2026-04-20T10:00:00Z"
    }
  ],
  "total": 2
}
```

### Spec: List Events — With Filters

**Given** o usuário está autenticado  
**When** envia request GET para `/api/events?status={filter}`  
**Then** o sistema filtra eventos pelo status especificado  
**And** retorna apenas eventos com aquele status  

**Valid Filters:**
- `upcoming` — eventos com data futura
- `past` — eventos com data passada
- `archived` — eventos arquivados

### Spec: List Events — Invalid Filter

**Given** o usuário está autenticado  
**When** envia.request com filtro de status inválido  
**Then** o sistema retorna 400 Bad Request  
**And** inclui mensagem de erro específica  

**Error Response:**
```json
{
  "error": "Invalid status filter",
  "message": "Allowed values: upcoming, past, archived"
}
```

### Spec: List Events — Unauthorized

**Given** o usuário não está autenticado (token inválido ou ausente)  
**When** tenta listar eventos  
**Then** o sistema retorna 401 Unauthorized  

---

## Acceptance Criteria

- [ ] Usuário autenticado consegue listar todos os seus eventos
- [ ] Eventos são ordenados por event_date decrescente
- [ ] Filtro "upcoming" retorna apenas eventos futuros
- [ ] Filtro "past" retorna apenas eventos passados
- [ ] Filtro "archived" retorna apenas eventos arquivados
- [ ] Retorna 401 se usuário não está autenticado
- [ ] Retorna 400 se filtro de status é inválido
- [ ] Retorna lista vazia se não há eventos para o filtro
- [ ] Inclui campo "total" com quantidade de eventos

---

## Dependencies

**Requires:**
- Authentication system (JWT validation)
- User model (to get owner_id)

**Provides:**
- Event listing for get-event-details feature
- Event listing for update-event feature
- Event listing for archive-event feature
- Event listing for delete-event feature