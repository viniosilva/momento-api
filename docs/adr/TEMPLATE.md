# ADR Template

Use this template for documenting new architectural decisions.

Copy and fill in all sections. Keep it concise but thorough.

---

# ADR-XXX: [Decision Title]

**Status**: Proposed | Accepted | Deprecated | Superseded

**Date**: YYYY-MM-DD

## Context

Why was this decision necessary? What problem needed solving?

What were the constraints and requirements? What was the state before?

## Decision

What decision was made? Keep it clear and concise.

Include code examples if helpful.

## Consequences

### Positive ✅

- Benefit 1
- Benefit 2
- Benefit 3

### Negative ❌

- Trade-off 1
- Drawback 1
- Complexity added

## Alternatives Considered

### Alternative 1: [Name]

Brief description. Why was it rejected?

```go
// Code example if applicable
```

- ✅ Pro 1
- ✅ Pro 2
- ❌ Con 1
- ❌ Con 2

### Alternative 2: [Name]

Brief description. Why was it rejected?

```go
// Code example if applicable
```

- ✅ Pros
- ❌ Cons

## Implementation Details

If applicable, provide implementation examples and patterns.

### Example 1

Code or detailed example

```go
// Implementation code
```

### Example 2

More details

```go
// More code
```

## Testing Strategy

If applicable, how to test this decision?

```go
func TestExample(t *testing.T) {
    // Test example
}
```

## Related Decisions

- ADR-XXX: [Description]
- ADR-YYY: [Description]

## References

- Link to article or documentation
- Link to external resource

---

## Guidelines for Filling This Template

### Status
- Start with `Proposed`
- Change to `Accepted` after team consensus
- Use `Deprecated` if no longer valid
- Use `Superseded` if replaced by newer ADR

### Context
- Explain the "why" before the "what"
- Avoid technical jargon — write for new developers
- What were the constraints?
- What was tried before?

### Decision
- 1-3 sentences, crystal clear
- This is the core of the ADR
- Link to code if it exists

### Consequences
- Be honest about both positive and negative impacts
- Mention performance, maintainability, scalability
- Trade-offs are important

### Alternatives
- Show that you considered other options
- Explain why each was rejected
- This justifies the chosen decision

### Implementation
- Provide concrete code examples
- Explain key patterns
- Make it easy to follow the decision

### Related Decisions
- Link to related ADRs
- Show how this fits in the bigger picture

## When to Create an ADR

Create an ADR when:
- Making a major architectural choice
- Choosing a pattern or library that affects many parts of the codebase
- Making a trade-off that deserves documentation
- Deciding against a popular tool or pattern (justify why)

Do NOT create for:
- Trivial changes
- Local refactorings
- Bug fixes

## Process

1. **Write**: Create ADR with Status `Proposed`
2. **Share**: Create pull request for review
3. **Discuss**: Team comments and feedback
4. **Refine**: Update based on feedback
5. **Accept**: Change Status to `Accepted` when consensus reached
6. **Implement**: Make the architectural decision
7. **Update**: If implementation differs, update ADR
8. **Review**: Later, deprecate or supersede if needed

## File Naming

```
ADR-XXX-title-in-kebab-case.md
```

Examples:
- ADR-001-clean-architecture-ddd.md
- ADR-002-http-routing-middleware.md
- ADR-003-cache-strategy.md

---

## Example ADRs

For complete examples, see:

- [ADR-001: Clean Architecture with Domain-Driven Design](./ADR-001-clean-architecture-ddd.md)
- [ADR-002: HTTP Routing, Input Validation & Middleware Chain](./ADR-002-http-routing-middleware.md)

---

**Questions?** See [ADR.md](./ADR.md) or [QUICK.md](./QUICK.md)
