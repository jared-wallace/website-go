# Phase 3: Admin Panel - Research

**Researched:** 2026-03-26
**Domain:** Go HTTP session auth, host-based routing, admin CRUD, split-pane editor, CSRF protection
**Confidence:** HIGH

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

**Subdomain Routing**
- D-01: Host-based mux in the single binary. Admin handlers serve requests where `Host` matches `admin.jared-wallace.com`; all other hosts route to blog handlers. Nginx already forwards both domains to :8080.
- D-02: Shared nautical design — admin reuses the beach bar aesthetic (base template, color palette, typography). Admin feels like the "back office" of the same beach bar.
- D-03: Unauthenticated requests to admin subdomain see a branded nautical-themed login page (not a redirect to the public blog).

**Editor Experience**
- D-04: Vanilla JS fetch for live preview. Editor textarea sends markdown to a server endpoint via debounced fetch() on keyup. Server renders with the same goldmark+bluemonday pipeline and returns HTML. Preview is always identical to published output.
- D-05: Manual save only. Explicit "Save Draft" and "Publish" buttons. Ctrl+S / Cmd+S keyboard shortcut. No autosave, no localStorage draft backup.
- D-06: On mobile/narrow screens, editor and preview stack vertically with a Write/Preview tab toggle. Full-width editing.
- D-07: Plain textarea with monospace font. No toolbar buttons (bold, italic, link, etc.). Admin knows markdown.

**Auth & Session Flow**
- D-08: Admin credentials stored as environment variables: `ADMIN_EMAIL` and `ADMIN_PASSWORD_HASH`. No users/admins table. Hash generated offline via `make hash-password` Makefile target. Single admin user per requirements.
- D-09: SCS v2.9.0 with pgxstore for Postgres-backed sessions. Session cookie carries HttpOnly + Secure + SameSite=Lax flags per success criteria.
- D-10: 24-hour session lifetime (inactivity-based expiry). Re-login required after 24 hours.
- D-11: Failed login shows generic "Invalid email or password" message (never reveals which is wrong). In-memory rate limiter: 5 attempts per minute per IP to prevent brute force.

**Admin Dashboard**
- D-12: Post list as a table (title, status, date, actions) with filter tabs: All / Published / Drafts / Deleted. Action links for Edit, Publish/Unpublish, Delete.
- D-13: Soft-delete recovery via "Restore" button in the Deleted tab. Restores post to draft status. No permanent delete in v1.
- D-14: Slug auto-generates from title (e.g., "My Post" -> "my-post") but admin can manually edit before saving.
- D-15: Delete action requires confirmation dialog. Publish/unpublish act immediately (easily reversible).

### Claude's Discretion
- CSRF implementation approach (check Go 1.26 stdlib CrossOriginProtection first, fall back to gorilla/csrf or manual tokens)
- Session middleware design and request context integration
- Admin template structure (separate admin base template extending shared styles, or full reuse of blog base.html)
- Table styling and responsive behavior on the dashboard
- Login form layout and error display
- Slug generation algorithm (Unicode handling, special character stripping)
- Debounce timing for editor preview (200-500ms range)
- Rate limiter implementation details (sync.Map, sliding window, etc.)

### Deferred Ideas (OUT OF SCOPE)
None — discussion stayed within phase scope.
</user_constraints>

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| ADMIN-01 | Admin can log in with email/password at admin.jared-wallace.com | D-08: bcrypt verify; D-09: SCS session; host-based routing dispatches to admin login page |
| ADMIN-02 | Admin session persists across browser refresh (Postgres-backed) | D-09: pgxstore stores session token in DB; survives process restart; SCS LoadAndSave middleware |
| ADMIN-03 | Admin can create posts with title, markdown body, and slug | New repository write methods (Create/Insert); service layer slug generation; goldmark pre-render on write |
| ADMIN-04 | Admin can edit existing posts | New repository UpdateByID; editor template re-uses same form with populated fields |
| ADMIN-05 | Admin can soft-delete posts (recoverable) | Model already has DeletedAt *time.Time; SoftDelete sets it; Restore clears it to NULL |
| ADMIN-06 | Admin can toggle posts between draft and published states | Published bool already on model; SetPublished(id, bool) repository method |
| ADMIN-07 | Admin can write in split-pane markdown editor with live preview | POST /admin/preview endpoint returns rendered HTML; debounced fetch() in vanilla JS |
</phase_requirements>

---

## Summary

Phase 3 builds entirely on the existing foundation. The Go binary already handles a single ServeMux; this phase adds a host-discriminating wrapper that routes `admin.jared-wallace.com` traffic to a new AdminHandler while all other hosts continue to the existing BlogHandler. The session layer (SCS v2.9.0 + pgxstore) and bcrypt dependency are already declared in go.mod as indirect imports — this phase promotes them to direct use.

