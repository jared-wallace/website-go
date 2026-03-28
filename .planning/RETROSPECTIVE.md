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

---

## Cross-Milestone Trends

### Process Evolution

| Milestone | Commits | Phases | Key Change |
|-----------|---------|--------|------------|
| v1.0 | 125 | 6 | Initial process -- established GSD workflow patterns |

### Cumulative Quality

| Milestone | Go LOC | Test Packages | Human Items Deferred |
|-----------|--------|---------------|---------------------|
| v1.0 | ~29,400 | 11 | 16 |

### Top Lessons (Verified Across Milestones)

1. Minimal dependencies + stdlib-first approach keeps the codebase simple and the binary small (9.9MB)
2. Render-on-write (not render-on-read) eliminates redundant markdown processing on every page view
