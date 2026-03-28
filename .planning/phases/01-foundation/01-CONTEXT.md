# Phase 1: Foundation - Context

**Gathered:** 2026-03-26
**Status:** Ready for planning

<domain>
## Phase Boundary

Deliver a compilable Go project scaffold with standard layout, a live Postgres connection with embedded goose migrations, a goldmark + bluemonday markdown rendering pipeline, an organized Makefile with dev workflow, and a passing GHA CI pipeline. No HTTP routes serving content yet — that's Phase 2.

</domain>

<decisions>
## Implementation Decisions

### Project Layout
- **D-01:** Flat `internal/` structure — `internal/server/`, `internal/handler/`, `internal/markdown/`, `internal/database/`, `internal/model/`. Add nesting only when a package actually needs sub-packages.
- **D-02:** Templates and static assets embedded via `go:embed` in `web/` package. Single binary deployment, no runtime file dependencies.
- **D-03:** Configuration via environment variables only (`os.Getenv` with helpers). No config files, no viper. 12-factor style with `envOr()` defaults and `mustEnv()` for required values.

### Markdown Pipeline
- **D-04:** All four goldmark extensions enabled: GFM tables + strikethrough, linkify, syntax highlighting (chroma), and YAML front matter (goldmark-meta).
- **D-05:** bluemonday `UGCPolicy()` for HTML sanitization. Allows standard formatting (bold, italic, links, images, code, tables, blockquotes) while stripping scripts, iframes, event handlers.

### Makefile & Dev Workflow
- **D-06:** Hot reload via `cosmtrek/air` for `make dev`. Watches `.go`, `.html`, `.css` files and auto-rebuilds. `.air.toml` committed to repo.
- **D-07:** Local Postgres via docker-compose only (`docker-compose.dev.yml`). `make dev-up` starts the container. Matches production topology.
- **D-08:** `golangci-lint` for `make lint` with `.golangci.yml` config. Enables errcheck, govet, staticcheck, gosimple, ineffassign, unused.

### GHA CI
- **D-09:** Postgres service container in CI. Tests run against a real database — catches migration and query issues.
- **D-10:** Single CI job running lint -> test -> build sequentially. Conserves free-tier minutes. Go module caching via `actions/setup-go` with `cache: true`.

### Claude's Discretion
- Specific Makefile target names and help text formatting
- `.golangci.yml` exact linter set (beyond the base six discussed)
- `internal/config/` package shape (struct fields, validation approach)
- Test file organization and naming conventions
- goose migration numbering format

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Tech Stack
- `.planning/research/STACK.md` — Authoritative dependency versions and rationale (pgx v5.9.1, goose v3.27.0, goldmark v1.8.2, scs v2.9.0, bcrypt v0.49.0)
- `.planning/research/ARCHITECTURE.md` — Project structure guidance and package organization patterns

### Project Context
- `.planning/PROJECT.md` — Core value, constraints, key decisions, infrastructure context
- `.planning/REQUIREMENTS.md` — FOUND-01 through FOUND-07 acceptance criteria
- `.planning/ROADMAP.md` — Phase 1 success criteria (5 criteria that must be TRUE)

### Infrastructure
- `.planning/research/PITFALLS.md` — Known pitfalls for Postgres EBS bind-mount, Docker multi-stage builds, goose embedding

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- None — greenfield project. Only `CLAUDE.md` exists in the repo.

### Established Patterns
- None yet. Phase 1 establishes the patterns all subsequent phases will follow.

### Integration Points
- `cmd/server/main.go` will be the entry point — wires config, database, and (later) HTTP server
- `db/migrations/` will hold goose SQL files embedded via `go:embed`
- `web/templates/` and `web/static/` will hold assets embedded via `go:embed`

</code_context>

<specifics>
## Specific Ideas

No specific requirements — open to standard approaches. The user consistently chose recommended/standard options, indicating a preference for conventional Go patterns over novel approaches.

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.

</deferred>

---

*Phase: 01-foundation*
*Context gathered: 2026-03-26*
