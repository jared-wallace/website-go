# Phase 1: Foundation - Research

**Researched:** 2026-03-26
**Domain:** Go project scaffold — module layout, pgx/goose database, goldmark+bluemonday markdown pipeline, Makefile dev workflow, GHA CI
**Confidence:** HIGH

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

- **D-01:** Flat `internal/` structure — `internal/server/`, `internal/handler/`, `internal/markdown/`, `internal/database/`, `internal/model/`. Add nesting only when a package actually needs sub-packages.
- **D-02:** Templates and static assets embedded via `go:embed` in `web/` package. Single binary deployment, no runtime file dependencies.
- **D-03:** Configuration via environment variables only (`os.Getenv` with helpers). No config files, no viper. 12-factor style with `envOr()` defaults and `mustEnv()` for required values.
- **D-04:** All four goldmark extensions enabled: GFM tables + strikethrough, linkify, syntax highlighting (chroma), and YAML front matter (goldmark-meta).
- **D-05:** bluemonday `UGCPolicy()` for HTML sanitization. Allows standard formatting (bold, italic, links, images, code, tables, blockquotes) while stripping scripts, iframes, event handlers.
- **D-06:** Hot reload via `cosmtrek/air` for `make dev`. Watches `.go`, `.html`, `.css` files and auto-rebuilds. `.air.toml` committed to repo.
- **D-07:** Local Postgres via docker-compose only (`docker-compose.dev.yml`). `make dev-up` starts the container. Matches production topology.
- **D-08:** `golangci-lint` for `make lint` with `.golangci.yml` config. Enables errcheck, govet, staticcheck, gosimple, ineffassign, unused.
- **D-09:** Postgres service container in CI. Tests run against a real database — catches migration and query issues.
- **D-10:** Single CI job running lint -> test -> build sequentially. Conserves free-tier minutes. Go module caching via `actions/setup-go` with `cache: true`.

### Claude's Discretion

- Specific Makefile target names and help text formatting
- `.golangci.yml` exact linter set (beyond the base six discussed)
- `internal/config/` package shape (struct fields, validation approach)
- Test file organization and naming conventions
- goose migration numbering format

### Deferred Ideas (OUT OF SCOPE)

None — discussion stayed within phase scope.
</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| FOUND-01 | Project follows standard Go layout (cmd/, internal/, etc.) | Directory structure pattern from ARCHITECTURE.md; D-01 locks the `internal/` shape |
| FOUND-02 | Postgres connection pool via pgx/v5 with health checks | pgx v5.9.1 + pgxpool; `pgxpool.New()` with `Ping()` health check pattern |
| FOUND-03 | Database migrations via goose with versioned SQL files | goose v3.27.0 with `go:embed` embedded SQL; `goose.Up()` at startup |
| FOUND-06 | Makefile with build, test, lint, run, docker, and migration targets | GNU Make 3.81 available; Air and goose CLI need installation instructions |
| FOUND-07 | GHA CI pipeline running lint, test, and build on push | `actions/setup-go` with cache; Postgres service container; single job |
</phase_requirements>

---

## Summary

Phase 1 establishes the skeleton that all subsequent phases hang on. The work falls into five parallel streams that converge into a single compilable binary: (1) Go module init and directory layout, (2) Postgres connection + embedded goose migrations, (3) goldmark + bluemonday markdown pipeline with tests, (4) Makefile with dev tooling, (5) GHA CI pipeline. No HTTP routes serving public content yet — those come in Phase 2.

The stack is fully locked by prior research (STACK.md, ARCHITECTURE.md). Every dependency version is verified against pkg.go.dev as of March 2026. The environment audit below shows that the dev machine has Go 1.23.3 (not the locked 1.26.1), Docker 28.1.1, golangci-lint v1.61.0, and Make 3.81. Air and goose CLI are absent and must be installed. The module will declare `go 1.26.1` in go.mod to match the Docker build environment; local dev on 1.23.3 may surface minor compat differences but the locked versions are all compatible with 1.23.x.

