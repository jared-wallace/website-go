.PHONY: build test lint run dev dev-up dev-down migrate docker hash-password deps check-deps deploy logs status help

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

## migrate: run goose migrations up
migrate:
	$(GO) run github.com/pressly/goose/v3/cmd/goose@latest -dir db/migrations postgres "$(DATABASE_URL)" up

## migrate-down: roll back last goose migration
migrate-down:
	$(GO) run github.com/pressly/goose/v3/cmd/goose@latest -dir db/migrations postgres "$(DATABASE_URL)" down

## migrate-status: show migration status
migrate-status:
	$(GO) run github.com/pressly/goose/v3/cmd/goose@latest -dir db/migrations postgres "$(DATABASE_URL)" status

## docker: build Docker image
docker:
	docker build -t website-go:latest .

## hash-password: generate bcrypt hash (usage: make hash-password PW=yourpassword)
hash-password:
	@go run ./cmd/hashpw "$(PW)"

## deps: install optional dev tools (golangci-lint, air)
deps:
	@echo "Installing golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0 2>/dev/null && echo "  ✓ golangci-lint" || echo "  ✗ golangci-lint (install manually)"
	@echo "Installing air (hot reload)..."
	@go install github.com/cosmtrek/air@latest 2>/dev/null && echo "  ✓ air" || echo "  ✗ air (install manually)"
	@echo "Done. goose runs via 'go run' — no install needed."

## check-deps: verify required tools are available
check-deps:
	@echo "Checking dependencies..."
	@command -v go >/dev/null 2>&1 && echo "  ✓ go ($$(go version | awk '{print $$3}'))" || echo "  ✗ go (required)"
	@command -v docker >/dev/null 2>&1 && echo "  ✓ docker" || echo "  ✗ docker (required for dev-up)"
	@command -v golangci-lint >/dev/null 2>&1 && echo "  ✓ golangci-lint" || echo "  ○ golangci-lint (optional — run 'make deps')"
	@command -v air >/dev/null 2>&1 && echo "  ✓ air" || echo "  ○ air (optional — run 'make deps')"

## deploy: deploy to production (run on EC2 after SSH)
deploy:
	@test -f /var/www/html/.env || \
		(echo "ERROR: /var/www/html/.env not found." && \
		 echo "Copy .env.example to /var/www/html/.env and fill in values." && \
		 exit 1)
	git pull
	docker compose build --no-cache
	docker compose up -d
	@echo "Deploy complete. Use 'make logs' to watch."

## logs: tail application logs
logs:
	docker compose logs -f app

## status: show container status
status:
	docker compose ps

## help: list available targets
help:
	@grep -E '^## ' Makefile | sed 's/## //' | column -t -s ':'
