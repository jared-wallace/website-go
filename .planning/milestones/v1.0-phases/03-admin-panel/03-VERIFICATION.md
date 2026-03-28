---
phase: 03-admin-panel
verified: 2026-03-27T00:00:00Z
status: human_needed
score: 5/5 must-haves verified
human_verification:
  - test: "Visit http://localhost:8080/admin/login (with ADMIN_HOST=localhost). Enter wrong password, see 'Invalid email or password.' Enter correct credentials, confirm redirect to /admin/posts."
    expected: "Login page renders with nautical branding. Wrong credentials show generic error. Correct credentials redirect. Session cookie has HttpOnly and SameSite=Lax in DevTools."
    why_human: "Cookie flag inspection, visual layout, and actual bcrypt timing cannot be verified programmatically without a live server."
  - test: "Log in, create a post titled 'Hello World' and watch the slug field. Type in the markdown body and observe the preview pane."
    expected: "Slug auto-populates as 'hello-world' as you type. Preview pane updates within ~300ms of each keystroke. Ctrl+S (or Cmd+S) triggers save."
    why_human: "JavaScript debounce behavior, live DOM updates, and keyboard shortcut handling require a browser."
  - test: "From the dashboard, publish a post, delete it via the confirmation dialog, then restore it. Check each filter tab (All, Published, Drafts, Deleted)."
    expected: "Status pill changes correctly. Native <dialog> appears on delete. Flash messages confirm each action. Filter tabs show correct subsets."
    why_human: "Visual confirmation of status pills, native dialog rendering, and flash message display require human eyes."
  - test: "Resize browser to mobile width (< 768px) in the editor."
    expected: "Editor shows Write/Preview tab toggle. Date column hidden in dashboard table. Toolbar buttons stack vertically."
    why_human: "Responsive layout behavior requires browser viewport manipulation."
  - test: "Log in, close the browser tab, reopen http://localhost:8080/admin/posts."
    expected: "Admin is still authenticated (Postgres-backed session survives tab close)."
    why_human: "Session persistence across browser sessions requires actual browser interaction."
---

# Phase 3: Admin Panel Verification Report

**Phase Goal:** The admin can securely log in at admin.jared-wallace.com and create, edit, publish, and soft-delete posts using a split-pane markdown editor with live preview.
**Verified:** 2026-03-27
**Status:** human_needed — all automated checks pass; 5 behavioral items require live browser verification
**Re-verification:** No — initial verification

---

## Goal Achievement

### Observable Truths (from ROADMAP.md Success Criteria)

| # | Truth | Status | Evidence |
|---|-------|--------|----------|
| 1 | Admin logs in with email/password; session persists (Postgres-backed); cookie carries HttpOnly + Secure + SameSite=Lax | VERIFIED | `auth.go`: bcrypt.CompareHashAndPassword, RenewToken, sessions.Put "authenticated". `main.go`: pgxstore.New, HttpOnly=true, SameSite=Lax, Secure=production-only. Cookie flags require human browser check. |
| 2 | Admin can create a new post with title, body, and slug; post saves to the database | VERIFIED | `editor.go:SavePost` calls `h.svc.Create`. `write.go:Create` calls `renderer.Render` then `r.pool.QueryRow` INSERT. `write.go:ErrSlugExists` handles unique constraint. |
| 3 | Admin can edit an existing post and soft-delete it (post is recoverable, not destroyed) | VERIFIED | `editor.go:SavePost` edit path calls `h.svc.Update`. `actions.go:DeletePost` calls `h.svc.SoftDelete`. `write.go:SoftDelete` sets `deleted_at = now()`, never destroys row. `actions.go:RestorePost` calls `h.svc.Restore`. |
| 4 | Admin can toggle a post between draft and published; only published posts appear on public routes | VERIFIED | `actions.go:PublishPost/UnpublishPost` call `h.svc.Publish/Unpublish`. `write.go:SetPublished` updates `published` column. Public `ListPublished` query filters `WHERE published = true AND deleted_at IS NULL`. |
| 5 | Admin writes in a split-pane editor with live markdown preview updating as they type | VERIFIED | `admin-editor.html` has `admin-layout` grid with `editor-pane` + `preview-pane`. `admin.js` fetches `/admin/preview` with 300ms debounce. `preview.go:Preview` calls `h.renderer.Render` and writes raw HTML. |

