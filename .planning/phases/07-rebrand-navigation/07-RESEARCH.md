# Phase 7: Rebrand + Navigation — Research

**Researched:** 2026-03-28
**Domain:** Go html/template string replacement, inline SVG, CSS anchor pattern, test assertion updates
**Confidence:** HIGH

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

- **D-01:** RSS icon appears in the **footer only**, next to the copyright line. Nav bar stays clean.
- **D-02:** Standard RSS broadcast SVG icon (the classic radiating-lines icon), color-matched to the site's nautical palette. Not a text link.
- **D-03:** Admin templates also rebrand to "The Wild Meridian". Nav becomes "The Wild Meridian -- Back Office", login heading updates. Full consistency, no split personality.
- **D-04:** Copyright uses dynamic `{{.Year}}` (not hardcoded). Text changes from `(c) {{.Year}} The Log` to `(c) {{.Year}} Jared Wallace`. Blog name drops out of the copyright line.

### Claude's Discretion

- RSS icon sizing and exact color values within the existing palette
- SVG icon implementation details (inline vs sprite)
- Any minor template formatting adjustments needed during the rename

### Deferred Ideas (OUT OF SCOPE)

None — discussion stayed within phase scope
</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| BRAND-01 | Blog title reads "The Wild Meridian" on all public pages (header, browser tab, etc.) | String replacements in base.html, list.html, post.html, 404.html confirmed and mapped |
| BRAND-02 | RSS feed title and metadata reflect "The Wild Meridian" | rss.go line 91 `buildRSSFeed()` Channel.Title is the single mutation point |
| BRAND-03 | XML sitemap, Open Graph meta tags, and any SEO surfaces reference "The Wild Meridian" | OG/Twitter meta in base.html default block and list.html override both contain "The Log" |
| BRAND-04 | Copyright footer reads "2026 Jared Wallace" | base.html line 51; `{{.Year}}` already in template, only string portion changes |
| NAV-01 | RSS feed is discoverable via a discreet link/icon in the nav or footer | New `.rss-link` anchor + inline SVG added to base.html footer; CSS rule added to main.css |
| NAV-02 | HTML head includes `<link rel="alternate" type="application/rss+xml">` | Already present at base.html line 12 — only the `title` attribute value needs updating |
</phase_requirements>

---

## Summary

Phase 7 is a surgical rename-and-add operation. There are no new Go packages, no schema changes, no new routes, and no new pages. Every task is either a string replacement in an existing file or the addition of a single new visual element (the RSS icon anchor) to the footer.

The codebase is clean: all 15 test packages pass, the string "The Log" appears in exactly the files listed in CONTEXT.md canonical refs, and the infrastructure for NAV-02 (`<link rel="alternate">`) is already present — it just needs its `title` attribute updated. Nothing about this phase requires external research beyond confirming the existing codebase patterns.

The one net-new implementation decision is the RSS broadcast SVG. The correct icon (circle at lower-left, two arcs radiating from it) does not exist in the project today. The shape is a well-known standard; the implementation pattern is to copy the inline SVG approach from the dark-mode toggle in `base.html` lines 29–44.

**Primary recommendation:** Execute as three sequential tasks — (1) string replacements across all templates and rss.go, (2) update test assertions, (3) add RSS icon anchor + CSS rule. Each task is independently verifiable with `go test ./...`.

---

## Standard Stack

### Core (no new dependencies — all stdlib or already in go.mod)

| Library | Version | Purpose | Why |
|---------|---------|---------|--------------|
| `html/template` stdlib | Go 1.26 | Template string replacement | Already in use; `{{define}}` / `{{block}}` inheritance already established |
| `encoding/xml` stdlib | Go 1.26 | RSS channel struct fields | Already in use in rss.go |
| Inline SVG | n/a | RSS broadcast icon | No icon library present; project uses hand-authored SVG (see dark-toggle pattern) |
| Custom CSS | n/a | `.rss-link` anchor style | All UI is hand-authored in `web/static/main.css` |

### No New Packages Required

This phase adds zero new dependencies. The go.mod file does not need to change.

---

## Architecture Patterns

