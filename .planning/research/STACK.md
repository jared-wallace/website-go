# Technology Stack

**Project:** website-go — personal blog server for jared-wallace.com
**Researched:** 2026-03-26 (Go stack), 2026-03-28 (CSS polish additions)
**Philosophy:** stdlib-first; add a dependency only when it earns its place

---

## Recommended Stack

### Runtime

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| Go | 1.26.1 | Language / runtime | Latest stable (released 2026-03-05); Go 1.22+ ServeMux rewrites eliminate the need for a router package; single static binary |

### HTTP Routing

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `net/http` stdlib | (Go 1.26) | HTTP server and routing | Go 1.22 added method-based patterns, wildcard path segments (`{id}`), and `r.PathValue()`. For a blog with ~15 routes this is sufficient without pulling in chi or gorilla/mux |

The standard `http.ServeMux` now handles `GET /posts/{slug}`, `POST /admin/posts/{id}`, etc. No external router is needed for this scope.

### Templating

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `html/template` stdlib | (Go 1.26) | Server-side HTML rendering | Automatic context-aware HTML escaping; supports template inheritance via `{{block}}` / `{{define}}`; zero dependencies; XSS protection by default |

Avoid `text/template` for HTML output — it does not escape. Avoid `templ` (compile-time code generation) unless the team wants a build-time step; html/template is sufficient for a blog.

### Markdown Rendering

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `github.com/yuin/goldmark` | v1.8.2 | Markdown → HTML | CommonMark-compliant; extensible AST; actively maintained; used by Hugo. Blackfriday is archived and not CommonMark-compliant. Goldmark ships with tables, strikethrough, task lists |

Extensions to include at wiring time:
- `goldmark/extension` — GFM tables, strikethrough, linkify
- `goldmark-meta` (`github.com/yuin/goldmark-meta`) — YAML front matter extraction for post metadata
- `goldmark-highlighting` or `github.com/alecthomas/chroma` — syntax highlighting in code fences (stretch goal)

### Database Driver

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `github.com/jackc/pgx/v5` | v5.9.1 | PostgreSQL driver + toolkit | `lib/pq` is maintenance-only. pgx v5 is faster (up to 70x in bulk ops), supports LISTEN/NOTIFY, COPY, and native `pgxpool`. Use the native pgx interface, not the `database/sql` adapter, to avoid feature loss |
| `github.com/jackc/pgx/v5/pgxpool` | (bundled with pgx v5) | Connection pooling | Single-instance blog needs only a small pool (max 10 connections); pgxpool handles lifetime and health checks |

Do not use GORM or sqlx. Raw SQL with pgx is readable, auditable, and avoids ORM magic for a schema this simple.

### Database Migrations

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `github.com/pressly/goose/v3` | v3.27.0 | Schema migrations | Sequential numbered SQL files (`00001_create_posts.sql`); supports embedded migrations via `go:embed`; single binary deployment; `goose up` in container entrypoint. golang-migrate is a valid alternative but goose has cleaner Go embedding support |

### Session Management

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `github.com/alexedwards/scs/v2` | v2.9.0 | HTTP session management | OWASP-aligned design; server-side session tokens (not cookie-payload); context-middleware pattern; has a `pgxstore` sub-package for Postgres-backed sessions. Gorilla/sessions stores payload in cookie by default which is less secure; SCS is also faster and smaller |
| `github.com/alexedwards/scs/pgxstore` | (bundled in scs repo) | Postgres session backend | Stores sessions in `sessions` table alongside app data; survives container restarts; eliminates need for Redis |

### Password Hashing

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `golang.org/x/crypto/bcrypt` | v0.49.0 | Admin password hash | The standard Go bcrypt implementation; not in `crypto/` stdlib but from the official `x/crypto` extended library (Anthropic-controlled). Use cost 12 minimum for 2025 hardware |

This is one dependency that cannot be replaced by stdlib — Go's `crypto/` package does not include bcrypt. `x/crypto` is maintained by the Go team and is effectively stdlib-adjacent.

### RSS Feed

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `encoding/xml` stdlib | (Go 1.26) | RSS 2.0 feed generation | An RSS feed is a straightforward XML document; defining a Go struct with XML annotations and using `xml.NewEncoder` requires zero dependencies. `gorilla/feeds` adds a dependency for ~50 lines of struct definitions |

Hand-roll RSS using `encoding/xml`. The RSS 2.0 spec is stable and the struct is small. This removes a dependency entirely.