The CSRF question (noted as a blocker in STATE.md) is resolved: `net/http.CrossOriginProtection` was introduced in Go 1.25 and is present in the local Go 1.26.0 runtime. It uses `Sec-Fetch-Site` and `Origin` headers — the same idiom as SameSite cookies — and requires zero new dependencies. For an admin panel with a single trusted origin (`https://admin.jared-wallace.com`), one `AddTrustedOrigin` call is sufficient.

The most complex part of this phase is the split-pane editor (ADMIN-07): a server-side preview endpoint reuses the existing `markdown.Renderer.Render()` function; the frontend is vanilla JS with a debounced `fetch()` call. No new JS frameworks needed.

**Primary recommendation:** Wire SCS LoadAndSave middleware globally, add host-based dispatch in main.go, extend the repository interface with write methods, and build the admin handler following the exact same constructor/template pattern as BlogHandler.

---

## Standard Stack

### Core (all already in go.mod)

| Library | Version | Purpose | Status |
|---------|---------|---------|--------|
| `github.com/alexedwards/scs/v2` | v2.9.0 | HTTP session management | indirect — promote to direct |
| `github.com/alexedwards/scs/pgxstore` | v0.0.0-20251002... | Postgres session backend | indirect — promote to direct |
| `golang.org/x/crypto/bcrypt` | v0.49.0 | Password verification | indirect — promote to direct |
| `net/http` CrossOriginProtection | Go 1.25+ (stdlib) | CSRF protection | zero dep, confirmed in Go 1.26.0 |

### No New External Dependencies Required

All libraries for Phase 3 are already present in go.mod. Promoting indirect dependencies to direct imports is done by adding `import` statements; `go mod tidy` updates the require block.

**Promote from indirect to direct:**
```bash
# These are already downloaded — no go get needed.
# Adding import paths in code + running go mod tidy is sufficient.
# To be explicit:
go get github.com/alexedwards/scs/v2@v2.9.0
go get github.com/alexedwards/scs/pgxstore
go get golang.org/x/crypto@v0.49.0
```

---

## Architecture Patterns

### Recommended Project Structure (additions for Phase 3)

```
internal/
├── config/
│   └── config.go              # ADD: AdminEmail, AdminPasswordHash, SessionSecret
├── handler/
│   ├── blog/                  # EXISTING — no changes
│   └── admin/                 # NEW
│       ├── handler.go         # Constructor, template parsing, render()
│       ├── auth.go            # LoginPage, LoginPost, Logout
│       ├── dashboard.go       # Dashboard (post list with filter tabs)
│       ├── editor.go          # NewPost, EditPost, SaveDraft, Publish
│       ├── actions.go         # SoftDelete, Restore, TogglePublished
│       └── preview.go         # POST /admin/preview endpoint
├── middleware/                 # NEW
│   └── auth.go                # RequireSession middleware
├── repository/
│   └── post/
│       ├── repository.go      # EXTEND: add write methods to interface
│       ├── queries.go         # EXISTING read queries — no changes
│       └── write.go           # NEW: Create, Update, SoftDelete, Restore, GetByID, SetPublished, ListAll
└── service/
    └── post/
        ├── service.go         # EXISTING
        ├── write.go           # NEW: Create, Update, SoftDelete, Restore, Publish, Unpublish, ListAll
        └── slug.go            # NEW: GenerateSlug function

db/migrations/
└── 00003_create_sessions.sql  # NEW: sessions table for pgxstore

web/
├── templates/
│   ├── admin-base.html        # NEW: admin base template (extends shared CSS/fonts)
│   ├── admin-login.html       # NEW
│   ├── admin-dashboard.html   # NEW
│   └── admin-editor.html      # NEW
└── static/
    ├── admin.css              # NEW: admin-specific styles (editor layout, table)
    └── admin.js               # NEW: debounced preview fetch, Ctrl+S handler, slug generation
```

### Pattern 1: Host-Based Dispatch in main.go

Go's `http.ServeMux` does not natively dispatch on `Host` header. The standard pattern is a top-level handler that inspects `r.Host` and delegates to the appropriate mux.

```go
// Source: verified pattern from Go stdlib net/http docs + project ARCHITECTURE.md
type hostRouter struct {
    blog  http.Handler
    admin http.Handler
}

func (hr *hostRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    host := r.Host
    // Strip port suffix for local dev (localhost:8080)
    if h, _, err := net.SplitHostPort(host); err == nil {
        host = h
    }
    if host == "admin.jared-wallace.com" || host == "localhost" {
        hr.admin.ServeHTTP(w, r)
        return
    }
    hr.blog.ServeHTTP(w, r)
}
```

