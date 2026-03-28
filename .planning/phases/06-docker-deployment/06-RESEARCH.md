# Phase 6: Docker + Deployment - Research

**Researched:** 2026-03-27
**Domain:** Docker multi-stage builds, docker-compose production setup, Go binary container patterns
**Confidence:** HIGH

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

- **D-01:** Manual SSH + build on-box. Workflow: SSH to EC2, git pull, `make deploy`. No container registry, no GHA CD pipeline.
- **D-02:** `make deploy` runs the full sequence: git pull, docker compose build, run migrations via entrypoint, docker compose up -d. One command, no steps to forget.
- **D-03:** Migrations run automatically in the container entrypoint before the server starts. Safe for single-instance (no concurrent deploys).
- **D-04:** Use the compiled binary's embedded migrations (go:embed), not the goose CLI. Add a `-migrate` flag or subcommand to the server binary. Entrypoint runs `./server -migrate` then `./server`. Keeps the container image minimal — no goose binary needed.
- **D-05:** Production env vars managed via a `.env` file on the EBS volume, referenced by docker-compose via `env_file`. Survives container rebuilds, never committed to git.
- **D-06:** Ship a `.env.example` with placeholder values. `make deploy` checks `.env` exists and fails with a helpful message if missing. No silent defaults for secrets.
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

### Deferred Ideas (OUT OF SCOPE)

None — discussion stayed within phase scope.

</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| FOUND-04 | Docker multi-stage build producing minimal container | Multi-stage Dockerfile pattern: `golang:1.26-alpine` build stage, `alpine:3.21` runtime stage; CGO_ENABLED=0 for static binary; templates/static embedded so final image contains only binary |
| FOUND-05 | docker-compose with app + Postgres sidecar, EBS volume mounts | `docker-compose.yml` with bind-mount `/var/www/html/pgdata:/var/lib/postgresql/data` and `/var/www/html/images:/var/www/html/images`; health check + depends_on; env_file |

</phase_requirements>

---

## Summary

Phase 6 ships a Go blog binary as a minimal Docker container with a Postgres sidecar. The existing codebase already embeds all templates and static assets via `go:embed` (`web/embed.go`), so the final Docker image only needs the compiled binary and CA certificates — no asset copying step beyond the binary. Migrations already run in `main.go` via `database.RunMigrations()`, which uses the embedded `db/migrations/embed.go` FS; D-04's `-migrate` flag is effectively pre-implemented at the application layer.

The two main deliverables are a `Dockerfile` (multi-stage, non-root runtime user) and `docker-compose.yml` (prod compose with EBS bind-mounts). Supporting pieces are: `.env.example` documenting all env vars from `config.go`, a `make deploy` target with `.env` guard, a `/health` HTTP endpoint, and optional `make logs` / `make status` targets. The `docker-compose.dev.yml` already exists and must not be changed.

The critical operational pitfall is the `chown 999:999` host prerequisite for the Postgres EBS bind-mount directory — the Postgres container runs as UID 999 and will refuse to start if the mounted directory is owned by root.

**Primary recommendation:** One `Dockerfile`, one `docker-compose.yml`, one `.env.example`, two Makefile targets (`deploy`, `logs`). Everything else is a single-file edit.

---

## Standard Stack

### Core

| Component | Version | Purpose | Why Standard |
|-----------|---------|---------|--------------|
| `golang:1.26-alpine` | 1.26 | Build stage base image | Matches project Go version; alpine reduces layer size |
| `alpine:3.21` | 3.21 | Runtime stage base image | Smallest image with CA certs + tzdata (needed for Postgres TLS); scratch omits CA certs and breaks pgx TLS handshake |
| `postgres:16-alpine` | 16-alpine | Postgres sidecar | Matches dev compose and CI; 16 is current LTS |
| Docker multi-stage | 28.x (local) | Container build | Build stage compiles; runtime stage copies binary only; final image ~20MB |
| docker compose v2 | v2.35.1 (local) | Orchestration | Already in use for dev; prod compose file is additive |

### Supporting

| Component | Version | Purpose | When to Use |
|-----------|---------|---------|-------------|
| `make deploy` | GNU Make 3.81 | Deploy automation | Wraps git pull + compose build + up; documents the runbook |
| `.env` / `env_file` | — | Secret injection | EBS-resident `.env` feeds all config vars into the app container |
| `pg_isready` health check | (built into postgres image) | Postgres liveness | Already used in dev compose; replicate in prod |

