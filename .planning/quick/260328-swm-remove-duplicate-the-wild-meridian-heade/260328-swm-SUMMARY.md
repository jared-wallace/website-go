---
phase: quick
plan: 260328-swm
subsystem: frontend
tags: [template, css, cleanup, polish]
dependency_graph:
  requires: []
  provides: [clean-home-page-header]
  affects: [web/templates/list.html, web/static/main.css, internal/handler/blog/handler_test.go]
tech_stack:
  added: []
  patterns: []
key_files:
  modified:
    - web/templates/list.html
    - web/static/main.css
    - internal/handler/blog/handler_test.go
decisions:
  - Removed TestListHero test entirely — test was asserting the presence of intentionally deleted markup, not a regression
metrics:
  duration: 3min
  completed_date: "2026-03-28"
  tasks_completed: 1
  files_changed: 3
---

# Quick Task 260328-swm: Remove Duplicate The Wild Meridian Header — Summary

**One-liner:** Deleted the redundant `.list-hero` div, its CSS rules, and the test asserting them so the site name appears exactly once via the nav bar.

## What Was Done

The home page (`/`) was rendering "The Wild Meridian" and "dispatches from the deep end" twice: once in the nav bar (from `base.html`) and again in a `.list-hero` div at the top of the post list in `list.html`. This task removed the redundant second instance entirely.

## Tasks Completed

| Task | Description | Commit |
|------|-------------|--------|
| 1 | Remove hero div, orphaned CSS, and TestListHero | 49dcc94 |

## Changes Made

**web/templates/list.html**
- Removed lines 19–22: the entire `.list-hero` div containing the `h1` and `p` tags

**web/static/main.css**
- Removed three rule blocks: `.list-hero` (margin), `.list-hero-title` (font/size), `.list-hero-tagline` (italic muted text)
- Net: 21 lines of CSS deleted, zero new classes added

**internal/handler/blog/handler_test.go**
- Removed `TestListHero` function (lines 307–327) — it tested for markup that was intentionally removed

## Verification

- `grep -r "list-hero" web/ internal/` returns no matches (CLEAN)
- `go test ./internal/handler/blog/...` — all 14 remaining tests pass

## Deviations from Plan

None — plan executed exactly as written.

## Known Stubs

None.

## Self-Check: PASSED

- web/templates/list.html: no .list-hero div present
- web/static/main.css: no .list-hero* rules present
- internal/handler/blog/handler_test.go: no TestListHero function present
- Commit 49dcc94: verified in git log
- All tests: PASS