### Established: Go Template Inheritance

All public templates use `{{template "base" .}}` at the top, then override named blocks:

- `{{define "title"}}` — overrides the `<title>` tag and the `{{block "title" .}}` default in base.html
- `{{define "meta"}}` — overrides the OG/Twitter `<head>` block
- `{{define "content"}}` — overrides the body content area

The base defaults (lines 6 and 13–23 of base.html) apply to any template that does NOT override the block. `list.html` overrides both `title` and `meta`; `post.html` overrides both; `404.html` overrides only `title`. This means:

- `base.html` line 6 default: `{{block "title" .}}The Log{{end}}` — used by any template that doesn't define a title block
- `base.html` lines 13–14: `og:title content="The Log"` — the OG default used by pages without a custom meta block

**All string occurrences of "The Log" that need replacement are enumerated below in the Copywriting Inventory section.**

### Established: Admin Template Inheritance

`admin-base.html` is a separate base (not inheriting from `base.html`). It has its own `{{define "base"}}` root. The `.site-name` text appears in two places:
1. The nav anchor text: `The Log -- Back Office` (admin-base.html line 16)
2. The login page `<h1>`: `The Log` (admin-login.html line 5, which overrides `{{define "nav"}}` to suppress the nav and renders its own heading)

### Established: Inline SVG with `currentColor`

The dark-mode toggle (base.html lines 29–44) uses:

```html
<button class="dark-toggle" id="dark-toggle" aria-label="Switch to dark mode">
  <svg ... stroke="currentColor" stroke-width="2" ... aria-hidden="true">
    ...
  </svg>
</button>
```

And CSS:

```css
.dark-toggle {
  color: var(--color-text-muted);
  min-width: 44px;
  min-height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.dark-toggle:hover {
  color: var(--color-accent);
}
```

The SVG inherits color from the parent element's `color` property via `stroke="currentColor"`. The RSS icon anchor must replicate this exact pattern — `<a>` owns the color, SVG inherits it.

### RSS Broadcast SVG Shape

The standard RSS broadcast icon consists of:
- A filled circle at the lower-left corner (position: `cx="3" cy="21" r="2"`, or similar)
- A short arc (inner) centered on the same lower-left origin
- A longer arc (outer) centered on the same origin

Feather-style stroke-only version (no fill, `stroke-width="2"`, 24×24 viewBox scaled down to 16×16 display):

```html
<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24"
     fill="none" stroke="currentColor" stroke-width="2"
     stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
  <path d="M4 11a9 9 0 0 1 9 9"></path>
  <path d="M4 4a16 16 0 0 1 16 16"></path>
  <circle cx="5" cy="19" r="1" fill="currentColor" stroke="none"></circle>
</svg>
```

This is the Feather `rss` icon verbatim. The viewBox stays at 0 0 24 24 while width/height are set to 16 — the browser scales it, so the stroke proportions remain correct.

### RSS Icon Anchor Pattern

```html
<!-- In base.html footer, right of the <p> copyright tag -->
<a href="/rss" class="rss-link" aria-label="RSS feed">
  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24"
       fill="none" stroke="currentColor" stroke-width="2"
       stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
    <path d="M4 11a9 9 0 0 1 9 9"></path>
    <path d="M4 4a16 16 0 0 1 16 16"></path>
    <circle cx="5" cy="19" r="1" fill="currentColor" stroke="none"></circle>
  </svg>
</a>
```

```css
/* In main.css under /* --- Footer --- */ */
.rss-link {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 44px;
  min-height: 44px;
  margin-left: 8px;
  color: var(--color-text-muted);
  text-decoration: none;
  transition: color 200ms ease;
}

.rss-link:hover {
  color: var(--color-accent);
}
```

---

## Copywriting Inventory

Complete enumeration of every "The Log" occurrence that must change, with exact file, line, old value, and new value. Verified by reading source files directly.

