# Architecture Patterns

**Domain:** Personal blog platform (Go, server-side rendered, single-admin)
**Researched:** 2026-03-28 (updated for v1.2 Shore Leave Polish milestone)
**Overall confidence:** HIGH — all findings based on direct inspection of the existing codebase

---

## Recommended Architecture

A three-layer monolith served from a single Go binary, behind Nginx + AWS ALB. The binary handles both the public blog (jared-wallace.com) and the admin panel (admin.jared-wallace.com) via host-based routing inside the same process. One docker-compose stack: the Go app + a Postgres sidecar.

```
[Browser]
    |
    v
[AWS ALB] -- TLS termination (ACM cert)
    |
    v
[Nginx on EC2 :443/:80] -- reverse proxy
    |
    v
[Go app on :8080]
    |
    |-- /              --> Public blog handlers
    |-- /admin/...     --> Admin handlers (session-gated)
    |-- /api/push      --> API push handler (token-gated)
    |-- /feed.xml      --> RSS handler
    |-- /uploads/...   --> Static file server (EBS-mounted images)
    |-- /static/...    --> Embedded static assets (CSS, JS, fonts)
    |
    v
[Postgres sidecar :5432]   (data on EBS bind-mount at /var/www/html/pgdata)
```

---

## v1.2 Integration Architecture: CSS and Template Changes

This section is the primary output for the Shore Leave Polish milestone. It answers the four integration questions directly.

### Question 1: Should texture/animation CSS go in main.css or separate files?

**Answer: Keep everything in `main.css`. Do not split.**

The static assets are embedded into the Go binary at build time via `//go:embed static` in `web/embed.go`. Each additional CSS file requires an additional `<link>` tag in `base.html` and in `admin-base.html`. For this milestone's scope — texture, animations, transitions, and a handful of component tweaks — the complexity of managing multiple files outweighs any organizational benefit.

`main.css` is currently 645 lines. Adding the v1.2 changes will push it to approximately 750-800 lines, which remains maintainable. CSS custom properties (`--color-bg`, `--color-accent`, etc.) are already defined in `:root` and `[data-theme="dark"]` blocks at the top of the file. New animation keyframes and texture declarations should be appended to the end of `main.css` in clearly labeled sections.

**Practical rule:** One file per audience. `main.css` for the public blog. `admin.css` for the admin panel. No new files.

**Section ordering for appended CSS:**

```
/* === EXISTING SECTIONS (lines 1-645) === */
/* ... */

/* --- Texture / Grain --- */
/* background noise overlay, body::before pseudo-element */

/* --- Transitions --- */
/* dark mode color transitions on :root, body, .post-card, etc. */

/* --- Animations --- */
/* @keyframes definitions + .page-enter, .post-card-animate classes */

/* --- Footer (expanded) --- */
/* replaces/extends existing .site-footer block at line 563 */

/* --- Homepage Hero --- */
/* .list-hero, .list-heading, etc. */
```

If the `.site-footer` block at line 563 needs significant structural changes, edit it in place rather than appending a duplicate block. Duplicate selectors in CSS cascade silently — the last one wins, which creates confusion.

---

### Question 2: How to structure footer template changes without breaking admin-base.html?

**Answer: The footer lives only in `base.html`. `admin-base.html` has no footer. They do not share a footer block. Changes to `base.html`'s footer cannot break `admin-base.html`.**

Verified by direct inspection:

- `web/templates/base.html` — contains `<footer class="site-footer">` (lines 50-64). This is the public blog footer.
- `web/templates/admin-base.html` — has NO footer element. It ends with the `{{block "content"}}` and a `<script>` tag.

Both templates define `{{define "base"}}` but they are loaded into separate template sets:

```
blog handler:  ParseFS(web.Templates, "templates/base.html",       "templates/"+page)
admin handler: ParseFS(web.Templates, "templates/admin-base.html", "templates/"+page)
```

