---
phase: 02-public-blog
verified: 2026-03-26T14:30:00Z
status: human_needed
score: 11/12 must-haves verified
human_verification:
  - test: "Verify mobile responsive layout (single-column card grid on phone-width viewport)"
    expected: "Card grid collapses to 1 column, tagline hidden in nav, pagination still accessible"
    why_human: "CSS @media breakpoints exist but mobile layout requires visual/browser inspection. 02-03 deferred this as localhost made device testing impractical. Confirm after Docker deployment or using browser DevTools."
---

# Phase 02: Public Blog Verification Report

**Phase Goal:** A reader visiting jared-wallace.com can browse, discover, and read published posts in a distinctive weathered beach bar design, on any device.
**Verified:** 2026-03-26T14:30:00Z
**Status:** human_needed (11/12 must-haves verified; 1 item requires visual confirmation)
**Re-verification:** No — initial verification

---

## Goal Achievement

### Observable Truths

| #  | Truth | Status | Evidence |
|----|-------|--------|----------|
| 1  | Post repository returns paginated published posts sorted by date descending | VERIFIED | `queries.go:19-20` — `ORDER BY created_at DESC LIMIT $1 OFFSET $2`; `pgxpool.Pool.Query` call confirmed |
| 2  | Post repository finds a single post by slug or returns ErrNotFound | VERIFIED | `queries.go:83` — `pgx.ErrNoRows` wrapped as `ErrNotFound`; sentinel test passes |
| 3  | Reading time calculation returns ceil(word_count / 200) | VERIFIED | `reading_time.go` uses `strings.Fields`; `TestReadingTime` passes (17 unit tests, 0 failures) |
| 4  | ToC extraction returns h2/h3 entries from rendered HTML, nil when fewer than 3 | VERIFIED | `toc.go` with `golang.org/x/net/html`; `TestExtractToC_FewHeadings` and `TestExtractToC_ThreeOrMore` pass |
| 5  | Excerpt generation produces plain text truncated from markdown body | VERIFIED | `excerpt.go` regex-strips markdown and truncates; `TestExcerpt_*` passes |
| 6  | Service layer computes pagination metadata (total pages, has prev/next) | VERIFIED | `list.go` — `math.Ceil`, `HasPrev`, `HasNext` computed; `TestListPublished_Pagination` passes |
| 7  | Reader sees a card grid of published posts with title, excerpt, tags, date, reading time | VERIFIED | `list.html` renders `range .Posts` with `post-card` article elements; handler calls `svc.ListPublished` |
| 8  | Reader can click a post title to view the full post at /posts/{slug} | VERIFIED | Route `GET /posts/{slug}` in `main.go:62`; `post.go` calls `svc.GetBySlug`; `post.html` renders `{{.RenderedHTML}}` |
| 9  | Reader sees nautical weathered beach bar design | VERIFIED | `main.css` 540 lines with `--color-bg`, CSS custom properties, card grid styles; Google Fonts in `base.html` |
| 10 | Reader can toggle dark mode; preference persists across page loads | VERIFIED | `main.js` — `dark-toggle` handler writes `localStorage`; `base.html` inline script reads `localStorage.getItem('theme')` before CSS load (flash prevention) |
| 11 | Reader sees a collapsible ToC on posts with 3+ headings | VERIFIED | `toc.html` block in `post.html`; `InjectHeadingIDs` wired in `get.go:31`; bug fixed in commit `3a3a728` |
| 12 | Reader sees themed "Lost at Sea" 404 for invalid URLs | VERIFIED | `404.html` contains "Lost at Sea" heading; catch-all route `GET /{path...}` registered in `main.go:63` |
| 13 | Reader sees numbered pagination when posts exceed 10 | VERIFIED | `list.go` — `PostsPerPage = 10`; `list.html` renders pagination links from `ListResult.TotalPages` |
| 14 | Reader experiences mobile-responsive design on any device | UNCERTAIN | CSS has `@media (max-width: 767px)` with `grid-template-columns: 1fr` and flex column overrides, but full mobile rendering not visually confirmed (02-03 deferred) |

