# Project Research Summary

**Project:** website-go — personal blog server for jared-wallace.com
**Domain:** Single-author personal blog, Go, server-side rendered, Dockerized, AWS EC2 + EBS
**Researched:** 2026-03-26
**Confidence:** HIGH

---

## Executive Summary

This is a well-understood domain: a single-author markdown blog with a web-based admin panel, deployed as a Dockerized Go binary behind Nginx and AWS ALB. The research consensus is strong across all four dimensions. Go 1.22+'s enhanced `net/http.ServeMux` eliminates the need for an external router; `html/template` with goldmark covers all rendering needs; pgx v5 + goose + scs v2 is the modern, minimal dependency stack for this exact deployment shape. The entire external dependency footprint is six packages. Everything else is stdlib.

The recommended architecture is a three-layer monolith (`handler → service → repository`) wired via constructor injection in `main.go`. A single Go binary handles both the public blog and the admin panel. Static assets and SQL migrations are embedded at compile time via `//go:embed`. The binary is the only deployable artifact beyond the Postgres sidecar, which stores its data on an EBS bind-mount.

The top risks are all preventable with a small number of upfront decisions: session cookie flags must be set correctly from day one (not retrofitted), the Postgres EBS directory must be chowned to UID 999 before first run, goldmark output must flow through a bluemonday sanitizer before being cast to `template.HTML`, and the ASG must be locked to `max_size = 1` before any production data is written. None of these are exotic — they are documented, one-time setup steps. The overall risk profile for this project is low.

---

## Key Findings

### Recommended Stack

Go 1.26.1 with a deliberate stdlib-first philosophy. The standard library covers routing (Go 1.22+ ServeMux with method + wildcard patterns), templating (html/template with context-aware escaping), RSS generation (encoding/xml), logging (log/slog), and static asset embedding (embed.FS). Six external dependencies complete the stack: goldmark (markdown), pgx v5 (Postgres driver + pool), goose (migrations), scs v2 (sessions), golang.org/x/crypto (bcrypt), and optionally goldmark-meta + goldmark-highlighting.

**Core technologies:**

| Technology | Purpose | Rationale |
|---|---|---|
| Go 1.26.1 | Language / runtime | Latest stable; Go 1.22 ServeMux eliminates router dependency |
| `net/http` ServeMux | Routing | Method + wildcard patterns sufficient for ~15 routes |
| `html/template` | Server-side rendering | Context-aware escaping; XSS protection by default |
| goldmark v1.8.2 | Markdown → HTML | CommonMark-compliant; extensible; used by Hugo |
| pgx/v5 v5.9.1 | Postgres driver + pool | Replaces maintenance-only lib/pq; native pgxpool included |
| goose v3.27.0 | Schema migrations | Embedded SQL files via `//go:embed`; clean Go integration |
| scs v2.9.0 | Session management | Server-side tokens; pgxstore sub-package; OWASP-aligned |
| `golang.org/x/crypto/bcrypt` | Password hashing | Cost 12 minimum; constant-time comparison built in |

Full dependency list and install commands: see `.planning/research/STACK.md`.

### Expected Features

The domain is well-established. Community consensus on what a personal dev blog must have is clear and consistent across sources.

**Must have (table stakes):**
- Markdown rendering with syntax-highlighted code blocks — the core contract
- Post listing with pagination and published/draft state
- Readable URL slugs and published dates
- Session-based admin auth (bcrypt + SCS)
- Web-based markdown editor with split-pane live preview
- Post CRUD with soft-delete
- RSS 2.0 feed
- Open Graph / social meta tags
- Mobile-responsive layout + custom 404 page

**Should have (differentiators):**
- Thumbs-up reaction counter (no account required, low moderation overhead)
- API push endpoint (`POST /api/push`) for local `.md` file workflow
- Table of contents auto-generation (goldmark AST walk)
- Estimated reading time (word count / 200 wpm)
- Dark mode (CSS `prefers-color-scheme` — zero JS required)
- Sitemap (`/sitemap.xml`)
- Canonical URL tag

**Defer to v2+:**
- LaTeX / KaTeX rendering (high complexity; front-matter-gated client-side load)
- Post series / related posts
- Tag listing pages (store `tags` field from day one; render later)
- Image optimization / WebP pipeline

**Explicitly out of scope (anti-features):**
Comments, user registration, OAuth, full-text search, newsletter, in-app analytics. See `.planning/research/FEATURES.md` for rationale on each.

### Architecture Approach