### Logging

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `log/slog` stdlib | (Go 1.26) | Structured logging | Added in Go 1.21; JSON or text output; leveled; zero dependencies. No need for zerolog or zap at blog scale |

### Infrastructure

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| Docker (multi-stage) | 27.x | Container build | Build stage: `golang:1.26-alpine`; runtime stage: `alpine:3.21` (not scratch — needs CA certs and timezone data for Postgres TLS). Set `CGO_ENABLED=0` for a static binary |
| docker-compose | v2.x | Local dev + prod orchestration | Sidecar Postgres container wired to EBS volume at `/var/www/html/postgres-data`; matches production topology in dev |
| Makefile | (GNU make) | Task runner | `make dev`, `make test`, `make build`, `make migrate-up`; documents tribal knowledge |

---

## CSS Design Polish Techniques (v1.2 Shore Leave Polish)

**Constraint:** Zero new dependencies. No build tools. No JavaScript animation libraries. Pure CSS + inline SVG only.

### 1. Background Noise / Grain Texture

**Technique:** SVG `feTurbulence` filter embedded inline in HTML, applied as a fixed overlay via CSS pseudo-element.

**Why this approach:** A data-URI SVG with `feTurbulence` generates Perlin noise at render time — no raster image download, no extra HTTP request, scales infinitely without quality loss. The pseudo-element pattern with `pointer-events: none` keeps all interaction intact.

**Pattern:**

Add a hidden SVG filter definition once in `base.html`, before `</body>`:

```html
<!-- Grain filter — zero dimensions, no layout impact -->
<svg width="0" height="0" aria-hidden="true" style="position:absolute">
  <filter id="grain" color-interpolation-filters="sRGB"
          x="0%" y="0%" width="100%" height="100%">
    <feTurbulence type="fractalNoise" baseFrequency="0.72" numOctaves="4"
                  stitchTiles="stitch" result="noise"/>
    <feColorMatrix type="saturate" values="0" in="noise" result="grayNoise"/>
    <feBlend in="SourceGraphic" in2="grayNoise" mode="multiply" result="blended"/>
    <feComposite in="blended" in2="SourceGraphic" operator="in"/>
  </filter>
</svg>
```

Apply in CSS as a full-page fixed overlay on `body::after`:

```css
body {
  position: relative; /* stacking context for the overlay */
}

body::after {
  content: "";
  position: fixed;
  inset: 0;               /* top/right/bottom/left: 0 shorthand */
  pointer-events: none;   /* clicks pass through */
  z-index: 9999;
  opacity: 0.035;         /* subtle — 3-4% is the target range */
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)'/%3E%3C/svg%3E");
  background-repeat: repeat;
  background-size: 256px 256px;
}
```

**Alternative (simpler, same visual result):** Use `mix-blend-mode: multiply` on a `position: fixed` `<div>` added at the end of `<body>`. The pseudo-element approach is cleaner because it requires no template markup change.

**Performance notes:**
- `numOctaves` must stay at 4 or below. Values above 4 double processing cost with negligible visual gain. Confirmed by Frontend Masters and Codrops documentation.
- `position: fixed` on the overlay element lets the browser promote it to its own GPU compositor layer (same as `will-change: transform`), avoiding repaint on scroll.
- `background-size: 256px` tiles a small SVG; the tile boundary is seamless because `stitchTiles="stitch"` is set.
- **Mobile impact:** The overlay is a fixed pseudo-element with no JS. GPU compositing handles it. On low-end Android devices the SVG filter re-renders on DOM paint; keep `opacity` at 0.03–0.05 and `baseFrequency` at 0.65–0.95 to minimize rasterization cost. If performance testing reveals an issue on mobile, fall back to `opacity: 0` on the overlay behind `@media (hover: none)`.

**Dark mode adjustment:** Grain opacity reads differently on dark backgrounds. Add:

```css
[data-theme="dark"] body::after {
  opacity: 0.06; /* slightly more visible on dark to maintain atmosphere */
  mix-blend-mode: screen; /* screen blends brighter on dark surfaces */
}
```

**Browser support:** `feTurbulence` is supported in all modern browsers (Chrome 8+, Firefox 3+, Safari 6+, Edge 12+). The inline SVG data-URI approach has been the standard technique since at least 2019 per CSS-Tricks. HIGH confidence.

---

### 2. Page Entry Animations (Fade-in + Card Stagger)

