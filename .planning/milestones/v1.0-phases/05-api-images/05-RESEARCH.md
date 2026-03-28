# Phase 5: API + Images - Research

**Researched:** 2026-03-27
**Domain:** Go file upload, MIME validation, bearer token auth, static file serving, per-handler timeouts
**Confidence:** HIGH

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

**Image Upload Workflow**
- D-01: Upload lives inside the post editor — a button above/beside the textarea. On upload, a markdown image tag (`![alt](/images/xyz.jpg)`) is auto-inserted at cursor position.
- D-02: Upload button only — no drag-and-drop, no clipboard paste.
- D-03: Accepted formats: JPEG and PNG only. Validated via magic bytes.
- D-04: Max file size: 5 MB per upload.

**Image Storage & Serving**
- D-05: Flat directory on EBS: `/var/www/html/images/{random-name}.{ext}`. No date or post bucketing.
- D-06: No database tracking of images. Markdown in post body IS the reference.
- D-07: URL scheme: `/images/{filename}`. Go serves files from the EBS images directory.

**API Push Endpoint**
- D-08: Bearer token managed as a single env var (`API_TOKEN`).
- D-09: Request body is raw markdown with `Content-Type: text/markdown`. Slug and title extracted from YAML front matter via goldmark-meta.
- D-10: Upsert by slug — if slug exists, update body/title; if not, create new post.
- D-11: Pushed posts land as drafts by default.

**Security Boundaries**
- D-12: Magic-byte validation via `net/http.DetectContentType` (reads first 512 bytes). Checks for `image/jpeg` and `image/png`.
- D-13: No rate limiting on the API push endpoint.
- D-14: Bump WriteTimeout to 30s on the image upload handler only via `http.ResponseController.SetWriteDeadline`.

### Claude's Discretion
- Random filename generation strategy (UUID, crypto/rand hex, etc.)
- How to wire the `/images/` file server into the existing mux
- Config struct extension for `API_TOKEN` and image directory path
- Front matter schema for the push endpoint (which YAML keys are required vs optional)
- Upload button placement and styling within the editor template
- Error response format for the API endpoint (plain text vs JSON)
- Whether the upload endpoint returns the inserted markdown tag or just the URL
- ReadTimeout adjustment for upload route (if needed alongside WriteTimeout bump)

### Deferred Ideas (OUT OF SCOPE)
None — discussion stayed within phase scope.
</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| ADMIN-08 | Admin can upload images and embed them in posts | File upload via multipart/form-data, MIME magic-byte check, EBS flat-dir storage, crypto/rand hex filename, `/images/` static file server, JS auto-insert markdown tag |
| ADMIN-09 | Admin can push .md files via API endpoint with bearer token auth | `POST /api/push`, Authorization header parsing, goldmark-meta front matter extraction, service.UpsertBySlug, 401 on missing/invalid token |
</phase_requirements>

---

## Summary

Phase 5 adds two independent capabilities to an already-working Go blog. Both are self-contained additions with clear integration points. The codebase has established patterns for everything needed: env-var config, session middleware (for the auth model), service layer, and vanilla JS in admin.js.

Image upload is a standard Go multipart form handler. The key correctness concerns are: (1) reading the first 512 bytes for magic-byte MIME detection before accepting the file, (2) generating a server-controlled random filename so client-supplied names are never used, and (3) writing to the EBS path that is guaranteed to exist at `/var/www/html/images/`. The WriteTimeout override via `http.ResponseController.SetWriteDeadline` handles the 5 MB upload case without touching the global server timeout.

The API push endpoint is a thin handler: parse `Authorization: Bearer <token>`, constant-time compare against `cfg.APIToken`, read raw body, extract front matter via the existing `markdown.RenderWithMeta`, and call an `UpsertBySlug` service method. The service method is the only new business logic — a conditional Create or Update based on slug existence.

**Primary recommendation:** Implement image upload and API push as separate, focused additions. No new packages required — every dependency is already in go.mod. The upsert logic is the only non-trivial piece: use a try-insert-then-update approach via the repository layer.

---

## Standard Stack

All dependencies are **already in go.mod**. No new packages are needed.

