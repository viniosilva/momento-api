---
issue: 00
feature: Events Management
group: Bootstrap/Setup
bootstrap: 00-initial-setup.md
---

### Status: DONE

# Initial Setup — Events Management

Este documento estabelece o bootstrap/base necessário para que todas as issues da feature Events Management possam ser executadas de forma independente. Todas as issues (01-06) dependem diretamente deste setup.

---

## Prerequisites

### Authentication System

**Necessário para todas as issues (1-6):**

- Sistema de autenticação JWT funcional
- Middleware de validação de token
- Endpoint de login/register

```typescript
// Middleware de autenticação esperado
interface AuthMiddleware {
  validateToken(token: string): Promise<User>
  getCurrentUser(request: Request): User | null
}
```

### User Model

**Coleção: `users`**

```json
{
  "_id": "user_uuid",
  "email": "user@example.com",
  "password_hash": "$2b$12$...",
  "created_at": "2026-04-21T10:00:00Z",
  "updated_at": "2026-04-21T10:00:00Z"
}
```

**Indices necessários:**
- `email` (unique)
- `_id` (index)

---

## Database Setup

### Events Collection

**Coleção: `events`**

```json
{
  "_id": "event_uuid",
  "owner_id": "user_uuid",
  "title": "Event Title",
  "description": "Event Description",
  "created_at": "2026-04-21T10:00:00Z",
  "updated_at": "2026-04-21T10:00:00Z",
  "deleted_at": null
}
```

### Indexes

| Index | Fields | Purpose |
|-------|--------|---------|
| `owner_id` | `{ owner_id: 1 }` | Listar eventos do usuário |

---

## Dependency Graph

```
┌─────────────────────────────────────────────────────────────────┐
│                    00-initial-setup                             │
│  ┌─────────────────┐    ┌─────────────────┐                   │
│  │ Authentication  │    │   User Model    │                   │
│  │   (JWT)        │    │   (users)       │                   │
│  └────────┬────────┘    └────────┬────────┘                   │
└──────────┼───────────────────────┼──────────────────────────────┘
           │                       │
           ▼                       ▼
┌─────────────────────────────────────────────────────────────────┐
│                      01-create-event                          │
│         (cria eventos — dados de entrada p/ todas)           │
└──────────────────────────┬──────────────────────────────────┘
                           │
        ┌──────────────────┼──────────────────┐
        ▼                  ▼                  ▼
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│02-list-events  │ │03-get-event-   │ │04-update-event  │
│                 │ │   details      │ │                 │
│ (lista eventos │ │ (obtém evento  │ │ (atualiza evento│
│  criados)      │ │  específico)   │ │  existente)     │
└─────────────────┘ └─────────────────┘ └─────────────────┘
        │                  │                  │
        ▼                  ▼                  ▼
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│05-archive-event │ │05-archive-event │ │05-archive-event │
│                 │ │                 │ │                 │
│ (arquiva       │ │ (arquiva        │ │ (arquiva        │
│  eventos)      │ │  eventos)       │ │  eventos)       │
└─────────────────┘ └─────────────────┘ └─────────────────┘
        │                  │                  │
        ▼                  ▼                  ▼
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│06-delete-event  │ │06-delete-event  │ │06-delete-event  │
│                 │ │                 │ │                 │
│ (deleta        │ │ (deleta         │ │ (deleta        │
│  eventos)      │ │  eventos)       │ │  eventos)       │
└─────────────────┘ └─────────────────┘ └─────────────────┘
```

### Dependency Summary

| Issue | Feature | Depends On |
|-------|---------|------------|
| 00 | Initial Setup | — (base) |
| 01 | Create Event | 00 (auth + user) |
| 02 | List Events | 00 (auth + user) + 01 (eventos existem) |
| 03 | Get Event Details | 00 (auth + user) + 01 (evento criado) |
| 04 | Update Event | 00 (auth + user) + 01 (evento criado) + 03 (valida ownership) |
| 05 | Archive Event | 00 (auth + user) + 01 (evento criado) + 03 (valida ownership) |
| 06 | Delete Event | 00 (auth + user) + 01 (evento criado) + 03 (valida ownership) |

---

## Test Data Strategy

### Minimal Test Data

Para cada teste funcionar independentemente, os seguintes dados devem estar disponíveis:

#### User Fixture (comum a todos)

```json
{
  "_id": "test-user-id",
  "email": "user@example.com",
  "password": "password123"
}
```

```json
{
  "_id": "other-user-id",
  "email": "other@example.com",
  "password": "password456"
}
```

#### Event Fixtures

| Fixture ID | owner_id | title |
|------------|----------|-------|
| `event-123` | test-user-id | Evento Teste |
| `event-456` | other-user-id | Outro Evento |
| `event-past` | test-user-id | Evento Passado |
| `event-archived` | test-user-id | Evento Arquivado |

### Test Execution Order

Cada issue pode ser executada independentemente desde que:

1. **Setup (seed) seja executado** antes dos testes:
   - Criar usuários de teste
   - Criar eventos de teste conforme necessário

2. **Cleanup seja executado** após cada teste:
   - Limpar eventos criados durante o teste
   - Manter usuários base (opcional)

### Independent Execution

**Para executar issue 01 (Create Event):**
- Requer: User fixture
- Não requer: Event fixtures

**Para executar issue 02 (List Events):**
- Requer: User fixture + Event fixtures (upcoming, past, archived)

**Para executar issue 03 (Get Event Details):**
- Requer: User fixture + Event fixture (event-123)

**Para executar issue 04 (Update Event):**
- Requer: User fixture + Event fixture (event-123)

**Para executar issue 05 (Archive Event):**
- Requer: User fixture + Event fixtures (upcoming para archive, archived para unarchive)

**Para executar issue 06 (Delete Event):**
- Requer: User fixture + Event fixture (event-123)

---

## API Endpoints

| Method | Endpoint | Feature |
|--------|----------|---------|
| POST | `/api/auth/register` | Auth |
| POST | `/api/auth/login` | Auth |
| POST | `/api/events` | 01-create-event |
| GET | `/api/events` | 02-list-events |
| GET | `/api/events/:eventId` | 03-get-event-details |
| PUT | `/api/events/:eventId` | 04-update-event |
| POST | `/api/events/:eventId/archive` | 05-archive-event |
| POST | `/api/events/:eventId/unarchive` | 05-archive-event |
| DELETE | `/api/events/:eventId` | 06-delete-event |