### Alternatives Considered

| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| `alpine:3.21` runtime | `scratch` | scratch is smaller but no CA certs; pgx TLS connections to Postgres fail without certs |
| shell entrypoint script | `-migrate` flag in binary | D-04 says binary flag; but migrations already run unconditionally in main.go — entrypoint may just exec the binary directly |
| bind-mount for images | Named volume | Named volumes can't easily be inspected/backed-up from the host; EBS bind-mount is explicit and maps to backup policy |

**Installation:** No new packages required. `golang:1.26-alpine` and `alpine:3.21` pull from Docker Hub on first `docker build`.

---

## Architecture Patterns

### Recommended Project Structure (new files only)

```
.
├── Dockerfile                    # Multi-stage build (new)
├── docker-compose.yml            # Prod compose (new)
├── docker-compose.dev.yml        # Dev Postgres only (EXISTING — do not touch)
├── .env.example                  # Template for EBS .env file (new)
├── .dockerignore                 # Exclude .env, bin/, test files from build context (new)
├── Makefile                      # Add deploy, logs, status targets (extend existing)
└── internal/
    └── handler/
        └── health/               # OR: add /health route directly in main.go
            └── health.go         # GET /health → 200 OK
```

### Pattern 1: Multi-Stage Dockerfile

**What:** Two FROM stages — builder compiles the Go binary with all dependencies, runtime copies only the binary.

**When to use:** Always for Go services. Keeps final image small; the Go toolchain (~500MB) never ships to production.

```dockerfile
# Source: https://docs.docker.com/language/golang/build-images/ (Go multi-stage pattern)
# Build stage
FROM golang:1.26-alpine AS builder

# Install git (needed for go modules that use git)
RUN apk add --no-cache git

WORKDIR /build

# Copy module files first for layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary — CGO_ENABLED=0 required for alpine:3.21 runtime
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/server ./cmd/server

# Runtime stage
FROM alpine:3.21

# CA certificates: required for pgx TLS connections to Postgres
# tzdata: required for time zone handling in SQL queries
RUN apk add --no-cache ca-certificates tzdata

# Non-root user — UID 1001 (avoids conflict with postgres UID 999)
RUN adduser -D -u 1001 appuser
USER appuser

WORKDIR /app
COPY --from=builder /build/bin/server .

EXPOSE 8080
CMD ["./server"]
```

**Note on UID choice:** Postgres container uses UID 999. App user should be a different UID (1001 is conventional for non-root app users). This avoids any accidental permission overlap on shared volume paths.

**Note on `-ldflags="-s -w"`:** Strips debug symbols and DWARF info, reducing binary size by ~30% with no runtime impact.

**Note on templates/static:** `web/embed.go` uses `go:embed templates` and `go:embed static`. These are baked into the binary at build time — no `COPY` of the `web/` directory is needed in the runtime stage. The binary is self-contained.

### Pattern 2: Production docker-compose.yml

**What:** App + Postgres with EBS bind-mounts, env_file, health checks, restart policies.

**When to use:** Production only. Dev uses `docker-compose.dev.yml` (Postgres sidecar only, named volumes).

```yaml
# Source: docker-compose reference + project decisions D-05, D-07, D-08
services:
  postgres:
    image: postgres:16-alpine
    restart: unless-stopped
    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: ${POSTGRES_DB}
    volumes:
      # EBS bind-mount — host dir must be: sudo chown 999:999 /var/www/html/pgdata
      - /var/www/html/pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER}"]
      interval: 5s
      timeout: 5s
      retries: 5

  app:
    build: .
    restart: unless-stopped
    ports:
      - "8080:8080"
    env_file:
      - /var/www/html/.env
    volumes:
      # Image uploads — EBS-backed, survives container restarts
      - /var/www/html/images:/var/www/html/images
    depends_on:
      postgres:
        condition: service_healthy
    healthcheck:
      test: ["CMD-SHELL", "wget -qO- http://localhost:8080/health || exit 1"]
      interval: 10s
      timeout: 5s
      retries: 3
      start_period: 30s
```

**Key decisions:**
- `depends_on: condition: service_healthy` — app container waits for Postgres to pass its health check before starting. This is the compose-native solution to startup ordering.
- `env_file: /var/www/html/.env` — absolute path to EBS-resident secrets file.
- `start_period: 30s` — gives the app time to run migrations before health checks begin.

