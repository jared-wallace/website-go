# Phase 9: CSS Foundation - Research

**Researched:** 2026-03-28
**Domain:** Pure CSS — grain texture, box-shadow layering, semi-transparent fills, border-radius, comment rebrand
**Confidence:** HIGH

## Summary

Phase 9 is a pure CSS change to `web/static/main.css`. Every implementation decision has been locked in CONTEXT.md — grain technique is the only open question, and the three candidate approaches (SVG filter, CSS `background-image` noise, base64 PNG) have well-understood tradeoffs. No HTML templates are touched. No new files are created (STATE.md constraint).

The work is five surgical edits: replace the file header comment, add resting `box-shadow` to `.post-card`, add `background` fill to `.tag-pill`, change `.reaction-btn` `border-radius`, and append a grain pseudo-element block. The largest risk is grain mobile performance — the `position: fixed` pseudo-element technique is already mandated in STATE.md and is the correct approach; the remaining question is which noise generation method to use.

**Primary recommendation:** Use an inline SVG `url()` data URI for the grain — no external file, no new CSS file, pure CSS, cross-browser, zero runtime cost. Append all new rules to `main.css` in a clearly labeled block.

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

- **D-01:** Grain intensity ~3-5% opacity (barely-there whisper)
- **D-02:** Grain tone is warm sand — brownish-tan speckle light mode, muted amber dark mode
- **D-03:** Implementation uses `position: fixed` pseudo-element (never `background-attachment: fixed`)
- **D-04:** Post card resting shadow: tight layer `1px/3px blur 0.08 opacity` + ambient layer `4px/12px blur 0.05 opacity`
- **D-05:** Light mode shadow color `rgba(44, 36, 24, ...)` warm brown; dark mode `rgba(0, 0, 0, ...)`
- **D-06:** Hover shadow values (`0 8px 24px`) remain unchanged
- **D-07:** Tag pill fill `rgba(44, 95, 122, 0.08)` light mode / `rgba(201, 168, 76, 0.12)` dark mode
- **D-08:** Existing `1px solid var(--color-divider)` border on tag pill stays
- **D-09:** `.reaction-btn` changes from `border-radius: 2rem` to `border-radius: 4px`
- **D-10:** CSS file header comment "The Log" → "The Wild Meridian"

### Claude's Discretion

- Exact CSS technique for grain generation (SVG filter noise vs CSS gradient noise vs base64 PNG)
- Precise opacity values within decided ranges (3% vs 5% — tune by eye)
- Shadow transition timing if adding `box-shadow` to existing `.post-card` transition property
- Whether tag pill needs padding adjustment after adding fill background

### Deferred Ideas (OUT OF SCOPE)

None — discussion stayed within phase scope.
</user_constraints>

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| ATMO-01 | Site displays a subtle grain/noise texture overlay on the page background in both light and dark modes | SVG `feTurbulence` filter via `body::before` fixed pseudo-element; two `[data-theme="dark"]` overrides |
| ATMO-02 | Post cards display a warm two-layer resting shadow that provides depth context for the existing hover lift | CSS `box-shadow` with two comma-separated layers on `.post-card`; dark mode override via `[data-theme="dark"] .post-card` |
| TYPO-01 | Tag pills display with a filled semi-transparent background visible in both light and dark modes | Add `background` property to `.tag-pill`; dark mode override via `[data-theme="dark"] .tag-pill` |
| TYPO-02 | Reaction button uses `border-radius: 4px` matching the site's design system | Single property change on `.reaction-btn` at line 611 |
| HOUSE-01 | CSS file header comment reads "The Wild Meridian" instead of "The Log" | String replacement in lines 1-4 of main.css |
</phase_requirements>

## Standard Stack

### Core

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| CSS Custom Properties | (browser-native) | Color palette, theming | Already in use; `var(--color-accent)`, `var(--color-bg)` etc. drive all colors |
| `position: fixed` pseudo-element | (browser-native) | Grain overlay | Mandated by STATE.md; avoids `background-attachment: fixed` scroll jank on iOS |
| SVG `feTurbulence` filter | (browser-native) | Procedural noise generation | No external file, inlined as data URI, GPU-composited, zero repaint cost |

### Supporting

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| `mix-blend-mode: overlay` | (browser-native) | Blend grain into surface | Allows warm tonal grain without opaque overlay covering content |
| `pointer-events: none` | (browser-native) | Grain layer passthrough | Prevents pseudo-element intercepting clicks |
| `will-change: transform` | (browser-native) | Compositing hint | Optional; only if profiling shows paint cost on low-end mobile |