The single most important correctness gate for this phase is the XSS test: a Go test that feeds `<script>alert(1)</script>` into the markdown pipeline and asserts no `<script>` tag appears in the output. This test encodes Decision D-05 as an executable contract and will catch any accidental re-ordering of goldmark → bluemonday.

**Primary recommendation:** Initialize the Go module with `go 1.26`, lay out the directory skeleton per D-01, wire dependencies in dependency order (config → db → migrations → markdown → main), and land the GHA CI pipeline before closing the phase.

---

## Standard Stack

### Core

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| Go | 1.26.1 | Language / runtime | Latest stable; `go:embed`, `log/slog`, enhanced ServeMux all in stdlib |
| `github.com/yuin/goldmark` | v1.8.2 | Markdown → HTML | CommonMark-compliant; extensible; Hugo uses it |
| `github.com/yuin/goldmark-meta` | latest (tracks goldmark) | YAML front matter | Official goldmark extension for post metadata extraction |
| `github.com/alecthomas/chroma/v2` | latest | Syntax highlighting | Used by goldmark-highlighting; pairs with goldmark-highlighting extension |
| `github.com/microcosm-cc/bluemonday` | latest | HTML sanitization | Purpose-built Go HTML sanitizer; UGCPolicy baseline |
| `github.com/jackc/pgx/v5` | v5.9.1 | Postgres driver + pool | lib/pq is maintenance-only; pgxpool bundled; faster |
| `github.com/pressly/goose/v3` | v3.27.0 | DB migrations | go:embed-friendly; sequential numbered SQL files |

### Supporting

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| `github.com/alexedwards/scs/v2` | v2.9.0 | Session management | Needed at Phase 3 (admin); wire the package now so migrations can include the sessions table |
| `golang.org/x/crypto` | v0.49.0 | bcrypt password hashing | Phase 3 usage; install now to keep go.sum clean |
| `log/slog` stdlib | Go 1.26 | Structured logging | Zero deps; `slog.New(slog.NewTextHandler(os.Stderr, nil))` |

### Alternatives Considered

| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| goldmark-meta | Manual front matter parsing | goldmark-meta integrates cleanly with goldmark's AST; hand-rolling is ~100 lines of error-prone regex |
| bluemonday | No sanitizer (trust admin content) | Eliminates XSS vector from API push endpoint compromise; D-05 locks this |
| golangci-lint | `go vet` + `staticcheck` separately | golangci-lint orchestrates multiple linters with caching; single binary |

**Installation:**
```bash
go get github.com/yuin/goldmark@v1.8.2
go get github.com/yuin/goldmark-meta@latest
go get github.com/yuin/goldmark-highlighting/v2@latest
go get github.com/alecthomas/chroma/v2@latest
go get github.com/microcosm-cc/bluemonday@latest
go get github.com/jackc/pgx/v5@v5.9.1
go get github.com/pressly/goose/v3@v3.27.0
go get github.com/alexedwards/scs/v2@v2.9.0
go get github.com/alexedwards/scs/pgxstore@latest
go get golang.org/x/crypto@v0.49.0

# CLI tools for dev workflow
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/cosmtrek/air@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

---

## Architecture Patterns

### Recommended Project Structure

```
website-go/
├── cmd/
│   └── server/
│       └── main.go           # Entry point — DI wiring only, ~50 lines
├── internal/
│   ├── config/               # Env-var loading: envOr(), mustEnv(), Config struct
│   ├── database/             # pgxpool connect, Ping health check, embedded migrations
│   ├── markdown/             # Goldmark wiring, bluemonday sanitize, Render(string) string
│   ├── handler/              # HTTP handlers (Phase 2+; stub package acceptable in Phase 1)
│   ├── model/                # Shared domain structs (Post, etc.)
│   └── server/               # http.Server construction with timeouts
├── db/
│   └── migrations/
│       └── 00001_initial_schema.sql
├── web/
│   ├── templates/            # html/template files (embedded, Phase 2+)
│   └── static/               # CSS, JS (embedded, Phase 2+)
├── .github/
│   └── workflows/
│       └── ci.yml
├── .air.toml                 # Air hot-reload config
├── .golangci.yml             # Linter config
├── docker-compose.dev.yml    # Local Postgres sidecar
├── Dockerfile                # Multi-stage (Phase 6; scaffold now)
├── Makefile
└── go.mod
```

### Pattern 1: Config via envOr / mustEnv

**What:** Thin helpers over `os.Getenv` — no third-party config library.
**When to use:** All configuration. Load once in `main.go`, pass `Config` struct down via constructors.

```go
// internal/config/config.go
package config