| File | Line | Old Value | New Value |
|------|------|-----------|-----------|
| `web/templates/base.html` | 6 | `The Log` (default title block) | `The Wild Meridian` |
| `web/templates/base.html` | 12 | `title="The Log"` (RSS alternate link) | `title="The Wild Meridian"` |
| `web/templates/base.html` | 14 | `content="The Log"` (og:title) | `content="The Wild Meridian"` |
| `web/templates/base.html` | 20 | `content="The Log"` (twitter:title) | `content="The Wild Meridian"` |
| `web/templates/base.html` | 27 | `The Log` (nav .site-name text) | `The Wild Meridian` |
| `web/templates/base.html` | 51 | `&copy; {{.Year}} The Log` | `&copy; {{.Year}} Jared Wallace` |
| `web/templates/list.html` | 3 | `The Log` (title define) | `The Wild Meridian` |
| `web/templates/list.html` | 6 | `content="The Log"` (og:title) | `content="The Wild Meridian"` |
| `web/templates/list.html` | 12 | `content="The Log"` (twitter:title) | `content="The Wild Meridian"` |
| `web/templates/post.html` | 3 | `&#8212; The Log` (title suffix) | `&#8212; The Wild Meridian` |
| `web/templates/404.html` | 3 | `&#8212; The Log` (title suffix) | `&#8212; The Wild Meridian` |
| `web/templates/admin-base.html` | 6 | `Admin{{end}} -- The Log` | `Admin{{end}} -- The Wild Meridian` |
| `web/templates/admin-base.html` | 16 | `The Log -- Back Office` | `The Wild Meridian -- Back Office` |
| `web/templates/admin-login.html` | 5 | `The Log` (h1.site-name) | `The Wild Meridian` |
| `internal/handler/blog/rss.go` | 91 | `"The Log"` (Channel.Title) | `"The Wild Meridian"` |
| `internal/handler/blog/rss_test.go` | 45 | `<title>The Log</title>` | `<title>The Wild Meridian</title>` |
| `internal/handler/blog/handler_test.go` | 175 | `og:title" content="The Log"` | `og:title" content="The Wild Meridian"` |

**Strings that must NOT change** (confirmed in source):
- `dispatches from the deep end` — tagline in base.html line 28, OG description in base.html line 15 and list.html line 7, RSS channel description in rss.go line 93
- `Back Office` — admin-login.html line 6 `<p class="login-subtitle">` stays unchanged
- `rssManagingEditor` constant in rss.go — already correct ("Jared Wallace"), no change needed

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| RSS icon SVG | Custom path math | Feather `rss` icon (copy paths verbatim) | The Feather icon set is already the project's implicit standard (dark-toggle uses Feather sun/moon icons); the RSS paths are well-tested across browsers |
| Color transitions | JS color manipulation | CSS `transition: color 200ms ease` | Matches the existing card/toggle timing already in main.css; pure CSS is sufficient |
| Touch target sizing | JS event area expansion | `min-width: 44px; min-height: 44px` on the anchor | Matches the `.dark-toggle` CSS pattern exactly |

---

## Common Pitfalls

### Pitfall 1: Forgetting the admin-base.html Title Block

**What goes wrong:** The admin base has `{{block "title" .}}Admin{{end}} -- The Log` — the suffix " -- The Log" is outside the block. Easy to miss the suffix when scanning for "The Log" if you only look at the block content.

**How to avoid:** The replacement is the entire attribute string on line 6: replace `Admin{{end}} -- The Log` with `Admin{{end}} -- The Wild Meridian`.

**Warning signs:** Admin panel browser tab still reads "Admin -- The Log" after the rename.

### Pitfall 2: post.html Title Pattern Is Different

**What goes wrong:** `list.html` has a standalone `{{define "title"}}The Log{{end}}`. `post.html` has `{{define "title"}}{{.Post.Title}} &#8212; The Log{{end}}` — "The Log" is a suffix after the post title, not a standalone value. A naive `sed` replace of `The Log` would produce correct results, but awareness of context helps when hand-editing.

**How to avoid:** Replace `&#8212; The Log` with `&#8212; The Wild Meridian` in post.html and 404.html.

### Pitfall 3: RSS Alternate Link Uses `title` Attribute (Not Content)

**What goes wrong:** The `<link rel="alternate">` uses `title="The Log"` — this is an HTML attribute, not a meta `content`. A search for `content="The Log"` misses it.