A three-layer monolith in a single Go binary. Handler packages translate HTTP to/from domain types; service packages hold business logic with no `net/http` imports; repository packages hold SQL with no business logic. Dependencies flow inward via constructor injection — no globals. The binary is deployed in a docker-compose stack alongside a Postgres sidecar, with EBS bind-mounted for both Postgres data and uploaded images. Static assets and migrations are embedded via `//go:embed`.

**Major components:**

| Component | Responsibility |
|---|---|
| `cmd/server/main.go` | DI wiring only; starts HTTP server; target ~50 lines |
| `internal/config` | Env-var loading at startup |
| `internal/handler/{blog,admin,api,feed}` | HTTP boundary; calls service layer; no SQL |
| `internal/middleware` | Auth guards, session loading, request logger, recovery |
| `internal/service/{post,auth}` | Business logic; markdown parsing; slug generation |
| `internal/repository/{post,session}` | SQL queries only; no `net/http` |
| `internal/markdown` | Goldmark + bluemonday pipeline; pure function |
| `db/migrations/` | Embedded SQL files; goose runs at startup |
| `web/templates/` + `web/static/` | Embedded into binary at compile time |

### Critical Pitfalls

1. **Session cookie missing `HttpOnly + Secure + SameSite=Lax`** — Set all three flags explicitly. Go's `http.Cookie` zero-values all security fields. Regenerate session ID on login (session fixation). Add CSRF token to admin forms. Address in the auth phase.

2. **Goldmark output cast to `template.HTML` without sanitization** — Run `bluemonday.UGCPolicy()` between goldmark and `template.HTML`. The pipeline is always: `markdown → goldmark.Convert() → bluemonday.Sanitize() → template.HTML`. Write a test that injects `<script>alert(1)</script>`. Address in the markdown rendering phase.

3. **Postgres EBS bind-mount wrong ownership** — The Postgres container runs as UID/GID 999. The host directory must be `chown 999:999` and `chmod 700` before first run. Document in `Makefile deploy` target. Address before any data is written to production.

4. **Image upload MIME spoofing / path traversal** — Validate magic bytes via `http.DetectContentType` (first 512 bytes). Generate a random server-side filename (UUID + extension). Never use the client-supplied filename. Enforce `http.MaxBytesReader` before multipart parse. Address before image upload ships.

5. **ASG `max_size > 1` with EBS-bound data** — EBS volumes are instance-bound. Set `max_size = 1` in Terraform. Set `delete_on_termination = false` on the data volume. Schedule daily EBS snapshots. Address before first production deployment.

---

## Implications for Roadmap

Based on the dependency graph in FEATURES.md, the build order in ARCHITECTURE.md, and the phase warnings in PITFALLS.md, the following phase structure is recommended. It is opinionated.

### Phase 1: Foundation
**Rationale:** Everything downstream depends on project layout, DB connection, schema, and the markdown pipeline. Lay this correctly or pay the refactor tax on every subsequent phase.
**Delivers:** Compilable binary with DB connection, migrations embedded and running, markdown renderer with bluemonday sanitizer, project directory structure matching architecture spec.
**Addresses:** Markdown rendering (table stakes), syntax highlighting, project structure.
**Avoids:** Pitfall 9 (everything-in-main), Pitfall 2 (XSS via unsanitized markdown), Pitfall 6 (DB connection leaks — establish `defer rows.Close()` pattern and `sqlclosecheck` lint from the start).
**Note:** Pre-render markdown on write and store rendered HTML in DB — do not render on every request. This is the correct anti-pattern avoidance per ARCHITECTURE.md.

### Phase 2: Public Blog (Read-Only)
**Rationale:** Public routes have no auth dependency and validate the full template/render/DB pipeline before the admin layer is added. Ship something readable first.
**Delivers:** Post listing with pagination, single post view, responsive layout, 404 handler, published/draft state filter, readable slugs, published dates.
**Addresses:** All table-stakes reading features.
**Avoids:** Pitfall 13 (draft post exposure on public routes — enforce `WHERE published = true` from day one), Pitfall 5 (set `http.Server` timeouts and graceful shutdown in this phase), Pitfall 11 (graceful shutdown paired with Pitfall 5).

### Phase 3: Admin Panel + Auth
**Rationale:** Auth is the prerequisite for all write operations. Session management and cookie security must be correct before any admin route is exposed.
**Delivers:** Session-based login, session middleware, admin post CRUD (create/edit/soft-delete), draft/publish workflow, split-pane web editor with live markdown preview.
**Addresses:** Session auth, post CRUD, draft/published workflow, web editor — all table stakes.
**Avoids:** Pitfall 1 (session cookie flags — `HttpOnly + Secure + SameSite=Lax`, session fixation prevention, CSRF tokens on all admin forms), Pitfall 10 (bcrypt cost 12).
**Implementation note:** HTMX is a natural fit for the split-pane preview via a `/preview` endpoint. However, HTMX is optional — a simple form + full-page reload preview also works and adds zero JS dependencies.