import (
    "fmt"
    "os"
)

type Config struct {
    DatabaseURL string
    Port        string
    AppEnv      string
}

func Load() Config {
    return Config{
        DatabaseURL: mustEnv("DATABASE_URL"),
        Port:        envOr("PORT", "8080"),
        AppEnv:      envOr("APP_ENV", "development"),
    }
}

func envOr(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}

func mustEnv(key string) string {
    v := os.Getenv(key)
    if v == "" {
        panic(fmt.Sprintf("required env var %q is not set", key))
    }
    return v
}
```

### Pattern 2: pgxpool Connect with Health Check

**What:** Open a pgxpool at startup, Ping to verify connectivity before accepting requests.
**When to use:** Database layer initialization in `main.go`.

```go
// internal/database/database.go
package database

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
    pool, err := pgxpool.New(ctx, dsn)
    if err != nil {
        return nil, fmt.Errorf("pgxpool.New: %w", err)
    }
    if err := pool.Ping(ctx); err != nil {
        pool.Close()
        return nil, fmt.Errorf("database ping failed: %w", err)
    }
    return pool, nil
}
```

### Pattern 3: Goose Embedded Migrations

**What:** Embed SQL migration files at compile time; run `goose.Up()` at startup before the HTTP server starts.
**When to use:** Database initialization in `main.go`, before any handler registration.

```go
// internal/database/migrations.go
package database

import (
    "context"
    "embed"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/pressly/goose/v3"
    "github.com/pressly/goose/v3/database"
)

//go:embed ../../db/migrations/*.sql
var migrationsFS embed.FS

func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
    goose.SetBaseFS(migrationsFS)
    if err := goose.SetDialect("postgres"); err != nil {
        return err
    }
    db, err := pool.Acquire(ctx)
    if err != nil {
        return err
    }
    defer db.Release()
    return goose.Up(db.Conn(), "db/migrations")
}
```

> Note: goose v3 has multiple migration runner APIs. Verify the exact API for running against a `pgxpool`-acquired connection in the goose v3.27.0 docs — the embed path must match the directory relative to the embed directive location.

### Pattern 4: Goldmark + Bluemonday Pipeline

**What:** Convert markdown to HTML with goldmark (all extensions), then sanitize with bluemonday `UGCPolicy()`. Return `template.HTML` for use in Go templates.
**When to use:** Post rendering — call at write time, store rendered HTML. Also tested directly in Phase 1.

```go
// internal/markdown/renderer.go
package markdown

import (
    "bytes"
    "html/template"

    "github.com/microcosm-cc/bluemonday"
    "github.com/yuin/goldmark"
    "github.com/yuin/goldmark-meta"
    "github.com/yuin/goldmark/extension"
    "github.com/yuin/goldmark/parser"
    "github.com/yuin/goldmark/renderer/html"
)

type Renderer struct {
    gm goldmark.Markdown
    bm *bluemonday.Policy
}

func NewRenderer() *Renderer {
    gm := goldmark.New(
        goldmark.WithExtensions(
            extension.GFM,
            extension.Linkify,
            meta.Meta,
            // goldmark-highlighting configured separately
        ),
        goldmark.WithRendererOptions(
            html.WithUnsafe(), // goldmark unsafe allows raw HTML — bluemonday sanitizes downstream
        ),
    )
    return &Renderer{
        gm: gm,
        bm: bluemonday.UGCPolicy(),
    }
}

