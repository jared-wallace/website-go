# Domain Pitfalls

**Domain:** Go blog server (personal blog + admin panel, Dockerized, behind AWS ALB + Nginx)
**Researched:** 2026-03-26
**Project:** website-go / jared-wallace.com

---

## Critical Pitfalls

Mistakes that cause rewrites, data loss, or security incidents.

---

### Pitfall 1: Session Cookie Missing Security Flags

**What goes wrong:** The admin session cookie is set without `HttpOnly`, `Secure`, and `SameSite=Lax` flags. JavaScript can read the cookie (`HttpOnly` missing), the cookie travels over plain HTTP to Nginx's internal 8080 listener (`Secure` ignored on internal traffic), and there is no CSRF defense (`SameSite` missing).

**Why it happens:** Go's `net/http` `http.Cookie` struct sets all security fields to their zero values by default. Nothing reminds you to set them. The server terminates TLS at ALB so the developer sees only HTTP internally and reasons "I don't need Secure." That reasoning is wrong — `Secure` also tells the browser not to send the cookie over cleartext links that bypass the proxy.

**Consequences:** XSS on any page can steal the session token. An admin click on a malicious link triggers a CSRF write action (new post, delete post, etc.).

**Prevention:**
```go
http.SetCookie(w, &http.Cookie{
    Name:     "session",
    Value:    token,
    Path:     "/",
    HttpOnly: true,
    Secure:   true,   // set even though ALB terminates TLS
    SameSite: http.SameSiteLaxMode,
    MaxAge:   86400,
})
```
Additionally, regenerate the session ID on login (destroy pre-auth session, issue new one) to prevent session fixation. Add a CSRF token for all state-mutating admin form submissions — Go 1.25+ includes `CrossOriginProtection` in stdlib, or use `github.com/gorilla/csrf` for earlier versions.

**Detection:** Missing flags are visible in browser DevTools → Application → Cookies. Any cookie without the lock icon is a red flag.

**Phase:** Address in the session auth phase (admin panel foundation).

---

### Pitfall 2: Markdown Rendered to `template.HTML` Without Sanitization

**What goes wrong:** Goldmark renders markdown to HTML, then the result is cast to `template.HTML` and injected into the Go template. Go's `html/template` trusts `template.HTML` values completely — it will not escape them. If the rendered output contains `<script>` tags or `javascript:` URLs sourced from the markdown, they execute in the reader's browser.

**Why it happens:** The developer correctly notes "I'm the only author, so I trust the content" and skips sanitization. This is sound logic for today's posts, but it creates two hidden risks: (1) if the admin session is ever compromised, an attacker can inject a post with a payload; (2) any future "API push" endpoint that accepts `.md` files provides an attack vector if the local machine is compromised or the API endpoint auth is weak.

**Consequences:** Stored XSS. A malicious post persisted in the database executes on every reader's browser.

**Prevention:** Run `bluemonday` after Goldmark, even for admin-authored content. Use `bluemonday.UGCPolicy()` as the baseline and add any blog-specific allowances (e.g., `<figure>`, `<figcaption>`, code block classes). The pipeline is:

```
markdown string → goldmark.Convert() → bluemonday.Sanitize() → template.HTML
```

Never set `unsafe: true` on the goldmark renderer unless you are certain bluemonday is running downstream of it.

**Detection:** Add a test post containing `<script>alert(1)</script>` — if the alert fires, sanitization is absent or mis-ordered.

**Phase:** Address in the markdown rendering phase, before any public-facing rendering ships.

---

### Pitfall 3: Postgres Data on EBS Bind-Mount with Wrong Permissions

**What goes wrong:** The `docker-compose.yml` uses a bind mount pointing at an EBS-backed directory (e.g., `/var/www/html/postgres-data`) for Postgres. On first startup everything works. After a `docker compose down && docker compose up`, or after a host reboot, Postgres refuses to start with: `FATAL: data directory "/var/lib/postgresql/data" has wrong ownership`.

**Why it happens:** The Postgres container runs as UID/GID 999 (`postgres`). The EBS directory on the host is owned by `root` (created by the OS or by `mkdir`). Docker bind mounts preserve host ownership — they do not remap UIDs the way named volumes might. Postgres strictly requires `0700` permissions owned by the user running the process.

