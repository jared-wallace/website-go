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
| `make migrate`  | Apply goose migrations (requires `goose` CLI and `DATABASE_URL`) |
| `make docker`   | Build Docker image                           |
| `make help`     | List all available targets                   |

## Project Structure

```
cmd/server/        — Application entrypoint and server wiring
internal/          — Core packages
  config/          — Environment and app configuration
  database/        — PostgreSQL connection and pool setup
  handler/         — HTTP handlers
  markdown/        — Goldmark rendering pipeline with bluemonday sanitization
  model/           — Domain types (Post, etc.)
  server/          — Server setup and middleware
db/migrations/     — Goose SQL migrations (embedded at build time via embed.go)
web/               — Frontend assets
  templates/       — Go html/template files
  static/          — CSS and static files
.github/workflows/ — CI pipeline (lint, test, build)
```

## Continuous Integration

`.github/workflows/ci.yml` runs on push to `main` and on pull requests. It targets Go 1.26 with a Postgres 16 sidecar on free-tier GitHub-hosted runners.

## Production Deployment

> Note: Dockerfile, deploy.sh, and docker-compose.prod.yml are Phase 6 deliverables and do not exist yet. The steps below document the intended deployment flow once Phase 6 is complete.

### Architecture Overview

```
Client
  -> Route53 (jared-wallace.com / www / admin)
  -> ALB (TLS termination via ACM wildcard cert, us-east-1)
  -> EC2 t4g.micro (ARM64, Amazon Linux 2023)
  -> Nginx :80 (systemd-managed reverse proxy)
  -> Go app container :8080
```

- **Infrastructure:** Managed in a separate `aws-infra` repo via Terraform.
- **ASG:** min=1, max=1, desired=1 — single instance for self-healing, not horizontal scaling.
- **Storage:** 10 GB gp3 EBS volume mounted at `/var/www/html` — persists across instance replacement.
- **Domains:** A/AAAA records for `jared-wallace.com`, `www`, and `admin` alias the ALB. `ssh.jared-wallace.com` is updated via Lambda + EventBridge when the ASG cycles.

### Deployment Steps

Once Phase 6 is complete:

```bash
# 1. SSH to the instance
ssh ec2-user@ssh.jared-wallace.com

# 2. Pull latest code
cd /var/www/html/app
git pull origin main

# 3. Build the Docker image (ARM64 target)
docker build -t website-go:latest .

# 4. Start the app via docker-compose (production config)
docker compose -f docker-compose.prod.yml up -d

# Nginx is already running via systemd and proxies :80 -> localhost:8080
```

### Persistent Storage Layout

All persistent data lives on the EBS volume and survives instance replacement:

| Path | Contents |
|------|----------|
| `/var/www/html/` | EBS mount point |
| `/var/www/html/app/` | Application code |
| `/var/www/html/postgres-data/` | Postgres data directory (planned) |
| `/var/www/html/images/` | Uploaded blog images (planned) |

### Caveats

- **ASG max_size:** Must be set to `1` with `delete_on_termination = false` on the EBS volume before any production data is written. Scaling to 2 instances while data is on EBS will cause split-brain.
- **Postgres directory ownership:** Before the first `docker compose up`, run `chown 999:999 /var/www/html/postgres-data` on the host. The Postgres container runs as UID 999 and will fail to start otherwise.
- **Phase 6 deliverables not yet created:** `Dockerfile`, `deploy.sh`, `docker-compose.prod.yml`.

## License

All rights reserved.
