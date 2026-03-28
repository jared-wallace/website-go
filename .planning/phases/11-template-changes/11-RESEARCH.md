# Phase 11: Template Changes - Research

**Researched:** 2026-03-28
**Domain:** Go html/template, inline SVG, CSS layout, WAI-ARIA landmark patterns
**Confidence:** HIGH

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

- **D-01:** Footer uses a side-by-side two-section layout: nav links on the left, copyright/utility on the right. Stacks vertically on mobile.
- **D-02:** Footer nav section contains About link and RSS icon only. No Home or Posts links.
- **D-03:** Personality phrase tone is wry and weathered — "Still anchored. Still writing."
- **D-04:** SVG rope style is twisted two-strand (classic nautical dock line).
- **D-05:** Rope color uses existing `--color-divider` CSS variable for both light and dark modes.
- **D-06:** SVG rope replaces only the footer `.rope-divider` instance. CSS class updated so future uses inherit the pattern, but scope is footer-only for this phase.
- **D-07:** Hero heading text is "The Wild Meridian" (h1) with "dispatches from the deep end" as subtitle.
- **D-08:** Hero is modest in size — Playfair Display h1, Lora italic tagline, standard spacing. Posts remain the visual focus.
- **D-09:** Top nav becomes: site name + tagline + dark toggle. About link removed entirely from nav bar.
- **D-10:** Nav gets `aria-label="Main navigation"`, footer nav gets `aria-label="Footer navigation"`.

### Claude's Discretion

- Exact CSS spacing and breakpoint for footer column stacking (breakpoint confirmed as `max-width: 767px` matching existing card-grid pattern)

### Deferred Ideas (OUT OF SCOPE)

None — discussion stayed within phase scope.
</user_constraints>

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| NAV-01 | About link is removed from the navigation bar and appears in the footer | base.html line 29 removal + footer nav insertion; no Go handler changes needed |
| NAV-02 | Footer displays a two-section layout with navigation links and copyright/utility area | `.footer-inner` flex layout; new CSS block in main.css |
| NAV-03 | Footer includes a short nautical personality phrase matching the site's voice | Static copy in base.html footer; `.footer-nav` section |
| NAV-04 | Footer navigation has proper `aria-label` and primary nav gets a matching label for landmark disambiguation | Two `<nav>` elements; each MUST have distinct `aria-label` per WCAG 2.1 |
| ATMO-03 | Rope dividers render as an inline SVG twisted rope pattern replacing the current dashed border | `<hr>` → `<svg class="rope-divider">` swap; `.rope-divider` CSS rule updated |
| TYPO-03 | List page displays a hero heading/tagline area above the post card grid | New `.list-hero` block inserted before `.card-grid` in list.html |
</phase_requirements>

---

## Summary

Phase 11 is a focused HTML template + CSS phase with zero Go backend changes. All six requirements are addressed through edits to two template files (`base.html`, `list.html`) and one CSS file (`main.css`). The work is surgical: two template files need localized edits, and main.css receives one new labeled block appended at the end.

The most technically interesting piece is the inline SVG twisted rope pattern. A twisted two-strand rope is achievable with a single SVG `<path>` using sinusoidal curves — two interleaving sine-wave strands that cross over each other at regular intervals. The SVG must be `aria-hidden="true"` (decorative), use `stroke="var(--color-divider)"` (CSS variable resolves correctly inside SVG when `color` or explicit `stroke` attributes reference the variable), and be `width: 100%` to fill the footer width.

The WAI-ARIA change is straightforward: the existing `<nav class="site-nav">` already works as a landmark; adding `aria-label="Main navigation"` is a one-attribute change. The new footer `<nav aria-label="Footer navigation">` is the second landmark, required by WCAG 2.1 SC 1.3.6 and 2.4.1 to avoid duplicate unlabeled landmarks.

**Primary recommendation:** Execute in three tasks — (1) nav restructure + footer rebuild in base.html + companion CSS, (2) rope SVG replacement, (3) list page hero. Each task is independently testable via `go test ./...`.

---

## Standard Stack

