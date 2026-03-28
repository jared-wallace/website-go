---
phase: 08-about-page
plan: "01"
subsystem: public-blog
tags: [about-page, navigation, content, sitemap, tdd]
dependency_graph:
  requires: []
  provides: [about-page, about-nav-link, about-sitemap-entry]
  affects: [web/templates/base.html, web/static/main.css, cmd/server/main.go]
tech_stack:
  added: [content package with go:embed string directive]
  patterns: [per-page template set, embedded markdown rendering, TDD red-green cycle]
key_files:
  created:
    - content/about.md
    - content/embed.go
    - web/templates/about.html
    - internal/handler/blog/about.go
    - internal/handler/blog/about_test.go
  modified:
    - internal/handler/blog/handler.go
    - internal/handler/blog/sitemap.go
    - web/templates/base.html
    - web/static/main.css
    - cmd/server/main.go
decisions:
  - "Used dedicated content/ package at repo root with go:embed string directive (not web/embed.go) — cleaner separation between web assets and content"
  - "Created renderer per request in AboutPage — simpler than adding renderer field to BlogHandler for a single handler; acceptable at blog scale"
metrics:
  duration: "3min"
  completed: "2026-03-28"
  tasks_completed: 2
  files_modified: 10
---

# Phase 08 Plan 01: About Page Summary

**One-liner:** Embedded markdown about page with go:embed string directive, goldmark rendering pipeline, and nautical nav link.

## Tasks Completed

| Task | Description | Commit | Files |
|------|-------------|--------|-------|
| 1 (RED) | Failing tests for AboutPage handler | b7332c4 | about_test.go |
| 1 (GREEN) | about.md content, about.html template, content embed.go, AboutPage handler | 9255932 | 5 files |
| 2 | Nav link, CSS styles, sitemap entry, route registration | e52ab77 | 4 files |

## What Was Built

A fully-integrated /about page for The Wild Meridian blog:

- **`content/about.md`** — static markdown content with personal bio, blog purpose, and "Thanks for reading" sign-off
- **`content/embed.go`** — dedicated `content` package using `//go:embed about.md` as a `string` variable, enabling direct pass-through to `renderer.Render()`
- **`web/templates/about.html`** — per-page template set following the same pattern as post.html; renders `{{.RenderedHTML}}` inside `.about-page` / `.post-body` container
- **`internal/handler/blog/about.go`** — `AboutPage` handler method; renders embedded markdown through goldmark+bluemonday pipeline
- **`web/templates/base.html`** — About nav link `<a href="/about" class="nav-link">About</a>` inserted between `.site-tagline` and `.dark-toggle`
- **`web/static/main.css`** — `.nav-link` (14px, 44px touch target, muted color, hover accent, focus-visible outline) and `.about-title` (36px Playfair Display, weight 700, line-height 1.15) rules
- **`internal/handler/blog/sitemap.go`** — /about added to `buildSitemap` with `monthly` changefreq and `0.5` priority; capacity bumped from `+1` to `+2`
- **`cmd/server/main.go`** — `GET /about` route registered on `blogMux` before catch-all 404

## Verification

All 6 plan-level checks pass:

1. `go build ./cmd/server/` — exits 0
2. `go test ./internal/handler/blog/ -v` — 18/18 tests pass including 4 new about tests
3. `grep -c 'nav-link' web/templates/base.html` — 1
4. `grep -c 'about-title' web/templates/about.html` — 1
5. `grep -c '/about' cmd/server/main.go` — 1
6. `grep -c '/about' internal/handler/blog/sitemap.go` — 1

## Decisions Made

1. **Dedicated `content/` package at repo root** rather than adding another embed variable to `web/embed.go`. Keeps content separate from web assets; the `//go:embed about.md` string directive is cleaner than `embed.FS` for a single file.

2. **Renderer created per-request in `AboutPage`** rather than adding a `md *markdown.Renderer` field to `BlogHandler`. The plan mentioned this as an optimization opportunity, but for a single static page with negligible traffic, creating a new renderer per request is simpler and avoids changing the `New()` signature and all call sites. Deferred as a possible future refactor if more embedded-markdown pages are added.

## Deviations from Plan

None — plan executed exactly as written. The plan's note about the renderer field was presented as an option; the simpler per-request approach was chosen deliberately and within the plan's described scope.

## Known Stubs

None. All content is real, wired, and renders correctly.

## Self-Check: PASSED

- `content/about.md` — FOUND
- `content/embed.go` — FOUND
- `web/templates/about.html` — FOUND
- `internal/handler/blog/about.go` — FOUND
- `internal/handler/blog/about_test.go` — FOUND
- Commit b7332c4 — FOUND (test RED)
- Commit 9255932 — FOUND (feat GREEN)
- Commit e52ab77 — FOUND (feat Task 2)
