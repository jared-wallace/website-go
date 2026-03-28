# Technology Stack

**Project:** website-go — personal blog server for jared-wallace.com
**Researched:** 2026-03-26
**Philosophy:** stdlib-first; add a dependency only when it earns its place

---

## Recommended Stack

### Runtime

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| Go | 1.26.1 | Language / runtime | Latest stable (released 2026-03-05); Go 1.22+ ServeMux rewrites eliminate the need for a router package; single static binary |

### HTTP Routing

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `net/http` stdlib | (Go 1.26) | HTTP server and routing | Go 1.22 added method-based patterns, wildcard path segments (`{id}`), and `r.PathValue()`. For a blog with ~15 routes this is sufficient without pulling in chi or gorilla/mux |

The standard `http.ServeMux` now handles `GET /posts/{slug}`, `POST /admin/posts/{id}`, etc. No external router is needed for this scope.

### Templating

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `html/template` stdlib | (Go 1.26) | Server-side HTML rendering | Automatic context-aware HTML escaping; supports template inheritance via `{{block}}` / `{{define}}`; zero dependencies; XSS protection by default |

Avoid `text/template` for HTML output — it does not escape. Avoid `templ` (compile-time code generation) unless the team wants a build-time step; html/template is sufficient for a blog.

### Markdown Rendering

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `github.com/yuin/goldmark` | v1.8.2 | Markdown → HTML | CommonMark-compliant; extensible AST; actively maintained; used by Hugo. Blackfriday is archived and not CommonMark-compliant. Goldmark ships with tables, strikethrough, task lists |

Extensions to include at wiring time:
- `goldmark/extension` — GFM tables, strikethrough, linkify
- `goldmark-meta` (`github.com/yuin/goldmark-meta`) — YAML front matter extraction for post metadata
- `goldmark-highlighting` or `github.com/alecthomas/chroma` — syntax highlighting in code fences (stretch goal)

### Database Driver

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `github.com/jackc/pgx/v5` | v5.9.1 | PostgreSQL driver + toolkit | `lib/pq` is maintenance-only. pgx v5 is faster (up to 70x in bulk ops), supports LISTEN/NOTIFY, COPY, and native `pgxpool`. Use the native pgx interface, not the `database/sql` adapter, to avoid feature loss |
| `github.com/jackc/pgx/v5/pgxpool` | (bundled with pgx v5) | Connection pooling | Single-instance blog needs only a small pool (max 10 connections); pgxpool handles lifetime and health checks |

Do not use GORM or sqlx. Raw SQL with pgx is readable, auditable, and avoids ORM magic for a schema this simple.

### Database Migrations

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `github.com/pressly/goose/v3` | v3.27.0 | Schema migrations | Sequential numbered SQL files (`00001_create_posts.sql`); supports embedded migrations via `go:embed`; single binary deployment; `goose up` in container entrypoint. golang-migrate is a valid alternative but goose has cleaner Go embedding support |

### Session Management

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `github.com/alexedwards/scs/v2` | v2.9.0 | HTTP session management | OWASP-aligned design; server-side session tokens (not cookie-payload); context-middleware pattern; has a `pgxstore` sub-package for Postgres-backed sessions. Gorilla/sessions stores payload in cookie by default which is less secure; SCS is also faster and smaller |
| `github.com/alexedwards/scs/pgxstore` | (bundled in scs repo) | Postgres session backend | Stores sessions in `sessions` table alongside app data; survives container restarts; eliminates need for Redis |

### Password Hashing

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `golang.org/x/crypto/bcrypt` | v0.49.0 | Admin password hash | The standard Go bcrypt implementation; not in `crypto/` stdlib but from the official `x/crypto` extended library (Anthropic-controlled). Use cost 12 minimum for 2025 hardware |

This is one dependency that cannot be replaced by stdlib — Go's `crypto/` package does not include bcrypt. `x/crypto` is maintained by the Go team and is effectively stdlib-adjacent.

### RSS Feed

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `encoding/xml` stdlib | (Go 1.26) | RSS 2.0 feed generation | An RSS feed is a straightforward XML document; defining a Go struct with XML annotations and using `xml.NewEncoder` requires zero dependencies. `gorilla/feeds` adds a dependency for ~50 lines of struct definitions |

Hand-roll RSS using `encoding/xml`. The RSS 2.0 spec is stable and the struct is small. This removes a dependency entirely.

### Logging

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `log/slog` stdlib | (Go 1.26) | Structured logging | Added in Go 1.21; JSON or text output; leveled; zero dependencies. No need for zerolog or zap at blog scale |

### Infrastructure

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| Docker (multi-stage) | 27.x | Container build | Build stage: `golang:1.26-alpine`; runtime stage: `alpine:3.21` (not scratch — needs CA certs and timezone data for Postgres TLS). Set `CGO_ENABLED=0` for a static binary |
| docker-compose | v2.x | Local dev + prod orchestration | Sidecar Postgres container wired to EBS volume at `/var/www/html/postgres-data`; matches production topology in dev |
| Makefile | (GNU make) | Task runner | `make dev`, `make test`, `make build`, `make migrate-up`; documents tribal knowledge |

