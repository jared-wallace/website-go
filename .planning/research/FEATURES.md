# Feature Research — v1.2 Shore Leave Polish

**Domain:** CSS/template design polish — nautical-themed personal blog
**Researched:** 2026-03-28
**Confidence:** HIGH (CSS techniques verified against MDN, css-tricks.com, joshwcomeau.com; existing codebase reviewed directly)

---

## Scope Note

This research covers the five design polish areas for "The Wild Meridian" v1.2 milestone. The existing CSS architecture (`web/static/main.css`) uses CSS custom properties for theming, a flat single-file structure, and `[data-theme="dark"]` for dark mode. All recommendations are constrained to CSS/template changes — no Go backend changes.

**Existing CSS comment to fix:** Line 2 reads `Design System — The Log` — this is the stale rebrand artifact. One-line fix.

---

## Area 1: Footer Navigation and Personality

### What Well-Designed Personal Blogs Do

Strong personal blog footers do three things simultaneously: provide utility navigation, signal personality (a phrase reinforcing brand voice), and close the page gracefully without competing with the content above. The footer is the last thing a reader sees — it should feel like the bartender saying goodnight, not a legal disclaimer.

The pattern that works for personal/editorial sites: a thin horizontal split with nav links on one side, utility items (copyright, RSS) on the other, and an optional personality line sitting quietly in between or below.

### Table Stakes

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| About link in footer | Readers who reach the footer want to know who wrote this | LOW | Move from `<nav>` to `<footer>` in `base.html`; remove from `.nav-link` |
| Copyright line | Legal expectation; already exists | LOW | `© {{.Year}} Jared Wallace` already present |
| RSS link | Readers who care about RSS look in footers | LOW | Already present; keep |
| Two-column footer layout | Separates nav from utility; feels considered | LOW | Left: nav links (About). Right: copyright + RSS. `display: flex; justify-content: space-between` on `.site-footer` |

### Differentiators

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| Personality tagline in footer | Reinforces brand voice; memorable last impression | LOW | One short italic line in `--font-body`, `--color-text-muted` — the existing "dispatches from the deep end" works, or a quieter nautical closer |
| Footer rope divider full opacity | Reinforces nautical motif as a section break | LOW | Existing rope divider already in footer; remove `opacity: 0.6` to let dashes read fully |

### Anti-Features

| Feature | Why Requested | Why Problematic | Alternative |
|---------|---------------|-----------------|-------------|
| Social media links | "Every footer has them" | Dilutes the handcrafted feel; links readers away | Omit entirely — a personal blog with no social links is a statement, not an oversight |
| Newsletter subscribe form | Engagement growth | No email system exists; RSS already covers it | RSS link is sufficient |
| Sitemap link | SEO advice from general-purpose guides | Machines read `/sitemap.xml` directly; surfacing it as a human link adds noise | Leave sitemap to `<head>` autodiscovery |

### CSS Dependency

Existing `.site-footer` uses `text-align: center` and a single centered `<p>`. Switching to a two-column layout requires changing `.site-footer` to `display: flex; justify-content: space-between; align-items: center` and removing `text-align: center`. The inner `<p>` and `.rss-link` need restructuring into two sibling elements. Template change required in `base.html`.

---

## Area 2: Background Texture for Weathered/Vintage Feel

### What Works

The canonical CSS-only grain technique: an SVG `feTurbulence` filter encoded as a data URI, applied as a pseudo-element overlay. Zero additional files, zero JavaScript, zero dependencies. Verified against ibelick.com and css-tricks.com — the technique is well-documented and works across all modern browsers.

**The exact pattern for this codebase:**

```css
body::before {
  content: "";
  position: fixed;
  inset: 0;
  pointer-events: none;
  z-index: 0;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 600 600'%3E%3Cfilter id='a'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='.65' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23a)'/%3E%3C/svg%3E");
  background-repeat: repeat;
  background-size: 182px;
  opacity: 0.035;
}

[data-theme="dark"] body::before {
  opacity: 0.06;
}
```