**Technique:** `@keyframes` fade-in on body/main, CSS custom property stagger delays on post cards. No JavaScript.

**Why this approach:** CSS `animation-fill-mode: backwards` ensures elements start invisible even before the animation fires (critical — without it, content flashes at full opacity during the delay). The `--animation-order` custom property pattern avoids writing N separate `nth-child` rules.

**Base fade-in (page-level):**

```css
@keyframes fade-in {
  from {
    opacity: 0;
    transform: translateY(8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.container {
  animation: fade-in 400ms ease-out both;
  /* `both` = shorthand for fill-mode forwards + backwards */
}
```

**Card stagger (list page):**

Each post card in the Go template receives an inline `--animation-order` variable. The `list.html` template loop already generates cards with `{{range $i, $p := .Posts}}`, so use `$i` to set the variable:

```html
<!-- In list.html -->
<a class="post-card" href="/posts/{{.Slug}}"
   style="--animation-order: {{$i}}">
  ...
</a>
```

```css
.post-card {
  animation: fade-in 350ms ease-out both;
  animation-delay: calc(var(--animation-order, 0) * 60ms);
}
```

At 60ms per card with a 2-column grid (4 cards visible), the last visible card enters at 180ms — well within the "feels instant" threshold. Beyond 8 cards the delay caps out at ~500ms; add a cap if pagination means large sets:

```css
/* Cap stagger at 6 steps so page 2 doesn't feel like a slow waterfall */
.post-card:nth-child(n+7) {
  animation-delay: calc(6 * 60ms);
}
```

**Accessibility — prefers-reduced-motion:**

This is mandatory, not optional. Users with vestibular disorders and motion sensitivity set this OS preference:

```css
@media (prefers-reduced-motion: reduce) {
  .container,
  .post-card {
    animation: none;
  }
}
```

Placing this at the end of the animation section ensures it overrides without specificity fights.

**Browser support:** `animation-fill-mode: both`, `@keyframes`, and CSS custom properties in `animation-delay: calc(...)` work in all modern browsers. The `--animation-order` pattern requires custom property support (Chrome 49+, Firefox 31+, Safari 9.1+). HIGH confidence.

---

### 3. Dark Mode Color Transitions

**Technique A (recommended): `transition` on `body` for structural colors.**

The simplest reliable approach: add a `transition` on `background-color` and `color` to `body`. Because all colors are CSS custom properties applied to elements, and those elements inherit from body, the transition cascades.

```css
body {
  transition:
    background-color 250ms ease,
    color 250ms ease;
}

/* Extend to surfaces that do NOT inherit body color directly */
.site-nav,
.post-card,
.toc,
.site-footer {
  transition:
    background-color 250ms ease,
    border-color 250ms ease,
    color 250ms ease;
}
```

**Important gotcha:** CSS custom properties (`--color-bg`, `--color-surface`) do not animate on their own — only concrete CSS properties like `background-color` transition. Listing every element that uses a background color is tedious but reliable.

**Technique B (modern, less boilerplate): `@property` with typed color variables.**

`@property` (Baseline 2024, ~94% browser support as of March 2026) registers a custom property with a type, enabling the browser to interpolate between color values:

```css
@property --color-bg {
  syntax: "<color>";
  inherits: true;
  initial-value: #F5F0E8;
}

@property --color-surface {
  syntax: "<color>";
  inherits: true;
  initial-value: #E8DFD0;
}

@property --color-accent {
  syntax: "<color>";
  inherits: true;
  initial-value: #2C5F7A;
}

/* Now a single transition rule on :root handles everything */
:root {
  transition:
    --color-bg 250ms ease,
    --color-surface 250ms ease,
    --color-accent 250ms ease,
    --color-text 250ms ease,
    --color-text-muted 250ms ease,
    --color-divider 250ms ease;
}
```

When `[data-theme="dark"]` is toggled, the browser interpolates the custom property values as colors, and every element consuming those variables gets the transition for free. Zero per-element transition rules needed.

**Recommendation:** Use Technique B (`@property`) as the primary approach. Provide Technique A as fallback commentary. The existing project already has a theme-toggle JS snippet in `base.html` that sets `data-theme`; both techniques work with this mechanism unchanged.

**One edge case to handle:** The inline `<script>` in `<head>` that sets `data-theme` before CSS loads prevents FOUC. When the page first loads in dark mode, `@property` transitions will not fire (good — no flash of light-mode colors). Transitions only fire on subsequent toggle clicks, which is the desired behavior.

