---
phase: 01-foundation
plan: 01
subsystem: database
tags: [go, pgx, goose, goldmark, bluemonday, scs, bcrypt, embed]

# Dependency graph
requires: []
provides:
  - Go module initialized with all Phase 1 dependencies
  - Standard project layout (cmd/, internal/, db/, web/)
  - Config package with 12-factor envOr/mustEnv helpers
  - Post domain model with soft-delete
  - Database Connect() with pgxpool health check
  - RunMigrations() via goose + pgx stdlib adapter
  - posts table schema with slug, soft-delete, and indexes
  - Graceful shutdown via signal.NotifyContext
  - go:embed wiring for templates and static assets
affects:
  - 01-02 (markdown pipeline — uses goldmark deps installed here)
  - 01-03 (makefile/CI — uses go.mod module path and build targets)
  - all future phases — directory layout and config package are foundations

# Tech tracking
tech-stack:
  added:
    - github.com/yuin/goldmark v1.8.2
    - github.com/yuin/goldmark-meta v1.1.0
    - github.com/yuin/goldmark-highlighting/v2
    - github.com/alecthomas/chroma/v2 v2.23.1
    - github.com/microcosm-cc/bluemonday v1.0.27
    - github.com/jackc/pgx/v5 v5.9.1
    - github.com/pressly/goose/v3 v3.27.0
    - github.com/alexedwards/scs/v2 v2.9.0
    - github.com/alexedwards/scs/pgxstore
    - golang.org/x/crypto v0.49.0
  patterns:
    - "12-factor config: envOr() for defaults, mustEnv() panics on missing required vars"
    - "DB pool injected via constructors — no global var DB *pgxpool.Pool"
    - "Embedded migrations in sibling package (db/migrations/embed.go) to avoid go:embed .. restriction"
    - "Graceful shutdown via signal.NotifyContext for SIGINT/SIGTERM"
    - "http.Server with explicit ReadTimeout/WriteTimeout/IdleTimeout — no bare ListenAndServe"

key-files:
  created:
    - go.mod
    - go.sum
    - cmd/server/main.go
    - internal/config/config.go
    - internal/config/config_test.go
    - internal/model/post.go
    - internal/server/server.go
    - web/embed.go
    - internal/database/database.go
    - internal/database/database_test.go
    - internal/database/migrations.go
    - db/migrations/00001_create_posts.sql
    - db/migrations/embed.go
  modified: []

key-decisions:
  - "go.mod uses 'go 1.26' directive without 'toolchain' to avoid strict enforcement on local Go 1.23"
  - "Migrations embedded in db/migrations/embed.go (sibling package) not internal/database — go:embed cannot traverse .. path components"
  - "goose.Up path arg is '.' when FS root is the migrations directory (FS contains *.sql directly)"
  - "web/templates and web/static use placeholder files (not .gitkeep) — go:embed requires at least one embeddable file per directory"

patterns-established:
  - "Pattern 1: Config loaded via config.Load() at startup, passed through constructors"
  - "Pattern 2: Database pool created once in main, passed to all consumers"
  - "Pattern 3: Embedded SQL migrations in package adjacent to SQL files"
  - "Pattern 4: Graceful shutdown wraps entire server lifecycle"

requirements-completed: [FOUND-01, FOUND-02, FOUND-03]

# Metrics
duration: 4min
completed: 2026-03-26
---

# Phase 1 Plan 01: Foundation Summary

**Go module scaffolded with pgx v5 pool, embedded goose migrations, posts schema, and 12-factor config — compilable binary connects to Postgres and runs migrations at startup**

## Performance

- **Duration:** ~4 min
- **Started:** 2026-03-26T11:40:00Z
- **Completed:** 2026-03-26T11:44:11Z
- **Tasks:** 2
- **Files modified:** 13 created

## Accomplishments

- Go module initialized with all Phase 1 dependencies in go.mod (goldmark, pgx, goose, scs, bcrypt, bluemonday)
- Standard project layout established: cmd/server/, internal/{config,database,model,server,handler,markdown}/, db/migrations/, web/{templates,static}/
- Config package with 12-factor envOr/mustEnv helpers, unit-tested (3 tests pass without Postgres)
- Database package: pgxpool.Connect with Ping health check + RunMigrations via goose embedded SQL
- Initial schema: posts table with BIGSERIAL id, slug UNIQUE, rendered_html, soft-delete, and partial indexes
- Binary compiles cleanly; `go build ./...` and `go vet ./...` both exit 0