Dark mode requires a higher opacity because dark surfaces absorb grain faster than light ones — the same 0.035 opacity is nearly invisible on `#1A1F2E`.

**Key tuning parameters:**
- `baseFrequency`: `0.65` produces fine grain; lower (`0.4`) = coarser, more tactile
- `numOctaves`: `3` is standard; `4` adds complexity at minimal GPU cost
- `opacity`: Stay under `0.06` in light mode; under `0.10` in dark mode. Anything above `0.12` reads as a rendering artifact, not a design choice

### Table Stakes

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| Subtle grain on `body` | "Weathered beach bar" requires texture; flat `#F5F0E8` reads as default web | LOW | Single CSS pseudo-element; no image files |
| Dark mode opacity adjustment | Dark surfaces absorb grain; needs compensation | LOW | `[data-theme="dark"] body::before { opacity: 0.06; }` |

### Differentiators

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| Slightly higher baseFrequency for cards | Cards feel like a different material — more worn, rougher texture | LOW | Apply a second pseudo-element variant to `.post-card::before` with `baseFrequency='.80'` and `opacity: 0.04` |

### Anti-Features

| Feature | Why Requested | Why Problematic | Alternative |
|---------|---------------|-----------------|-------------|
| Heavy grain (opacity > 0.12) | More dramatic atmosphere | Looks like a GPU artifact; makes text harder to read; breaks legibility | Tune `baseFrequency` for more texture without increasing opacity |
| External grain PNG/WebP | Easier to adjust in design tools | Extra HTTP request; does not adapt to dark mode; more bytes on every pageload | SVG data URI is inline, adaptive, zero-dependency |
| `filter: contrast() brightness()` boost | CSS-Tricks shows this as an option for more vivid grain | Causes color banding on long gradients; browser rendering inconsistency between Blink and WebKit | Opacity tuning is more reliable |

### CSS Dependency

The existing `body` has no `::before` or `::after`. No z-index conflicts: all content is either at `z-index: auto` or `z-index: 100` (nav). A `position: fixed; inset: 0; z-index: 0; pointer-events: none` pseudo-element is safe. No template changes needed.

---

## Area 3: Card Depth and Visual Weight

### What Works

Per Josh W. Comeau's research (joshwcomeau.com/css/designing-shadows), the key insight is that a single `box-shadow` always looks synthetic. Physical objects cast two overlapping shadows simultaneously: a short crisp contact shadow (key light) and a diffuse ambient shadow (fill light). Two-layer shadows are the minimum for a natural result.

**For a blog card at rest — "paper resting on warm wood":**

```css
box-shadow:
  0 1px 2px rgba(44, 36, 24, 0.08),
  0 2px 8px rgba(44, 36, 24, 0.06);
```

**On hover — card lifts off surface:**

```css
box-shadow:
  0 4px 8px rgba(44, 36, 24, 0.10),
  0 8px 24px rgba(44, 36, 24, 0.12);
```

The existing CSS has no resting shadow — cards are flat until hover. Adding a resting shadow gives them weight and makes the hover lift feel like picking up an object rather than a CSS trick.

**Shadow hue matters.** The existing background `--color-bg: #F5F0E8` is warm brown-beige. Using `rgba(44, 36, 24, ...)` — the RGB equivalent of `--color-text` (`#2C2418`) — produces warm shadows that sit naturally against the surface. Pure black (`rgba(0,0,0,...)`) on warm backgrounds reads as grey and disconnected.

**Dark mode:** Warm hue matching is less critical on dark surfaces. Use `rgba(0, 0, 0, 0.25)` resting, `rgba(0, 0, 0, 0.45)` hover — dark mode can handle higher opacity because background contrast absorbs it.

### Rope Divider Strengthening

