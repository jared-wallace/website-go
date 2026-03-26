---
phase: 2
slug: public-blog
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-26
---

# Phase 2 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test (stdlib) |
| **Config file** | none — existing `go test ./...` infrastructure from Phase 1 |
| **Quick run command** | `go test ./internal/handler/... ./internal/service/...` |
| **Full suite command** | `make test` |
| **Estimated runtime** | ~5 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./internal/handler/... ./internal/service/...`
- **After every plan wave:** Run `make test`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 5 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 02-01-01 | 01 | 1 | BLOG-01 | unit | `go test ./internal/service/...` | ❌ W0 | ⬜ pending |
| 02-01-02 | 01 | 1 | BLOG-02 | unit | `go test ./internal/handler/...` | ❌ W0 | ⬜ pending |
| 02-01-03 | 01 | 1 | BLOG-03 | unit | `go test ./internal/handler/...` | ❌ W0 | ⬜ pending |
| 02-02-01 | 02 | 1 | BLOG-04 | manual | browser check | n/a | ⬜ pending |
| 02-02-02 | 02 | 1 | BLOG-05 | unit | `go test ./internal/handler/...` | ❌ W0 | ⬜ pending |
| 02-03-01 | 03 | 2 | BLOG-08 | unit | `go test ./internal/service/...` | ❌ W0 | ⬜ pending |
| 02-03-02 | 03 | 2 | BLOG-11 | manual | browser check | n/a | ⬜ pending |
| 02-03-03 | 03 | 2 | BLOG-12 | unit | `go test ./internal/handler/...` | ❌ W0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `internal/service/post_service_test.go` — stubs for post listing, pagination, reading time, excerpt
- [ ] `internal/handler/post_handler_test.go` — stubs for route handling, slug lookup, 404, ToC
- [ ] Test fixtures: sample posts with markdown content for rendering verification

*Existing `go test` infrastructure from Phase 1 covers framework needs.*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Nautical design renders correctly on mobile/desktop | BLOG-04 | Visual layout verification requires browser | Open listing page in Chrome DevTools, toggle mobile/desktop viewports, verify card grid and typography |
| Dark mode toggle and persistence | BLOG-11 | localStorage + CSS interaction requires browser | Click toggle, verify palette swap, refresh page, verify persisted state |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 5s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
