# Start Feature - Practical Guide

## Worktree Setup

```bash
git worktree add ../feature-{name} -b feature-{name}
# Example: git worktree add ../feature-auth-login -b feature/auth-login
```

## References Quick Map

| What | Where | Purpose |
|------|-------|---------|
| Stack & Architecture | `.cursorrules` | Overview |
| Clean Architecture layers | `.cursor/rules/go_backend_architecture.md` | Implementation |
| HTTP handlers | `.cursor/rules/go_backend_ports.md` | Routes |
| Tests | `.cursor/rules/go_backend_tests.md` | Unit tests with mocks |
| Domain design | `.cursor/rules/go_backend_domain.md` + `go_backend_object_calisthenics.md` | Entities & rules |
| App services | `.cursor/rules/go_backend_app.md` | Business logic |
| Product context | `docs/product.md` | User journey & vision |
| Tech decisions | `docs/architecture.md` | Storage, auth, infra |

## Feature Start Checklist

- [ ] Create worktree: `git worktree add ../feature-{name} -b feature-{name}`
- [ ] Understand product: read `docs/product.md`
- [ ] Check tech decisions: read `docs/architecture.md`
- [ ] Design domain entities, VOs, and errors first
- [ ] Implement app layer (services)
- [ ] Implement ports layer (handlers)
- [ ] Write tests BEFORE merging

## Tests Priority

```
domain > app > adapters > ports
```

**Target**: 80% coverage

**Generate mocks**:
```bash
make mock
```

**Run tests**:
```bash
make test
make coverage
```