The existing rope divider: `border-top: 2px dashed var(--color-divider); opacity: 0.6`. The `opacity: 0.6` is doing the wrong job — it fades the motif rather than integrating it. Removing the opacity fade and letting the dashes show at full `--color-divider` color makes the nautical motif visible without requiring any template changes.

The 2px thickness is correct — 1px dashes can disappear on high-DPI screens.

### Table Stakes

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| Resting shadow on cards | Cards with no resting state feel flat; hover lift has no context | LOW | Two-layer shadow, warm hue; no template changes |
| Warm shadow hue in light mode | `rgba(0,0,0,...)` shadows look grey/cold on warm `#F5F0E8` background | LOW | Use `rgba(44, 36, 24, ...)` matching `--color-text` |
| Rope divider at full opacity | Dashed border at `opacity: 0.6` is too subtle to read as a motif | LOW | Remove `opacity: 0.6` from `.rope-divider` — one-line CSS change |

### Differentiators

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| Distinct resting vs hover elevation | Makes the lift feel physical; shadow grows as card rises | LOW | Different two-layer values at rest vs hover (hover already has `translateY(-4px)` — add matching shadow increase) |
| Left-border accent peek on card hover | Nautical bar-ledge feel; color accent visible on focus | LOW | `border-left: 3px solid var(--color-accent)` on `.post-card:hover` instead of the uniform `1px solid var(--color-divider)` |

### Anti-Features

| Feature | Why Requested | Why Problematic | Alternative |
|---------|---------------|-----------------|-------------|
| Spread value > 0 on `box-shadow` | Creates "glow" or "halo" effect | Looks too digital for a weathered aesthetic | Use offset + blur only; zero spread |
| `filter: drop-shadow()` instead of `box-shadow` | Sometimes recommended for irregular shapes | Cannot be transitioned cleanly alongside `transform` — causes compositing layer promotion issues | Stay with `box-shadow` which already has a transition wired in existing CSS |

### CSS Dependency

The existing `.post-card` already has `transition: transform 200ms ease, box-shadow 200ms ease` — the shadow transition is pre-wired. Adding resting shadow is a one-line change. No template changes needed for resting shadow or rope divider fixes.

---

## Area 4: Page Entry Animations

### What Feels Organic vs Corporate

**Corporate:** Large Y offsets (`translateY(40px)+`), bounce/spring easings, durations over 500ms, stagger delays over 100ms, animations that fire on every page load without user preference detection. SaaS landing page energy.

**Organic:** Small Y offsets (12–20px), `ease-out` easing (fast arrival, slow settle), 250–350ms duration, stagger delays 40–70ms between cards, opacity as the primary effect with Y offset as secondary. The Cloud Four stagger study uses 0.025s increments as a reference point — that maps to 25ms, which at 300ms duration feels like birds settling on a wire rather than a loading spinner.

**Implementation for post card stagger (pure CSS, no JS):**

```css
@keyframes card-enter {
  from { opacity: 0; transform: translateY(14px); }
  to   { opacity: 1; transform: translateY(0); }
}

.post-card {
  animation: card-enter 300ms ease-out both;
}

.post-card:nth-child(1) { animation-delay:   0ms; }
.post-card:nth-child(2) { animation-delay:  50ms; }
.post-card:nth-child(3) { animation-delay: 100ms; }
.post-card:nth-child(4) { animation-delay: 150ms; }
.post-card:nth-child(5) { animation-delay: 180ms; }
.post-card:nth-child(6) { animation-delay: 210ms; }
```

Six rules cover a 2-column grid of 3 rows. Beyond that the stagger is imperceptible anyway — cap delay at ~250ms total.

**`prefers-reduced-motion` is mandatory, not optional polish.** WCAG 2.1 SC 2.3.3 (Animation from Interactions) requires motion animations to be suppressible. Per MDN, the value is `reduce`, not "none" — meaning you can still use fade-only transitions, just not translations.

```css
@media (prefers-reduced-motion: reduce) {
  .post-card {
    animation: none;
  }
  body {
    transition: none;
  }
}
```