**How to avoid:** base.html line 12 replacement is `title="The Log"` → `title="The Wild Meridian"`.

### Pitfall 4: Tests Will Fail Until Assertion Strings Are Updated

**What goes wrong:** Making template/rss.go changes without updating the test assertions causes `go test ./internal/handler/blog/...` to fail immediately. This is expected — tests are the verification mechanism — but it blocks the green test gate.

**How to avoid:** Update test assertions in the same task as the source files they test, or in an immediately sequential task before the green-gate check.

### Pitfall 5: RSS Icon Missing `aria-hidden` on SVG Itself

**What goes wrong:** Screen readers encounter both the `aria-label` on the `<a>` and any text content of the SVG, causing double-announcement.

**How to avoid:** SVG gets `aria-hidden="true"`, label lives on the `<a>` as `aria-label="RSS feed"` — matches the dark-toggle pattern exactly.

---

## Runtime State Inventory

| Category | Items Found | Action Required |
|----------|-------------|-----------------|
| Stored data | None — "The Log" is not stored as a DB key, user_id, or collection name; it is display copy in templates and Go constants | None |
| Live service config | None — no external service config (n8n, Datadog, etc.) stores the blog name | None |
| OS-registered state | None — no Task Scheduler tasks, systemd units, or pm2 configs embed "The Log" | None |
| Secrets/env vars | None — no .env variables or SOPS keys reference the blog name | None |
| Build artifacts | None — no compiled artifacts embed the blog name; Docker image is rebuilt fresh | None |

"The Log" exists only as string literals in source files. This is a pure code-edit rename with no data migration required.

---

## Environment Availability

Step 2.6: All dependencies for this phase are already available — it is a pure code/template edit phase using tools already confirmed operational.

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Go toolchain | `go test ./...` validation | Yes | 1.26.x (per CLAUDE.md) | — |
| `web/static/main.css` | CSS rule addition | Yes | n/a (file exists) | — |
| `web/templates/*.html` | String replacements | Yes | all files confirmed readable | — |

No missing dependencies.

---

## Validation Architecture

### Test Framework

| Property | Value |
|----------|-------|
| Framework | Go stdlib `testing` package |
| Config file | none — `go test` discovers tests by convention |
| Quick run command | `go test ./internal/handler/blog/...` |
| Full suite command | `go test ./...` |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| BRAND-01 | Browser title and nav render "The Wild Meridian" | integration (template render) | `go test ./internal/handler/blog/... -run TestListOGMeta` | Yes — handler_test.go |
| BRAND-02 | RSS channel title is "The Wild Meridian" | unit | `go test ./internal/handler/blog/... -run TestServeRSS` | Yes — rss_test.go line 45 |
| BRAND-03 | OG meta tag og:title contains "The Wild Meridian" | integration | `go test ./internal/handler/blog/... -run TestListOGMeta` | Yes — handler_test.go line 175 |
| BRAND-04 | Copyright footer renders `{{.Year}} Jared Wallace` | manual smoke | `go build ./cmd/server && curl http://localhost:8080/` | n/a — no dedicated test |
| NAV-01 | RSS icon rendered in footer HTML | manual smoke | `go build ./cmd/server && curl http://localhost:8080/ | grep rss-link` | n/a — no dedicated test |
| NAV-02 | `<link rel="alternate" type="application/rss+xml">` in rendered head | integration | `go test ./internal/handler/blog/... -run TestListOGMeta` (body includes head) | Yes — indirectly covered |

### Sampling Rate

- **Per task commit:** `go test ./internal/handler/blog/...`
- **Per wave merge:** `go test ./...`
- **Phase gate:** Full suite green before `/gsd:verify-work`

### Wave 0 Gaps

None — existing test infrastructure covers all phase requirements. No new test files or fixtures needed. The existing `TestServeRSS` and `TestListOGMeta` tests become the phase gate once their assertion strings are updated.

---

## Code Examples

### Pattern: Template Block String Replacement (base.html)