Wire in main.go:
```go
// Wrap admin mux with CrossOriginProtection
cop := http.NewCrossOriginProtection()
cop.AddTrustedOrigin("https://admin.jared-wallace.com")

router := &hostRouter{
    blog:  blogMux,
    admin: sessionManager.LoadAndSave(cop.Handler(adminMux)),
}
srv := server.New(cfg.Port, router)
```

### Pattern 2: SCS Session Wiring

```go
// Source: pkg.go.dev/github.com/alexedwards/scs/v2 (verified)
sessionManager := scs.New()
sessionManager.Store = pgxstore.New(pool)
sessionManager.Lifetime = 24 * time.Hour
sessionManager.Cookie.HttpOnly = true
sessionManager.Cookie.Secure = true
sessionManager.Cookie.SameSite = http.SameSiteLaxMode
sessionManager.Cookie.Name = "admin_session"
```

`LoadAndSave` middleware must wrap the entire admin mux. It loads the session from the store on every request and saves changes before the response is written.

### Pattern 3: Login Handler with Session Fixation Prevention

```go
// Source: scs v2 official docs - RenewToken prevents session fixation
func (h *AdminHandler) LoginPost(w http.ResponseWriter, r *http.Request) {
    email := r.FormValue("email")
    password := r.FormValue("password")

    if !h.rateLimiter.Allow(r.RemoteAddr) {
        // return 429 or generic error — don't reveal rate limit specifics
        h.renderLogin(w, http.StatusTooManyRequests, "Too many attempts. Try again later.")
        return
    }

    // Constant-time comparison via bcrypt; always compare even on email mismatch
    if !h.verifyCredentials(email, password) {
        h.renderLogin(w, http.StatusUnauthorized, "Invalid email or password.")
        return
    }

    // RenewToken before privilege change — prevents session fixation
    if err := h.sessions.RenewToken(r.Context()); err != nil {
        h.renderError(w, err)
        return
    }
    h.sessions.Put(r.Context(), "authenticated", true)
    http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
}
```

### Pattern 4: RequireSession Middleware

```go
// Source: scs v2 documentation pattern
func RequireSession(sm *scs.SessionManager, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !sm.GetBool(r.Context(), "authenticated") {
            http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

Apply to all admin routes except `/admin/login`:
```go
adminMux := http.NewServeMux()
adminMux.HandleFunc("GET /admin/login", admin.LoginPage)
adminMux.HandleFunc("POST /admin/login", admin.LoginPost)
adminMux.Handle("GET /admin/posts", RequireSession(sm, http.HandlerFunc(admin.Dashboard)))
adminMux.Handle("GET /admin/posts/new", RequireSession(sm, http.HandlerFunc(admin.NewPost)))
// etc.
```

### Pattern 5: Sessions Table Migration

```sql
-- db/migrations/00003_create_sessions.sql
-- +goose Up
CREATE TABLE sessions (
    token TEXT PRIMARY KEY,
    data  BYTEA       NOT NULL,
    expiry TIMESTAMPTZ NOT NULL
);
CREATE INDEX sessions_expiry_idx ON sessions (expiry);

-- +goose Down
DROP TABLE IF EXISTS sessions;
```

Source: pgxstore official docs (pkg.go.dev/github.com/alexedwards/scs/pgxstore) — this is the exact schema pgxstore expects.

### Pattern 6: Admin Handler Constructor (mirrors BlogHandler)

```go
// Follows internal/handler/blog/handler.go pattern exactly
type AdminHandler struct {
    svc         *postservice.Service
    sessions    *scs.SessionManager
    renderer    *markdown.Renderer
    rateLimiter *RateLimiter
    funcMap     template.FuncMap
    templates   map[string]*template.Template
    adminEmail  string
    adminHash   []byte // bcrypt hash loaded once at startup
}