**Consequences:** Database fails to start. All writes since the last `docker compose down` are lost if you wipe the directory to fix permissions. Worst case: data loss on instance replacement.

**Prevention:**
1. Create the host directory with explicit ownership before first run:
   ```bash
   sudo mkdir -p /var/www/html/postgres-data
   sudo chown 999:999 /var/www/html/postgres-data
   sudo chmod 700 /var/www/html/postgres-data
   ```
2. Document this in the deployment runbook and Makefile `deploy` target.
3. Alternatively, use a Docker named volume with a custom driver pointing to the EBS mount path — this lets Docker manage ownership internally.
4. Set `PGDATA` in the compose file to avoid the default `/var/lib/postgresql/data/pgdata` subdirectory confusion.

**Detection:** `docker compose logs postgres` on startup. Any mention of "wrong ownership" or "permission denied" pointing at `/var/lib/postgresql/data` is this pitfall.

**Phase:** Address in the Docker Compose / infrastructure phase before any data is written to production.

---

### Pitfall 4: Image Upload Accepts MIME Type from Client Header

**What goes wrong:** The upload handler reads `r.Header.Get("Content-Type")` or the multipart file's MIME type provided by the browser to decide whether an upload is a valid image. Browsers let users (and scripts) set any Content-Type. An attacker uploads an HTML file with `Content-Type: image/png` — it lands on disk, and if directly served by Nginx or Go's file server, a victim's browser renders it as HTML, executing any embedded JavaScript.

**Why it happens:** The extension check (`strings.HasSuffix(filename, ".png")`) and Content-Type header check are both trivially spoofed. They feel sufficient during development but are not.

**Consequences:** Stored XSS via uploaded file. If the file server serves with the wrong Content-Type, the browser executes the payload.

**Prevention:**
1. Validate magic bytes server-side using Go's `net/http.DetectContentType`, which reads the first 512 bytes:
   ```go
   buf := make([]byte, 512)
   n, _ := file.Read(buf)
   contentType := http.DetectContentType(buf[:n])
   // Allow only image/jpeg, image/png, image/gif, image/webp
   ```
2. Generate a random filename on the server side (e.g., `uuid + extension`). Discard the original filename. Never use the client-supplied filename as a filesystem path — this prevents path traversal.
3. Enforce a maximum file size via `http.MaxBytesReader` before the multipart parse.
4. Ensure Nginx or the Go file server sets an explicit `Content-Type` header based on the server-assigned extension, not client-supplied metadata.
5. Store images under a path that is not executable (no CGI, no Go handler serving from it with dynamic type detection).

**Detection:** Upload a text file renamed to `evil.png`. If it is served and the browser renders it as HTML, the check is absent.

**Phase:** Address before the image upload feature ships. Do not defer "cleanup later."

---

## Moderate Pitfalls

---

### Pitfall 5: Go `http.Server` With No Timeouts (Slowloris / Resource Exhaustion)

**What goes wrong:** Using `http.ListenAndServe(addr, mux)` or constructing `http.Server{}` without setting `ReadTimeout`, `WriteTimeout`, and `IdleTimeout`. The server holds connections open indefinitely for slow or malicious clients, exhausting goroutines and file descriptors.

**Why it happens:** `http.ListenAndServe` is the canonical "getting started" call. It has no timeout parameters. The convenience is a trap for production use.

**Consequences:** Slowloris-style DoS attacks work. A crawler that opens many slow connections can take the server down without needing a high request rate.

**Prevention:**
```go
srv := &http.Server{
    Addr:         ":8080",
    Handler:      mux,
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  60 * time.Second,
}
srv.ListenAndServe()
```
Note: if any route streams a response (e.g., a future SSE endpoint), `WriteTimeout` must be adjusted or set per-connection via `http.ResponseController`.

**Detection:** `go tool pprof` goroutine profile shows goroutines blocked in `net.(*conn).serve` that are not making progress.

**Phase:** Address when the HTTP server is first constructed. One struct, one fix.

---

### Pitfall 6: `database/sql` Connection Leaks from Unclosed `*sql.Rows`

**What goes wrong:** A handler calls `db.QueryContext(ctx, ...)`, iterates some rows, then returns early on an error without calling `rows.Close()`. The connection is checked out from the pool but never returned. Under load, the pool exhausts and every subsequent query blocks or times out.