func (r *Renderer) Render(src string) template.HTML {
    var buf bytes.Buffer
    ctx := parser.NewContext()
    if err := r.gm.Convert([]byte(src), &buf, parser.WithContext(ctx)); err != nil {
        return template.HTML("<!-- render error -->")
    }
    sanitized := r.bm.SanitizeBytes(buf.Bytes())
    return template.HTML(sanitized)
}
```

### Pattern 5: http.Server with Timeouts

**What:** Explicit `ReadTimeout`, `WriteTimeout`, `IdleTimeout` on the server struct.
**When to use:** Always — `http.ListenAndServe` has no timeout parameters and is a production trap.

```go
srv := &http.Server{
    Addr:         ":" + cfg.Port,
    Handler:      mux,
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  60 * time.Second,
}
```

### Pattern 6: Graceful Shutdown

**What:** Intercept `SIGINT`/`SIGTERM` and call `srv.Shutdown()` with a deadline so in-flight requests complete.
**When to use:** Alongside the `http.Server` construction.

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()
go srv.ListenAndServe()
<-ctx.Done()
shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
srv.Shutdown(shutdownCtx)
```

### Anti-Patterns to Avoid

- **Global `var DB *pgxpool.Pool`:** Hidden dependency, untestable, race-prone. Pass pool via constructor injection.
- **`http.ListenAndServe(addr, mux)` bare call:** No timeout parameters — exposes the server to Slowloris. Always use `http.Server` struct.
- **Casting goldmark output to `template.HTML` directly without bluemonday:** Bypasses sanitization; stored XSS risk even for admin-authored content (API push vector).
- **`goose.SetBaseFS` called without matching the embed path prefix:** Common mis-configuration — the path passed to `goose.Up()` must reflect the directory as seen from the embed FS root.
- **Everything in `main.go`:** Phase 1 is where the `internal/` boundary gets established. Don't defer structure.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| HTML sanitization | Custom regex/tag stripper | bluemonday | Tag stripping via regex misses event handlers, javascript: URLs, data: URIs — it is a solved problem with a long CVE tail |
| DB migrations | Custom "run SQL files at startup" loop | goose v3 | Migration ordering, idempotency, version tracking, rollback support — 50+ edge cases to handle correctly |
| Session management | Cookie encrypt/decrypt, server-side token store | scs v2 | OWASP session management is a specification document, not a blog post — SCS implements it correctly |
| Structured logging | `fmt.Printf("[INFO] %s", msg)` | `log/slog` | JSON output, leveled filtering, zero overhead at below-threshold levels — free in stdlib |
| Config validation | `if os.Getenv("X") == "" { os.Exit(1) }` scattered across packages | `internal/config` with `mustEnv` | Centralized, testable, documented — scattered panics are untraceable |

**Key insight:** Go's stdlib-first philosophy pays off here. The only essential external dependencies for Phase 1 are goldmark (markdown), bluemonday (sanitizer), pgx (Postgres), and goose (migrations). Everything else — logging, config, HTTP server, embedding — uses stdlib.

---

## Common Pitfalls

### Pitfall 1: goldmark `html.WithUnsafe()` Without bluemonday Downstream

**What goes wrong:** `html.WithUnsafe()` is needed to allow raw HTML pass-through in markdown (e.g., `<figure>` tags). Without it, goldmark strips raw HTML. But if `WithUnsafe()` is set and bluemonday is absent or called before goldmark, the raw HTML — including any injected `<script>` tags — reaches the browser.

**Why it happens:** Developers add `WithUnsafe()` to fix a rendering issue and forget the sanitizer is the safety net.

**How to avoid:** The pipeline order is fixed: `goldmark.Convert() → bluemonday.Sanitize() → template.HTML`. Add a test (the XSS test in success criterion 3) that asserts `<script>` does not survive the pipeline.

**Warning signs:** Any markdown test that contains raw HTML passes without the `<script>` being stripped.

---

### Pitfall 2: goose embed path mismatch

**What goes wrong:** `goose.SetBaseFS(migrationsFS)` is set, but `goose.Up(db, "db/migrations")` uses a path that doesn't match the directory structure inside the embed FS.

