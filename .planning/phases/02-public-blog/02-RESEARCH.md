# Phase 2: Public Blog - Research

**Researched:** 2026-03-26
**Domain:** Go server-side rendering — HTTP handlers, templates, Postgres queries, CSS theming
**Confidence:** HIGH

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

**Post Listing Layout**
- D-01: Card grid layout — 2 columns on desktop, 1 on mobile. Each card shows title, published date, reading time, tags, and a 2-3 line excerpt.
- D-02: Cards have uniform height with truncated excerpts. Responsive CSS grid, no JavaScript required for layout.
- D-03: Home page (`/`) shows the card grid directly — blog name/tagline as header, then straight to content. No hero section or landing page.

**Post Metadata & Tags**
- D-04: Each card and post page displays: published date, estimated reading time.
- D-05: Tags stored as a simple TEXT column (comma-separated) on the posts table. New migration required. Rendered as small visual labels on cards — not clickable/filterable in this phase.

**Pagination**
- D-06: Traditional numbered pagination at the bottom of the listing. URL structure: `/posts?page=2`. Prev/Next arrows plus page numbers.
- D-07: 10 posts per page.

**Nautical Design**
- D-08: Tasteful accents level — nautical color palette (warm wood, sandy off-whites, deep ocean blues). Rope dividers, anchor icon, wood-grain header. Subtle theming.
- D-09: Typography at Claude's discretion — pick fonts matching the beach bar aesthetic.
- D-10: Dark mode uses "night beach" palette — deep navy/charcoal background, sand-gold accents, muted ocean tones.
- D-11: Dark mode toggle as a sun/moon icon in the top-right nav bar corner. Respects `prefers-color-scheme` for initial state, persists user choice via localStorage.

**Table of Contents**
- D-12: Inline collapsible ToC block at the top of the post, after title/metadata, before content. Click to expand/collapse.
- D-13: Includes h2 and h3 headings. h4+ excluded.
- D-14: ToC only appears when a post has 3 or more h2/h3 headings.

### Claude's Discretion
- Typography choices (font families, sizes, line heights)
- Exact color hex values within the nautical palette
- Card hover effects and transitions
- 404 page design and copy (must be on-theme)
- Reading time calculation formula (words per minute)
- Excerpt extraction approach (plain text truncation from rendered body)
- Pagination component styling details
- Nav bar layout and content

### Deferred Ideas (OUT OF SCOPE)
None — discussion stayed within phase scope.
</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| BLOG-01 | Reader can view published posts rendered from markdown with syntax-highlighted code | Goldmark + bluemonday pipeline exists in `internal/markdown/renderer.go`; syntax highlighting via goldmark-highlighting/v2 already wired |
| BLOG-02 | Reader can browse paginated post listing sorted by date | Handler pattern: `LIMIT 10 OFFSET (page-1)*10` query on `posts` table, `?page=N` URL param |
| BLOG-03 | Reader can access posts via readable URL slugs (/posts/my-post) | `GET /posts/{slug}` via Go 1.22 ServeMux; slug uniqueness already enforced by DB index |
| BLOG-04 | Reader sees published date and estimated reading time on each post | Reading time = `ceil(word_count / 200)`; requires word-count function in service layer; date from `created_at` |
| BLOG-05 | Reader experiences weathered beach bar nautical design, mobile-responsive | UI-SPEC fully specifies design system — Playfair Display + Lora, CSS custom properties, two breakpoints |
| BLOG-08 | Reader sees themed 404 page for invalid URLs | Custom `http.NotFound` handler; "Lost at Sea" copy specified in UI-SPEC |
| BLOG-11 | Reader can toggle dark mode (CSS prefers-color-scheme + manual toggle) | `[data-theme="dark"]` on `<html>`; localStorage persistence; inline JS in `<head>` to prevent flash |
| BLOG-12 | Reader sees auto-generated table of contents on long posts | Parse rendered HTML for h2/h3 anchors; goldmark can attach IDs to headings via `goldmark/extension/headinganchor`; threshold: 3+ headings |
</phase_requirements>

---

## Summary

