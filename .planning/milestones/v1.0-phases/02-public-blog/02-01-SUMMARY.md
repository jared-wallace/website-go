---
phase: 02-public-blog
plan: "01"
subsystem: database
tags: [pgx, postgres, goldmark, html-parser, pagination, toc, excerpt]

# Dependency graph
requires:
  - phase: 01-foundation
    provides: pgxpool.Pool, markdown renderer, goose migrations, model.Post struct
provides:
  - Post repository interface (ListPublished, CountPublished, FindBySlug) backed by pgx
  - Tags migration (00002_add_tags_to_posts.sql)
  - model.Post.Tags field
  - Service layer: ListPublished with pagination math, GetBySlug with enriched PostDetail
  - ReadingTime, ExtractToC, Excerpt, ParseTags utility functions
affects: [02-02-handlers, 02-03-templates, 03-admin]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "Repository interface injected into Service — enables mock testing without DB"
    - "ErrNotFound sentinel exported from repository; re-exported from service for caller convenience"
    - "RenderedHTML cast to template.HTML in service layer, never in template — avoids double-escape"
    - "ToC returns nil when < 3 headings — callers treat nil as no-ToC signal"
    - "Tags stored as comma-separated TEXT, split at service layer not repository"

key-files:
  created:
    - db/migrations/00002_add_tags_to_posts.sql
    - internal/repository/post/repository.go
    - internal/repository/post/queries.go
    - internal/repository/post/repository_test.go
    - internal/service/post/service.go
    - internal/service/post/list.go
    - internal/service/post/get.go
    - internal/service/post/reading_time.go
    - internal/service/post/toc.go
    - internal/service/post/excerpt.go
    - internal/service/post/service_test.go
  modified:
    - internal/model/post.go

key-decisions:
  - "ErrNotFound re-exported from service/post so handlers import one package only"
  - "ExtractToC uses golang.org/x/net/html (already indirect dep) — no new dependency"
  - "Excerpt operates on raw markdown, ReadingTime on raw markdown — not rendered HTML"
  - "PostsPerPage = 10 (D-07); pagination clamps page to [1, totalPages]"

patterns-established:
  - "Repository pattern: interface + unexported struct + New() constructor"
  - "TDD: RED (test file only) → GREEN (implementation) → full test pass before commit"
  - "Service re-exports repository sentinel errors for single-import ergonomics"

requirements-completed: [BLOG-02, BLOG-03, BLOG-04, BLOG-12]

# Metrics
duration: 4min
completed: 2026-03-26
---

# Phase 02 Plan 01: Data and Service Layer Summary

**pgx-backed post repository with pagination, plus service layer computing reading time, ToC extraction, excerpt stripping, and tag parsing — all tested without a database.**

## Performance

- **Duration:** ~4 min
- **Started:** 2026-03-26T13:01:29Z
- **Completed:** 2026-03-26T13:04:47Z
- **Tasks:** 2
- **Files modified:** 12

## Accomplishments

- Tags migration added as a clean goose Up/Down DDL file
- Repository interface with 3 methods backed by parameterized pgx queries (LIMIT/OFFSET, pgx.ErrNoRows wrapping)
- Service layer with pagination math: TotalPages, HasPrev, HasNext, page clamping
- ReadingTime: ceil(words/200), minimum 1
- ExtractToC: golang.org/x/net/html walk collecting h2/h3 with id attrs, nil when < 3 headings
- Excerpt: regex markdown stripping (links, bold, italic, code, headings, HR), rune-aware truncation
- ParseTags: comma-split, trim, filter empty, returns non-nil empty slice
- 17 unit tests all passing with zero DB dependency; 1 sentinel test in repository package

## Task Commits

Each task was committed atomically:

1. **Task 1: Migration, model update, and post repository** - `bc3ca2a` (feat)
2. **Task 2: Post service with pagination, reading time, ToC, and excerpt** - `f81b2b3` (feat)

_Note: Both tasks used TDD (RED → GREEN) within a single commit per task._

## Files Created/Modified

- `db/migrations/00002_add_tags_to_posts.sql` - Tags column DDL (goose Up/Down)
- `internal/model/post.go` - Added Tags string field after RenderedHTML
- `internal/repository/post/repository.go` - Repository interface, ErrNotFound sentinel, New() constructor
- `internal/repository/post/queries.go` - ListPublished (ORDER BY created_at DESC LIMIT/OFFSET), CountPublished, FindBySlug
- `internal/repository/post/repository_test.go` - Sentinel unit test (no build tag); integration test stubs guarded by //go:build integration
- `internal/service/post/service.go` - Service struct with repo injection, New() constructor
- `internal/service/post/list.go` - PostsPerPage=10, ListResult, PostSummary, ListPublished pagination
- `internal/service/post/get.go` - PostDetail, GetBySlug, ErrNotFound re-export
- `internal/service/post/reading_time.go` - ReadingTime(body string) int
- `internal/service/post/toc.go` - ToCEntry, ExtractToC using golang.org/x/net/html
- `internal/service/post/excerpt.go` - Excerpt (regex-strip + truncate), ParseTags
- `internal/service/post/service_test.go` - 17 unit tests covering all public functions

## Decisions Made

- ErrNotFound re-exported from service/post so HTTP handlers only need to import one package
- ExtractToC uses golang.org/x/net/html which was already an indirect dependency — no new module added
- Both Excerpt and ReadingTime operate on raw markdown body (not rendered HTML), consistent with RESEARCH.md Pattern 9
- Integration tests guarded with `//go:build integration` tag; only sentinel test runs in CI without a database

## Deviations from Plan

None — plan executed exactly as written.

## Issues Encountered

None.

## User Setup Required

None — no external service configuration required.

## Next Phase Readiness

- Repository and service layers are fully tested and ready for handler wiring in Plan 02-02
- HTTP handlers can import `internal/service/post` and check `errors.Is(err, postservice.ErrNotFound)` for 404 responses
- Tags migration (00002) must be run against the database before handlers go live

---
*Phase: 02-public-blog*
*Completed: 2026-03-26*
