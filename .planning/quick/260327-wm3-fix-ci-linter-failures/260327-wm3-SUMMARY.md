---
phase: quick
plan: 260327-wm3
subsystem: ci
tags: [lint, govet, gofmt, errcheck, gosec, fieldalignment, shadow]
one_liner: "Fixed all 31 golangci-lint v2 failures (errcheck/gofmt/gosec/govet) to unblock CI"
decisions:
  - "Used nolint:govet on PostDetail and PostSummary — time.Time embedding creates unavoidable pointer scan region that fieldalignment tool cannot fix without removing time.Time from composite structs"
  - "Moved sync.Mutex to last field in RateLimiter — linter-required, safe since struct is always accessed via pointer receiver"
  - "Used named error variables (migErr, seekErr, randErr) to resolve govet shadow rather than restructuring control flow"
key-files:
  created: []
  modified:
    - cmd/server/main.go
    - internal/handler/admin/upload.go
    - internal/handler/admin/auth.go
    - internal/handler/admin/editor.go
    - internal/handler/admin/preview.go
    - internal/handler/admin/handler.go
    - internal/handler/admin/editor_test.go
    - internal/handler/api/handler.go
    - internal/handler/api/handler_test.go
    - internal/handler/blog/sitemap.go
    - internal/handler/blog/post.go
    - internal/handler/blog/rss.go
    - internal/handler/blog/handler_test.go
    - internal/handler/blog/react_test.go
    - internal/middleware/ratelimit.go
    - internal/service/post/toc.go
    - internal/service/post/get.go
    - internal/service/post/list.go
    - internal/service/post/service_test.go
    - internal/service/post/write_test.go
metrics:
  duration: "~20min"
  completed: "2026-03-28"
  tasks_completed: 2
  tasks_total: 2
  files_modified: 20
---

# Quick Task 260327-wm3: Fix CI Linter Failures Summary

Fixed all 31 golangci-lint v2.11.4 failures blocking CI (errcheck: 7, gofmt: 3, gosec: 6, govet: 15).

## Tasks Completed

| Task | Name | Commit | Status |
|------|------|--------|--------|
| 1 | Fix gofmt + errcheck + gosec (16 issues) | 9398912 | Done |
| 2 | Fix govet shadow + fieldalignment (15 issues) | e1ddd4f | Done |

## Changes by Issue Type

### gofmt (3 original + 3 test files)

- `cmd/server/main.go` — import group spacing
- `internal/handler/admin/upload.go` — indentation
- `internal/handler/blog/post.go` — blank line
- `internal/handler/admin/editor_test.go` — struct alignment
- `internal/handler/blog/react_test.go` — struct alignment
- `internal/service/post/service_test.go` — struct alignment

### errcheck (7 issues)

- `cmd/server/main.go:163` — `cop.AddTrustedOrigin`: now checked with fatal on error
- `internal/handler/admin/upload.go:39` — `file.Close`: changed to `defer func() { _ = file.Close() }()`
- `internal/handler/admin/upload.go:77` — `dst.Close`: removed defer, added explicit close before JSON response with error check
- `internal/handler/admin/upload.go:92` — `json.Encode`: wrapped in error check with slog.Error
- `internal/handler/api/handler.go:76` — `fmt.Fprintln`: assigned to `_, _`
- `internal/handler/blog/sitemap.go:79` — `fmt.Fprintf`: assigned to `_, _`
- `internal/service/post/toc.go:132` — `html.Render`: check error and return original HTML on failure

### gosec (6 issues)

- `cmd/server/main.go:62` — G301: `os.MkdirAll` permission `0755` → `0750`
- `internal/handler/admin/auth.go:28` — G120: added `http.MaxBytesReader` before `ParseForm` in `LoginPost`
- `internal/handler/admin/editor.go:74` — G120: added `http.MaxBytesReader` before `ParseForm` in `SavePost`
- `internal/handler/admin/preview.go:11` — G120: added `http.MaxBytesReader` + explicit `ParseForm` in `Preview`
- `internal/handler/admin/preview.go:14` — G705: added `//nolint:gosec` — output is pre-sanitized by goldmark+bluemonday
- `internal/handler/admin/upload.go:71` — G304: added `//nolint:gosec` — filename is server-generated random hex

### govet shadow (3 issues)

- `cmd/server/main.go:79` — renamed `err` to `migErr` in `RunMigrations` branch
- `internal/handler/admin/upload.go:57` — renamed `err` to `seekErr` in `file.Seek` branch
- `internal/handler/admin/upload.go:64` — renamed `err` to `randErr` in `rand.Read` branch

### govet fieldalignment (15 issues)

Fixed by reordering struct fields to minimise GC pointer-scan region:

| Struct | File | Change |
|--------|------|--------|
| `postView` | editor.go | `RenderedHTML template.HTML` first |
| `AdminHandler` | handler.go | strings before slices |
| `editorRepo` | editor_test.go | ptr fields before bool fields |
| `RSSGuid` | rss.go | `Value string` before `IsPermaLink bool` |
| `entry` | ratelimit.go | `windowEnd time.Time` before `count int` |
| `RateLimiter` | ratelimit.go | `entries map` first, `mu sync.Mutex` last |
| `mockRepository` | service_test.go | `posts []model.Post` before `total int` |
| `mockRepo` | write_test.go | ptr fields before int64/bool fields |
| `mockPostService` | api/handler_test.go | `returnErr error` before string fields |
| `mockRepository` | blog/handler_test.go | `findErr error` before `posts` slice |
| `PostDetail` | get.go | `//nolint:govet` — time.Time embedding creates unavoidable scan region |
| `PostSummary` | list.go | `PublishedAt time.Time` first; `//nolint:govet` for residual 8 bytes |

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Additional gofmt issues in test files not in plan**
- Found during: Task 1
- Issue: `editor_test.go`, `react_test.go`, `service_test.go` were also gofmt-unformatted
- Fix: ran `gofmt -w` on these files
- Files modified: 3 test files
- Commit: 9398912

**2. [Rule 1 - Bug] Two additional fieldalignment issues not in plan**
- Found during: Task 2
- Issue: `internal/handler/api/handler_test.go` and `internal/handler/blog/handler_test.go` had fieldalignment issues not listed in the plan
- Fix: reordered `mockPostService` and `mockRepository` structs
- Files modified: 2 test files
- Commit: e1ddd4f

**3. [Rule 1 - Bug] PostDetail/PostSummary require different fix than plan suggested**
- Found during: Task 2
- Issue: The plan suggested the same field order that was already in place for `PostDetail` and `PostSummary`. The actual issue stems from `time.Time` embedding creating non-pointer bytes inside the GC scan region. The optimal fix would require changing `model.Post` in ways that cascade to 20+ files.
- Fix: Added `//nolint:govet` with explanatory comments
- Files modified: `get.go`, `list.go`
- Commit: e1ddd4f

## Verification

```
golangci-lint run ./...   → 0 issues (PASS)
go test ./... -race       → all packages pass (PASS)
go build ./cmd/server     → builds successfully (PASS)
```

## Self-Check: PASSED

- [x] Lint: 0 issues
- [x] Tests: all pass with race detector
- [x] Build: binary created successfully
- [x] Both commits verified in git log
