# Phase 4: Distribution - Context

**Gathered:** 2026-03-27
**Status:** Ready for planning

<domain>
## Phase Boundary

Deliver four distribution features for published posts: an RSS 2.0 feed at /rss, Open Graph meta tags for social sharing, a sitemap at /sitemap.xml with robots.txt, and a thumbs-up reaction system. No new admin features, no image upload, no API endpoints — those are Phase 5.

</domain>

<decisions>
## Implementation Decisions

### RSS Feed
- **D-01:** Most recent 25 published posts in the feed. Full post content in `<description>` (per BLOG-09).
- **D-02:** Author name "Jared Wallace" in `<managingEditor>` and per-item `<author>` fields.
- **D-03:** Tags mapped to `<category>` elements per RSS item.
- **D-04:** Auto-discovery `<link rel="alternate" type="application/rss+xml">` added to base.html so every page advertises the feed.

### Open Graph & Social Previews
- **D-05:** OG description sourced from the existing post excerpt (2-3 line excerpt from card listings). Automatic, no manual override field needed.
- **D-06:** Static site-wide fallback OG image (branded nautical graphic) used on all pages until Phase 5 adds per-post images. Image served from `/static/`.
- **D-07:** Twitter card type `summary` (small square image + title + description). Works well with the static fallback.
- **D-08:** Homepage gets site-level OG tags: `og:title="The Log"`, `og:description="dispatches from the deep end"`, `og:image=fallback`.
- **D-09:** Individual post pages get per-post OG tags: `og:title={post title}`, `og:description={post excerpt}`, `og:image=fallback`.

### Thumbs-Up Reactions
- **D-10:** Button placed below the post content, before the footer. Shows current count next to the button.
- **D-11:** One thumbs-up per reader per post (binary, not multi-clap).
- **D-12:** Rate limiting via IP-based server-side check (one per IP per post per 24h, stored in reactions table with IP hash + post ID) plus localStorage flag to hide/disable the button client-side after tapping.
- **D-13:** Tap feedback: count increments with a subtle CSS animation (bounce or color fill). No JS animation library.

### Sitemap & Crawlers
- **D-14:** Sitemap includes all published post URLs plus the homepage URL.
- **D-15:** `/robots.txt` handler returns `Sitemap: https://jared-wallace.com/sitemap.xml` plus standard `Allow` directive.

### Claude's Discretion
- RSS feed `<title>`, `<link>`, `<description>` channel-level metadata (use "The Log" / "dispatches from the deep end")
- OG fallback image design (nautical-themed, appropriate dimensions for `summary` card)
- Thumbs-up button icon and CSS animation specifics
- Sitemap `<changefreq>` and `<priority>` values
- Reactions table schema (IP hashing approach, index design)
- Whether to use a `like_count` column on posts table vs. COUNT query on reactions table
- robots.txt additional directives (if any)

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Tech Stack
- `.planning/research/STACK.md` — Authoritative dependency versions; `encoding/xml` for RSS and sitemap (no external dep)

### Project Context
- `.planning/PROJECT.md` — Core value, constraints, design direction
- `.planning/REQUIREMENTS.md` — BLOG-06, BLOG-07, BLOG-09, BLOG-10 acceptance criteria
- `.planning/ROADMAP.md` — Phase 4 success criteria (4 criteria that must be TRUE)

### Prior Phases
- `.planning/phases/01-foundation/01-CONTEXT.md` — Project layout, markdown pipeline, CI choices
- `.planning/phases/02-public-blog/02-CONTEXT.md` — Nautical design, template patterns, excerpt extraction, card layout, per-page template sets
- `.planning/phases/03-admin-panel/03-CONTEXT.md` — Host-based routing, shared design, service/handler patterns

### Infrastructure
- `.planning/research/PITFALLS.md` — Known pitfalls for Postgres and Docker

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/service/post/list.go` — `ListPublished()` fetches paginated published posts with rendered HTML. RSS feed can reuse this with a higher limit (25).
- `internal/service/post/get.go` — `GetBySlug()` returns `PostDetail` with excerpt. OG tags can use this data directly.
- `internal/model/post.go` — Post struct has Title, Slug, Body, RenderedHTML, Tags, CreatedAt, Published. Needs no changes for RSS/OG/sitemap. Reactions will need a new model or column.
- `web/templates/base.html` — `<head>` block needs OG meta tags and RSS auto-discovery link. Currently has no `{{block "meta" .}}` — will need one.
- `web/templates/post.html` — Individual post template. Needs OG overrides and thumbs-up button placement.
- `internal/handler/blog/handler.go` — `BlogHandler` with constructor DI, per-page template sets, centralized `render()`. New handlers for /rss, /sitemap.xml, /robots.txt follow the same pattern.
- `web/static/main.js` — Exists for dark mode toggle. Thumbs-up localStorage logic and fetch() call can go here.
- `web/static/main.css` — Nautical CSS with dark mode. Thumbs-up button and animation styles go here.

### Established Patterns
- Per-page template sets: parse base.html + page.html separately (avoids block name collisions)
- Constructor-based DI in main.go
- `encoding/xml` structs for XML generation (per STACK.md)
- `go:embed` for templates and static assets
- `log/slog` for structured logging

### Integration Points
- `cmd/server/main.go` — Register new routes: GET /rss, GET /sitemap.xml, GET /robots.txt, POST /posts/{slug}/react
- `db/migrations/` — New migration for reactions table (00004_create_reactions.sql)
- `web/templates/base.html` — Add RSS auto-discovery link and OG meta block to `<head>`
- `internal/handler/blog/` — New handler methods for RSS, sitemap, robots.txt, and reaction endpoint
- `internal/repository/post/` — May need new query for "most recent 25 published" and reaction CRUD

</code_context>

<specifics>
## Specific Ideas

No specific requirements — open to standard approaches. User consistently chose recommended options, confirming preference for conventional, simple patterns.

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.

</deferred>

---

*Phase: 04-distribution*
*Context gathered: 2026-03-27*
