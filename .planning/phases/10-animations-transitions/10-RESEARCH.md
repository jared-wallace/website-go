# Phase 10: Animations & Transitions - Research

**Researched:** 2026-03-28
**Domain:** CSS animations, CSS transitions, flash-of-unstyled-content (FOUC) prevention, reduced-motion accessibility
**Confidence:** HIGH

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

**Page Load Fade-In (ANIM-02)**
- D-01: Whisper-level fade — opacity-only, ~250ms, `ease-out`. No translate or slide.
- D-02: Applied to `.main-content` (or equivalent content wrapper). Navigation and chrome appear instantly.

**Post Card Stagger (ANIM-03)**
- D-03: Fast cascade with 75ms intervals between cards, capped at 6 cards. Cards 7+ appear with no delay.
- D-04: Total cascade duration ~375ms. Uses the same whisper-level opacity fade as the page load animation for visual consistency.

**Dark Mode Transition Speed (ANIM-04)**
- D-05: Quick blend at 250ms ease. Matches the existing 200ms timing on card and button transitions for consistency.
- D-06: Transitions applied to `background-color`, `color`, and `border-color` across all themed surfaces. Only active when `.theme-ready` class is present.

**Flash Prevention / .theme-ready Gate (ANIM-05)**
- D-07: Inline `<script>` in `<head>`, placed BEFORE the CSS `<link>`. Reads `localStorage('theme')` and sets `data-theme="dark"` before any CSS paints.
- D-08: `.theme-ready` class added to `<html>` via `requestAnimationFrame` callback — enables transitions only after first paint is complete.
- D-09: CSS transition rules scoped behind `.theme-ready` selector so they are inert during initial page load.

**Reduced Motion Safety (ANIM-01)**
- D-10: `prefers-reduced-motion: reduce` media query disables ALL animations: new page fade-in, card stagger, dark mode transitions, AND existing `reaction-bounce` keyframe.
- D-11: Zero animation for reduced-motion users — no degraded/shortened compromise. Instant state changes only.

### Claude's Discretion
- Exact easing function choices within the "ease-out" / "ease" family
- Whether card stagger reuses the page fade-in keyframe or defines its own
- CSS `@property` approach for custom property animation (ANIM-04) vs standard transition approach
- Whether `.theme-ready` gate script consolidates with existing inline script in `base.html`
- Precise selector list for dark mode transitions

### Deferred Ideas (OUT OF SCOPE)
None — discussion stayed within phase scope
</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| ANIM-01 | All animations (existing `reaction-bounce` + new) wrapped in `prefers-reduced-motion` guards | Confirmed pattern: single `@media (prefers-reduced-motion: reduce)` block at end of main.css; existing `reaction-bounce` at line 653 needs guard added |
| ANIM-02 | Main content fades in on page load with subtle opacity transition | `@keyframes fade-in` + `animation` on `.main-content`; nav/chrome excluded; 250ms ease-out confirmed in UI-SPEC |
| ANIM-03 | Post cards on list page stagger entrance with CSS animation delays | `:nth-child(1-6)` selectors with `animation-delay` increments of 75ms; cards 7+ inherit the keyframe with no delay; uses same `fade-in` keyframe as ANIM-02 |
| ANIM-04 | Dark mode toggle produces smooth color transition via CSS `@property` or standard transition | Standard `transition` approach recommended over `@property` — see Architecture Patterns; `.theme-ready` gate prevents flash |
| ANIM-05 | Dark mode transitions gated behind `.theme-ready` class added by JS post-load | `requestAnimationFrame` callback in existing inline `<head>` script; CSS rules scoped as `.theme-ready selector { transition: ... }` |
</phase_requirements>

---

## Summary

Phase 10 is a pure CSS + minimal JS phase. No new libraries, no new files, no new routes. All changes land in three existing files: `web/static/main.css`, `web/templates/base.html`, and potentially `web/static/main.js` (though current analysis suggests main.js needs no changes).

The technical work divides into four independent blocks: (1) the `@keyframes fade-in` + `.main-content` fade, (2) the card stagger rules using `:nth-child`, (3) the `.theme-ready`-gated dark mode transition rules, and (4) the `prefers-reduced-motion` guard wrapping all of the above plus the existing `reaction-bounce`.

