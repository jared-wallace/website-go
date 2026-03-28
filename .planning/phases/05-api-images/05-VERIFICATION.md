---
phase: 05-api-images
verified: 2026-03-28T01:00:00Z
status: passed
score: 4/4 must-haves verified
re_verification: false
---

# Phase 5: API + Images Verification Report

**Phase Goal:** The admin can upload images for embedding in posts and push .md files from a local machine via a bearer-token-authenticated API endpoint.
**Verified:** 2026-03-28T01:00:00Z
**Status:** passed
**Re-verification:** No -- initial verification

## Goal Achievement

### Observable Truths

| # | Truth | Status | Evidence |
|---|-------|--------|----------|
| 1 | Admin uploads a JPEG or PNG via the admin panel; image is stored on the EBS volume with a server-generated random filename; client-supplied filename is never used | VERIFIED | `upload.go` uses `crypto/rand` for 32-hex filenames, `DetectContentType` for magic-byte validation. `upload_test.go` TestUploadImage_ValidJPEG confirms file on disk with correct name length. TestUploadImage_RandomFilename confirms two uploads produce different names. Editor template has upload button, admin.js has `insertAtCursor` wired to `/admin/images/upload`. |
| 2 | A file with a spoofed MIME extension (e.g., .jpg with HTML content) is rejected at the magic-byte check | VERIFIED | `upload.go` line 48: `http.DetectContentType(buf[:n])` sniffs magic bytes. Returns 415 "only JPEG and PNG accepted" when MIME does not match. TestUploadImage_SpoofedMIME passes with HTML content named `evil.jpg`. |
| 3 | `POST /api/push` with a valid bearer token accepts a .md file body and creates or upserts a post by slug | VERIFIED | `handler.go` (api package) extracts slug/title from YAML front matter via `RenderWithMeta`, calls `svc.UpsertBySlug`. `write.go` UpsertBySlug creates with `Published: false` or updates preserving tags. `main.go` line 128: `blogMux.Handle("POST /api/push", requireToken(...))`. TestPushPost_ValidToken confirms 200 + correct slug/title. TestUpsertBySlug_NewPost confirms draft creation. TestUpsertBySlug_ExistingPost confirms upsert. |
| 4 | `POST /api/push` with a missing or invalid token returns 401; no post is created | VERIFIED | `apitoken.go` uses `subtle.ConstantTimeCompare`, returns 401 for missing/invalid/empty tokens. Empty config token disables endpoint entirely. TestRequireAPIToken_{MissingHeader,InvalidToken,EmptyBearer,EmptyConfigToken} all pass. TestPushPost_NoToken and TestPushPost_InvalidToken confirm 401 end-to-end through middleware+handler chain. |

**Score:** 4/4 truths verified

### Required Artifacts

| Artifact | Expected | Status | Details |
|----------|----------|--------|---------|
| `internal/handler/admin/upload.go` | UploadImage handler with MIME validation | VERIFIED | 107 lines. Has `DetectContentType`, `crypto/rand`, `file.Seek(0, io.SeekStart)`, `SetWriteDeadline`, `extensionFromMIME`. Wired in `main.go` line 151. |
| `internal/handler/admin/upload_test.go` | Unit tests for upload edge cases | VERIFIED | 183 lines. 6 tests: ValidJPEG, ValidPNG, SpoofedMIME, TooLarge, RandomFilename, NoFile. All pass. |
| `internal/config/config.go` | ImageDir and APIToken config fields | VERIFIED | `ImageDir string` (line 19), `APIToken string` (line 18), both loaded via `envOr`. |
| `internal/middleware/apitoken.go` | RequireAPIToken with constant-time comparison | VERIFIED | 39 lines. Uses `subtle.ConstantTimeCompare`. Empty token disables endpoint. |
| `internal/middleware/apitoken_test.go` | 5 middleware tests | VERIFIED | ValidToken, MissingHeader, InvalidToken, EmptyBearer, EmptyConfigToken. All pass. |
| `internal/handler/api/handler.go` | APIHandler with PushPost method | VERIFIED | 77 lines. Uses `LimitReader(maxBody+1)` for exact size enforcement, extracts slug/title from meta, calls `UpsertBySlug`. |
| `internal/handler/api/handler_test.go` | 6 API handler tests | VERIFIED | ValidToken, NoToken, InvalidToken, NoSlug, NoTitle, BodyTooLarge. All pass with middleware chained. |
| `internal/service/post/write.go` | UpsertBySlug service method | VERIFIED | Lines 43-67. Creates with `Published: false`, updates preserving `existing.Tags`. |
| `internal/service/post/upsert_test.go` | 3 upsert tests | VERIFIED | NewPost (draft check), ExistingPost (tags preserved), FindError (propagation). All pass. |
| `web/templates/admin-editor.html` | Upload Image button in editor | VERIFIED | Contains `id="upload-image-btn"`, `class="action-link"`, `id="upload-image-input"` with `accept="image/jpeg,image/png"`, `id="upload-error"`. |
| `web/static/admin.js` | Upload fetch + insertAtCursor | VERIFIED | Contains `fetch('/admin/images/upload'`, `insertAtCursor`, `data.markdownTag`, `'Uploading\u2026'`. |
| `cmd/server/main.go` | All routes wired | VERIFIED | `os.MkdirAll(cfg.ImageDir)` (line 62), `http.FileServer(http.Dir(cfg.ImageDir))` (line 113), `GET /images/` on both muxes (lines 117, 152), `POST /admin/images/upload` (line 151), `POST /api/push` (line 128), `RequireAPIToken` (line 127). |