**Why it happens:** The `//go:embed` directive captures files relative to the Go source file's package directory. If the migration package is `internal/database/` and the migrations are in `db/migrations/`, the embed path must be `../../db/migrations/*.sql` but the `goose.Up` path argument is the path *within* the FS, which depends on the embed directive's anchor.

**How to avoid:** Use `//go:embed db/migrations` from the project root package (`main` or a package at the root), or carefully match the relative path. Verify with a test that calls `goose.Up` against an in-memory Postgres and confirms the `goose_db_version` table exists.

**Warning signs:** `goose: no migration files found` at startup despite SQL files existing in the repository.

---

### Pitfall 3: go.mod declares 1.26 but local Go is 1.23

**What goes wrong:** `go.mod` has `go 1.26.1` (Docker build target). The developer runs `go build` locally with Go 1.23.3. Some toolchain features (workspace, new stdlib additions) may behave differently or be absent.

**Why it happens:** The locked version (1.26.1) matches Docker; the dev machine predates it. This is a known gap, not a blocker.

**How to avoid:** Do not use any 1.24+ specific stdlib APIs (verify on go.dev/doc/go1.24) that would break on 1.23. The locked stack avoids this — all dependencies target 1.21+ minimum. Add a `go version` check in Makefile `build` target that warns if local Go < 1.23.

**Warning signs:** `go: go.mod requires go >= 1.26; found go 1.23.3` if the toolchain directive is strict. Avoid using the `toolchain` directive (separate from `go` directive) to prevent this.

---

### Pitfall 4: golangci-lint version mismatch with Go version

**What goes wrong:** golangci-lint v1.61.0 (installed on dev machine) was built with Go 1.23.3. When CI uses a newer Go or a different golangci-lint version, linter behavior and available linters differ.

**Why it happens:** golangci-lint versions are tightly coupled to Go versions. CI and local must agree.

**How to avoid:** Pin golangci-lint version in GHA workflow: `golangci/golangci-lint-action@v6` with `version: v1.61.0`. Match the version in the Makefile install target.

**Warning signs:** CI lint passes but local lint fails (or vice versa) due to different available linters.

---

### Pitfall 5: GHA Postgres service not ready before go test

**What goes wrong:** The GHA workflow starts the Postgres service container and immediately runs `go test`. Postgres takes 2-5 seconds to initialize. The first test that dials the database fails with `connection refused`.

**Why it happens:** GHA service containers start asynchronously. The `services:` block doesn't guarantee the service is accepting connections when the first `run:` step executes.

**How to avoid:** Add a health check in the GHA service definition:
```yaml
services:
  postgres:
    image: postgres:16-alpine
    env:
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: website_test
    options: >-
      --health-cmd pg_isready
      --health-interval 10s
      --health-timeout 5s
      --health-retries 5
```
GHA waits for the health check to pass before running steps.

**Warning signs:** Intermittent CI failures on database connection in the first test, especially on cold runners.

---

## Code Examples

### go.mod initial state

```go
// go.mod
module github.com/jared-wallace/website-go

go 1.26
```

### Main.go wiring skeleton

```go
// cmd/server/main.go
package main

import (
    "context"
    "log/slog"
    "os"

    "github.com/jared-wallace/website-go/internal/config"
    "github.com/jared-wallace/website-go/internal/database"
)

func main() {
    log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))

    cfg := config.Load()

    ctx := context.Background()
    pool, err := database.Connect(ctx, cfg.DatabaseURL)
    if err != nil {
        log.Error("database connection failed", "err", err)
        os.Exit(1)
    }
    defer pool.Close()

    if err := database.RunMigrations(ctx, pool); err != nil {
        log.Error("migrations failed", "err", err)
        os.Exit(1)
    }

    log.Info("server starting", "port", cfg.Port)
    // Phase 2: wire HTTP server here
}
```

### Makefile skeleton

