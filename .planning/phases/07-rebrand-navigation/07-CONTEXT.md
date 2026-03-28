# Phase 7: Rebrand + Navigation - Context

**Gathered:** 2026-03-28
**Status:** Ready for planning

<domain>
## Phase Boundary

Every public surface reflects "The Wild Meridian" identity instead of "The Log", and readers can discover the RSS feed via a visible icon. This is a systematic rebrand (string replacements across templates, RSS handler, OG/Twitter meta) plus adding an RSS icon to the footer.

</domain>

<decisions>
## Implementation Decisions

### RSS Discoverability
- **D-01:** RSS icon appears in the **footer only**, next to the copyright line. Nav bar stays clean.
- **D-02:** Standard RSS broadcast SVG icon (the classic radiating-lines icon), color-matched to the site's nautical palette. Not a text link.

### Admin Panel Branding
- **D-03:** Admin templates also rebrand to "The Wild Meridian". Nav becomes "The Wild Meridian -- Back Office", login heading updates. Full consistency, no split personality.

### Copyright Footer
- **D-04:** Copyright uses dynamic `{{.Year}}` (not hardcoded). Text changes from `(c) {{.Year}} The Log` to `(c) {{.Year}} Jared Wallace`. Blog name drops out of the copyright line.

### Claude's Discretion
- RSS icon sizing and exact color values within the existing palette
- SVG icon implementation details (inline vs sprite)
- Any minor template formatting adjustments needed during the rename

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Templates (all contain "The Log" strings to replace)
- `web/templates/base.html` -- Title, OG/Twitter meta, nav `.site-name`, copyright footer, RSS alternate link title
- `web/templates/list.html` -- Title and OG/Twitter meta overrides
- `web/templates/post.html` -- Title block
- `web/templates/404.html` -- Title block
- `web/templates/admin-base.html` -- Title and nav (rebrand to "The Wild Meridian -- Back Office")
- `web/templates/admin-login.html` -- Heading

### RSS Handler
- `internal/handler/blog/rss.go` -- RSS channel title on line 91, description on line 93

### Tests (assertions reference "The Log")
- `internal/handler/blog/rss_test.go` -- Line 45: `<title>The Log</title>`
- `internal/handler/blog/handler_test.go` -- Line 175: `og:title" content="The Log"`

### Requirements
- `.planning/REQUIREMENTS.md` -- BRAND-01 through BRAND-04, NAV-01, NAV-02

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `base.html` already has `<link rel="alternate" type="application/rss+xml">` in the head (NAV-02 partially satisfied -- just needs title update)
- Dark mode toggle SVG pattern in nav can be followed for RSS icon SVG implementation
- `rope-divider` CSS class on footer `<hr>` establishes the footer visual pattern

### Established Patterns
- Templates use Go `html/template` with `{{define}}` / `{{block}}` inheritance
- OG/Twitter meta tags are overridden per-page via `{{block "meta" .}}`
- `.Year` is passed as template data for dynamic copyright year
- RSS feed is built via Go structs with XML annotations in `rss.go`

### Integration Points
- Footer in `base.html` is the single point for the RSS icon addition (all pages inherit it)
- RSS handler `buildRSSFeed()` function is the single point for feed metadata changes
- Test files will need assertion string updates to match new branding

</code_context>

<specifics>
## Specific Ideas

No specific requirements -- open to standard approaches for the rebrand. The RSS icon should follow the same inline SVG pattern used by the dark mode toggle.

</specifics>

<deferred>
## Deferred Ideas

None -- discussion stayed within phase scope

</deferred>

---

*Phase: 07-rebrand-navigation*
*Context gathered: 2026-03-28*
