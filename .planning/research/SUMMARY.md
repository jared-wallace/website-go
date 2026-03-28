# Project Research Summary

**Project:** website-go — personal blog server for jared-wallace.com
**Domain:** Server-side rendered Go blog with nautical design system (The Wild Meridian)
**Researched:** 2026-03-26 (backend platform), 2026-03-28 (v1.2 CSS polish update)
**Confidence:** HIGH

---

## Executive Summary

This is a well-understood domain: a single-author markdown blog with a web-based admin panel, deployed as a Dockerized Go binary behind Nginx and AWS ALB. The research consensus is strong across all four dimensions. Go 1.22+'s enhanced `net/http.ServeMux` eliminates the need for an external router; `html/template` with goldmark covers all rendering needs; pgx v5 + goose + scs v2 is the modern, minimal dependency stack for this exact deployment shape. The entire external dependency footprint is six packages (plus bluemonday, which is security-critical and must be added). Everything else is stdlib.

The recommended architecture is a three-layer monolith (`handler → service → repository`) wired via constructor injection in `main.go`. A single Go binary handles both the public blog and the admin panel. Static assets and SQL migrations are embedded at compile time via `//go:embed`. The binary is the only deployable artifact beyond the Postgres sidecar, which stores its data on an EBS bind-mount.

The v1.2 Shore Leave Polish milestone adds five CSS design areas to the existing site without any Go backend changes: background grain texture, card depth and shadow, page entry animations, dark mode color transitions, and a footer redesign. All five areas are achievable with pure CSS and minor template changes — zero new dependencies. The top risks are all preventable with a small number of upfront decisions: session cookie flags correctly set from day one, Postgres EBS directory chowned to UID 999 before first run, goldmark output flowing through bluemonday before being cast to `template.HTML`, ASG locked to `max_size = 1`, and dark mode transitions gated behind a deferred `.theme-ready` JS class to prevent page-load flash. None of these are exotic — they are documented, one-time setup steps.

---

## Key Findings

### Recommended Stack

Go 1.26.1 with a deliberate stdlib-first philosophy. The standard library covers routing (Go 1.22+ ServeMux with method + wildcard patterns), templating (html/template with context-aware escaping), RSS generation (encoding/xml), logging (log/slog), and static asset embedding (embed.FS). Six external dependencies complete the stack: goldmark (markdown), pgx v5 (Postgres driver + pool), goose (migrations), scs v2 (sessions), golang.org/x/crypto (bcrypt), and goldmark-meta. The v1.2 CSS polish layer adds zero new dependencies — all grain texture, animation, and transition work is pure CSS and inline SVG.

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
| bluemonday (add to go.mod) | HTML sanitization | goldmark output must be sanitized before `template.HTML` cast |

Full dependency list and install commands: see `.planning/research/STACK.md`.

**CSS design (v1.2 — zero new dependencies):**
- Background grain: inline SVG `feTurbulence` as data-URI on `body::before` pseudo-element
- Card depth: two-layer `box-shadow` with warm hue (`rgba(44, 36, 24, ...)`) — not pure black
- Entry animations: CSS `@keyframes` + `nth-child` stagger delays; no JavaScript
- Dark mode transitions: CSS `@property` on `:root` (Baseline 2024) gated behind `.theme-ready` JS class
- Footer: semantic `<footer>` with labeled `<nav aria-label="Footer navigation">`

### Expected Features

**Must have (table stakes) — backend platform:**
- Markdown rendering with syntax-highlighted code blocks
- Post listing with pagination and published/draft state
- Session-based admin auth (bcrypt + SCS)
- Web-based markdown editor with split-pane live preview
- Post CRUD with soft-delete, readable URL slugs
- RSS 2.0 feed
- Open Graph / social meta tags
- Mobile-responsive layout + custom 404 page

**Must have (table stakes) — v1.2 CSS polish:**
- Background grain texture on `body` — flat `#F5F0E8` reads as default web without it
- Card resting box-shadow — cards with no resting state have no context for the hover lift
- `prefers-reduced-motion` guard on all animations — WCAG 2.1 SC 2.3.3; non-negotiable
- Tag pill filled background — transparent pills are nearly invisible in dark mode (legibility fix)
- Footer About link + layout — current footer lacks navigation utility

