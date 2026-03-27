---
phase: 03-admin-panel
plan: 04
subsystem: ui
tags: [go, html-template, markdown, admin, editor, javascript, css]

requires:
  - phase: 03-03
    provides: dashboard, post actions, session auth, mock patterns for tests

provides:
  - Split-pane markdown editor with NewPost/EditPost/SavePost handlers
  - Live preview endpoint (POST /admin/preview) returning raw HTML
  - admin.js: debounced preview (300ms), Ctrl+S shortcut, slug auto-generation, mobile tab toggle
  - admin-editor.html: full editor template with title/slug/tags/toolbar
  - admin.css editor additions: 50/50 grid desktop, stacked mobile with tab bar
  - Handler tests: TestNewPost, TestEditPost, TestSavePost, TestPreview

affects:
  - phase 04 (any future feature work referencing editor behavior)
  - deploy (admin.js served as static asset)

tech-stack:
  added: []
  patterns:
    - "Renderer interface in admin handler package mirrors service.Renderer pattern — mock-friendly without importing markdown package"
    - "Editor template uses postView struct to pass template.HTML RenderedHTML safely without double-escaping"
    - "Preview endpoint writes raw goldmark+bluemonday output directly (no template) to avoid double-encoding"
    - "admin.js uses IIFE and var declarations (no ES6+ module syntax) for broadest browser compatibility"

key-files:
  created:
    - internal/handler/admin/editor.go
    - internal/handler/admin/preview.go
    - internal/handler/admin/editor_test.go
    - web/static/admin.js
  modified:
    - internal/handler/admin/handler.go
    - web/templates/admin-editor.html
    - web/static/admin.css

key-decisions:
  - "Renderer interface added to admin package so handler tests can pass a mock renderer without importing markdown package (aligns with service layer pattern)"
  - "postView struct used in EditPost to expose RenderedHTML as template.HTML type, preventing html/template double-escaping"
  - "Preview writes h.renderer.Render() output directly to ResponseWriter — no template execution layer"

patterns-established:
  - "Handler interface for dependencies (Renderer): accept interface not concrete type in New() signature"

requirements-completed: [ADMIN-03, ADMIN-04, ADMIN-07]

duration: 7min
completed: 2026-03-27
---

# Phase 03 Plan 04: Admin Editor Summary

**Split-pane markdown editor with live preview (debounced 300ms fetch to /admin/preview), Ctrl+S shortcut, slug auto-generation, and mobile Write/Preview tab toggle — completing the admin CRUD loop**

## Performance

- **Duration:** 7 min
- **Started:** 2026-03-27T04:28:22Z
- **Completed:** 2026-03-27T04:35:30Z
- **Tasks:** 2 of 3 (Task 3 is checkpoint:human-verify — see below)
- **Files modified:** 7

## Accomplishments
- NewPost, EditPost, SavePost, and Preview handlers fully implemented (replaced 501 stubs)
- Preview endpoint returns raw HTML from goldmark+bluemonday pipeline without template wrapping (prevents double-encoding)
- admin.js provides all client-side editor features: 300ms debounced preview, Ctrl+S/Cmd+S save, slug auto-generation mirroring Go's GenerateSlug algorithm, mobile Write/Preview tab toggle
- admin-editor.html replaces placeholder stub with full split-pane layout (50/50 grid desktop, stacked mobile)
- Handler tests: 6 new tests covering NewPost, EditPost (valid and invalid ID), SavePost (draft and publish), and Preview
- All 16 admin handler tests pass

## Task Commits

1. **Task 1: Editor handlers, preview endpoint, handler tests, and admin.js** - `d67934b` (feat)
2. **Task 2: Editor template and CSS layout** - `b7086b0` (feat)
3. **Task 3: Visual and functional verification** - PENDING (checkpoint:human-verify — not executed by this agent)

## Files Created/Modified
- `internal/handler/admin/editor.go` - NewPost, EditPost, SavePost handlers with slug auto-gen and ErrSlugExists handling
- `internal/handler/admin/preview.go` - Preview endpoint (raw HTML output, no template)
- `internal/handler/admin/editor_test.go` - 6 handler tests for editor and preview
- `internal/handler/admin/handler.go` - Removed 501 stubs; added Renderer interface; New() accepts Renderer instead of *markdown.Renderer
- `web/static/admin.js` - Debounced preview, Ctrl+S, slug gen, mobile tabs (IIFE, vanilla JS)
- `web/templates/admin-editor.html` - Full split-pane editor template replacing placeholder stub
- `web/static/admin.css` - Editor layout additions: split-pane grid, tab bar, monospace textarea, toolbar, mobile responsive