The highest-risk item is the `.theme-ready` gate: the inline script in `base.html` line 7 already handles `localStorage` → `data-theme` initialization. Extending that script with a `requestAnimationFrame` call is the correct path — inserting a second inline script would be messy and the consolidation was explicitly called out as the preferred approach in the UI-SPEC.

**Primary recommendation:** Write all four CSS blocks in sequence at the bottom of `main.css` under a clearly labeled `/* === Animations & Transitions === */` section header, and extend the existing inline `<head>` script in `base.html` with the `requestAnimationFrame(.theme-ready)` call.

---

## Standard Stack

### Core

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| CSS `@keyframes` + `animation` | CSS3 / all browsers | Entry animations | Native — no dependency, universally supported |
| CSS `transition` | CSS3 / all browsers | Dark mode color blend | Native — no dependency, universally supported |
| `prefers-reduced-motion` media query | CSS Media Queries Level 5 | Accessibility guard | WCAG 2.1 SC 2.3.3 (AAA) and 2.3.1 (A) — industry standard |
| `requestAnimationFrame` | Web API | `.theme-ready` gate timing | Fires after first paint, exactly what is needed |

### No New Dependencies

Per REQUIREMENTS.md "Out of Scope" and STATE.md decisions, JavaScript animation libraries are explicitly out of scope for this milestone. This phase introduces zero new npm/Go module dependencies.

### Installation

```bash
# No installation required — pure CSS and native browser APIs
```

---

## Architecture Patterns

### Recommended CSS Block Structure (append to main.css)

```
web/static/main.css
└─ [existing content, lines 1-677]
   └─ /* === Animations & Transitions === */  (new — appended)
      ├─ @keyframes fade-in
      ├─ .main-content fade-in rule
      ├─ .card-grid .post-card:nth-child(1-6) stagger rules
      ├─ .theme-ready transition rules (all themed surfaces)
      └─ @media (prefers-reduced-motion: reduce) { ... }  (LAST — overrides everything above)
```

### Pattern 1: Opacity-Only Fade-In Keyframe

**What:** A `@keyframes fade-in` that animates from `opacity: 0` to `opacity: 1`. No transform, no translate. Applied to `.main-content` and reused for card stagger.

**When to use:** Entry animations where the goal is "alive" not "dramatic." Opacity-only is hardware-accelerated on the compositor thread — no layout or paint cost.

**Why opacity-only matters:** D-01 explicitly excludes translate/slide. Opacity changes are composited separately from layout, so they never cause reflow. This is the correct performance choice for a blog.

```css
/* Source: MDN Web Docs - @keyframes, CSS Animations spec */
@keyframes fade-in {
  from { opacity: 0; }
  to   { opacity: 1; }
}

.main-content {
  animation: fade-in 250ms ease-out both;
}
```

The `both` fill-mode ensures the element starts at `opacity: 0` before the animation fires (prevents a flash of the final state before the keyframe begins).

### Pattern 2: CSS nth-child Stagger (no JavaScript)

**What:** Assign increasing `animation-delay` values using `:nth-child` selectors. Cap at 6. Cards beyond 6 get the animation with no delay.

**When to use:** Any list where items should enter sequentially. The cap prevents the 7th+ card from having a long wait.

```css
/* Source: MDN Web Docs - :nth-child, animation-delay */

/* Base: all cards get the animation (cards 7+ get no delay) */
.card-grid .post-card {
  animation: fade-in 250ms ease-out both;
}

/* Stagger: first 6 only */
.card-grid .post-card:nth-child(1) { animation-delay: 0ms;   }
.card-grid .post-card:nth-child(2) { animation-delay: 75ms;  }
.card-grid .post-card:nth-child(3) { animation-delay: 150ms; }
.card-grid .post-card:nth-child(4) { animation-delay: 225ms; }
.card-grid .post-card:nth-child(5) { animation-delay: 300ms; }
.card-grid .post-card:nth-child(6) { animation-delay: 375ms; }
```

**Important:** The base `.card-grid .post-card` rule must set `animation-delay: 0ms` (or omit it, which defaults to 0ms). The `:nth-child` rules override. Cards 7+ inherit `0ms` from the base rule — they are NOT excluded from the keyframe, they just have no delay.

