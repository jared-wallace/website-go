# Phase 3: Admin Panel - Context

**Gathered:** 2026-03-26
**Status:** Ready for planning

<domain>
## Phase Boundary

Deliver a secure admin panel at admin.jared-wallace.com where the admin can log in with email/password, and create, edit, publish, soft-delete, and restore posts using a split-pane markdown editor with live preview. No image upload, no API push endpoint, no public-facing features — those are Phase 4/5.

</domain>

<decisions>
## Implementation Decisions

### Subdomain Routing
- **D-01:** Host-based mux in the single binary. Admin handlers serve requests where `Host` matches `admin.jared-wallace.com`; all other hosts route to blog handlers. Nginx already forwards both domains to :8080.
- **D-02:** Shared nautical design — admin reuses the beach bar aesthetic (base template, color palette, typography). Admin feels like the "back office" of the same beach bar.
- **D-03:** Unauthenticated requests to admin subdomain see a branded nautical-themed login page (not a redirect to the public blog).

### Editor Experience
- **D-04:** Vanilla JS fetch for live preview. Editor textarea sends markdown to a server endpoint via debounced fetch() on keyup. Server renders with the same goldmark+bluemonday pipeline and returns HTML. Preview is always identical to published output.
- **D-05:** Manual save only. Explicit "Save Draft" and "Publish" buttons. Ctrl+S / Cmd+S keyboard shortcut. No autosave, no localStorage draft backup.
- **D-06:** On mobile/narrow screens, editor and preview stack vertically with a Write/Preview tab toggle. Full-width editing.
- **D-07:** Plain textarea with monospace font. No toolbar buttons (bold, italic, link, etc.). Admin knows markdown.

### Auth & Session Flow
- **D-08:** Admin credentials stored as environment variables: `ADMIN_EMAIL` and `ADMIN_PASSWORD_HASH`. No users/admins table. Hash generated offline via `make hash-password` Makefile target. Single admin user per requirements.
- **D-09:** SCS v2.9.0 with pgxstore for Postgres-backed sessions. Session cookie carries HttpOnly + Secure + SameSite=Lax flags per success criteria.
- **D-10:** 24-hour session lifetime (inactivity-based expiry). Re-login required after 24 hours.
- **D-11:** Failed login shows generic "Invalid email or password" message (never reveals which is wrong). In-memory rate limiter: 5 attempts per minute per IP to prevent brute force.

### Admin Dashboard
- **D-12:** Post list as a table (title, status, date, actions) with filter tabs: All / Published / Drafts / Deleted. Action links for Edit, Publish/Unpublish, Delete.
- **D-13:** Soft-delete recovery via "Restore" button in the Deleted tab. Restores post to draft status. No permanent delete in v1.
- **D-14:** Slug auto-generates from title (e.g., "My Post" -> "my-post") but admin can manually edit before saving.
- **D-15:** Delete action requires confirmation dialog. Publish/unpublish act immediately (easily reversible).

### Claude's Discretion
- CSRF implementation approach (check Go 1.26 stdlib CrossOriginProtection first, fall back to gorilla/csrf or manual tokens)
- Session middleware design and request context integration
- Admin template structure (separate admin base template extending shared styles, or full reuse of blog base.html)
- Table styling and responsive behavior on the dashboard
- Login form layout and error display
- Slug generation algorithm (Unicode handling, special character stripping)
- Debounce timing for editor preview (200-500ms range)
- Rate limiter implementation details (sync.Map, sliding window, etc.)

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Tech Stack
- `.planning/research/STACK.md` — Authoritative dependency versions: SCS v2.9.0, pgxstore, bcrypt v0.49.0, goldmark v1.8.2
- `.planning/research/ARCHITECTURE.md` — Project structure guidance and package organization patterns

### Project Context
- `.planning/PROJECT.md` — Core value, constraints, key decisions (session auth over OAuth, subdomain admin panel)
- `.planning/REQUIREMENTS.md` — ADMIN-01 through ADMIN-07 acceptance criteria
- `.planning/ROADMAP.md` — Phase 3 success criteria (5 criteria that must be TRUE)

### Prior Phases
- `.planning/phases/01-foundation/01-CONTEXT.md` — Project layout (flat internal/), config via env vars, CI choices
- `.planning/phases/02-public-blog/02-CONTEXT.md` — Nautical design decisions, template patterns, handler/service/repository structure

### Infrastructure
- `.planning/research/PITFALLS.md` — Known pitfalls for Postgres EBS bind-mount, Docker builds

### Blockers (from STATE.md)
- Go 1.26 stdlib may include CrossOriginProtection — confirm before reaching for gorilla/csrf
- Postgres sessions table migration must be created as part of this phase

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/handler/blog/handler.go` — BlogHandler pattern: constructor-based DI, templates parsed once at startup into map, centralized render() method. Admin handler should follow the same pattern.
- `internal/service/post/service.go` — Post service with ListPublished, GetBySlug. Extend with Create, Update, Delete, Publish, Unpublish, ListAll (including drafts/deleted).
- `internal/repository/post/repository.go` — Repository interface pattern. Extend with write operations (Create, Update, SoftDelete, Restore, GetByID).
- `internal/model/post.go` — Post struct already has DeletedAt (*time.Time) for soft-delete and Published (bool) for draft/publish toggle.
- `internal/markdown/renderer.go` — Goldmark + bluemonday pipeline with Render(). Reuse for editor preview endpoint.
- `internal/config/config.go` — Environment-based config. Extend with AdminEmail, AdminPasswordHash, SessionSecret.
- `web/templates/base.html` — Base template with block "content". Admin pages can extend or create admin-specific base.
- `web/static/style.css` — Nautical CSS. Admin pages share the palette and typography.

### Established Patterns
- Per-page template sets: parse base.html + page.html separately per handler (avoids block name collisions)
- Constructor-based dependency injection in main.go
- Interface-based repository (mockable for tests)
- FuncMap injected at template parse time
- `go:embed` for templates and static assets
- `log/slog` for structured logging

### Integration Points
- `cmd/server/main.go` — Add host-based routing, SCS session manager init, admin handler wiring
- `db/migrations/` — Add sessions table migration (00003_create_sessions.sql)
- `web/templates/` — Add admin templates (login, dashboard, editor)
- `web/static/` — Add admin-specific CSS/JS (editor preview logic)
- `go.mod` — SCS and bcrypt already declared as indirect deps; need direct import

</code_context>

<specifics>
## Specific Ideas

No specific requirements — open to standard approaches. User consistently chose recommended options across all areas, confirming preference for conventional, simple patterns. The admin panel should be functional and straightforward — a tool for writing, not a showpiece.

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.

</deferred>

---

*Phase: 03-admin-panel*
*Context gathered: 2026-03-26*