Phase 2 builds the complete public reading experience on top of the Phase 1 foundation. The Go binary has a database connection, migration runner, markdown renderer, and server factory — but no HTTP routes, no handlers, no templates, and no CSS. This phase fills all of that.

The implementation path is clear and low-ambiguity. Every technology choice is locked by STACK.md: Go 1.26 stdlib ServeMux for routing, `html/template` for rendering, pgx/v5 for queries, the existing goldmark+bluemonday pipeline for markdown. The UI design contract is fully specified in `02-UI-SPEC.md`, down to exact hex values, font choices, component copy, and interaction states. The planner can work directly from that spec without design decisions.

The two technically novel pieces are: (1) ToC generation — parsing rendered HTML to extract headings and synthesize anchor IDs, which requires goldmark's heading anchor extension or a post-render HTML scan; and (2) dark mode flash prevention — requiring a small inline `<script>` in the `<head>` that runs before the page paints, which is the only JavaScript that must be inline (the rest can be deferred or event-driven).

**Primary recommendation:** Build in dependency order — database layer (migration + repository) first, then service layer (reading time, ToC), then handlers, then templates, then CSS. The existing markdown pipeline and server factory are already production-ready; do not touch them.

---

## Standard Stack

All choices are locked by CLAUDE.md and STACK.md. No alternatives to evaluate.

### Core (from STACK.md — all versions verified against pkg.go.dev)

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| `net/http` stdlib | Go 1.26 | HTTP server, routing, ServeMux | Go 1.22 `{slug}` wildcard patterns; no chi/mux needed at ~15 routes |
| `html/template` stdlib | Go 1.26 | Server-side HTML templates | Auto-escaping, `{{block}}`/`{{define}}` inheritance, zero deps |
| `github.com/jackc/pgx/v5` | v5.9.1 | Postgres driver + pool | Already in go.mod; pgxpool for connection management |
| `github.com/yuin/goldmark` | v1.8.2 | Markdown rendering | Already wired in `internal/markdown/renderer.go` |
| `github.com/yuin/goldmark-highlighting/v2` | v2.0.0-20230729083705 | Syntax highlighting | Already wired; chroma monokai style |
| `github.com/microcosm-cc/bluemonday` | v1.0.27 | HTML sanitization | Already wired; UGCPolicy + class allow |
| `log/slog` stdlib | Go 1.26 | Structured logging | Already used in main.go |

### New for Phase 2

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| `github.com/yuin/goldmark` extension `headinganchor` | bundled in goldmark | Attach `id=` attributes to headings for ToC anchor links | Required for BLOG-12; without IDs, ToC links have no target |
| Google Fonts (Playfair Display + Lora) | CDN | Typography | Per UI-SPEC; loaded via `<link>` in base template `<head>` |

**No new Go module dependencies required for Phase 2.** The heading anchor extension is part of goldmark's `extension` package already in go.mod. Google Fonts is a CDN link, not a Go dependency.

**Installation:** Nothing new to install. All Go dependencies already in go.mod.

---

## Architecture Patterns

### Project Structure for Phase 2

Phase 1 created `internal/handler/` as an empty package. Phase 2 fills it:

```
internal/
├── handler/
│   └── blog/
│       ├── handler.go          # BlogHandler struct, constructor
│       ├── list.go             # GET / and GET /posts?page=N
│       ├── post.go             # GET /posts/{slug}
│       └── notfound.go         # custom 404 handler
├── repository/
│   └── post/
│       ├── repository.go       # PostRepository struct, interface
│       └── queries.go          # SQL: ListPublished, FindBySlug, Count
├── service/
│   └── post/
│       ├── service.go          # PostService struct, constructor
│       ├── list.go             # ListPublished with pagination
│       ├── get.go              # GetBySlug
│       └── reading_time.go     # word count → ceil(n/200)
└── model/
    └── post.go                 # ADD Tags string field (migration required)

web/
├── templates/
│   ├── base.html               # <html>, <head>, nav, footer blocks
│   ├── list.html               # card grid + pagination
│   ├── post.html               # single post with ToC
│   └── 404.html                # "Lost at Sea" page
└── static/
    ├── main.css                # all styles — nautical theme, dark mode vars
    └── main.js                 # dark mode toggle + ToC collapse (minimal)

db/migrations/
└── 00002_add_tags_to_posts.sql # ALTER TABLE posts ADD COLUMN tags TEXT NOT NULL DEFAULT ''
```