### Phase 4: Distribution + Engagement
**Rationale:** With posts existing and admin working, distribution features can be added independently. These are small, well-scoped additions.
**Delivers:** RSS 2.0 feed, Open Graph / canonical URL meta tags, thumbs-up reaction counter, sitemap.
**Addresses:** RSS, OG tags, thumbs-up, sitemap — all either table stakes or high-value differentiators.
**Avoids:** Pitfall 13 (RSS draft exposure — same `published = true` filter applies here).
**Note:** RSS can be hand-rolled with `encoding/xml` stdlib per STACK.md. Do not reach for `gorilla/feeds`.

### Phase 5: API Push + Image Upload
**Rationale:** Both features touch security-sensitive surfaces (bearer token auth, file upload) and should be validated carefully. Deferred until the core blog is stable.
**Delivers:** `POST /api/push` endpoint with bearer token for local `.md` workflow; image upload with magic-byte validation, random server-side filenames, and EBS persistence.
**Addresses:** API push differentiator, image support.
**Avoids:** Pitfall 4 (image upload MIME spoofing and path traversal — magic bytes, random filenames, `MaxBytesReader`).

### Phase 6: Docker + Deployment
**Rationale:** Wrap a working binary in production-grade infrastructure. Doing this last means the infrastructure is tested against real code, not a prototype.
**Delivers:** Multi-stage Dockerfile (golang:1.26-alpine build, alpine:3.21 runtime), docker-compose with Postgres sidecar, Makefile task runner, Nginx config with correct `proxy_set_header` directives, EBS bind-mount with correct ownership documented.
**Avoids:** Pitfall 3 (Postgres EBS ownership — `chown 999:999`, documented in Makefile), Pitfall 8 (ASG `max_size=1`, `delete_on_termination=false`, snapshot schedule), Pitfall 12 (Nginx `proxy_set_header Host $host`), Pitfall 7 (X-Forwarded-For — read rightmost IP).

### Phase 7: Enhancements (stretch / post-v1)
**Rationale:** Low priority but real value. Add after the blog is live and content is accumulating.
**Delivers:** Estimated reading time, table of contents, dark mode (CSS-only), post cover images, tag listing pages, LaTeX/KaTeX (front-matter-gated).

### Phase Ordering Rationale

- Foundation before everything: the `internal/` package structure, DB migration strategy, and markdown pipeline are load-bearing. Retrofitting these is expensive.
- Public read before admin write: validates the render → template pipeline end-to-end with the simpler code path before auth complexity is introduced.
- Auth before any write operation: session security pitfalls are the highest-impact security risk in this project.
- Distribution features after content pipeline: RSS and OG tags depend on published posts and slugs existing.
- Image upload deferred to Phase 5: it has the most security surface area (MIME validation, path traversal) and should be built after core flows are stable.
- Docker last: deploying a working binary, not a moving target.

### Research Flags

**Phases with well-documented patterns (no deeper research needed):**
- Phase 1 (Foundation): Standard Go project layout + pgx + goose — extensively documented.
- Phase 2 (Public Blog): Standard CRUD + html/template — no exotic patterns.
- Phase 3 (Admin + Auth): SCS v2 + bcrypt — Alex Edwards' documentation covers this exhaustively.
- Phase 4 (Distribution): RSS 2.0 spec is stable; OG tags are trivial.
- Phase 6 (Docker + Deployment): Standard patterns; pitfalls documented and preventable.

**Phases that may benefit from targeted research during planning:**
- Phase 3 (web editor): If HTMX is chosen for the live preview endpoint, verify the HTMX + Go integration pattern. HTMX is not in the dependency list yet — confirm whether a minimal vanilla JS solution is preferable to avoid another JS dependency.
- Phase 5 (API push): Decide on the exact request format (multipart vs. JSON body with base64) and upsert strategy (slug-based `ON CONFLICT`) before implementation.

---

## Disagreements Between Research Files

One naming inconsistency exists between ARCHITECTURE.md and STACK.md. It does not affect behavior but must be resolved before implementation:

| Topic | STACK.md says | ARCHITECTURE.md says | Recommendation |
|---|---|---|---|
| DB layer package | Use native pgx/v5 interface; "do not use sqlx" | Repository docs reference `sqlx` in comments (e.g., `repository/post` note says "pgx or sqlx") | **Use native pgx/v5.** STACK.md is explicit and authoritative. The architecture's repository pattern works identically with pgx. The `sqlx` references in ARCHITECTURE.md are vestigial. |
| Migration tool | goose v3 | ARCHITECTURE.md says "golang-migrate" in one comment | **Use goose v3.** Both tools are valid; goose has better `//go:embed` support (per STACK.md). ARCHITECTURE.md's reference to golang-migrate is inconsistent with STACK.md. |

---

## Confidence Assessment

| Area | Confidence | Notes |
|---|---|---|
| Stack | HIGH | All package versions verified against pkg.go.dev; Go version confirmed via go.dev/doc/devel/release |
| Features | HIGH | Established domain; strong community consensus; feature list matches multiple reference blogs |
| Architecture | HIGH | Well-documented patterns; three-layer monolith for a Go blog is idiomatic and well-sourced |
| Pitfalls | HIGH | All pitfalls sourced from OWASP, official Go documentation, and production incident reports |

**Overall confidence: HIGH**

### Gaps to Address

- **HTMX decision:** FEATURES.md mentions HTMX for the live preview editor. STACK.md does not list HTMX as a dependency. This is a deliberate open question — decide during Phase 3 planning whether HTMX (one JS file) is worth it vs. a vanilla JS debounced fetch. Either path is valid; the decision should be explicit.
- **bluemonday dependency:** PITFALLS.md recommends `github.com/microcosm-cc/bluemonday` for HTML sanitization. STACK.md does not list it. It should be added to the dependency list before Phase 1 implementation starts. This is a security-critical addition, not optional.
- **Alpine version:** STACK.md notes alpine:3.21 as "assumed current — verify at build time." Confirm actual current version when the Dockerfile is written in Phase 6.
- **CSRF library:** PITFALLS.md notes Go 1.25+ includes `CrossOriginProtection` in stdlib; as of Go 1.26 this should be confirmed before reaching for `gorilla/csrf`.

---

## Sources

### Primary (HIGH confidence)
- [Go Release History](https://go.dev/doc/devel/release) — Go 1.26.1 confirmed latest stable
- [Routing Enhancements for Go 1.22](https://go.dev/blog/routing-enhancements) — ServeMux method + wildcard patterns
- [goldmark on pkg.go.dev](https://pkg.go.dev/github.com/yuin/goldmark) — v1.8.2
- [pgx on pkg.go.dev](https://pkg.go.dev/github.com/jackc/pgx/v5) — v5.9.1
- [scs on pkg.go.dev](https://pkg.go.dev/github.com/alexedwards/scs/v2) — v2.9.0
- [goose on pkg.go.dev](https://pkg.go.dev/github.com/pressly/goose/v3) — v3.27.0
- [golang.org/x/crypto/bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) — v0.49.0
- [bluemonday on pkg.go.dev](https://pkg.go.dev/github.com/microcosm-cc/bluemonday) — HTML sanitizer
- [OWASP Session Management Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html)
- [OWASP File Upload Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/File_Upload_Cheat_Sheet.html)
- [The complete guide to Go net/http timeouts — Cloudflare](https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/)
- [RSS 2.0 specification](https://validator.w3.org/feed/docs/rss2.html)

### Secondary (MEDIUM confidence)
- [How I developed a markdown blog in Go and HTMX — fluxsec.red](https://fluxsec.red/how-I-developed-a-markdown-blog-with-go-and-HTMX)
- [Building a Markdown Blog Engine with Go and Fiber — dasroot.net (2026)](https://dasroot.net/posts/2026/01/building-markdown-blog-engine-go-fiber/)
- [Go Project Structure: Practices & Patterns — glukhov.org (2025)](https://www.glukhov.org/post/2025/12/go-project-structure/)
- [The Repository Pattern in Go — Three Dots Labs](https://threedots.tech/post/repository-pattern-in-go/)
- [SCS Session Manager — Alex Edwards](https://www.alexedwards.net/blog/scs-session-manager)
- [Which Go Router Should I Use? — Alex Edwards](https://www.alexedwards.net/blog/which-go-router-should-i-use)
- [pgx vs lib/pq — Preslav Rachev](https://preslav.me/2022/05/13/pq-or-pgx-choosing-the-right-postgresql-golang-driver/)
- [ASG with stateful Docker containers — Portworx](https://portworx.com/blog/auto-scaling-groups-ebs-docker/)
- [Docker Volumes in Production — blog.shukebeta.com](https://blog.shukebeta.com/2024/10/23/docker-volumes-in-production-a-practical-guide-to-named-volumes-vs-bind-mounts/)

---

*Research completed: 2026-03-26*
*Ready for roadmap: yes*