### Core (existing, confirmed in go.mod)
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| `net/http` stdlib | (Go 1.26) | Multipart upload parsing, static file server, timeout control | `ParseMultipartForm`, `http.FileServer`, `http.ResponseController` are all stdlib |
| `crypto/rand` stdlib | (Go 1.26) | Cryptographically secure random filename bytes | Standard pattern; no UUID library needed for 16-byte hex name |
| `encoding/hex` stdlib | (Go 1.26) | Encode random bytes as hex string for filename | Cleaner than fmt.Sprintf("%x", b) for this use case |
| `github.com/yuin/goldmark-meta` | v1.1.0 | Extract YAML front matter (slug, title) from pushed .md | Already in go.mod; `RenderWithMeta` already exists in `internal/markdown/renderer.go` |
| `strings` stdlib | (Go 1.26) | Authorization header parsing (`strings.TrimPrefix`) | No helper package needed |

### No New Dependencies Required
The complete go.mod already contains everything Phase 5 needs. Adding zero new entries to go.mod is the correct outcome.

---

## Architecture Patterns

### New Files and Modification Points

```
cmd/server/main.go                     # Add /images/ file server, POST /api/push route, upload route
internal/config/config.go             # Add APIToken, ImageDir fields
internal/handler/admin/upload.go      # New: UploadImage handler
internal/handler/api/handler.go       # New: API handler package (or add to admin package — see below)
internal/service/post/write.go        # Add UpsertBySlug method
web/templates/admin-editor.html       # Add upload button HTML
web/static/admin.js                   # Add upload fetch + cursor-insert logic
```

### Pattern 1: Bearer Token Middleware

Follow the same shape as `middleware.RequireSession`. The middleware reads the `Authorization` header, strips the `Bearer ` prefix, and does a constant-time comparison against `cfg.APIToken`. Returns 401 on mismatch.

```go
// Source: mirrors internal/middleware/auth.go pattern
func RequireAPIToken(token string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            raw := r.Header.Get("Authorization")
            got := strings.TrimPrefix(raw, "Bearer ")
            // subtle.ConstantTimeCompare prevents timing attacks
            if subtle.ConstantTimeCompare([]byte(got), []byte(token)) != 1 {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

Use `crypto/subtle.ConstantTimeCompare` — it is already an indirect dependency via `golang.org/x/crypto`. Import it directly from `crypto/subtle` (stdlib, no new dep).

### Pattern 2: Image Upload Handler

```go
// Source: net/http stdlib — ParseMultipartForm + DetectContentType
func (h *AdminHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
    // 1. Extend write deadline for this handler only (D-14)
    rc := http.NewResponseController(w)
    _ = rc.SetWriteDeadline(time.Now().Add(30 * time.Second))

    // 2. Parse multipart — 5 MB max in memory (D-04)
    const maxBytes = 5 << 20 // 5 MB
    if err := r.ParseMultipartForm(maxBytes); err != nil {
        http.Error(w, "file too large or bad request", http.StatusBadRequest)
        return
    }

    file, header, err := r.FormFile("image")
    if err != nil { /* ... */ }
    defer file.Close()

    // 3. Read first 512 bytes for MIME sniff (D-12)
    buf := make([]byte, 512)
    n, err := file.Read(buf)
    if err != nil && err != io.EOF { /* ... */ }
    mime := http.DetectContentType(buf[:n])
    if mime != "image/jpeg" && mime != "image/png" {
        http.Error(w, "only JPEG and PNG accepted", http.StatusUnsupportedMediaType)
        return
    }

    // 4. Seek back to start (multipart.File implements io.Seeker)
    if _, err := file.Seek(0, io.SeekStart); err != nil { /* ... */ }

    // 5. Generate random filename (D-05, discretion)
    ext := extensionFromMIME(mime) // ".jpg" or ".png"
    b := make([]byte, 16)
    if _, err := cryptorand.Read(b); err != nil { /* ... */ }
    filename := hex.EncodeToString(b) + ext

    // 6. Write to EBS path
    dst, err := os.Create(filepath.Join(h.imageDir, filename))
    if err != nil { /* ... */ }
    defer dst.Close()
    if _, err := io.Copy(dst, file); err != nil { /* ... */ }

    // 7. Return markdown tag or URL (discretion: return JSON with url + markdownTag)
    url := "/images/" + filename
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "url":         url,
        "markdownTag": fmt.Sprintf("![alt](%s)", url),
    })
}
```

**Key detail:** `multipart.File` implements `io.ReadSeeker`. After reading 512 bytes for MIME sniff, seek back to 0 before copying to disk. Forgetting this produces a truncated file — a common pitfall.

### Pattern 3: Static Image File Server (/images/)

Wire directly in `main.go` on both muxes (blog needs to serve images too, since posts embed `/images/...` URLs). The images directory lives on EBS, NOT embedded in the binary.

```go
// Source: net/http stdlib — http.FileServer + http.StripPrefix
// In main.go, BEFORE building hostRouter:
imageDir := cfg.ImageDir // e.g. "/var/www/html/images"
imageServer := http.StripPrefix("/images/", http.FileServer(http.Dir(imageDir)))