**Should have (differentiators):**
- Thumbs-up reaction counter (no account required)
- API push endpoint (`POST /api/push`) for local `.md` workflow
- Table of contents auto-generation (goldmark AST walk)
- Card stagger animation + page fade-in (pure CSS, no JS)
- Dark mode color transitions (JS-gated post-load)
- Homepage hero heading above card grid
- Rope divider at full opacity (remove suppressing `opacity: 0.6`)

**Defer to v2+:**
- Clickable tag filtering (requires backend routing, out of scope for CSS milestone)
- LaTeX / KaTeX rendering
- Post series / related posts
- Image optimization / WebP pipeline

**Explicitly out of scope (anti-features):**
Comments, user registration, OAuth, full-text search, newsletter, in-app analytics, social media links in footer. See `.planning/research/FEATURES.md` for rationale on each.

### Architecture Approach

A three-layer monolith in a single Go binary. Handler packages translate HTTP to/from domain types; service packages hold business logic with no `net/http` imports; repository packages hold SQL with no business logic. Dependencies flow inward via constructor injection — no globals. The binary is deployed in a docker-compose stack alongside a Postgres sidecar, with EBS bind-mounted for both Postgres data and uploaded images. Static assets and migrations are embedded via `//go:embed`.

For v1.2 CSS changes, the architecture is explicit: one CSS file per audience (`main.css` for the public blog, `admin.css` for admin — no new files), new sections appended with labeled block comments, all color values use `var(--color-*)` tokens (never hardcoded hex). Template changes are isolated to `base.html` and `list.html` — `admin-base.html` is unaffected because the two template sets are parsed independently via separate `ParseFS` calls.

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

**Backend security (highest consequence):**

1. **Session cookie missing `HttpOnly + Secure + SameSite=Lax`** — Set all three flags explicitly. Go's `http.Cookie` zero-values all security fields. Regenerate session ID on login (session fixation). Add CSRF token to admin forms. Address in the auth phase.

2. **Goldmark output cast to `template.HTML` without sanitization** — Run `bluemonday.UGCPolicy()` between goldmark and `template.HTML`. Pipeline: `markdown → goldmark.Convert() → bluemonday.Sanitize() → template.HTML`. Write a test injecting `<script>alert(1)</script>`. Address in the markdown rendering phase.

3. **Image upload MIME spoofing / path traversal** — Validate magic bytes via `http.DetectContentType` (first 512 bytes). Generate a random server-side filename (UUID + extension). Never use the client-supplied filename. Enforce `http.MaxBytesReader` before multipart parse.

4. **Postgres EBS bind-mount wrong ownership** — `chown 999:999 && chmod 700` on host dir before first run. Document in `Makefile deploy` target.

5. **ASG `max_size > 1` with EBS-bound data** — Set `max_size = 1` in Terraform. Set `delete_on_termination = false` on the data volume. Schedule daily EBS snapshots.

**CSS polish (v1.2 — correctness-critical):**

6. **Dark mode transitions fire on page load and flash white (Pitfall 15)** — Gate all CSS color transitions behind `.theme-ready` class added via `window.addEventListener('load', ...)` in `main.js`. Without this, every dark-mode page load produces a visible white flash.

7. **Mobile scroll lag from grain texture (Pitfall 14)** — Never use `background-attachment: fixed`. Use `position: fixed; inset: 0` on `body::before` pseudo-element. iOS Safari and Android Chrome disable fixed-background compositing during scroll.

8. **`prefers-reduced-motion` not respected (Pitfall 17)** — Add nuclear override block at bottom of `main.css`: `animation-duration: 0.01ms !important` (not `animation: none`, which breaks `fill-mode: both` and leaves elements permanently invisible). Fix the existing `reaction-bounce` animation in the same pass.

---

## Implications for Roadmap

The project has two distinct implementation contexts: the backend platform (phases 1–6 below, substantially built per the existing codebase) and the v1.2 CSS polish milestone (phases 7–9 below). Both are covered by research.

