# ADR Quick Reference

Quick guide to understand and use ADRs in the project.

## What is an ADR in 30 seconds?

A document that records **WHY** we made an architectural decision, not just WHAT.

- 📚 Keeps history of technical evolution
- 🧠 Helps onboard new developers
- 🛡️ Defends choices against future critique

## Current ADRs

| # | Topic | Decision |
|---|-------|----------|
| 001 | Architecture | 4 layers (Domain, Application, Infrastructure, Presentation) + DDD patterns |
| 002 | HTTP & Validation | `http.ServeMux` (no framework) + custom middleware + domain sanitization |

## ADR Structure

```
ADR-XXX: [Title]
├── Status: Accepted ✅
├── Context: Why?
├── Decision: What?
├── Consequences: Pros/Cons?
├── Alternatives: Why not?
└── Implementation: How?
```

## One-Minute Per ADR

### ADR-001: Clean Architecture + DDD

**Problem**: Need scalable, testable structure

**Solution**: 4 layers (each layer knows only layers below)
```
Presentation
    ↓
Application
    ↓
Domain (zero dependencies)
    ↑
Infrastructure
```

**Patterns**: Value Objects (type-safe validation), Repository Pattern (DB independence)

**Benefits**: Testability, maintainability, independence

**Trade-off**: More boilerplate

---

### ADR-002: HTTP Routing + Middleware + Validation

**Problem**: Need routing, middleware, input validation

**Solution**: 
- Routing: `http.ServeMux` (Go 1.22+, zero dependencies)
- Middleware: Custom chain pattern
- Validation: Structural (middleware) + content (domain)

**Execution order**: Recovery → Logging → Sanitization → Auth → Handler

**Benefits**: No dependencies, clear, secure, XSS prevention

**Trade-off**: Manual middleware setup

---

## Create New ADR

1. **Copy template**
   ```bash
   cp docs/adr/TEMPLATE.md docs/adr/ADR-003-title.md
   ```

2. **Fill fields**: Context, Decision, Consequences, Alternatives

3. **Start with**: `Status: Proposed`

4. **Share**: Team review

5. **Finalize**: Change to `Status: Accepted` after consensus

6. **Update index**: Link in ADR.md

## Decision Flow

```
Architectural Problem
    ↓
Research & Propose (ADR, Status: Proposed)
    ↓
Team Discussion
    ↓
Consensus?
    ├─ Yes → Accept (Status: Accepted)
    └─ No  → Iterate
    ↓
Implement & Document
    ↓
Later?: Deprecate or Supersede with new ADR
```

## Files

- 📖 **[ADR.md](./ADR.md)** — Overview & current decisions
- 🚀 **[QUICK.md](./QUICK.md)** — This file
- 📝 **[TEMPLATE.md](./TEMPLATE.md)** — Template for new ADRs
- 001️⃣ **[ADR-001-clean-architecture-ddd.md](./ADR-001-clean-architecture-ddd.md)** — 4 layers + DDD
- 002️⃣ **[ADR-002-http-routing-middleware.md](./ADR-002-http-routing-middleware.md)** — HTTP + middleware + validation

## Data Flow Example

```
HTTP Request
    ↓
[Middleware] Validate structure (JSON, size, depth)
    ↓
[Handler] Parse request
    ↓
[Service] Call business logic
    ↓
[Domain] NewNoteContent(userInput)
    ├─ Trim whitespace
    ├─ Validate length
    ├─ Remove XSS ← ✅ Guaranteed safe
    └─ Type-safe wrapper
    ↓
[Domain] NewNote(id, safeContent)
    ↓
[Repo] Save to MongoDB
```

## Key Takeaways

| Concept | Benefit |
|---------|---------|
| **4 Layers** | Clear separation, testability |
| **Value Objects** | Type-safe, guaranteed valid |
| **Repository Pattern** | Database independence |
| **Middleware Chain** | Zero dependencies, clear order |
| **Domain Sanitization** | XSS prevention guaranteed |

## Next Steps

1. Read [ADR-001](./ADR-001-clean-architecture-ddd.md) for architecture overview
2. Read [ADR-002](./ADR-002-http-routing-middleware.md) for HTTP & validation details
3. When making architectural decisions, create new ADR using [TEMPLATE.md](./TEMPLATE.md)

---

**Last Updated**: 2026-03-28

For full details, see [ADR.md](./ADR.md)