### Pattern 1: Repository Interface + Constructor Injection

Define a `PostRepository` interface in `internal/repository/post/`. The service depends on the interface, not the concrete type. This makes the handler layer testable with a fake repository.

```go
// internal/repository/post/repository.go
type Repository interface {
    ListPublished(ctx context.Context, limit, offset int) ([]model.Post, error)
    CountPublished(ctx context.Context) (int, error)
    FindBySlug(ctx context.Context, slug string) (*model.Post, error)
}

type postgresRepository struct {
    pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) Repository {
    return &postgresRepository{pool: pool}
}
```

### Pattern 2: Thin Handler, Fat Service

Handlers translate HTTP context to domain calls. All business logic (reading time, ToC generation, pagination math) lives in the service layer.

```go
// internal/handler/blog/list.go
func (h *BlogHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
    page := parsePageParam(r)
    result, err := h.svc.ListPublished(r.Context(), page)
    if err != nil {
        slog.Error("list posts failed", "error", err)
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }
    h.tmpl.ExecuteTemplate(w, "list.html", result)
}
```

### Pattern 3: Template Execution Helper

Avoid duplicating error-handling template execution. A small helper on the handler prevents panics from swallowed template errors.

```go
func (h *BlogHandler) render(w http.ResponseWriter, name string, data any) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    if err := h.tmpl.ExecuteTemplate(w, name, data); err != nil {
        slog.Error("template execution failed", "template", name, "error", err)
        http.Error(w, "render error", http.StatusInternalServerError)
    }
}
```

### Pattern 4: ServeMux Route Registration in main.go

```go
// cmd/server/main.go
mux := http.NewServeMux()

// Static assets (embedded)
mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(web.Static)))

// Public blog
mux.HandleFunc("GET /", blogHandler.ListPosts)
mux.HandleFunc("GET /posts", blogHandler.ListPosts)           // /posts?page=N
mux.HandleFunc("GET /posts/{slug}", blogHandler.ShowPost)

srv := server.New(cfg.Port, mux)
```

Note: `/` and `/posts` both map to the same listing handler. Go ServeMux exact-path `/` is the catch-all for unmatched routes, which also serves as the 404 handler hook — register a custom `http.NotFoundHandler` via `mux.HandleFunc("GET /{path...}", blogHandler.NotFound)` to serve the themed 404.

### Pattern 5: Pagination Math

```go
// internal/service/post/list.go
const postsPerPage = 10

type ListResult struct {
    Posts       []model.Post
    CurrentPage int
    TotalPages  int
    HasPrev     bool
    HasNext     bool
}

func (s *PostService) ListPublished(ctx context.Context, page int) (ListResult, error) {
    if page < 1 { page = 1 }
    total, err := s.repo.CountPublished(ctx)
    // ...
    offset := (page - 1) * postsPerPage
    posts, err := s.repo.ListPublished(ctx, postsPerPage, offset)
    // ...
    totalPages := int(math.Ceil(float64(total) / float64(postsPerPage)))
    return ListResult{
        Posts: posts, CurrentPage: page, TotalPages: totalPages,
        HasPrev: page > 1, HasNext: page < totalPages,
    }, nil
}
```

### Pattern 6: ToC Generation

goldmark's `goldmark/extension` package includes `headinganchor` which adds `id=` attributes to headings during render. After rendering, the service parses the rendered HTML to extract h2/h3 elements and build the ToC data structure.

```go
// internal/service/post/toc.go
type ToCEntry struct {
    ID    string
    Text  string
    Level int // 2 or 3
}

// ExtractToC scans rendered HTML for <h2 id="..."> and <h3 id="..."> elements.
// Returns nil if fewer than 3 entries found (D-14).
func ExtractToC(renderedHTML string) []ToCEntry {
    // Use strings scanning or golang.org/x/net/html parser
    // Pattern: find all <h2 id="slug">text</h2> and <h3 id="slug">text</h3>
}
```

