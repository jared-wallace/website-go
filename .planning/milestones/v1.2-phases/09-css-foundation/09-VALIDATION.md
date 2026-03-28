---
phase: 9
slug: css-foundation
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-28
---

# Phase 9 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test + browser visual inspection |
| **Config file** | none — CSS-only phase, no new test infra needed |
| **Quick run command** | `grep -c "The Wild Meridian" web/static/main.css` |
| **Full suite command** | `go test ./... && grep -c "The Log" web/static/main.css \| grep -q "^0$"` |
| **Estimated runtime** | ~5 seconds |

---

## Sampling Rate

- **After every task commit:** Run `grep -c "The Wild Meridian" web/static/main.css`
- **After every plan wave:** Run `go test ./... && grep -c "The Log" web/static/main.css`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 5 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 09-01-01 | 01 | 1 | HOUSE-01 | grep | `grep "The Wild Meridian" web/static/main.css` | ✅ | ⬜ pending |
| 09-01-02 | 01 | 1 | HOUSE-01 | grep | `grep -c "The Log" web/static/main.css` (expect 0) | ✅ | ⬜ pending |
| 09-01-03 | 01 | 1 | ATMO-01 | grep | `grep "feTurbulence" web/static/main.css` | ✅ | ⬜ pending |
| 09-01-04 | 01 | 1 | ATMO-02 | grep | `grep "box-shadow" web/static/main.css` | ✅ | ⬜ pending |
| 09-01-05 | 01 | 1 | TYPO-01 | grep | `grep "tag" web/static/main.css` (background rules) | ✅ | ⬜ pending |
| 09-01-06 | 01 | 1 | TYPO-02 | grep | `grep "border-radius.*4px" web/static/main.css` | ✅ | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

*Existing infrastructure covers all phase requirements. CSS-only changes verified via grep and visual inspection.*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Grain texture visible on background | ATMO-01 | Visual appearance requires browser rendering | Load site in browser, confirm subtle noise texture on background in both light and dark modes |
| Grain does not cause scroll lag on mobile | ATMO-01 | Performance requires real device testing | Open on mobile device/emulator, scroll rapidly, confirm no jank |
| Post card shadow gives visual weight | ATMO-02 | Visual depth perception is subjective | View post listing, confirm cards have warm two-layer shadow at rest |
| Tag pills legible in both modes | TYPO-01 | Color contrast is visual | Toggle light/dark mode, confirm tag text readable against filled background |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 5s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