### Core

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| `html/template` stdlib | Go 1.26 | Template rendering | Project standard; XSS-safe by default |
| `main.css` | — | All styles | Project rule: no new CSS files for v1.2 |
| Inline SVG | — | Rope divider graphic | No JS, no external assets, CSS variable compatible |

### Supporting

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| Playfair Display | (Google Fonts, already loaded) | Hero h1 | Display headings; already in `<head>` |
| Lora | (Google Fonts, already loaded) | Hero tagline, footer phrase | Body text; already in `<head>` |

### Alternatives Considered

| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| Inline SVG rope | CSS background SVG data URI | CSS variables don't resolve in data URIs; inline is the only reliable approach for `var(--color-divider)` |
| Inline SVG rope | CSS `repeating-linear-gradient` | Cannot replicate twisted rope appearance |
| New `.footer-nav-link` class | Reuse `.nav-link` | `.nav-link` already has correct 14px/muted/accent-hover/44px-touch-target styles — no duplication needed |

---

## Architecture Patterns

### File Change Inventory

```
web/templates/base.html    — 2 edits (nav + footer)
web/templates/list.html    — 1 edit (hero block insertion)
web/static/main.css        — 2 edits (rope-divider rule update + new Phase 11 block)
```

No Go source files change. No new files created.

### Pattern 1: Footer Two-Section Layout

**What:** Replace the current single-column centered footer with a flex row split into `.footer-nav` (left) and `.footer-copyright` (right). Wraps at 767px.

**When to use:** Any time a footer needs two regions with responsive stacking.

**Target HTML structure in base.html:**
```html
<footer class="site-footer">
  <svg class="rope-divider" ...></svg>
  <div class="footer-inner">
    <nav aria-label="Footer navigation">
      <a href="/about" class="nav-link">About</a>
      <a href="/rss" class="rss-link" aria-label="RSS feed"><!-- svg --></a>
    </nav>
    <div class="footer-copyright">
      <p class="footer-phrase">Still anchored. Still writing.</p>
      <p>&copy; {{.Year}} Jared Wallace</p>
    </div>
  </div>
</footer>
```

**CSS additions (new Phase 11 block at end of main.css):**
```css
/* === Phase 11: Template Changes === */

.footer-inner {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 24px;
}

@media (max-width: 767px) {
  .footer-inner {
    flex-direction: column;
    align-items: center;
    gap: 16px;
  }
}

.footer-copyright {
  text-align: right;
}

@media (max-width: 767px) {
  .footer-copyright {
    text-align: center;
  }
}

.footer-phrase {
  margin: 0 0 4px;
  font-style: italic;
  font-size: 14px;
  color: var(--color-text-muted);
}

.list-hero {
  margin-bottom: 32px;
}

.list-hero-title {
  font-family: var(--font-display);
  font-size: 28px;
  font-weight: 700;
  line-height: 1.2;
  margin: 0 0 8px;
}

.list-hero-tagline {
  font-family: var(--font-body);
  font-style: italic;
  font-size: 18px;
  font-weight: 400;
  line-height: 1.5;
  color: var(--color-text-muted);
  margin: 0;
}
```

### Pattern 2: Rope SVG Divider

**What:** A twisted two-strand rope path rendered in SVG, replacing `<hr class="rope-divider">`.

**Key constraint:** `stroke` must reference `var(--color-divider)` not a hex value, so dark mode works without a separate selector.

**SVG approach:** Two interleaving sinusoidal `<path>` elements. Each strand is a cubic bezier sequence producing a wave. The strands are offset by half a period to create the crossing/twist visual. A simplified but effective pattern:

```svg
<svg class="rope-divider" xmlns="http://www.w3.org/2000/svg"
     width="100%" height="12" aria-hidden="true" preserveAspectRatio="none">
  <!-- Strand 1: peaks up, then down -->
  <path d="M0,6 C10,1 20,11 30,6 C40,1 50,11 60,6 C70,1 80,11 90,6
           C100,1 110,11 120,6 C130,1 140,11 150,6 C160,1 170,11 180,6
           C190,1 200,11 210,6 C220,1 230,11 240,6 C250,1 260,11 270,6
           C280,1 290,11 300,6"
        fill="none" stroke="var(--color-divider)" stroke-width="2"
        stroke-linecap="round"/>
  <!-- Strand 2: offset half period — peaks down, then up -->
  <path d="M0,6 C10,11 20,1 30,6 C40,11 50,1 60,6 C70,11 80,1 90,6
           C100,11 110,1 120,6 C130,11 140,1 150,6 C160,11 170,1 180,6
           C190,11 200,1 210,6 C220,11 230,1 240,6 C250,11 260,1 270,6
           C280,11 290,1 300,6"
        fill="none" stroke="var(--color-divider)" stroke-width="2"
        stroke-linecap="round" opacity="0.7"/>
</svg>
```

**CSS rule update** (replace existing `.rope-divider` block at line 445):
```css
/* --- Rope Divider --- */

.rope-divider {
  display: block;
  width: 100%;
  height: 12px;
  margin: 24px 0;
  overflow: visible;
}
```

The old `border-top: 2px dashed` and `opacity: 0.6` are removed. The SVG itself carries the visual.

**Note on `.site-footer .rope-divider`:** The existing override at line 585 (`margin-bottom: 16px`) will need review — with the new two-section layout, the rope sits above `.footer-inner`, so `margin: 24px 0` from the base rule is appropriate. The `.site-footer .rope-divider` override can be removed.

### Pattern 3: List Page Hero

**What:** Unconditional hero heading block inserted before the conditional `.card-grid` / `.empty-state` block.

**Template change in list.html:**
```html
{{define "content"}}
<div class="container">
  <div class="list-hero">
    <h1 class="list-hero-title">The Wild Meridian</h1>
    <p class="list-hero-tagline">dispatches from the deep end</p>
  </div>
  {{if not .Posts}}
  ...existing empty state...
```

**Accessibility note:** This `<h1>` must be the only `<h1>` on the list page. Post card titles use `<h2>` — confirmed in existing `list.html` line 34 (`<h2 class="post-card-title">`). No heading hierarchy violation.

### Pattern 4: Nav aria-label

**What:** One-attribute addition to the existing `<nav>` element.

**Before (base.html line 26):**
```html
<nav class="site-nav">
```

**After:**
```html
<nav class="site-nav" aria-label="Main navigation">
```

No CSS change. No handler change. The dark mode transition list in main.css references `.site-nav` by class — `aria-label` addition does not affect selector specificity.

### Anti-Patterns to Avoid

- **`stroke="currentColor"` for rope SVG:** `currentColor` inherits the element's `color` property. The footer text color is `var(--color-text-muted)`, which could make the rope too faint or too dark depending on mode. Use `stroke="var(--color-divider)"` directly.
- **Hiding About link with `display: none`:** D-09 specifies full removal from nav, not hiding. A hidden element still renders in accessibility trees and could confuse screen readers.
- **Inline `style` attributes:** Project rule: no inline styles; all colors via CSS custom properties.
- **Creating a separate `.footer-nav-link` class:** The `.nav-link` class already provides all required styles (14px, muted color, accent hover, 44px touch target). Duplication wastes lines.
- **Placing hero inside the `{{if .Posts}}` block:** The hero is unconditional (D-08). It renders even on the empty state page.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Two-landmark nav disambiguation | Custom JS focus management | `aria-label` attributes on `<nav>` elements | WCAG 2.1 native landmark labeling; browsers + screen readers handle it |
| Dark-mode-aware SVG color | Duplicate SVG + CSS dark selector | `stroke="var(--color-divider)"` | CSS custom properties resolve inside inline SVG; no duplication needed |
| Responsive footer stacking | JS resize observer | `@media (max-width: 767px)` CSS | CSS-only; matches existing breakpoint in codebase |

**Key insight:** This phase contains no new technical primitives — it composes existing patterns (flex layout, media queries, CSS variables, `aria-label`) in new locations.

---

## Common Pitfalls

