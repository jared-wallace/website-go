---
phase: 02-public-blog
plan: "02"
subsystem: ui
tags: [go, html-template, css, javascript, dark-mode, nautical-design, blog-ui]

requires:
  - phase: 02-01
    provides: PostService with ListPublished/GetBySlug, PostSummary, PostDetail, ToCEntry types

provides:
  - Nautical "weathered beach bar" HTML templates (base, list, post, 404)
  - CSS design system with light/dark custom properties using CSS variables
  - JavaScript dark mode toggle with localStorage persistence and ToC collapse
  - BlogHandler HTTP layer connecting service to templates
  - Full HTTP server with routes registered and serving traffic

affects: [02-03, 03-admin]

tech-stack:
  added: []
  patterns:
    - Per-page template sets (base.html + page.html) to avoid block name collisions in html/template
    - Render helper injects Year for footer; per-request data map keeps templates clean
    - Dark mode via CSS custom properties + data-theme attribute on html element
    - Flash prevention via inline script before CSS link in head

key-files:
  created:
    - web/templates/base.html
    - web/templates/list.html
    - web/templates/post.html
    - web/templates/404.html
    - web/static/main.css
    - web/static/main.js
    - internal/handler/blog/handler.go
    - internal/handler/blog/list.go
    - internal/handler/blog/post.go
    - internal/handler/blog/notfound.go
    - internal/handler/blog/handler_test.go
  modified:
    - cmd/server/main.go

key-decisions:
  - "Per-page template sets: parse base.html + page.html into separate template.Template instances to prevent block name collisions when all pages define 'content' and 'title' blocks"
  - "ExecuteTemplate calls 'base' named template (not page name) after parsing base+page pair into each set"

patterns-established:
  - "Template caching pattern: BlogHandler.templates map[string]*template.Template keyed by page filename"
  - "Render helper pattern: inject Year and Content-Type before WriteHeader to prevent double-write"
  - "Route pattern: GET /{$} exact home, GET /posts list, GET /posts/{slug} post, GET /{path...} catch-all 404"

requirements-completed: [BLOG-01, BLOG-02, BLOG-03, BLOG-04, BLOG-05, BLOG-08, BLOG-11, BLOG-12]

duration: 32min
completed: 2026-03-26
---

# Phase 2 Plan 02: Public Blog UI Summary

**Nautical "weathered beach bar" blog UI with CSS design system, dark mode toggle, ToC collapse, and full HTTP stack serving traffic at localhost:8080**

## Performance

- **Duration:** 32 min
- **Started:** 2026-03-26T13:09:04Z
- **Completed:** 2026-03-26T13:40:51Z
- **Tasks:** 3
- **Files modified:** 12

## Accomplishments

- Full nautical design system: 540-line CSS with light/dark custom properties, card grid, typography, pagination, and all component styles
- Dark mode toggle with localStorage persistence, flash prevention via inline script before CSS
- BlogHandler connecting post service to HTML templates with per-page template caching
- Full HTTP server wired end-to-end: pool -> repo -> service -> handler -> mux -> ListenAndServe

## Task Commits

Each task was committed atomically:

1. **Task 1: HTML templates, CSS design system, and JavaScript** - `3406d2f` (feat)
2. **Task 2: Blog HTTP handlers** - `63df1d9` (feat)
3. **Task 3: Wire routes in main.go** - `3beec1d` (feat)

**Plan metadata:** (docs commit follows)

## Files Created/Modified

- `web/templates/base.html` - Master layout: dark mode flash prevention, Google Fonts, nav with sun/moon toggle, footer with Year
- `web/templates/list.html` - Card grid with tags, excerpt, metadata, numbered pagination, empty state
- `web/templates/post.html` - Single post with ToC collapse, rope divider, post body
- `web/templates/404.html` - "Lost at Sea" themed 404 with anchor SVG, CTA button
- `web/static/main.css` - 540-line nautical design system with CSS custom properties for light/dark mode
- `web/static/main.js` - Dark mode toggle + ToC collapse; no external dependencies
- `internal/handler/blog/handler.go` - BlogHandler struct, New() with per-page template sets, render helper
- `internal/handler/blog/list.go` - ListPosts: page param parsing, service call, template render
- `internal/handler/blog/post.go` - ShowPost: PathValue slug, ErrNotFound -> 404
- `internal/handler/blog/notfound.go` - NotFound handler returning HTTP 404
- `internal/handler/blog/handler_test.go` - Tests: 404 status, empty list 200, not-found slug 404
- `cmd/server/main.go` - Full server wiring with all routes, graceful shutdown

## Decisions Made

- **Per-page template sets:** Go's `html/template` overwrites block names when parsing multiple files into one set; if list.html and post.html both define `{{define "content"}}`, the last one parsed wins for all pages. Fixed by parsing `base.html + page.html` into a separate `*template.Template` per page.
- **ExecuteTemplate("base", data):** Each per-page template set is executed via the "base" named template (not the page file name), which then resolves the "content" and "title" blocks from the paired page file.
- **Worktree rebase onto new:** The worktree branch `worktree-agent-acc24b4b` was based off the pre-Phase-1 SAVEPOINT. Rebased onto `new` branch to get Phase 1 work before starting Phase 2.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Per-page template sets to fix block name collision**
- **Found during:** Task 2 (Blog HTTP handlers)
- **Issue:** Go `html/template` ParseFS with all pages in one set causes the last-parsed `{{define "content"}}` to win for all pages. The plan's example called `ExecuteTemplate(w, "list.html", data)` but the plan also noted "adjust accordingly" if needed.
- **Fix:** Parsed each page as `base.html + page.html` into a separate `*template.Template`; stored in `map[string]*template.Template` keyed by page name. Handler's `render()` calls `ExecuteTemplate(w, "base", data)`.
- **Files modified:** internal/handler/blog/handler.go
- **Verification:** All 3 handler tests pass; `go build ./...` succeeds
- **Committed in:** `63df1d9` (Task 2 commit)

---

**Total deviations:** 1 auto-fixed (Rule 1 - bug in template rendering approach)
**Impact on plan:** Required fix for correct rendering. All acceptance criteria met.

## Issues Encountered

- Worktree branch was rebased onto `new` before starting — not a deviation, just setup.

## Known Stubs

None — all data flows from service layer through to templates. No placeholder text or hardcoded empty values in rendered content.

## Next Phase Readiness

- HTTP server fully wired and compiles; ready for Phase 2 Plan 03 (RSS feed, sitemap, or additional features)
- BlogHandler, templates, and CSS design system established as the UI foundation for Phase 3 admin work
- No blockers

## Self-Check: PASSED

- FOUND: web/templates/base.html
- FOUND: web/templates/list.html
- FOUND: web/templates/post.html
- FOUND: web/templates/404.html
- FOUND: web/static/main.css
- FOUND: web/static/main.js
- FOUND: internal/handler/blog/handler.go
- FOUND: internal/handler/blog/list.go
- FOUND: internal/handler/blog/post.go
- FOUND: internal/handler/blog/notfound.go
- FOUND: internal/handler/blog/handler_test.go
- FOUND: cmd/server/main.go
- FOUND commit 3406d2f (Task 1 — templates/CSS/JS)
- FOUND commit 63df1d9 (Task 2 — handlers)
- FOUND commit 3beec1d (Task 3 — main.go)
- FOUND commit 3bc5ce8 (docs — SUMMARY/STATE/REQUIREMENTS)

---
*Phase: 02-public-blog*
*Completed: 2026-03-26*