**Why it happens:** Go's `database/sql` requires explicit row closing. Forgetting a `defer rows.Close()` immediately after `db.Query` is a classic Go mistake.

**Consequences:** Connection pool exhaustion under any non-trivial load. All database operations hang. Service appears up but returns no data.

**Prevention:**
```go
rows, err := db.QueryContext(ctx, query, args...)
if err != nil {
    return err
}
defer rows.Close() // always, immediately after the nil check

for rows.Next() { ... }
return rows.Err()
```
Also: use `db.ExecContext` (not `db.QueryContext`) for INSERT/UPDATE/DELETE. Using `QueryContext` for non-SELECT statements leaks the result set.

**Detection:** Monitor `db.Stats().OpenConnections` in a `/metrics` or `/healthz` endpoint. A monotonically increasing count under stable load indicates a leak.

**Phase:** Address during the database layer phase. Add a linting rule (e.g., `sqlclosecheck` via `golangci-lint`) to catch this automatically.

---

### Pitfall 7: X-Forwarded-For Header Blindly Trusted for Rate Limiting or Logging

**What goes wrong:** The Go server reads `r.Header.Get("X-Forwarded-For")` to log client IPs or enforce per-IP rate limits. An attacker connecting directly to port 8080 (bypassing ALB) or even through ALB can prepend a spoofed IP to the header.

**Why it happens:** AWS ALB appends the real client IP to `X-Forwarded-For` — it does not replace the header. So a request with `X-Forwarded-For: 1.2.3.4` from a client becomes `X-Forwarded-For: 1.2.3.4, <real-client-ip>`. If the server reads `[0]` it reads the spoofed value.

**Consequences:** Rate-limit bypass. Inaccurate access logs. IP-based allow/deny lists are circumvented.

**Prevention:** Read from the rightmost (last) IP in `X-Forwarded-For`, which is the one appended by the most recently trusted proxy (ALB). For this stack: trust only the rightmost entry added by ALB. If Nginx also adds one, trust the second-to-rightmost. Document the exact proxy chain so the parsing logic matches the deployment topology.

**Detection:** In logs, look for obviously spoofed IPs (e.g., `127.0.0.1`, `10.0.0.1`) appearing as the apparent client IP.

**Phase:** Address when implementing request logging and any rate-limiting middleware.

---

### Pitfall 8: ASG Scale-In Terminates the Instance That Holds the EBS Volume

**What goes wrong:** The ASG is configured with `min 1, max 2`. Under a brief traffic spike it scales to 2. The EBS volume (containing Postgres data and uploaded images) is attached to instance A. ASG terminates instance A on scale-in. Instance B has no EBS. Data is gone.

**Why it happens:** EBS volumes are AZ-bound and instance-bound. They do not follow instances in an ASG automatically. The architecture decision to use EBS for storage is sound for `desired 1` but fragile the moment a second instance exists.

**Consequences:** Complete data loss. Postgres data and uploaded images on the detached volume are unrecoverable unless a snapshot exists.

**Prevention:**
1. Lock the ASG to `max 1` in Terraform to eliminate scale-out. Self-healing (replace unhealthy instance) is the only permitted scaling event. This is documented as the intent in PROJECT.md.
2. Implement a lifecycle hook or `aws:autoscaling:lifecyclehook` that detaches and re-attaches the EBS volume to the replacement instance on scale events.
3. Set the EBS volume to `delete_on_termination = false` in the ASG launch template so the volume survives instance replacement.
4. Take daily EBS snapshots (EventBridge + Lambda or AWS Backup) as a recovery baseline.

**Detection:** Review the ASG launch template in Terraform — confirm `delete_on_termination = false` on the data volume. Confirm `max_size = 1` in the ASG resource.

**Phase:** Address before the first production deployment. This is infrastructure, not application code.

---

### Pitfall 9: Go Project Structure — Everything in `main` or a Single Package

**What goes wrong:** The server starts as one `main.go` file with all handlers, database calls, and business logic inline. It grows to 2000 lines. Testing becomes impossible because everything is tightly coupled to the `http.Request`/`http.ResponseWriter` interface.

**Why it happens:** Go's simplicity encourages starting in `main`. There is no framework enforcing structure. Refactoring later requires touching imports across the whole codebase.

**Consequences:** Untestable code. No clear boundary between HTTP transport layer and application logic. Database queries mixed with template rendering.

