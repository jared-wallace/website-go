# Requirements: website-go

**Defined:** 2026-03-26
**Core Value:** A reader visits jared-wallace.com and reads well-rendered markdown blog posts with images in a distinctive, memorable design.

## v1 Requirements

Requirements for initial release. Each maps to roadmap phases.

### Foundation

- [ ] **FOUND-01**: Project follows standard Go layout (cmd/, internal/, etc.)
- [ ] **FOUND-02**: Postgres connection pool via pgx/v5 with health checks
- [ ] **FOUND-03**: Database migrations via goose with versioned SQL files
- [ ] **FOUND-04**: Docker multi-stage build producing minimal container
- [ ] **FOUND-05**: docker-compose with app + Postgres sidecar, EBS volume mounts
- [ ] **FOUND-06**: Makefile with build, test, lint, run, docker, and migration targets
- [ ] **FOUND-07**: GHA CI pipeline running lint, test, and build on push

### Blog

- [ ] **BLOG-01**: Reader can view published posts rendered from markdown with syntax-highlighted code
- [ ] **BLOG-02**: Reader can browse paginated post listing sorted by date
- [ ] **BLOG-03**: Reader can access posts via readable URL slugs (/posts/my-post)
- [ ] **BLOG-04**: Reader sees published date and estimated reading time on each post
- [ ] **BLOG-05**: Reader experiences weathered beach bar nautical design, mobile-responsive
- [ ] **BLOG-06**: Reader sees proper OG meta tags when links are shared
- [ ] **BLOG-07**: Search engines can discover all posts via /sitemap.xml
- [ ] **BLOG-08**: Reader sees themed 404 page for invalid URLs
- [ ] **BLOG-09**: Reader can subscribe via RSS feed at /rss with full post content
- [ ] **BLOG-10**: Reader can give thumbs-up reaction on posts (rate-limited, no auth required)
- [ ] **BLOG-11**: Reader can toggle dark mode (CSS prefers-color-scheme + manual toggle)
- [ ] **BLOG-12**: Reader sees auto-generated table of contents on long posts

### Admin

- [ ] **ADMIN-01**: Admin can log in with email/password at admin.jared-wallace.com
- [ ] **ADMIN-02**: Admin session persists across browser refresh (Postgres-backed)
- [ ] **ADMIN-03**: Admin can create posts with title, markdown body, and slug
- [ ] **ADMIN-04**: Admin can edit existing posts
- [ ] **ADMIN-05**: Admin can soft-delete posts (recoverable)
- [ ] **ADMIN-06**: Admin can toggle posts between draft and published states
- [ ] **ADMIN-07**: Admin can write in split-pane markdown editor with live preview
- [ ] **ADMIN-08**: Admin can upload images and embed them in posts
- [ ] **ADMIN-09**: Admin can push .md files via API endpoint with bearer token auth

## v2 Requirements

Deferred to future release. Tracked but not in current roadmap.

### Enhancements

- **ENH-01**: LaTeX / KaTeX math rendering in posts (front-matter gated)
- **ENH-02**: Post cover images with OG image integration
- **ENH-03**: Post series / related posts linking
- **ENH-04**: Canonical URL override for syndicated content

## Out of Scope

Explicitly excluded. Documented to prevent scope creep.

| Feature | Reason |
|---------|--------|
| Comments system | Moderation overhead; thumbs-up covers lightweight engagement |
| User registration / reader accounts | Single-author blog, zero reader value |
| OAuth / social login | One admin user; bcrypt + sessions is sufficient |
| Full-text search | Pagination + browser search sufficient at blog scale |
| Newsletter / email delivery | RSS covers follow-me use case; revisit on demand |
| In-app analytics | Use Plausible/Fathom externally instead |
| Image optimization / CDN | Overkill for single-instance personal blog |
| WYSIWYG editor | Markdown philosophy; split-pane preview is sufficient |
| Multi-author support | Personal blog; revisit if needed |

## Traceability

Which phases cover which requirements. Updated during roadmap creation.

| Requirement | Phase | Status |
|-------------|-------|--------|
| *(populated by roadmapper)* | | |

**Coverage:**
- v1 requirements: 28 total
- Mapped to phases: 0
- Unmapped: 28 ⚠️

---
*Requirements defined: 2026-03-26*
*Last updated: 2026-03-26 after initial definition*