func New(svc *postservice.Service, sm *scs.SessionManager, r *markdown.Renderer, cfg config.Config) *AdminHandler {
    // Parse templates once at startup — same pattern as BlogHandler
    pages := []string{"admin-login.html", "admin-dashboard.html", "admin-editor.html"}
    templates := make(map[string]*template.Template, len(pages))
    for _, page := range pages {
        tmpl := template.Must(
            template.New("").Funcs(funcMap).ParseFS(
                web.Templates,
                "templates/admin-base.html",
                "templates/"+page,
            ),
        )
        templates[page] = tmpl
    }
    // ...
}
```

### Pattern 7: Repository Write Methods

Extend the existing `Repository` interface with write operations:

```go
type Repository interface {
    // Existing read methods (unchanged)
    ListPublished(ctx context.Context, limit, offset int) ([]model.Post, error)
    CountPublished(ctx context.Context) (int, error)
    FindBySlug(ctx context.Context, slug string) (*model.Post, error)

    // NEW: write methods for admin
    FindByID(ctx context.Context, id int64) (*model.Post, error)
    ListAll(ctx context.Context) ([]model.Post, error)  // includes drafts + deleted
    Create(ctx context.Context, p model.Post) (*model.Post, error)
    Update(ctx context.Context, p model.Post) error
    SoftDelete(ctx context.Context, id int64) error
    Restore(ctx context.Context, id int64) error       // clears deleted_at, sets published=false
    SetPublished(ctx context.Context, id int64, published bool) error
}
```

### Pattern 8: Slug Generation

Pure function, no external dependency needed:

```go
// internal/service/post/slug.go
func GenerateSlug(title string) string {
    // 1. Normalize Unicode to ASCII (golang.org/x/text/unicode/norm — already indirect)
    // 2. Lowercase
    // 3. Replace non-alphanumeric runs with "-"
    // 4. Trim leading/trailing "-"
    // 5. Deduplicate "--"
    s := strings.ToLower(title)
    reg := regexp.MustCompile(`[^a-z0-9]+`)
    s = reg.ReplaceAllString(s, "-")
    s = strings.Trim(s, "-")
    if s == "" {
        return fmt.Sprintf("post-%d", time.Now().Unix())
    }
    return s
}
```

For Unicode titles (accented characters), `golang.org/x/text/transform` + `norm.NFD` can be used — x/text is already an indirect dependency via x/net.

### Pattern 9: In-Memory Rate Limiter (5 attempts/minute per IP)

Simple fixed-window counter using `sync.Map`. Sufficient for single-instance; no Redis needed.

```go
type RateLimiter struct {
    mu      sync.Mutex
    entries sync.Map // key: IP string, value: *entry
}

type entry struct {
    count     int
    windowEnd time.Time
}

// Allow returns true if the IP is under the rate limit.
func (rl *RateLimiter) Allow(ip string, limit int, window time.Duration) bool {
    now := time.Now()
    val, _ := rl.entries.LoadOrStore(ip, &entry{windowEnd: now.Add(window)})
    e := val.(*entry)
    rl.mu.Lock()
    defer rl.mu.Unlock()
    if now.After(e.windowEnd) {
        e.count = 0
        e.windowEnd = now.Add(window)
    }
    if e.count >= limit {
        return false
    }
    e.count++
    return true
}
```

Extract real IP from `r.RemoteAddr`; strip port. For production behind Nginx, trust `X-Forwarded-For` or `X-Real-IP` header (Nginx is already configured to proxy to :8080).

### Pattern 10: CSRF via CrossOriginProtection

Go 1.25+ stdlib, confirmed available in Go 1.26.0 locally:

```go
// Source: pkg.go.dev/net/http#CrossOriginProtection (Go 1.26.0)
cop := http.NewCrossOriginProtection()
// Allow the admin subdomain to make cross-origin POST requests
// (needed for fetch() from admin.jared-wallace.com to same server)
cop.AddTrustedOrigin("https://admin.jared-wallace.com")
// Wrap admin mux only — blog mux is read-only (GET), CrossOriginProtection
// allows all safe methods unconditionally.
adminHandler = cop.Handler(adminMux)
```

Key behavior: requests without `Sec-Fetch-Site` or `Origin` headers (curl, server-to-server) pass through — non-browser API clients are unaffected. Safe methods (GET/HEAD/OPTIONS) are always allowed.

### Pattern 11: Live Preview Endpoint

```go
// POST /admin/preview — returns rendered HTML fragment, not a full page
func (h *AdminHandler) Preview(w http.ResponseWriter, r *http.Request) {
    body := r.FormValue("body")
    html := h.renderer.Render(body)
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprint(w, string(html))
}
```

Frontend debounce (admin.js):
```javascript
let previewTimer;
document.getElementById('editor-body').addEventListener('input', function() {
    clearTimeout(previewTimer);
    previewTimer = setTimeout(() => {
        const formData = new FormData();
        formData.append('body', this.value);
        fetch('/admin/preview', { method: 'POST', body: formData })
            .then(r => r.text())
            .then(html => {
                document.getElementById('preview-pane').innerHTML = html;
            });
    }, 300); // 300ms debounce — within the 200-500ms range from D-04
});
```

### Pattern 12: Admin Template Structure

Admin pages use a separate `admin-base.html` that shares CSS/fonts with `base.html` but has a different nav (shows "Admin Panel", logout link) and no public footer. This avoids cluttering the public `base.html` with admin-only markup.

```html
{{define "base"}}<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>{{block "title" .}}Admin{{end}} — The Log</title>
  <!-- Same fonts and CSS as public base -->
  <link href="https://fonts.googleapis.com/css2?family=Playfair+Display..." rel="stylesheet">
  <link rel="stylesheet" href="/static/main.css">
  <link rel="stylesheet" href="/static/admin.css">
