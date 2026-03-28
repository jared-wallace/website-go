---
gsd_state_version: 1.0
milestone: v1.0
milestone_name: milestone
status: verifying
stopped_at: Completed 06-01-PLAN.md
last_updated: "2026-03-28T02:07:48.490Z"
last_activity: 2026-03-28
progress:
  total_phases: 6
  completed_phases: 4
  total_plans: 17
  completed_plans: 14
  percent: 80
---

# Project State

## Project Reference

See: .planning/PROJECT.md (updated 2026-03-26)

**Core value:** A reader visits jared-wallace.com and reads well-rendered markdown blog posts with images in a distinctive, memorable design.
**Current focus:** Phase 04 — distribution

## Current Position

Phase: 06 (docker-deployment) — EXECUTING
Plan: 2 of 2
Status: Executing Phase 06
Last activity: 2026-03-28 -- Completed 06-01-PLAN.md

Progress: [████████░░] 82%

## Performance Metrics

**Velocity:**

- Total plans completed: 0
- Average duration: -
- Total execution time: 0 hours

**By Phase:**

| Phase | Plans | Total | Avg/Plan |
|-------|-------|-------|----------|
| - | - | - | - |

**Recent Trend:**

- Last 5 plans: -
- Trend: -

*Updated after each plan completion*
| Phase 01-foundation P01 | 4 | 2 tasks | 13 files |
| Phase 01-foundation P02 | 2 | 1 tasks | 4 files |
| Phase 01-foundation P03 | 3 | 2 tasks | 7 files |
| Phase 02-public-blog P01 | 4min | 2 tasks | 12 files |
| Phase 02-public-blog P02 | 32min | 3 tasks | 12 files |
| Phase 03-admin-panel P01 | 6min | 2 tasks | 18 files |
| Phase 03-admin-panel P03 | 10 | 2 tasks | 6 files |
| Phase 03-admin-panel P04 | 7min | 2 tasks | 7 files |
| Phase 04-distribution P01 | 8min | 2 tasks | 6 files |
| Phase 04-distribution P02 | 3min | 2 tasks | 6 files |
| Phase 04-distribution P03 | 4min | 2 tasks | 13 files |
| Phase 05-api-images P01 | 4min | 2 tasks | 7 files |
| Phase 05 P02 | 4min | 2 tasks | 8 files |
| Phase 06 P01 | 2min | 2 tasks | 4 files |

## Accumulated Context

### Decisions

Decisions are logged in PROJECT.md Key Decisions table.
Recent decisions affecting current work:

- [Pre-planning]: Use native pgx/v5 (not sqlx); use goose v3 (not golang-migrate) — per STACK.md authority over ARCHITECTURE.md comments
- [Pre-planning]: bluemonday is a required dependency (security-critical); must be added before Phase 1 implementation starts
- [Pre-planning]: HTMX vs vanilla JS for split-pane preview is an open question — decide during Phase 3 planning
- [Pre-planning]: Postgres sessions table migration must exist before Phase 3 (admin) work begins
- [Phase 01-01]: go.mod uses 'go 1.26' without toolchain directive — avoids strict enforcement on local Go 1.23 while targeting Docker build environment
- [Phase 01-01]: Migrations embedded in db/migrations/embed.go sibling package — go:embed forbids '..' path components so embed must live adjacent to SQL files
- [Phase 01-02]: goldmark → bluemonday pipeline order locked; html.WithUnsafe() safe only with bluemonday downstream
- [Phase 01-02]: bluemonday UGCPolicy extended with AllowAttrs("class") for syntax highlighting CSS classes
- [Phase 01-foundation]: golangci-lint v1.61.0 pinned in CI; local lint fails on Go 1.23 due to goose v3.27.0 requiring Go 1.25 — CI uses Go 1.26 so CI lint passes
- [Phase 01-foundation]: Single CI job (lint->test->build) to conserve free-tier GHA minutes per D-10
- [Phase 02-public-blog]: ErrNotFound re-exported from service/post so HTTP handlers import one package only
- [Phase 02-public-blog]: ExtractToC uses golang.org/x/net/html (already indirect dep) — no new dependency added
- [Phase 02-public-blog]: Per-page template sets in html/template: parse base.html + page.html separately per page to avoid block name collisions across pages
- [Phase 02-public-blog]: ExecuteTemplate calls 'base' named template after parsing base+page pair into each template set
- [Phase 03-admin-panel]: Service.New accepts Renderer interface (not *markdown.Renderer) — mock-friendly without importing markdown package in tests
- [Phase 03-admin-panel]: pgxstore re-added as direct dep in Plan 02 when session store is wired — go mod tidy correctly removes it when unused
- [Phase 03-admin-panel]: In-memory filter over ListAll chosen over separate DB queries per tab — blog scale makes this simpler
- [Phase 03-admin-panel]: successRepo separate from mockRepository — action tests need write ops to succeed; auth tests need them to fail
- [Phase 03-admin-panel]: Renderer interface added to admin handler package so tests pass mocks without importing markdown — aligns with service layer pattern
- [Phase 03-admin-panel]: Preview endpoint writes goldmark+bluemonday output directly to ResponseWriter (no template) to prevent double-encoding
- [Phase 03-admin-panel]: postView struct in EditPost exposes RenderedHTML as template.HTML type to prevent html/template from escaping pre-rendered HTML
- [Phase 04-distribution]: CDATA type implements xml.Marshaler to prevent double-escaping of HTML in RSS descriptions
- [Phase 04-distribution]: managingEditor uses jaredwallace@jared-wallace.com (Jared Wallace) per D-02 and RSS 2.0 spec
- [Phase 04-02]: Excerpt(p.Body, 200) added to PostDetail.GetBySlug — OG description computed at service layer, not template layer
- [Phase 04-02]: list.html OG meta override is explicit duplicate of base defaults — intentional independence from future base.html changes
- [Phase 04-distribution]: IP hashed via SHA-256 before storage in reactions table — never store raw IP; RowsAffected()==0 from ON CONFLICT DO NOTHING detects duplicate without error
- [Phase 04-distribution]: Reaction JS kept inside existing IIFE in main.js; Service exposes AddReaction/CountReactions as thin pass-throughs to maintain handler-to-service dependency pattern
- [Phase 05-01]: MaxBytesReader used instead of ParseMultipartForm maxMemory for hard 5MB upload limit enforcement
- [Phase 05-01]: APIToken config field added alongside ImageDir to avoid re-modifying config.go in Plan 02
- [Phase 05-01]: SetImageDir test helper on AdminHandler for t.TempDir injection in upload tests
- [Phase 05]: Upserter interface in api handler decouples from concrete Service for testability
- [Phase 05]: LimitReader reads maxBody+1 bytes then checks length for exact size enforcement without silent truncation
- [Phase 05]: API push route on blogMux (not adminMux) -- bearer token is the auth gate, not host routing
- [Phase 06]: 9.9MB final image via alpine:3.21 + stripped static binary -- well under 30MB target
- [Phase 06]: /health returns bare 200 with no DB ping -- Postgres liveness handled by compose pg_isready

### Pending Todos

None yet.

### Blockers/Concerns

- Phase 6 (Docker): Postgres EBS bind-mount requires `chown 999:999` on host before first run — document in Makefile deploy target (Pitfall 3)
- Phase 6 (Docker): ASG must have `max_size = 1` and `delete_on_termination = false` before any production data is written (Pitfall 8)
- Phase 3 (Admin): Go 1.26 stdlib may include CrossOriginProtection — confirm before reaching for gorilla/csrf during Phase 3 planning

## Session Continuity

Last session: 2026-03-28T02:07:48.487Z
Stopped at: Completed 06-01-PLAN.md
Resume file: None