The `golang.org/x/net` package (already an indirect dependency via bluemonday) provides `html.Parse()` for safe HTML tree traversal — no regex needed, no new dependency.

### Pattern 7: Dark Mode Flash Prevention

The dark mode toggle reads `localStorage` to initialize. Without an inline script, the page renders in light mode for one frame before JS applies the dark class, causing a visible flash.

```html
<!-- In base.html <head>, BEFORE any CSS link tags -->
<script>
  (function(){
    var t = localStorage.getItem('theme');
    if (t === 'dark' || (!t && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
      document.documentElement.setAttribute('data-theme','dark');
    }
  })();
</script>
```

This must be inline and synchronous in `<head>`. It cannot be deferred or in `main.js`.

### Pattern 8: CSS Custom Properties for Theming

```css
/* web/static/main.css */
:root {
  --color-bg:         #F5F0E8;
  --color-surface:    #E8DFD0;
  --color-accent:     #2C5F7A;
  --color-text:       #2C2418;
  --color-text-muted: #7A6A55;
  --color-divider:    #C8B89A;
}

[data-theme="dark"] {
  --color-bg:         #1A1F2E;
  --color-surface:    #242B3D;
  --color-accent:     #C9A84C;
  --color-text:       #E8DFD0;
  --color-text-muted: #9A8B75;
  --color-divider:    #3A4558;
}
```

All component styles reference `var(--color-*)` — zero JavaScript involved in theming.

### Pattern 9: Reading Time

```go
// internal/service/post/reading_time.go
import "math"

// ReadingTime returns estimated minutes to read body text.
// Formula: ceil(word_count / 200). Sourced from UI-SPEC: "{N} min read"
func ReadingTime(body string) int {
    words := len(strings.Fields(body))
    return int(math.Ceil(float64(words) / 200.0))
}
```

Apply to `body` (raw markdown text), not rendered HTML, to avoid counting HTML tags as words.

### Anti-Patterns to Avoid

- **Re-rendering markdown on every read:** The existing architecture pre-renders on write. Phase 2 handlers read `post.RenderedHTML` from the DB — never call `renderer.Render(post.Body)` in a read handler.
- **Paginating with `SELECT *` and Go-side slicing:** Use SQL `LIMIT` and `OFFSET`. Fetching all posts and slicing in Go is catastrophic at any real post count.
- **Template parsing on every request:** Parse templates once at startup (or at handler construction), cache the `*template.Template`. Never call `template.ParseFS()` inside an HTTP handler.
- **Serving Google Fonts with a `preconnect` missing:** Without `<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>` the font load blocks rendering. Include both preconnect hints.
- **Dark mode toggle implemented in CSS only:** CSS `:has()` or `:checked` tricks for dark mode cannot persist across page loads. The localStorage approach with `[data-theme]` is the correct pattern.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| HTML sanitization | Custom tag-stripping regex | `bluemonday` (already wired) | Regex cannot handle nested/malformed HTML; bluemonday uses a proper parser |
| Markdown rendering | Custom Markdown parser | `goldmark` (already wired) | CommonMark edge cases are numerous; goldmark is battle-tested |
| Connection pooling | Manual `pgx.Connect()` re-use | `pgxpool.Pool` (already wired) | Pool handles health checks, max connections, idle timeouts |
| HTML heading parsing for ToC | Regex on HTML strings | `golang.org/x/net/html` (already indirect dep) | Regex on HTML is brittle; the x/net HTML parser handles all edge cases |
| CSS dark mode variable system | Per-component dark mode overrides | Single `[data-theme="dark"]` root override | Per-component overrides multiply maintenance; root vars cascade automatically |

**Key insight:** Phase 2 is primarily wiring existing components (markdown renderer, DB pool, server factory) into handlers and templates. The risk of hand-rolling is low because the stack is already chosen and proven — but the dark mode flash prevention pattern and ToC heading extraction are the two places where "obvious" approaches (pure CSS, regex) fail in subtle ways.