### Key Link Verification

| From | To | Via | Status | Details |
|------|----|-----|--------|---------|
| `web/static/admin.js` | `/admin/images/upload` | fetch POST with FormData | WIRED | `fetch('/admin/images/upload', { method: 'POST', body: formData })` at line 134 |
| `cmd/server/main.go` | `internal/handler/admin/upload.go` | adminMux route registration | WIRED | `adminMux.Handle("POST /admin/images/upload", requireAuth(http.HandlerFunc(adminH.UploadImage)))` at line 151 |
| `cmd/server/main.go` | `http.FileServer` | /images/ static file server on both muxes | WIRED | `blogMux.Handle("GET /images/{path...}", imageServer)` line 117; `adminMux.Handle("GET /images/{path...}", imageServer)` line 152 |
| `cmd/server/main.go` | `internal/handler/api/handler.go` | blogMux route + RequireAPIToken | WIRED | `blogMux.Handle("POST /api/push", requireToken(http.HandlerFunc(apiH.PushPost)))` at line 128 |
| `internal/handler/api/handler.go` | `internal/service/post/write.go` | svc.UpsertBySlug call | WIRED | `h.svc.UpsertBySlug(r.Context(), title, slug, string(raw))` at line 69 |
| `internal/middleware/apitoken.go` | `crypto/subtle` | ConstantTimeCompare | WIRED | `subtle.ConstantTimeCompare([]byte(got), []byte(token))` at line 31 |

### Behavioral Spot-Checks

| Behavior | Command | Result | Status |
|----------|---------|--------|--------|
| Project compiles | `go build ./...` | Exit 0 | PASS |
| All upload tests pass | `go test ./internal/handler/admin/... -run TestUploadImage` | 6/6 PASS | PASS |
| All token middleware tests pass | `go test ./internal/middleware/... -run TestRequireAPIToken` | 5/5 PASS | PASS |
| All upsert tests pass | `go test ./internal/service/post/... -run TestUpsertBySlug` | 3/3 PASS | PASS |
| All API handler tests pass | `go test ./internal/handler/api/...` | 6/6 PASS | PASS |
| Full test suite passes | `go test ./...` | All packages PASS | PASS |

### Requirements Coverage

| Requirement | Source Plan | Description | Status | Evidence |
|-------------|------------|-------------|--------|----------|
| ADMIN-08 | 05-01-PLAN | Admin can upload images and embed them in posts | SATISFIED | Upload handler validates MIME via magic bytes, generates random filenames, stores on EBS. Editor UI has upload button that inserts markdown tag. /images/ served on both muxes. 6 tests pass. |
| ADMIN-09 | 05-02-PLAN | Admin can push .md files via API endpoint with bearer token auth | SATISFIED | POST /api/push on blogMux with RequireAPIToken middleware. UpsertBySlug creates drafts or updates preserving tags. Constant-time token comparison. 14 tests pass across middleware, service, and handler. |

### Anti-Patterns Found

| File | Line | Pattern | Severity | Impact |
|------|------|---------|----------|--------|
| (none) | - | - | - | No anti-patterns detected in any phase files |

### Human Verification Required

### 1. Image Upload End-to-End

**Test:** Log into admin panel, create/edit a post, click "Upload Image", select a JPEG, verify the markdown tag appears at cursor position and the image renders in the preview pane.
**Expected:** Image uploads successfully, markdown `![alt](/images/xxx.jpg)` inserted at cursor, preview shows the image.
**Why human:** Requires visual confirmation of cursor position behavior and live preview rendering.

### 2. Spoofed File Rejection UX

**Test:** Attempt to upload an HTML file renamed to .jpg via the editor Upload Image button.
**Expected:** Error message "Only JPEG and PNG are accepted." appears next to the upload button; no file is stored.
**Why human:** Requires visual confirmation of error display styling and user experience.

### 3. API Push via curl

**Test:** Run `curl -H "Authorization: Bearer $TOKEN" -d @test-post.md https://jared-wallace.com/api/push` with a valid .md file containing YAML front matter with slug and title.
**Expected:** Returns "ok", post appears in admin dashboard as a draft.
**Why human:** Requires running server with database and valid API_TOKEN configured.

### Gaps Summary

No gaps found. All 4 success criteria from the roadmap are verified through code inspection and passing tests. All artifacts exist, are substantive (no stubs), and are properly wired. Both requirements (ADMIN-08, ADMIN-09) are satisfied. The full test suite passes with no regressions.

---

_Verified: 2026-03-28T01:00:00Z_
_Verifier: Claude (gsd-verifier)_