### Phase 1: Foundation
**Rationale:** Everything downstream depends on project layout, DB connection, schema, and the markdown pipeline. Lay this correctly or pay the refactor tax on every subsequent phase.
**Delivers:** Compilable binary with DB connection, migrations embedded and running, markdown renderer with bluemonday sanitizer, project directory structure matching architecture spec.
**Addresses:** Markdown rendering (table stakes), syntax highlighting, project structure.
**Avoids:** Pitfall 9 (everything-in-main), Pitfall 2 (XSS via unsanitized markdown), Pitfall 6 (DB connection leaks — establish `defer rows.Close()` pattern and `sqlclosecheck` lint from the start).

### Phase 2: Public Blog (Read-Only)
**Rationale:** Public routes have no auth dependency and validate the full template/render/DB pipeline before the admin layer is added. Ship something readable first.
**Delivers:** Post listing with pagination, single post view, responsive layout, 404 handler, published/draft state filter, readable slugs, published dates.
**Avoids:** Pitfall 13 (draft post exposure — enforce `WHERE published = true` from day one), Pitfall 5 (set `http.Server` timeouts), Pitfall 11 (graceful shutdown).

### Phase 3: Admin Panel + Auth
**Rationale:** Auth is the prerequisite for all write operations. Session management and cookie security must be correct before any admin route is exposed.
**Delivers:** Session-based login, session middleware, admin post CRUD (create/edit/soft-delete), draft/publish workflow, split-pane web editor.
**Avoids:** Pitfall 1 (session cookie flags — `HttpOnly + Secure + SameSite=Lax`, session fixation, CSRF tokens), Pitfall 10 (bcrypt cost 12).
**Research flag:** If HTMX is chosen for live preview, verify HTMX + Go integration pattern before implementation. HTMX is not yet in the dependency list — a vanilla JS debounced fetch is a valid alternative with zero added dependencies.

### Phase 4: Distribution + Engagement
**Rationale:** With posts existing and admin working, distribution features can be added independently.
**Delivers:** RSS 2.0 feed, Open Graph / canonical URL meta tags, thumbs-up reaction counter, sitemap.
**Avoids:** Pitfall 13 (RSS draft exposure — same `published = true` filter).
**Note:** RSS hand-rolled with `encoding/xml` stdlib; do not reach for `gorilla/feeds`.

### Phase 5: API Push + Image Upload
**Rationale:** Both features touch security-sensitive surfaces and should be built after core flows are stable.
**Delivers:** `POST /api/push` with bearer token; image upload with magic-byte validation, random server-side filenames, EBS persistence.
**Avoids:** Pitfall 4 (image upload MIME spoofing and path traversal).

### Phase 6: Docker + Deployment
**Rationale:** Wrap a working binary in production-grade infrastructure. Doing this last means infrastructure is tested against real code.
**Delivers:** Multi-stage Dockerfile, docker-compose with Postgres sidecar, Makefile task runner, Nginx config, EBS bind-mount with correct ownership documented.
**Avoids:** Pitfall 3 (Postgres EBS ownership), Pitfall 8 (ASG `max_size=1`, `delete_on_termination=false`), Pitfall 12 (Nginx `proxy_set_header Host $host`), Pitfall 7 (X-Forwarded-For rightmost IP).

### Phase 7: CSS Foundation and Component Fixes (v1.2)
**Rationale:** All changes are CSS-only and additive to existing selectors. Highest impact-to-risk ratio in v1.2. No template changes — these can be verified independently in browser.
**Delivers:** Background grain texture, warm card two-layer shadows, filled tag pills, rope divider at full opacity, stale CSS comment fix ("The Log" → "The Wild Meridian").
**Avoids:** Pitfall 14 (mobile scroll lag — `position: fixed` pseudo-element, not `background-attachment: fixed`), Pitfall 19 (blend mode inversion in dark mode — separate `mix-blend-mode` per theme).

