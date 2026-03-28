# Phase 11: Template Changes - Context

**Gathered:** 2026-03-28
**Status:** Ready for planning

<domain>
## Phase Boundary

Restructure the navigation and footer layout: move the About link from the top nav to a new two-section footer, add a nautical personality phrase, apply aria-labels for screen reader landmark disambiguation, replace the dashed CSS rope divider with an inline SVG twisted rope pattern, and add a hero heading/tagline above the post card grid on the list page.

</domain>

<decisions>
## Implementation Decisions

### Footer Layout (NAV-02)
- **D-01:** Footer uses a side-by-side two-section layout: nav links on the left, copyright/utility on the right. Stacks vertically on mobile.
- **D-02:** Footer nav section contains About link and RSS icon only. No Home or Posts links.

### Footer Voice (NAV-03)
- **D-03:** Personality phrase tone is wry and weathered — think "Built on salt air and stubbornness" or "Still anchored. Still writing."

### Claude's Discretion
- Exact wording of the nautical personality phrase (must be wry/weathered tone, short)
- Exact CSS spacing and breakpoint for footer column stacking

### Rope Divider (ATMO-03)
- **D-04:** SVG rope style is twisted two-strand (classic nautical dock line).
- **D-05:** Rope color uses existing `--color-divider` CSS variable for both light and dark modes.
- **D-06:** SVG rope replaces only the footer `.rope-divider` instance. The CSS class is updated so future uses inherit the pattern, but scope is footer-only for this phase.

### List Page Hero (TYPO-03)
- **D-07:** Hero heading text is "The Wild Meridian" (h1) with "dispatches from the deep end" as subtitle.
- **D-08:** Hero is modest in size — Playfair Display h1, Lora italic tagline, standard spacing. Posts remain the visual focus.

### Nav Restructuring (NAV-01, NAV-04)
- **D-09:** Top nav becomes: site name + tagline + dark toggle. About link removed entirely from nav bar.
- **D-10:** Nav gets `aria-label="Main navigation"`, footer nav gets `aria-label="Footer navigation"`.

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Templates
- `web/templates/base.html` — Contains both the `<nav>` and `<footer>` that need restructuring
- `web/templates/list.html` — List page template where hero heading is added

### Styles
- `web/static/main.css` — All CSS including `.site-nav`, `.site-footer`, `.rope-divider`, `.nav-link` rules

### Requirements
- `.planning/REQUIREMENTS.md` — NAV-01 through NAV-04, ATMO-03, TYPO-03 acceptance criteria

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `.rope-divider` CSS class already exists (line 445 of main.css) — needs SVG replacement not new class
- `.nav-link` styles exist (line 93 of main.css) — can be reused for footer links
- Playfair Display and Lora fonts already loaded via Google Fonts in base.html
- `.site-tagline` class exists in nav — can inform hero tagline styling

### Established Patterns
- Go `html/template` with `{{block}}` / `{{define}}` inheritance from base.html
- CSS variables for all colors (--color-divider, --color-text-muted, --color-accent, etc.)
- Dark mode via `[data-theme="dark"]` selector with `.theme-ready` gating for transitions
- No CSS files beyond main.css — all changes go in main.css with labeled block comments

### Integration Points
- `base.html` line 26-46: `<nav class="site-nav">` — remove About link, add aria-label
- `base.html` line 50-64: `<footer class="site-footer">` — restructure to two-section layout with nav element
- `list.html` line 18: Start of `{{define "content"}}` — insert hero heading before card grid
- `.Year` template variable already available in footer context

</code_context>

<specifics>
## Specific Ideas

No specific requirements — open to standard approaches within the decisions above.

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.

</deferred>

---

*Phase: 11-template-changes*
*Context gathered: 2026-03-28*