No npm packages. No build step. No new files.

**Grain technique decision (Claude's discretion):**

Three approaches considered:

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| SVG `feTurbulence` data URI | Zero files, procedural, GPU layer, ~400 bytes | Slightly verbose CSS | **USE THIS** |
| `background-image` CSS noise hack (repeating gradients) | Pure CSS, no data | Produces banding, not true noise | Avoid |
| Base64 PNG grain | Smallest CSS, true raster noise | External file OR 2–5KB base64 blob in CSS, not in STATE.md's "no new files" spirit | Fallback only |

SVG filter is the correct choice: procedural, contained in CSS, zero-dependency, and consistent across browsers that matter.

## Architecture Patterns

### File Layout

All changes land in one file: `web/static/main.css`. No new files per STATE.md constraint.

```
web/static/main.css
  Lines 1–4     : HOUSE-01 — header comment rebrand
  Lines 183–206 : ATMO-02 — .post-card resting shadow (edit existing rules)
  Lines 217–226 : TYPO-01 — .tag-pill fill (edit existing rule)
  Lines 608–621 : TYPO-02 — .reaction-btn border-radius (edit one line)
  End of file   : ATMO-01 — grain pseudo-element block (append)
```

### Pattern 1: Grain Pseudo-Element (ATMO-01)

**What:** `body::before` fixed full-viewport layer with SVG turbulence filter, `pointer-events: none`, low opacity, warm blend mode.

**When to use:** Any CSS-only ambient texture overlay on a themed surface.

```css
/* Source: MDN — SVG filter primitives / CSS background with data URI */

/* --- Background Grain Texture (ATMO-01) --- */

body::before {
  content: '';
  position: fixed;
  inset: 0;
  z-index: 9999;
  pointer-events: none;
  opacity: 0.04; /* tune within 0.03–0.05 range */
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='200' height='200'%3E%3Cfilter id='grain'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.65' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23grain)' opacity='1'/%3E%3C/svg%3E");
  background-repeat: repeat;
  background-size: 200px 200px;
}
```

Dark mode — the default brownish noise reads warmly on light backgrounds; on dark backgrounds it needs a tint adjustment. Since the SVG renders grayscale noise, the opacity and blend mode handle the visual difference. A `[data-theme="dark"] body::before` rule can reduce opacity slightly (the dark bg makes noise more visible):

```css
/* Source: Established pattern — grain visibility increases on dark backgrounds */

[data-theme="dark"] body::before {
  opacity: 0.03; /* dark surfaces amplify noise perception; dial back */
}
```

**Z-index note:** `z-index: 9999` places grain above content but `pointer-events: none` makes it invisible to interaction. Verify against `.site-nav` which has `z-index: 100` — grain sits above nav visually but does not block it. If nav interaction fails, lower grain z-index to `1000` (still above normal content stack).

### Pattern 2: Two-Layer Resting Shadow (ATMO-02)

**What:** Add `box-shadow` property to existing `.post-card` rule. Comma-separated layers: tight shadow first, ambient second.

```css
/* Source: CONTEXT.md D-04, D-05 — locked values */

/* Edit existing .post-card rule (lines 183–192) — add box-shadow: */
.post-card {
  /* ... existing properties ... */
  box-shadow:
    0 1px 3px rgba(44, 36, 24, 0.08),
    0 4px 12px rgba(44, 36, 24, 0.05);
}

/* Dark mode override — warm brown disappears on dark bg */
[data-theme="dark"] .post-card {
  box-shadow:
    0 1px 3px rgba(0, 0, 0, 0.2),
    0 4px 12px rgba(0, 0, 0, 0.15);
}
```

**Transition:** `.post-card` already has `transition: transform 200ms ease, box-shadow 200ms ease`. The resting shadow will animate to hover shadow automatically — no transition change needed.

### Pattern 3: Tag Pill Fill (TYPO-01)

**What:** Change `background: transparent` to `background: rgba(44, 95, 122, 0.08)` on `.tag-pill`. Add dark mode override.

```css
/* Source: CONTEXT.md D-07, D-08 */

/* Edit existing .tag-pill rule (lines 217–226) */
.tag-pill {
  /* ... existing properties ... */
  background: rgba(44, 95, 122, 0.08); /* faint ocean blue */
}

[data-theme="dark"] .tag-pill {
  background: rgba(201, 168, 76, 0.12); /* faint warm gold */
}
```

**Padding note (Claude's discretion):** Current padding is `2px 4px`. After adding the fill, this may feel tight. Consider `2px 6px` — adds 2px horizontal breathing room without disrupting grid layout. Not mandatory; check by eye.

### Pattern 4: Reaction Button Radius (TYPO-02)

**What:** Single-line change in `.reaction-btn` (line 611).

```css
/* Before */
border-radius: 2rem;

/* After */
border-radius: 4px;
```

### Pattern 5: Header Comment Rebrand (HOUSE-01)

**What:** Lines 1–4 of `main.css`. String replacement.

```css
/* Before */
/* ===================================================================
   Design System — The Log (weathered beach bar nautical theme)

/* After */
/* ===================================================================
   Design System — The Wild Meridian (weathered beach bar nautical theme)
```

### Anti-Patterns to Avoid

- **`background-attachment: fixed` for grain:** Triggers repaint on every scroll frame on iOS/Safari. STATE.md mandates `position: fixed` pseudo-element instead.
- **`filter: url(#grain)` referencing an inline `<svg>` in HTML:** Requires template changes. This phase is CSS-only.
- **New CSS file for grain:** STATE.md: "No new CSS files for v1.2 — all changes appended to `main.css` with labeled block comments."
- **Grain on `html::before` instead of `body::before`:** `html` pseudo-elements sometimes clip under navigation. `body::before` is more reliable with the existing sticky nav (`z-index: 100`).
- **High `baseFrequency` values (> 0.9):** Produces overly fine noise that appears as visual static rather than warm grain. `0.65` with `fractalNoise` gives the warmest, most organic texture.

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Procedural noise generation | A JS canvas grain renderer or raster image | SVG `feTurbulence` data URI in CSS | GPU composited, zero JS, zero files, zero latency |
| Dark mode color logic | JavaScript that reads computed styles | CSS `[data-theme="dark"]` overrides | Already the established pattern in this codebase |
| Shadow "warmth" calculation | Custom color math | rgba with existing `(44, 36, 24)` palette values from CONTEXT.md | Values already derived and locked |

## Runtime State Inventory

Step 2.5: SKIPPED — this is a pure CSS edit phase, not a rename/refactor/migration. No runtime state is affected by modifying CSS properties and comments.

## Common Pitfalls

### Pitfall 1: Grain Z-Index Blocking Interaction

**What goes wrong:** Grain pseudo-element with `z-index: 9999` sits above clickable elements. Clicks pass through the `pointer-events: none` layer but focus rings or browser accessibility tools may behave oddly.

**Why it happens:** `pointer-events: none` disables mouse events but some browser accessibility overlays layer above everything.

**How to avoid:** Set `pointer-events: none` explicitly. Test keyboard tab navigation through nav, cards, and reaction button after implementation.

**Warning signs:** Click events not firing on `.post-card` links; reaction button not responding.

### Pitfall 2: Grain Visible as Hard Tile Edge

**What goes wrong:** The 200×200 SVG tile produces a visible seam.

**Why it happens:** `feTurbulence` with `stitchTiles='stitch'` is designed to prevent seams but some browsers have rendering quirks at integer tile boundaries.

**How to avoid:** Use `stitchTiles='stitch'` (already in the template above). If seams still appear, increase tile size to 300×300.

**Warning signs:** Regular grid pattern visible at 200px intervals on a zoomed-out view.

### Pitfall 3: Dark Mode Grain Too Harsh

**What goes wrong:** 4% opacity on `#1A1F2E` dark background appears much more intense than on `#F5F0E8` light background — the low-luminance surface amplifies the noise contrast.

**Why it happens:** Perceived noise intensity scales with background luminance. Same opacity reads differently across light/dark.

**How to avoid:** Use separate opacity values: `0.04` light mode, `0.03` dark mode (already in the code examples above).

**Warning signs:** On dark mode the background looks "dirty" rather than "warm."

### Pitfall 4: Tag Pill Overflow After Fill

**What goes wrong:** Adding `background` to `.tag-pill` with existing `2px 4px` padding makes pills feel cramped — text touches the colored edge.

**Why it happens:** Transparent background doesn't emphasize padding; filled background makes tight padding obvious.

**How to avoid:** Increase horizontal padding to `2px 6px` or `2px 8px` at implementation time. This is in Claude's discretion scope.

**Warning signs:** Tag text sits flush against pill edge in the rendered page.

### Pitfall 5: Resting Shadow Conflicts with `.post-card:hover` Transition

**What goes wrong:** If `box-shadow` is NOT in the existing transition list, the rest-to-hover shadow change will be instantaneous instead of animated.

**Why it happens:** `transition: transform 200ms ease, box-shadow 200ms ease` is already present (line 191) — but only if `box-shadow` is in the property list.

**How to avoid:** Verify the transition value includes `box-shadow` before or when adding the resting shadow. It does (confirmed from `main.css` line 191). No change needed.

**Warning signs:** Shadow jumps on hover rather than animating smoothly.

## Code Examples

Verified patterns referencing actual file structure:

### HOUSE-01 — Header Comment (line 1-4)

```css
/* ===================================================================
   Design System — The Wild Meridian (weathered beach bar nautical theme)
   Light and dark mode via CSS custom properties.
   =================================================================== */
```

### ATMO-02 — Post Card Resting Shadow (edit lines 183-192, add dark override)

```css
.post-card {
  background-color: var(--color-surface);
  border: 1px solid var(--color-divider);
  border-radius: 4px;
  padding: 24px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  transition: transform 200ms ease, box-shadow 200ms ease;
  box-shadow:
    0 1px 3px rgba(44, 36, 24, 0.08),
    0 4px 12px rgba(44, 36, 24, 0.05);
}

[data-theme="dark"] .post-card {
  box-shadow:
    0 1px 3px rgba(0, 0, 0, 0.2),
    0 4px 12px rgba(0, 0, 0, 0.15);
}
```

### TYPO-01 — Tag Pill Fill (edit lines 217-226, add dark override)

```css
.tag-pill {
  display: inline-block;
  padding: 2px 6px; /* increased from 4px horizontal for legibility with fill */
  border: 1px solid var(--color-divider);
  border-radius: 2px;
  font-size: 14px;
  color: var(--color-text-muted);
  background: rgba(44, 95, 122, 0.08);
  line-height: 1.4;
}

[data-theme="dark"] .tag-pill {
  background: rgba(201, 168, 76, 0.12);
}
```

### TYPO-02 — Reaction Button Radius (edit line 611)

```css
/* Change only this property in .reaction-btn */
border-radius: 4px;
```

### ATMO-01 — Grain Block (append to end of file)

```css
/* --- Background Grain Texture (ATMO-01) --- */

body::before {
  content: '';
  position: fixed;
  inset: 0;
  z-index: 9999;
  pointer-events: none;
  opacity: 0.04;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='200' height='200'%3E%3Cfilter id='grain'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.65' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23grain)' opacity='1'/%3E%3C/svg%3E");
  background-repeat: repeat;
  background-size: 200px 200px;
}

[data-theme="dark"] body::before {
  opacity: 0.03;
}
```

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Raster PNG grain (external file) | SVG `feTurbulence` data URI | ~2020 (browser support matured) | Zero HTTP requests, procedural, CSS-contained |
| `background-attachment: fixed` | `position: fixed` pseudo-element | iOS Safari 13+ (2019) | Fixes scroll jank on mobile; Safari never properly supported `background-attachment: fixed` |
| CSS hack noise (repeating gradients) | SVG filter noise | Always suboptimal | Gradients produce banding; feTurbulence produces true stochastic noise |

**Deprecated/outdated:**

- `background-attachment: fixed` for parallax/fixed textures on mobile: broken on iOS Safari, causes scroll performance regression on Android. Confirmed by STATE.md and MDN.

## Open Questions

1. **Grain opacity fine-tuning**
   - What we know: 3–5% range is locked; `0.04` is the midpoint starting value
   - What's unclear: Whether the SVG noise at `baseFrequency=0.65` reads as "warm sand" or "cool neutral" at render time
   - Recommendation: Implement at `0.04`, view in browser at 100% zoom on both light and dark modes before committing. Tune within the locked range.

2. **Z-index ceiling for grain layer**
   - What we know: Nav is `z-index: 100`; grain at `9999` sits above everything
   - What's unclear: Whether any future Phase 10/11 elements (modals, overlays) will need z-index above 9999
   - Recommendation: Use `z-index: 9999` for now; document in Phase 10/11 research that grain occupies this slot.

## Environment Availability

Step 2.6: SKIPPED — this phase is pure CSS edits to `main.css`. No external tools, CLIs, databases, or runtimes beyond the existing Go dev server are required. Visual verification requires a browser (implicitly available).

## Validation Architecture

### Test Framework

| Property | Value |
|----------|-------|
| Framework | Go `testing` stdlib |
| Config file | none (stdlib, driven by `go test`) |
| Quick run command | `go test ./... -count=1` |
| Full suite command | `make test` (`go test ./... -v -race`) |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| ATMO-01 | Grain pseudo-element renders in both themes | manual | — | N/A |
| ATMO-02 | Post card resting shadow visible before hover | manual | — | N/A |
| TYPO-01 | Tag pill fill visible in light and dark modes | manual | — | N/A |
| TYPO-02 | Reaction button `border-radius: 4px` | manual + grep | `grep -n 'border-radius' web/static/main.css` | ✅ |
| HOUSE-01 | Header comment contains "The Wild Meridian", no "The Log" | manual + grep | `grep -c 'The Log' web/static/main.css` (must return 0) | ✅ |

**Note on automated testing for CSS:** CSS visual regressions require a browser; Go unit tests operate at the HTTP/handler layer and do not render CSS. The appropriate validation for this phase is manual browser inspection in both light and dark modes. HOUSE-01 and TYPO-02 have trivial grep-verifiable outcomes that can be scripted.

### Sampling Rate

- **Per task commit:** `grep -c 'The Log' web/static/main.css` returns 0 (HOUSE-01); `grep 'border-radius: 4px' web/static/main.css` matches `.reaction-btn` block (TYPO-02)
- **Per wave merge:** Visual browser check in light + dark mode; `make test` green (Go tests unaffected by CSS)
- **Phase gate:** All 5 success criteria confirmed in browser before `/gsd:verify-work`

### Wave 0 Gaps

None — no test infrastructure gaps. Go test suite exists and passes. CSS changes do not require new test files. Visual verification is manual by nature.

## Project Constraints (from CLAUDE.md)

Directives the planner must verify compliance with:

| Directive | Impact on Phase 9 |
|-----------|-------------------|
| Go with minimal dependencies — avoid large frameworks | Not applicable (pure CSS phase) |
| Must run as Docker container on port 8080 behind Nginx/ALB | Not applicable |
| All persistent data on EBS at /var/www/html | Not applicable |
| Leverage `frontend-design` skill for template/UI work | CSS changes are UI work — follow skill aesthetics guidance: texture, depth, grain overlays explicitly listed as valid techniques |
| No new CSS files for v1.2 (STATE.md) | All changes append to `main.css` |
| `position: fixed` pseudo-element for grain (STATE.md) | Grain technique locked |
| Code changes through GSD workflow | Enforced at orchestration level |
| Run code through `/simplify` skill before presenting | Applies to any code shown to user |

## Sources

### Primary (HIGH confidence)

- `web/static/main.css` (project file, read directly) — exact line numbers for all edit targets verified
- `.planning/phases/09-css-foundation/09-CONTEXT.md` — all locked decisions, color values, opacity ranges
- `.planning/STATE.md` — `position: fixed` mandate, no new CSS files constraint
- MDN Web Docs — `feTurbulence`, `position: fixed`, `box-shadow`, CSS custom properties (well-established browser primitives)

### Secondary (MEDIUM confidence)

- iOS Safari `background-attachment: fixed` scroll jank: widely reported, confirmed in STATE.md decision rationale
- SVG `feTurbulence` `stitchTiles='stitch'` for seamless tiling: documented SVG filter primitive behavior

### Tertiary (LOW confidence)

- Specific opacity perception difference between `0.03` and `0.04` on dark vs light backgrounds — tunable by eye; the range is correct, the exact value needs visual confirmation

## Metadata

**Confidence breakdown:**

- Standard stack: HIGH — pure browser-native CSS, no third-party libraries, all techniques verified
- Architecture: HIGH — exact line numbers confirmed from reading `main.css`, all edit targets identified
- Pitfalls: HIGH — grain z-index and mobile performance issues are well-documented browser behaviors
- Grain technique choice: HIGH — SVG feTurbulence is the established approach for CSS-only procedural noise

**Research date:** 2026-03-28
**Valid until:** 2026-09-28 (stable browser primitives; CSS Custom Properties and SVG filters are not changing)