blogMux.Handle("GET /images/{path...}", imageServer)
adminMux.Handle("GET /images/{path...}", imageServer)
```

Use `http.Dir` (not `http.FileServerFS`) because the images directory is on the EBS volume, not embedded. `http.FileServer(http.Dir(...))` serves from a real filesystem path.

**Security:** `http.FileServer` handles path traversal prevention internally (it calls `path.Clean` and rejects `..`). No custom sanitization needed. (HIGH confidence — Go stdlib behavior.)

### Pattern 4: API Push Endpoint + UpsertBySlug

```go
// internal/handler/api/handler.go (new package)
func (h *APIHandler) PushPost(w http.ResponseWriter, r *http.Request) {
    body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20)) // 1 MB limit for .md files
    if err != nil { /* ... */ }

    _, meta := h.renderer.RenderWithMeta(string(body))
    slug, _ := meta["slug"].(string)
    title, _ := meta["title"].(string)
    if slug == "" {
        http.Error(w, "front matter must include 'slug'", http.StatusBadRequest)
        return
    }
    if title == "" {
        title = slug // fallback
    }

    if err := h.svc.UpsertBySlug(r.Context(), title, slug, string(body)); err != nil {
        http.Error(w, "failed to upsert post", http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "ok")
}
```

**UpsertBySlug service method:** Call `repo.FindBySlug` — if found, call `repo.Update`; if `ErrNotFound`, call `repo.Create` with `published=false` (D-11). The service layer owns the conditional; the handler stays thin.

### Pattern 5: JS Upload + Cursor Insert

Follow the existing admin.js IIFE pattern. The upload fetch calls `POST /admin/images/upload` with `FormData`. On success, insert the returned `markdownTag` at the textarea cursor position using `selectionStart`/`selectionEnd`.

```javascript
// Source: existing admin.js pattern (vanilla JS, IIFE)
function insertAtCursor(textarea, text) {
  var start = textarea.selectionStart;
  var end = textarea.selectionEnd;
  textarea.value = textarea.value.slice(0, start) + text + textarea.value.slice(end);
  textarea.selectionStart = textarea.selectionEnd = start + text.length;
  textarea.dispatchEvent(new Event('input')); // triggers live preview debounce
}
```

### Where Does the API Handler Live?

Two valid options:

**Option A:** New `internal/handler/api/` package — mirrors the existing `handler/blog/` and `handler/admin/` pattern. Cleanest separation. Recommended.

**Option B:** Add `PushPost` directly to `internal/handler/admin/handler.go` — saves one package but mixes session-auth and token-auth handlers. Avoid.

Use Option A.

### Anti-Patterns to Avoid

- **Using client-supplied filename:** `header.Filename` is attacker-controlled. Never use it for the stored file path or URL. Generate server-side always.
- **MIME check from Content-Type header:** The client sets this header. Use `http.DetectContentType` on actual bytes.
- **Serving images from embedded FS:** Images are on EBS, not compiled into the binary. Use `http.Dir`, not `fs.FS`.
- **Skipping the Seek after MIME sniff:** Reading 512 bytes advances the file cursor. Failing to `Seek(0, io.SeekStart)` before `io.Copy` produces a file missing its first 512 bytes.
- **Using `==` for token comparison:** Use `crypto/subtle.ConstantTimeCompare` to prevent timing oracle attacks.
- **Placing `/api/push` behind admin host check:** The API endpoint should be reachable at the public domain too (or at minimum, not host-restricted). The bearer token is the gate, not the host.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| MIME type detection | Custom byte comparison | `net/http.DetectContentType` | WhatWG MIME sniff spec, handles edge cases, stdlib |
| Random filename | Math.random or timestamp | `crypto/rand` + `encoding/hex` | Cryptographically secure, collision-resistant |
| Path traversal protection | Custom path cleaning | `http.FileServer` built-in | `path.Clean` already applied by stdlib |
| Constant-time string compare | `==` operator | `crypto/subtle.ConstantTimeCompare` | Prevents timing attacks on token comparison |
| Front matter parsing | Custom YAML parser | `goldmark-meta` (already in dep tree) | Already wired in `renderer.go` as `RenderWithMeta` |

---

## Common Pitfalls

### Pitfall 1: Forgetting Seek After MIME Sniff
**What goes wrong:** File stored on disk is truncated by 512 bytes; images render as corrupted.
**Why it happens:** `file.Read(buf)` advances the cursor. `io.Copy` starts from the current position.
**How to avoid:** Call `file.Seek(0, io.SeekStart)` immediately after MIME detection, before `io.Copy`.
**Warning signs:** Images open as partially-corrupt or appear smaller than uploaded size.

### Pitfall 2: CSRF Protection on Upload Endpoint
**What goes wrong:** The existing `cop.Handler(adminMux)` applies `http.CrossOriginProtection` to all admin routes. The upload endpoint is a POST — it will be covered.
**Why it happens:** The upload uses `fetch` with a `FormData` body, same as the preview endpoint. The COP header check passes when the JS `fetch` call is same-origin.
**How to avoid:** No action needed — the existing COP wrapping handles this correctly for same-origin JS. Do NOT add `application/json` Content-Type to the fetch upload call; `multipart/form-data` is correct and COP allows it from the same origin.
**Warning signs:** 403 errors from upload fetch; inspect `Origin` header in request.

### Pitfall 3: ImageDir Not Created on First Boot
**What goes wrong:** The upload handler calls `os.Create(filepath.Join(imageDir, filename))` and fails if `/var/www/html/images` does not exist.
**Why it happens:** EBS volume is mounted but subdirectory is not pre-created.
**How to avoid:** Call `os.MkdirAll(cfg.ImageDir, 0755)` during startup (in `main.go` after config load), before starting the server. Idempotent — safe to call if directory already exists.
**Warning signs:** Upload returns 500; server logs show `no such file or directory`.

### Pitfall 4: API Push Endpoint Host Routing
**What goes wrong:** `POST /api/push` placed only on `adminMux` — only reachable from `admin.jared-wallace.com`, which is not what a CLI push client expects.
**Why it happens:** Reflexive assumption that "API" belongs on the admin host.
**How to avoid:** Add the route to `blogMux` (public domain) OR create a third mux for API routes, OR handle it in the `hostRouter.ServeHTTP` before dispatching to admin/blog. Simpler: add to `blogMux` — it has its own bearer-token guard and doesn't need session infrastructure.
**Warning signs:** `curl -X POST https://jared-wallace.com/api/push` returns 404.