### Phase 8: Animation and Transitions (v1.2)
**Rationale:** Animations have CSS dependencies on Phase 7 (card shadows give animated cards a surface to lift from). Dark mode transitions require the one v1.2 change that touches both CSS and JavaScript.
**Delivers:** Card stagger animation (pure CSS `nth-child` delays), page fade-in, dark mode color transitions gated behind `.theme-ready` JS class, `prefers-reduced-motion` guard covering all animations including existing `reaction-bounce`.
**Avoids:** Pitfall 15 (load flash — `.theme-ready` deferred class), Pitfall 16 (CLS regression — CSS-only initial state, `fill-mode: both`, transform+opacity only), Pitfall 17 (reduced motion — `animation-duration: 0.01ms !important` nuclear override).
**Gate:** Lighthouse CLS ≤ 0.1 before and after; physical mobile device scroll test for grain texture.

### Phase 9: Template Changes (v1.2)
**Rationale:** CSS must be in place before HTML changes land. The About link must appear in the footer (this phase) before it is removed from the nav (also this phase) — do not strand navigation.
**Delivers:** Footer redesign with About link + personality tagline + two-column layout, homepage hero heading above card grid, About removed from primary nav.
**Avoids:** Pitfall 18 (ARIA landmark collisions — `aria-label` on both `<nav>` elements; footer stays direct child of `<body>`).
**Gate:** Axe DevTools scan after footer expansion; zero landmark violations required.

### Phase Ordering Rationale

- Foundation before everything: the `internal/` package structure, DB migration strategy, and markdown pipeline are load-bearing. Retrofitting costs compound.
- Public read before admin write: validates the render → template pipeline with the simpler code path before auth complexity is introduced.
- Auth before any write operation: session security pitfalls are the highest-impact security risk in this project.
- CSS foundation (Phase 7) before animations (Phase 8): card shadows give animated cards a surface to lift from; z-index stacking between the fixed grain pseudo-element and animated cards must be verified before animation layers are added.
- Dark mode transitions (Phase 8) before template changes (Phase 9): any new footer elements need to be covered by the `.theme-ready` transition system from day one.
- Template changes last in v1.2: CSS is already in effect when HTML classes are emitted; unused CSS has zero user impact, but missing CSS with present HTML causes unstyled flash.

### Research Flags

**Phases with well-documented patterns (no deeper research needed):**
- Phase 1 (Foundation): Standard Go project layout + pgx + goose — extensively documented.
- Phase 2 (Public Blog): Standard CRUD + html/template — no exotic patterns.
- Phase 3 (Admin + Auth): SCS v2 + bcrypt — Alex Edwards' documentation covers this exhaustively.
- Phase 4 (Distribution): RSS 2.0 spec is stable; OG tags are trivial.
- Phase 6 (Docker + Deployment): Standard patterns; pitfalls documented and preventable.
- Phase 7 (CSS Foundation): All techniques verified against MDN, CSS-Tricks, Frontend Masters. HIGH confidence.
- Phase 8 (Animations): `@keyframes` stagger and `prefers-reduced-motion` patterns MDN-documented; `.theme-ready` deferred class pattern verified against two dedicated blog posts.
- Phase 9 (Templates): Go template composition is standard; template isolation confirmed via direct codebase inspection; ARIA nav labeling is W3C spec.

**Phases that may benefit from targeted research during planning:**
- Phase 3 (web editor live preview): Decide HTMX vs. vanilla JS debounced fetch before implementation. HTMX is not yet in the dependency list.
- Phase 5 (API push): Decide on exact request format (multipart vs. JSON body with base64) and upsert strategy (slug-based `ON CONFLICT`) before implementation.

---

## Disagreements Between Research Files

One naming inconsistency exists between ARCHITECTURE.md and STACK.md. It does not affect behavior but must be resolved before implementation:

| Topic | STACK.md says | ARCHITECTURE.md says | Recommendation |
|---|---|---|---|
| DB layer package | Use native pgx/v5 interface; "do not use sqlx" | Repository docs reference `sqlx` in some comments | **Use native pgx/v5.** STACK.md is explicit and authoritative. |
| Migration tool | goose v3 | One comment references golang-migrate | **Use goose v3.** Both are valid; goose has better `//go:embed` support. |

