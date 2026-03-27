.PHONY: build test lint run dev dev-up dev-down migrate docker help

BINARY := bin/server
GO     := go

## build: compile the server binary
build:
	$(GO) build -o $(BINARY) ./cmd/server

## test: run all tests with race detector
test:
	$(GO) test ./... -v -race

## lint: run golangci-lint
lint:
	golangci-lint run ./...

## run: build and run the server
run: build
	./$(BINARY)

## dev: start server with hot reload (requires: go install github.com/cosmtrek/air@latest)
dev:
	air -c .air.toml

## dev-up: start local Postgres via docker-compose
dev-up:
	docker compose -f docker-compose.dev.yml up -d

## dev-down: stop local Postgres
dev-down:
	docker compose -f docker-compose.dev.yml down

## migrate: run goose migrations (requires: go install github.com/pressly/goose/v3/cmd/goose@latest)
migrate:
	goose -dir db/migrations postgres "$(DATABASE_URL)" up

## docker: build Docker image
docker:
	docker build -t website-go:latest .

## hash-password: generate bcrypt hash (usage: make hash-password PW=yourpassword)
hash-password:
	@go run ./cmd/hashpw "$(PW)"

## help: list available targets
help:
	@grep -E '^## ' Makefile | sed 's/## //' | column -t -s ':'
