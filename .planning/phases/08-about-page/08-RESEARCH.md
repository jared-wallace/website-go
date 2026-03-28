# Phase 8: About Page - Research

**Researched:** 2026-03-28
**Domain:** Go html/template about page — static embedded markdown, nav link, existing chrome reuse
**Confidence:** HIGH

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions
- **D-01:** "About" link appears in the nav bar between the tagline and the dark mode toggle. This is the first navigation link on the site.
- **D-02:** Plain text link styled to match the nav palette. No icon — clean and unobtrusive, letting the site name stay dominant.
- **D-03:** About page covers both personal bio and blog purpose — who you are and why The Wild Meridian exists.
- **D-04:** Nautical-flavored voice matching the site's existing "dispatches from the deep end" aesthetic. Not pirate cosplay, but thematically consistent.
- **D-05:** Claude drafts the initial `about.md` content. Software engineer basics — keep it generic and light.
- **D-06:** `about.md` lives in the repo (e.g., `content/about.md`) and is compiled into the binary via `go:embed`. Zero runtime file dependencies. Edits require rebuild + redeploy, which is acceptable for a rarely-changing about page.

### Claude's Discretion
- Exact placement of `content/` directory (could be `web/content/`, `content/`, etc. — follow existing embed patterns)
- Template structure for the about page (follows per-page template set pattern from existing pages)
- CSS styling details for the About nav link
- About page heading and section structure within the markdown

### Deferred Ideas (OUT OF SCOPE)
None — discussion stayed within phase scope.
</user_constraints>

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| ABOUT-01 | User can navigate to an about page from the main site navigation | Nav link in `base.html`, `GET /about` route in `blogMux`, `BlogHandler.AboutPage` method |
| ABOUT-02 | About page renders content from a static markdown file on disk | `go:embed` pattern for `content/about.md`; `markdown.Renderer.Render()` already exists |
| ABOUT-03 | About page matches the existing nautical design (header, footer, dark mode) | `base.html` chrome is inherited by all templates — reuse `.container.container-narrow` and `.post-body` |
</phase_requirements>

---

## Summary

Phase 8 is a low-risk feature addition that composes existing, fully-functional building blocks. The blog handler already parses per-page template sets, the markdown renderer is ready to consume any string, and the `go:embed` pattern is established in `web/embed.go`. Nothing needs to be invented — this phase wires three existing systems together with one new route, one new handler method, one new template, one new CSS class, one new content file, and one nav link addition.

The critical layout constraint is that `.dark-toggle` uses `margin-left: auto` to push itself to the far right of the flex nav. The About nav link must be inserted in the HTML before the toggle — it will naturally sit left of the auto-margin. No reordering of flex items or new margin tricks are required; the insertion point is purely a DOM ordering concern.

The `content/` directory (sibling to `web/`) is the recommended embed location. It keeps content separate from template/static assets and follows the established pattern where `web/embed.go` handles embed declarations. A second `//go:embed content` variable can be added to the same file.

**Primary recommendation:** Add `content/about.md` with `go:embed`, expose it as a package-level `embed.FS`, call `markdown.Renderer.Render()` in the new `AboutPage` handler, and wire the route before the catch-all 404.

---

## Standard Stack

No new dependencies. Phase 8 uses only what is already imported.

### Core (all existing)
| Library | Purpose | Location |
|---------|---------|----------|
| `html/template` stdlib | Template rendering | `internal/handler/blog/handler.go` |
| `internal/markdown.Renderer` | Markdown to `template.HTML` | `internal/markdown/renderer.go` |
| `embed` stdlib | Binary-embedded content files | `web/embed.go` (to extend) |
| `net/http` stdlib | Route registration | `cmd/server/main.go` |

**Installation:** None required.

---

## Architecture Patterns

### Recommended Project Structure (additions only)

```
content/
└── about.md            # NEW — embedded markdown content for about page

web/
├── embed.go            # MODIFY — add Content embed.FS var
└── templates/
    └── about.html      # NEW — per-page template set (base + about)

internal/handler/blog/
└── about.go            # NEW — BlogHandler.AboutPage method
```

### Pattern 1: Per-Page Template Set (established)

**What:** Each HTML page parses `base.html` + its own template file into an independent `*template.Template`. This prevents `{{define}}` block name collisions across pages.

**When to use:** Every new public-facing page.

**Example (from handler.go lines 48-58):**
```go
// Source: internal/handler/blog/handler.go
pages := []string{"list.html", "post.html", "404.html"}
templates := make(map[string]*template.Template, len(pages))
for _, page := range pages {
    tmpl := template.Must(
        template.New("").Funcs(funcMap).ParseFS(
            web.Templates,
            "templates/base.html",
            "templates/"+page,
        ),
    )
    templates[page] = tmpl
}
```

