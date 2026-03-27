---
phase: 03-admin-panel
plan: "03"
subsystem: dashboard
tags: [dashboard, post-management, filter-tabs, status-pills, soft-delete, publish, templates, css]
dependency_graph:
  requires:
    - internal/handler/admin/handler.go
    - internal/service/post/write.go
    - internal/model/post.go
    - web/templates/admin-base.html
    - web/static/admin.css
  provides:
    - internal/handler/admin/dashboard.go
    - internal/handler/admin/actions.go
    - internal/handler/admin/actions_test.go
    - web/templates/admin-dashboard.html
    - web/static/admin.css (extended)
  affects:
    - internal/handler/admin/handler.go (stubs removed)
tech_stack:
  added: []
  patterns:
    - In-memory filter over ListAll result — avoids N+1 queries; post counts small at blog scale
    - PopString for flash messages — SCS PopString clears flash after first read (Pitfall 8)
    - req.SetPathValue for httptest — Go 1.22 net/http pattern for injecting path values in tests
    - successRepo mock pattern — separate mock with all write ops returning nil for action handler tests
    - Native <dialog> element — no JS dependency for confirmation modal; showModal() / close() browser APIs
key_files:
  created:
    - internal/handler/admin/dashboard.go
    - internal/handler/admin/actions.go
    - internal/handler/admin/actions_test.go
  modified:
    - internal/handler/admin/handler.go
    - web/templates/admin-dashboard.html
    - web/static/admin.css
decisions:
  - "In-memory filter over ListAll chosen over separate DB queries per tab — blog scale makes this simpler and eliminates 4 extra repository methods"
  - "successRepo separate from mockRepository in handler_test.go — action tests need write ops to succeed; auth tests need them to fail. Clean separation avoids conditional logic in mock"
  - "req.SetPathValue used in httptest to simulate Go 1.22 ServeMux path extraction"
metrics:
  duration: "10 minutes"
  completed: "2026-03-26"
  tasks: 2
  files_changed: 6
---

# Phase 03 Plan 03: Admin Dashboard Summary

**One-liner:** Post management dashboard with in-memory filter tabs (All/Published/Drafts/Deleted), action handlers for soft-delete/restore/publish/unpublish with flash messages, native dialog confirmation, and 10 passing handler tests.

## Tasks Completed

| # | Name | Commit | Files |
|---|------|--------|-------|
| 1 | Dashboard handler, post action handlers, and action handler tests | 1d2ae03 | dashboard.go, actions.go, actions_test.go, handler.go |
| 2 | Dashboard template and CSS styles | 701d5c3 | admin-dashboard.html, admin.css |

## What Was Built

### Dashboard Handler (`internal/handler/admin/dashboard.go`)

`Dashboard(w, r)` reads `filter` from query string, calls `h.svc.ListAll`, and filters the result in-memory:
- `""` / `"all"`: posts where `DeletedAt == nil` (excludes deleted)
- `"published"`: `DeletedAt == nil AND Published == true`
- `"drafts"`: `DeletedAt == nil AND Published == false`
- `"deleted"`: `DeletedAt != nil`

Reads flash messages via `PopString` (per Pitfall 8 — clears on read), renders `admin-dashboard.html` with `Posts`, `Filter`, `FlashSuccess`, `FlashError`.

### Action Handlers (`internal/handler/admin/actions.go`)

All four handlers follow the same pattern: parse `{id}` via `r.PathValue("id")`, call service method, set flash, redirect 303.

| Handler | Service Call | Flash Message | Redirect |
|---------|-------------|---------------|----------|
| `DeletePost` | `svc.SoftDelete` | "Post deleted. Restore it from the Deleted tab." | `/admin/posts?filter=deleted` |
| `RestorePost` | `svc.Restore` | "Post restored to drafts." | `/admin/posts?filter=drafts` |
| `PublishPost` | `svc.Publish` | "Post published." | `/admin/posts` |
| `UnpublishPost` | `svc.Unpublish` | "Post unpublished." | `/admin/posts` |