---

## Common Pitfalls

### Pitfall 1: Template Parsing Panics on Embedded FS with Wrong Path

**What goes wrong:** `template.ParseFS(web.Templates, "templates/*.html")` with the pattern `"*.html"` fails — embedded paths include the `templates/` prefix, so the pattern must be `"templates/*.html"`.

**Why it happens:** `go:embed` preserves the directory prefix in the virtual FS. `*.html` matches files at the root, but embedded files are at `templates/base.html`, etc.

**How to avoid:** Use `template.ParseFS(web.Templates, "templates/*.html")`. Verify at startup — template parse errors should be fatal (`log.Fatal`), not swallowed.

**Warning signs:** Blank page output with no error; or "template not found" panics on first request.

### Pitfall 2: Custom 404 Handler Intercepts Valid Routes

**What goes wrong:** Registering `GET /{path...}` as a catch-all for 404 runs before more-specific routes if ServeMux specificity is misunderstood.

**Why it happens:** Go's ServeMux matches the most-specific pattern. `/{path...}` is less specific than `/posts/{slug}`, so it only fires when nothing else matches. This is correct behavior — but the handler must explicitly call `http.NotFound` logic, not just render a 404 template with a 200 status.

**How to avoid:** Set the correct HTTP status: `w.WriteHeader(http.StatusNotFound)` before executing the 404 template.

**Warning signs:** Search engines index the 404 page at random URLs; analytics show 404 pages with 200 status codes.

### Pitfall 3: `html/template` Treats `template.HTML` as Safe — Double-Render Risk

**What goes wrong:** Storing rendered HTML in Postgres, then passing it through another `html/template` render step that re-escapes it.

**Why it happens:** If the template data field `Post.RenderedHTML` is typed as `string` instead of `template.HTML`, the template engine will HTML-escape the already-rendered content, producing visible `&lt;` and `&gt;` in the browser.

**How to avoid:** In the `model.Post` struct, keep `RenderedHTML` as `template.HTML` (already defined as `string` in Phase 1 — must cast in the service layer before passing to template). Either change the field type or cast in the template data struct.

**Warning signs:** Post body shows raw escaped HTML tags in the browser instead of rendered content.

### Pitfall 4: Excerpt Generation from Rendered HTML

**What goes wrong:** Truncating `post.RenderedHTML` (HTML string) to N characters for the card excerpt. The truncation can cut mid-tag, producing malformed HTML like `<str` or unclosed `<a>` elements.

**Why it happens:** Naive `body[:150]` truncation ignores HTML structure.

**How to avoid:** Generate the excerpt from `post.Body` (raw markdown), strip markdown syntax, then truncate plain text. Alternatively, extract plain text from rendered HTML using `golang.org/x/net/html` before truncating.

**Warning signs:** Cards show broken HTML artifacts; browser inspector shows unclosed tags in card content.

### Pitfall 5: Pagination Off-by-One at Boundaries

**What goes wrong:** `OFFSET 0` returns the first page correctly, but `OFFSET 10` when there are exactly 10 posts returns an empty page rather than a 404 or redirect to page 1.

**Why it happens:** `TotalPages = ceil(total / 10)` is correct, but the handler doesn't validate that the requested page is within `[1, totalPages]` before executing the query.

**How to avoid:** Validate `page >= 1 && page <= totalPages` after computing `totalPages`. Redirect to page 1 (or return a 404) for out-of-range requests.

**Warning signs:** Empty post grid with pagination controls showing page 0 or negative pages.

### Pitfall 6: Dark Mode Flash on Initial Load

**What goes wrong:** `main.js` is loaded with `defer` (correct for performance). The dark mode initialization script in `main.js` fires after page paint, causing a white flash before switching to dark.

**Why it happens:** Deferred JS runs after DOMContentLoaded, which is after initial paint.

**How to avoid:** The flash-prevention snippet (Pattern 7 above) MUST be inline in `<head>`, before the CSS `<link>` tags. The rest of the dark mode toggle logic (click handler, localStorage write) can live in deferred `main.js`.