### Pitfall 5: WriteTimeout vs. ResponseController Ordering
**What goes wrong:** Server-wide `WriteTimeout` (10s) fires before `SetWriteDeadline` (30s) takes effect.
**Why it happens:** `http.ResponseController.SetWriteDeadline` EXTENDS the deadline from the connection's perspective, but if the server-wide `WriteTimeout` has already elapsed, there is nothing to extend.
**How to avoid:** Call `rc.SetWriteDeadline(time.Now().Add(30 * time.Second))` at the very start of the upload handler, before any blocking I/O. The existing server WriteTimeout is 10s — for 5 MB on a slow connection this may still be tight. Consider raising the server-wide WriteTimeout to 30s and accepting the trade-off, OR keep it at 10s and rely on the per-handler extension (Go docs confirm this works when called before the deadline is exceeded).
**Warning signs:** Large uploads return 504 or connection reset; client sees timeout on slow network.

### Pitfall 6: goldmark-meta Key Casing
**What goes wrong:** Front matter key `Slug` (capitalized) works in some parsers but `meta["slug"]` returns nil in goldmark-meta.
**Why it happens:** goldmark-meta uses the YAML key exactly as written. YAML keys are case-sensitive.
**How to avoid:** Document the required front matter schema: keys must be lowercase (`slug`, `title`). Enforce in handler: if `meta["slug"]` is nil, return 400.
**Warning signs:** `UpsertBySlug` called with empty slug, logs show "front matter must include slug" on valid-looking files.

---

## Code Examples

