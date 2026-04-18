# Workflow - Implementation Sequence

## Required Order (layer by layer)

1. **Domain**: Value Objects → Entities → Constants → Tests
2. **App**: DTOs (`dto.go`) → Ports (`port.go`) → Tests (`*_test.go`)
3. **Adapters**: Indexes → Repository → External Services + Tests
4. **App**: Service (`*_service.go`) → Verify tests (`make test`)
5. **Ports**: Request/Response → Port → Handler → Router → Tests
6. **Orchestration**: Inject in `cmd/api/main.go` → `make swag` → `make run`
7. **Finalization**: `make mock` → `make test` → Verify coverage

## Useful Commands
```bash
make          # Install deps
make run      # Run app
make test     # Run tests
make mock     # Generate mocks
make swag     # Generate Swagger
```

## New Feature Checklist

1. Domain (VOs → Entities → Tests)
2. App (DTOs → Ports → Tests)
3. Adapters (Indexes → Repository → External Services + Tests)
4. App (Service → Verify tests)
5. Ports (Req/Res → Handler → Router → Tests)
6. Orchestration (main.go → Swagger)
7. Finalization (mock → test → coverage)