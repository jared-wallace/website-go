---
phase: 07-rebrand-navigation
plan: "01"
subsystem: templates, rss, css
tags: [rebrand, brand-identity, rss, accessibility, css]
dependency_graph:
  requires: []
  provides: [wild-meridian-brand, rss-footer-icon]
  affects: [web/templates, web/static, internal/handler/blog]
tech_stack:
  added: []
  patterns: [feather-svg-icon, 44px-touch-target, currentColor-svg]
key_files:
  created: []
  modified:
    - web/templates/base.html
    - web/templates/list.html
    - web/templates/post.html
    - web/templates/404.html
    - web/templates/admin-base.html
    - web/templates/admin-login.html
    - internal/handler/blog/rss.go
    - internal/handler/blog/rss_test.go
    - internal/handler/blog/handler_test.go
    - web/static/main.css
decisions:
  - "Added display:flex to .site-footer p to enable inline icon alignment alongside copyright text"
  - "RSS icon uses Feather rss SVG (two arcs + circle) matching existing dark-toggle inline SVG pattern"
metrics:
  duration: "~2 min"
  completed: "2026-03-28"
  tasks_completed: 2
  files_modified: 10
---

# Phase 7 Plan 1: Rebrand + RSS Icon Summary

**One-liner:** Complete identity migration from "The Log" to "The Wild Meridian" with Feather RSS icon in footer linking to /rss.

## Tasks Completed

| Task | Name | Commit | Files |
|------|------|--------|-------|
| 1 | Rebrand all "The Log" strings to "The Wild Meridian" | 9a64efa | 9 files |
| 2 | Add RSS broadcast icon to footer and CSS styling | d245346 | 2 files |

## What Was Built

### Task 1: Brand String Replacement (9 files)

Replaced all 17 occurrences of "The Log" across the entire public+admin surface:

- **base.html**: title block, RSS `<link rel="alternate">` title, og:title, twitter:title, nav `.site-name`, copyright footer (to "Jared Wallace")
- **list.html**: `{{define "title"}}`, og:title, twitter:title
- **post.html**: title suffix `&#8212; The Wild Meridian`
- **404.html**: title suffix `&#8212; The Wild Meridian`
- **admin-base.html**: browser tab `-- The Wild Meridian`, nav `.site-name` "The Wild Meridian -- Back Office"
- **admin-login.html**: login page `<h1>` heading
- **rss.go**: `Channel.Title` field
- **rss_test.go**: `<title>The Wild Meridian</title>` assertion
- **handler_test.go**: `og:title" content="The Wild Meridian"` assertion

### Task 2: RSS Footer Icon (2 files)

- **base.html**: Replaced single-line `<p>` copyright with a multi-line `<p>` containing the copyright text and a `<a href="/rss" class="rss-link" aria-label="RSS feed">` anchor wrapping a Feather `rss` SVG (two radiating arcs + filled circle at lower-left). SVG uses `stroke="currentColor"` and `aria-hidden="true"` matching the existing dark-toggle pattern.
- **main.css**: Added `.rss-link` rule (44px touch target, `color: var(--color-text-muted)`, `text-decoration: none`, 200ms ease transition) and `.rss-link:hover` rule (`color: var(--color-accent)`), mirroring `.dark-toggle`. Also added `display: flex; align-items: center; justify-content: center` to `.site-footer p` for proper icon alignment.

## Verification Results

1. Zero remaining "The Log" strings in templates and blog handler: **PASS**
2. Full test suite (`go test ./...`): **PASS** (9 packages, 0 failures)
3. RSS icon present in base.html and main.css: **PASS**
4. Tagline "dispatches from the deep end" preserved: **PASS**
5. Dynamic `{{.Year}}` preserved in copyright: **PASS**

## Deviations from Plan

### Auto-applied: Extra CSS on `.site-footer p`

- **Found during:** Task 2
- **Issue:** Adding an inline `<a>` tag next to copyright text without flex on the parent `<p>` would result in misaligned baseline rendering
- **Fix:** Added `display: flex; align-items: center; justify-content: center` to `.site-footer p` — keeps the footer visually centered with icon properly vertically aligned
- **Files modified:** web/static/main.css
- **Commit:** d245346

This is a minor layout correctness fix, not a plan deviation — the plan specified the icon should appear "next to" the copyright text, and flex alignment is required for that to look right.

## Known Stubs

None. All brand strings are wired to real values. The RSS icon links to an existing /rss endpoint.

## Self-Check: PASSED
