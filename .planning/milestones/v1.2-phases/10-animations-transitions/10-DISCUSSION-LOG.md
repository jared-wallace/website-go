# Phase 10: Animations & Transitions - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-03-28
**Phase:** 10-animations-transitions
**Areas discussed:** Entry animation feel, Dark mode transition speed, Stagger cap & rhythm, Flash prevention timing

---

## Entry Animation Feel

| Option | Description | Selected |
|--------|-------------|----------|
| Whisper | 200-300ms opacity-only fade. Readers barely notice it consciously but the page feels "alive". No translate/slide. | :white_check_mark: |
| Gentle lift | 400-500ms fade + subtle 8-12px translate-up. Content rises into place like a slow tide. More noticeable, more personality. | |
| Dramatic entrance | 600ms+ with larger translate (20px+). Unmistakable motion — the page announces itself. Risks feeling slow on repeat visits. | |

**User's choice:** Whisper (Recommended)
**Notes:** Opacity-only fade at ~250ms ease-out. Consistent with the site's "barely there" design philosophy established in Phase 9 (grain texture at 3-5% opacity).

---

## Dark Mode Transition Speed

| Option | Description | Selected |
|--------|-------------|----------|
| Quick blend | 200-300ms. Feels responsive and intentional — colors shift smoothly but don't linger. Matches existing 200ms transitions. | :white_check_mark: |
| Languid morph | 500-600ms. Dreamy crossfade — the whole page slowly breathes into the new palette. More theatrical, can feel sluggish. | |
| Instant with fade | Background swaps instantly, text/accents fade over 150ms. Snappy but less cohesive. | |

**User's choice:** Quick blend (Recommended)
**Notes:** 250ms ease across background-color, color, and border-color. Maintains consistency with the existing 200ms timing baseline.

---

## Stagger Cap & Rhythm

| Option | Description | Selected |
|--------|-------------|----------|
| Fast cascade, cap at 6 | 75ms between cards, max 6 staggered. Cards 7+ appear instantly. Quick domino effect. | :white_check_mark: |
| Lazy wave, cap at 4 | 150ms between cards, max 4 staggered. Slower, more deliberate reveal. | |
| You decide | Claude picks timing and cap based on post count and fade-in feel. | |

**User's choice:** Fast cascade, cap at 6
**Notes:** Total cascade ~375ms. Pairs well with the whisper fade-in — quick enough to not feel like waiting.

---

## Flash Prevention Timing

| Option | Description | Selected |
|--------|-------------|----------|
| Inline script in `<head>` | Tiny inline script reads localStorage and sets data-theme BEFORE CSS. Adds .theme-ready after first paint via requestAnimationFrame. Zero flash. | :white_check_mark: |
| DOMContentLoaded | Add .theme-ready in main.js IIFE after DOM ready. Simpler but brief flash window. | |
| You decide | Claude picks the most robust approach for the template structure. | |

**User's choice:** Inline script in `<head>` (Recommended)
**Notes:** Belt-and-suspenders approach. Script before CSS link guarantees no flash. requestAnimationFrame ensures transitions don't fire until after first paint.

---

## Claude's Discretion

- Exact easing function choices within ease-out / ease family
- Whether card stagger reuses page fade-in keyframe or defines its own
- CSS @property vs standard transition for dark mode custom property animation
- Whether .theme-ready gate consolidates with existing inline theme scripts
- Precise selector list for dark mode transitions

## Deferred Ideas

None — discussion stayed within phase scope