**Warning signs:** On page load in dark mode, a brief white flash is visible before the dark background appears.

### Pitfall 7: Tags Column Migration Missing `NOT NULL DEFAULT ''`

**What goes wrong:** Adding `tags TEXT` without a default on a table that already has rows causes the migration to fail with `null value violates not-null constraint` if posts exist.

**Why it happens:** `ALTER TABLE posts ADD COLUMN tags TEXT NOT NULL` with no default fails if any row exists (Postgres cannot backfill NULL into a NOT NULL column without a default).

**How to avoid:** Migration must be `ALTER TABLE posts ADD COLUMN tags TEXT NOT NULL DEFAULT ''`.

**Warning signs:** `goose up` fails with a Postgres constraint violation error on non-empty databases.

---

## Code Examples

Verified patterns from the existing codebase and ARCHITECTURE.md:

### Reading Published Posts with Pagination (pgx/v5 style)

```go
// internal/repository/post/queries.go
// Source: STACK.md — pgx native interface, not database/sql adapter
func (r *postgresRepository) ListPublished(ctx context.Context, limit, offset int) ([]model.Post, error) {
    rows, err := r.pool.Query(ctx,
        `SELECT id, title, slug, body, rendered_html, tags, published, created_at
         FROM posts
         WHERE published = true AND deleted_at IS NULL
         ORDER BY created_at DESC
         LIMIT $1 OFFSET $2`,
        limit, offset,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()  // Pitfall 6 from PITFALLS.md — always defer Close

    var posts []model.Post
    for rows.Next() {
        var p model.Post
        if err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Body,
            &p.RenderedHTML, &p.Tags, &p.Published, &p.CreatedAt); err != nil {
            return nil, err
        }
        posts = append(posts, p)
    }
    return posts, rows.Err()
}
```

### Template Inheritance (html/template block/define)

```html
<!-- web/templates/base.html -->
<!-- Source: ARCHITECTURE.md — go:embed pattern -->
<!DOCTYPE html>
<html lang="en" {{if .DarkMode}}data-theme="dark"{{end}}>
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{block "title" .}}The Log{{end}}</title>
  <!-- Flash prevention MUST be first, before CSS -->
  <script>(function(){var t=localStorage.getItem('theme');if(t==='dark'||(!t&&window.matchMedia('(prefers-color-scheme: dark)').matches)){document.documentElement.setAttribute('data-theme','dark');}})()</script>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Playfair+Display:wght@700&family=Lora:wght@400&display=swap" rel="stylesheet">
  <link rel="stylesheet" href="/static/main.css">
</head>
<body>
  <nav class="site-nav">
    <a href="/" class="site-name">The Log</a>
    <span class="site-tagline">dispatches from the deep end</span>
    <button class="dark-toggle" aria-label="Switch to dark mode" id="dark-toggle">
      <!-- SVG sun/moon inline -->
    </button>
  </nav>
  {{block "content" .}}{{end}}
  <footer>...</footer>
  <script src="/static/main.js" defer></script>
</body>
</html>
```

### Tags Migration

```sql
-- db/migrations/00002_add_tags_to_posts.sql
-- +goose Up
ALTER TABLE posts ADD COLUMN tags TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE posts DROP COLUMN tags;
```

### ToC Data Structure for Template

