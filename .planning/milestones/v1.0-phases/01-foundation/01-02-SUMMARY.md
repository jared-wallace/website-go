---
phase: 01-foundation
plan: 02
subsystem: markdown
tags: [goldmark, bluemonday, markdown, html-sanitization, xss-protection, syntax-highlighting, gfm, yaml-frontmatter]

# Dependency graph
requires: []
provides:
  - "Goldmark + bluemonday rendering pipeline (internal/markdown package)"
  - "Renderer struct with Render() and RenderWithMeta() methods"
  - "XSS sanitization via bluemonday UGCPolicy downstream of goldmark"
  - "YAML front matter extraction for post metadata"
affects: [phase-02-content, phase-03-admin, all-phases-using-markdown]

# Tech tracking
tech-stack:
  added:
    - "github.com/yuin/goldmark v1.8.2"
    - "github.com/yuin/goldmark-meta v1.1.0"
    - "github.com/yuin/goldmark-highlighting/v2"
    - "github.com/alecthomas/chroma/v2"
    - "github.com/microcosm-cc/bluemonday v1.0.27"
  patterns:
    - "goldmark.Convert() THEN bluemonday.Sanitize() — order is critical for XSS safety"
    - "html.WithUnsafe() on goldmark is safe only because bluemonday runs downstream"
    - "Renderer struct wraps both goldmark.Markdown and bluemonday.Policy for single call surface"
    - "template.HTML return type signals to Go's html/template that output is pre-sanitized"
    - "AllowAttrs(class).OnElements(code,span,pre) added to preserve syntax highlighting CSS classes"

key-files:
  created:
    - "internal/markdown/renderer.go"
    - "internal/markdown/renderer_test.go"
    - "go.mod"
    - "go.sum"
  modified: []

key-decisions:
  - "Pipeline order is goldmark → bluemonday (never reverse); html.WithUnsafe() is safe only because bluemonday sanitizes downstream"
  - "bluemonday AllowAttrs(class).OnElements(code,span,pre) added to preserve syntax highlighting CSS classes through sanitization"
  - "module name: github.com/jared-wallace/website-go"

patterns-established:
  - "TDD: write all tests first (RED), then implement (GREEN), no behavior changed in refactor"
  - "markdown.NewRenderer() as constructor — callers never build goldmark or bluemonday directly"
  - "Render() for HTML-only output, RenderWithMeta() when YAML front matter is needed"

requirements-completed: [FOUND-01]

# Metrics
duration: 2min
completed: 2026-03-26
---

# Phase 1 Plan 02: Markdown Rendering Pipeline Summary

**Goldmark + bluemonday pipeline converting markdown to sanitized template.HTML with GFM tables, strikethrough, linkify, syntax highlighting (monokai/chroma), and YAML front matter extraction — all 10 XSS and correctness tests green**

## Performance

- **Duration:** 2 min
- **Started:** 2026-03-26T11:40:14Z
- **Completed:** 2026-03-26T11:41:56Z
- **Tasks:** 1 (TDD: RED + GREEN)
- **Files modified:** 4 (go.mod, go.sum, renderer.go, renderer_test.go)

## Accomplishments

- Markdown rendering pipeline with all four required goldmark extensions (GFM, linkify, syntax highlighting, YAML front matter)
- bluemonday UGCPolicy sanitization downstream of goldmark — strips script, iframe, event handlers
- 10 tests covering every specified behavior: XSS stripping, GFM tables, strikethrough, code fences, link auto-detection, front matter extraction, iframe stripping, event handler stripping, empty input
- `go build ./...` clean after adding 5 new dependencies

## Task Commits

Each task was committed atomically:

1. **Task 1: Markdown renderer with goldmark extensions and bluemonday sanitization** - `1b09a23` (feat)

**Plan metadata:** (pending docs commit)

_Note: TDD task — RED (no source file) confirmed failing, then GREEN implemented in single commit_

## Files Created/Modified

- `internal/markdown/renderer.go` - Renderer struct, NewRenderer(), Render(), RenderWithMeta()
- `internal/markdown/renderer_test.go` - 10 behavior tests including critical XSS gate
- `go.mod` - Module initialized as github.com/jared-wallace/website-go, go 1.23.3
- `go.sum` - Dependency checksums for goldmark, bluemonday, chroma, and transitive deps

## Decisions Made

- **Pipeline order locked:** goldmark.Convert() then bluemonday.SanitizeBytes() — html.WithUnsafe() on goldmark is safe because bluemonday cleans up downstream. This order is enforced by test design (XSS test catches any accidental reversal).
- **AllowAttrs("class") on code/span/pre:** bluemonday UGCPolicy strips class attributes by default, which would destroy syntax highlighting CSS classes from chroma. Added a policy extension to preserve them.
- **Module name:** `github.com/jared-wallace/website-go` — standard Go module naming for the repo.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 2 - Missing Critical] Added AllowAttrs("class") to bluemonday policy for syntax highlighting**
- **Found during:** Task 1 (implementation review of plan notes)
- **Issue:** Plan explicitly warned "If syntax highlighting produces class attributes on code spans, confirm bluemonday UGCPolicy preserves them -- if not, add `bm.AllowAttrs("class").OnElements("code", "span", "pre")`"
- **Fix:** Applied the prescribed fix preemptively — added `bm.AllowAttrs("class").OnElements("code", "span", "pre")` to the policy
- **Files modified:** internal/markdown/renderer.go
- **Verification:** Code fence test passes, syntax highlighting class attributes preserved
- **Committed in:** 1b09a23 (Task 1 commit)

---

**Total deviations:** 1 auto-applied (plan-prescribed policy extension for syntax highlighting class preservation)
**Impact on plan:** Required for syntax highlighting to work correctly in Phase 2. No scope creep.

## Issues Encountered

None — plan executed cleanly. TDD RED phase confirmed by build failure (no source file), GREEN phase passed all 10 tests on first implementation.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- `internal/markdown` package is production-ready: tested, sanitized, typed as `template.HTML`
- Ready for Phase 1 Plan 03 (project scaffold, Makefile, CI)
- Any Phase 2+ handler that renders a post can call `renderer.Render(post.Body)` directly

## Self-Check: PASSED

- internal/markdown/renderer.go: FOUND
- internal/markdown/renderer_test.go: FOUND
- .planning/phases/01-foundation/01-02-SUMMARY.md: FOUND
- Commit 1b09a23: FOUND
- go test ./internal/markdown/... -v -count=1: ALL 10 TESTS PASS
