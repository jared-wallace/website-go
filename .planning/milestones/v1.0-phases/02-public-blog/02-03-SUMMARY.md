---
phase: 02-public-blog
plan: 03
status: complete
started: 2026-03-26T21:00:00Z
completed: 2026-03-27T02:05:00Z
---

# Plan 02-03: Visual Verification — Summary

## What Happened

Human visual verification of the complete public blog against the UI-SPEC design contract.

## Issues Found & Fixed

1. **Static file serving broken** — `embed.FS` path mismatch caused `main.css` and `main.js` to return 404 (browser showed `NS_ERROR_CORRUPTED_CONTENT`). Fixed with `fs.Sub(web.Static, "static")` in `cmd/server/main.go`.

2. **Table of Contents not rendering** — `ExtractToC` required `id` attributes on headings, but rendered HTML lacked them. Added `InjectHeadingIDs` to slugify heading text into `id` attributes, and updated `ExtractToC` to fall back to generated slugs.

## Verification Results

| Item | Status |
|------|--------|
| Home page card grid (2-col desktop) | ✓ |
| Card hover lift effect | ✓ |
| Pagination (12 posts, 2 pages) | ✓ |
| Single post at /posts/{slug} | ✓ |
| ToC on posts with 3+ headings | ✓ (after fix) |
| ToC collapse/expand toggle | ✓ |
| Dark mode toggle + persistence | ✓ |
| No-flash on dark mode reload | ✓ |
| 404 "Lost at Sea" themed page | ✓ |
| Playfair Display + Lora typography | ✓ (after static fix) |
| Mobile responsive (not fully tested — localhost limitation) | ○ Deferred |

## Key Files

### Modified
- `cmd/server/main.go` — `fs.Sub` fix for static file serving
- `internal/service/post/toc.go` — `InjectHeadingIDs`, `slugify`, updated `ExtractToC`
- `internal/service/post/get.go` — wire `InjectHeadingIDs` into `GetBySlug`

## Deviations

- Mobile verification deferred — running on localhost makes device testing impractical. Will verify after Docker deployment (Phase 6).

## Self-Check: PASSED
