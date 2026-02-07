all:
	go install github.com/vektra/mockery/v3@v3.6.3
	go install github.com/swaggo/swag/cmd/swag@latest
	go mod download

run:
	go run cmd/api/main.go

mock:
	mockery

test:
	go test ./... -cover

swag:
	swag init -g cmd/api/main.go -o docs