```makefile
# Makefile
.PHONY: build test lint run dev dev-up dev-down migrate help

BINARY := bin/server
GO     := go
AIR    := air

## build: compile the server binary
build:
	$(GO) build -o $(BINARY) ./cmd/server

## test: run all tests
test:
	$(GO) test ./... -v -race

## lint: run golangci-lint
lint:
	golangci-lint run ./...

## run: build and run the server
run: build
	./$(BINARY)

## dev: start server with hot reload (requires air)
dev:
	$(AIR) -c .air.toml

## dev-up: start local Postgres via docker-compose
dev-up:
	docker compose -f docker-compose.dev.yml up -d

## dev-down: stop local Postgres
dev-down:
	docker compose -f docker-compose.dev.yml down

## migrate: run goose migrations up
migrate:
	goose -dir db/migrations postgres "$(DATABASE_URL)" up

## docker: build Docker image
docker:
	docker build -t website-go:latest .

## help: list available targets
help:
	@grep -E '^## ' Makefile | sed 's/## //' | column -t -s ':'
```

### GHA CI workflow

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main, new]
  pull_request:
    branches: [main]

jobs:
  ci:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:16-alpine
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: website_test
          POSTGRES_USER: postgres
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432

    env:
      DATABASE_URL: postgres://postgres:postgres@localhost:5432/website_test?sslmode=disable

    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.26'
          cache: true

      - name: Lint
        uses: golangci/golangci-lint-action@v6
        with:
          version: v1.61.0

      - name: Test
        run: go test ./... -v -race

      - name: Build
        run: go build ./cmd/server
```

### XSS test for markdown pipeline

```go
// internal/markdown/renderer_test.go
package markdown_test

import (
    "strings"
    "testing"

    "github.com/jared-wallace/website-go/internal/markdown"
)

func TestRender_XSSStripped(t *testing.T) {
    r := markdown.NewRenderer()
    input := `Hello <script>alert(1)</script> world`
    output := string(r.Render(input))
    if strings.Contains(output, "<script>") {
        t.Errorf("expected <script> to be stripped, got: %s", output)
    }
}

func TestRender_BasicMarkdown(t *testing.T) {
    r := markdown.NewRenderer()
    output := string(r.Render("**bold** text"))
    if !strings.Contains(output, "<strong>bold</strong>") {
        t.Errorf("expected bold rendering, got: %s", output)
    }
}
```

### .air.toml

```toml
# .air.toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/server"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", ".git"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "html", "css"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
```

### .golangci.yml

```yaml
# .golangci.yml
linters:
  enable:
    - errcheck
    - govet
    - staticcheck
    - gosimple
    - ineffassign
    - unused
    - sqlclosecheck    # catches unclosed sql.Rows (Pitfall 6 from PITFALLS.md)
    - gosec            # security issues
    - gofmt
    - goimports

linters-settings:
  errcheck:
    check-type-assertions: true
  govet:
    enable-all: true

issues:
  exclude-rules:
    - path: "_test.go"
      linters:
        - errcheck
```

### docker-compose.dev.yml

```yaml
# docker-compose.dev.yml
services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: website
      POSTGRES_PASSWORD: website
      POSTGRES_DB: website_dev
    ports:
      - "5432:5432"
    volumes:
      - postgres_dev_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U website"]
      interval: 5s
      timeout: 5s
      retries: 5

volumes:
  postgres_dev_data:
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| `blackfriday` markdown | `goldmark` | 2020 (Gitea migration) | CommonMark compliance; extensible AST |
| `lib/pq` Postgres driver | `pgx/v5` | 2022 (lib/pq maintenance-only) | 70x bulk op speed; native pooling |
| `gorilla/sessions` | `alexedwards/scs` | ~2022 | Server-side tokens; OWASP alignment |
| `golang-migrate` | `goose/v3` | Ongoing preference | Cleaner `go:embed` support |
| `log` stdlib | `log/slog` | Go 1.21 (2023) | Structured, leveled, zero-dep |
| `chi` router | `net/http` ServeMux | Go 1.22 (2024) | Method + path patterns now in stdlib |