Because each page gets its own independent `*template.Template` set (not a shared global registry), the two `{{define "base"}}` definitions never collide. The blog handler's template set contains `base.html`; the admin handler's set contains `admin-base.html`. They are entirely separate objects in memory.

**What this means for the footer expansion:**

Edit `base.html`'s `<footer>` block freely. Add nav links, personality copy, About link. The admin panel is unaffected in every sense — different template file, different parsed template set, no shared state.

**One caution:** The `.Year` data binding in the footer (`{{.Year}}`) is injected by the blog handler's `render()` function. Any new template variables added to the footer (e.g., a `{{.SiteTagline}}`) must also be injected in `render()`. Keep the footer to static HTML + existing bindings where possible to avoid a Go-side change.

---

### Question 3: What is the right order to implement these changes to avoid conflicts?

**Answer: CSS custom properties first, then structural/layout CSS, then animations, then template HTML changes.**

The CSS cascade has real ordering dependencies here. CSS custom property changes affect every selector that references them, so those go first. Animations reference `transform` and `opacity` which should not conflict with anything, but they should come after the structural rules they animate.

**Recommended implementation sequence:**

```
Step 1 — CSS comment fix (no risk, do first)
  main.css line 3: "The Log" → "The Wild Meridian"

Step 2 — CSS custom properties: add transition values to :root
  Add --transition-color: 200ms ease to :root
  Add transition: color, background-color, border-color to body and key components
  Touching :root affects everything — verify dark mode toggle still works

Step 3 — Reaction button radius fix
  .reaction-btn: border-radius: 2rem → border-radius: 4px
  Isolated single-property change, no cascade effects

Step 4 — Post card depth and tag visibility
  .post-card: strengthen box-shadow baseline
  .tag-pill: adjust border, background, or color for visibility
  These are additive — adding properties to existing selectors

Step 5 — Rope divider visual strengthening
  .rope-divider: style changes (border style, color, opacity)
  Used in footer and post pages — verify both contexts look correct

Step 6 — Background texture (new CSS + pseudo-element)
  body::before with SVG noise or CSS gradient overlay
  Add texture CSS section at end of main.css
  Test with dark mode — texture must work in both themes

Step 7 — Page entry animations (new keyframes + classes)
  Add @keyframes at end of main.css
  Add .post-card animation class
  Animation classes are applied via JS or HTML class attributes — do not affect
  existing layout until the classes are actually added to templates

Step 8 — Template: footer expansion (base.html)
  Expand <footer> block in base.html
  Add About link, nav structure, personality copy
  Add corresponding CSS for new footer structure

Step 9 — Template: homepage heading/hero (list.html)
  Add hero/heading markup above .card-grid in list.html
  Add .list-hero CSS to support it
  If the heading needs a data binding, update blog handler's render() call

Step 10 — Template: remove About from nav (base.html)
  Remove <a href="/about" class="nav-link">About</a> from <nav> in base.html
  This step depends on Step 8 (About must exist in footer before removal from nav)
```

**Why this order matters:**

- Steps 2-5 are CSS-only, no template changes. They can be verified independently in the browser without touching any `.html` files.
- Step 6 (texture) interacts with `--color-bg` — do it after dark mode transitions are confirmed working.
- Step 7 (animations) adds CSS classes; those classes produce no visible effect until the HTML templates actually emit those classes. This decoupling is safe — adding unused CSS classes to `main.css` has zero user-visible impact.
- Steps 8-10 change HTML. Do them last so CSS is already in place when the markup lands.
- Step 10 (remove About from nav) must come after Step 8 (add About to footer). Don't strand the About link.

---

### Question 4: Template inheritance considerations for the homepage heading addition?

**Answer: There is no template inheritance to worry about. The heading change is isolated to `list.html`.**