**Score:** 5/5 truths verified

---

## Required Artifacts

| Artifact | Expected | Status | Details |
|----------|----------|--------|---------|
| `internal/repository/post/write.go` | Write query implementations for admin CRUD | VERIFIED | 7 methods: FindByID, ListAll, Create, Update, SoftDelete, Restore, SetPublished. ErrSlugExists on code 23505. Real pgx queries, not static returns. |
| `internal/service/post/write.go` | Service layer admin operations | VERIFIED | 8 methods. Create/Update call renderer.Render before repo. All others delegate to repo. |
| `internal/service/post/slug.go` | Slug generation from title | VERIFIED | GenerateSlug: lowercase, regex replace, trim, unix timestamp fallback. 6 test cases pass. |
| `internal/middleware/auth.go` | RequireSession middleware | VERIFIED | Checks sm.GetBool "authenticated", redirects to /admin/login with 303. |
| `internal/middleware/ratelimit.go` | In-memory rate limiter | VERIFIED | Fixed-window counter per IP, sync.Mutex, port stripping. 5 tests pass including window reset. |
| `db/migrations/00003_create_sessions.sql` | Sessions table for pgxstore | VERIFIED | CREATE TABLE sessions (token, data, expiry), sessions_expiry_idx. |
| `internal/config/config.go` | Admin env var loading | VERIFIED | AdminEmail, AdminPasswordHash, AdminHost, SessionSecret all via envOr (non-panicking). |
| `internal/handler/admin/handler.go` | AdminHandler constructor, render method | VERIFIED | Parses all 3 template pages. Pre-computes dummyHash at startup. funcMap with formatDate/currentYear. |
| `internal/handler/admin/auth.go` | Login, LoginPost, Logout handlers | VERIFIED | Constant-time bcrypt check with dummyHash fallback. RenewToken before setting authenticated. extractIP checks X-Real-IP header. |
| `internal/handler/admin/dashboard.go` | Dashboard with filter logic | VERIFIED | ListAll + in-memory filter by tab. PopString for flash messages. |
| `internal/handler/admin/actions.go` | DeletePost, RestorePost, PublishPost, UnpublishPost | VERIFIED | All parse path ID, call service method, set flash message, redirect. Flash copy matches UI-SPEC verbatim. |
| `internal/handler/admin/editor.go` | NewPost, EditPost, SavePost | VERIFIED | Create/Update paths wired. ErrSlugExists handled. RenderedHTML cast to template.HTML to prevent double-escaping. |
| `internal/handler/admin/preview.go` | Live preview endpoint | VERIFIED | Calls h.renderer.Render, writes raw HTML with Content-Type text/html. No template wrapping. |
| `cmd/server/main.go` | Host router, SCS wiring, admin routes | VERIFIED | hostRouter dispatches by Host header. All 10 admin routes registered and wrapped with requireAuth. CrossOriginProtection applied. |
| `web/templates/admin-base.html` | Base template with nav | VERIFIED | "The Log -- Back Office" nav, "Go Ashore" logout, FlashSuccess/FlashError blocks, admin.css + admin.js refs. |
| `web/templates/admin-login.html` | Nautical login form | VERIFIED | Centered form, "Sign In" CTA, labeled email/password fields, role="alert" on error. |
| `web/templates/admin-dashboard.html` | Post table with filter tabs | VERIFIED | filter-tabs, admin-table, status-pill (published/draft/deleted), confirm-dialog, empty states for all 4 filter tabs, "New Post" button. |
| `web/templates/admin-editor.html` | Split-pane editor | VERIFIED | admin-layout grid, editor-body textarea, preview-content div, tab-write/tab-preview, Save Draft + Publish buttons, aria-label + aria-live attributes. |
| `web/static/admin.css` | Admin styling | VERIFIED | --color-danger token, .admin-login, .form-field, .cta-button, .admin-layout (grid-template-columns: 1fr 1fr), .admin-tab-bar, .admin-toolbar, .filter-tab.active, .status-pill variants, .confirm-dialog, .empty-state. |
| `web/static/admin.js` | Client-side interactivity | VERIFIED | fetch('/admin/preview') with 300ms setTimeout, ctrlKey/metaKey shortcut, slugManuallyEdited logic, tab-write/tab-preview toggle. |
| `cmd/hashpw/main.go` | bcrypt hash generator CLI | VERIFIED | bcrypt.GenerateFromPassword cost 12. Outputs $2a$12... prefix. Makefile `hash-password` target wired. |