**Dark mode color transitions:** The existing mode toggle fires instantly — no CSS transition. Adding smooth theme transitions requires `transition: background-color 200ms ease, color 200ms ease, border-color 200ms ease` on `body`. The gotcha: this transition must not be active during the initial theme-on-load script in `base.html`, or every page load will flash an animation as the theme applies.

The standard fix: enable transitions via a `.transitions-ready` class added by `main.js` after the initial paint settles, or by wrapping the transition declarations in a selector that only applies after the inline script has run. The existing `base.html` inline script sets `data-theme` before paint — the transition should only be added by `main.js` on `DOMContentLoaded`.

### Table Stakes

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| `prefers-reduced-motion` guard | Accessibility requirement; vestibular disorders affect 70M+ people | LOW | Single `@media` block wrapping all `animation` and motion `transition` declarations |
| No flash-of-animation on initial load | Dark mode loads via inline script; transitions must not fire on load | LOW | Add transition via JS after initial paint; not in base CSS |

### Differentiators

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| Card stagger on list page load | Cards settling into view — organic, tactile, memorable | LOW | Pure CSS `nth-child` delay rules; no JS; 6 rules covers a full grid page |
| Dark mode color transitions | Theme switch feels smooth, considered | LOW | `transition` on `body` enabled via JS after initial paint |
| Post page fade-in | Post content arriving like a page turning | LOW | Single `animation: card-enter 300ms ease-out both` on `.post-page` container |

### Anti-Features

| Feature | Why Requested | Why Problematic | Alternative |
|---------|---------------|-----------------|-------------|
| Scroll-triggered animations (IntersectionObserver) | "Modern" | Adds JS complexity; blog posts have minimal content requiring scroll-reveal; the fold is not the problem on a blog listing | CSS load animations are sufficient |
| Spring/bounce easing (`cubic-bezier` with values > 1) | "Alive" feel | Reads as app-like, not aged/weathered — overshoot belongs on Notion, not a beach bar chalkboard | `ease-out` with 14px offset achieves the same sense of arrival without the bounce |
| Parallax scrolling | Nautical depth metaphor | Performance regression on mobile; vestibular trigger; fights `prefers-reduced-motion` | Background grain provides depth without motion |
| Skeleton loading screens | Professional polish | Zero value when page loads from server-rendered Go templates in under 100ms on EBS | Server rendering eliminates the problem skeleton screens solve |

### CSS Dependency

The existing CSS has one `@keyframes` block (`reaction-bounce`). Adding `card-enter` is purely additive. The dark mode transition requires a JS change to `main.js` (adding `.transitions-ready` class post-paint) and a CSS change scoping transitions to `.transitions-ready body`. This is the most involved change in v1.2 — low risk, but requires touching both CSS and JS.

---

## Area 5: Tag Pill Visibility at Small Sizes

### What Is Wrong Now

The existing `.tag-pill`:
- `padding: 2px 4px` — vertically too tight at 14px font-size; label feel, not pill feel
- `background: transparent` — three layers of low contrast compound: transparent + muted border + muted text
- In dark mode especially: `--color-text-muted: #9A8B75` on a transparent surface over `--color-surface: #242B3D` = barely legible

Per Smart Interface Design Patterns research: at sizes under 16px, pills require either a filled background or a high-contrast border. Relying on text color alone at small sizes fails. The weathered theme can use a warm fill without departing from the palette.

**The fix:**

```css
.tag-pill {
  padding: 3px 8px;                       /* was 2px 4px */
  border: 1px solid var(--color-divider);
  border-radius: 3px;                     /* was 2px — minor */
  font-size: 12px;                        /* was 14px — smaller font with more padding = same height, better legibility */
  color: var(--color-text-muted);
  background-color: var(--color-surface); /* was transparent */
  line-height: 1.4;
}
```

The `font-size: 12px` with larger padding maintains a similar overall pill height while actually improving readability — this matches the Bootstrap `badge-sm` / Tailwind `badge-xs` sizing pattern where small badges explicitly downsize font rather than shrinking padding.