Go's `html/template` does not have classical inheritance — it uses block/define composition. Each page is a completely independent `*template.Template` parsed from exactly two files: `base.html` + the page file. The `{{block "content"}}` in `base.html` is overridden by `{{define "content"}}` in the page template.

The homepage listing page is `web/templates/list.html`. The heading addition goes inside `{{define "content"}}` in that file, above the `.card-grid` div. No other template is affected.

**The only Go-side consideration:** If the homepage heading needs dynamic data (e.g., a tagline pulled from config), it must be passed via the `data` map in the blog handler's `render()` call for `list.html`. If it is static HTML text, no Go changes are needed — just add the markup.

**Current `list.html` content structure:**

```
{{define "content"}}
<div class="container">
  {{if not .Posts}}
    <div class="empty-state">...</div>
  {{else}}
    <div class="card-grid">...</div>
    <nav class="pagination">...</nav>
  {{end}}
</div>
{{end}}
```

**Proposed structure after heading addition:**

```
{{define "content"}}
<div class="container">
  <header class="list-hero">
    <h1 class="list-heading">dispatches from the deep end</h1>
  </header>
  {{if not .Posts}}
    ...
```

The heading sits inside `.container`, above the conditional post grid. Static text — no handler changes required.

---

## File Change Map

Every file that changes in v1.2, with the nature of the change:

| File | Change Type | What Changes |
|------|-------------|--------------|
| `web/static/main.css` | Edit in-place + append | Fix "The Log" comment; reaction-btn radius; card depth; tag pills; rope divider; dark mode transitions; texture section; animation section; footer expansion CSS; hero/heading CSS |
| `web/templates/base.html` | Edit | Expand `<footer>` block; remove About `<a>` from `<nav>` |
| `web/templates/list.html` | Edit | Add hero/heading header above `.card-grid` |
| `web/static/admin.css` | No change | Admin panel is isolated from all v1.2 work |
| `web/templates/admin-base.html` | No change | No footer, no public nav; unaffected |
| `internal/handler/blog/handler.go` | No change (unless dynamic data needed in hero) | Only if new `.HeroTagline` or similar binding is required |

---

## Component Boundaries

| Component | Responsibility | Communicates With |
|-----------|---------------|-------------------|
| `cmd/server/main.go` | Process entry point, wires dependencies, starts HTTP server | All internal packages |
| `internal/config` | Loads config from env vars or file at startup | `main.go` only |
| `internal/router` | Registers all routes, applies middleware chains | Handlers, middleware |
| `internal/middleware` | Auth guard, session loading, request logger, recovery | Router, sessions |
| `internal/handler/blog` | Public post list, single post, thumbs-up | PostService |
| `internal/handler/admin` | Login, editor, post CRUD, image upload | PostService, SessionStore |
| `internal/handler/api` | `POST /api/push` for local-to-server .md push | PostService |
| `internal/handler/feed` | Generates RSS XML | PostService |
| `internal/service/post` | Business logic: parse markdown, frontmatter, slugs | PostRepository |
| `internal/service/auth` | bcrypt verify, session create/destroy | SessionStore, DB |
| `internal/repository/post` | SQL queries: posts table CRUD | pgx pool |
| `internal/repository/session` | Session store backed by Postgres table | pgx pool |
| `internal/markdown` | Goldmark conversion, sanitization, syntax highlight | PostService |
| `db/migrations/` | SQL migration files (embedded via `//go:embed`) | goose at startup |
| `web/templates/` | html/template files for blog and admin views | Handlers |
| `web/static/` | CSS, JS, fonts (embedded into binary via `//go:embed`) | Router |
| `docker-compose.yml` | Defines `app` + `db` services, bind-mounts EBS path | EC2 host |

---

## Data Flow

### Public: Reader Loads a Post

```
Browser GET /posts/my-slug
  -> Router -> blog.Handler.ShowPost
  -> PostService.GetBySlug(slug)
    -> PostRepository.FindBySlug(slug)   [SQL SELECT]
    <- Post{Body: "# Hello..."}
  -> markdown.Render(post.Body)          [Goldmark HTML]
  -> template.Execute(w, data)           [html/template]
Browser receives rendered HTML page
```

