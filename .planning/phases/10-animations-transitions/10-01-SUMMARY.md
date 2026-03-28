---
phase: 10-animations-transitions
plan: 01
subsystem: ui
tags: [css, animations, dark-mode, reduced-motion, transitions]

# Dependency graph
requires:
  - phase: 09-css-foundation
    provides: main.css with base styles, dark mode variables, and .post-card transition rules

provides:
  - "@keyframes fade-in with .container page load animation (ANIM-02)"
  - "Card stagger on .card-grid .post-card with 6-step 75ms interval delays (ANIM-03)"
  - ".theme-ready gate via requestAnimationFrame preventing flash-of-transitions on page load (ANIM-05)"
  - "Dark mode color blend transitions on all themed surfaces (ANIM-04)"
  - "prefers-reduced-motion guard disabling all animations including pre-existing reaction-bounce (ANIM-01)"

affects: [11-typography-polish, future-css-phases]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - ".theme-ready CSS gate pattern: JS adds class after first paint, enabling transitions only for user-initiated changes"
    - "prefers-reduced-motion guard as last block in animation section, covers both new and legacy animations"
    - "Card stagger via :nth-child selectors with explicit animation-delay at 75ms intervals, capped at 6"

key-files:
  created: []
  modified:
    - web/static/main.css
    - web/templates/base.html

key-decisions:
  - "Gated dark mode transitions behind .theme-ready to prevent flash-of-transition on page load — JS adds the class via requestAnimationFrame after first paint"
  - "250ms duration chosen as human-perceptible but sub-threshold for 'slow' — based on RESEARCH.md D-01/D-02 findings"
  - "Card stagger capped at 6 cards (375ms max total) to avoid long waits on short post lists"
  - ".theme-ready .post-card transition must include all 5 properties (transform, box-shadow, background-color, color, border-color) to avoid clobbering existing hover lift effect"
  - "Reduced-motion guard placed as final block and explicitly covers reaction-bounce (pre-existing unguarded animation)"

patterns-established:
  - "Theme-ready gate: requestAnimationFrame callback adds .theme-ready class to <html>, CSS transitions use .theme-ready selector prefix"
  - "Reduced-motion last: @media (prefers-reduced-motion: reduce) block always final in animation section, covers all animations including pre-existing ones"
  - "Card stagger: .card-grid .post-card:nth-child(N) selectors with explicit animation-delay values"

requirements-completed: [ANIM-01, ANIM-02, ANIM-03, ANIM-04, ANIM-05]

# Metrics
duration: 10min
completed: 2026-03-28
---

# Phase 10 Plan 01: Animations & Transitions Summary

**CSS motion system with page fade-in, 6-card stagger, .theme-ready-gated dark mode blends, flash prevention via requestAnimationFrame, and prefers-reduced-motion guard covering all animations including pre-existing reaction-bounce**

## Performance

- **Duration:** ~10 min
- **Completed:** 2026-03-28
- **Tasks:** 2 (1 auto + 1 human-verify checkpoint)
- **Files modified:** 2

## Accomplishments

- Added `@keyframes fade-in` with `.container` page-load animation at 250ms ease-out (ANIM-02)
- Added `.card-grid .post-card` stagger via :nth-child selectors with 75ms intervals for first 6 cards (ANIM-03)
- Implemented `.theme-ready` gate pattern: inline script extended with `requestAnimationFrame` to add the class after first paint, enabling dark mode transitions only for user-initiated toggles — preventing flash-of-transitions on load (ANIM-04, ANIM-05)
- Added `@media (prefers-reduced-motion: reduce)` block as final CSS block covering all new animations and the pre-existing unguarded `reaction-bounce` (ANIM-01)
- Preserved existing `.post-card` hover transitions (transform + box-shadow) by including all 5 properties in the `.theme-ready .post-card` rule

## Task Commits

1. **Task 1: Add animation CSS and extend theme-ready gate script** - `61175ff` (feat)
2. **Task 2: Visual verification of motion system** - user-approved checkpoint

## Files Modified

- `web/static/main.css` - Appended animations section with fade-in keyframe, .container animation, card stagger, .theme-ready dark mode transitions, and prefers-reduced-motion guard
- `web/templates/base.html` - Extended inline script with `requestAnimationFrame` callback adding `.theme-ready` to `<html>` after first paint

## Decisions Made

- **theme-ready gate pattern:** Dark mode transitions gated behind `.theme-ready` to prevent firing during initial render. JS adds class via `requestAnimationFrame` after first paint.
- **250ms duration:** Sub-perceptible threshold per RESEARCH.md findings.
- **Stagger cap at 6 cards:** 375ms total max prevents long waits on short post lists.
- **All-5-property .post-card transition:** Explicitly lists transform, box-shadow, background-color, color, border-color to avoid clobbering hover effect.
- **Reduced-motion covers reaction-bounce:** Pre-existing animation brought into WCAG 2.3.3 compliance.

## Deviations from Plan

None.

## Issues Encountered

None.

---
*Phase: 10-animations-transitions | Completed: 2026-03-28*