### Pattern 3: make deploy Target

**What:** Single command that performs the full deploy sequence with guards.

```makefile
## deploy: deploy to production (run on EC2 via SSH)
deploy:
	@test -f /var/www/html/.env || (echo "ERROR: /var/www/html/.env not found. Copy .env.example and fill in values."; exit 1)
	git pull
	docker compose build --no-cache
	docker compose up -d
	@echo "Deploy complete. Run 'make logs' to tail output."

## logs: tail application logs
logs:
	docker compose logs -f app

## status: show running container status
status:
	docker compose ps
```

**Note on `--no-cache`:** On a single-instance blog, build cache is less valuable than correctness. `--no-cache` guarantees a clean build on each deploy. If build time becomes an issue (it won't for this binary size), drop the flag.

### Pattern 4: /health Endpoint

**What:** Simple HTTP 200 OK response. No DB ping required.

**Rationale:** The Postgres health check is handled at the compose layer (`pg_isready`). The app's `/health` endpoint just proves the binary is alive and listening. A DB ping adds complexity and can cause health checks to fail during brief Postgres restarts — restarting both containers unnecessarily.

```go
// Add to blogMux in main.go (simplest integration)
blogMux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
})
```

### Pattern 5: .env.example

**What:** Documents all required and optional env vars with placeholder values and comments.

```bash
# Required
DATABASE_URL=postgres://website:CHANGEME@postgres:5432/website_prod?sslmode=disable
ADMIN_EMAIL=admin@example.com
ADMIN_PASSWORD_HASH=$2a$12$CHANGEME_RUN_make_hash-password_PW=yourpassword

# Optional (defaults shown)
PORT=8080
APP_ENV=production
ADMIN_HOST=admin.jared-wallace.com
SESSION_SECRET=
API_TOKEN=
IMAGE_DIR=/var/www/html/images

# Postgres sidecar (referenced in docker-compose.yml)
POSTGRES_USER=website
POSTGRES_PASSWORD=CHANGEME
POSTGRES_DB=website_prod
```

**Note on DATABASE_URL hostname:** Inside docker-compose, the Postgres hostname is `postgres` (the service name), not `localhost`. This is a common first-time mistake.

### Pattern 6: .dockerignore

**What:** Excludes files from the build context to speed up `docker build`.

```
.env
.env.example
bin/
.git/
.planning/
*.md
docker-compose*.yml
Makefile
.air.toml
```

### Anti-Patterns to Avoid

- **Running migrations as a separate entrypoint step when they already run in main.go:** D-04 says "add a `-migrate` flag"; however, `main.go` already calls `database.RunMigrations()` unconditionally before starting the server. The flag approach would duplicate this. **Resolution:** The entrypoint is simply `CMD ["./server"]` — migrations run automatically. D-04 is already implemented by the existing code.
- **Using `latest` tags for base images:** Pin `postgres:16-alpine`, `golang:1.26-alpine`, `alpine:3.21` to avoid silent breakage on compose pull.
- **Storing `.env` in git:** `.env` must be in `.gitignore`; `.env.example` (with dummy values) is committed.
- **Bind-mounting the entire EBS volume into the app container:** Mount only `/var/www/html/images`, not `/var/www/html`. Keeps Postgres data and secrets outside the app container's writable surface.
- **Using `depends_on` without `condition: service_healthy`:** `depends_on` alone only waits for the container to start, not for Postgres to be ready. Without `condition: service_healthy`, the app races the database and migration fails.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Startup ordering | Sleep loops in entrypoint | `depends_on: condition: service_healthy` | Compose handles this natively with pg_isready; sleep loops are fragile and hard to tune |
| Secrets management | Env vars baked into image | `env_file` on EBS volume | Secrets in image layers are visible to anyone with image access |
| Migration on startup | Custom shell migration script | Existing `database.RunMigrations()` in main.go | Already implemented; calling it again from shell would double-apply or create race |
| Health check binary | wget install in runtime stage | `wget` already present in `alpine:3.21` base | Alpine includes wget; no extra RUN step needed |

**Key insight:** The application already does the hard part. Migrations run on startup, assets are embedded, config loads from env. Docker is just packaging what already works.

---

## Runtime State Inventory

Phase 6 is greenfield infrastructure (new files). No rename or data migration is involved.

| Category | Items Found | Action Required |
|----------|-------------|-----------------|
| Stored data | None — Phase 6 creates the container layer, does not rename stored data | None |
| Live service config | None — no live service exists yet; this phase creates it | None |
| OS-registered state | None | None |
| Secrets/env vars | `/var/www/html/.env` must be created on EC2 before first deploy | Human task: create from `.env.example` |
| Build artifacts | `bin/server` — excluded from Docker build context via `.dockerignore` | None; Docker build compiles fresh |

---

## Common Pitfalls

### Pitfall 1: Postgres bind-mount permission error

**What goes wrong:** Postgres container starts then immediately exits with "initdb: error: directory '/var/lib/postgresql/data' exists but is not empty or could not be created". More commonly: `chown: /var/lib/postgresql/data: Permission denied`.

**Why it happens:** The Postgres container process runs as UID 999 (the `postgres` system user inside the image). If the host directory `/var/www/html/pgdata` is owned by root (uid 0), Postgres cannot write to it.

**How to avoid:** Before the first deploy, on the EC2 host:
```bash
sudo mkdir -p /var/www/html/pgdata
sudo chown 999:999 /var/www/html/pgdata
```
The `make deploy` target should check for this directory and warn if missing or wrong ownership.

**Warning signs:** Postgres container exits immediately; `docker compose logs postgres` shows permission denied or initdb errors.

### Pitfall 2: DATABASE_URL uses localhost instead of service name

**What goes wrong:** App container starts, tries to connect to `localhost:5432`, gets "connection refused" even though Postgres container is running.

**Why it happens:** Inside docker-compose, each service has its own network namespace. `localhost` inside the app container is the app container. Postgres is reachable as `postgres` (the compose service name).

**How to avoid:** `.env.example` must use `postgres` as the hostname: `DATABASE_URL=postgres://...@postgres:5432/...`. Document this prominently in `.env.example` comments.

**Warning signs:** App starts, logs "database ping failed: dial tcp 127.0.0.1:5432: connect: connection refused".

### Pitfall 3: Migrations running before Postgres accepts connections

**What goes wrong:** App container starts before Postgres is ready to accept connections. `database.RunMigrations()` panics or returns an error.

**Why it happens:** Docker starts containers in dependency order, but "started" != "ready". Postgres needs a few seconds for initdb or recovery before accepting connections.

**How to avoid:** `depends_on: condition: service_healthy` makes Docker wait for pg_isready to succeed before starting the app container. This is already in the recommended compose pattern above.

**Warning signs:** App crashes with "database ping failed" on first deploy; works fine on restart.

### Pitfall 4: CA certificates missing in scratch/minimal runtime

**What goes wrong:** pgx TLS connection to Postgres fails with "x509: certificate signed by unknown authority" or "tls: failed to verify certificate".

**Why it happens:** A `scratch` runtime image has no CA certificate bundle. `alpine:3.21` includes `/etc/ssl/certs/ca-certificates.crt`.

**How to avoid:** Use `alpine:3.21` as the runtime base, not `scratch`. Already in the recommended Dockerfile.

**Warning signs:** TLS errors on `DATABASE_URL` connections; does not occur with `sslmode=disable` (development setting).

### Pitfall 5: Embedded assets not included in binary

**What goes wrong:** Server starts but returns 404 for all static assets and template rendering fails.

**Why it happens:** `go:embed` works at compile time from the directory structure relative to the Go source file. If `COPY . .` in the Dockerfile doesn't include `web/templates/` and `web/static/`, the embed is empty.

**How to avoid:** The recommended `COPY . .` in the build stage copies the entire repo including `web/`. No additional COPY is needed. Verify with `docker run --rm website-go:latest ls -la /app` — only `server` binary should be there; embedded assets are inside the binary.

**Warning signs:** 404 for `/static/` routes; template execute errors in logs.

### Pitfall 6: .env checked into git

**What goes wrong:** Production database URL, admin password hash, and API token end up in version control history permanently.

**Why it happens:** Accidentally running `git add .` without a `.gitignore` entry.

**How to avoid:** Confirm `.env` is in `.gitignore`. Only `.env.example` (with placeholder values) is committed. The `make deploy` target reads from `/var/www/html/.env` (absolute EBS path), not from the repo.

---

## Code Examples

Verified patterns from the existing codebase and Docker official docs.

### Dockerfile: Full Multi-Stage Build

```dockerfile
# Source: Docker multi-stage Go docs + STACK.md (golang:1.26-alpine, alpine:3.21)

# ── Build stage ──────────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git

WORKDIR /build

# Layer cache: download modules before copying source
COPY go.mod go.sum ./
RUN go mod download

# Copy full source (includes web/templates, web/static for go:embed)
COPY . .

# Static binary — no CGO, strip symbols for size
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o bin/server \
    ./cmd/server

# ── Runtime stage ─────────────────────────────────────────────────────────────
FROM alpine:3.21

# ca-certificates: pgx TLS connections to Postgres
# tzdata: time zone data for SQL timestamp handling
RUN apk add --no-cache ca-certificates tzdata

# Non-root user (UID 1001 — distinct from postgres UID 999)
RUN adduser -D -u 1001 appuser
USER appuser

WORKDIR /app
COPY --from=builder /build/bin/server .

EXPOSE 8080
CMD ["./server"]
```

### /health Handler (add to blogMux in main.go)

```go
// Source: existing main.go pattern (blogMux.HandleFunc)
// Add after other blogMux registrations, before catch-all 404
blogMux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
})
```

### Makefile deploy Target

```makefile
## deploy: deploy to production (run on EC2 after SSH)
deploy:
	@test -f /var/www/html/.env || \
		(echo "ERROR: /var/www/html/.env not found."; \
		 echo "Copy .env.example to /var/www/html/.env and fill in values."; \
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
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| `docker-compose` (v1, standalone binary) | `docker compose` (v2, plugin) | Docker Desktop 3.4+ / 2021 | Command changed from hyphen to space; v2 is default everywhere now |
| `depends_on: [postgres]` (starts only) | `depends_on: condition: service_healthy` | docker-compose v3.9 / 2021 | Proper startup ordering without sleep loops |
| ENV in Dockerfile for secrets | `env_file` in compose / runtime injection | Security baseline | Secrets in image layers are visible in `docker inspect` and image history |

**Deprecated/outdated:**
- `docker-compose.yml` version: field: Compose v2 ignores the `version:` key; omit it to avoid confusion.
- `links:` directive: Replaced by default bridge networking in compose; services reach each other by service name automatically.

---

## Open Questions

1. **D-04 vs existing RunMigrations behavior**
   - What we know: `database.RunMigrations()` already runs unconditionally in `main.go` before server start (lines 79-83). D-04 specifies a `-migrate` flag in the binary.
   - What's unclear: D-04 may have been written before realizing the migration already runs on startup; or it may intend a separate "migrate only, no server" mode for manual recovery scenarios.
   - Recommendation: Planner should clarify. Simplest path: `CMD ["./server"]` is sufficient since migrations run on startup. If a `-migrate` flag is wanted for ops tooling, it's additive (a few lines in main.go).

2. **EBS .env file absolute path in env_file**
   - What we know: docker-compose `env_file` supports absolute paths.
   - What's unclear: Whether `/var/www/html/.env` is acceptable for local dev if someone accidentally runs `docker compose up` (without `-f docker-compose.dev.yml`).
   - Recommendation: `.env.example` note and Makefile guard make this explicit; low risk with the dev/prod file separation (D-07).

---

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Docker | `docker build`, `make deploy` | ✓ | 28.1.1 | — |
| docker compose | `make dev-up`, `make deploy` | ✓ | v2.35.1 | — |
| GNU Make | All Makefile targets | ✓ | 3.81 | Run commands manually |
| `golang:1.26-alpine` (Docker Hub) | Build stage | Not cached locally | Pulled on `docker build` | — |
| `alpine:3.21` (Docker Hub) | Runtime stage | Not cached locally | Pulled on `docker build` | — |
| `postgres:16-alpine` (Docker Hub) | Sidecar | Cached locally (16-alpine) | 16-alpine available | — |

**Missing dependencies with no fallback:** None. All required tools are present.

**Missing dependencies with fallback:** `golang:1.26-alpine` and `alpine:3.21` not cached locally but will pull automatically on first `docker build`. EC2 must have outbound internet access to Docker Hub (standard).

---

## Validation Architecture

### Test Framework

| Property | Value |
|----------|-------|
| Framework | Go testing stdlib (built-in) |
| Config file | None — `go test ./...` discovers all `*_test.go` files |
| Quick run command | `go test ./... -count=1` |
| Full suite command | `go test ./... -v -race` |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| FOUND-04 | Docker image builds without error | smoke | `docker build -t website-go:test . && docker image inspect website-go:test` | ❌ Wave 0 (manual / CI step) |
| FOUND-04 | Final image contains only binary (no Go toolchain) | smoke | `docker run --rm website-go:test ls /app` | ❌ Wave 0 (manual) |
| FOUND-04 | Binary runs as non-root user | smoke | `docker run --rm website-go:test id` | ❌ Wave 0 (manual) |
| FOUND-05 | `docker compose up` starts both containers | smoke | `docker compose up -d && docker compose ps` | ❌ Wave 0 (manual) |
| FOUND-05 | App serves /health within 30 seconds | smoke | `docker compose up -d && sleep 10 && curl -f http://localhost:8080/health` | ❌ Wave 0 (manual) |
| FOUND-05 | EBS data dir survives container restart | integration | `docker compose down && docker compose up -d` — verify data persists | ❌ Wave 0 (manual) |
| — | /health endpoint returns 200 OK | unit | `go test ./... -run TestHealth` (or handler test) | ❌ Wave 0 |

**Note:** FOUND-04 and FOUND-05 are primarily infrastructure/smoke tests that cannot be fully automated with `go test`. Automated verification is `docker build` + `docker run` commands. The `/health` handler is the one testable unit.

### Sampling Rate

- **Per task commit:** `go test ./... -count=1` (existing suite, ensures /health handler doesn't break routing)
- **Per wave merge:** `go test ./... -v -race` (full suite with race detector)
- **Phase gate:** Full suite green + `docker build` succeeds + `docker compose up` starts healthy

### Wave 0 Gaps

- [ ] `/health` handler test — add to `internal/handler/blog/handler_test.go` or new file covering `GET /health` → 200
- [ ] No new framework install needed — Go stdlib testing already in place

---

## Project Constraints (from CLAUDE.md)

| Directive | Implication for This Phase |
|-----------|---------------------------|
| Go stdlib-first; minimal deps | No new Go dependencies needed for this phase |
| Must run as Docker container on port 8080 | `EXPOSE 8080`, `PORT=8080` default in `.env.example` |
| All persistent data at /var/www/html | Bind-mounts for pgdata and images both under /var/www/html |
| GHA CI on free tier | No new CI jobs; `make deploy` is SSH-based, not GHA |
| Docker: `golang:1.26-alpine` build, `alpine:3.21` runtime | Locked in STACK.md — use exactly these tags |
| `CGO_ENABLED=0` for static binary | Required in Dockerfile build stage |

---

## Sources

### Primary (HIGH confidence)

- Docker multi-stage build docs — https://docs.docker.com/language/golang/build-images/
- docker-compose `depends_on` condition reference — https://docs.docker.com/compose/compose-file/05-services/#depends_on
- `.planning/research/STACK.md` — base image versions (golang:1.26-alpine, alpine:3.21, postgres:16-alpine)
- `web/embed.go`, `db/migrations/embed.go` — existing embed patterns (read directly)
- `internal/config/config.go` — complete env var inventory (read directly)
- `cmd/server/main.go` lines 79-83 — `database.RunMigrations()` already runs on startup
- `internal/database/migrations.go` — goose embedded migration implementation

### Secondary (MEDIUM confidence)

- Alpine package availability (ca-certificates, tzdata, wget) — standard Alpine package index; verified by project STACK.md noting alpine:3.21 for CA certs
- Postgres UID 999 in official postgres:16-alpine image — widely documented, consistent across postgres Docker image releases

### Tertiary (LOW confidence)

- None — all claims in this document are grounded in codebase inspection or official documentation.

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — base images from STACK.md, versions from existing go.mod and docker-compose.dev.yml
- Architecture: HIGH — Dockerfile pattern from official Go Docker docs; compose patterns from official compose spec; env var list read directly from config.go
- Pitfalls: HIGH — Pitfalls 1-3 documented in STATE.md as known blockers; Pitfalls 4-6 from standard Docker/Go operational knowledge

**Research date:** 2026-03-27
**Valid until:** 2026-06-27 (stable tech; Docker and compose API changes slowly)