**Why not use CSS custom properties for delay?** Go's `html/template` renders the card list via `{{range .Posts}}` without index access by default. There is no `:nth-child` CSS variable injection from the server. Pure CSS `:nth-child` selectors are the correct zero-JS approach.

### Pattern 3: .theme-ready Gate for Dark Mode Transitions

**What:** CSS transition rules on themed surfaces are scoped behind the `.theme-ready` class on `<html>`. This class is added only after first paint, so the page-load `data-theme` initialization never triggers the transitions.

**When to use:** Any site with dark mode that reads `localStorage` on load. Without the gate, the browser fires transitions during the initial render when `data-theme="dark"` is set, causing a brief flash or color sweep.

```css
/* Source: Known pattern for preventing dark mode FOUC with CSS transitions */

/* Themed surfaces — transitions enabled only after .theme-ready is set */
.theme-ready body,
.theme-ready .site-nav,
.theme-ready .post-card,
.theme-ready .site-footer,
.theme-ready .dark-toggle,
.theme-ready .nav-link,
.theme-ready .reaction-btn,
.theme-ready .tag-pill,
.theme-ready .toc,
.theme-ready .post-body pre,
.theme-ready .post-body p code,
.theme-ready .post-body li code {
  transition: background-color 250ms ease, color 250ms ease, border-color 250ms ease;
}
```

**Specificity note:** The `.post-card` already has `transition: transform 200ms ease, box-shadow 200ms ease` (main.css line 191). The `.theme-ready .post-card` rule adds to the transition list — CSS `transition` properties on the same element merge when specified on separate selectors at different specificity levels, BUT explicit transition shorthand on a higher-specificity selector will override. Use a *comma-extended* approach on `.theme-ready .post-card` to preserve the existing hover transitions:

```css
.theme-ready .post-card {
  transition:
    transform 200ms ease,
    box-shadow 200ms ease,
    background-color 250ms ease,
    color 250ms ease,
    border-color 250ms ease;
}
```

This is the critical integration point. If you use the blanket `.theme-ready *` approach, specificity will likely clobber the existing `.post-card` hover transitions. Use per-element rules.

### Pattern 4: .theme-ready Script (extend existing inline script)

**What:** Extend the existing `<script>` tag in `base.html` line 7 with a `requestAnimationFrame` callback that adds `.theme-ready` to `<html>`.

**When to use:** The existing script already runs before CSS loads (it is before the `<link>` element). The `requestAnimationFrame` call inside it defers the `.theme-ready` class addition until after first paint.

```html
<!-- Source: MDN Web Docs - requestAnimationFrame, localStorage -->
<script>(function(){
  var t = localStorage.getItem('theme');
  if (t === 'dark' || (!t && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    document.documentElement.setAttribute('data-theme', 'dark');
  }
  requestAnimationFrame(function() {
    document.documentElement.classList.add('theme-ready');
  });
})()</script>
```

This replaces the existing single-line script. The `requestAnimationFrame` fires after the browser has committed the first frame — the `data-theme` attribute is already set before any paint, so the transition gate opens only for subsequent user-initiated toggles.

### Pattern 5: prefers-reduced-motion Guard

**What:** A single `@media (prefers-reduced-motion: reduce)` block at the bottom of the new CSS section that zeros out all animation and transition values set in this phase, plus wraps the existing `reaction-bounce`.

**When to use:** Always, as the last block. Placing it last ensures maximum override specificity.

```css
/* Source: MDN Web Docs - prefers-reduced-motion, WCAG 2.1 */
@media (prefers-reduced-motion: reduce) {
  .main-content {
    animation: none;
  }

  .card-grid .post-card,
  .card-grid .post-card:nth-child(1),
  .card-grid .post-card:nth-child(2),
  .card-grid .post-card:nth-child(3),
  .card-grid .post-card:nth-child(4),
  .card-grid .post-card:nth-child(5),
  .card-grid .post-card:nth-child(6) {
    animation: none;
    animation-delay: 0ms;
  }

  .theme-ready body,
  .theme-ready .site-nav,
  .theme-ready .post-card,
  /* ... all themed surfaces ... */ {
    transition: none;
  }

  /* Fix existing unguarded reaction-bounce (main.css line 649) */
  .reaction-btn.bounce .reaction-icon {
    animation: none;
  }
}
```