### Pitfall 1: CSS Variable in SVG `stroke` Attribute

**What goes wrong:** Developer uses a hex color like `stroke="#C8B89A"` thinking it's simpler. Dark mode breaks — the rope stays light-mode color.

**Why it happens:** SVGs inside HTML inherit the CSS cascade including custom properties. `stroke="var(--color-divider)"` works on inline SVGs. It does NOT work on external SVGs loaded via `<img>` or `<object>`, but inline SVG is the approach here (D-06).

**How to avoid:** Use `stroke="var(--color-divider)"` in the SVG path attributes. Verified: CSS custom properties in SVG presentation attributes work in all modern browsers when the SVG is inline in HTML.

**Confidence:** HIGH — this is standard CSS behavior for inline SVGs.

### Pitfall 2: Duplicate Unlabeled `<nav>` Landmarks

**What goes wrong:** Developer adds the footer `<nav>` but forgets `aria-label` on either the primary nav or the footer nav. Screen reader users hear "navigation" twice with no way to distinguish which is which.

**Why it happens:** It's easy to add the footer `<nav>` element without remembering WCAG 2.1 requires each landmark instance to be uniquely labeled when there are multiple of the same type.

**How to avoid:** Both `<nav>` elements get `aria-label` simultaneously in the same commit. NAV-04 acceptance criteria explicitly calls this out.

**Warning signs:** Running axe DevTools or NVDA on the page and hearing unlabeled "navigation" twice.

### Pitfall 3: `<h1>` Conflict on List Page

**What goes wrong:** Adding `<h1>` to the hero block while forgetting that post card titles or another template element also uses `<h1>`.

**Why it happens:** Template inheritance means `base.html` doesn't declare an `<h1>` — but other templates might.

**How to avoid:** Confirmed: existing `list.html` uses `<h2>` for post card titles (line 34). The about page has its own `<h1>` but lives in a separate template. The new hero `<h1>` on list.html will be the only `<h1>` on that page.

### Pitfall 4: `preserveAspectRatio` on Width-100% SVG

**What goes wrong:** SVG rope renders squished or distorted at narrow viewports because the viewBox aspect ratio is enforced.

**Why it happens:** The default `preserveAspectRatio="xMidYMid meet"` tries to preserve the path proportions.

**How to avoid:** Set `preserveAspectRatio="none"` on the SVG element so the path stretches horizontally to fill `width: 100%`. The rope wave will tile naturally at any width.

### Pitfall 5: `.site-footer p` Rule Interfering with New Markup

**What goes wrong:** The existing rule `.site-footer p { display: flex; justify-content: center; }` (main.css line 589) will apply to new `<p>` elements inside the footer — including `.footer-phrase` and the copyright `<p>`.

**Why it happens:** The rule was written for the single-line footer layout where the only `<p>` held the copyright + RSS icon flex row.

**How to avoid:** Replace the generic `.site-footer p` flex rule with a specific `.footer-copyright-line` class (or scope it appropriately) so only the copyright row gets flex treatment. Alternatively, remove `.site-footer p` from main.css and add the flex only where needed. The Phase 11 CSS block can override it.

---

## Code Examples

### NAV-01 / NAV-04: Nav Restructure (base.html)

```html
<!-- Source: base.html lines 26-46, modified per D-09, D-10 -->
<nav class="site-nav" aria-label="Main navigation">
  <a href="/" class="site-name">The Wild Meridian</a>
  <span class="site-tagline">dispatches from the deep end</span>
  <button class="dark-toggle" id="dark-toggle" aria-label="Switch to dark mode">
    <!-- existing sun/moon SVGs unchanged -->
  </button>
</nav>
```

### NAV-02 / NAV-03 / NAV-04: Footer Rebuild (base.html)