---

## Alternatives Considered

| Category | Recommended | Alternative | Why Not |
|----------|-------------|-------------|---------|
| Router | `net/http` ServeMux (stdlib) | `go-chi/chi`, `gorilla/mux` | Not needed — Go 1.22 ServeMux handles method+path patterns. Chi is excellent but adds a dependency for functionality now in stdlib |
| Templating | `html/template` (stdlib) | `templ` | Templ requires a code-generation build step and additional toolchain dependency. html/template is sufficient for a blog |
| Markdown | goldmark | blackfriday v2 | Blackfriday is not CommonMark-compliant and is effectively unmaintained. Gitea migrated away from it in 2020 |
| DB Driver | pgx v5 | `lib/pq`, GORM | lib/pq is maintenance-only. GORM is an ORM — unnecessary complexity for a blog schema |
| Migrations | goose v3 | golang-migrate | Both are valid; goose has cleaner `go:embed` support for embedded SQL migrations in a single binary |
| Sessions | scs v2 | gorilla/sessions, roll-your-own | gorilla/sessions stores payload in cookie (larger surface area); gorilla toolkit was briefly unmaintained. SCS is focused, secure, and has a pgxstore |
| RSS | `encoding/xml` stdlib | gorilla/feeds | gorilla/feeds is ~50 lines of struct defs; not worth the dependency |
| Logging | `log/slog` stdlib | zerolog, zap | At blog scale, slog is fast enough and avoids a dep. Revisit if structured log ingestion becomes a requirement |
| Image storage | EBS volume (`/var/www/html/images`) | S3 | Per PROJECT.md decision: EBS is sufficient for a single-blog, avoids AWS SDK dependency |

---

## Complete Dependency List

```
github.com/yuin/goldmark           v1.8.2
github.com/yuin/goldmark-meta      (latest, tracks goldmark)
github.com/jackc/pgx/v5            v5.9.1
github.com/alexedwards/scs/v2      v2.9.0
golang.org/x/crypto                v0.49.0
github.com/pressly/goose/v3        v3.27.0
```

Six external dependencies total. Everything else is stdlib.

---

## Installation

```bash
go get github.com/yuin/goldmark@v1.8.2
go get github.com/yuin/goldmark-meta@latest
go get github.com/jackc/pgx/v5@v5.9.1
go get github.com/alexedwards/scs/v2@v2.9.0
go get github.com/alexedwards/scs/pgxstore@latest
go get golang.org/x/crypto@v0.49.0
go get github.com/pressly/goose/v3@v3.27.0

# CLI tool for running migrations
go install github.com/pressly/goose/v3/cmd/goose@latest
```

---

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Go version (1.26.1) | HIGH | Verified via go.dev/doc/devel/release (2026-03-05) |
| net/http ServeMux routing | HIGH | Verified via go.dev/blog/routing-enhancements; official blog |
| goldmark v1.8.2 | HIGH | Verified via pkg.go.dev (published 2026-03-25) |
| pgx v5.9.1 | HIGH | Verified via pkg.go.dev (published 2026-03-22) |
| scs v2.9.0 | HIGH | Verified via pkg.go.dev (published 2025-04-17) |
| goose v3.27.0 | HIGH | Verified via pkg.go.dev (published 2026-02-22) |
| golang.org/x/crypto v0.49.0 | HIGH | Verified via pkg.go.dev (published 2026-03-11) |
| RSS via encoding/xml | HIGH | Standard library, well-documented pattern |
| Docker alpine runtime base | MEDIUM | Alpine 3.21 assumed current; verify at build time |

---

## Sources

- [Go Release History](https://go.dev/doc/devel/release) — Go 1.26.1 confirmed latest
- [Routing Enhancements for Go 1.22](https://go.dev/blog/routing-enhancements) — Official Go blog on ServeMux improvements
- [goldmark on pkg.go.dev](https://pkg.go.dev/github.com/yuin/goldmark) — v1.8.2, March 2026
- [pgx on pkg.go.dev](https://pkg.go.dev/github.com/jackc/pgx/v5) — v5.9.1, March 2026
- [scs on pkg.go.dev](https://pkg.go.dev/github.com/alexedwards/scs/v2) — v2.9.0
- [goose on pkg.go.dev](https://pkg.go.dev/github.com/pressly/goose/v3) — v3.27.0
- [golang.org/x/crypto/bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) — v0.49.0
- [The Go Ecosystem in 2025 — JetBrains GoLand Blog](https://blog.jetbrains.com/go/2025/11/10/go-language-trends-ecosystem-2025/)
- [SCS Session Manager — Alex Edwards](https://www.alexedwards.net/blog/scs-session-manager)
- [pgx vs lib/pq — Preslav Rachev](https://preslav.me/2022/05/13/pq-or-pgx-choosing-the-right-postgresql-golang-driver/)
- [Go net/http ServeMux is All You Need — DEV Community](https://dev.to/leapcell/gos-httpservemux-is-all-you-need-1mam)
- [Gitea migration from blackfriday to goldmark](https://github.com/go-gitea/gitea/pull/9533)
