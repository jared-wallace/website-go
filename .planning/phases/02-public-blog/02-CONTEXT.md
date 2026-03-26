# Phase 2: Public Blog - Context

**Gathered:** 2026-03-26
**Status:** Ready for planning

<domain>
## Phase Boundary

Deliver the full public reading experience: a card-grid post listing with numbered pagination, individual post pages at `/posts/{slug}`, a weathered beach bar nautical design with dark mode, auto-generated table of contents on long posts, and a themed 404 page. No admin functionality, no auth, no content creation — those are Phase 3.

</domain>

<decisions>
## Implementation Decisions

### Post Listing Layout
- **D-01:** Card grid layout — 2 columns on desktop, 1 on mobile. Each card shows title, published date, reading time, tags, and a 2-3 line excerpt.
- **D-02:** Cards have uniform height with truncated excerpts. Responsive CSS grid, no JavaScript required for layout.
- **D-03:** Home page (`/`) shows the card grid directly — blog name/tagline as header, then straight to content. No hero section or landing page.

### Post Metadata & Tags
- **D-04:** Each card and post page displays: published date, estimated reading time.
- **D-05:** Tags stored as a simple TEXT column (comma-separated) on the posts table. New migration required. Rendered as small visual labels on cards — not clickable/filterable in this phase.

### Pagination
- **D-06:** Traditional numbered pagination at the bottom of the listing. URL structure: `/posts?page=2`. Prev/Next arrows plus page numbers.
- **D-07:** 10 posts per page.

### Nautical Design
- **D-08:** Tasteful accents level — nautical color palette (warm wood, sandy off-whites, deep ocean blues) as the foundation. Rope dividers, anchor icon, wood-grain header. Subtle theming that doesn't overpower content.
- **D-09:** Typography at Claude's discretion — pick fonts that match the beach bar aesthetic and optimize readability.
- **D-10:** Dark mode uses "night beach" palette — deep navy/charcoal background, sand-gold accents, muted ocean tones. Same vibe, different lighting.
- **D-11:** Dark mode toggle as a sun/moon icon in the top-right nav bar corner. Respects `prefers-color-scheme` for initial state, persists user choice via localStorage.

### Table of Contents
- **D-12:** Inline collapsible ToC block at the top of the post, after title/metadata, before content. Click to expand/collapse.
- **D-13:** Includes h2 and h3 headings. h4+ excluded.
- **D-14:** ToC only appears when a post has 3 or more h2/h3 headings.

### Claude's Discretion
- Typography choices (font families, sizes, line heights)
- Exact color hex values within the nautical palette
- Card hover effects and transitions
- 404 page design and copy (must be on-theme)
- Reading time calculation formula (words per minute)
- Excerpt extraction approach (plain text truncation from rendered body)
- Pagination component styling details
- Nav bar layout and content

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Tech Stack
- `.planning/research/STACK.md` — Authoritative dependency versions and rationale (goldmark, bluemonday, pgx)
- `.planning/research/ARCHITECTURE.md` — Project structure guidance and package organization patterns

### Project Context
- `.planning/PROJECT.md` — Core value, constraints, design direction ("weathered bar by the beach")
- `.planning/REQUIREMENTS.md` — BLOG-01 through BLOG-05, BLOG-08, BLOG-11, BLOG-12 acceptance criteria
- `.planning/ROADMAP.md` — Phase 2 success criteria (5 criteria that must be TRUE)

### Prior Phase
- `.planning/phases/01-foundation/01-CONTEXT.md` — Phase 1 decisions (project layout, markdown pipeline, CI choices)

### Infrastructure
- `.planning/research/PITFALLS.md` — Known pitfalls for Postgres EBS bind-mount, Docker builds

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/markdown/renderer.go` — Goldmark + bluemonday pipeline with `Render()` and `RenderWithMeta()`. Front matter extraction available for metadata.
- `internal/model/post.go` — Post struct with Title, Slug, Body, RenderedHTML, Published, CreatedAt, UpdatedAt, DeletedAt. Needs `Tags` field added.
- `internal/server/server.go` — HTTP server factory with safe timeouts and graceful shutdown.
- `internal/database/` — pgx connection pool and goose migration runner.
- `internal/config/` — Environment-based configuration loader.
- `db/migrations/00001_create_posts.sql` — Posts table with slug uniqueness and published index.

### Established Patterns
- `go:embed` for templates and static assets (web/embed.go)
- `html/template` for server-side rendering
- `log/slog` for structured logging
- Environment-based config with `envOr()` defaults
- Flat `internal/` package structure (no deep nesting)

### Integration Points
- `cmd/server/main.go` — Has TODO for Phase 2 HTTP handler wiring. Needs ServeMux, route registration, and `srv.ListenAndServe()`.
- `web/templates/` — Currently has placeholder.html; needs real templates (base layout, post listing, single post, 404).
- `web/static/` — Currently has placeholder.css; needs real CSS with nautical design.
- `internal/handler/` — Empty package; will contain HTTP handlers for post listing, single post, static pages.
- New migration needed for tags column on posts table.

</code_context>

<specifics>
## Specific Ideas

No specific requirements — open to standard approaches. User consistently chose recommended options, confirming preference for conventional patterns.

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.

</deferred>

---

*Phase: 02-public-blog*
*Context gathered: 2026-03-26*