### Admin: Publish a Post via Web Editor

```
Admin POST /admin/posts (form submit)
  -> middleware.RequireSession            [reject if no valid session]
  -> admin.Handler.CreatePost
  -> PostService.Create(form data)
    -> markdown.ParseFrontmatter(body)   [extract title, date, tags]
    -> PostRepository.Insert(post)       [SQL INSERT]
  -> redirect to /admin/posts
```

### API: Push Markdown from Local Machine

```
curl -X POST /api/push -H "Authorization: Bearer $TOKEN" --data-binary @post.md
  -> middleware.RequireAPIToken           [static token from config]
  -> api.Handler.Push
  -> PostService.Upsert(markdown bytes)
    -> markdown.ParseFrontmatter(body)
    -> PostRepository.Upsert(post)       [INSERT ... ON CONFLICT DO UPDATE]
  -> 200 OK
```

---

## Patterns to Follow

### Pattern: Append New CSS Sections With Block Comments

New CSS sections for v1.2 should be appended to the end of `main.css` in clearly labeled comment blocks. Do not interleave new selectors into existing sections unless editing an existing selector in-place.

```css
/* --- Texture / Grain --- */

body::before {
  content: '';
  position: fixed;
  inset: 0;
  /* noise overlay */
  pointer-events: none;
  z-index: -1;
}
```

### Pattern: Per-Page Template Sets (Existing — Do Not Break)

Each page is parsed as an independent `*template.Template` from exactly two files. This prevents `{{define}}` block name collisions across pages. Do not change the parsing logic in `handler.go`.

```go
// blog/handler.go — this pattern is correct and must not change
tmpl := template.New("").Funcs(funcMap).ParseFS(
    web.Templates,
    "templates/base.html",
    "templates/"+page,
)
```

### Pattern: CSS Custom Properties for Theme-Aware Design

All color values must use `var(--color-*)` tokens, not hardcoded hex. New v1.2 selectors — texture overlays, animated borders, hero headings — must all reference custom properties so dark mode works automatically.

```css
/* Correct */
.list-hero {
  border-bottom: 1px solid var(--color-divider);
  color: var(--color-text-muted);
}

/* Wrong — breaks dark mode */
.list-hero {
  border-bottom: 1px solid #C8B89A;
  color: #7A6A55;
}
```

### Pattern: Dependency Injection via Constructor (not global state)

Wire dependencies at startup in `main.go`. Pass them down. No `init()` globals.

---

## Anti-Patterns to Avoid

### Anti-Pattern: Duplicate CSS Selectors for "Override"

**What:** Appending a second `.site-footer { ... }` block to override the existing one.
**Why bad:** Silent cascade — both blocks apply, order-dependent. Hard to reason about which properties are active. When `main.css` reaches 1000+ lines this becomes a maintenance trap.
**Instead:** Edit the existing `.site-footer` block in-place for structural property changes. Append only genuinely new selectors (`.footer-nav`, `.footer-col`, etc.).

### Anti-Pattern: Adding a Third Stylesheet for Animations

**What:** Creating `web/static/animations.css` and linking it from `base.html`.
**Why bad:** Requires a new `<link>` tag in both `base.html` AND `admin-base.html` (since admin loads `main.css` too). Extra HTTP request. The `//go:embed` pattern handles it fine but the organizational gain is nil — animations.css would be 50-80 lines.
**Instead:** Append a labeled `/* --- Animations --- */` section to `main.css`.

### Anti-Pattern: Hardcoded Colors in New CSS

**What:** Using `#1A1F2E` or `rgba(44, 36, 24, 0.12)` in new v1.2 selectors.
**Why bad:** Dark mode is implemented via CSS custom property swap on `[data-theme="dark"]`. Hardcoded values bypass the theme system entirely and will look broken in one mode.
**Instead:** Use `var(--color-bg)`, `var(--color-surface)`, `var(--color-divider)`, etc. throughout.

