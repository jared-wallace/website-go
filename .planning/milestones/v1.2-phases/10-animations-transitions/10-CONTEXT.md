# Phase 10: Animations & Transitions - Context

**Gathered:** 2026-03-28
**Status:** Ready for planning

<domain>
## Phase Boundary

The site's motion system is complete, safe, and capable. Entry animations delight users who prefer motion, dark mode transitions are smooth, and users who prefer reduced motion experience zero flashing or layout shift. Changes span `main.css` (animations, keyframes, reduced-motion guards) and `base.html` (inline theme script). The existing `main.js` dark mode toggle is updated to work with the `.theme-ready` gate.

</domain>

<decisions>
## Implementation Decisions

### Page Load Fade-In (ANIM-02)
- **D-01:** Whisper-level fade — opacity-only, ~250ms, `ease-out`. No translate or slide. Readers feel the page is "alive" without consciously noticing the animation.
- **D-02:** Applied to `.main-content` (or equivalent content wrapper). Navigation and chrome appear instantly.

### Post Card Stagger (ANIM-03)
- **D-03:** Fast cascade with 75ms intervals between cards, capped at 6 cards. Cards 7+ appear with no delay.
- **D-04:** Total cascade duration ~375ms. Uses the same whisper-level opacity fade as the page load animation for visual consistency.

### Dark Mode Transition Speed (ANIM-04)
- **D-05:** Quick blend at 250ms ease. Matches the existing 200ms timing on card and button transitions for consistency.
- **D-06:** Transitions applied to `background-color`, `color`, and `border-color` across all themed surfaces. Only active when `.theme-ready` class is present (prevents page-load flash).

### Flash Prevention / .theme-ready Gate (ANIM-05)
- **D-07:** Inline `<script>` in `<head>`, placed BEFORE the CSS `<link>`. Reads `localStorage('theme')` and sets `data-theme="dark"` before any CSS paints.
- **D-08:** `.theme-ready` class added to `<html>` via `requestAnimationFrame` callback — enables transitions only after first paint is complete.
- **D-09:** CSS transition rules are scoped behind `.theme-ready` selector so they are inert during initial page load.

### Reduced Motion Safety (ANIM-01)
- **D-10:** `prefers-reduced-motion: reduce` media query disables ALL animations: the new page fade-in, card stagger, dark mode transitions, AND the existing `reaction-bounce` keyframe.
- **D-11:** Reduced-motion users see instant state changes with zero animation — no degraded/shortened animation compromise.

### Claude's Discretion
- Exact easing function choices within the "ease-out" / "ease" family
- Whether card stagger reuses the page fade-in keyframe or defines its own
- CSS `@property` approach for custom property animation (ANIM-04) vs standard transition approach — choose based on browser support and complexity
- Whether `.theme-ready` gate script also handles the existing `data-theme` initialization (consolidating the two inline scripts if one exists in `base.html`)
- Precise selector list for dark mode transitions (which elements get explicit transition rules vs inheriting)

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### CSS (primary target file)
- `web/static/main.css` — All animation CSS lands here. Line 191: existing `.post-card` transitions. Line 633: `.reaction-btn` transitions. Lines 650-660: `reaction-bounce` keyframe (needs reduced-motion guard). Lines 8-27: `:root` and `[data-theme="dark"]` custom properties.

### JavaScript
- `web/static/main.js` — Dark mode toggle handler (lines 2-18). Must work with `.theme-ready` gate. Reaction bounce class toggle (line 51).

### Templates
- `web/templates/base.html` — `<head>` section where inline theme script goes. Must be placed BEFORE CSS `<link>`.

### Requirements
- `.planning/REQUIREMENTS.md` — ANIM-01, ANIM-02, ANIM-03, ANIM-04, ANIM-05

### Prior Decisions
- `.planning/STATE.md` — "Accumulated Context > Decisions" section: `.theme-ready` gate decided, `prefers-reduced-motion` guard scope decided, no new CSS files rule

### Prior Phase Context
- `.planning/phases/09-css-foundation/09-CONTEXT.md` — Phase 9 established shadow transitions on `.post-card`, tag pill fills, and grain overlay. Phase 10 animations must not conflict with these.

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- Existing `transition: transform 200ms ease, box-shadow 200ms ease` on `.post-card` — dark mode transitions append to this
- `reaction-bounce` keyframe already defined — just needs `prefers-reduced-motion` guard wrapping
- Dark mode toggle in `main.js` already reads/writes `localStorage('theme')` and toggles `data-theme`

### Established Patterns
- Dark mode uses `[data-theme="dark"]` attribute on `<html>` element
- Transitions use 200ms ease as the baseline timing
- `main.js` is a single IIFE — new behavior integrates into existing structure
- Inline scripts in `base.html` `<head>` set theme before paint (pattern to follow for `.theme-ready`)

### Integration Points
- `base.html` `<head>`: Add inline script before CSS link for flash prevention
- `main.css`: Add `@keyframes fade-in`, card stagger rules, dark mode transition rules, `prefers-reduced-motion` media query
- `main.js`: Dark mode toggle may need adjustment if `.theme-ready` affects transition behavior
- Existing `reaction-bounce` animation: Wrap in reduced-motion guard (currently unprotected)

</code_context>

<specifics>
## Specific Ideas

No specific requirements beyond the decided values. The overall motion philosophy is "whisper" — animations should make the site feel polished without drawing attention to themselves. A visitor should think "this site feels nice" not "look at those animations."

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope

</deferred>

---

*Phase: 10-animations-transitions*
*Context gathered: 2026-03-28*
