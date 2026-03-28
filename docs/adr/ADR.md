# Architecture Decision Records (ADRs)

Documentation of significant architectural decisions made in the project.

## What is an ADR?

An **Architecture Decision Record** is a document that records a significant architectural decision, its context, rationale, and consequences.

### Why ADRs?

- 📚 **History**: Understand evolution of architecture
- 🧠 **Onboarding**: New developers understand design choices
- 🛡️ **Justification**: Defend decisions against future reviews
- 🔍 **Traceability**: Clear trail of "why" not just "what"

### Standard Format

Each ADR follows this template:

```
# ADR-XXX: [Decision Title]

**Status**: Proposed | Accepted | Deprecated | Superseded

## Context
Why was this decision necessary?

## Decision
What was decided?

## Consequences
### Positive ✅
- Benefit 1
- Benefit 2

### Negative ❌
- Trade-off 1
- Drawback 1

## Alternatives Considered
### Alternative 1: [Name]
Description, why rejected?

### Alternative 2: [Name]
Description, why rejected?

## Related Decisions
- ADR-XXX: [Related title]
```

### Current ADRs

- [ADR-001: Clean Architecture with Domain-Driven Design](./ADR-001-clean-architecture-ddd.md) ✅
- [ADR-002: HTTP Routing, Input Validation & Middleware Chain](./ADR-002-http-routing-middleware.md) ✅

### How to Use ADRs

1. **When to create**: When making significant architectural decisions
2. **Status**: Start with `Proposed`, change to `Accepted` after approval
3. **Numbering**: Increment sequentially (ADR-001, ADR-002, etc)
4. **Format**: Use [TEMPLATE.md](./TEMPLATE.md) as basis
5. **Share**: Include in code review for team discussion

### Status Legend

- ✅ **Accepted** — Decision approved and in active use
- 📝 **Proposed** — Awaiting team discussion/approval
- ⚠️ **Deprecated** — No longer valid, replaced
- 🔄 **Superseded** — Replaced by newer ADR

---

**Next step**: Read [ADR-001](./ADR-001-clean-architecture-ddd.md) or [QUICK.md](./QUICK.md)