**Browser support:** `@property` — Baseline 2024 (all modern browsers since July 2024). Older browsers fall back gracefully: colors switch instantly with no transition, dark mode still works. HIGH confidence.

---

### 4. Rope / Nautical SVG Dividers

**Technique:** Inline SVG `<path>` elements with a stroke that mimics rope — sinusoidal or hand-drawn wave with stroke-dasharray for knot texture. No images, no external assets.

**Why inline SVG over CSS `border-style: dashed`:** The existing `.rope-divider` is `border-top: 2px dashed`. This is flat and unconvincing. An inline SVG path gives control over stroke-width, stroke-dasharray rhythm, and color via `currentColor` (inherits from CSS), so it respects both light and dark mode automatically.

**Pattern A — Wavy rope line (simple, high fidelity):**

```html
<!-- Drop-in replacement for <hr class="rope-divider"> -->
<div class="rope-divider" aria-hidden="true">
  <svg xmlns="http://www.w3.org/2000/svg"
       viewBox="0 0 800 20" preserveAspectRatio="none"
       width="100%" height="20">
    <path
      d="M0,10 C50,2 100,18 150,10 C200,2 250,18 300,10
         C350,2 400,18 450,10 C500,2 550,18 600,10
         C650,2 700,18 750,10 C800,2 800,10 800,10"
      fill="none"
      stroke="currentColor"
      stroke-width="2.5"
      stroke-linecap="round"
      stroke-dasharray="8 4"
      opacity="0.5"/>
  </svg>
</div>
```

The `stroke-dasharray="8 4"` (8px dash, 4px gap) on the wavy path creates the illusion of twisted rope strands. `preserveAspectRatio="none"` with `width="100%"` makes it span the container width at any viewport.

**Pattern B — Double-strand rope (more nautical character):**

```html
<div class="rope-divider" aria-hidden="true">
  <svg xmlns="http://www.w3.org/2000/svg"
       viewBox="0 0 800 24" preserveAspectRatio="none"
       width="100%" height="24">
    <!-- Strand 1: offset up -->
    <path
      d="M0,8 C40,2 80,14 120,8 C160,2 200,14 240,8
         C280,2 320,14 360,8 C400,2 440,14 480,8
         C520,2 560,14 600,8 C640,2 680,14 720,8 C760,2 800,8 800,8"
      fill="none" stroke="currentColor" stroke-width="2"
      stroke-linecap="round" stroke-dasharray="6 3" opacity="0.6"/>
    <!-- Strand 2: offset down, phase-shifted -->
    <path
      d="M0,16 C40,22 80,10 120,16 C160,22 200,10 240,16
         C280,22 320,10 360,16 C400,22 440,10 480,16
         C520,22 560,10 600,16 C640,22 680,10 720,16 C760,22 800,16 800,16"
      fill="none" stroke="currentColor" stroke-width="2"
      stroke-linecap="round" stroke-dasharray="6 3" opacity="0.6"/>
  </svg>
</div>
```

The two phase-shifted strands read as a twisted rope. Color via `currentColor` means it inherits `--color-divider` automatically.

**CSS for the divider container:**

```css
.rope-divider {
  /* Replace the existing border-top rule entirely */
  border: none;
  margin: 24px 0;
  line-height: 0;  /* collapse inline-block gap */
  color: var(--color-divider);
}
```

**Positioning in footer:** The footer `.rope-divider` used `border-top`. With SVG it becomes a block element; the footer layout already wraps it correctly (the existing `margin-bottom: 16px` on `.site-footer .rope-divider` stays valid).

**Browser support:** Inline SVG with `<path>`, `stroke-dasharray`, and `currentColor` is universally supported (Chrome 4+, Firefox 3+, Safari 3.2+). HIGH confidence.

**What NOT to do:** Do not fetch an external SVG file for a divider. An HTTP request for a decorative 400-byte SVG is not worth it. Do not use CSS `background-image: url(rope.svg)` — that also requires a network fetch and loses `currentColor` inheritance.

---

## Alternatives Considered