**Deprecated / outdated:**
- `blackfriday v2`: Not CommonMark-compliant, effectively unmaintained since ~2020. Gitea migrated away.
- `lib/pq`: Maintenance-only since 2022. pgx v5 is the community standard.
- `gorilla/sessions`: Still functional but stores payload in cookie by default; briefly had maintenance concerns.

---

## Open Questions

1. **goose v3 API for pgxpool-acquired connection**
   - What we know: goose v3.27.0 supports `go:embed` and Postgres. The standard API takes a `*database/sql.DB`.
   - What's unclear: The cleanest way to use goose with a `pgxpool.Pool` (which is not a `*sql.DB`). Options include: acquiring a stdlib-compatible connection via `pgxpool`'s `database/sql` adapter, or using `pgx/stdlib.OpenDBFromPool`.
   - Recommendation: Use `pgx/stdlib.OpenDBFromPool(pool)` to get a `*sql.DB` for goose only, keeping the main app on native pgx. Verify in goose v3 docs.

2. **golangci-lint v1.61.0 compatibility with go 1.26 module**
   - What we know: Dev machine has golangci-lint v1.61.0 (built with Go 1.23.3). Go module will declare 1.26.
   - What's unclear: Whether v1.61.0 will correctly analyze Go 1.26 code or produce false positives on new syntax.
   - Recommendation: Pin GHA CI to v1.61.0 to match local. If issues arise, upgrade golangci-lint via `go install` before the phase closes.

3. **goldmark-highlighting API in v2**
   - What we know: `github.com/yuin/goldmark-highlighting/v2` is the maintained version with chroma v2 support.
   - What's unclear: Exact import path and initialization pattern for the v2 API with custom chroma style.
   - Recommendation: Use the default chroma style (`monokai` or `github`) for Phase 1. Style customization is Phase 2 scope.

---

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|-------------|-----------|---------|----------|
| Go | All | Yes | 1.23.3 (local) / 1.26.1 (Docker) | — |
| Docker | dev-up, docker target | Yes | 28.1.1 | — |
| docker compose (v2) | dev-up, dev-down | Yes | 2.35.1 | — |
| GNU Make | All Makefile targets | Yes | 3.81 | — |
| golangci-lint | make lint | Yes | v1.61.0 | GHA golangci-lint-action installs correct version in CI |
| Air (hot reload) | make dev | No | — | `go run ./cmd/server` as fallback; install via `go install github.com/cosmtrek/air@latest` |
| goose CLI | make migrate | No | — | `go run github.com/pressly/goose/v3/cmd/goose` as inline fallback; install via `go install` |
| Postgres (local) | Tests, make dev | No (service) | — | docker compose dev-up provides it |

**Missing dependencies with no fallback:**
- None — all missing tools have installation paths or compose-based fallbacks.

**Missing dependencies with fallback:**
- Air: Install via `go install github.com/cosmtrek/air@latest`. The Makefile `dev` target should document this. Fallback: `go run ./cmd/server`.
- goose CLI: Install via `go install github.com/pressly/goose/v3/cmd/goose@latest`. Makefile `migrate` target should document this. Fallback: call goose programmatically from main (already done via embedded migrations).
- Local Postgres: Provided by `docker compose -f docker-compose.dev.yml up -d`. Not a missing dependency — it's the designed approach.

---

## Validation Architecture

### Test Framework

| Property | Value |
|----------|-------|
| Framework | `go test` (stdlib) — no external test framework |
| Config file | none (standard Go test convention) |
| Quick run command | `go test ./internal/markdown/... -v` |
| Full suite command | `go test ./... -v -race` |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|--------------|
| FOUND-01 | Standard Go layout compiles | smoke | `go build ./...` | Wave 0 — create directory structure |
| FOUND-02 | pgxpool connects and pings | integration | `go test ./internal/database/... -v` (requires Postgres) | Wave 0 — create test file |
| FOUND-03 | goose migrations run without error | integration | `go test ./internal/database/... -run TestMigrations -v` | Wave 0 — create test file |
| FOUND-06 | Makefile targets are functional | smoke | `make build && make lint && make test` | Wave 0 — create Makefile |
| FOUND-07 | GHA CI reports green | e2e / CI | Push to branch; observe Actions tab | Wave 0 — create .github/workflows/ci.yml |