```html
<!-- Source: base.html lines 50-64, rebuilt per D-01, D-02, D-03, D-10 -->
<footer class="site-footer">
  <svg class="rope-divider" xmlns="http://www.w3.org/2000/svg"
       width="100%" height="12" aria-hidden="true" preserveAspectRatio="none">
    <path d="M0,6 C10,1 20,11 30,6 ..." fill="none" stroke="var(--color-divider)" stroke-width="2" stroke-linecap="round"/>
    <path d="M0,6 C10,11 20,1 30,6 ..." fill="none" stroke="var(--color-divider)" stroke-width="2" stroke-linecap="round" opacity="0.7"/>
  </svg>
  <div class="footer-inner">
    <nav aria-label="Footer navigation">
      <a href="/about" class="nav-link">About</a>
      <a href="/rss" class="rss-link" aria-label="RSS feed">
        <!-- existing RSS SVG unchanged -->
      </a>
    </nav>
    <div class="footer-copyright">
      <p class="footer-phrase">Still anchored. Still writing.</p>
      <p class="footer-copyright-line">&copy; {{.Year}} Jared Wallace</p>
    </div>
  </div>
</footer>
```

### TYPO-03: List Page Hero (list.html)

```html
<!-- Source: list.html — insert before existing {{if not .Posts}} block -->
{{define "content"}}
<div class="container">
  <div class="list-hero">
    <h1 class="list-hero-title">The Wild Meridian</h1>
    <p class="list-hero-tagline">dispatches from the deep end</p>
  </div>
  {{if not .Posts}}
  ...
```

### ATMO-03: Rope Divider CSS Update (main.css)

