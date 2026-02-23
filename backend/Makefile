.PHONY: build run migrate-up migrate-down clean test

include .env
export

BINARY_NAME=budget-tracker

build:
	go build -o $(BINARY_NAME) ./cmd/main.go

run: build
	./$(BINARY_NAME)

migrate-up:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" up

migrate-down:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" down

clean:
	rm -f $(BINARY_NAME)

test:
	go test ./...