**Dark mode contrast check:** `--color-surface: #242B3D` against `--color-bg: #1A1F2E` background has perceptible separation. The `--color-text-muted: #9A8B75` on `#242B3D` delivers approximately 3:1 contrast ratio — adequate for decorative non-interactive labels; below the 4.5:1 needed for interactive elements.

### Table Stakes

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| Filled pill background | Transparent pills vanish on surfaces with similar hue | LOW | `background-color: var(--color-surface)` — one custom property, already defined |
| Adequate padding to feel pill-like | `2px 4px` feels like a debug label, not a designed element | LOW | `3px 8px` matches standard badge component conventions |

### Differentiators

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| Warm tag background in light mode | Tags feel like paper labels on wood — material, not digital | LOW | `background-color: var(--color-divider)` at low opacity, or a dedicated `--color-tag-bg` custom property |
| Accent border/color on tag hover (when clickable) | Confirms interactivity when tags become filterable in a future milestone | LOW | Defer until tags are actually clickable; don't add hover styles to non-interactive elements |

### Anti-Features

| Feature | Why Requested | Why Problematic | Alternative |
|---------|---------------|-----------------|-------------|
| Accent-filled tags (solid accent background) | Bold, high-visibility | Ocean blue / gold filled tags on every card creates visual noise that competes with post titles | Reserve accent fill for interactive selected/active states only |
| Tag icons or emoji | More personality | Unreadable at 12px; tags are labels, not illustrations | Clean typography and visible background deliver personality without decoration |
| Per-tag color coding | Helps distinguish categories | Requires either a fixed color map (brittle) or generated colors (inconsistent with the warm palette) | Single neutral tag style is more coherent at small blog scale |

### CSS Dependency

Purely additive changes to `.tag-pill`. No template changes. No new custom properties required (though `--color-tag-bg` is an option for easier tuning).

---

## Bonus: Homepage Heading / Hero Area

Not in the five research questions but included because it has CSS dependencies on Areas 3 and 4 above.

The existing list page has no `<h1>` above the card grid. A minimal hero establishes context without pretending to be a marketing page. Playfair Display at 36–42px, followed by a short italic subtitle in `--font-body --color-text-muted`:

```html
<div class="page-hero">
  <h1 class="page-hero-title">The Wild Meridian</h1>
  <p class="page-hero-subtitle">dispatches from the deep end</p>
</div>
```

The card stagger (Area 4) should have the hero fade in first (`animation-delay: 0`) with cards following at `50ms+` intervals. This creates a reading hierarchy: the establishment shot lands first, then the content settles in behind it. Template change in `list.html`; CSS additions for `.page-hero-title` and `.page-hero-subtitle`.

---

## Feature Dependencies

```
Background texture (Area 2)
    └── independent — CSS ::before only, no deps

Card resting shadow (Area 3)
    └── enhances card stagger (Area 4)
        shadows give cards a surface to lift off from; stagger without
        resting shadow = cards materializing from nothing

Dark mode transitions (Area 4)
    └── requires prefers-reduced-motion guard
    └── requires JS transition-enablement in main.js after initial paint
        (base.html inline script runs before CSS transitions should be active)

Tag pill fix (Area 5)
    └── independent — CSS only, no deps

Footer expansion + About move (Area 1)
    └── requires About link removed from nav in base.html
    └── independent of all other areas — do as atomic template change
```

### Critical Dependency Note

**Dark mode transitions + initial page load:** The inline `<script>` in `base.html` applies `data-theme="dark"` before the first paint. If CSS `transition` is declared on `body` unconditionally, every page load triggers a color-flash animation as the theme sets. The fix is to add transition declarations only after `DOMContentLoaded` fires in `main.js` — either via a `.transitions-ready` class or by directly appending a `<style>` element. This is the one change in v1.2 that touches both CSS and JavaScript.

