---
phase: 03-admin-panel
plan: "01"
subsystem: data-layer
tags: [repository, service, middleware, auth, migrations, slug, tdd]
dependency_graph:
  requires: []
  provides:
    - internal/repository/post/write.go
    - internal/service/post/write.go
    - internal/service/post/slug.go
    - internal/middleware/auth.go
    - internal/middleware/ratelimit.go
    - db/migrations/00003_create_sessions.sql
    - cmd/hashpw/main.go
  affects:
    - internal/config/config.go
    - internal/repository/post/repository.go
    - internal/service/post/service.go
    - cmd/server/main.go
tech_stack:
  added:
    - scs/v2 v2.9.0 (direct — session middleware)
    - golang.org/x/crypto v0.49.0 (direct — bcrypt for hashpw)
  patterns:
    - TDD red-green for slug and service write tests
    - Repository interface extended without breaking existing read paths
    - Renderer interface in service layer enables testable markdown pre-rendering
    - Fixed-window rate limiter with sync.Mutex, per-IP, port-stripping
key_files:
  created:
    - db/migrations/00003_create_sessions.sql
    - internal/repository/post/write.go
    - internal/service/post/slug.go
    - internal/service/post/slug_test.go
    - internal/service/post/write.go
    - internal/service/post/write_test.go
    - internal/middleware/auth.go
    - internal/middleware/auth_test.go
    - internal/middleware/ratelimit.go
    - internal/middleware/ratelimit_test.go
    - cmd/hashpw/main.go
  modified:
    - internal/config/config.go
    - internal/repository/post/repository.go
    - internal/service/post/service.go
    - internal/service/post/service_test.go
    - internal/handler/blog/handler_test.go
    - cmd/server/main.go
    - Makefile
    - go.mod / go.sum
decisions:
  - "Service.New now accepts a Renderer interface (not *markdown.Renderer) — allows mock in tests without importing the markdown package"
  - "pgxstore removed from go.mod by go mod tidy — not imported yet; Plan 02 will re-add as direct dep when wiring session store"
  - "RequireSession test uses two-request session cookie pattern (set then access) — avoids httptest session store setup complexity"
metrics:
  duration: "6 minutes"
  completed: "2026-03-27"
  tasks: 2
  files_changed: 18
---

# Phase 03 Plan 01: Admin Data Infrastructure Summary

**One-liner:** Repository writes, service admin layer, slug generation, sessions migration, RequireSession middleware, fixed-window rate limiter, and bcrypt hashpw CLI — all the plumbing Plans 02-04 depend on.

## Tasks Completed

| # | Name | Commit | Files |
|---|------|--------|-------|
| 1 | Config, migration, repository writes, service writes, slug generation | 9606377 | config.go, 00003_create_sessions.sql, repository.go, write.go (repo), service.go, write.go (svc), slug.go, slug_test.go, write_test.go, service_test.go, handler_test.go, main.go |
| 2 | RequireSession middleware, rate limiter, hashpw tool | cd8c9cd | auth.go, auth_test.go, ratelimit.go, ratelimit_test.go, cmd/hashpw/main.go, Makefile, go.mod, go.sum |

## What Was Built

### Config Extension

`internal/config/config.go` now carries four admin fields loaded via `envOr` (non-panicking): `AdminEmail`, `AdminPasswordHash`, `AdminHost` (defaults to `admin.jared-wallace.com`), and `SessionSecret`. The server logs an INFO message when `AdminEmail` is empty so the public blog continues to serve while the admin panel is implicitly disabled.

### Sessions Migration

`db/migrations/00003_create_sessions.sql` creates the `sessions` table required by pgxstore (token, data, expiry) and a `sessions_expiry_idx` for efficient TTL cleanup.

### Repository Interface — 10 Methods

The `Repository` interface in `internal/repository/post/repository.go` now exposes 10 methods (3 existing + 7 new: `FindByID`, `ListAll`, `Create`, `Update`, `SoftDelete`, `Restore`, `SetPublished`). All 7 are implemented on `*postgresRepository` in `write.go`. `Create` returns `ErrSlugExists` (wrapping Postgres error code 23505) for unique-slug conflicts.

