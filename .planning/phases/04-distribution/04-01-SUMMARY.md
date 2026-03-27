---
phase: 04-distribution
plan: "01"
subsystem: blog-handler
tags: [rss, sitemap, robots, xml, feed, seo]
dependency_graph:
  requires: []
  provides: [rss-feed, sitemap, robots-txt]
  affects: [cmd/server/main.go, internal/handler/blog, internal/service/post]
tech_stack:
  added: []
  patterns: [encoding/xml-struct-marshal, CDATA-wrapper, xml-handler-method]
key_files:
  created:
    - internal/service/post/feed.go
    - internal/handler/blog/rss.go
    - internal/handler/blog/sitemap.go
    - internal/handler/blog/rss_test.go
    - internal/handler/blog/sitemap_test.go
  modified:
    - cmd/server/main.go
decisions:
  - "baseURL constant defined in rss.go and referenced in sitemap.go — same package, no import needed"
  - "CDATA type implements xml.Marshaler to wrap RenderedHTML in CDATA sections, preventing double-escaping"
  - "ListSlugsForSitemap uses repo.ListPublished(ctx, 10000, 0) — no new repo interface method needed at blog scale"
  - "managingEditor uses jaredwallace@jared-wallace.com (Jared Wallace) per D-02 and RSS 2.0 spec"
metrics:
  duration: "~8 minutes"
  completed: "2026-03-27"
  tasks_completed: 2
  files_changed: 6
requirements_satisfied: [BLOG-09, BLOG-07]
---

# Phase 4 Plan 1: RSS Feed, Sitemap, and robots.txt Summary

**One-liner:** RSS 2.0 feed with CDATA HTML content, Sitemap 0.9 XML, and robots.txt served from three new BlogHandler methods wired in blogMux.

## What Was Built

Added machine-readable distribution endpoints to the public blog:

1. **`GET /rss`** — RSS 2.0 feed of the 25 most recent published posts. Full rendered HTML content in `<![CDATA[...]]>` sections prevents double-escaping. Channel metadata uses "The Log" / "dispatches from the deep end". Per-item `<category>` elements map post tags per D-03.

2. **`GET /sitemap.xml`** — Sitemap 0.9 XML listing homepage + all published post URLs as absolute `https://jared-wallace.com/...` URLs.

3. **`GET /robots.txt`** — Plain text with `Sitemap:` directive pointing to the sitemap endpoint per D-15.

## Tasks Completed

| Task | Name | Commit | Files |
|------|------|--------|-------|
| 1 (RED) | Failing tests | 0f40da0 | rss_test.go, sitemap_test.go |
| 1 (GREEN) | Service + handler implementation | 6e8969c | feed.go, rss.go, sitemap.go |
| 2 | Wire routes in main.go | 6e7ada2 | cmd/server/main.go |

## Deviations from Plan

None — plan executed exactly as written.

## Known Stubs

None — all endpoints return real data from the service/repo layer. Mock in tests is intentional and correct.

## Self-Check: PASSED

- All 7 files verified present on disk
- All 3 task commits verified in git log (0f40da0, 6e8969c, 6e7ada2)
- Full test suite green: `go test ./...` exits 0, no failures