**Card shadow + card stagger:** Do both together or neither. A card that animates in from `opacity: 0` but has no resting shadow looks like it materialized in space. The resting shadow gives the card a physical context before it arrives.

---

## Prioritization for v1.2

### Do First (High Impact, Low Risk — CSS Only)

| Feature | User Value | Implementation Cost | Priority |
|---------|------------|---------------------|----------|
| Background grain texture | HIGH — establishes the weathered atmosphere | LOW | P1 |
| Card resting shadow (warm hue) | HIGH — grounds the card grid physically | LOW | P1 |
| Tag pill background fill | MEDIUM — legibility fix, currently nearly invisible in dark mode | LOW | P1 |
| Rope divider full opacity | LOW — small motif improvement | LOW | P1 |
| CSS comment rebrand fix (line 2) | LOW — internal hygiene | LOW | P1 |

### Do Second (Moderate Impact, Requires Template or JS Touch)

| Feature | User Value | Implementation Cost | Priority |
|---------|------------|---------------------|----------|
| Card stagger animation + prefers-reduced-motion | MEDIUM — delight, not function | LOW | P2 |
| Dark mode color transitions (JS-gated) | MEDIUM — polish | LOW-MEDIUM | P2 |
| Footer About move + personality line | MEDIUM — reduces nav clutter; adds brand voice | LOW | P2 |
| Homepage hero heading | MEDIUM — orientation for first-time visitors | LOW | P2 |
| Reaction button radius consistency (pill → 4px) | LOW — consistency | LOW | P2 |

### Defer or Skip

| Feature | User Value | Implementation Cost | Priority |
|---------|------------|---------------------|----------|
| Rope divider template-level character upgrade | LOW | MEDIUM | P3 |
| Per-card grain texture variation | LOW — imperceptible to most readers | LOW | P3 |
| Clickable tag filtering | MEDIUM | MEDIUM (requires routing) | Out of scope v1.2 |

---

## Sources

- [Grainy Gradients — CSS-Tricks](https://css-tricks.com/grainy-gradients/) — SVG feTurbulence technique; HIGH confidence
- [Creating Grainy Backgrounds with CSS — ibelick.com](https://ibelick.com/blog/create-grainy-backgrounds-with-css) — data URI implementation with exact CSS; HIGH confidence
- [Designing Beautiful Shadows in CSS — Josh W. Comeau](https://www.joshwcomeau.com/css/designing-shadows/) — layered shadows, hue matching, elevation system; HIGH confidence
- [Staggered Animations with CSS Custom Properties — Cloud Four](https://cloudfour.com/thinks/staggered-animations-with-css-custom-properties/) — `--index` pattern, 0.025s increment reference; HIGH confidence
- [prefers-reduced-motion — MDN](https://developer.mozilla.org/en-US/docs/Web/CSS/Reference/At-rules/@media/prefers-reduced-motion) — accessibility requirement; HIGH confidence
- [Design Accessible Animation — Pope Tech (2025)](https://blog.pope.tech/2025/12/08/design-accessible-animation-and-movement/) — WCAG 2.1 SC 2.3.3 guidance; HIGH confidence
- [Badges vs Pills vs Chips vs Tags — Smart Interface Design Patterns](https://smart-interface-design-patterns.com/articles/badges-chips-tags-pills/) — pill sizing and legibility best practices; MEDIUM confidence
- [10 Modern Footer UX Patterns for 2026 — Eleken](https://www.eleken.co/blog-posts/footer-ux) — footer layout patterns; MEDIUM confidence
- [Dark Mode in CSS — CSS-Tricks](https://css-tricks.com/a-complete-guide-to-dark-mode-on-the-web/) — transition timing with theme switch; HIGH confidence
- Existing codebase: `web/static/main.css`, `web/templates/base.html` — reviewed directly; HIGH confidence

---

*Feature research for: v1.2 Shore Leave Polish — The Wild Meridian design system*
*Researched: 2026-03-28*
