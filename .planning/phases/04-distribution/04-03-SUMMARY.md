---
phase: 04-distribution
plan: 03
subsystem: api
tags: [go, postgres, reactions, javascript, css, sha256, json]

# Dependency graph
requires:
  - phase: 04-02
    provides: OG meta, post detail view (ShowPost) already wired
provides:
  - Thumbs-up reaction system: POST /posts/{slug}/react endpoint
  - Reaction count displayed on post page HTML
  - localStorage guard preventing re-tap on page reload
  - CSS bounce animation on successful reaction
affects: [04-distribution, 06-docker]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - Service pass-through for repo methods the handler needs directly
    - IP hashing via SHA-256 before DB storage (privacy-preserving rate-limit)
    - ON CONFLICT DO NOTHING with RowsAffected() == 0 to detect duplicate without error

key-files:
  created:
    - db/migrations/00004_create_reactions.sql
    - internal/model/reaction.go
    - internal/repository/post/reactions.go
    - internal/handler/blog/react.go
    - internal/handler/blog/react_test.go
  modified:
    - internal/repository/post/repository.go
    - internal/handler/blog/post.go
    - internal/service/post/service.go
    - cmd/server/main.go
    - web/templates/post.html
    - web/static/main.js
    - web/static/main.css
    - internal/handler/blog/handler_test.go
    - internal/service/post/service_test.go
    - internal/service/post/write_test.go
    - internal/handler/admin/handler_test.go
    - internal/handler/admin/actions_test.go
    - internal/handler/admin/editor_test.go

key-decisions:
  - "IP hashed via SHA-256 before storage — never store raw IP; RowsAffected()==0 from ON CONFLICT DO NOTHING detects duplicate without error"
  - "Service exposes AddReaction/CountReactions as thin pass-throughs — keeps handler-to-service dependency pattern consistent"
  - "Reaction JS kept inside the existing IIFE in main.js — no second script block needed"

patterns-established:
  - "Pattern: ON CONFLICT DO NOTHING + RowsAffected for upsert-with-detection in pgx v5"
  - "Pattern: X-Real-IP with RemoteAddr fallback for client IP behind Nginx proxy"

requirements-completed: [BLOG-10]

# Metrics
duration: 4min
completed: 2026-03-27
---

# Phase 4 Plan 03: Thumbs-Up Reaction System Summary

**POST /posts/{slug}/react with SHA-256 IP hashing, ON CONFLICT DO NOTHING deduplication, JSON response, localStorage guard, and CSS bounce animation**

## Performance

- **Duration:** 4 min
- **Started:** 2026-03-27T18:16:29Z
- **Completed:** 2026-03-27T18:19:54Z
- **Tasks:** 2
- **Files modified:** 13

## Accomplishments

- Reactions table migration with unique index on (post_id, ip_hash) enforcing one reaction per IP per post at the DB level
- React handler returning JSON `{count, already_reacted}` with SHA-256 IP privacy hashing and X-Real-IP/RemoteAddr extraction
- Post page HTML includes reaction button with live count, JS fetch + localStorage guard, CSS bounce animation

## Task Commits

Each task was committed atomically:

1. **Task 1: Migration, model, repository interface, and implementation** - `6ef9f43` (feat)
2. **Task 2: React handler, ShowPost count, route, template, JS, and CSS** - `4793ebe` (feat)

**Plan metadata:** (docs commit — see below)

## Files Created/Modified

- `db/migrations/00004_create_reactions.sql` - reactions table, unique index on (post_id, ip_hash), cascade delete
- `internal/model/reaction.go` - Reaction struct
- `internal/repository/post/reactions.go` - AddReaction (ON CONFLICT DO NOTHING) and CountReactions Postgres implementations
- `internal/repository/post/repository.go` - AddReaction and CountReactions added to Repository interface
- `internal/handler/blog/react.go` - React handler: slug lookup, IP hash, AddReaction, CountReactions, JSON response
- `internal/handler/blog/react_test.go` - TestReact, TestReactDuplicate, TestReactNotFound, TestPostReactionCount
- `internal/handler/blog/post.go` - ShowPost fetches and passes ReactionCount to template
- `internal/service/post/service.go` - AddReaction and CountReactions pass-throughs
- `cmd/server/main.go` - POST /posts/{slug}/react route registered
- `web/templates/post.html` - reaction-bar with reaction-btn and reaction-count span
- `web/static/main.js` - fetch + localStorage guard + bounce animation + animationend cleanup
- `web/static/main.css` - .reaction-bar, .reaction-btn, .reacted/:disabled state, reaction-bounce keyframes
- Multiple test files - all mock repositories updated with AddReaction/CountReactions stubs

## Decisions Made

- IP hashed via SHA-256 before storage — never store raw IP; privacy-preserving rate-limit
- RowsAffected() == 0 from ON CONFLICT DO NOTHING cleanly detects duplicate without needing a SELECT
- Service exposes AddReaction/CountReactions as thin pass-throughs to maintain handler-to-service pattern
- Reaction JS placed inside the existing IIFE in main.js — no new script block needed

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] Updated all mock repositories in admin and service/post test files**
- **Found during:** Task 1 (extending Repository interface)
- **Issue:** Adding AddReaction and CountReactions to the Repository interface caused build failures in admin handler tests (actions_test.go, editor_test.go, handler_test.go) and service/post/write_test.go — none of those mocks had the new methods
- **Fix:** Added stub implementations (returning false/0/nil) to all four additional mock repositories
- **Files modified:** internal/handler/admin/handler_test.go, internal/handler/admin/actions_test.go, internal/handler/admin/editor_test.go, internal/service/post/write_test.go
- **Verification:** `go test ./...` all pass
- **Committed in:** 6ef9f43 (Task 1 commit)

---

**Total deviations:** 1 auto-fixed (1 blocking)
**Impact on plan:** Necessary consequence of interface extension — all mocks must implement the full interface. No scope creep.

## Issues Encountered

None beyond the mock update deviation above.

## User Setup Required

None — no external service configuration required. The reactions table migration runs automatically on next container start via goose.

## Next Phase Readiness

- Phase 04 all 3 plans complete — distribution features (RSS, OG meta, reactions) are done
- Ready for Phase 05 or Phase 06 (Docker deployment)
- Note: reactions table will be created on next `goose up` / container start

---
*Phase: 04-distribution*
*Completed: 2026-03-27*
