# website-go

A personal blog platform for [jared-wallace.com](https://jared-wallace.com), built as a Go web server with a "weathered beach bar" nautical design. The public blog lives at jared-wallace.com, with an admin panel at admin.jared-wallace.com for writing and managing markdown posts.

Deployed as a Docker container behind an AWS ALB + Nginx reverse proxy.

## Tech Stack

| Layer | Technology | Notes |
|-------|-----------|-------|
| Language | Go 1.26 | Single static binary, stdlib routing (`net/http`) |
| Templates | `html/template` | Server-side rendering with XSS protection |
| Markdown | [goldmark](https://github.com/yuin/goldmark) v1.8.2 | CommonMark-compliant, GFM extensions |
| Database | PostgreSQL 16 via [pgx](https://github.com/jackc/pgx) v5 | Native driver with connection pooling |
| Migrations | [goose](https://github.com/pressly/goose) v3 | Embedded SQL migrations via `go:embed` |
| Sessions | [scs](https://github.com/alexedwards/scs) v2 | Postgres-backed, OWASP-aligned |
| Auth | `golang.org/x/crypto/bcrypt` | Cost 12 |
| RSS | `encoding/xml` stdlib | Zero-dependency feed generation |
| Logging | `log/slog` stdlib | Structured, leveled |
| CI | GitHub Actions | Lint + test + build on free tier |

## Project Structure

```
cmd/
  server/main.go          # Application entrypoint
  hashpw/main.go           # CLI tool: bcrypt password generation
internal/
  config/                  # Environment variable parsing
  database/                # Connection pool + embedded migrations
  handler/
    admin/                 # Admin panel (dashboard, editor, uploads)
    api/                   # JSON API (token-authenticated post push)
    blog/                  # Public blog (list, post, RSS)
  markdown/                # Goldmark rendering pipeline
  middleware/              # Auth, API token, rate limiting
  model/                   # Domain types (Post, Tag, Reaction)
  repository/              # Database queries (pgx)
  server/                  # HTTP server wiring + routing
  service/                 # Business logic layer
web/
  static/                  # CSS, JS, images
  templates/               # HTML templates (base, list, post, admin-*)
db/migrations/             # Goose SQL migration files
  00001_create_posts.sql
  00002_add_tags_to_posts.sql
  00003_create_sessions.sql
  00004_create_reactions.sql
```

## Quick Start

### Prerequisites

- Go 1.26+
- Docker and Docker Compose
- make

### Setup

```bash
# Install optional dev tools (golangci-lint, air)
make deps

# Start local Postgres
make dev-up

# Configure environment
cp .env.example .env.local
# Edit .env.local — at minimum, set:
#   DATABASE_URL=postgres://website:website@localhost:5432/website_dev?sslmode=disable

# Source your env and run migrations + dev server
export $(grep -v '^#' .env.local | xargs)
make migrate
make dev    # hot reload on :8080
```

Visit http://localhost:8080 for the public blog.

To access the admin panel, set `ADMIN_EMAIL` and `ADMIN_PASSWORD_HASH`:

```bash
# Generate a password hash
make hash-password PW=yourpassword

# Export the hash and email, then start the server
export ADMIN_EMAIL=you@example.com
export ADMIN_PASSWORD_HASH='$2a$12$...'  # output from above
export ADMIN_HOST=localhost
```

### Tear Down

```bash
make dev-down
```

## Configuration

All configuration is via environment variables. See `.env.example` for the full template.

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DATABASE_URL` | Yes | -- | Postgres connection string |
| `PORT` | No | `8080` | HTTP listen port |
| `APP_ENV` | No | `development` | `development` or `production` |
| `ADMIN_EMAIL` | No | -- | Admin login email (admin disabled if empty) |
| `ADMIN_PASSWORD_HASH` | No | -- | Bcrypt hash (generate via `make hash-password PW=...`) |
| `ADMIN_HOST` | No | `admin.jared-wallace.com` | Hostname that routes to admin panel |
| `SESSION_SECRET` | No | -- | 32+ char secret for session cookies |
| `API_TOKEN` | No | -- | Bearer token for JSON API (API disabled if empty) |
| `IMAGE_DIR` | No | `/var/www/html/images` | Uploaded image storage path |
| `POSTGRES_USER` | No | -- | Used by docker-compose Postgres sidecar |
| `POSTGRES_PASSWORD` | No | -- | Used by docker-compose Postgres sidecar |
| `POSTGRES_DB` | No | -- | Used by docker-compose Postgres sidecar |

## Makefile Targets

Run `make help` to list all targets. Summary:

| Target | Description |
|--------|-------------|
| `make build` | Compile to `./bin/server` (static binary, CGO off) |
| `make run` | Build and run |
| `make test` | Run tests with race detector |
| `make lint` | Run golangci-lint |
| `make dev` | Hot reload via [air](https://github.com/cosmtrek/air) on :8080 |
| `make dev-up` | Start local Postgres (docker-compose) |
| `make dev-down` | Stop local Postgres |
| `make migrate` | Run Goose migrations up |
| `make migrate-down` | Roll back last migration |
| `make migrate-status` | Show migration status |
| `make docker` | Build Docker image (`website-go:latest`) |
| `make deploy` | Production deploy (run on EC2) |
| `make logs` | Tail container logs |
| `make status` | Show container status |
| `make deps` | Install golangci-lint + air |
| `make check-deps` | Verify required tools are available |
| `make hash-password` | Generate bcrypt hash: `make hash-password PW=secret` |

## Testing

```bash
# Run all tests with race detection
make test

# Lint (errcheck, govet, staticcheck, sqlclosecheck, gosec)
make lint
```

CI runs automatically on push to `main`/`new` and on PRs to `main` via GitHub Actions. The pipeline lints, tests (with a Postgres service container), and builds.

## Docker

### Build

```bash
make docker
# Produces: website-go:latest
```

The Dockerfile uses a two-stage build:
1. **Builder:** `golang:1.26-alpine` -- compiles a static binary with `CGO_ENABLED=0`
2. **Runtime:** `alpine:3.21` -- minimal image with CA certs and timezone data

The container runs as `appuser` (UID 1001) and listens on port 8080.

### Production Compose

`docker-compose.yml` runs the app alongside a Postgres 16 sidecar. Persistent data (database + images) lives on an EBS volume mounted at `/var/www/html`.

## Deployment

The application runs on AWS behind an ALB with HTTPS termination.

### Architecture

```
Route53 (jared-wallace.com)
  -> ALB (HTTPS, ACM cert for *.jared-wallace.com)
    -> ASG (1x t4g.micro, Amazon Linux 2023 ARM64)
      -> Docker (app container + Postgres sidecar)
        -> EBS 10GB gp3 @ /var/www/html (pgdata, images, .env)
```

Infrastructure is managed via Terraform in the `../aws-infra` directory (S3 backend, `us-east-1`).

### First-Time Infrastructure Setup

1. `cd ../aws-infra && terraform init && terraform apply`
2. Create an EC2 key pair in the AWS console
3. SSH to the new instance
4. Clone the repo to `/var/www/html/app`
5. Create `/var/www/html/.env` from `.env.example` with production values
6. `chown 999:999 /var/www/html/pgdata` (Postgres container UID)
7. `make deploy`

### Subsequent Deploys (on EC2)

```bash
ssh ec2-user@<instance-ip>
cd /var/www/html/app
make deploy   # git pull, docker compose build --no-cache, docker compose up -d
make logs     # watch output
make status   # verify containers are healthy
```

## Database

PostgreSQL 16 with pgx v5.9.1 native driver and pgxpool connection pooling.

Migrations are embedded in the binary via `go:embed` and run automatically on startup. They can also be run manually:

```bash
make migrate          # apply pending migrations
make migrate-down     # roll back last migration
make migrate-status   # check current state
```

### Schema

| Table | Purpose |
|-------|---------|
| `posts` | Blog posts with soft-delete support |
| `sessions` | SCS server-side session storage |
| `reactions` | Thumbs-up reactions (IP-hashed for uniqueness) |

Tags are stored as a column on the posts table (migration 00002).

## License

This is a personal project. All rights reserved.