</head>
<body class="admin-body">
  <nav class="admin-nav">
    <a href="/admin/posts" class="site-name">The Log — Back Office</a>
    <a href="/admin/logout" class="logout-link">Go Ashore</a>
  </nav>
  {{block "content" .}}{{end}}
  <script src="/static/admin.js" defer></script>
</body>
</html>
{{end}}
```

### Anti-Patterns to Avoid

- **Storing password hash in a database table:** D-08 locks this to env vars. Never create an admins table.
- **Using gorilla/csrf:** CrossOriginProtection is available in stdlib Go 1.25+. No external CSRF library needed.
- **SCS LoadAndSave wrapping only some routes:** LoadAndSave must wrap the entire admin handler tree so the session is available to RequireSession middleware on every request, including the login page (to detect already-authenticated users).
- **Returning different errors for wrong email vs wrong password:** D-11 requires a single generic message. Constant-time bcrypt.CompareHashAndPassword must run even when email doesn't match (to prevent timing attacks that reveal valid email addresses).
- **Rendering markdown on every GET /posts/{slug}:** The existing architecture pre-renders at write time. Admin edit should re-render and store `rendered_html` when saving, not on every read.
- **Not calling RenewToken after login:** Session fixation vulnerability. Always call `sessionManager.RenewToken(ctx)` before `sessionManager.Put(ctx, "authenticated", true)`.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Session token generation/storage | Custom cookie+DB scheme | SCS v2 + pgxstore | Server-side tokens, OWASP-aligned, automatic cleanup goroutine |
| Password hashing | md5/sha256/custom | bcrypt (golang.org/x/crypto/bcrypt) | bcrypt has adaptive cost factor; never use fast hashes for passwords |
| CSRF protection | Manual token forms | `net/http.CrossOriginProtection` | Stdlib since Go 1.25; origin-header based; zero dependencies |
| Markdown rendering | Custom parser | goldmark + bluemonday (already wired) | `internal/markdown.Renderer.Render()` already exists and is tested |
| Slug uniqueness check | Timestamp suffix | DB UNIQUE constraint on `slug` column | Already exists in schema; let the DB enforce it, handle the error gracefully |

**Key insight:** All the hard security primitives (sessions, password hashing, CSRF) are already decided and mostly implemented — this phase is primarily about adding admin routes, templates, and write-path SQL.

---

## Common Pitfalls

### Pitfall 1: SCS Lifetime vs IdleTimeout Confusion
**What goes wrong:** Setting `Lifetime` and expecting inactivity reset — it does not. `Lifetime` is absolute; `IdleTimeout` resets on activity.
**Why it happens:** The two fields have different semantics. D-10 says "24-hour inactivity-based expiry."
**How to avoid:** Set `sessionManager.IdleTimeout = 24 * time.Hour` (not `Lifetime`) to get inactivity-based expiry. Set `Lifetime = 0` (no absolute limit) or a long value like 30 days as a backstop.
**Warning signs:** Users getting logged out exactly 24 hours after first login regardless of activity.

### Pitfall 2: pgxstore Sessions Table Not Created Before SCS Starts
**What goes wrong:** Server starts, pgxstore.New(pool) succeeds, but first session write fails with "relation 'sessions' does not exist."
**Why it happens:** pgxstore does not auto-create the table.
**How to avoid:** The `00003_create_sessions.sql` migration must exist and run before the server accepts requests. Goose runs at startup via `database.RunMigrations()` in main.go — migration file order is critical.
**Warning signs:** 500 errors on first POST /admin/login.

### Pitfall 3: SCS Cookie SameSite=Strict Breaks Login Redirects
**What goes wrong:** After redirect from POST /admin/login to /admin/posts, session cookie is not sent because SameSite=Strict blocks cross-site navigations — including same-origin redirects initiated from form POST.
**Why it happens:** SameSite=Strict suppresses the cookie on top-level cross-site navigations (e.g., a bookmark link from another domain). More subtly, it can break redirect chains in some browsers.
**How to avoid:** Use `http.SameSiteLaxMode` (D-09 explicitly requires Lax, not Strict). Lax allows the cookie on top-level GET navigations (including redirect targets) but blocks cross-site POST.
**Warning signs:** Infinite redirect loop on login.

### Pitfall 4: Host-Based Routing and Localhost Dev
**What goes wrong:** `r.Host` is `localhost:8080` in local dev, not `admin.jared-wallace.com`. Admin routes become unreachable.
**Why it happens:** Go includes the port in `r.Host` when no Host header strips it.
**How to avoid:** Strip the port using `net.SplitHostPort` before comparing. In local dev, you can either: (a) add a `/etc/hosts` alias and `ADMIN_HOST=localhost` env var, or (b) use a URL path prefix (`/admin/...`) for local dev and rely on host-based routing only in production. The recommended approach is stripping the port and matching on a configurable `ADMIN_HOST` env var defaulting to `admin.jared-wallace.com`.
**Warning signs:** 404 on all admin routes in local dev.

### Pitfall 5: Double-Encoding of Admin Preview HTML
**What goes wrong:** Preview pane shows literal `&lt;h1&gt;` instead of rendered heading.
**Why it happens:** The preview endpoint returns `template.HTML` but the JavaScript sets `innerHTML` — if the server also wraps the response in an HTML template that escapes it, double-encoding occurs.
**How to avoid:** The `/admin/preview` endpoint must write the raw HTML string directly (no template), with `Content-Type: text/html`. The JS sets `element.innerHTML = html` (raw text response). Never pass the preview HTML through a Go template's `{{.}}` — use `{{.HTML}}` with `template.HTML` type or bypass templates entirely for this endpoint.
**Warning signs:** Raw HTML tags visible in preview pane.

### Pitfall 6: bcrypt Timing Attack on Email Check
**What goes wrong:** Login returns immediately when email is wrong (fast string compare), but takes ~100ms when email is right but password is wrong (bcrypt). Attacker can enumerate valid emails.
**Why it happens:** Early-return on email mismatch skips the expensive bcrypt comparison.
**How to avoid:** Always run `bcrypt.CompareHashAndPassword` regardless of whether the email matched. Use a dummy hash for the "email not found" branch:
```go
// Run bcrypt even when email is wrong (constant-time behavior)
hashToCheck := h.adminHash
if email != h.adminEmail {
    hashToCheck = dummyHash // bcrypt hash of a random string, computed once at startup
}
bcrypt.CompareHashAndPassword(hashToCheck, []byte(password))
```
**Warning signs:** Login endpoint responds measurably faster for unknown emails.

### Pitfall 7: Slug Collision on Create
**What goes wrong:** Two posts with similar titles generate the same slug; INSERT fails with UNIQUE constraint violation.
**Why it happens:** Slug generation is deterministic; "My Post" and "my post" both produce "my-post".
**How to avoid:** Catch the pgx unique violation error (`pgconn.PgError` with Code "23505") and return a validation error to the admin: "Slug already in use — please edit the slug field." The admin can manually modify the slug in the editor form.
**Warning signs:** 500 error on save instead of a friendly form validation message.

### Pitfall 8: Flash Messages Require PopString, Not GetString
**What goes wrong:** "Post saved!" flash message persists across multiple page loads.
**Why it happens:** `GetString` reads the value without removing it; the flash message persists in the session.
**How to avoid:** Use `sessionManager.PopString(ctx, "flash")` for one-time messages. `Pop*` methods retrieve and delete in one atomic operation.
**Warning signs:** Flash messages appearing on every dashboard page load after an action.

---

## Code Examples

### Config Extension
```go
// internal/config/config.go
type Config struct {
    DatabaseURL       string
    Port              string
    AppEnv            string
    AdminEmail        string // ADMIN_EMAIL env var
    AdminPasswordHash string // ADMIN_PASSWORD_HASH env var (bcrypt hash)
    SessionSecret     string // SESSION_SECRET env var (32+ byte random string)
    AdminHost         string // ADMIN_HOST env var, default "admin.jared-wallace.com"
}
```

### Makefile hash-password Target
```makefile
## hash-password: generate a bcrypt hash for ADMIN_PASSWORD_HASH (usage: make hash-password PW=yourpassword)
hash-password:
	@go run ./cmd/hashpw "$(PW)"