### Random Hex Filename Generation
```go
// Source: crypto/rand + encoding/hex stdlib
import (
    cryptorand "crypto/rand"
    "encoding/hex"
)

func randomFilename(ext string) (string, error) {
    b := make([]byte, 16) // 16 bytes = 32 hex chars
    if _, err := cryptorand.Read(b); err != nil {
        return "", err
    }
    return hex.EncodeToString(b) + ext, nil
}
```

### MIME Extension Helper
```go
// extensionFromMIME maps detected MIME type to file extension.
// Only called after magic-byte validation, so only two cases needed.
func extensionFromMIME(mime string) string {
    switch mime {
    case "image/jpeg":
        return ".jpg"
    case "image/png":
        return ".png"
    default:
        return ""
    }
}
```

### Config Extension Pattern
```go
// Source: matches existing internal/config/config.go pattern
type Config struct {
    // ... existing fields ...
    APIToken string // envOr("API_TOKEN", "")
    ImageDir string // envOr("IMAGE_DIR", "/var/www/html/images")
}

// In Load():
APIToken: envOr("API_TOKEN", ""),
ImageDir: envOr("IMAGE_DIR", "/var/www/html/images"),
```

### UpsertBySlug Service Method
```go
// Source: mirrors existing Create/Update pattern in internal/service/post/write.go
func (s *Service) UpsertBySlug(ctx context.Context, title, slug, body string) error {
    rendered := s.renderer.Render(body)
    existing, err := s.repo.FindBySlug(ctx, slug)
    if err != nil {
        if errors.Is(err, postrepo.ErrNotFound) {
            // Create as draft (D-11)
            p := model.Post{
                Title:        title,
                Slug:         slug,
                Body:         body,
                RenderedHTML: string(rendered),
                Published:    false,
            }
            _, err = s.repo.Create(ctx, p)
            return err
        }
        return err
    }
    // Update existing
    p := model.Post{
        ID:           existing.ID,
        Title:        title,
        Slug:         slug,
        Body:         body,
        RenderedHTML: string(rendered),
        Tags:         existing.Tags, // preserve tags on push
    }
    return s.repo.Update(ctx, p)
}
```

### Required Front Matter Schema
```yaml
---
slug: my-post-slug      # required — used as upsert key
title: My Post Title    # required — stored as post title
---
```