### Anti-Pattern: Adding Template Variables Without Handler Updates

**What:** Adding `{{.SomeNewValue}}` in `base.html` or `list.html` without passing that key in the `data` map from the Go handler.
**Why bad:** Go's `html/template` silently renders missing map keys as an empty string by default (no error). The bug shows up visually, not as a build or runtime error.
**Instead:** Either use static HTML text (no binding needed) or add the key to `render()`'s `data` map before using `{{.Key}}` in a template.

---

## Directory Layout

```
website-go/
├── cmd/
│   └── server/
│       └── main.go           # Entry point, DI wiring
├── internal/
│   ├── config/               # Env-var config loading
│   ├── handler/
│   │   ├── blog/             # Public read handlers
│   │   ├── admin/            # Admin CRUD handlers
│   │   ├── api/              # Push API handler
│   │   └── feed/             # RSS handler
│   ├── middleware/           # Session loader, auth guards, logger
│   ├── service/
│   │   ├── post/             # Post business logic
│   │   └── auth/             # Login, bcrypt
│   ├── repository/
│   │   ├── post/             # posts table queries (pgx)
│   │   └── session/          # sessions table (scs pgxstore)
│   └── markdown/             # Goldmark wrapper, sanitizer
├── db/
│   └── migrations/           # *.sql files (embedded)
├── web/
│   ├── embed.go              # go:embed declarations for Templates + Static
│   ├── templates/
│   │   ├── base.html         # PUBLIC blog base (nav + footer) — v1.2 changes here
│   │   ├── admin-base.html   # ADMIN base (nav only, no footer) — unchanged
│   │   ├── list.html         # Post listing — v1.2 heading addition here
│   │   ├── post.html         # Single post view
│   │   ├── about.html        # About page
│   │   └── 404.html          # Themed 404
│   └── static/
│       ├── main.css          # Public blog styles — primary v1.2 edit target
│       ├── admin.css         # Admin panel styles — unchanged in v1.2
│       ├── main.js           # Dark mode toggle (66 lines) — unchanged in v1.2
│       └── admin.js          # Admin JS (163 lines) — unchanged in v1.2
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── go.mod
```

---

## Scalability Considerations

This project targets a single-instance deployment. The architecture choices reflect that.

| Concern | At current scale (1 instance) | If ever scaling out |
|---------|-------------------------------|---------------------|
| Sessions | Postgres-backed scs store — works with 1 instance | Already shareable across instances (no sticky sessions needed) |
| Image uploads | EBS bind-mount on single EC2 | Would need S3 or shared NFS — known limitation, documented out of scope |
| DB connections | pgxpool (max 10) | Sufficient; PgBouncer if needed later |
| Caching | None needed — low traffic blog | Add in-process cache for rendered HTML if traffic warrants |
| Static assets | Embedded in binary — served from memory | Already optimal; no CDN needed at blog scale |

---

## Sources

- Direct inspection of `web/templates/base.html` (2026-03-28)
- Direct inspection of `web/templates/admin-base.html` (2026-03-28)
- Direct inspection of `web/static/main.css` (645 lines, 2026-03-28)
- Direct inspection of `internal/handler/blog/handler.go` — ParseFS pattern confirmed
- Direct inspection of `web/embed.go` — go:embed declarations confirmed
- [Go html/template documentation — {{block}} and {{define}}](https://pkg.go.dev/html/template)
- [CSS Custom Properties — MDN](https://developer.mozilla.org/en-US/docs/Web/CSS/--*)
- [go:embed — Go documentation](https://pkg.go.dev/embed)

---
*Architecture research for: CSS/template integration patterns — v1.2 Shore Leave Polish*
*Researched: 2026-03-28*