```css
/* --- Rope Divider --- */

/* Replace dashed border with inline-SVG-compatible block rule.
   The <hr> is replaced by <svg class="rope-divider"> in base.html. */
.rope-divider {
  display: block;
  width: 100%;
  height: 12px;
  margin: 24px 0;
  overflow: visible;
}
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| `<hr>` with `border-top: dashed` | Inline `<svg>` with sinusoidal path | Phase 11 | Visual upgrade; CSS variable color compatibility maintained |
| Single-region centered footer | Two-section flex footer with `<nav>` | Phase 11 | Structural improvement; About link moved here |
| Unlabeled `<nav>` | `aria-label="Main navigation"` | Phase 11 | WCAG 2.1 compliance; screen reader landmark clarity |

**Deprecated/outdated after this phase:**
- `.site-footer p { display: flex; justify-content: center; }` — too broad; replaced by scoped `.footer-copyright-line` class

---

## Open Questions

1. **SVG path length vs. viewport**
   - What we know: the SVG uses `width="100%"` with a fixed 300-unit viewBox path and `preserveAspectRatio="none"`, so the path stretches to fill any width.
   - What's unclear: at very narrow viewports (<320px) the wave period will look compressed. Acceptable for this phase given the "modest decorative" intent.
   - Recommendation: Accept minor compression. If the rope needs to look natural at all widths, use a repeating pattern via `<pattern>` + `<rect>`, but that adds complexity beyond the D-06 scope.

2. **`.site-footer p` selector collision**
   - What we know: existing rule at line 589 applies flex to all `<p>` in footer; new layout adds multiple `<p>` elements.
   - What's unclear: whether removing the rule would regress any other page.
   - Recommendation: Add a targeted `.footer-copyright-line` class for the copyright row and remove the generic `.site-footer p` flex rule in the Phase 11 CSS block.

---

## Environment Availability

Step 2.6: SKIPPED (no external dependencies — pure template and CSS changes, no new tools required).

---

## Validation Architecture

### Test Framework

| Property | Value |
|----------|-------|
| Framework | Go `testing` package (stdlib) |
| Config file | none — `go test` discovers `*_test.go` files automatically |
| Quick run command | `go test ./internal/handler/blog/... -run TestList` |
| Full suite command | `go test ./...` |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| NAV-01 | About link absent from nav HTML | unit | `go test ./internal/handler/blog/... -run TestNavAboutLinkRemoved` | ❌ Wave 0 |
| NAV-02 | Footer contains `.footer-inner` with two child sections | unit | `go test ./internal/handler/blog/... -run TestFooterTwoSection` | ❌ Wave 0 |
| NAV-03 | Footer contains personality phrase text | unit | `go test ./internal/handler/blog/... -run TestFooterPersonalityPhrase` | ❌ Wave 0 |
| NAV-04 | Primary nav has `aria-label="Main navigation"`; footer nav has `aria-label="Footer navigation"` | unit | `go test ./internal/handler/blog/... -run TestNavAriaLabels` | ❌ Wave 0 |
| ATMO-03 | `<svg class="rope-divider">` present; no `<hr class="rope-divider">` | unit | `go test ./internal/handler/blog/... -run TestRopeDividerSVG` | ❌ Wave 0 |
| TYPO-03 | List page response contains `.list-hero` with h1 and tagline | unit | `go test ./internal/handler/blog/... -run TestListHero` | ❌ Wave 0 |

**Note:** The existing `TestListPostsEmpty` and `TestListOGMeta` tests in `handler_test.go` parse rendered HTML via `strings.Contains` — the new tests follow the same pattern. No new test infrastructure needed; just new test functions in `handler_test.go`.

### Sampling Rate

- **Per task commit:** `go test ./internal/handler/blog/... -run TestList`
- **Per wave merge:** `go test ./...`
- **Phase gate:** `go test ./...` green before `/gsd:verify-work`

### Wave 0 Gaps

- [ ] `internal/handler/blog/handler_test.go` — add `TestNavAboutLinkRemoved`, `TestFooterTwoSection`, `TestFooterPersonalityPhrase`, `TestNavAriaLabels`, `TestRopeDividerSVG`, `TestListHero` functions. All use the existing `newTestHandler` + `strings.Contains` pattern already established in the file.

---

## Project Constraints (from CLAUDE.md)

| Directive | Applies to This Phase |
|-----------|----------------------|
| Go stdlib preferred, minimal dependencies | Confirmed — zero new dependencies in this phase |
| No large frameworks; stdlib where reasonable | Confirmed — html/template throughout |
| `net/http` ServeMux routing | No new routes in this phase |
| All persistent data on EBS at `/var/www/html` | Not applicable — no data changes |
| Leverage `frontend-design` skill for template/UI work | No `.claude/skills/` directory found in this project — applying design system conventions from UI-SPEC.md directly |
| No new CSS files for v1.2 | Confirmed — all additions appended to `main.css` with labeled block comment |
| No inline styles — all colors via CSS custom properties | Confirmed — SVG stroke uses `var(--color-divider)` |
| Dark mode via `[data-theme="dark"]` selector with `.theme-ready` gating | No new animated elements; existing gating in place from Phase 10 |
| GSD workflow enforcement — no direct repo edits outside a GSD workflow | Research only; implementation will follow `/gsd:execute-phase` |

---

## Sources

### Primary (HIGH confidence)

- Official WCAG 2.1 — SC 1.3.6 and 2.4.1 landmark labeling requirement for multiple same-type landmarks
- MDN Web Docs — CSS custom properties (`var()`) in inline SVG `stroke` attributes: properties resolve in SVG presentation attributes when SVG is inline in HTML
- `web/templates/base.html` — current nav/footer markup (read directly)
- `web/templates/list.html` — current list page template (read directly)
- `web/static/main.css` — all existing CSS rules including `.rope-divider` (line 445), `.site-footer` (line 578), `.site-nav` (line 55) (read directly)
- `.planning/phases/11-template-changes/11-CONTEXT.md` — locked decisions D-01 through D-10
- `.planning/phases/11-template-changes/11-UI-SPEC.md` — approved UI contract with precise typography, spacing, and color values

### Secondary (MEDIUM confidence)

- `internal/handler/blog/handler_test.go` — existing test pattern using `strings.Contains` for HTML assertion; new tests follow this pattern

### Tertiary (LOW confidence)

- None — all findings verified from project source files and official specifications.

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — project files read directly; no uncertainty
- Architecture: HIGH — all patterns drawn from existing codebase; UI-SPEC approved
- Pitfalls: HIGH — CSS variable in SVG and `.site-footer p` collision are directly observable from main.css source

**Research date:** 2026-03-28
**Valid until:** 2026-04-28 (stable template/CSS domain; no external API dependencies)
