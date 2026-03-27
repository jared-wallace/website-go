# Phase 5: API + Images - Context

**Gathered:** 2026-03-27
**Status:** Ready for planning

<domain>
## Phase Boundary

Deliver two capabilities: (1) image upload from the admin editor with magic-byte validation and EBS storage, and (2) a bearer-token-authenticated API endpoint for pushing .md files from a local machine. No new public-facing features, no media library UI, no image optimization — those are future work.

</domain>

<decisions>
## Implementation Decisions

### Image Upload Workflow
- **D-01:** Upload lives inside the post editor — a button above/beside the textarea. On upload, a markdown image tag (`![alt](/images/xyz.jpg)`) is auto-inserted at cursor position.
- **D-02:** Upload button only — no drag-and-drop, no clipboard paste. Consistent with the "plain textarea, no toolbar" decision from Phase 3 (D-07).
- **D-03:** Accepted formats: JPEG and PNG only. Validated via magic bytes.
- **D-04:** Max file size: 5 MB per upload.

### Image Storage & Serving
- **D-05:** Flat directory on EBS: `/var/www/html/images/{random-name}.{ext}`. No date or post bucketing. Blog scale won't hit filesystem limits.
- **D-06:** No database tracking of images. The markdown in post body IS the reference. No images table, no migration needed for this.
- **D-07:** URL scheme: `/images/{filename}`. Go serves files from the EBS images directory. Simple, readable, standard blog convention.

### API Push Endpoint
- **D-08:** Bearer token managed as a single env var (`API_TOKEN`), consistent with existing credential pattern (ADMIN_EMAIL, ADMIN_PASSWORD_HASH from Phase 3 D-08). Rotate by redeploying with new value.
- **D-09:** Request body is raw markdown with `Content-Type: text/markdown`. Slug and title extracted from YAML front matter via goldmark-meta (already in the dependency tree). Simplest CLI usage: `curl -H "Authorization: Bearer $TOKEN" -d @post.md`.
- **D-10:** Upsert by slug — if slug exists, update body/title; if not, create new post. Idempotent and safe to re-push.
- **D-11:** Pushed posts land as drafts by default. Admin publishes from the web UI when ready.

### Security Boundaries
- **D-12:** Magic-byte validation via Go stdlib `net/http.DetectContentType` (reads first 512 bytes, sniffs MIME type). Checks for JPEG (FF D8 FF) and PNG (89 50 4E 47) signatures. No external deps.
- **D-13:** No rate limiting on the API push endpoint. Bearer token is a strong enough gate for a single-admin blog.
- **D-14:** Bump WriteTimeout to 30s on the image upload handler only (per-handler timeout override). Current 10s server-wide timeout stays for all other routes.

### Claude's Discretion
- Random filename generation strategy (UUID, crypto/rand hex, etc.)
- How to wire the `/images/` file server into the existing mux (http.StripPrefix + http.FileServer or custom handler)
- Config struct extension for `API_TOKEN` and image directory path
- Front matter schema for the push endpoint (which YAML keys are required vs optional)
- Upload button placement and styling within the editor template
- Error response format for the API endpoint (plain text vs JSON)
- Whether the upload endpoint returns the inserted markdown tag or just the URL
- ReadTimeout adjustment for upload route (if needed alongside WriteTimeout bump)

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Tech Stack
- `.planning/research/STACK.md` -- Authoritative dependency versions: goldmark v1.8.2, goldmark-meta for YAML front matter, pgx v5.9.1
- `.planning/research/ARCHITECTURE.md` -- Project structure guidance and package organization patterns

### Project Context
- `.planning/PROJECT.md` -- Core value, constraints, EBS storage at /var/www/html, minimal deps philosophy
- `.planning/REQUIREMENTS.md` -- ADMIN-08 (image upload), ADMIN-09 (API push) acceptance criteria

### Prior Phase Context
- `.planning/phases/03-admin-panel/03-CONTEXT.md` -- Admin handler patterns, editor decisions (D-04 vanilla JS, D-07 plain textarea, D-08 env var credentials)
- `.planning/phases/04-distribution/04-CONTEXT.md` -- D-06 static fallback OG image (Phase 5 enables per-post images)

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/config/config.go`: Env-var config pattern — extend with `API_TOKEN` and `ImageDir` fields
- `internal/middleware/auth.go`: RequireSession middleware — new bearer token middleware follows same pattern
- `internal/handler/admin/handler.go`: AdminHandler with editor templates — extend with upload button and handler
- `internal/handler/admin/editor.go`: Editor page — upload button integrates here
- `internal/markdown/renderer.go`: goldmark + bluemonday pipeline — already handles front matter via goldmark-meta
- `internal/middleware/ratelimit.go`: RateLimiter — available if API rate limiting is ever reconsidered

### Established Patterns
- Repository -> Service -> Handler layering (post CRUD follows this)
- Config via env vars with `mustEnv()` / `envOr()` helpers
- SCS session-based auth for admin routes; new API auth is token-based (different middleware)
- template.HTML for pre-rendered content prevents double-escaping
- Per-page template sets parsed with admin-base.html

### Integration Points
- `cmd/server/main.go`: Mux setup — add `/images/` file server, `POST /api/push` route, upload handler on admin mux
- `internal/config/config.go`: Add `APIToken` and `ImageDir` fields
- `internal/server/server.go`: WriteTimeout may need per-handler override for upload route
- `web/templates/admin-editor.html`: Add upload button UI
- `web/static/js/main.js`: Add upload fetch logic (follows existing vanilla JS pattern from Phase 3 D-04)

</code_context>

<specifics>
## Specific Ideas

No specific requirements -- open to standard approaches.

</specifics>

<deferred>
## Deferred Ideas

None -- discussion stayed within phase scope.

</deferred>

---

*Phase: 05-api-images*
*Context gathered: 2026-03-27*