**Score:** 13/14 truths verified (1 uncertain, requires human)

---

### Required Artifacts

#### Plan 02-01 (Data and Service Layer)

| Artifact | Expected | Status | Details |
|----------|----------|--------|---------|
| `db/migrations/00002_add_tags_to_posts.sql` | Tags column DDL | VERIFIED | Contains `ALTER TABLE posts ADD COLUMN tags TEXT NOT NULL DEFAULT ''`; goose Up/Down present |
| `internal/model/post.go` | Post struct with Tags field | VERIFIED | `Tags string` field present after RenderedHTML |
| `internal/repository/post/repository.go` | Repository interface + ErrNotFound + New() | VERIFIED | `type Repository interface`, `var ErrNotFound`, `func New(pool *pgxpool.Pool) Repository` all confirmed |
| `internal/service/post/service.go` | Service struct with repo injection + New() | VERIFIED | `type Service struct`, `repo postrepo.Repository` field, `func New(repo postrepo.Repository) *Service` |
| `internal/service/post/list.go` | PostsPerPage, ListResult, pagination logic | VERIFIED | `const PostsPerPage = 10`, `type ListResult struct`, `math.Ceil` pagination confirmed |
| `internal/service/post/get.go` | PostDetail, template.HTML cast | VERIFIED | `type PostDetail struct`, `template.HTML` cast in `RenderedHTML` field, `InjectHeadingIDs` wired |
| `internal/service/post/reading_time.go` | ReadingTime function | VERIFIED | `func ReadingTime`, `strings.Fields` confirmed |
| `internal/service/post/toc.go` | ToCEntry, ExtractToC, InjectHeadingIDs | VERIFIED | All three present; `golang.org/x/net/html` import confirmed |
| `internal/service/post/excerpt.go` | Excerpt, ParseTags | VERIFIED | Both functions confirmed |
| `internal/service/post/service_test.go` | Unit tests (min 80 lines) | VERIFIED | 236 lines, 17 tests covering ReadingTime/ToC/Excerpt/ParseTags/pagination; all pass |

#### Plan 02-02 (UI Layer)

| Artifact | Expected | Status | Details |
|----------|----------|--------|---------|
| `web/templates/base.html` | HTML skeleton with dark mode flash prevention | VERIFIED | `localStorage.getItem` present; inline script before CSS link |
| `web/templates/list.html` | Card grid with pagination | VERIFIED | `post-card` class, `range .Posts`, pagination markup |
| `web/templates/post.html` | Single post with ToC | VERIFIED | `toc` block present, `{{.RenderedHTML}}` in `post-body` div |
| `web/templates/404.html` | "Lost at Sea" themed 404 | VERIFIED | "Lost at Sea" heading confirmed |
| `web/static/main.css` | Nautical design system (min 150 lines) | VERIFIED | 540 lines; `--color-bg` custom property in both light and dark contexts |
| `web/static/main.js` | Dark mode toggle + ToC collapse | VERIFIED | `dark-toggle` event listener confirmed |
| `internal/handler/blog/handler.go` | BlogHandler + New() | VERIFIED | `type BlogHandler struct`, per-page template cache, `func New(svc *postservice.Service)` |
| `cmd/server/main.go` | HTTP server with routes + ListenAndServe | VERIFIED | `mux.HandleFunc` for 4 routes; `bloghandler.New(svc)` construction; `fs.Sub` static serving fix |

---

### Key Link Verification

#### Plan 02-01 Links

| From | To | Via | Status | Details |
|------|----|-----|--------|---------|
| `internal/service/post/service.go` | `internal/repository/post/repository.go` | Repository interface injection | WIRED | `repo postrepo.Repository` field at line 13; `func New(repo postrepo.Repository)` at line 17 |
| `internal/repository/post/queries.go` | `internal/model/post.go` | `model.Post` in query results | WIRED | `[]model.Post` return type, `var p model.Post` scan targets at lines 14, 28, 30 |

