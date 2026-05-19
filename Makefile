PATH := $(shell go env GOPATH)/bin:$(PATH)

all:
	go install github.com/vektra/mockery/v3@v3.6.3
	go install github.com/swaggo/swag/cmd/swag@latest
	env GOFLAGS= go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.19.1
	go mod download

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-up-cli:
	migrate -path migrations -database "$(DATABASE_URL)" up $(filter-out $@,$(MAKECMDGOALS))

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down $(filter-out $@,$(MAKECMDGOALS))

migrate-down-all:
	migrate -path migrations -database "$(DATABASE_URL)" down -all

migrate-force:
	@read -p "Enter version to force: " version; \
	migrate -path migrations -database "$(DATABASE_URL)" force $$version

migrate-version:
	migrate -path migrations -database "$(DATABASE_URL)" version

minio-bootstrap:
	docker run --rm \
	--network container:momento-api-minio-1 \
	--entrypoint /bin/sh minio/mc -c "mc alias set local http://localhost:9000 momento_admin momento_admin && mc mb local/momento || true && mc anonymous set public local/momento"

run:
	go run cmd/api/main.go

mock:
	mockery

test:
	go test ./... -cover

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o output.html

swag:
	swag init -g cmd/api/main.go -o docs