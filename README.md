# website-go

Personal blog platform for jared-wallace.com — a Go web server with a "weathered beach bar" nautical theme.

## Tech Stack

- **Go** — stdlib `net/http` routing, `html/template` rendering, `log/slog` logging
- **PostgreSQL** — pgx/v5 driver with pgxpool; goose v3 migrations; scs/v2 sessions
- **Markdown** — goldmark (CommonMark) with goldmark-meta, goldmark-highlighting, bluemonday sanitization
- **Auth** — bcrypt password hashing, server-side sessions via scs + pgxstore

## Prerequisites

- Go 1.23+ (local dev) or Go 1.26 (Docker build)
- PostgreSQL 15+
- Docker and docker-compose (for local Postgres sidecar)
- `golangci-lint` (for linting)

## Getting Started

```bash
# 1. Clone the repo
git clone https://github.com/jared-wallace/website-go.git
cd website-go

# 2. Start local Postgres
make dev-up

# 3. Set DATABASE_URL
export DATABASE_URL="postgres://website:website@localhost:5432/website_dev?sslmode=disable"

# 4. Run migrations
make migrate

# 5. Start the server
make run

# 6. Visit http://localhost:8080
```

## Development

| Target      | Description                                      |
|-------------|--------------------------------------------------|
| `make run`  | Build and run the server                         |
| `make dev`  | Hot reload with `air` (install: `go install github.com/cosmtrek/air@latest`) |
| `make test` | Run all tests with race detector                 |
| `make lint` | Run golangci-lint                                |
| `make dev-up`   | Start local Postgres via docker-compose      |
| `make dev-down` | Stop local Postgres                          |
| `make migrate`  | Apply goose migrations (requires `goose` CLI) |
| `make docker`   | Build Docker image                           |
| `make help`     | List all available targets                   |

## Project Structure

```
cmd/server/     — Application entrypoint and server wiring
internal/       — Core packages: config, database, handler, markdown, model, server
db/migrations/  — Goose SQL migrations (embedded into the binary at build time)
web/            — Templates and static assets
```

## License

All rights reserved.