### Anti-Patterns to Avoid

- **`@property` for dark mode color transitions:** CSS `@property` allows animating custom properties directly, but browser support as of 2026 is Chromium-only (Firefox and Safari have partial support). Standard `transition` on concrete properties (`background-color`, `color`, `border-color`) has universal support. Given the small number of selectors involved, standard transitions are the correct approach. CONTEXT.md explicitly marks `@property` as Claude's discretion — choose standard transitions.

- **Second inline `<script>` tag for `.theme-ready`:** Adding a second `<script>` block after the CSS link defeats the purpose — it would run after CSS loads, not before. Extend the existing script (line 7 of `base.html`).

- **`.theme-ready *` wildcard selector:** This clobbers existing transitions (e.g., `.post-card` hover lift). Use explicit per-element selectors.

- **`animation-fill-mode: forwards` without `backwards`:** Using `both` is correct. `forwards` alone does not prevent a pre-animation flash because the element renders at its natural state before the first keyframe fires. `both` applies `from` before the delay expires.

- **Missing `opacity: 1` visibility on reduced-motion:** When `animation: none` is set, the element renders at its natural `opacity: 1` (since the keyframe is gone). However, if `.main-content` had a static `opacity: 0` applied elsewhere, reduced-motion users would see a blank page. Confirm no static opacity is set on `.main-content`.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Flash-of-wrong-theme | Custom debounce/delay script | `requestAnimationFrame` in existing `<head>` script | rAF fires after first paint by definition — timing is exact |
| Card entrance sequencing | JavaScript `IntersectionObserver` + `setTimeout` stagger | CSS `:nth-child` + `animation-delay` | Zero JS, composited on GPU, no scroll event overhead |
| Dark mode color animation | JavaScript class toggling + `setTimeout` cleanup | CSS `transition` gated by `.theme-ready` | CSS transitions are declarative and GPU-composited |

**Key insight:** This entire phase is achievable without touching `main.js` at all. The existing dark mode toggle already writes `data-theme` correctly. CSS transitions respond to attribute changes without JS involvement. The only JS change is one extra line in the existing inline `<head>` script.

---

## Common Pitfalls

### Pitfall 1: Post-Card Transition Clobbering

**What goes wrong:** Adding a new `transition` shorthand on `.post-card` inside `.theme-ready` silently removes the existing `transform` and `box-shadow` transitions, breaking the hover lift animation.

**Why it happens:** CSS `transition` shorthand replaces the entire transition list — it does not append. The existing rule at line 191 sets `transition: transform 200ms ease, box-shadow 200ms ease`. A new rule targeting `.theme-ready .post-card { transition: background-color 250ms ease, ... }` would win by specificity and delete the hover transitions.

**How to avoid:** On `.theme-ready .post-card`, spell out ALL five transition values: `transform 200ms ease, box-shadow 200ms ease, background-color 250ms ease, color 250ms ease, border-color 250ms ease`.

**Warning signs:** Hover over a post card in browser DevTools after implementation — the card should still lift. If it doesn't lift, the transition clobber has occurred.

### Pitfall 2: animation-fill-mode Flash

**What goes wrong:** Cards appear at full opacity for a frame before the animation starts, then briefly disappear to `opacity: 0`, then fade in. Visible as a flicker.

**Why it happens:** If `animation-fill-mode` is not set to `both` (or at minimum `backwards`), the element renders at its natural state before the `from` keyframe is applied. With `animation-delay` on cards 2–6, this window is 75ms–375ms — perceptible.

**How to avoid:** Always specify `animation: fade-in 250ms ease-out both` — the `both` keyword is the fourth value in the shorthand (name duration timing-function fill-mode).

**Warning signs:** Cards briefly flash visible, then fade in. Most visible on card 6 which has a 375ms delay.

### Pitfall 3: Dark Mode Flash on Hard Reload

**What goes wrong:** Returning visitor sees a brief white flash before dark mode is applied on page load.

**Why it happens:** The inline `<script>` block runs synchronously before CSS parses, but if the script is placed AFTER the `<link rel="stylesheet">` tag, CSS has already begun computing styles with default (light) values.

**How to avoid:** The inline script MUST be before the `<link rel="stylesheet">` element. The existing `base.html` already has this correct order (script on line 7, CSS link on line 11). Preserve this order.