---

## Key Link Verification

| From | To | Via | Status | Details |
|------|----|-----|--------|---------|
| `cmd/server/main.go` | `internal/handler/admin/handler.go` | `adminhandler.New()` constructor | WIRED | Line 112: `adminH := adminhandler.New(svc, sessionManager, renderer, rl, cfg)` |
| `cmd/server/main.go` | `internal/middleware/auth.go` | `middleware.RequireSession` wrapping admin routes | WIRED | Line 119: `requireAuth := middleware.RequireSession(sessionManager)` used on all 7 protected routes |
| `internal/handler/admin/auth.go` | `golang.org/x/crypto/bcrypt` | `bcrypt.CompareHashAndPassword` | WIRED | Line 54: `err := bcrypt.CompareHashAndPassword(hashToCheck, []byte(password))` |
| `internal/handler/admin/dashboard.go` | `internal/service/post/write.go` | `svc.ListAll()` | WIRED | Line 13: `posts, err := h.svc.ListAll(r.Context())` |
| `internal/handler/admin/actions.go` | `internal/service/post/write.go` | `svc.SoftDelete/Restore/Publish/Unpublish` | WIRED | Lines 23, 40, 57, 74: all delegate to matching service methods |
| `internal/handler/admin/editor.go` | `internal/service/post/write.go` | `svc.Create` and `svc.Update` | WIRED | Lines 90, 129: create and edit paths both call correct service methods |
| `internal/handler/admin/preview.go` | `internal/markdown/renderer.go` | `h.renderer.Render()` | WIRED | Line 12: `rendered := h.renderer.Render(body)` |
| `web/static/admin.js` | `/admin/preview` | debounced `fetch()` on textarea input | WIRED | Line 18: `fetch('/admin/preview', { method: 'POST', body: formData })` with 300ms debounce |
| `internal/service/post/write.go` | `internal/repository/post/write.go` | `s.repo.*` delegation | WIRED | All 8 service methods call matching repo methods. Create/Update call renderer.Render first. |

---

## Data-Flow Trace (Level 4)

| Artifact | Data Variable | Source | Produces Real Data | Status |
|----------|---------------|--------|-------------------|--------|
| `admin-dashboard.html` | `.Posts` | `h.svc.ListAll` -> `r.pool.Query` SELECT with no WHERE filter | Yes — DB query returns all rows | FLOWING |
| `admin-editor.html` | `.Post` | `h.svc.GetByID` -> `r.pool.QueryRow` SELECT WHERE id = $1 | Yes — fetches specific post by ID | FLOWING |
| `preview-content` div | live HTML | POST /admin/preview -> `h.renderer.Render(body)` | Yes — goldmark+bluemonday pipeline | FLOWING |

---

## Behavioral Spot-Checks

| Behavior | Command | Result | Status |
|----------|---------|--------|--------|
| Build compiles clean | `go build ./...` | exit 0, no output | PASS |
| All package tests pass | `go test ./...` | 11 packages: all ok/skipped, 0 failures | PASS |
| hashpw produces bcrypt output | `go run ./cmd/hashpw "testpass"` | `$2a$12...` | PASS |
| Admin handler tests: auth, actions, editor, preview | `go test ./internal/handler/admin/... -v -race` | 16 tests, all PASS | PASS |
| Middleware tests: RequireSession, RateLimiter | `go test ./internal/middleware/... -v -race` | 7 tests, all PASS | PASS |
| Service write tests | `go test ./internal/service/post/... -v -race` | 7 write tests + slug tests, all PASS | PASS |

---

## Requirements Coverage