| Category | Recommended | Alternative | Why Not |
|----------|-------------|-------------|---------|
| Router | `net/http` ServeMux (stdlib) | `go-chi/chi`, `gorilla/mux` | Not needed — Go 1.22 ServeMux handles method+path patterns. Chi is excellent but adds a dependency for functionality now in stdlib |
| Templating | `html/template` (stdlib) | `templ` | Templ requires a code-generation build step and additional toolchain dependency. html/template is sufficient for a blog |
| Markdown | goldmark | blackfriday v2 | Blackfriday is not CommonMark-compliant and is effectively unmaintained. Gitea migrated away from it in 2020 |
| DB Driver | pgx v5 | `lib/pq`, GORM | lib/pq is maintenance-only. GORM is an ORM — unnecessary complexity for a blog schema |
| Migrations | goose v3 | golang-migrate | Both are valid; goose has cleaner `go:embed` support for embedded SQL migrations in a single binary |
| Sessions | scs v2 | gorilla/sessions, roll-your-own | gorilla/sessions stores payload in cookie (larger surface area); gorilla toolkit was briefly unmaintained. SCS is focused, secure, and has a pgxstore |
| RSS | `encoding/xml` stdlib | gorilla/feeds | gorilla/feeds is ~50 lines of struct defs; not worth the dependency |
| Logging | `log/slog` stdlib | zerolog, zap | At blog scale, slog is fast enough and avoids a dep. Revisit if structured log ingestion becomes a requirement |
| Image storage | EBS volume (`/var/www/html/images`) | S3 | Per PROJECT.md decision: EBS is sufficient for a single-blog, avoids AWS SDK dependency |
| Grain texture | Inline SVG data-URI pseudo-element | PNG image file, canvas, JS library | PNG requires extra HTTP request; canvas requires JS; external library violates minimal-deps constraint |
| Animation | CSS `@keyframes` + `animation-delay` | GSAP, Framer Motion, AOS.js | All require JavaScript bundles. CSS-only `@keyframes` with `animation-fill-mode: both` is sufficient for blog-scale entry animations |
| Dark mode transition | CSS `@property` typed variables | Per-element `transition` on every selector | `@property` on `:root` handles all elements in one rule; Baseline 2024 means ~94% support; graceful degradation (instant switch) on older browsers |
| Rope divider | Inline SVG path with `stroke-dasharray` | CSS `border: dashed`, PNG image | CSS dashed border is flat/unconvincing; PNG loses `currentColor` and requires HTTP fetch |

---

## What NOT to Use

| Avoid | Why | Use Instead |
|-------|-----|-------------|
| GSAP / Anime.js / AOS.js | External JS dependency for what 20 lines of CSS handles | `@keyframes` + `animation-delay: calc(...)` |
| Lottie animations | JSON bundle + JS runtime — wildly over-engineered for a blog | Static inline SVG |
| `will-change: transform` on every animated element | Causes excessive GPU layer promotion, higher memory use on mobile | Use only on the grain overlay `fixed` element where compositing is genuinely needed |
| `numOctaves` > 4 on `feTurbulence` | Exponential render cost, negligible visual gain past 4 | Keep at 4; use `baseFrequency` to tune grain size instead |
| `@keyframes` on `color` or `background-color` in keyframes for theme switch | Bypasses `prefers-color-scheme` and user preference; creates flash | Use CSS `transition` (respects reduced-motion) or `@property` |
| Sass/PostCSS build step | Violates project constraint of no build tools | CSS custom properties with `calc()` cover all stagger/theming needs in vanilla CSS |
| External SVG file for divider | HTTP request for decorative asset; loses `currentColor` | Inline SVG in template HTML |

---

## Complete Dependency List

```
github.com/yuin/goldmark           v1.8.2
github.com/yuin/goldmark-meta      (latest, tracks goldmark)
github.com/jackc/pgx/v5            v5.9.1
github.com/alexedwards/scs/v2      v2.9.0
golang.org/x/crypto                v0.49.0
github.com/pressly/goose/v3        v3.27.0
```

Six external dependencies total. Everything else is stdlib. CSS polish adds zero new dependencies.

---

## Installation

```bash
go get github.com/yuin/goldmark@v1.8.2
go get github.com/yuin/goldmark-meta@latest
go get github.com/jackc/pgx/v5@v5.9.1
go get github.com/alexedwards/scs/v2@v2.9.0
go get github.com/alexedwards/scs/pgxstore@latest
go get golang.org/x/crypto@v0.49.0
go get github.com/pressly/goose/v3@v3.27.0

# CLI tool for running migrations
go install github.com/pressly/goose/v3/cmd/goose@latest
```

CSS polish requires no installation — all changes are to `web/static/main.css` and `web/templates/base.html`.

