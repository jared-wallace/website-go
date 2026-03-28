---
phase: 04-distribution
plan: 02
subsystem: social-meta
tags: [og-meta, twitter-card, rss-autodiscovery, templates, go]
dependency_graph:
  requires: []
  provides: [og-meta-tags, twitter-card-tags, rss-autodiscovery-link, og-fallback-image]
  affects: [web/templates/base.html, web/templates/post.html, web/templates/list.html]
tech_stack:
  added: []
  patterns: [html/template block/define override, go image/png for asset generation]
key_files:
  created:
    - web/static/og-fallback.png
  modified:
    - internal/service/post/get.go
    - internal/handler/blog/post.go
    - web/templates/base.html
    - web/templates/post.html
    - web/templates/list.html
    - internal/handler/blog/handler_test.go
decisions:
  - Excerpt(p.Body, 200) added to PostDetail.GetBySlug so OG description is computed at service layer (not template layer)
  - list.html OG meta override has same content as base.html defaults — intentional, makes list page explicit and independent
  - 400x400 solid-color PNG (#2C5F7A) generated via Python3 as ImageMagick not available in build environment
metrics:
  duration: "2m 43s"
  completed: "2026-03-27"
  tasks: 2
  files: 6
---

# Phase 4 Plan 02: Open Graph Meta Tags and RSS Auto-Discovery Summary

OG/Twitter Card meta tags on all public pages, dynamic per-post tags using PostDetail.Excerpt, RSS auto-discovery link, and a 400x400 branded fallback image.

## Tasks Completed

| # | Name | Commit | Files |
|---|------|--------|-------|
| 1 | Add Excerpt to PostDetail and pass from ShowPost | be2efc0 | internal/service/post/get.go, internal/handler/blog/post.go |
| 2 | Add OG meta block, RSS auto-discovery, fallback image, OG meta tests | f92c21f | web/templates/base.html, web/templates/post.html, web/templates/list.html, web/static/og-fallback.png, internal/handler/blog/handler_test.go |

## What Was Built

**PostDetail.Excerpt field** — Added `Excerpt string` to the `PostDetail` struct in `internal/service/post/get.go`. The value is computed using the pre-existing `Excerpt(p.Body, 200)` function, giving a clean 200-char plain-text excerpt for use in OG description tags. The `ShowPost` handler passes this as `"Excerpt"` to the post.html template.

**base.html meta block** — Added RSS auto-discovery link (`<link rel="alternate" type="application/rss+xml" title="The Log" href="/rss">`) and a `{{block "meta" .}}` extensibility block with site-level OG+Twitter defaults. Every page now advertises the RSS feed and has fallback social meta tags.

**post.html OG override** — `{{define "meta"}}` block with dynamic `{{.Post.Title}}`, `{{.Excerpt}}`, `{{.Post.Slug}}` values. Twitter card type `summary`. Canonical URL uses production domain with post slug.

**list.html OG override** — `{{define "meta"}}` block with explicit site-level values. Same content as base.html defaults, but independent — future changes to base.html won't accidentally affect list page.

**og-fallback.png** — 400x400 PNG, solid #2C5F7A (ocean accent color), generated via Python3 `image/png` approach. Served at `/static/og-fallback.png`.

**Test coverage** — `TestPostOGMeta` and `TestListOGMeta` added to `internal/handler/blog/handler_test.go`. Both verify meta tag presence in full handler response bodies using the existing `newTestHandler` + `mockRepository` pattern.

## Verification Results

- `go build ./...` — PASS
- `go test ./internal/handler/blog/... -run "TestPostOGMeta|TestListOGMeta" -v` — PASS (both tests green)
- All existing blog handler tests pass (TestNotFound, TestListPostsEmpty, TestShowPostNotFound)
- `test -f web/static/og-fallback.png` — PASS (400x400 PNG, 8-bit RGB)

## Deviations from Plan

None — plan executed exactly as written.

Note: `go test ./...` shows pre-existing failures in `internal/handler/blog` for `ServeRSS` and `ServeSitemap` — these are TDD red-phase stubs committed by the parallel Plan 01 agent and are expected to fail until Plan 01's green phase completes. Unrelated to this plan.

## Known Stubs

None — all OG meta values are wired from real data sources (post title, service-computed excerpt, post slug). The fallback image is intentionally static for Phase 4 (per D-06; per-post images are Phase 5 scope).

## Self-Check: PASSED

- be2efc0 exists: FOUND
- f92c21f exists: FOUND
- internal/service/post/get.go contains `Excerpt string`: FOUND
- web/templates/base.html contains `block "meta"`: FOUND
- web/templates/post.html contains `define "meta"`: FOUND
- web/templates/list.html contains `define "meta"`: FOUND
- web/static/og-fallback.png exists: FOUND
- handler_test.go contains `TestPostOGMeta`: FOUND
- handler_test.go contains `TestListOGMeta`: FOUND