**Prevention:** Start with a flat but intentional layout:
```
cmd/server/main.go      — wires dependencies, starts server
internal/handler/       — HTTP handlers only; call into service layer
internal/service/       — business logic; no http.* imports
internal/store/         — database queries; no http.* imports
internal/model/         — shared types / domain structs
internal/template/      — template loading and rendering helpers
static/                 — CSS, JS, images
templates/              — HTML templates
```
The `internal/` boundary prevents future accidental coupling. Keep `main.go` under ~50 lines — it should only call `service.Run()` or equivalent. Avoid generic package names like `utils` or `helpers`.

**Detection:** If a `store` package imports `net/http`, something has gone wrong.

**Phase:** Address at project initialization, before any feature code is written.

---

## Minor Pitfalls

---

### Pitfall 10: bcrypt Cost Factor Left at `bcrypt.DefaultCost` (10)

**What goes wrong:** `bcrypt.DefaultCost` is 10, set in 2011. Modern GPUs (e.g., RTX 5090 in 2025) are approximately 65% faster than the previous generation. Cost 10 produces a hash in ~100ms on server hardware but can be cracked much faster offline.

**Prevention:** Use cost 12 (approximately 400ms per hash on typical server CPU) for a single-admin blog. This is imperceptible to the one person who logs in occasionally, and meaningfully increases offline cracking cost. Use `bcrypt.CompareHashAndPassword` for comparison — it is constant-time by construction, so no additional timing-attack protection is needed beyond using this function.

**Phase:** Set when implementing the admin login handler.

---

### Pitfall 11: `http.ListenAndServe` Called in `main` Without Graceful Shutdown

**What goes wrong:** `os.Interrupt` (Ctrl-C or `docker stop`) kills the process immediately. In-flight requests are dropped mid-response. Any open database transactions are not rolled back cleanly.

**Prevention:** Use `signal.NotifyContext` + `srv.Shutdown(ctx)` with a short deadline (5–10 seconds). This is a 20-line addition.

**Phase:** Address when the HTTP server is first constructed. Pair with the timeout fix (Pitfall 5).

---

### Pitfall 12: Nginx Forwards `Host` Header from ALB, Not Original Client Host

**What goes wrong:** Nginx is configured with `proxy_pass http://localhost:8080` but does not forward `Host` or sets it to `localhost`. The Go server cannot distinguish requests for `jared-wallace.com` from requests for `admin.jared-wallace.com` if it relies on the `Host` header for routing.

**Prevention:** In the Nginx config:
```nginx
proxy_set_header Host              $host;
proxy_set_header X-Real-IP         $remote_addr;
proxy_set_header X-Forwarded-For   $proxy_add_x_forwarded_for;
proxy_set_header X-Forwarded-Proto $scheme;
```
In the Go server, read `r.Host` for subdomain routing, but verify it matches expected values (never trust arbitrary host headers for security decisions).

**Detection:** Make a request and log `r.Host` on the Go side. If it prints `localhost:8080`, the Nginx config needs the `proxy_set_header Host $host` line.

**Phase:** Address when wiring up the Nginx config in the deployment phase.

---

### Pitfall 13: RSS Feed Exposes Draft Posts or Internal Slugs

**What goes wrong:** The RSS handler queries all posts without a `published = true` filter. Draft posts (saved but not published) appear in the feed and are indexed by RSS readers before the author intends.

**Prevention:** Enforce `WHERE published = true AND published_at <= NOW()` on all public-facing queries, including the RSS feed. Add an integration test that seeds a draft and a published post, then asserts the RSS endpoint returns only the published one.

**Phase:** Address during the RSS feature implementation.

---

## CSS Design Polish Pitfalls (v1.2 Shore Leave Polish)

These pitfalls are specific to the v1.2 milestone: adding CSS noise textures, dark mode transitions, page entry animations, and footer redesign to the existing site.

---

### Pitfall 14: CSS Noise Texture via `background-attachment: fixed` Kills Mobile Scrolling

**What goes wrong:**
A noise or grain overlay implemented with `background-attachment: fixed` causes full-page repaints on every scroll frame. On desktop this is merely sluggish. On iOS Safari and Android Chrome it drops to single-digit fps, because those browsers disable fixed-background compositing when the address bar animates during scroll — forcing a CPU repaint instead of a GPU composite. The bug is documented and intentional — browsers will not fix it.