**Warning signs:** Open DevTools > Network tab, throttle to Slow 3G, hard-reload in dark mode. A white flash before dark background appears indicates the script is too late.

### Pitfall 4: Reduced-Motion Guard Incompleteness

**What goes wrong:** The existing `reaction-bounce` at line 649–658 is currently unguarded. If the new `prefers-reduced-motion` block only wraps the new animations, the existing bounce still fires for reduced-motion users.

**Why it happens:** The `reaction-bounce` keyframe predates this phase. Easy to overlook when writing a new animation section.

**How to avoid:** The `@media (prefers-reduced-motion: reduce)` block must include `.reaction-btn.bounce .reaction-icon { animation: none; }` alongside the new animation guards.

**Warning signs:** Set `prefers-reduced-motion: reduce` in OS settings, click the thumbs-up reaction button on a post — if the icon bounces, the guard is missing.

### Pitfall 5: .main-content Selector Does Not Exist

**What goes wrong:** The fade-in animation is applied to `.main-content` but that class is not present in the HTML templates.

**Why it happens:** Looking at `list.html`, the content block wraps in `<div class="container">` — there is no `.main-content` wrapper. The CONTEXT.md says "`.main-content` (or equivalent content wrapper)."

**How to avoid:** Audit the templates before writing selectors. Options: (a) use `.container` as the target (already present in all page templates), or (b) add a `.main-content` wrapper to `base.html`'s `{{block "content"}}` area. Option (a) is simpler — no template change needed. Option (b) requires touching `base.html` but gives a semantically cleaner animation hook.

**Warning signs:** Animation defined in CSS but no visible fade-in in browser. DevTools will show the rule is parsed but matching zero elements.

---

## Code Examples

Verified patterns from official sources / existing codebase:

### Existing Script to Extend (base.html line 7 — current)

```html
<script>(function(){var t=localStorage.getItem('theme');if(t==='dark'||(!t&&window.matchMedia('(prefers-color-scheme: dark)').matches)){document.documentElement.setAttribute('data-theme','dark');}})()</script>
```

### Extended Script (with .theme-ready gate)

```html
<!-- Source: MDN requestAnimationFrame, extending existing base.html script -->
<script>(function(){
  var t=localStorage.getItem('theme');
  if(t==='dark'||(!t&&window.matchMedia('(prefers-color-scheme: dark)').matches)){
    document.documentElement.setAttribute('data-theme','dark');
  }
  requestAnimationFrame(function(){
    document.documentElement.classList.add('theme-ready');
  });
})()</script>
```

This can be kept as a single minified line if desired — the content is what matters, not formatting.

### Complete Animation CSS Block (to append to main.css)

```css
/* === Animations & Transitions === */

/* --- Entry Animation Keyframe --- */

@keyframes fade-in {
  from { opacity: 0; }
  to   { opacity: 1; }
}

/* --- Page Load Fade-In (ANIM-02) --- */

.container {
  animation: fade-in 250ms ease-out both;
}

/* --- Post Card Stagger (ANIM-03) --- */

.card-grid .post-card {
  animation: fade-in 250ms ease-out both;
}

.card-grid .post-card:nth-child(1) { animation-delay: 0ms;   }
.card-grid .post-card:nth-child(2) { animation-delay: 75ms;  }
.card-grid .post-card:nth-child(3) { animation-delay: 150ms; }
.card-grid .post-card:nth-child(4) { animation-delay: 225ms; }
.card-grid .post-card:nth-child(5) { animation-delay: 300ms; }
.card-grid .post-card:nth-child(6) { animation-delay: 375ms; }

/* --- Dark Mode Color Transitions — .theme-ready Gate (ANIM-04, ANIM-05) --- */

.theme-ready body,
.theme-ready .site-nav,
.theme-ready .site-footer,
.theme-ready .dark-toggle,
.theme-ready .nav-link,
.theme-ready .reaction-btn,
.theme-ready .tag-pill,
.theme-ready .toc,
.theme-ready .post-body pre,
.theme-ready .post-body p code,
.theme-ready .post-body li code {
  transition: background-color 250ms ease, color 250ms ease, border-color 250ms ease;
}

/* Post card: preserve existing hover transitions alongside new color transitions */
.theme-ready .post-card {
  transition:
    transform 200ms ease,
    box-shadow 200ms ease,
    background-color 250ms ease,
    color 250ms ease,
    border-color 250ms ease;
}

/* --- Reduced Motion Guard (ANIM-01) — MUST be last --- */

@media (prefers-reduced-motion: reduce) {
  .container {
    animation: none;
  }

  .card-grid .post-card,
  .card-grid .post-card:nth-child(1),
  .card-grid .post-card:nth-child(2),
  .card-grid .post-card:nth-child(3),
  .card-grid .post-card:nth-child(4),
  .card-grid .post-card:nth-child(5),
  .card-grid .post-card:nth-child(6) {
    animation: none;
    animation-delay: 0ms;
  }

  .theme-ready body,
  .theme-ready .site-nav,
  .theme-ready .post-card,
  .theme-ready .site-footer,
  .theme-ready .dark-toggle,
  .theme-ready .nav-link,
  .theme-ready .reaction-btn,
  .theme-ready .tag-pill,
  .theme-ready .toc,
  .theme-ready .post-body pre,
  .theme-ready .post-body p code,
  .theme-ready .post-body li code {
    transition: none;
  }

  /* Existing reaction-bounce — currently unguarded at line 649 */
  .reaction-btn.bounce .reaction-icon {
    animation: none;
  }
}
```