## Decisions Made
- Used `Renderer` interface in `handler.go` (instead of `*markdown.Renderer`) so editor tests can pass mock renderers without importing the markdown package — consistent with the service layer's existing `Renderer` interface pattern.
- `postView` struct in `EditPost` exposes `RenderedHTML` as `template.HTML` so the preview pane renders pre-computed HTML without escaping in the template.
- Preview endpoint writes directly to `http.ResponseWriter` (no template execution) to avoid double-encoding the already-sanitized goldmark+bluemonday output.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Fixed nil renderer panic in TestPreview**
- **Found during:** Task 1 (editor handler tests)
- **Issue:** All three test setups (handler_test, actions_test, editor_test) pass `nil` for the renderer parameter in `admin.New()`. This was fine while Preview wasn't tested, but `preview.go` calls `h.renderer.Render()` on a nil `*markdown.Renderer`, causing a nil pointer dereference.
- **Fix:** Refactored `handler.go` to define a `Renderer` interface and changed `New()` to accept `Renderer` instead of `*markdown.Renderer`. Updated `editor_test.go` to pass `noopRendererEditor{}` (which satisfies `Renderer`) instead of `nil`. Production call site in `main.go` passes `*markdown.Renderer` which satisfies the interface — no change needed there.
- **Files modified:** `internal/handler/admin/handler.go`, `internal/handler/admin/editor_test.go`
- **Verification:** `go build ./...` passes; `TestPreviewReturnsHTML` passes
- **Committed in:** d67934b (Task 1 commit)

---

**Total deviations:** 1 auto-fixed (Rule 1 - bug)
**Impact on plan:** Necessary for test correctness. No scope creep — the interface refactor aligns with the existing service layer pattern noted in STATE.md decisions.

## Issues Encountered
None beyond the deviation above.

## Checkpoint: Task 3 Pending Human Verification

Task 3 (`type="checkpoint:human-verify"`) requires end-to-end manual verification of the complete admin panel. This task is a blocking gate — it cannot be automated. The orchestrator must present the following verification steps to the user:

**Prerequisites:**
1. Start dev environment: `make dev-up`
2. Set env vars:
   ```
   export ADMIN_EMAIL="admin@test.com"
   export ADMIN_PASSWORD_HASH="$(make hash-password PW=testpass123)"
   export ADMIN_HOST="localhost"
   export DATABASE_URL="postgres://postgres:postgres@localhost:5432/website?sslmode=disable"
   ```
3. Run server: `make run` or `go run ./cmd/server`

**Verification checklist (12 steps):**
1. Login page at http://localhost:8080/admin/login with nautical branding; wrong password shows error
2. Dashboard empty state shows "Nothing here yet." and "New Post" button
3. Create post: title auto-generates slug, live preview updates within 300ms, Save Draft redirects with flash
4. Dashboard shows post with "Draft" status pill
5. Publish action changes status to "Published"
6. Edit: fields pre-populated, body editable
7. Delete: confirm dialog, redirect to Deleted tab with flash
8. Restore: from Deleted tab, restores to drafts
9. Mobile: Write/Preview tab toggle visible, table date column hidden
10. Ctrl+S / Cmd+S: saves as draft
11. Logout: "Go Ashore" redirects to login
12. Session persistence: re-open tab stays logged in; cookie `admin_session` is HttpOnly, SameSite=Lax

**Resume signal:** Type "approved" or describe issues to fix.

## Known Stubs
None — all plan artifacts are fully implemented. Task 3 is a human verification checkpoint (not a code stub).

## Next Phase Readiness
- Admin CRUD loop complete pending human verification (Task 3)
- All handler tests pass (16 tests, -race clean)
- `go build ./...` clean
- After human verification approves, Phase 03 is ready for transition

---
*Phase: 03-admin-panel*
*Completed: 2026-03-27*

## Self-Check: PASSED

- editor.go: FOUND
- preview.go: FOUND
- editor_test.go: FOUND
- admin.js: FOUND
- admin-editor.html: FOUND
- Commit d67934b: FOUND
- Commit b7086b0: FOUND