**Why it happens:**
`background-attachment: fixed` tells the browser the background stays in viewport coordinates, which breaks the scroll compositor path. Developers reach for it because it looks correct in DevTools desktop emulation, which does not replicate this rendering path.

**How to avoid:**
Never use `background-attachment: fixed` for the noise overlay. Instead, apply the texture via a `::before` pseudo-element on `body` with `position: fixed; inset: 0; z-index: -1; pointer-events: none`. This creates a single composited layer that costs one paint at load time and zero repaints on scroll:

```css
body::before {
  content: '';
  position: fixed;
  inset: 0;
  z-index: -1;
  pointer-events: none;
  background-image: url("data:image/svg+xml,..."); /* inline SVG feTurbulence */
  opacity: 0.04;
  mix-blend-mode: multiply;
}

[data-theme="dark"] body::before {
  mix-blend-mode: overlay;
  opacity: 0.06;
}
```

Use an inline SVG `feTurbulence` filter as the background source — not a PNG. The SVG compresses to ~900 bytes and requires no additional HTTP request. Keep `opacity` at 0.03–0.06 max.

**Warning signs:**
- Any `background-attachment: fixed` on `body` or a full-bleed container
- Chrome DevTools "Show paint flashing" turns the entire viewport red on scroll

**Phase to address:** Background texture phase. Must be implemented correctly on first pass — retrofitting this after animations are layered in disrupts z-index stacking.

---

### Pitfall 15: Dark Mode Transitions Fire on Page Load and Flash White

**What goes wrong:**
Adding `transition: background-color 300ms ease, color 300ms ease` to `:root` or `body` looks great when the user clicks the dark mode toggle. But on page load, the browser fires the transition from the default light values to the dark values — producing a 300ms white-to-dark flash that is more jarring than no transition at all. This happens even though the blocking script in `base.html` `<head>` sets `data-theme="dark"` before first paint.

**Why it happens:**
CSS transitions do not distinguish "set before paint by a blocking script" from "set by user interaction." Any attribute change triggers registered transitions, including the initial one applied during HTML parsing.

**How to avoid:**
Gate all color transitions behind a `.theme-ready` class that is added by JavaScript only after the page has fully loaded. The class must never be present in the initial HTML:

```css
/* No transitions by default */
:root { --color-bg: #F5F0E8; /* ... */ }

/* Transitions activate only after JS adds .theme-ready */
.theme-ready * {
  transition: background-color 250ms ease, color 250ms ease,
              border-color 250ms ease;
}
```

```js
/* In main.js, after existing dark toggle code */
window.addEventListener('load', function() {
  document.documentElement.classList.add('theme-ready');
});
```

The `load` event fires after all subresources are parsed — the class is never present during the blocking script's initial `data-theme` set.

**Warning signs:**
- Any `transition` on `:root`, `html`, or `body` without a `.theme-ready` guard
- White flash visible when navigating between pages with dark mode active
- DevTools Timeline shows a full-page style recalculation immediately after DOMContentLoaded

**Phase to address:** Dark mode transitions phase. Must be addressed before any other color transitions are added anywhere in the cascade.

---

### Pitfall 16: Page Entry Animations Cause Cumulative Layout Shift (CLS)

**What goes wrong:**
Card stagger animations that animate `margin`, `padding`, `height`, or absolute positioning properties cause CLS and tank Lighthouse scores. Even correct `transform` + `opacity` animations cause CLS if the initial `opacity: 0` state is set via JavaScript after paint — the element is visible for one frame before the animation resets it. `animation-fill-mode: none` (the default) also leaves elements invisible after animation completes.

**Why it happens:**
- Developers set `opacity: 0` in a JS `onload` handler rather than in CSS, which fires after the browser has already painted the element visible
- Stagger code uses CSS custom property delays without `animation-fill-mode: both`, so elements snap back to their initial (invisible) state after animating
- "Card rise" effects animate `box-shadow` or `border` instead of `transform`

**How to avoid:**
Set the start state in CSS, not JavaScript. Use only `transform` and `opacity`. Always use `animation-fill-mode: both`. Gate everything behind `prefers-reduced-motion: no-preference` so the initial state is visible by default:

```css
/* Elements are visible by default — motion is progressive enhancement */
@media (prefers-reduced-motion: no-preference) {
  .post-card {
    opacity: 0;
    animation: card-enter 400ms ease both;
  }
  .post-card:nth-child(2) { animation-delay: 80ms; }
  .post-card:nth-child(3) { animation-delay: 160ms; }
  .post-card:nth-child(4) { animation-delay: 240ms; }
}

@keyframes card-enter {
  from { opacity: 0; transform: translateY(12px); }
  to   { opacity: 1; transform: translateY(0); }
}
```

The `prefers-reduced-motion: no-preference` wrapper means `.post-card` has no `opacity: 0` unless motion is allowed — users with reduced motion see fully visible cards instantly.

**Warning signs:**
- `opacity: 0` in a `<script>` block rather than in CSS
- Animations on `margin-top`, `top`, `height`, or `max-height`
- Missing `animation-fill-mode: both` or `fill-mode: both` on entrance animations
- Lighthouse CLS score degrades after animations are added (target: ≤ 0.1)

**Phase to address:** Page entry animations phase. Run Lighthouse before and after as the gating check.

---

### Pitfall 17: `prefers-reduced-motion` Not Respected — Existing and New Animations

**What goes wrong:**
The existing `.reaction-bounce` animation (in `main.css`, lines 640–645) and all new stagger/fade animations will run for users who have enabled "Reduce Motion" in their OS accessibility settings. This is a WCAG 2.3.3 failure and can trigger vestibular disorder symptoms.

**Why it happens:**
Reduced-motion handling is an afterthought. The existing `reaction-bounce` has no reduced-motion guard. New animations added for v1.2 will repeat this pattern if not deliberately prevented.

**How to avoid:**
Add a single override block at the bottom of `main.css` that covers all motion — both the existing animation and any new ones:

```css
@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
}
```

Do not use `animation: none` — this breaks `animation-fill-mode: both` and leaves elements with `opacity: 0` permanently invisible. Setting duration to `0.01ms` preserves fill-mode while making motion imperceptible.

**Warning signs:**
- Any `@keyframes` block without a corresponding reduced-motion guard
- The existing `reaction-bounce` animation ships without a fix alongside the new animations

**Phase to address:** Earliest animation phase — bundle the `reaction-bounce` fix into the same CSS edit that adds new animations.

---

### Pitfall 18: Footer Expansion Breaks ARIA Landmark Semantics

**What goes wrong:**
Adding a navigation block inside the footer without proper labeling creates ambiguous ARIA landmarks. If the footer gains a `<nav>` element and the primary site nav in the header also has no `aria-label`, screen reader users hear two identical "navigation" landmark announcements and cannot distinguish between them. Additionally, if the footer is structurally moved inside a `<main>` or `<div>` wrapper, it loses its implicit `contentinfo` landmark role.

**Why it happens:**
Footer redesigns focus on visual layout (flex columns, personality block) and miss that the `<footer>` element has an implicit ARIA role that depends on its position in the DOM. Adding unlabeled `<nav>` elements inside it is the most common regression.

**How to avoid:**
Label the footer nav distinctly from the primary nav. The primary `<nav>` in `base.html` currently has no `aria-label` — acceptable when it is the only nav on the page. Once a second `<nav>` is added to the footer, both require distinct labels:

```html
<!-- In base.html: add aria-label to primary nav -->
<nav class="site-nav" aria-label="Primary navigation">
  ...
</nav>

<!-- In the footer -->
<footer class="site-footer">
  <nav aria-label="Footer navigation">
    <a href="/about">About</a>
    <a href="/rss">RSS</a>
  </nav>
  <p>&copy; {{.Year}} Jared Wallace ...</p>
</footer>
```

Keep the `<footer>` as a direct child of `<body>` — do not nest it inside `<main>` or a wrapper `<div>`.

**Warning signs:**
- Two `<nav>` elements in the page with no `aria-label` on either
- `<footer>` moved inside a layout wrapper or `<main>` tag during the redesign
- Axe DevTools reports "Landmark region must have accessible name" after the footer update

**Phase to address:** Footer redesign phase. Run an Axe DevTools scan before and after as the pass/fail criterion.

---

### Pitfall 19: `mix-blend-mode` on Noise Overlay Inverts in Dark Mode

**What goes wrong:**
A noise overlay using `mix-blend-mode: multiply` looks correct on the light sandy `#F5F0E8` background. In dark mode (`#1A1F2E`), `multiply` absorbs all remaining light and makes the background pitch-black or introduces a muddy color cast. The reverse — `mix-blend-mode: screen` — washes out the light mode background.