```
`cmd/hashpw/main.go` is a one-file helper that calls `bcrypt.GenerateFromPassword([]byte(pw), 12)` and prints the result.

### bcrypt Verification (startup loading)
```go
// In AdminHandler.New():
hash, err := base64.StdEncoding.DecodeString(cfg.AdminPasswordHash)
// or simply store as string; bcrypt.CompareHashAndPassword accepts []byte
adminHash := []byte(cfg.AdminPasswordHash)
```

### SQL: ListAll (all posts including drafts and deleted)
```go
const qListAll = `
    SELECT id, title, slug, body, rendered_html, tags, published, created_at, updated_at, deleted_at
    FROM posts
    ORDER BY created_at DESC`
```

### SQL: Create Post
```go
const qCreate = `
    INSERT INTO posts (title, slug, body, rendered_html, tags, published)
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING id, created_at, updated_at`
```

### SQL: SoftDelete / Restore
```go
const qSoftDelete = `UPDATE posts SET deleted_at = now(), updated_at = now() WHERE id = $1`
const qRestore    = `UPDATE posts SET deleted_at = NULL, published = false, updated_at = now() WHERE id = $1`
```

### SQL: SetPublished
```go
const qSetPublished = `UPDATE posts SET published = $1, updated_at = now() WHERE id = $2 AND deleted_at IS NULL`
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| gorilla/csrf token forms | `net/http.CrossOriginProtection` | Go 1.25 (2025) | Zero dependency; origin-header based; no hidden form fields needed |
| gorilla/sessions (cookie store) | SCS v2 (server-side) | Ongoing best practice | Server-side session token; can invalidate server-side; no cookie payload bloat |
| lib/pq | pgx v5 | pgx became dominant ~2022 | Faster, native interface, LISTEN/NOTIFY support |