Current (to replace):
```html
<title>{{block "title" .}}The Log{{end}}</title>
<link rel="alternate" type="application/rss+xml" title="The Log" href="/rss">
<meta property="og:title" content="The Log">
<meta name="twitter:title" content="The Log">
<a href="/" class="site-name">The Log</a>
<p>&copy; {{.Year}} The Log</p>
```

Target (after replacement):
```html
<title>{{block "title" .}}The Wild Meridian{{end}}</title>
<link rel="alternate" type="application/rss+xml" title="The Wild Meridian" href="/rss">
<meta property="og:title" content="The Wild Meridian">
<meta name="twitter:title" content="The Wild Meridian">
<a href="/" class="site-name">The Wild Meridian</a>
<p>&copy; {{.Year}} Jared Wallace</p>
```

### Pattern: RSS Icon Anchor in Footer

Replace the current footer block in base.html:
```html
<footer class="site-footer">
  <hr class="rope-divider">
  <p>&copy; {{.Year}} The Log</p>
</footer>
```

With:
```html
<footer class="site-footer">
  <hr class="rope-divider">
  <p>
    &copy; {{.Year}} Jared Wallace
    <a href="/rss" class="rss-link" aria-label="RSS feed">
      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24"
           fill="none" stroke="currentColor" stroke-width="2"
           stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
        <path d="M4 11a9 9 0 0 1 9 9"></path>
        <path d="M4 4a16 16 0 0 1 16 16"></path>
        <circle cx="5" cy="19" r="1" fill="currentColor" stroke="none"></circle>
      </svg>
    </a>
  </p>
</footer>
```

### Pattern: CSS Rule Addition (main.css, under `/* --- Footer --- */`)

Append after the existing `.site-footer p` rule:
```css
.rss-link {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 44px;
  min-height: 44px;
  margin-left: 8px;
  color: var(--color-text-muted);
  text-decoration: none;
  transition: color 200ms ease;
}

.rss-link:hover {
  color: var(--color-accent);
}
```

### Pattern: RSS Handler Constant Update (rss.go)

Current (line 91):
```go
Title: "The Log",
```

Target:
```go
Title: "The Wild Meridian",
```

---

## Project Constraints (from CLAUDE.md)

| Directive | Impact on This Phase |
|-----------|---------------------|
| Go with minimal dependencies — prefer stdlib | Confirmed: zero new deps required |
| Must run as Docker container on port 8080 | No impact — no infrastructure changes |
| All persistent data on EBS at /var/www/html | No impact — no storage changes |
| Design: leverage `frontend-design` skill | No skill directory present; UI-SPEC.md has been approved and serves as the design contract |
| GSD Workflow Enforcement before file edits | All file edits proceed through execute-phase |
| Run code through `/simplify` skill before presenting | Planner should include simplify review as a verification step |

---

## Sources

### Primary (HIGH confidence)

- Direct file inspection: `web/templates/base.html` — all "The Log" occurrences enumerated from source
- Direct file inspection: `web/templates/list.html`, `post.html`, `404.html`, `admin-base.html`, `admin-login.html` — confirmed occurrences
- Direct file inspection: `internal/handler/blog/rss.go` — Channel.Title line 91 confirmed
- Direct file inspection: `internal/handler/blog/rss_test.go` — assertion line 45 confirmed
- Direct file inspection: `internal/handler/blog/handler_test.go` — assertion line 175 confirmed
- Direct file inspection: `web/static/main.css` — footer CSS section and dark-toggle pattern confirmed
- `go test ./...` execution — all 15 packages pass before any changes

### Secondary (MEDIUM confidence)

- Feather icon set RSS icon paths — standard paths used by thousands of projects; no official URL needed, paths are geometric definitions

---

## Metadata

**Confidence breakdown:**
- Copywriting inventory: HIGH — every occurrence verified by reading source files
- RSS icon SVG: HIGH — Feather rss icon is a stable, widely-used open-source asset
- CSS pattern: HIGH — modelled directly on the dark-toggle implementation in the same file
- Test impact: HIGH — test files read directly, assertions on known lines

**Research date:** 2026-03-28
**Valid until:** Indefinite — this is a rename of static strings in a stable codebase; no external dependencies, no API versions, no framework updates affect this research
