---
phase: 7
slug: rebrand-navigation
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-28
---

# Phase 7 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | Go stdlib `testing` package |
| **Config file** | none — `go test` discovers tests by convention |
| **Quick run command** | `go test ./internal/handler/blog/...` |
| **Full suite command** | `go test ./...` |
| **Estimated runtime** | ~5 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./internal/handler/blog/...`
- **After every plan wave:** Run `go test ./...`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 5 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 07-01-01 | 01 | 1 | BRAND-01 | integration | `go test ./internal/handler/blog/... -run TestListOGMeta` | ✅ | ⬜ pending |
| 07-01-02 | 01 | 1 | BRAND-02 | unit | `go test ./internal/handler/blog/... -run TestServeRSS` | ✅ | ⬜ pending |
| 07-01-03 | 01 | 1 | BRAND-03 | integration | `go test ./internal/handler/blog/... -run TestListOGMeta` | ✅ | ⬜ pending |
| 07-01-04 | 01 | 1 | BRAND-04 | manual smoke | `curl http://localhost:8080/ \| grep "Jared Wallace"` | n/a | ⬜ pending |
| 07-01-05 | 01 | 1 | NAV-01 | manual smoke | `curl http://localhost:8080/ \| grep rss-link` | n/a | ⬜ pending |
| 07-01-06 | 01 | 1 | NAV-02 | integration | `go test ./internal/handler/blog/... -run TestListOGMeta` | ✅ | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

Existing infrastructure covers all phase requirements. No new test files or fixtures needed.

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Copyright footer renders "Jared Wallace" | BRAND-04 | No existing test covers footer text rendering | Build and curl localhost:8080, grep for "Jared Wallace" |
| RSS icon visible in footer | NAV-01 | Visual element — no test for SVG presence | Build and curl localhost:8080, grep for "rss-link" class |

---

## Validation Sign-Off

- [x] All tasks have `<automated>` verify or Wave 0 dependencies
- [x] Sampling continuity: no 3 consecutive tasks without automated verify
- [x] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [x] Feedback latency < 5s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