Optional fields (ignored by push endpoint but valid in goldmark-meta):
```yaml
tags: go, web           # not extracted by push; admin sets tags in UI
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| gorilla/csrf for CSRF | `http.CrossOriginProtection` | Go 1.26 stdlib | No external dep needed; already wired in main.go |
| Server-wide timeout only | Per-handler `ResponseController.SetWriteDeadline` | Go 1.20 | Upload handler gets 30s; rest stays at 10s |
| `lib/pq` driver | `pgx/v5` | Already in project | No impact for Phase 5 — no new DB work |

---

## Environment Availability

No new external dependencies. The EBS volume path `/var/www/html/images` does not exist yet (it is a runtime directory, not a source artifact). Handled by `os.MkdirAll` at startup.

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| EBS volume at /var/www/html | Image storage | Runtime-only | — | Dev: use local path via IMAGE_DIR env var |
| `crypto/rand` | Random filename | stdlib | Go 1.26 | — |
| `net/http.DetectContentType` | MIME validation | stdlib | Go 1.26 | — |
| `http.ResponseController` | Per-handler timeout | stdlib | Go 1.20+ | — |

---

## Validation Architecture

### Test Framework
| Property | Value |
|----------|-------|
| Framework | Go stdlib `testing` |
| Config file | none (standard `go test ./...`) |
| Quick run command | `go test ./internal/handler/admin/... ./internal/handler/api/... ./internal/middleware/... ./internal/service/post/... -run TestUpload\|TestPush\|TestBearer` |
| Full suite command | `go test ./...` |

### Phase Requirements → Test Map
| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| ADMIN-08 | Valid JPEG upload returns 200 + markdown tag | unit | `go test ./internal/handler/admin/... -run TestUploadImage_ValidJPEG` | ❌ Wave 0 |
| ADMIN-08 | Spoofed MIME (HTML with .jpg ext) rejected with 415 | unit | `go test ./internal/handler/admin/... -run TestUploadImage_SpoofedMIME` | ❌ Wave 0 |
| ADMIN-08 | Oversized file (>5 MB) rejected with 400 | unit | `go test ./internal/handler/admin/... -run TestUploadImage_TooLarge` | ❌ Wave 0 |
| ADMIN-08 | Server-generated filename used, not client filename | unit | `go test ./internal/handler/admin/... -run TestUploadImage_RandomFilename` | ❌ Wave 0 |
| ADMIN-09 | Valid bearer token + .md body creates draft post | unit | `go test ./internal/handler/api/... -run TestPushPost_ValidToken` | ❌ Wave 0 |
| ADMIN-09 | Missing token returns 401, no post created | unit | `go test ./internal/handler/api/... -run TestPushPost_NoToken` | ❌ Wave 0 |
| ADMIN-09 | Invalid token returns 401, no post created | unit | `go test ./internal/handler/api/... -run TestPushPost_InvalidToken` | ❌ Wave 0 |
| ADMIN-09 | Repeated push of same slug upserts (no duplicate) | unit | `go test ./internal/service/post/... -run TestUpsertBySlug` | ❌ Wave 0 |
| ADMIN-09 | Missing slug in front matter returns 400 | unit | `go test ./internal/handler/api/... -run TestPushPost_NoSlug` | ❌ Wave 0 |

### Sampling Rate
- **Per task commit:** `go test ./internal/handler/admin/... ./internal/handler/api/... ./internal/service/post/...`
- **Per wave merge:** `go test ./...`
- **Phase gate:** Full suite green before `/gsd:verify-work`

### Wave 0 Gaps
- [ ] `internal/handler/admin/upload_test.go` — covers ADMIN-08 upload cases
- [ ] `internal/handler/api/handler_test.go` — covers ADMIN-09 push cases
- [ ] `internal/service/post/upsert_test.go` — covers UpsertBySlug logic

---

## Open Questions

1. **Upload endpoint on which host?**
   - What we know: Upload is initiated from admin editor JS. Admin editor lives on `admin.jared-wallace.com`.
   - What's clear: Upload handler belongs on `adminMux`, behind `RequireSession`. No ambiguity.
   - Resolution: `POST /admin/images/upload` on `adminMux`, wrapped with `requireAuth`.

2. **API push endpoint on which host?**
   - What we know: `POST /api/push` is called from a CLI on a local machine, not a browser. The user's machine points to `jared-wallace.com` (the public domain).
   - Recommendation: Add to `blogMux` with `RequireAPIToken` middleware. This way `curl https://jared-wallace.com/api/push` works naturally.

3. **Tags on push?**
   - What we know: D-09 says slug and title are extracted from front matter. Tags are not mentioned.
   - Recommendation: Preserve existing tags on upsert (shown in `UpsertBySlug` example above). Admin sets tags from the UI. Document this in the API spec comment.

---

## Sources

### Primary (HIGH confidence)
- Go stdlib `net/http` — `DetectContentType`, `ParseMultipartForm`, `FileServer`, `ResponseController` — all verified in stdlib docs
- Go stdlib `crypto/rand` + `encoding/hex` — verified at [pkg.go.dev/crypto/rand](https://pkg.go.dev/crypto/rand)
- `github.com/yuin/goldmark-meta` v1.1.0 — already in go.mod; `RenderWithMeta` already implemented in `internal/markdown/renderer.go`
- Codebase: `internal/config/config.go`, `internal/middleware/auth.go`, `internal/handler/admin/handler.go`, `internal/service/post/write.go` — all read directly

### Secondary (MEDIUM confidence)
- [How to Use the http.ResponseController Type — Alex Edwards](https://www.alexedwards.net/blog/how-to-use-the-http-responsecontroller-type) — SetWriteDeadline per-handler pattern
- [How to detect the Content Type of a file in Go — freshman.tech](https://freshman.tech/snippets/go/file-content-type/) — DetectContentType usage
- [Go crypto/rand example — go.dev](https://go.dev/src/crypto/rand/example_test.go) — random bytes pattern

### Tertiary (LOW confidence — not needed; all findings HIGH/MEDIUM)
None.

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — zero new dependencies; all verified in existing go.mod and stdlib
- Architecture: HIGH — directly derived from reading the actual codebase
- Pitfalls: HIGH — MIME seek pitfall and CSRF interaction are well-known Go patterns; host routing pitfall derived from reading main.go directly

**Research date:** 2026-03-27
**Valid until:** 2026-05-01 (stable stdlib patterns; no fast-moving libraries)