---

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Go version (1.26.1) | HIGH | Verified via go.dev/doc/devel/release (2026-03-05) |
| net/http ServeMux routing | HIGH | Verified via go.dev/blog/routing-enhancements; official blog |
| goldmark v1.8.2 | HIGH | Verified via pkg.go.dev (published 2026-03-25) |
| pgx v5.9.1 | HIGH | Verified via pkg.go.dev (published 2026-03-22) |
| scs v2.9.0 | HIGH | Verified via pkg.go.dev (published 2025-04-17) |
| goose v3.27.0 | HIGH | Verified via pkg.go.dev (published 2026-02-22) |
| golang.org/x/crypto v0.49.0 | HIGH | Verified via pkg.go.dev (published 2026-03-11) |
| RSS via encoding/xml | HIGH | Standard library, well-documented pattern |
| Docker alpine runtime base | MEDIUM | Alpine 3.21 assumed current; verify at build time |
| SVG feTurbulence grain technique | HIGH | Verified via CSS-Tricks, Frontend Masters, FreeCodeCamp; technique has been standard since 2019 |
| CSS @keyframes stagger with custom properties | HIGH | Verified via CSS-Tricks official article and MDN |
| @property for dark mode transitions | HIGH | Verified via MDN; Baseline 2024 (July 2024); ~94% browser support |
| Inline SVG path rope divider | HIGH | Standard SVG spec, universally supported; `currentColor` inheritance is CSS spec |
| Grain texture mobile performance | MEDIUM | `position: fixed` + `pointer-events: none` is the established safe pattern; actual frame budget varies by device; recommend testing on real low-end Android |

---

## Sources

- [Go Release History](https://go.dev/doc/devel/release) — Go 1.26.1 confirmed latest
- [Routing Enhancements for Go 1.22](https://go.dev/blog/routing-enhancements) — Official Go blog on ServeMux improvements
- [goldmark on pkg.go.dev](https://pkg.go.dev/github.com/yuin/goldmark) — v1.8.2, March 2026
- [pgx on pkg.go.dev](https://pkg.go.dev/github.com/jackc/pgx/v5) — v5.9.1, March 2026
- [scs on pkg.go.dev](https://pkg.go.dev/github.com/alexedwards/scs/v2) — v2.9.0
- [goose on pkg.go.dev](https://pkg.go.dev/github.com/pressly/goose/v3) — v3.27.0
- [golang.org/x/crypto/bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) — v0.49.0
- [The Go Ecosystem in 2025 — JetBrains GoLand Blog](https://blog.jetbrains.com/go/2025/11/10/go-language-trends-ecosystem-2025/)
- [SCS Session Manager — Alex Edwards](https://www.alexedwards.net/blog/scs-session-manager)
- [pgx vs lib/pq — Preslav Rachev](https://preslav.me/2022/05/13/pq-or-pgx-choosing-the-right-postgresql-postgresql-golang-driver/)
- [Go net/http ServeMux is All You Need — DEV Community](https://dev.to/leapcell/gos-httpservemux-is-all-you-need-1mam)
- [Gitea migration from blackfriday to goldmark](https://github.com/go-gitea/gitea/pull/9533)
- [Grainy Gradients — Frontend Masters Blog](https://frontendmasters.com/blog/grainy-gradients/) — feTurbulence technique, performance notes, numOctaves guidance
- [How to Create Grainy CSS Backgrounds Using SVG Filters — FreeCodeCamp](https://www.freecodecamp.org/news/grainy-css-backgrounds-using-svg-filters/) — inline SVG data-URI approach, overlay implementation
- [SVG Filter Effects: Creating Texture with feTurbulence — Codrops](https://tympanus.net/codrops/2019/02/19/svg-filter-effects-creating-texture-with-feturbulence/) — baseFrequency/numOctaves reference
- [Different Approaches for Creating a Staggered Animation — CSS-Tricks](https://css-tricks.com/different-approaches-for-creating-a-staggered-animation/) — CSS custom property stagger pattern
- [CSS @property — MDN](https://developer.mozilla.org/en-US/docs/Web/CSS/@property) — Baseline 2024 confirmation, typed custom property animation
- [prefers-reduced-motion — MDN](https://developer.mozilla.org/en-US/docs/Web/CSS/@media/prefers-reduced-motion) — accessibility media query
- [Using CSS transitions — MDN](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_transitions/Using_CSS_transitions) — transition property reference
