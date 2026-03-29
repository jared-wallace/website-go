---
phase: quick
plan: 260328-tea
subsystem: frontend/templates
tags: [bug-fix, templates, css, tests]
dependency_graph:
  requires: []
  provides: [clean-about-page-single-h1]
  affects: [web/templates/about.html, web/static/main.css, internal/handler/blog/about_test.go]
tech_stack:
  added: []
  patterns: [markdown-as-single-source-of-truth-for-page-headings]
key_files:
  modified:
    - web/templates/about.html
    - web/static/main.css
    - internal/handler/blog/about_test.go
decisions:
  - Markdown h1 is the authoritative heading source; templates must not duplicate it
metrics:
  duration: 3min
  completed: 2026-03-28
---

# Quick Task 260328-tea: Fix Duplicate About Heading Summary

**One-liner:** Removed template-hardcoded `<h1 class="about-title">About</h1>` and orphaned `.about-title` CSS so the goldmark-rendered markdown `# About` is the sole heading.

## Tasks Completed

| # | Task | Commit | Files |
|---|------|--------|-------|
| 1 | Remove duplicate heading, stale hr, orphaned CSS, update test | f672102 | about.html, main.css, about_test.go |

## What Was Done

The about page had two `<h1>` elements: one hardcoded in the template and one rendered from the markdown's `# About` front matter. The fix:

1. **web/templates/about.html** — Deleted `<h1 class="about-title">About</h1>` and `<hr class="rope-divider">`. The content block now renders only `{{.RenderedHTML}}`.
2. **web/static/main.css** — Removed the `.about-title` rule block (7 lines including the `/* --- About Page --- */` section comment, which was left empty).
3. **internal/handler/blog/about_test.go** — Replaced the `"body contains about-title CSS class"` subtest with `"body contains About heading from markdown"`, asserting `<h1>About</h1>` is present in the rendered output.

## Verification

```
grep -c '<h1' web/templates/about.html      # 0 (pass)
grep -c 'rope-divider' web/templates/about.html  # 0 (pass)
grep -c 'about-title' web/static/main.css   # 0 (pass)
go test ./internal/handler/blog/ -run TestAboutPage -v  # all 6 subtests PASS
```

## Deviations from Plan

None — plan executed exactly as written.

## Known Stubs

None.

## Self-Check: PASSED

- `web/templates/about.html` — exists, contains 0 `<h1` tags
- `web/static/main.css` — exists, contains 0 `.about-title` references
- `internal/handler/blog/about_test.go` — exists, subtest renamed
- Commit f672102 — verified in git log