**Deprecated/outdated:**
- gorilla/csrf: Still maintained but CrossOriginProtection is idiomatic for Go 1.25+ projects with no need for custom token forms.
- gorilla/sessions CookieStore: Stores session data in cookie payload — avoid for security-sensitive sessions.

---

## Open Questions

1. **Local dev admin host routing**
   - What we know: `r.Host` will be `localhost:8080` locally; production is `admin.jared-wallace.com`
   - What's unclear: Should local dev use host-based routing (requires `/etc/hosts` alias) or a separate path prefix?
   - Recommendation: Add `ADMIN_HOST` env var defaulting to `admin.jared-wallace.com`; local dev sets `ADMIN_HOST=localhost` in `.env`. Strip port with `net.SplitHostPort` before comparing. This is the minimal-friction approach.

2. **SESSION_SECRET usage**
   - What we know: SCS with pgxstore does not require a server-side signing secret for the session token itself — the token is an opaque random string stored server-side
   - What's unclear: Whether an additional HMAC-signing layer on the session cookie name/value is needed
   - Recommendation: SCS pgxstore does not need a secret (unlike JWT or cookie-store approaches). Skip `SESSION_SECRET` — the Postgres-backed store is the security boundary. If added later, SCS supports `HMAC` cookie signing via a separate config.

3. **Admin mux route for /admin/preview with CSRF**
   - What we know: `CrossOriginProtection` protects all POST routes; the preview endpoint is also a POST
   - What's unclear: Whether the browser's fetch() from the editor page sends a `Sec-Fetch-Site: same-origin` header when calling `/admin/preview` on the same host
   - Recommendation: Since both the editor page and the preview endpoint are on `admin.jared-wallace.com`, the browser sends `Sec-Fetch-Site: same-origin` — no trusted origin config needed for preview. The `AddTrustedOrigin` call is only needed if preview were called cross-origin.

---

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Go | Runtime | Yes | 1.26.0 | — |
| Postgres | pgxstore sessions, post writes | Requires `make dev-up` | — (Docker) | In-memory scs store for unit tests |
| bcrypt (`x/crypto`) | Password verification | Already in go.mod | v0.49.0 | — |
| SCS v2 + pgxstore | Sessions | Already in go.mod | v2.9.0 | — |
| `net/http.CrossOriginProtection` | CSRF | stdlib Go 1.25+ | Go 1.26.0 | — |

**Missing dependencies with no fallback:** None. All dependencies are present.

**Note:** `CrossOriginProtection` is confirmed present in Go 1.26.0 (local runtime). The STATE.md blocker "Go 1.26 stdlib may include CrossOriginProtection" is resolved — it was introduced in Go 1.25, is available in 1.26, and functions correctly for this use case.

---

## Validation Architecture

### Test Framework
| Property | Value |
|----------|-------|
| Framework | Go standard `testing` package + `net/http/httptest` |
| Config file | None (go test ./...) |
| Quick run command | `go test ./internal/handler/admin/... -v -race` |
| Full suite command | `go test ./... -v -race` |