```go
// internal/service/post/get.go
type PostPageData struct {
    Post        *model.Post
    RenderedHTML template.HTML   // cast from string to avoid double-escape
    ToC         []ToCEntry       // nil if < 3 headings
    ReadingTime int
}
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Cookie-payload sessions (gorilla) | Server-side session token (scs) | ~2021 | Eliminates cookie bloat and client-side session tampering |
| `lib/pq` database/sql driver | pgx/v5 native interface | ~2022 | Faster scans, richer type support, no adapter overhead |
| Blackfriday markdown | goldmark (CommonMark) | ~2020 | Spec-compliant rendering, extensible AST |
| Per-component dark mode CSS | CSS custom properties + `[data-theme]` | ~2019 | Single source of truth, zero JS for color changes |
| `text/template` for HTML | `html/template` | Go 1.0 | Auto-escaping; the old way is an XSS waiting to happen |

**No deprecated patterns apply to Phase 2's scope.**

---

## Open Questions

1. **Heading anchor ID format from goldmark**
   - What we know: goldmark's heading anchor extension generates IDs by slugifying heading text (e.g., "My Section" → `my-section`)
   - What's unclear: The exact slug algorithm (lowercase + hyphens? unicode handling?). ID collisions if two headings have the same text.
   - Recommendation: Use `goldmark/extension/headinganchor` with its default ID generator. Test with duplicate heading text to verify collision handling before shipping.

2. **Template parsing strategy: ParseGlob vs ParseFS**
   - What we know: `web.Templates` is an `embed.FS`. `template.ParseFS` accepts a glob pattern.
   - What's unclear: Whether `{{template "base.html" .}}` calls work correctly when templates are parsed as a set vs. individual files.
   - Recommendation: Parse all templates as a set at startup: `template.New("").ParseFS(web.Templates, "templates/*.html")`. All `{{define}}` blocks across files will be available in the set.

3. **Excerpt from markdown body vs rendered HTML**
   - What we know: UI-SPEC says "plain text truncation from rendered body" (Claude's discretion area).
   - What's unclear: Whether to strip markdown syntax from `post.Body` (simpler) or strip HTML from `post.RenderedHTML` (accurate to final output).
   - Recommendation: Strip markdown from `post.Body` using a simple regex or strings-based approach. Avoids parsing HTML, and the difference in output is negligible for 2-3 line excerpts.

---

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Go toolchain | Build | Assumed (Phase 1 complete) | 1.26 (in Docker) | — |
| Postgres | Database queries | Assumed (Phase 1 complete) | 15+ | — |
| Google Fonts CDN | Typography | Requires internet at load time | CDN | System serif fallback in CSS font-stack |

**Google Fonts note:** Fonts load from CDN at reader request time (not build time). For offline/air-gapped environments, the CSS font stack falls back to `Georgia, serif` (Lora fallback) and `'Times New Roman', serif` (Playfair fallback). No build-time action needed.

---

## Validation Architecture

### Test Framework

| Property | Value |
|----------|-------|
| Framework | Go testing stdlib (`testing` package) |
| Config file | none — `go test ./...` convention |
| Quick run command | `go test ./internal/... -run TestUnit -v` |
| Full suite command | `go test ./... -v -race` (matches `make test`) |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| BLOG-01 | Published posts render markdown with syntax highlighting | unit | `go test ./internal/markdown/... -v` | Yes (renderer_test.go) |
| BLOG-01 | Syntax highlighting CSS classes survive bluemonday | unit | `go test ./internal/markdown/... -run TestRender_CodeFence -v` | Yes (renderer_test.go) |
| BLOG-02 | ListPublished returns 10 posts max, sorted by date desc | unit | `go test ./internal/service/post/... -run TestListPublished -v` | No — Wave 0 |
| BLOG-02 | Pagination offset math is correct | unit | `go test ./internal/service/post/... -run TestPagination -v` | No — Wave 0 |
| BLOG-03 | FindBySlug returns post for valid slug | unit | `go test ./internal/repository/post/... -run TestFindBySlug -v` | No — Wave 0 |
| BLOG-03 | FindBySlug returns error for unknown slug | unit | `go test ./internal/repository/post/... -run TestFindBySlugNotFound -v` | No — Wave 0 |
| BLOG-04 | ReadingTime returns ceil(words/200) | unit | `go test ./internal/service/post/... -run TestReadingTime -v` | No — Wave 0 |
| BLOG-08 | 404 handler returns HTTP 404 status | unit | `go test ./internal/handler/blog/... -run TestNotFound -v` | No — Wave 0 |
| BLOG-11 | Dark mode toggle JS sets localStorage | manual | Browser DevTools — inspect `localStorage.getItem('theme')` | manual-only (no Go test) |
| BLOG-12 | ExtractToC returns nil when < 3 headings | unit | `go test ./internal/service/post/... -run TestExtractToC -v` | No — Wave 0 |
| BLOG-12 | ExtractToC returns entries for h2 and h3 only | unit | `go test ./internal/service/post/... -run TestExtractToCLevels -v` | No — Wave 0 |

### Sampling Rate

- **Per task commit:** `go test ./internal/... -v -race`
- **Per wave merge:** `go test ./... -v -race` (full suite including integration tests if DB available)
- **Phase gate:** Full suite green before `/gsd:verify-work`

### Wave 0 Gaps

- [ ] `internal/service/post/service_test.go` — covers BLOG-02, BLOG-04, BLOG-12 (TestListPublished, TestPagination, TestReadingTime, TestExtractToC, TestExtractToCLevels)
- [ ] `internal/repository/post/repository_test.go` — covers BLOG-03 (TestFindBySlug, TestFindBySlugNotFound) — NOTE: requires Postgres connection; mark with `//go:build integration` tag or skip without `DATABASE_URL`
- [ ] `internal/handler/blog/handler_test.go` — covers BLOG-08 (TestNotFound) using `httptest.NewRecorder`