**Why it happens:**
Developers prototype in one color mode, ship it, and discover the other mode is broken after deploy.

**How to avoid:**
Use separate blend modes per theme. Test both modes before committing:

```css
body::before {
  mix-blend-mode: multiply; /* correct for light sandy bg */
  opacity: 0.04;
}

[data-theme="dark"] body::before {
  mix-blend-mode: overlay;  /* or soft-light for dark backgrounds */
  opacity: 0.06;
}
```

Correct opacity range: 0.03–0.05 light mode, 0.05–0.08 dark mode. Above 0.08 reads as grime.

**Warning signs:**
- Single `mix-blend-mode` with no dark-mode override
- Dark mode body background appears darker than `--color-bg: #1A1F2E` should produce
- Light mode fine, dark mode screenshots show a solid dark vignette over content

**Phase to address:** Background texture phase — treat dark mode verification as a required step, not a post-launch check.

---

## Phase-Specific Warnings

| Phase Topic | Likely Pitfall | Mitigation |
|-------------|----------------|------------|
| HTTP server bootstrap | No timeouts, no graceful shutdown (Pitfalls 5, 11) | Set `ReadTimeout`, `WriteTimeout`, `IdleTimeout`; add shutdown handler |
| Project structure init | Everything in `main`, no `internal/` (Pitfall 9) | Lay out `cmd/`, `internal/handler/`, `internal/service/`, `internal/store/` before any feature code |
| Admin session auth | Missing cookie flags, session fixation, CSRF (Pitfall 1) | `HttpOnly + Secure + SameSite`, regenerate session ID on login, CSRF token on forms |
| Admin login handler | Low bcrypt cost (Pitfall 10) | `bcrypt.GenerateFromPassword(pw, 12)` |
| Markdown rendering | XSS via `template.HTML` (Pitfall 2) | goldmark → bluemonday → `template.HTML`, write a test |
| Database layer | Unclosed rows, pool exhaustion (Pitfall 6) | `defer rows.Close()` pattern; `sqlclosecheck` linter |
| Image upload | Client MIME spoofing, path traversal (Pitfall 4) | Magic bytes via `http.DetectContentType`, random server-side filename |
| Docker Compose wiring | Postgres data dir permissions (Pitfall 3) | `chown 999:999` on host dir before first run; document in Makefile |
| Nginx config | Host header not forwarded (Pitfall 12) | `proxy_set_header Host $host` in every server block |
| Deployment | EBS data loss on ASG scale event (Pitfall 8) | `max_size = 1`, `delete_on_termination = false`, EBS snapshot schedule |
| Request logging / middleware | X-Forwarded-For spoofing (Pitfall 7) | Read rightmost IP, document trusted proxy chain |
| RSS feed | Draft post exposure (Pitfall 13) | Filter by `published = true` on all public queries |
| Background texture (v1.2) | Mobile scroll lag from `background-attachment: fixed` (Pitfall 14) | Use `position: fixed` pseudo-element; verify on a physical mobile device |
| Dark mode transitions (v1.2) | Load flash from transition firing on theme init (Pitfall 15) | `.theme-ready` class guard; `window.addEventListener('load', ...)` in `main.js` |
| Page entry animations (v1.2) | CLS regression and invisible elements after animation (Pitfall 16) | CSS-only initial state; `animation-fill-mode: both`; Lighthouse CLS ≤ 0.1 |
| Any animation (v1.2) | `prefers-reduced-motion` not respected (Pitfall 17) | Add nuclear reduced-motion override block; fix existing `reaction-bounce` in same pass |
| Footer redesign (v1.2) | Duplicate unlabeled nav landmarks, contentinfo lost (Pitfall 18) | `aria-label` on both `<nav>` elements; footer stays direct child of `<body>` |
| Noise texture dark mode (v1.2) | `mix-blend-mode: multiply` inverts dark bg (Pitfall 19) | Separate blend mode override for `[data-theme="dark"] body::before` |

---

## "Looks Done But Isn't" Checklist (v1.2 CSS Polish)