### Phase Requirements to Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| ADMIN-01 | Login with valid credentials returns redirect + session | unit (httptest) | `go test ./internal/handler/admin/... -run TestLoginPost -v` | No — Wave 0 |
| ADMIN-01 | Login with invalid credentials returns 401 + generic message | unit (httptest) | `go test ./internal/handler/admin/... -run TestLoginPostInvalid -v` | No — Wave 0 |
| ADMIN-02 | Session persists: request with valid cookie passes RequireSession | unit (httptest + mock scs) | `go test ./internal/middleware/... -run TestRequireSession -v` | No — Wave 0 |
| ADMIN-03 | Create post: valid form saves to repo, redirects | unit (httptest + mock repo) | `go test ./internal/handler/admin/... -run TestCreatePost -v` | No — Wave 0 |
| ADMIN-03 | Slug generated from title matches expected pattern | unit | `go test ./internal/service/post/... -run TestGenerateSlug -v` | No — Wave 0 |
| ADMIN-04 | Edit post: form pre-populated with existing post data | unit (httptest + mock repo) | `go test ./internal/handler/admin/... -run TestEditPost -v` | No — Wave 0 |
| ADMIN-05 | SoftDelete sets deleted_at; Restore clears it | unit (mock repo) | `go test ./internal/service/post/... -run TestSoftDelete -v` | No — Wave 0 |
| ADMIN-06 | TogglePublished flips Published field | unit (mock repo) | `go test ./internal/service/post/... -run TestTogglePublished -v` | No — Wave 0 |
| ADMIN-07 | Preview endpoint renders markdown and returns HTML fragment | unit (httptest) | `go test ./internal/handler/admin/... -run TestPreview -v` | No — Wave 0 |

### Sampling Rate
- **Per task commit:** `go test ./internal/handler/admin/... -v -race`
- **Per wave merge:** `go test ./... -v -race`
- **Phase gate:** Full suite green before `/gsd:verify-work`

### Wave 0 Gaps

All admin test files are new. Minimum required before implementation:

- [ ] `internal/handler/admin/handler_test.go` — test helper, mock setup following `internal/handler/blog/handler_test.go` pattern
- [ ] `internal/middleware/auth_test.go` — RequireSession middleware tests
- [ ] `internal/service/post/write_test.go` — slug generation, soft-delete, publish toggle
- [ ] `internal/repository/post/write_test.go` — SQL write queries (requires live Postgres, mark with `//go:build integration`)

---

## Project Constraints (from CLAUDE.md)

| Directive | Impact on Phase 3 |
|-----------|------------------|
| Go minimal dependencies — prefer stdlib | Use `net/http.CrossOriginProtection` (stdlib) over gorilla/csrf; all new deps are already in go.mod |
| Must run as Docker container on port 8080 | Session cookie `Secure=true` requires TLS termination at Nginx/ALB — already in place for production |
| All persistent data on EBS at /var/www/html | Sessions table lives in Postgres sidecar — already on EBS bind-mount |
| Leverage `frontend-design` skill for template/UI work | Admin templates should follow the same nautical aesthetic; skill directory not present but CONTEXT.md D-02 captures the design decision |
| Run `/simplify` on all code changes | Applies to all code written in this phase |

---

## Sources

### Primary (HIGH confidence)
- `pkg.go.dev/net/http#CrossOriginProtection` — Go 1.25+ CSRF protection, confirmed in Go 1.26.0
- `pkg.go.dev/github.com/alexedwards/scs/v2` — SCS SessionManager API, cookie config, middleware
- `pkg.go.dev/github.com/alexedwards/scs/pgxstore` — pgxstore constructor, sessions table schema
- `internal/handler/blog/handler.go` — Existing constructor/template pattern to mirror exactly
- `internal/repository/post/repository.go` — Existing Repository interface to extend
- `internal/model/post.go` — Confirmed: DeletedAt *time.Time and Published bool already exist
- `go.mod` — Confirmed: scs, pgxstore, x/crypto all present as indirect deps

### Secondary (MEDIUM confidence)
- [A modern approach to preventing CSRF in Go — Alex Edwards](https://www.alexedwards.net/blog/preventing-csrf-in-go) — background context on CrossOriginProtection
- [CSRF Protection in Go 1.25 — Samuel Adebayo](https://www.samueladebayo.dev/posts/golang-cross-origin-protection/) — usage patterns
- [Sliding Window Rate Limiting — Arpit Bhayani](https://arpitbhayani.me/blogs/sliding-window-ratelimiter/) — rate limiter algorithm details

### Tertiary (LOW confidence)
- WebSearch results on per-IP rate limiter patterns — general patterns verified against Go stdlib behavior

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — all packages already in go.mod; versions confirmed
- Architecture: HIGH — mirrors existing Phase 2 patterns exactly; no novel decisions
- CSRF approach: HIGH — CrossOriginProtection confirmed in Go 1.26.0 locally
- Pitfalls: HIGH — based on direct code reading and official SCS docs
- Rate limiter: MEDIUM — pattern is standard but specific implementation is at Claude's discretion

**Research date:** 2026-03-26
**Valid until:** 2026-09-26 (stable libraries; re-verify if SCS or pgx major version bumps)
