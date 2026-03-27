---
phase: 03-admin-panel
plan: "02"
subsystem: auth-ui
tags: [auth, sessions, host-routing, templates, css, scs, bcrypt, csrf]
dependency_graph:
  requires:
    - internal/middleware/auth.go
    - internal/middleware/ratelimit.go
    - internal/config/config.go
    - internal/service/post/service.go
  provides:
    - internal/handler/admin/handler.go
    - internal/handler/admin/auth.go
    - internal/handler/admin/handler_test.go
    - web/templates/admin-base.html
    - web/templates/admin-login.html
    - web/templates/admin-dashboard.html
    - web/templates/admin-editor.html
    - web/static/admin.css
  affects:
    - cmd/server/main.go
    - go.mod / go.sum
tech_stack:
  added:
    - github.com/alexedwards/scs/pgxstore v0.0.0-20251002162104-209de6e426de (direct -- Postgres session store)
  patterns:
    - hostRouter pattern: dispatches by Host header, strips port via net.SplitHostPort
    - Constant-time login: dummyHash pre-computed at startup, always call bcrypt.CompareHashAndPassword
    - CrossOriginProtection (stdlib Go 1.26) wraps admin mux for CSRF protection
    - SCS LoadAndSave wraps CrossOriginProtection wraps adminMux
    - Per-page template sets (admin-base.html + page template) mirrors BlogHandler pattern
    - testSetup pattern: handler and its session manager share same SCS instance in tests
key_files:
  created:
    - internal/handler/admin/handler.go
    - internal/handler/admin/auth.go
    - internal/handler/admin/handler_test.go
    - web/templates/admin-base.html
    - web/templates/admin-login.html
    - web/templates/admin-dashboard.html
    - web/templates/admin-editor.html
    - web/static/admin.css
  modified:
    - cmd/server/main.go
    - go.mod
    - go.sum
decisions:
  - "pgxstore re-added as direct dep (as predicted in Plan 01) when session store wired in cmd/server/main.go"
  - "testSetup pattern: handler and wrapper share the same SCS instance so h.sessions.GetBool works on the same context data the middleware injects"
  - "http.NewCrossOriginProtection() confirmed available in Go 1.26 stdlib -- no gorilla/csrf needed"
metrics:
  duration: "4 minutes"
  completed: "2026-03-27"
  tasks: 2
  files_changed: 10
---

# Phase 03 Plan 02: Admin Auth Flow and Template Foundation Summary

**One-liner:** SCS session auth with constant-time bcrypt login, host-based routing, CrossOriginProtection CSRF, admin base template with "Go Ashore" nav, nautical login page, and admin.css foundation — all auth plumbing Plans 03-04 depend on.

## Tasks Completed

| # | Name | Commit | Files |
|---|------|--------|-------|
| 1 | Admin handler, auth handlers, host router, main.go wiring | 982b7fb | handler.go, auth.go, handler_test.go, main.go, go.mod, go.sum |
| 2 | Admin base template, login template, stub templates, admin CSS | b9341d7 | admin-base.html, admin-login.html, admin-dashboard.html, admin-editor.html, admin.css |

## What Was Built

### AdminHandler (`internal/handler/admin/handler.go`)

Mirrors BlogHandler pattern exactly: per-page template sets (admin-base.html + page), funcMap with `formatDate`/`currentYear`, `render` method. Constructor pre-computes `dummyHash` via bcrypt cost 12 at startup so constant-time comparison is guaranteed regardless of whether the email matches.

Stub methods (returning 501) for all Dashboard/Editor/CRUD handlers — Plans 03-03 and 03-04 replace these.

### Auth Handlers (`internal/handler/admin/auth.go`)

`LoginPage`: checks `sm.GetBool(ctx, "authenticated")` and redirects immediately if already authed; pops `flash_error` from session for replay after redirect.

`LoginPost`: constant-time credential check (always runs bcrypt, uses dummyHash when email mismatches), rate-limits by IP via `extractIP` (X-Real-IP → RemoteAddr fallback), calls `RenewToken` before setting authenticated flag to prevent session fixation, redirects 303 to `/admin/posts`.

`Logout`: `sm.Destroy`, redirect 303 to `/admin/login`.

### Host Router (`cmd/server/main.go`)