| Requirement | Source Plans | Description | Status | Evidence |
|-------------|-------------|-------------|--------|---------|
| ADMIN-01 | 03-01, 03-02 | Admin can log in with email/password at admin.jared-wallace.com | SATISFIED | auth.go: bcrypt comparison, RenewToken, session.Put authenticated. Rate limiting enforced. |
| ADMIN-02 | 03-01, 03-02 | Admin session persists across browser refresh (Postgres-backed) | SATISFIED | main.go: pgxstore.New(pool), SCS IdleTimeout 24h, Lifetime 30d, HttpOnly, SameSite=Lax. |
| ADMIN-03 | 03-01, 03-03, 03-04 | Admin can create posts with title, markdown body, and slug | SATISFIED | editor.go:SavePost -> svc.Create -> repo.Create with real INSERT. ErrSlugExists handled. |
| ADMIN-04 | 03-01, 03-04 | Admin can edit existing posts | SATISFIED | editor.go:SavePost edit path -> svc.Update -> repo.Update. EditPost pre-populates form fields. |
| ADMIN-05 | 03-01, 03-03 | Admin can soft-delete posts (recoverable) | SATISFIED | actions.go:DeletePost -> svc.SoftDelete sets deleted_at. RestorePost clears it. Row never destroyed. |
| ADMIN-06 | 03-01, 03-03 | Admin can toggle posts between draft and published states | SATISFIED | actions.go:PublishPost/UnpublishPost -> svc.Publish/Unpublish -> repo.SetPublished. |
| ADMIN-07 | 03-04 | Admin can write in split-pane markdown editor with live preview | SATISFIED | admin-editor.html: admin-layout grid. admin.js: debounced fetch to /admin/preview. preview.go: raw HTML response. |

All 7 ADMIN-0x requirements for Phase 3 are SATISFIED. No orphaned requirements found (REQUIREMENTS.md traceability table matches).

---

## Anti-Patterns Found

| File | Pattern | Severity | Impact |
|------|---------|----------|--------|
| `handler_test.go`, `actions_test.go` | `errors.New("not implemented")` in mock struct methods | INFO | Test mocks only — these are unexercised mock branches for interface compliance. Not in production code paths. |
| `admin-editor.html` | HTML `placeholder` attribute on input/textarea | INFO | Standard HTML UX hint attributes, not code stubs. |

No blocker or warning anti-patterns found. The "not implemented" strings appear exclusively in test mock structs for interface satisfaction, not in production handler code.

---

## Human Verification Required

### 1. Session Cookie Security Flags

**Test:** Open browser DevTools > Application > Cookies after logging in at http://localhost:8080/admin/login (with ADMIN_HOST=localhost env var). Use correct credentials.
**Expected:** Cookie named `admin_session` is present with HttpOnly=true, SameSite=Lax. Secure=false in dev (http), Secure=true when APP_ENV=production.
**Why human:** Cookie flag inspection requires a live browser session. Code sets these flags but enforcement is browser-verified.

### 2. Live Preview (300ms Debounce)

**Test:** In the editor, type in the markdown body field. Watch the right pane.
**Expected:** Preview pane updates within approximately 300ms of typing stopping. Rendered markdown (bold, headings, code blocks) appears. No update fires on every keystroke (debounce is working).
**Why human:** JavaScript timing and DOM update behavior requires a browser.

### 3. Slug Auto-Generation + Manual Override

**Test:** On the New Post page, type a title like "Hello World". Then manually edit the slug field. Continue editing the title.
**Expected:** Slug auto-populates as "hello-world" while untouched. After manual edit, slug no longer auto-updates when title changes.
**Why human:** DOM event sequencing and the slugManuallyEdited flag logic require browser interaction.

### 4. Ctrl+S Keyboard Shortcut

**Test:** In the editor with unsaved changes, press Ctrl+S (Mac: Cmd+S).
**Expected:** Form submits as "draft". Browser does not trigger its native Save Page dialog.
**Why human:** Keyboard event interception requires a browser to confirm e.preventDefault() works correctly.

### 5. Mobile Responsive Layout

**Test:** Open the editor and dashboard in browser DevTools device emulation or a physical mobile device.
**Expected:** Editor shows Write/Preview tab toggle instead of side-by-side panes. Dashboard hides the date column. Toolbar buttons stack vertically. Minimum touch target size (44px) is met.
**Why human:** Responsive CSS and layout require viewport manipulation.

---

## Gaps Summary

No automated gaps. All 5 phase success criteria are verified at the code level. The 7 ADMIN requirements are all satisfied with real implementations (no stubs in production code paths). `go build ./...` and all tests pass cleanly. The remaining items are human verification of browser-dependent behaviors (cookie flags, JavaScript interactivity, responsive layout) which cannot be programmatically confirmed without a running server and browser.

---

_Verified: 2026-03-27_
_Verifier: Claude (gsd-verifier)_