### Service Write Layer

`internal/service/post/write.go` exposes 8 admin operations. `Create` and `Update` call `renderer.Render(body)` before storing, pre-computing `RenderedHTML`. The `Renderer` interface in `service.go` keeps the markdown package out of test imports.

### Slug Generation

`internal/service/post/slug.go` — `GenerateSlug` lowercases, replaces non-alphanumeric runs with hyphens, trims, and falls back to `post-{unix_timestamp}` for empty input. Regex compiled once at package level.

### RequireSession Middleware

`internal/middleware/auth.go` — checks `sm.GetBool(r.Context(), "authenticated")`; redirects to `/admin/login` (303) if false. Fully composable: returns `func(http.Handler) http.Handler`.

### Rate Limiter

`internal/middleware/ratelimit.go` — fixed-window counter per IP. `Allow(ip)` strips port, resets count after window expiry, blocks when `count >= limit`. Thread-safe via `sync.Mutex`.

### hashpw CLI Tool

`cmd/hashpw/main.go` — `go run ./cmd/hashpw <password>` prints a bcrypt cost-12 hash to stdout. Wired into `make hash-password PW=<password>`.

## Test Coverage

- `slug_test.go` — 6 cases covering normal slugification, empty input fallback
- `write_test.go` — 7 tests using mockRepo + mockRenderer; verifies delegation, error propagation, and RenderedHTML population for Create/Update
- `auth_test.go` — 2 tests: unauthenticated (303 redirect), authenticated (200 pass-through via two-request session cookie pattern)
- `ratelimit_test.go` — 4 tests: limit enforcement, window reset, IP isolation, port stripping

All tests pass with `-race` flag.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 2 - Missing critical functionality] Extended handler_test.go and service_test.go mocks to satisfy new Repository interface**
- **Found during:** Task 1 — extending Repository interface broke existing mock implementations in two test files
- **Issue:** `handler_test.go::mockRepository` and `service_test.go::mockRepository` only implemented 3 methods; adding 7 new interface methods caused compile failure
- **Fix:** Added the 7 new methods (returning `errors.New("not implemented")`) to both mocks; added `noopRenderer` to both test packages; updated `postservice.New(repo)` calls to `postservice.New(repo, noopRenderer{})`
- **Files modified:** `internal/handler/blog/handler_test.go`, `internal/service/post/service_test.go`
- **Commit:** 9606377

**2. [Rule 1 - Bug] pgxstore dropped from go.mod after go mod tidy**
- **Found during:** Task 2 — `go mod tidy` removed `scs/pgxstore` because no Go file imports it yet
- **Issue:** pgxstore was `// indirect` before; since no code in this plan imports it, tidy correctly removes it
- **Fix:** Accepted removal as correct — Plan 02 will re-add it as direct dependency when wiring the session store in `cmd/server/main.go`
- **Impact:** None for this plan; Plans 02+ will restore it

## Known Stubs

None. All implementations are complete — no hardcoded empty values or placeholders that flow to rendering.

## Self-Check: PASSED

- [x] `internal/repository/post/write.go` — exists
- [x] `internal/service/post/write.go` — exists
- [x] `internal/service/post/slug.go` — exists
- [x] `internal/middleware/auth.go` — exists
- [x] `internal/middleware/ratelimit.go` — exists
- [x] `db/migrations/00003_create_sessions.sql` — exists
- [x] `cmd/hashpw/main.go` — exists
- [x] Commit 9606377 — verified via git log
- [x] Commit cd8c9cd — verified via git log
- [x] `go build ./...` — passes
- [x] `go test ./internal/service/post/... -v -race` — PASS
- [x] `go test ./internal/middleware/... -v -race` — PASS
- [x] `go run ./cmd/hashpw "testpass"` — outputs `$2a$12$...`