`hostRouter` struct with `ServeHTTP` dispatching by `r.Host` (port-stripped via `net.SplitHostPort`). Admin traffic goes to `sessionManager.LoadAndSave(cop.Handler(adminMux))`; all other traffic goes to `blogMux`.

### Session Wiring

`pgxstore.New(pool)` for Postgres-backed sessions. `IdleTimeout: 24h`, `Lifetime: 30*24h`, `SameSite: Lax`, `HttpOnly: true`, `Secure: cfg.AppEnv == "production"`.

### CrossOriginProtection

`http.NewCrossOriginProtection()` (Go 1.26 stdlib) wraps `adminMux`. `AddTrustedOrigin("https://" + cfg.AdminHost)`. Resolves STATE.md blocker: "confirm before reaching for gorilla/csrf during Phase 3 planning."

### Templates

`admin-base.html`: nautical nav "The Log -- Back Office" + "Go Ashore" logout button, FlashSuccess/FlashError sections, admin.css + admin.js references.

`admin-login.html`: centered login form with Playfair Display "The Log" heading, "Back Office" subtitle, labeled email/password fields, error flash, "Sign In" CTA button.

Stub templates `admin-dashboard.html` and `admin-editor.html` provide the minimum `{{define "content"}}` blocks so `AdminHandler.New()` parses cleanly at startup.

### admin.css

`--color-danger` token (#9B1C1C light / #EF4444 dark), `.admin-login` centering layout, `.form-field` with explicit label/input styles, `.cta-button` with `min-height: 44px`, `.logout-link` with `min-height: 44px; min-width: 44px`, `.flash-error` and `.flash-success`, `.admin-nav` flex layout.

## Test Coverage

- `TestLoginPageRendersForm` — GET /admin/login returns 200, body contains "Sign In"
- `TestLoginPostInvalidCredentials` — POST with wrong credentials returns 401, body contains "Invalid email or password"
- `TestLoginPostValidCredentials` — POST with correct credentials returns 303 to /admin/posts
- `TestLogout` — POST /admin/logout returns 303 to /admin/login

All tests pass with `-race` flag. Tests use in-memory SCS store (no database required). The `testSetup` pattern ensures handler and LoadAndSave wrapper share the same SCS instance — required for session context to work correctly.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Test isolation: handler and session wrapper must share same SCS instance**
- **Found during:** Task 1 test execution — `panic: scs: no session data in context`
- **Issue:** Tests created a fresh `newTestSessionManager()` for the `LoadAndSave` wrapper while the handler used a different SCS instance from `newTestHandler()`. SCS panics when session data is absent because the middleware that loads it wasn't used.
- **Fix:** Introduced `testSetup` struct that wires handler and session manager from the same SCS instance; `serve()` helper calls `ts.sm.LoadAndSave` on the handler. Mirrors the runtime behavior where `adminHandler = sessionManager.LoadAndSave(cop.Handler(adminMux))`.
- **Files modified:** `internal/handler/admin/handler_test.go`
- **Commit:** 982b7fb (included in same commit since tests were being written for the first time)

## Known Stubs

- `web/templates/admin-dashboard.html` — stub template with `{{define "content"}}` placeholder; Plan 03-03 replaces with full dashboard
- `web/templates/admin-editor.html` — stub template with `{{define "content"}}` placeholder; Plan 03-04 replaces with full editor
- `AdminHandler.Dashboard`, `NewPost`, `EditPost`, `SavePost`, `DeletePost`, `RestorePost`, `PublishPost`, `UnpublishPost`, `Preview` — all return 501 Not Implemented; Plans 03-03 and 03-04 replace these

These stubs are intentional and required: `AdminHandler.New()` parses all templates at startup. The stubs prevent parse failures while their full implementations are deferred to later plans.

## Self-Check: PASSED

- [x] `internal/handler/admin/handler.go` — exists
- [x] `internal/handler/admin/auth.go` — exists
- [x] `internal/handler/admin/handler_test.go` — exists
- [x] `web/templates/admin-base.html` — exists
- [x] `web/templates/admin-login.html` — exists
- [x] `web/templates/admin-dashboard.html` — exists
- [x] `web/templates/admin-editor.html` — exists
- [x] `web/static/admin.css` — exists
- [x] Commit 982b7fb — verified via git log
- [x] Commit b9341d7 — verified via git log
- [x] `go build ./...` — passes
- [x] `go test ./internal/handler/admin/... -v -race` — PASS (4 tests)