## Task Commits

Each task was committed atomically:

1. **Task 1: Initialize Go module, directory skeleton, and all dependencies** - `e8ec780` (feat)
2. **Task 2: Database connection, embedded migrations, and initial schema** - `7ebd1f8` (feat)

## Files Created/Modified

- `go.mod` / `go.sum` — module definition with all Phase 1 dependencies
- `cmd/server/main.go` — entry point wiring config → DB → migrations → graceful shutdown
- `internal/config/config.go` — Config struct, Load(), envOr(), mustEnv()
- `internal/config/config_test.go` — 3 unit tests (panic on missing DB URL, defaults, override)
- `internal/model/post.go` — Post struct with DeletedAt soft-delete
- `internal/server/server.go` — New() with timeouts, GracefulShutdown() via signal.NotifyContext
- `web/embed.go` — go:embed for templates and static FS
- `internal/database/database.go` — Connect() with pgxpool.New + Ping
- `internal/database/database_test.go` — integration tests with testing.Short() guards
- `internal/database/migrations.go` — RunMigrations() via goose + stdlib.OpenDBFromPool
- `db/migrations/00001_create_posts.sql` — posts table DDL
- `db/migrations/embed.go` — embed.FS for SQL files (sibling package pattern)

## Decisions Made

- `go 1.26` in go.mod without a `toolchain` directive — avoids strict version enforcement on the local Go 1.23.3 toolchain while targeting the Docker build environment's Go 1.26
- Migrations embedded in `db/migrations/embed.go` as a sibling package rather than from `internal/database/` — Go's `go:embed` directive forbids path components containing `..`, so the embed FS must live alongside the files being embedded
- `goose.Up(sqlDB, ".")` path is `"."` because `migrations.FS` is rooted at the migrations directory itself (embed pattern `*.sql`), placing SQL files at the root of the FS
- web/templates and web/static received placeholder files (not `.gitkeep`) — `go:embed` requires at least one non-hidden embeddable file per embedded directory

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Embed path workaround for go:embed `..` restriction**
- **Found during:** Task 2 (migrations.go authoring)
- **Issue:** Plan specified `//go:embed ../../db/migrations/*.sql` in `internal/database/migrations.go` but Go forbids `..` in embed paths
- **Fix:** Created `db/migrations/embed.go` as a dedicated package exporting `var FS embed.FS` with `//go:embed *.sql`. `RunMigrations` imports this package and passes `migrations.FS` to `goose.SetBaseFS()`.
- **Files modified:** `db/migrations/embed.go` (created), `internal/database/migrations.go` (uses import)
- **Verification:** `go build ./...` exits 0
- **Committed in:** `7ebd1f8`

**2. [Rule 1 - Bug] web/static and web/templates need embeddable files, not .gitkeep**
- **Found during:** Task 1 verification (`go build ./...`)
- **Issue:** `go:embed` reports "contains no embeddable files" for directories containing only `.gitkeep` (hidden dot-files are excluded by embed)
- **Fix:** Added `web/static/placeholder.css` and `web/templates/placeholder.html` as minimal placeholder files. Removed `.gitkeep` from those dirs.
- **Files modified:** `web/static/placeholder.css`, `web/templates/placeholder.html`
- **Verification:** `go build ./...` exits 0
- **Committed in:** `e8ec780`

---

**Total deviations:** 2 auto-fixed (both Rule 1 — bugs in plan assumptions about Go toolchain behavior)
**Impact on plan:** Both fixes necessary for correct compilation. No scope creep.

## Issues Encountered

None beyond the auto-fixed embed path issues above.

## User Setup Required

None - no external service configuration required for this plan.

## Next Phase Readiness

- Plan 02 (markdown pipeline) can proceed: all goldmark/bluemonday deps are installed
- Plan 03 (Makefile/CI) can proceed: module path `github.com/jared-wallace/website-go` and build targets are stable
- Integration tests in `internal/database/` need `DATABASE_URL` set to run — requires `make dev-up` (docker-compose Postgres) per D-07

---

*Phase: 01-foundation*
*Completed: 2026-03-26*