---

## Confidence Assessment

| Area | Confidence | Notes |
|---|---|---|
| Stack | HIGH | All package versions verified against pkg.go.dev; Go version confirmed via go.dev/doc/devel/release; CSS techniques verified against MDN and CSS-Tricks |
| Features | HIGH | Established domain; v1.2 research based on direct inspection of existing codebase (`main.css` line numbers referenced throughout FEATURES.md) |
| Architecture | HIGH | Direct codebase inspection of all relevant files; template parsing pattern confirmed empirically; CSS file isolation confirmed |
| Pitfalls | HIGH | Backend pitfalls sourced from OWASP cheat sheets and Cloudflare engineering blog; CSS pitfalls from MDN spec and documented browser bugs (Mozilla Bugzilla #90198) |

**Overall confidence: HIGH**

### Gaps to Address

- **bluemonday not yet in go.mod:** PITFALLS.md identifies this as security-critical (Pitfall 2). Must be added before any markdown rendering ships publicly. Add `github.com/microcosm-cc/bluemonday` to the dependency list.
- **HTMX decision:** FEATURES.md mentions HTMX for live preview. STACK.md does not list it. Decide during Phase 3 planning whether HTMX (one JS file) is preferable to vanilla JS. Either path is valid; the decision should be explicit.
- **Alpine version at build time:** STACK.md notes alpine:3.21 as "assumed current — verify at build time." Confirm actual current version when the Dockerfile is written.
- **CSRF library:** PITFALLS.md notes Go 1.25+ includes `CrossOriginProtection` in stdlib. As of Go 1.26, confirm before reaching for `gorilla/csrf`.
- **Grain texture mobile performance:** The `position: fixed` pseudo-element pattern is the established safe approach, but STACK.md recommends testing on a physical low-end Android device before shipping. This is a verification gap, not a technique gap.

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
- [MDN — prefers-reduced-motion, @property, CSS transitions](https://developer.mozilla.org) — animation and transition specifications
- Direct codebase inspection: `web/static/main.css`, `web/templates/base.html`, `web/templates/admin-base.html`, `internal/handler/blog/handler.go`, `web/embed.go`

### Secondary (MEDIUM confidence)
- [Josh W. Comeau — Designing Beautiful Shadows](https://www.joshwcomeau.com/css/designing-shadows/) — layered shadow technique, warm hue matching
- [Frontend Masters — Grainy Gradients](https://frontendmasters.com/blog/grainy-gradients/) — feTurbulence performance notes, numOctaves guidance
- [CSS-Tricks — Staggered Animations](https://css-tricks.com/different-approaches-for-creating-a-staggered-animation/) — `--animation-order` custom property stagger
- [Blog of Maxime Heckel — Fixing dark mode flash](https://blog.maximeheckel.com/posts/switching-off-the-lights-part-2-fixing-dark-mode-flashing-on-servered-rendered-website/) — `.theme-ready` deferred class pattern
- [web.dev — Optimize CLS](https://web.dev/articles/optimize-cls) — safe vs unsafe animation properties
- [SCS Session Manager — Alex Edwards](https://www.alexedwards.net/blog/scs-session-manager)
- [pgx vs lib/pq — Preslav Rachev](https://preslav.me/2022/05/13/pq-or-pgx-choosing-the-right-postgresql-golang-driver/)
- [ASG with stateful Docker containers — Portworx](https://portworx.com/blog/auto-scaling-groups-ebs-docker/)
- [W3C WAI — Landmark Regions](https://www.w3.org/WAI/ARIA/apg/practices/landmark-regions/) — multiple nav labeling requirement

### Tertiary (source-verified, cross-check at implementation)
- [Mozilla Bugzilla #90198](https://bugzilla.mozilla.org/show_bug.cgi?id=90198) — fixed-background repaints on scroll (intentional browser behavior)
- [RSS 2.0 specification](https://validator.w3.org/feed/docs/rss2.html)

---

*Research completed: 2026-03-28 (updated from 2026-03-26 initial)*
*Ready for roadmap: yes*
