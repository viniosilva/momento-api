# Workflow - Implementation Sequence

## Required Order (layer by layer)

1. **Domain**: Value Objects → Entities → Constants
2. **App**: DTOs (`dto.go`) → Ports (`port.go`)
3. **Adapters**: Indexes → Repository → External Services
4. **App**: Service (`*_service.go`)
5. **Ports**: Request/Response → Port → Handler → Router
6. **Orchestration**: Inject in `cmd/api/main.go` → `make swag` → `make run`

## Useful Commands
```bash
make          # Install deps
make run      # Run app
make test     # Run tests
make coverage # View coverage report
make mock     # Generate mocks
make swag     # Generate Swagger
```

## See Also
- @.cursor/rules/go_backend_tests.md