**Note on `.container` vs `.main-content`:** The example above uses `.container` because that is the actual wrapping class present in all page templates. If the plan adds a `.main-content` wrapper div instead, substitute that selector. The tradeoff: `.container` is simpler (no template edit) but also wraps the content on non-list pages. `.main-content` is semantically cleaner but requires one template change.

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| `prefers-reduced-motion` ignored | Required per WCAG 2.1 SC 2.3.3 | 2018 (WCAG 2.1) | Reduced-motion guard is non-optional for accessible sites |
| Inline cookie-based theme init | `localStorage` + `data-theme` attribute | Established pattern | Current codebase already uses this correctly |
| CSS `@property` for custom property animation | Standard `transition` on concrete properties | `@property` still limited browser support (2026) | Standard transitions preferred for this use case |
| JavaScript-driven stagger (setTimeout) | Pure CSS `:nth-child` + `animation-delay` | CSS animations mature ~2015 | No JS overhead, GPU-composited |

**Deprecated/outdated:**
- `CSS animation-fill-mode` omission: Older tutorials skip `fill-mode: both` — this causes the flash described in Pitfall 2. Always include `both`.

---

## Open Questions

1. **`.container` vs `.main-content` as fade-in target**
   - What we know: `.container` exists in all templates. `.main-content` does not exist yet.
   - What's unclear: Whether fading `.container` is visually correct — it also wraps nav-adjacent content on some pages. The post page and about page both use `.container` directly.
   - Recommendation: Use `.container` to avoid a template change. If it fades elements the designer doesn't want animated, a `.main-content` wrapper can be added in the same plan as a Wave 0 task.

2. **Whether `.toc` and `post-body` code elements need dark mode transitions**
   - What we know: These elements use `var(--color-surface)` and `var(--color-divider)` which change with `data-theme`.
   - What's unclear: Whether the visual change on these smaller elements is noticeable enough to warrant explicit transition rules.
   - Recommendation: Include them in the selector list. The transition cost is negligible and visual consistency is better.

---

## Environment Availability

Step 2.6: SKIPPED — this phase makes no external tool calls, installs no packages, and requires no databases or services. All changes are to static CSS and an HTML template. The existing Go server serves these files unchanged.

---

## Validation Architecture

### Test Framework

| Property | Value |
|----------|-------|
| Framework | Go `testing` (stdlib) + `net/http/httptest` |
| Config file | none — standard `go test ./...` |
| Quick run command | `go test ./internal/handler/blog/... -run TestList -v` |
| Full suite command | `go test ./... -count=1` |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| ANIM-01 | CSS contains `prefers-reduced-motion: reduce` block covering all animations | Manual — CSS inspection | n/a | n/a |
| ANIM-02 | `.container` (or `.main-content`) has `animation: fade-in` in CSS | Manual — browser visual | n/a | n/a |
| ANIM-03 | `.card-grid .post-card:nth-child(1-6)` have `animation-delay` values | Manual — browser DevTools | n/a | n/a |
| ANIM-04 | Dark mode toggle produces color blend, not instant switch | Manual — browser visual | n/a | n/a |
| ANIM-05 | No white flash on dark-mode hard reload | Manual — browser network throttle | n/a | n/a |