**For this phase:** Add `"about.html"` to the `pages` slice — that's the only change to `handler.go` initialization.

### Pattern 2: go:embed FS Variable (established)

**What:** Declare a package-level `embed.FS` variable with a `//go:embed` directive. The FS is read-only and populated at compile time.

**When to use:** Any file that should be bundled into the binary.

**Example (from web/embed.go):**
```go
// Source: web/embed.go
//go:embed templates
var Templates embed.FS

//go:embed static
var Static embed.FS
```

**For this phase:** Add a `Content` variable to `web/embed.go`:
```go
//go:embed content
var Content embed.FS
```

The `content/` directory must be a sibling of `web/` at the repo root, OR placed inside `web/content/` if the embed declaration lives in `web/embed.go`. Since `web/embed.go` is the embed hub, `web/content/about.md` with `//go:embed content` is the cleanest path — it keeps all embed declarations in one file without adding a new package.

### Pattern 3: Handler Method with Embedded Content

**What:** Read the embedded file once at handler construction time (or at first request — either works since it's immutable), render via `markdown.Renderer.Render()`, pass rendered HTML to the template.

**When to use:** Static content pages.

**Example (new file — follows ShowPost structure):**
```go
// Source: internal/handler/blog/about.go (new)
func (h *BlogHandler) AboutPage(w http.ResponseWriter, r *http.Request) {
    h.render(w, http.StatusOK, "about.html", map[string]interface{}{
        "RenderedHTML": h.aboutHTML,
    })
}
```

`h.aboutHTML` is a `template.HTML` field populated at `New()` time by reading from the embedded FS and calling `h.renderer.Render(src)`.

### Pattern 4: Route Registration (established)

**What:** Register `GET /about` on `blogMux` before the catch-all.

**When to use:** Every new public route.

**Example (from main.go lines 118-133):**
```go
// Source: cmd/server/main.go
blogMux.HandleFunc("GET /about", blog.AboutPage)
// ... must appear before the catch-all:
blogMux.HandleFunc("GET /{path...}", blog.NotFound)
```

### Pattern 5: About Template Structure

**What:** Minimal template that defines `title` and `content` blocks, delegating all chrome to `base.html`.

**Example (new file — follows post.html structure):**
```html
{{template "base" .}}

{{define "title"}}About — The Wild Meridian{{end}}

{{define "content"}}
<div class="container container-narrow">
  <article class="about-page">
    <h1 class="about-title">About</h1>
    <hr class="rope-divider">
    <div class="post-body">{{.RenderedHTML}}</div>
  </article>
</div>
{{end}}
```

No `{{define "meta"}}` override needed — the base template's default OG meta is acceptable for a static about page.

### Anti-Patterns to Avoid

- **Reading from disk at runtime:** `about.md` must be embedded via `go:embed`, not `os.ReadFile`. Runtime disk reads break the containerized deploy model.
- **Parsing templates per request:** Templates are parsed once at `New()` and cached in `h.templates`. Do not parse in the handler method.
- **Using `RenderWithMeta` unnecessarily:** The about page needs no YAML front matter. Use `Render()` not `RenderWithMeta()` — simpler and avoids front matter leaking into rendered output.
- **Registering route after catch-all:** `GET /about` must appear before `GET /{path...}` in `blogMux` registration. Go 1.22 ServeMux uses most-specific match, but explicit registration before the catch-all is the project convention.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Markdown → HTML | Custom string replacement | `internal/markdown.Renderer.Render()` | Goldmark + bluemonday pipeline already handles GFM, XSS sanitization, syntax highlighting |
| Template inheritance | Custom include system | `html/template` `{{template "base" .}}` | Already working, already tested |
| File embedding | `os.ReadFile` at runtime | `go:embed` + `embed.FS` | Binary portability; no filesystem deps at runtime |
| Dark mode toggle | New JS | Inherited from `base.html` | Already in `main.js`, no per-page wiring needed |

---

## Common Pitfalls

### Pitfall 1: `margin-left: auto` on dark toggle displaces About link
**What goes wrong:** `.dark-toggle` has `margin-left: auto` which pushes it to the far right. If the About `<a>` is placed after the toggle in the HTML, it will appear to the right of or behind the toggle button.
**Why it happens:** Flexbox `margin-left: auto` absorbs all remaining space, pushing subsequent siblings to the right edge — or the About link ends up after the toggle visually.
**How to avoid:** Insert the About `<a>` in `base.html` between `.site-tagline` and the `.dark-toggle` button — before the button in DOM order. The `margin-left: auto` on `.dark-toggle` still pushes the toggle right, and the About link sits naturally between the tagline and that auto-margin gap.
**Warning signs:** About link appears to the right of the moon/sun icon in the nav bar.

### Pitfall 2: `about.md` front matter leaks into rendered HTML
**What goes wrong:** If `RenderWithMeta` is used and the markdown file has YAML front matter, goldmark renders the front matter block as literal text rather than stripping it.
**Why it happens:** `goldmark-meta` strips front matter from the parse tree but only when accessed via `meta.Get(ctx)`. The raw front matter still passes through `Convert()` as a code block if the context isn't used correctly.
**How to avoid:** Do not add YAML front matter to `about.md`. Use `Render()` not `RenderWithMeta()`. The about page has no metadata needs (no publish date, no tags).
**Warning signs:** `---` and YAML keys appear at the top of the rendered about page.

### Pitfall 3: Template map key mismatch causes 500
**What goes wrong:** Handler calls `h.render(w, 200, "about.html", data)` but `"about.html"` was not added to the `pages` slice in `New()`, so `h.templates["about.html"]` is nil. The `render()` method logs "unknown template page" and returns 500.
**Why it happens:** Two-place registration — `pages` slice in `New()` AND `HandleFunc` in `main.go`.
**How to avoid:** Add `"about.html"` to the `pages` slice in `handler.go New()` at the same time as adding the route. The handler panics at startup (via `template.Must`) if the template file is missing — which is the desired early failure mode.
**Warning signs:** `/about` returns 500 with "unknown template page" in logs.

### Pitfall 4: `go:embed` path must be relative to the file declaring it
**What goes wrong:** `//go:embed content` in `web/embed.go` looks for `web/content/`, not `content/` at the repo root.
**Why it happens:** `go:embed` paths are resolved relative to the package directory of the file containing the directive.
**How to avoid:** Place `content/about.md` at `web/content/about.md` if the embed directive lives in `web/embed.go`. OR create a new `content/embed.go` package at the repo root. The simpler path is `web/content/about.md` with an existing embed hub.
**Warning signs:** `go build` fails with "pattern content: directory prefix content does not exist".

### Pitfall 5: Sitemap omission
**What goes wrong:** `/about` is accessible but not listed in `sitemap.xml`. Search engines may not discover it.
**Why it happens:** `buildSitemap` in `sitemap.go` iterates only published posts — it has no awareness of static pages.
**How to avoid:** Add a `SitemapURL` for `/about` in `buildSitemap()`, prepended alongside the homepage entry. Priority `0.6`, `changefreq: "yearly"`.
**Warning signs:** `GET /sitemap.xml` response does not contain `/about`.

---

## Code Examples

### Adding `about.html` to the template set (handler.go)
```go
// Source: internal/handler/blog/handler.go — modify pages slice
pages := []string{"list.html", "post.html", "404.html", "about.html"}
```

### Embedding content directory (web/embed.go)
```go
// Source: web/embed.go — add Content var
//go:embed content
var Content embed.FS
```

### Reading embedded markdown at construction (about.go)
```go
// Source: internal/handler/blog/about.go (new file)
package blog

import (
    "log/slog"
    "net/http"
)

// AboutPage handles GET /about, rendering the embedded about.md as HTML.
func (h *BlogHandler) AboutPage(w http.ResponseWriter, r *http.Request) {
    h.render(w, http.StatusOK, "about.html", map[string]interface{}{
        "RenderedHTML": h.aboutHTML,
    })
}
```

`h.aboutHTML` is a `template.HTML` field on `BlogHandler`, populated in `New()`:
```go
// In New(), after funcMap setup:
raw, err := web.Content.ReadFile("content/about.md")
if err != nil {
    panic("about.md not found in embedded FS: " + err.Error())
}
aboutHTML := renderer.Render(string(raw))
```

This requires passing `renderer *markdown.Renderer` to `New()` — check if the current `New()` signature takes the service only. The post service holds a reference to the renderer internally; for the about page, pass the renderer directly to `New()` or expose it from the service.

**Alternative:** Pre-render `aboutHTML` inside the `postservice.Service` as a static field. Either approach is valid — the simpler is passing renderer to `blog.New()`.

### Nav link insertion (base.html)
```html
<!-- Between .site-tagline and .dark-toggle -->
<span class="site-tagline">dispatches from the deep end</span>
<a href="/about" class="nav-link">About</a>
<button class="dark-toggle" ...>
```

### New CSS class (main.css)
```css
/* About nav link — sits between tagline and dark toggle */
.nav-link {
  font-size: 14px;
  color: var(--color-text-muted);
  text-decoration: none;
  margin-left: 16px;
  min-height: 44px;
  display: inline-flex;
  align-items: center;
  white-space: nowrap;
}

.nav-link:hover {
  color: var(--color-accent);
}

.nav-link:focus-visible {
  outline: 2px solid var(--color-accent);
  outline-offset: 2px;
}
```

### Sitemap update (sitemap.go)
```go
// In buildSitemap(), alongside homepage entry:
urls = append(urls, SitemapURL{
    Loc:        baseURL + "/about",
    ChangeFreq: "yearly",
    Priority:   "0.6",
})
```

### Route registration (main.go)
```go
// Before the catch-all:
blogMux.HandleFunc("GET /about", blog.AboutPage)
blogMux.HandleFunc("GET /{path...}", blog.NotFound) // catch-all last
```

---

## Environment Availability

Step 2.6: SKIPPED — Phase 8 is purely code/template/CSS changes with no external dependencies beyond the existing Go toolchain.

---

## Validation Architecture

### Test Framework
| Property | Value |
|----------|-------|
| Framework | Go stdlib `testing` package |
| Config file | none — `go test ./...` discovers tests by convention |
| Quick run command | `go test ./internal/handler/blog/... -run TestAbout` |
| Full suite command | `go test ./...` |

All existing tests pass (verified: `ok github.com/jared-wallace/website-go/internal/handler/blog 0.381s`).

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| ABOUT-01 | `GET /about` returns HTTP 200 | unit | `go test ./internal/handler/blog/... -run TestAboutPageStatus` | ❌ Wave 0 |
| ABOUT-01 | Nav contains `href="/about"` | unit | `go test ./internal/handler/blog/... -run TestAboutNavLink` | ❌ Wave 0 |
| ABOUT-02 | Response body contains rendered markdown content | unit | `go test ./internal/handler/blog/... -run TestAboutPageContent` | ❌ Wave 0 |
| ABOUT-03 | Response contains `site-nav`, `site-footer`, `dark-toggle` chrome | unit | `go test ./internal/handler/blog/... -run TestAboutPageChrome` | ❌ Wave 0 |

Tests follow the established `package blog_test` pattern using `httptest.NewRecorder()` and `newTestHandler()` — no new test infrastructure needed.

### Sampling Rate
- **Per task commit:** `go test ./internal/handler/blog/... -run TestAbout`
- **Per wave merge:** `go test ./...`
- **Phase gate:** Full suite green before `/gsd:verify-work`

### Wave 0 Gaps
- [ ] `internal/handler/blog/about_test.go` — covers ABOUT-01 (200 status, nav link), ABOUT-02 (content), ABOUT-03 (chrome presence)

*(No new framework install needed — existing `testing` + `net/http/httptest` covers all cases.)*

---

## Sources

### Primary (HIGH confidence)
- `web/embed.go` — verified embed.FS pattern; `//go:embed templates` and `//go:embed static`
- `internal/handler/blog/handler.go` — verified per-page template set pattern, `pages` slice, `render()` method
- `internal/handler/blog/post.go` — verified ShowPost handler structure (data map, render call)
- `web/templates/base.html` — verified nav DOM order: `.site-name`, `.site-tagline`, `.dark-toggle`; confirmed `margin-left: auto` behavior
- `web/templates/post.html` — verified `{{template "base" .}}` + `{{define "content"}}` pattern
- `web/static/main.css` lines 55-89 — verified `.site-nav` flex layout, `.dark-toggle` margin-left: auto, `.site-tagline` font-size/color
- `cmd/server/main.go` lines 115-133 — verified blogMux route registration order, catch-all placement
- `internal/handler/blog/sitemap.go` — verified `buildSitemap()` structure for `/about` addition
- `internal/markdown/renderer.go` — verified `Render(src string) template.HTML` signature

### Secondary (MEDIUM confidence)
- `go test ./...` run — confirmed all 9 test packages pass, test framework state validated
- `08-UI-SPEC.md` — UI contract for `.nav-link` CSS spec, `.about-title`, page structure

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — no new dependencies; all libraries verified in place
- Architecture: HIGH — all patterns verified from existing source files
- Pitfalls: HIGH — derived from direct inspection of live CSS, template structure, and embed mechanism
- Test map: HIGH — pattern matches existing `handler_test.go` exactly

**Research date:** 2026-03-28
**Valid until:** 2026-04-28 (stable Go project; no fast-moving ecosystem concerns)
