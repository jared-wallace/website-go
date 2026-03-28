# Phase 6: Docker + Deployment - Context

**Gathered:** 2026-03-27
**Status:** Ready for planning

<domain>
## Phase Boundary

Ship the application as a minimal Docker container in a docker-compose stack with a Postgres sidecar, correctly bind-mounted to the EBS volume at /var/www/html, deployable behind the existing Nginx + ALB via a single `make deploy` command. No CI/CD pipeline changes, no registry infrastructure, no orchestration beyond docker-compose.

</domain>

<decisions>
## Implementation Decisions

### Deployment Workflow
- **D-01:** Manual SSH + build on-box. Workflow: SSH to EC2, git pull, `make deploy`. No container registry, no GHA CD pipeline. Fits single-instance blog with zero extra infra.
- **D-02:** `make deploy` runs the full sequence: git pull, docker compose build, run migrations via entrypoint, docker compose up -d. One command, no steps to forget.

### Migration Strategy
- **D-03:** Migrations run automatically in the container entrypoint before the server starts. Safe for single-instance (no concurrent deploys).
- **D-04:** Use the compiled binary's embedded migrations (go:embed), not the goose CLI. Add a `-migrate` flag or subcommand to the server binary. Entrypoint runs `./server -migrate` then `./server`. Keeps the container image minimal — no goose binary needed.

### Environment & Secrets
- **D-05:** Production env vars managed via a `.env` file on the EBS volume, referenced by docker-compose via `env_file`. Survives container rebuilds, never committed to git.
- **D-06:** Ship a `.env.example` with placeholder values. `make deploy` checks `.env` exists and fails with a helpful message if missing. No silent defaults for secrets.

### Compose File Strategy
- **D-07:** Separate files: keep `docker-compose.dev.yml` (Postgres-only, named volumes) for local dev, add `docker-compose.yml` for prod (app + Postgres, EBS bind-mounts, .env file). Clean separation, no override confusion.
- **D-08:** Prod compose includes `restart: unless-stopped` on both containers. Health check on app via HTTP `/health` endpoint, Postgres via `pg_isready`. Docker auto-restarts crashed containers.

### Claude's Discretion
- Dockerfile multi-stage build structure (build stage base image, runtime stage base image, COPY strategy)
- Entrypoint script implementation (shell script vs binary flags)
- `/health` endpoint implementation (simple 200 OK vs DB ping)
- Non-root user setup in the container (UID/GID choice)
- Makefile deploy target implementation details (error handling, output formatting)
- `.env.example` placeholder values and documentation comments
- Whether to add a `make logs` or `make status` convenience target
- Docker build cache optimization strategy
- Container logging configuration (stdout/stderr vs file)

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Tech Stack
- `.planning/research/STACK.md` -- Authoritative dependency versions and Docker base image guidance (golang:1.26-alpine build, alpine:3.21 runtime, CGO_ENABLED=0)

### Project Context
- `.planning/PROJECT.md` -- Core value, constraints, EBS volume at /var/www/html, port 8080 behind Nginx/ALB, minimal deps philosophy
- `.planning/REQUIREMENTS.md` -- FOUND-04 (Docker multi-stage build), FOUND-05 (docker-compose with app + Postgres sidecar)

### Infrastructure Context
- `../aws-infra` -- Terraform repo with ALB, ASG, Nginx, EBS volume configuration (referenced in PROJECT.md Context section)

### Existing Codebase
- `Makefile` -- Current targets to extend with deploy, logs, status
- `docker-compose.dev.yml` -- Existing dev Postgres setup (named volumes, pg_isready health check pattern)
- `.github/workflows/ci.yml` -- CI pipeline (no changes needed, but reference for Go version and build command)
- `internal/config/config.go` -- Env var pattern (mustEnv/envOr) that .env file must satisfy
- `db/migrations/embed.go` -- Embedded migrations pattern to reuse for entrypoint migration

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `db/migrations/embed.go`: go:embed migrations — reuse for entrypoint `-migrate` flag implementation
- `docker-compose.dev.yml`: Postgres health check pattern (pg_isready) — replicate in prod compose
- `Makefile`: Existing `docker` target builds image — extend with deploy, logs, status targets
- `internal/config/config.go`: All env vars documented with defaults — basis for `.env.example`

### Established Patterns
- Config via env vars with `mustEnv()`/`envOr()` — `.env` file maps directly to this
- Single binary at `./cmd/server` with `go build -o bin/server ./cmd/server`
- Templates and static files embedded or served from `web/` directory
- Postgres 16-alpine used in dev compose and CI — use same in prod

### Integration Points
- `cmd/server/main.go`: Add `-migrate` flag handling before server start
- `Makefile`: Add `deploy`, potentially `logs` and `status` targets
- Root directory: New `Dockerfile`, `docker-compose.yml`, `.env.example`
- `internal/server/server.go` or new handler: `/health` endpoint for Docker health check

</code_context>

<specifics>
## Specific Ideas

No specific requirements -- open to standard approaches.

</specifics>

<deferred>
## Deferred Ideas

None -- discussion stayed within phase scope.

</deferred>

---

*Phase: 06-docker-deployment*
*Context gathered: 2026-03-27*
