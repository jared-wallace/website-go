---
phase: 05-api-images
plan: 02
subsystem: api
tags: [bearer-token, api-push, upsert, front-matter, constant-time-compare]

# Dependency graph
requires:
  - phase: 05-api-images
    plan: 01
    provides: APIToken config field, image upload handler pattern, renderer with RenderWithMeta
  - phase: 03-admin-panel
    provides: Service.Create/Update, RequireSession middleware pattern, mockRepo in write_test.go
provides:
  - RequireAPIToken middleware with constant-time bearer token comparison
  - POST /api/push endpoint on blogMux for CLI push-to-publish
  - UpsertBySlug service method (create draft or update existing, preserving tags)
  - APIHandler package under internal/handler/api
affects: [06-docker]

# Tech tracking
tech-stack:
  added: []
  patterns: [crypto/subtle.ConstantTimeCompare for bearer tokens, io.LimitReader(maxBody+1) for exact size enforcement, interface-based Upserter for handler testability]

key-files:
  created:
    - internal/middleware/apitoken.go
    - internal/middleware/apitoken_test.go
    - internal/handler/api/handler.go
    - internal/handler/api/handler_test.go
    - internal/service/post/upsert_test.go
  modified:
    - internal/service/post/write.go
    - internal/service/post/write_test.go
    - cmd/server/main.go

key-decisions:
  - "Upserter interface in api handler decouples from concrete Service -- tests use lightweight mock"
  - "LimitReader(maxBody+1) then len check gives exact size enforcement without silent truncation"
  - "mockRepo gains findErr pointer field for independent FindBySlug error control in upsert tests"
  - "Route on blogMux (not adminMux) so curl to jared-wallace.com/api/push works without admin host routing"

patterns-established:
  - "API middleware pattern: RequireAPIToken(token) returns func(http.Handler) http.Handler, same shape as RequireSession"
  - "Service interface slicing: handler defines minimal Upserter interface rather than depending on full *Service"

requirements-completed: [ADMIN-09]

# Metrics
duration: 4min
completed: 2026-03-28
---

# Phase 5 Plan 2: API Push Endpoint Summary

**Bearer-token authenticated POST /api/push endpoint with YAML front matter extraction, UpsertBySlug service method, and constant-time token comparison**

## Performance

- **Duration:** 4 min
- **Started:** 2026-03-28T00:34:49Z
- **Completed:** 2026-03-28T00:38:34Z
- **Tasks:** 2
- **Files modified:** 8

## Accomplishments
- Bearer token middleware with crypto/subtle constant-time comparison prevents timing oracle attacks
- UpsertBySlug creates new posts as drafts (published=false), updates existing posts preserving admin-set tags
- API handler extracts slug and title from YAML front matter, falls back to slug as title
- 1 MB body limit enforced without silent truncation
- 14 new tests across middleware, service, and handler layers (all passing)

## Task Commits

Each task was committed atomically:

1. **Task 1: Bearer token middleware, UpsertBySlug service method, and tests** - `f13860f` (test: RED), `933ab83` (feat: GREEN)
2. **Task 2: API push handler, tests, and main.go wiring** - `7149e8c` (test: RED), `d8392b1` (feat: GREEN)

_Note: Both tasks followed TDD with separate RED and GREEN commits._

## Files Created/Modified
- `internal/middleware/apitoken.go` - RequireAPIToken middleware with constant-time bearer token comparison
- `internal/middleware/apitoken_test.go` - 5 tests: valid, missing, invalid, empty bearer, empty config token
- `internal/handler/api/handler.go` - APIHandler with PushPost method, front matter extraction, 1 MB limit
- `internal/handler/api/handler_test.go` - 6 tests: valid token, no token, invalid token, no slug, no title, body too large
- `internal/service/post/write.go` - Added UpsertBySlug method (create draft or update existing)
- `internal/service/post/upsert_test.go` - 3 tests: new post, existing post, find error propagation
- `internal/service/post/write_test.go` - Added findErr field to mockRepo for independent FindBySlug error control
- `cmd/server/main.go` - Wired POST /api/push on blogMux with RequireAPIToken middleware

## Decisions Made
- Used interface slicing (Upserter) in API handler rather than depending on full *Service -- keeps test mocks minimal
- LimitReader reads maxBody+1 bytes then checks length, avoiding silent truncation of exactly-at-limit payloads
- Added findErr pointer field to mockRepo so upsert tests can return ErrNotFound from FindBySlug without affecting Create
- Route placed on blogMux (not adminMux) per plan -- bearer token is the auth gate, not host-based routing

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Fixed LimitReader silent truncation**
- **Found during:** Task 2 (GREEN phase)
- **Issue:** Plan specified `io.LimitReader(r.Body, 1<<20)` but LimitReader silently truncates at limit -- a 1 MB+1 byte body would be read as exactly 1 MB and accepted
- **Fix:** Read `maxBody+1` bytes, then check `len(raw) > maxBody` to detect oversized bodies
- **Files modified:** internal/handler/api/handler.go
- **Verification:** TestPushPost_BodyTooLarge passes with 1 MB+1 byte body rejected at 400
- **Committed in:** d8392b1

---

**Total deviations:** 1 auto-fixed (1 bug)
**Impact on plan:** Fix necessary for correctness. No scope creep.

## Issues Encountered
None beyond the auto-fixed deviation above.

## User Setup Required
None - API_TOKEN environment variable already documented from Plan 01. Set it to enable the push endpoint; leave empty to disable.

## Known Stubs
None - all data paths are fully wired.

## Next Phase Readiness
- Push-to-publish API complete; `curl -H "Authorization: Bearer $TOKEN" -d @post.md https://jared-wallace.com/api/push` is fully functional
- Phase 06 (Docker) should ensure API_TOKEN is included in docker-compose environment variables

## Self-Check: PASSED

All 7 files verified present. All 4 commit hashes verified in git log.

---
*Phase: 05-api-images*
*Completed: 2026-03-28*
