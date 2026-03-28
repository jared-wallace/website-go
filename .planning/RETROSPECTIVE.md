# Project Retrospective

*A living document updated after each milestone. Lessons feed forward into future planning.*

## Milestone: v1.0 -- MVP

**Shipped:** 2026-03-28
**Phases:** 6 | **Plans:** 17 | **Commits:** 125

### What Was Built
- Full-stack Go blog platform with public reader experience and admin panel
- Nautical "weathered beach bar" design system with dark mode
- Content pipeline: markdown editor, API push, RSS, sitemap, OG tags, reactions
- Docker deployment with Postgres sidecar and EBS persistence

### What Worked
- Per-page template sets in html/template avoided block name collision bugs that would have been painful to debug
- goldmark + bluemonday pipeline order (render unsafe, then sanitize) gave both full rendering and XSS protection
- TDD-style plan execution (failing tests first, then implementation) caught integration issues early
- Minimal dependency philosophy paid off -- stdlib ServeMux, slog, encoding/xml all sufficient at blog scale
- pgx/v5 native interface over database/sql adapter avoided feature loss and was cleaner to use

### What Was Inefficient
- Some SUMMARY.md files had empty one-liner frontmatter, causing garbage in automated milestone extraction
- ROADMAP.md progress table became stale (showed Phase 2 at 1/3 and Phase 6 at 0/2 when both were complete)
- Phase 3 missing two SUMMARY.md files (03-03, 03-04) despite plans being marked complete
- Nyquist validation files created but never executed -- overhead without value

### Patterns Established
- Host-based routing via stdlib ServeMux for admin subdomain separation
- Renderer interface pattern for test mocks (avoid importing markdown package in handler tests)
- Service layer as the single source of rendered HTML (render on write, not on read)
- SHA-256 IP hashing for privacy-preserving deduplication (never store raw IPs)
- CDATA xml.Marshaler for RSS HTML content (prevent double-escaping)

### Key Lessons
1. Keep SUMMARY.md frontmatter populated -- automated tooling depends on it for milestone completion
2. Update ROADMAP.md progress table after each plan, not just at phase boundaries
3. Human verification items accumulate fast when building without a running server -- plan a deploy-and-verify phase
4. In-memory filtering beats separate DB queries at blog scale -- resist premature optimization

### Cost Observations
- Model mix: primarily opus for planning/execution, sonnet for verification/integration checks
- Notable: 6 phases completed in ~3 days of development time

## Milestone: v1.1 -- The Wild Meridian

**Shipped:** 2026-03-28
**Phases:** 2 | **Plans:** 2 | **Timeline:** ~1 hour

### What Was Built
- Full rebrand from "The Log" to "The Wild Meridian" across all public surfaces (header, RSS, OG, sitemap, admin)
- Discreet RSS broadcast icon in footer with 44px touch target and autodiscovery link tag
- /about page with embedded markdown content, nautical template, nav link, and sitemap entry

### What Worked
- Small, focused milestone shipped fast -- 2 phases, 2 plans, done in an hour
- go:embed string directive for about.md (dedicated content/ package) was cleaner than embedding under web/
- Phase 8 executor agent completed all work autonomously with zero intervention
- Verification caught the stale REQUIREMENTS.md traceability table (BRAND/NAV still "Pending" despite Phase 7 passing)

### What Was Inefficient
- SUMMARY.md one-liner frontmatter still empty -- same issue as v1.0, automated extraction returned garbage
- ROADMAP.md progress table showed Phase 7 as "0/1 In progress" and Phase 8 as "0/1 Not started" despite both being complete on disk
- REQUIREMENTS.md traceability not auto-updated for Phase 7 requirements (BRAND-01-04, NAV-01-02 stayed "Pending")

### Patterns Established
- Dedicated content/ package with go:embed for static pages -- scales to future static pages without touching web/
- .nav-link CSS class pattern with 44px touch target, hover/focus-visible states -- reusable for future nav items

### Key Lessons
1. Milestone tooling (`summary-extract`) needs a fallback when one-liner frontmatter is empty -- extract from SUMMARY body instead
2. For tiny milestones (2-3 phases), skip the audit -- overhead exceeds value when everything fits in one context window
3. Progress table staleness is a recurring issue -- consider making gsd-tools update it atomically with plan completion

### Cost Observations
- Model mix: opus orchestrator, sonnet executor + verifier
- Notable: entire milestone (plan + execute + verify + complete) in a single session

---

## Cross-Milestone Trends

### Process Evolution

| Milestone | Commits | Phases | Key Change |
|-----------|---------|--------|------------|
| v1.0 | 125 | 6 | Initial process -- established GSD workflow patterns |
| v1.1 | ~10 | 2 | Small milestone -- fast execution, same tooling gaps surfaced |

### Cumulative Quality

| Milestone | Go LOC | Test Packages | Human Items Deferred |
|-----------|--------|---------------|---------------------|
| v1.0 | ~29,400 | 11 | 16 |
| v1.1 | ~29,600 | 11 | 3 (visual checks) |

### Top Lessons (Verified Across Milestones)

1. Minimal dependencies + stdlib-first approach keeps the codebase simple and the binary small (9.9MB)
2. Render-on-write (not render-on-read) eliminates redundant markdown processing on every page view
3. ROADMAP progress table and REQUIREMENTS traceability drift between phases -- needs atomic tooling updates (v1.0 + v1.1)
