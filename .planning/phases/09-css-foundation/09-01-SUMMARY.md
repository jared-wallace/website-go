---
phase: 09-css-foundation
plan: "01"
subsystem: frontend
tags: [css, atmosphere, grain-texture, shadows, tag-pills, rebrand]
dependency_graph:
  requires: []
  provides: [grain-texture, post-card-shadows, tag-pill-fills, reaction-radius-fix, header-rebrand]
  affects: [web/static/main.css]
tech_stack:
  added: [SVG feTurbulence inline data URI for grain texture]
  patterns: [CSS pseudo-element overlay, position fixed with pointer-events none, dark mode overrides]
key_files:
  created: []
  modified:
    - web/static/main.css
decisions:
  - "Background grain uses position: fixed pseudo-element (never background-attachment: fixed) to avoid iOS scroll jank"
  - "Grain z-index: 9999 with pointer-events: none — above all content but click-through"
  - "Dark mode grain opacity reduced to 0.03 (vs 0.04 light) because dark surfaces amplify noise perception"
  - "Post card resting shadow uses warm brown rgba(44, 36, 24) in light mode, neutral black in dark mode"
  - "Tag pill fills use ocean blue rgba(44, 95, 122, 0.08) light / gold rgba(201, 168, 76, 0.12) dark"
  - "Reaction button changed from 2rem pill radius to 4px squared corners per design system"
  - "No new CSS files created — all changes appended to or edited within existing main.css"
metrics:
  duration: "~10min"
  completed: "2026-03-28"
  tasks_completed: 3
  files_modified: 1
---

# Phase 09 Plan 01: CSS Atmosphere & Component Fixes Summary

**One-liner:** Five surgical CSS edits to main.css — grain texture overlay, post card resting shadows, tag pill tinted fills, reaction button radius fix, and header comment rebrand from "The Log" to "The Wild Meridian."

## Tasks Completed

| Task | Description | Commit | Files |
|------|-------------|--------|-------|
| 1 (auto) | Rebrand header comment, add post card shadows, tag pill fills, reaction button radius | 772a873 | web/static/main.css |
| 2 (auto) | Append SVG feTurbulence grain texture overlay | 1a27523 | web/static/main.css |
| 3 (checkpoint) | Visual verification in browser — approved | — | — |

## What Was Built

Five CSS atmosphere improvements to `web/static/main.css`, all pure CSS with no template modifications:

- **HOUSE-01** — Header comment rebranded from "The Log" to "The Wild Meridian"
- **ATMO-01** — Background grain texture via `body::before` pseudo-element with inline SVG `feTurbulence` data URI. Uses `position: fixed` for mobile scroll performance, `pointer-events: none` for click-through, and separate light (0.04) / dark (0.03) opacities
- **ATMO-02** — Post card two-layer resting shadow (`0 1px 3px` tight + `0 4px 12px` ambient) with warm brown in light mode and neutral black in dark mode. Existing hover shadow unchanged per D-06
- **TYPO-01** — Tag pill filled backgrounds: ocean blue `rgba(44, 95, 122, 0.08)` light, gold `rgba(201, 168, 76, 0.12)` dark. Padding corrected to 4px 8px (on 4-point grid). Existing border preserved
- **TYPO-02** — Reaction button `border-radius` changed from `2rem` (pill) to `4px` (squared corners matching design system)

## Verification

All automated checks pass:

1. `grep -c "The Log" web/static/main.css` — returns 0
2. `grep -c "The Wild Meridian" web/static/main.css` — returns 1
3. `grep "feTurbulence" web/static/main.css` — matches
4. `grep "pointer-events: none" web/static/main.css` — matches
5. `grep "0 1px 3px rgba(44, 36, 24, 0.08)" web/static/main.css` — matches
6. `grep "rgba(44, 95, 122, 0.08)" web/static/main.css` — matches
7. `grep -c "background-attachment" web/static/main.css` — returns 0
8. Visual checkpoint approved by user (both light and dark modes)

## Decisions Made

1. **`position: fixed` for grain overlay** — `background-attachment: fixed` causes scroll jank on iOS. The pseudo-element approach with `pointer-events: none` is the safe pattern.

2. **z-index: 9999 for grain** — Higher than all content including `.site-nav` (z-index: 100), but entirely non-interactive via `pointer-events: none`. No click-through issues observed.

3. **Separate dark mode shadow palette** — Warm brown `rgba(44, 36, 24)` disappears against dark backgrounds, so dark mode post card shadows use neutral `rgba(0, 0, 0)` instead.

## Deviations from Plan

None — plan executed exactly as written across all three tasks.

## Known Stubs

None. All CSS is production-ready.

## Self-Check: PASSED

- `web/static/main.css` — MODIFIED (only file touched)
- Commit 772a873 — FOUND (Task 1: rebrand + component fixes)
- Commit 1a27523 — FOUND (Task 2: grain texture)
- Visual checkpoint — APPROVED (Task 3)