Error path (invalid ID or service failure): sets `flash_error` "Something went wrong. Try again." and redirects to `/admin/posts`.

### handler.go cleanup

Removed stub implementations for Dashboard, DeletePost, RestorePost, PublishPost, UnpublishPost. Kept NewPost, EditPost, SavePost, Preview stubs for Plan 03-04.

### Dashboard Template (`web/templates/admin-dashboard.html`)

Full replacement of the Plan 02 stub. Implements:
- **Filter tabs**: All / Published / Drafts / Deleted with `active` class on current tab
- **Post table**: Title (focal point — Playfair Display 18px/700), Status pill, Date, Actions columns
- **Status pills**: `.status-published` (accent bg), `.status-draft` (border only), `.status-deleted` (danger border)
- **Action links**: Edit (link), Publish/Unpublish (form POST), Delete button (opens dialog)
- **Native `<dialog>`**: `showModal()` / `close()` — no JS dependency for confirmation
- **Empty states**: Per-tab copy — "No published posts.", "No drafts.", "Nothing deleted.", "Nothing here yet."
- **New Post CTA**: Top-right `.cta-button` linking to `/admin/posts/new`

### admin.css extensions

Added 170 lines of dashboard styles: `.admin-dashboard` layout, `.filter-tabs` with active underline, `.admin-table` with column widths, `.post-title-link` focal point, `.status-pill` variants, `.action-link` with danger variant, `.inline-form`, `.confirm-dialog` with backdrop, `.cta-button--secondary` and `.cta-button--danger` button variants, `.empty-state`, responsive media query hiding date column on mobile.

## Test Coverage

10 tests total (6 new + 4 existing auth tests):

| Test | What It Verifies |
|------|-----------------|
| `TestDeletePost` | 303 redirect to `/admin/posts?filter=deleted` |
| `TestRestorePost` | 303 redirect to `/admin/posts?filter=drafts` |
| `TestPublishPost` | 303 redirect to `/admin/posts` |
| `TestUnpublishPost` | 303 redirect to `/admin/posts` |
| `TestDeletePostInvalidID` | Invalid ID → 303 redirect to `/admin/posts` |
| `TestDashboardRenders` | GET returns 200 |

All tests pass with `-race` flag. `successRepo` mock introduced for action tests — write ops return `nil` by default, unlike the `mockRepository` in auth tests which returns errors.

## Deviations from Plan

None — plan executed exactly as written.

## Known Stubs

None for this plan's goals. The dashboard is fully wired to `svc.ListAll` and all action methods. Remaining stubs in handler.go (NewPost, EditPost, SavePost, Preview) are intentional and documented — Plan 03-04 replaces them.

## Self-Check: PASSED

- [x] `internal/handler/admin/dashboard.go` — exists, contains `func (h *AdminHandler) Dashboard`, `r.URL.Query().Get("filter")`, `h.svc.ListAll`, `PopString`
- [x] `internal/handler/admin/actions.go` — exists, contains all 4 action handlers and correct flash copy
- [x] `internal/handler/admin/actions_test.go` — exists, contains all 6 required test functions
- [x] `web/templates/admin-dashboard.html` — exists, contains `filter-tabs`, `admin-table`, `status-pill`, `confirm-dialog`, `Delete this post?`, `Keep Post`, `Confirm Delete`, all empty state messages, `New Post`
- [x] `web/static/admin.css` — exists, contains `.filter-tab.active`, `.status-published`, `.status-draft`, `.confirm-dialog`, `.action-link--danger`, `.empty-state`
- [x] Commit 1d2ae03 — verified via git log
- [x] Commit 701d5c3 — verified via git log
- [x] `go build ./...` — passes
- [x] `go test ./internal/handler/admin/... -v -race` — PASS (10 tests)