#### Plan 02-02 Links

| From | To | Via | Status | Details |
|------|----|-----|--------|---------|
| `internal/handler/blog/handler.go` | `internal/service/post/service.go` | Service injection in constructor | WIRED | `svc *postservice.Service` at line 21; `func New(svc *postservice.Service)` |
| `internal/handler/blog/list.go` | `web/templates/list.html` | `ExecuteTemplate` via render helper | WIRED | `h.render(w, ..., "list.html", ...)` calls `tmpl.ExecuteTemplate(w, "base", data)` |
| `cmd/server/main.go` | `internal/handler/blog/handler.go` | Route registration via blogHandler | WIRED | `blog := bloghandler.New(svc)` at line 50; all 4 `mux.HandleFunc` calls use `blog.*` methods |
| `web/templates/base.html` | `web/static/main.css` | Stylesheet link | WIRED | `<link rel="stylesheet" href="/static/main.css">` at line 11 |

---

### Data-Flow Trace (Level 4)

| Artifact | Data Variable | Source | Produces Real Data | Status |
|----------|---------------|--------|-------------------|--------|
| `web/templates/list.html` | `.Posts` (range) | `handler/blog/list.go` calls `svc.ListPublished(ctx, page)` | Yes — `repo.ListPublished` executes `SELECT ... ORDER BY created_at DESC LIMIT $1 OFFSET $2` via `pgxpool.Pool.Query` | FLOWING |
| `web/templates/post.html` | `.RenderedHTML` | `handler/blog/post.go` calls `svc.GetBySlug(ctx, slug)` | Yes — `repo.FindBySlug` executes `SELECT ... WHERE slug = $1` via `pgxpool.Pool.QueryRow`; HTML rendered by goldmark at write time | FLOWING |
| `web/templates/list.html` | `.TotalPages`, `.HasNext`, `.HasPrev` | `service/post/list.go` — `math.Ceil(total / PostsPerPage)` | Yes — `repo.CountPublished` executes `SELECT COUNT(*)` | FLOWING |
| `web/templates/post.html` | `.ToC` | `service/post/get.go` — `InjectHeadingIDs` then `ExtractToC` | Yes — parsed from real RenderedHTML; nil if < 3 headings | FLOWING |

---

### Behavioral Spot-Checks

| Behavior | Command | Result | Status |
|----------|---------|--------|--------|
| All unit tests pass (no DB needed) | `go test ./... -count=1 -short` | All 5 test packages pass; 0 failures | PASS |
| Full binary compiles | `go build ./...` | Exit 0 | PASS |
| Vet clean | `go vet ./...` | Exit 0 (no issues) | PASS |
| ErrNotFound sentinel test | `go test ./internal/repository/post/... -run TestErrNotFound -v` | `TestErrNotFound_Sentinel PASS` | PASS |
| Service unit tests | `go test ./internal/service/post/... -v -count=1` | 17 tests, all PASS | PASS |
| Handler tests | `go test ./internal/handler/blog/... -v -count=1` | 3 tests (404, empty list, not-found slug) all PASS | PASS |
| Static file serving fix | `grep "fs.Sub" cmd/server/main.go` | `staticFS, err := fs.Sub(web.Static, "static")` at line 54 | PASS |
| ToC heading ID injection | `grep "InjectHeadingIDs" internal/service/post/get.go` | Found at line 31 | PASS |

---

### Requirements Coverage

