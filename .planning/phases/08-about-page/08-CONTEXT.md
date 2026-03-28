# Phase 8: About Page - Context

**Gathered:** 2026-03-28
**Status:** Ready for planning

<domain>
## Phase Boundary

Readers can navigate to and read a personal about page that matches the site's nautical design. The page renders content from a static markdown file embedded in the binary, accessible via a new nav link. No new capabilities beyond what ABOUT-01, ABOUT-02, and ABOUT-03 require.

</domain>

<decisions>
## Implementation Decisions

### Navigation Placement
- **D-01:** "About" link appears in the nav bar **between the tagline and the dark mode toggle**. This is the first navigation link on the site.
- **D-02:** Plain text link styled to match the nav palette. No icon — clean and unobtrusive, letting the site name stay dominant.

### Content & Tone
- **D-03:** About page covers **both personal bio and blog purpose** — who you are and why The Wild Meridian exists.
- **D-04:** Nautical-flavored voice matching the site's existing "dispatches from the deep end" aesthetic. Not pirate cosplay, but thematically consistent.
- **D-05:** Claude drafts the initial `about.md` content. Software engineer basics — keep it generic and light.

### Markdown File Location
- **D-06:** `about.md` lives in the repo (e.g., `content/about.md`) and is compiled into the binary via `go:embed`. Zero runtime file dependencies. Edits require rebuild + redeploy, which is acceptable for a rarely-changing about page.

### Claude's Discretion
- Exact placement of `content/` directory (could be `web/content/`, `content/`, etc. — follow existing embed patterns)
- Template structure for the about page (follows per-page template set pattern from existing pages)
- CSS styling details for the About nav link
- About page heading and section structure within the markdown

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Templates (inherit from for about page)
- `web/templates/base.html` — Nav bar structure (line 26-44), footer, dark mode toggle, `{{block}}` pattern
- `web/templates/post.html` — Example of a content page template defining blocks from base
- `web/templates/list.html` — Example of per-page template set pattern

### Markdown Rendering
- `internal/markdown/renderer.go` — goldmark+bluemonday pipeline with `Render()` and `RenderWithMeta()` methods

### Route Registration
- `cmd/server/main.go` — Lines 116-133: blogMux route pattern (use `HandleFunc("GET /about", ...)`)

### Existing Static Embeds
- `web/static/` — Existing `go:embed` pattern for static assets (CSS, JS)
- `cmd/server/main.go` — Check how `staticFS` is embedded for the embed pattern to follow

### Requirements
- `.planning/REQUIREMENTS.md` — ABOUT-01, ABOUT-02, ABOUT-03

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/markdown/renderer.go`: `Render()` method converts markdown string to `template.HTML` — can render the embedded about.md directly
- `web/templates/base.html`: Full nautical chrome (nav, footer, dark mode) — about template just needs to define `content` and `title` blocks
- Dark mode toggle SVG in nav: establishes the right-side nav element pattern that About link will sit before

### Established Patterns
- Per-page template sets: each page parses `base.html` + its own template file independently (avoids block name collisions)
- `go:embed` for static assets: `web/static/` is already embedded as `staticFS` in main.go
- Route registration: `blogMux.HandleFunc("GET /path", handler)` with catch-all 404 at the end
- Template data includes `.Year` for dynamic copyright

### Integration Points
- `base.html` nav: Add "About" link element between tagline span and dark toggle button
- `cmd/server/main.go`: Register `GET /about` route on blogMux (before the catch-all 404)
- Blog handler package (`internal/handler/blog/`): Add `AboutPage` handler method
- Sitemap handler: Consider adding `/about` to the sitemap

</code_context>

<specifics>
## Specific Ideas

No specific requirements — open to standard approaches. The about page should feel like a natural part of the existing site, not a bolted-on afterthought.

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope

</deferred>

---

*Phase: 08-about-page*
*Context gathered: 2026-03-28*
