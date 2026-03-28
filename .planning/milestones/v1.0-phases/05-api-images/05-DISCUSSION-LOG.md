# Phase 5: API + Images - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md -- this log preserves the alternatives considered.

**Date:** 2026-03-27
**Phase:** 05-api-images
**Areas discussed:** Image upload workflow, Image storage & serving, API push endpoint, Security boundaries

---

## Image Upload Workflow

### Upload UI Placement

| Option | Description | Selected |
|--------|-------------|----------|
| In the editor | Upload button in the post editor. On upload, markdown image tag auto-inserted at cursor position. | :white_check_mark: |
| Separate media page | Dedicated /admin/media page listing all uploaded images. Admin copies URL manually. | |
| Both | Editor inline upload plus separate media library page. | |

**User's choice:** In the editor
**Notes:** Keeps the writing flow unbroken.

### Image Formats

| Option | Description | Selected |
|--------|-------------|----------|
| JPEG + PNG only | Covers 99% of blog image needs. Matches success criteria. | :white_check_mark: |
| JPEG + PNG + WebP | Adds WebP for smaller file sizes. | |
| JPEG + PNG + WebP + GIF | Full format support including animated GIFs. | |

**User's choice:** JPEG + PNG only
**Notes:** None.

### Max File Size

| Option | Description | Selected |
|--------|-------------|----------|
| 5 MB | Generous for blog images, prevents accidental huge uploads. | :white_check_mark: |
| 10 MB | More headroom for high-res photography. | |
| No limit | Trust the admin, only natural constraint is disk space. | |

**User's choice:** 5 MB
**Notes:** None.

### Upload Trigger

| Option | Description | Selected |
|--------|-------------|----------|
| Upload button only | Simple button, file picker, upload, markdown tag inserted. | :white_check_mark: |
| Button + drag-and-drop | Button plus drag-drop onto textarea. | |
| Button + paste from clipboard | Button plus Ctrl+V image paste support. | |

**User's choice:** Upload button only
**Notes:** Consistent with Phase 3 "plain textarea, no toolbar" decision.

---

## Image Storage & Serving

### Disk Layout

| Option | Description | Selected |
|--------|-------------|----------|
| Flat directory | /var/www/html/images/{random-name}.jpg -- all in one folder. | :white_check_mark: |
| Year-month buckets | /var/www/html/images/2026/03/{random-name}.jpg -- organized by date. | |
| Per-post directories | /var/www/html/images/posts/{slug}/{random-name}.jpg -- grouped by post. | |

**User's choice:** Flat directory
**Notes:** Blog scale won't hit filesystem limits.

### Database Tracking

| Option | Description | Selected |
|--------|-------------|----------|
| No DB tracking | Images are just files on disk. Markdown body is the reference. | :white_check_mark: |
| Images table | Track filename, original name, upload date, post ID in Postgres. | |
| You decide | Claude picks. | |

**User's choice:** No DB tracking
**Notes:** No migration needed.

### URL Scheme

| Option | Description | Selected |
|--------|-------------|----------|
| /images/{filename} | Simple, readable, standard blog convention. | :white_check_mark: |
| /uploads/{filename} | Distinguishes uploads from static assets. | |
| /static/uploads/{filename} | Nested under /static/ with CSS/JS. | |

**User's choice:** /images/{filename}
**Notes:** None.

---

## API Push Endpoint

### Bearer Token Management

| Option | Description | Selected |
|--------|-------------|----------|
| Single env var | API_TOKEN env var, like ADMIN_EMAIL/ADMIN_PASSWORD_HASH. | :white_check_mark: |
| Generated + stored in DB | Admin generates tokens from panel, stored hashed in Postgres. | |
| You decide | Claude picks. | |

**User's choice:** Single env var
**Notes:** Consistent with existing credential pattern. Rotate by redeploying.

### Request Body Format

| Option | Description | Selected |
|--------|-------------|----------|
| Raw markdown body | Content-Type: text/markdown. Slug/title from YAML front matter. | :white_check_mark: |
| JSON envelope | JSON with slug, title, body, publish fields. | |
| Multipart form | Multipart with markdown file + optional metadata. | |

**User's choice:** Raw markdown body
**Notes:** Simplest CLI usage: curl -d @post.md.

### Upsert Behavior

| Option | Description | Selected |
|--------|-------------|----------|
| Upsert by slug | If slug exists, update; if not, create. Idempotent. | :white_check_mark: |
| Create only | Reject if slug exists. Prevents accidental overwrites. | |
| Separate endpoints | POST for create, PUT for update. RESTful but doubles surface. | |

**User's choice:** Upsert by slug
**Notes:** Matches success criteria.

### Default Publish State

| Option | Description | Selected |
|--------|-------------|----------|
| Draft by default | Pushed posts land as drafts. Publish from web UI. | :white_check_mark: |
| Front matter controlled | published: true/false in YAML front matter. | |
| Always published | Pushed posts go live immediately. | |

**User's choice:** Draft by default
**Notes:** Safer -- prevents accidental publishing.

---

## Security Boundaries

### Magic-Byte Validation

| Option | Description | Selected |
|--------|-------------|----------|
| http.DetectContentType | Go stdlib, reads first 512 bytes, sniffs MIME type. | :white_check_mark: |
| Manual header check | Read first 4-8 bytes, compare against known magic bytes. | |
| You decide | Claude picks. | |

**User's choice:** http.DetectContentType
**Notes:** Simple, no deps, meets success criteria.

### API Rate Limiting

| Option | Description | Selected |
|--------|-------------|----------|
| No rate limit | Bearer token is strong enough gate for single-admin blog. | :white_check_mark: |
| Basic rate limit | Reuse existing RateLimiter, e.g. 60 req/min per token. | |
| You decide | Claude decides based on threat model. | |

**User's choice:** No rate limit
**Notes:** None.

### Upload Timeout Handling

| Option | Description | Selected |
|--------|-------------|----------|
| Bump WriteTimeout for upload route | Extend to 30s on upload handler only. | :white_check_mark: |
| Keep current timeouts | 10s WriteTimeout is fine for 5MB on modern connections. | |
| You decide | Claude evaluates. | |

**User's choice:** Bump WriteTimeout for upload route
**Notes:** Per-handler timeout override, other routes keep 10s.

---

## Claude's Discretion

- Random filename generation strategy
- /images/ file server wiring
- Config struct extension
- Front matter schema for push endpoint
- Upload button placement and styling
- Error response format for API
- Upload endpoint response (markdown tag vs just URL)
- ReadTimeout adjustment for upload route

## Deferred Ideas

None -- discussion stayed within phase scope.