**Note on FOUND-02/03:** Integration tests require a live Postgres. For local: `make dev-up` first. For CI: GHA Postgres service container (D-09). The success criterion "markdown pipeline XSS test" (criterion 3) is the most valuable unit test in Phase 1 — pure function, no Postgres required.

### Sampling Rate

- **Per task commit:** `go build ./...` (zero external deps)
- **Per wave merge:** `go test ./... -race` (requires local Postgres for DB tests)
- **Phase gate:** Full suite green (including CI) before `/gsd:verify-work`

### Wave 0 Gaps

All test files are missing — this is a greenfield project. The following must be created before implementation tasks can produce passing tests:

- [ ] `internal/markdown/renderer_test.go` — XSS test (FOUND-01 proxy; pure function, no Postgres)
- [ ] `internal/database/database_test.go` — connection + migration tests (FOUND-02, FOUND-03)
- [ ] `.github/workflows/ci.yml` — GHA CI (FOUND-07)
- [ ] `Makefile` — dev workflow (FOUND-06)
- [ ] `go.mod` + `go.sum` — module initialization (prerequisite for all)

---

## Project Constraints (from CLAUDE.md)

Directives from `./CLAUDE.md` that the planner must enforce:

| Directive | Requirement |
|-----------|-------------|
| Tech stack | Go with minimal dependencies — avoid large frameworks, prefer stdlib where reasonable |
| Infrastructure | Must run as Docker container on port 8080 behind existing Nginx/ALB |
| Budget | GHA CI must work on free tier (no paid GitHub features) |
| Storage | All persistent data must live on EBS volume at /var/www/html |
| Design | Leverage `frontend-design` skill for all template/UI work (Phase 2+ concern) |
| Code quality | All code changes must pass `/simplify` skill review before presenting to user |
| GSD workflow | Use GSD entry points for all repo edits — no direct edits outside a GSD workflow |

**Phase 1 specific:** No HTTP routes serving content = no templates rendered = `frontend-design` skill not required in Phase 1. The `/simplify` pass applies to all code written during execution.

---

## Sources

### Primary (HIGH confidence)

- `STACK.md` (`.planning/research/STACK.md`) — All dependency versions verified against pkg.go.dev March 2026
- `ARCHITECTURE.md` (`.planning/research/ARCHITECTURE.md`) — Directory layout, dependency injection patterns, data flow
- `PITFALLS.md` (`.planning/research/PITFALLS.md`) — XSS pipeline ordering, goose embedding, Postgres EBS pitfalls
- `01-CONTEXT.md` — Locked decisions D-01 through D-10
- [Go 1.22 Routing Enhancements](https://go.dev/blog/routing-enhancements) — Official Go blog

### Secondary (MEDIUM confidence)

- [GHA service containers docs](https://docs.github.com/en/actions/use-cases-and-examples/creating-an-example-workflow#adding-a-service-container) — health check pattern for Postgres service
- [golangci-lint-action v6](https://github.com/golangci/golangci-lint-action) — GHA action pinning behavior

### Tertiary (LOW confidence)

- goose v3 pgxpool adapter pattern — training data; verify against goose v3.27.0 docs before implementing

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — all versions verified in STACK.md against pkg.go.dev March 2026
- Architecture: HIGH — locked by D-01 through D-10 in CONTEXT.md; backed by ARCHITECTURE.md
- Pitfalls: HIGH — sourced from PITFALLS.md with specific CVE and issue references
- GHA CI pattern: MEDIUM — standard pattern; GHA service container health check verified against official docs
- goose pgxpool adapter: LOW — needs verification against goose v3.27.0 API docs before implementation

**Research date:** 2026-03-26
**Valid until:** 2026-04-26 (stable stack; goose/goldmark versions worth re-checking if delayed)
