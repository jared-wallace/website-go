# Phase 9: CSS Foundation - Context

**Gathered:** 2026-03-28
**Status:** Ready for planning

<domain>
## Phase Boundary

The site's visual atmosphere is elevated through pure CSS changes to `main.css` — background grain texture, post card resting shadows, tag pill legibility, reaction button radius consistency, and CSS comment rebrand. No HTML templates are modified in this phase.

</domain>

<decisions>
## Implementation Decisions

### Background Grain Texture (ATMO-01)
- **D-01:** Grain intensity is barely-there whisper (~3-5% opacity). Readers should feel warmth subconsciously without being able to pinpoint the texture while reading.
- **D-02:** Grain tone is warm sand — brownish-tan speckle in light mode, muted amber in dark mode. Reinforces the beach bar palette (not neutral monochrome).
- **D-03:** Implementation uses `position: fixed` pseudo-element (already decided in STATE.md — never `background-attachment: fixed` for mobile scroll performance).

### Post Card Resting Shadow (ATMO-02)
- **D-04:** Gentle float depth — two-layer shadow at rest. Tight shadow (~1px/3px blur, 0.08 opacity) plus ambient shadow (~4px/12px blur, 0.05 opacity). Cards sit slightly above the page like postcards on a table.
- **D-05:** Shadow color is warm brown tint in light mode using existing `rgba(44, 36, 24, ...)` palette. Dark mode uses neutral `rgba(0, 0, 0, ...)` since warm tones disappear on dark backgrounds.
- **D-06:** Hover shadow remains unchanged (existing `0 8px 24px` values). The rest-to-hover jump should feel noticeable but not dramatic.

### Tag Pill Fill (TYPO-01)
- **D-07:** Fill derives from the accent color at low opacity — `rgba(44, 95, 122, 0.08)` in light mode (faint ocean blue), `rgba(201, 168, 76, 0.12)` in dark mode (faint warm gold).
- **D-08:** Existing `1px solid var(--color-divider)` border stays. The fill adds legibility without replacing the border treatment.

### Reaction Button Radius (TYPO-02)
- **D-09:** Change `.reaction-btn` from `border-radius: 2rem` (pill shape) to `border-radius: 4px` to match the site's design system. One-line change.

### CSS Comment Rebrand (HOUSE-01)
- **D-10:** Replace CSS file header comment "The Log" with "The Wild Meridian". Mechanical find-and-replace.

### Claude's Discretion
- Exact CSS technique for grain generation (SVG filter noise vs CSS gradient noise vs base64 PNG — choose based on performance and browser support)
- Precise opacity values within the decided ranges (e.g., 3% vs 5% for grain — tune by eye)
- Shadow transition timing if adding `box-shadow` to the existing `.post-card` transition property
- Whether tag pill needs any padding adjustment after adding the fill background

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### CSS (primary target file)
- `web/static/main.css` — All changes land here. Lines 1-4: header comment (HOUSE-01). Lines 183-206: `.post-card` and hover shadows (ATMO-02). Lines 217-226: `.tag-pill` (TYPO-01). Lines 608-621: `.reaction-btn` with `border-radius: 2rem` (TYPO-02).

### Color Palette Reference
- `web/static/main.css` lines 8-27 — `:root` and `[data-theme="dark"]` custom properties define the full color palette

### Requirements
- `.planning/REQUIREMENTS.md` — ATMO-01, ATMO-02, TYPO-01, TYPO-02, HOUSE-01

### Prior Decisions
- `.planning/STATE.md` — "Accumulated Context > Decisions" section contains v1.2-specific constraints (no new CSS files, grain pseudo-element technique, reduced-motion guards)

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- CSS custom properties (`:root` / `[data-theme="dark"]`) provide the full color palette — grain and shadows should use these values or their rgba equivalents
- Existing hover shadow on `.post-card:hover` uses `rgba(44, 36, 24, 0.12)` — resting shadow uses same color family at lower opacity for consistency

### Established Patterns
- Dark mode overrides use `[data-theme="dark"] .selector` pattern
- `transition` properties already on `.post-card` (transform 200ms ease, box-shadow 200ms ease) — shadow addition will animate naturally
- `.tag-pill` currently has `background: transparent` — changing to semi-transparent fill is non-breaking

### Integration Points
- All changes are in `web/static/main.css` — no other files need modification
- Grain pseudo-element attaches to `body::before` or similar — must not interfere with existing layout or z-index stacking
- Verify grain renders correctly over both `--color-bg` and `--color-surface` areas

</code_context>

<specifics>
## Specific Ideas

No specific requirements beyond the decided values — open to standard approaches for implementation. The grain should be the kind of thing you'd only notice if someone pointed it out.

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope

</deferred>

---

*Phase: 09-css-foundation*
*Context gathered: 2026-03-28*
