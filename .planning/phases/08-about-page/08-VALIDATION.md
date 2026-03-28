---
phase: 8
slug: about-page
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-28
---

# Phase 8 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | Go stdlib `testing` package |
| **Config file** | none — `go test ./...` discovers tests by convention |
| **Quick run command** | `go test ./internal/handler/blog/... -run TestAbout` |
| **Full suite command** | `go test ./...` |
| **Estimated runtime** | ~2 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./internal/handler/blog/... -run TestAbout`
- **After every plan wave:** Run `go test ./...`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 5 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 08-01-01 | 01 | 1 | ABOUT-01 | unit | `go test ./internal/handler/blog/... -run TestAboutPageStatus` | ❌ W0 | ⬜ pending |
| 08-01-02 | 01 | 1 | ABOUT-01 | unit | `go test ./internal/handler/blog/... -run TestAboutNavLink` | ❌ W0 | ⬜ pending |
| 08-01-03 | 01 | 1 | ABOUT-02 | unit | `go test ./internal/handler/blog/... -run TestAboutPageContent` | ❌ W0 | ⬜ pending |
| 08-01-04 | 01 | 1 | ABOUT-03 | unit | `go test ./internal/handler/blog/... -run TestAboutPageChrome` | ❌ W0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `internal/handler/blog/about_test.go` — stubs for ABOUT-01, ABOUT-02, ABOUT-03

*Existing `testing` + `net/http/httptest` infrastructure covers all cases. No new framework install needed.*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Nav link visual placement (between tagline and dark toggle) | ABOUT-01 | CSS layout order is a visual concern | Inspect nav in browser; About link sits between tagline and moon/sun icon |
| Nautical design consistency | ABOUT-03 | Subjective visual match | Compare about page chrome to blog list page — header, footer, dark mode toggle present and styled identically |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 5s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
