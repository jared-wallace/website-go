---
phase: 05-api-images
plan: 01
subsystem: api
tags: [image-upload, multipart, mime-validation, file-server, admin-editor]

# Dependency graph
requires:
  - phase: 03-admin-panel
    provides: AdminHandler struct, editor template, admin.js IIFE, requireAuth middleware
provides:
  - UploadImage handler with magic-byte MIME validation and random filenames
  - ImageDir and APIToken config fields
  - /images/ static file server on blog and admin muxes
  - Upload Image button and insertAtCursor JS in admin editor
affects: [05-api-images plan 02, 06-docker]

# Tech tracking
tech-stack:
  added: []
  patterns: [MaxBytesReader for upload size enforcement, crypto/rand hex filenames, http.DetectContentType magic-byte validation]

key-files:
  created:
    - internal/handler/admin/upload.go
    - internal/handler/admin/upload_test.go
  modified:
    - internal/config/config.go
    - internal/handler/admin/handler.go
    - web/templates/admin-editor.html
    - web/static/admin.js
    - cmd/server/main.go

key-decisions:
  - "MaxBytesReader used instead of ParseMultipartForm maxMemory for hard 5MB limit enforcement"
  - "APIToken config field added alongside ImageDir to avoid re-modifying config.go in Plan 02"
  - "SetImageDir test helper method on AdminHandler for t.TempDir injection"

patterns-established:
  - "Upload handler pattern: MaxBytesReader -> ParseMultipartForm -> FormFile -> DetectContentType sniff -> Seek(0) -> crypto/rand filename -> os.Create"
  - "Image file server on both muxes via shared imageServer variable"

requirements-completed: [ADMIN-08]

# Metrics
duration: 4min
completed: 2026-03-28
---

# Phase 5 Plan 1: Image Upload Summary

**Admin image upload with magic-byte MIME validation, random hex filenames, EBS storage, and editor UI integration**

## Performance

- **Duration:** 4 min
- **Started:** 2026-03-28T00:26:14Z
- **Completed:** 2026-03-28T00:31:09Z
- **Tasks:** 2
- **Files modified:** 7

## Accomplishments
- Upload handler validates JPEG/PNG via magic bytes, rejects spoofed MIME with 415 and oversized with 400
- Server-generated 32-char hex filenames (never trusts client filename)
- /images/ file server wired on both blog and admin muxes for public posts and admin preview
- Editor has Upload Image button that triggers file picker and auto-inserts markdown image tag at cursor
- 6 unit tests covering all upload edge cases

## Task Commits

Each task was committed atomically:

1. **Task 1: Upload handler with MIME validation and tests** - `914dead` (test: RED), `bc6b891` (feat: GREEN)
2. **Task 2: Editor UI, upload JS, image file server, and main.go wiring** - `da18faa` (feat)

_Note: Task 1 was TDD with separate RED and GREEN commits._

## Files Created/Modified
- `internal/handler/admin/upload.go` - UploadImage handler with MIME validation, random filename, JSON response
- `internal/handler/admin/upload_test.go` - 6 tests: ValidJPEG, ValidPNG, SpoofedMIME, TooLarge, RandomFilename, NoFile
- `internal/config/config.go` - Added ImageDir and APIToken fields
- `internal/handler/admin/handler.go` - Added imageDir field, SetImageDir test helper
- `web/templates/admin-editor.html` - Upload Image button, hidden file input, error span
- `web/static/admin.js` - Upload fetch logic, insertAtCursor function
- `cmd/server/main.go` - os.MkdirAll, imageServer, /images/ routes, upload route

## Decisions Made
- Used http.MaxBytesReader for hard 5MB enforcement (ParseMultipartForm maxMemory only controls memory vs temp-file threshold, not total size)
- Added APIToken config field preemptively for Plan 02 to avoid re-modifying config.go
- Full 8-byte PNG signature needed in tests for http.DetectContentType recognition

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Fixed upload size enforcement mechanism**
- **Found during:** Task 1 (GREEN phase)
- **Issue:** Plan specified `r.ParseMultipartForm(5<<20)` for size limiting, but this only sets the memory threshold for multipart parsing -- it does not reject oversized files
- **Fix:** Added `http.MaxBytesReader(w, r.Body, maxUpload+512)` before ParseMultipartForm to enforce a hard size limit
- **Files modified:** internal/handler/admin/upload.go
- **Verification:** TestUploadImage_TooLarge passes with 6MB file rejected at 400
- **Committed in:** bc6b891

**2. [Rule 1 - Bug] Fixed PNG test magic bytes**
- **Found during:** Task 1 (GREEN phase)
- **Issue:** Test used 4-byte PNG header `{0x89, 0x50, 0x4E, 0x47}` but http.DetectContentType requires the full 8-byte PNG signature to identify image/png
- **Fix:** Extended to full signature `{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}`
- **Files modified:** internal/handler/admin/upload_test.go
- **Verification:** TestUploadImage_ValidPNG passes
- **Committed in:** bc6b891

---

**Total deviations:** 2 auto-fixed (2 bugs)
**Impact on plan:** Both fixes necessary for correctness. No scope creep.

## Issues Encountered
None beyond the auto-fixed deviations above.

## User Setup Required
None - no external service configuration required. IMAGE_DIR defaults to /var/www/html/images.

## Next Phase Readiness
- Upload handler ready; Plan 02 (API push endpoint) can wire APIToken config field already in place
- Docker phase (06) should ensure /var/www/html/images directory exists on EBS volume

---
*Phase: 05-api-images*
*Completed: 2026-03-28*
