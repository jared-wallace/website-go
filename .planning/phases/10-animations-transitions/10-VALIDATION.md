---
phase: 10
slug: animations-transitions
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-28
---

# Phase 10 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test + browser manual verification |
| **Config file** | none — CSS-only phase, no Go test infrastructure changes needed |
| **Quick run command** | `grep -c '@keyframes\|animation:\|transition:' static/css/main.css` |
| **Full suite command** | `make test` |
| **Estimated runtime** | ~5 seconds |

---

## Sampling Rate

- **After every task commit:** Run `grep -c '@keyframes\|animation:\|transition:' static/css/main.css`
- **After every plan wave:** Run `make test`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 5 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 10-01-01 | 01 | 1 | ANIM-05 | grep | `grep -c 'prefers-reduced-motion' static/css/main.css` | ✅ | ⬜ pending |
| 10-01-02 | 01 | 1 | ANIM-01 | grep | `grep '@keyframes fadeIn' static/css/main.css` | ❌ W0 | ⬜ pending |
| 10-01-03 | 01 | 1 | ANIM-02 | grep | `grep 'animation-delay' static/css/main.css` | ❌ W0 | ⬜ pending |
| 10-01-04 | 01 | 1 | ANIM-03 | grep | `grep 'transition.*color\|transition.*background' static/css/main.css` | ❌ W0 | ⬜ pending |
| 10-01-05 | 01 | 1 | ANIM-04 | manual | Browser check: no FOUC on page load | N/A | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

*Existing infrastructure covers all phase requirements — this is a CSS-only phase with no new test framework needed.*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| No white flash on dark mode page load | ANIM-04 | Requires visual browser check with `prefers-color-scheme: dark` | 1. Set OS to dark mode 2. Hard-refresh page 3. Verify no white flash before dark theme renders |
| Stagger animation visible on post list | ANIM-02 | Timing/visual effect requires human eye | 1. Load blog index page 2. Observe post cards appearing sequentially |
| Dark mode toggle smooth transition | ANIM-03 | Visual smoothness is subjective | 1. Toggle dark/light mode 2. Verify color change is gradual, not instant |
| Reduced motion disables all animation | ANIM-05 | Requires OS accessibility setting | 1. Enable `prefers-reduced-motion: reduce` in OS 2. Verify zero animation including reaction bounce |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 5s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
