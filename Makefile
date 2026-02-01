run:
	go run cmd/api/main.go

mock:
	mockery

test:
	go test ./... -cover