- [ ] **Noise texture on mobile:** Scroll test on a physical iOS or Android device — no visible lag, no full-page paint flash in Chrome DevTools "Show paint flashing"
- [ ] **Dark mode transition — no load flash:** Open a fresh browser tab with dark mode active. Zero visible white flash before the page is interactive
- [ ] **Reduced motion:** Enable OS "Reduce Motion" setting — all animations are imperceptible, no elements are left in `opacity: 0` state permanently
- [ ] **CLS after animations:** Lighthouse CLS score ≤ 0.1 before and after adding card stagger
- [ ] **Footer ARIA:** Axe DevTools reports zero landmark or label violations after footer expansion
- [ ] **RSS icon location:** RSS icon remains in the footer's top-level visible area regardless of expanded layout — not hidden inside a collapsed section
- [ ] **Copyright year:** `{{.Year}}` template variable still renders correctly after footer HTML restructure
- [ ] **Noise texture blend mode:** Both light and dark mode screenshots show texture visible but not dominant — background color reads as intended in both modes
- [ ] **Reaction button radius:** Visual check that `.reacted` and `.bounce` state buttons also display the new `4px` radius (they inherit, but verify in browser)
- [ ] **CSS comment rebrand:** `grep -r "The Log" web/static/` returns zero results

---

## Sources

- [OWASP Session Management Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html)
- [Building a Secure Session Manager in Go](https://themsaid.com/building-secure-session-manager-in-go)
- [OWASP File Upload Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/File_Upload_Cheat_Sheet.html)
- [bluemonday - Go Packages](https://pkg.go.dev/github.com/microcosm-cc/bluemonday)
- [The complete guide to Go net/http timeouts - Cloudflare](https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/)
- [Postgres data volume wrong ownership - Docker forums](https://forums.docker.com/t/data-directory-var-lib-postgresql-data-pgdata-has-wrong-ownership/17963)
- [ASG with stateful Docker containers - Portworx](https://portworx.com/blog/auto-scaling-groups-ebs-docker/)
- [CSS-Tricks: Grainy Gradients](https://css-tricks.com/grainy-gradients/) — SVG feTurbulence technique, browser inconsistencies
- [Frontend Masters: Grainy Gradients](https://frontendmasters.com/blog/grainy-gradients/) — implementation approach
- [freeCodeCamp: Grainy CSS Backgrounds using SVG Filters](https://www.freecodecamp.org/news/grainy-css-backgrounds-using-svg-filters/) — SVG vs PNG weight comparison
- [CSS-Tricks: Fixed Background Attachment Hack](https://css-tricks.com/the-fixed-background-attachment-hack/) — pseudo-element workaround for mobile performance
- [Mozilla Bugzilla #90198](https://bugzilla.mozilla.org/show_bug.cgi?id=90198) — fixed-background repaints on scroll (intentional browser behavior)
- [web.dev: Optimize CLS](https://web.dev/articles/optimize-cls) — safe vs unsafe animation properties, transform-only rule
- [Motion.dev: Web Animation Performance Tier List](https://motion.dev/magazine/web-animation-performance-tier-list) — compositor-only properties (transform, opacity, filter)
- [dev.to: Light/dark mode avoid flickering on reload](https://dev.to/ayc0/light-dark-mode-avoid-flickering-on-reload-1567) — blocking script pattern for dark mode
- [Blog of Maxime Heckel: Fixing dark mode flash on server-rendered sites](https://blog.maximeheckel.com/posts/switching-off-the-lights-part-2-fixing-dark-mode-flashing-on-servered-rendered-website/) — `.theme-ready` deferred class pattern
- [MDN: prefers-reduced-motion](https://developer.mozilla.org/en-US/docs/Web/CSS/Reference/At-rules/@media/prefers-reduced-motion) — media query specification
- [W3C WAI: C39 — Using prefers-reduced-motion](https://www.w3.org/WAI/WCAG22/Techniques/css/C39.html) — WCAG technique
- [MDN: ARIA navigation role](https://developer.mozilla.org/en-US/docs/Web/Accessibility/ARIA/Reference/Roles/navigation_role) — landmark labeling requirements when multiple nav elements exist
- [MDN: ARIA contentinfo role](https://developer.mozilla.org/en-US/docs/Web/Accessibility/ARIA/Reference/Roles/contentinfo_role) — footer landmark rules, position requirements
- [W3C WAI: Landmark Regions](https://www.w3.org/WAI/ARIA/apg/practices/landmark-regions/) — multiple nav labeling requirement