**Note:** These are all CSS/visual behaviors. Go unit tests cannot verify CSS properties or browser rendering. All five requirements are manual-verification items.

### Sampling Rate

- **Per task commit:** `go test ./... -count=1` — ensure no Go compilation breaks from template changes
- **Per wave merge:** `go test ./... -count=1`
- **Phase gate:** Go tests green + manual browser checklist before `/gsd:verify-work`

### Wave 0 Gaps

The existing Go test suite covers handler and template rendering — no new test files are needed for this phase. CSS animation validation is inherently manual.

- [ ] Manual browser checklist (not a test file gap — these are manual steps for the verifier):
  - [ ] ANIM-02: Observe `.container` fade-in on hard reload (light mode)
  - [ ] ANIM-03: Observe card stagger on list page with 2+ posts
  - [ ] ANIM-04: Toggle dark mode — verify smooth color blend
  - [ ] ANIM-05: Hard reload in dark mode — verify no white flash (throttle to Slow 3G in DevTools)
  - [ ] ANIM-01: Set `prefers-reduced-motion: reduce` in OS — verify no animations on reload or toggle
  - [ ] ANIM-01: Click reaction button with reduced motion — verify no bounce

---

## Sources

### Primary (HIGH confidence)

- MDN Web Docs — CSS `@keyframes`, `animation`, `animation-fill-mode`, `animation-delay`: fundamental CSS spec
- MDN Web Docs — CSS `transition` shorthand: property merging behavior, specificity rules
- MDN Web Docs — `prefers-reduced-motion` media query: browser support universal as of 2021+
- MDN Web Docs — `requestAnimationFrame`: fires after first paint, exact semantics confirmed
- Codebase — `web/static/main.css` lines 191, 633, 649–658: existing transition/animation patterns verified by direct read
- Codebase — `web/templates/base.html` line 7: existing inline script verified by direct read
- Codebase — `web/templates/list.html`: `.container` confirmed as the wrapping element, no `.main-content` class present

### Secondary (MEDIUM confidence)

- CSS `@property` browser support assessment: based on known Chromium-only status for registered custom properties animation as of late 2025; standard transitions confirmed as the correct alternative
- WCAG 2.1 SC 2.3.3 (Animation from Interactions, Level AAA): `prefers-reduced-motion` guard is the standard implementation

### Tertiary (LOW confidence)

- None — all findings for this phase are verifiable via the codebase directly or established CSS spec behavior

---

## Project Constraints (from CLAUDE.md)

| Directive | Applies to Phase 10? |
|-----------|---------------------|
| Go with minimal dependencies — avoid large frameworks | Yes — no new dependencies introduced |
| All persistent data lives on EBS volume at `/var/www/html` | Not applicable — CSS/template phase |
| Docker container on port 8080 | Not applicable — no server changes |
| Design: use `frontend-design` skill | Note: `.claude/skills/` directory was not found in this project. No skill directory exists. Research applied general CSS animation best practices. |
| GSD workflow enforcement before file edits | Yes — planner will produce tasks through GSD |
| No new CSS files for v1.2 — all changes to `main.css` | Yes — all CSS appended to `main.css` |
| CSS changes before template changes | Enforced by wave ordering in the plan |
| Dark mode transitions gated behind `.theme-ready` | Yes — core requirement of this phase |
| `prefers-reduced-motion` guard ships alongside new animations AND fixes existing `reaction-bounce` | Yes — ANIM-01 |

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — pure CSS/browser APIs, no library choices to get wrong
- Architecture: HIGH — patterns verified against existing codebase file contents; no assumptions about what classes exist
- Pitfalls: HIGH — each pitfall is traceable to a specific existing code line or confirmed CSS behavior
- Transition specificity concern: HIGH — post-card transition clobber is a well-known CSS shorthand gotcha, verified against current line 191

**Research date:** 2026-03-28
**Valid until:** 2026-09-28 (CSS spec is stable; `@property` browser support is the only area to recheck if timelines shift)
