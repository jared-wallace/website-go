---
phase: 11
slug: template-changes
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-28
---

# Phase 11 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | Go `testing` package (stdlib) |
| **Config file** | none — `go test` discovers `*_test.go` files automatically |
| **Quick run command** | `go test ./internal/handler/blog/... -run TestList` |
| **Full suite command** | `go test ./...` |
| **Estimated runtime** | ~5 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./internal/handler/blog/... -run TestList`
- **After every plan wave:** Run `go test ./...`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 5 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 11-00-01 | 00 | 0 | NAV-01 | unit | `go test ./internal/handler/blog/... -run TestNavAboutLinkRemoved` | ❌ W0 | ⬜ pending |
| 11-00-02 | 00 | 0 | NAV-02 | unit | `go test ./internal/handler/blog/... -run TestFooterTwoSection` | ❌ W0 | ⬜ pending |
| 11-00-03 | 00 | 0 | NAV-03 | unit | `go test ./internal/handler/blog/... -run TestFooterPersonalityPhrase` | ❌ W0 | ⬜ pending |
| 11-00-04 | 00 | 0 | NAV-04 | unit | `go test ./internal/handler/blog/... -run TestNavAriaLabels` | ❌ W0 | ⬜ pending |
| 11-00-05 | 00 | 0 | ATMO-03 | unit | `go test ./internal/handler/blog/... -run TestRopeDividerSVG` | ❌ W0 | ⬜ pending |
| 11-00-06 | 00 | 0 | TYPO-03 | unit | `go test ./internal/handler/blog/... -run TestListHero` | ❌ W0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `internal/handler/blog/handler_test.go` — six new test functions following existing `strings.Contains` pattern
- No new test infrastructure needed — existing test setup covers all requirements

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| SVG rope renders visually as twisted two-strand pattern | ATMO-03 | Visual appearance cannot be verified by string matching | Open browser, inspect footer rope divider visually |
| Footer stacks vertically on mobile | NAV-02 | Responsive layout requires viewport testing | Resize browser below 767px, verify column stacking |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 5s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