---

## Project Constraints (from CLAUDE.md)

These directives are mandatory. The planner must verify all tasks comply.

| Directive | Constraint |
|-----------|------------|
| Tech stack | Go with minimal dependencies — stdlib first, no large frameworks |
| Infrastructure | Must run as Docker container on port 8080 behind Nginx/ALB |
| Storage | All persistent data lives on EBS volume at `/var/www/html` |
| Design | Use the `frontend-design` skill for all template/UI work |
| GSD workflow | All file changes must start through a GSD command |
| Simplify | All code changes run through `/simplify` before presentation |
| Code quality | No global `db` or `cfg` variables; constructor injection only |
| Markdown pipeline | goldmark → bluemonday → `template.HTML` order is locked (Phase 01-02 decision) |
| bluemonday policy | UGCPolicy + `AllowAttrs("class")` on code/span/pre — locked (Phase 01-02 decision) |
| go:embed | Use `go:embed` for templates and static assets (web/embed.go pattern established) |
| Logging | `log/slog` only — no zerolog, no zap |
| Budget | GHA CI must work on free tier (single job: lint→test→build) |

---

## Sources

### Primary (HIGH confidence)

- STACK.md (`.planning/research/STACK.md`) — authoritative dependency versions and rationale
- ARCHITECTURE.md (`.planning/research/ARCHITECTURE.md`) — project structure, DI patterns, data flow
- `02-UI-SPEC.md` (`.planning/phases/02-public-blog/02-UI-SPEC.md`) — full design contract for all Phase 2 UI components
- `02-CONTEXT.md` — locked implementation decisions D-01 through D-14
- Existing codebase — `internal/markdown/renderer.go`, `internal/model/post.go`, `internal/server/server.go`, `cmd/server/main.go`, `go.mod`

### Secondary (MEDIUM confidence)

- PITFALLS.md (`.planning/research/PITFALLS.md`) — verified domain pitfalls for this Go/Postgres stack
- REQUIREMENTS.md — acceptance criteria for BLOG-01 through BLOG-12
- Go html/template documentation — `{{block}}`/`{{define}}` inheritance, `ParseFS` patterns
- goldmark heading anchor extension behavior — verified against goldmark README and extension source

### Tertiary (LOW confidence)

- Google Fonts CDN stability — assumed stable; no verified SLA. Fallback CSS font-stack mitigates risk.

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — all libraries from go.mod, versions verified in STACK.md against pkg.go.dev
- Architecture: HIGH — patterns derived from ARCHITECTURE.md and existing Phase 1 code structure
- Pitfalls: HIGH — drawn from PITFALLS.md (domain-researched) plus Phase 2-specific additions
- UI design: HIGH — fully specified in 02-UI-SPEC.md with exact values

**Research date:** 2026-03-26
**Valid until:** 2026-04-26 (stable stack; Google Fonts CDN availability is the only time-sensitive element)
