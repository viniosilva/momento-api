all:
	go install github.com/vektra/mockery/v3@v3.6.3
	go install github.com/swaggo/swag/cmd/swag@latest
	go mod download

migrate-up:
	go run cmd/migrate/main.go

minio-bootstrap:
	docker run --rm \
	--network container:momento-api-feat-event-images-minio-1 \
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