| Requirement | Source Plan | Description | Status | Evidence |
|-------------|------------|-------------|--------|----------|
| BLOG-01 | 02-02 | Reader can view published posts rendered from markdown with syntax-highlighted code | SATISFIED | `post.html` renders `{{.RenderedHTML}}`; `RenderedHTML` is goldmark-rendered HTML stored in DB; `template.HTML` cast prevents double-escape |
| BLOG-02 | 02-01, 02-02 | Reader can browse paginated post listing sorted by date | SATISFIED | `queries.go` `ORDER BY created_at DESC LIMIT/OFFSET`; `list.go` pagination math; `list.html` pagination UI |
| BLOG-03 | 02-01, 02-02 | Reader can access posts via readable URL slugs (/posts/my-post) | SATISFIED | Route `GET /posts/{slug}` in `main.go`; `r.PathValue("slug")` in `post.go`; `slug` column in DB |
| BLOG-04 | 02-01, 02-02 | Reader sees published date and estimated reading time on each post | SATISFIED | `PostSummary.ReadingTime` and `PostSummary.PublishedAt` in `list.go`; both rendered in `list.html` card and `post.html` |
| BLOG-05 | 02-02 | Reader experiences weathered beach bar nautical design, mobile-responsive | PARTIAL | Design system confirmed (540-line CSS, nautical color palette, Playfair Display + Lora fonts); `@media (max-width: 767px)` breakpoints exist for single-column grid; mobile visual rendering deferred to post-deploy |
| BLOG-08 | 02-02 | Reader sees themed 404 page for invalid URLs | SATISFIED | `404.html` "Lost at Sea"; catch-all route `GET /{path...}` returns HTTP 404; `TestNotFound` passes |
| BLOG-11 | 02-02 | Reader can toggle dark mode (CSS prefers-color-scheme + manual toggle) | SATISFIED | `main.js` `dark-toggle` handler; `base.html` inline script handles `prefers-color-scheme: dark` and `localStorage` persistence; `data-theme` attribute on `<html>` |
| BLOG-12 | 02-01, 02-02 | Reader sees auto-generated table of contents on long posts | SATISFIED | `ExtractToC` (nil < 3 headings); `InjectHeadingIDs` adds id attrs; `post.html` renders ToC block; collapse toggle in `main.js` |

**Orphaned requirements check:** BLOG-06, BLOG-07, BLOG-09, BLOG-10 are mapped to later phases in REQUIREMENTS.md — none are orphaned for Phase 2.

---

### Anti-Patterns Found

| File | Pattern | Severity | Impact |
|------|---------|----------|--------|
| None found | — | — | — |

No TODOs, FIXMEs, placeholder text, empty implementations, or hardcoded empty data arrays found in any phase 2 artifact. The two bugs found during human verification (static file serving and ToC heading IDs) were both fixed and committed in `3a3a728` before phase completion.

---

### Human Verification Required

#### 1. Mobile Responsive Layout

**Test:** Open the running blog in a browser, activate Chrome DevTools device toolbar, select an iPhone-sized viewport (375px wide). Navigate to the home page and a post page.
**Expected:**
- Card grid collapses to single column
- Nav tagline is hidden (CSS `display: none` at mobile breakpoint)
- Pagination controls remain accessible and usable
- Post body text is readable without horizontal scrolling
- ToC collapse/expand still functional
**Why human:** CSS `@media (max-width: 767px)` breakpoints confirmed in source, and `grid-template-columns: 1fr` confirmed. However, visual correctness of the full mobile layout (nav, post body width, typography scale) cannot be verified programmatically. The 02-03 summary explicitly deferred this to post-Docker-deployment verification.

---

### Gaps Summary

No blocking gaps. All 13 verifiable truths are confirmed by code inspection and passing tests. The one outstanding item (mobile visual layout) is a human verification task, not a code defect — the CSS implementation is complete and structurally correct.

The two bugs discovered during Plan 02-03 human verification (static file serving via `fs.Sub` and ToC heading ID injection via `InjectHeadingIDs`) were both resolved and committed before phase sign-off. Both fixes are confirmed present in the codebase.

---

*Verified: 2026-03-26T14:30:00Z*
*Verifier: Claude (gsd-verifier)*
