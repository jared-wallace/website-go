---
phase: 3
slug: admin-panel
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-26
---

# Phase 3 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | Go standard `testing` package + `net/http/httptest` |
| **Config file** | None (`go test ./...`) |
| **Quick run command** | `go test ./internal/handler/admin/... -v -race` |
| **Full suite command** | `go test ./... -v -race` |
| **Estimated runtime** | ~15 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./internal/handler/admin/... -v -race`
- **After every plan wave:** Run `go test ./... -v -race`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 15 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 03-01-01 | 01 | 1 | ADMIN-01 | unit (httptest) | `go test ./internal/handler/admin/... -run TestLoginPost -v` | No — Wave 0 | ⬜ pending |
| 03-01-02 | 01 | 1 | ADMIN-01 | unit (httptest) | `go test ./internal/handler/admin/... -run TestLoginPostInvalid -v` | No — Wave 0 | ⬜ pending |
| 03-01-03 | 01 | 1 | ADMIN-02 | unit (httptest) | `go test ./internal/middleware/... -run TestRequireSession -v` | No — Wave 0 | ⬜ pending |
| 03-02-01 | 02 | 2 | ADMIN-03 | unit (httptest) | `go test ./internal/handler/admin/... -run TestCreatePost -v` | No — Wave 0 | ⬜ pending |
| 03-02-02 | 02 | 2 | ADMIN-03 | unit | `go test ./internal/service/post/... -run TestGenerateSlug -v` | No — Wave 0 | ⬜ pending |
| 03-02-03 | 02 | 2 | ADMIN-04 | unit (httptest) | `go test ./internal/handler/admin/... -run TestEditPost -v` | No — Wave 0 | ⬜ pending |
| 03-02-04 | 02 | 2 | ADMIN-05 | unit | `go test ./internal/service/post/... -run TestSoftDelete -v` | No — Wave 0 | ⬜ pending |
| 03-02-05 | 02 | 2 | ADMIN-06 | unit | `go test ./internal/service/post/... -run TestTogglePublished -v` | No — Wave 0 | ⬜ pending |
| 03-03-01 | 03 | 2 | ADMIN-07 | unit (httptest) | `go test ./internal/handler/admin/... -run TestPreview -v` | No — Wave 0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `internal/handler/admin/handler_test.go` — test helper, mock setup following `internal/handler/blog/handler_test.go` pattern
- [ ] `internal/middleware/auth_test.go` — RequireSession middleware tests
- [ ] `internal/service/post/write_test.go` — slug generation, soft-delete, publish toggle
- [ ] `internal/repository/post/write_test.go` — SQL write queries (requires live Postgres, mark with `//go:build integration`)

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Split-pane editor layout | ADMIN-07 | Visual layout verification | Open editor in browser; verify textarea left, preview right; resize to mobile width, verify vertical stack with tab toggle |
| Session cookie flags | ADMIN-02 | Requires browser dev tools | Login; open DevTools > Application > Cookies; verify HttpOnly, Secure, SameSite=Lax |
| Nautical design consistency | D-02 | Visual design verification | Compare admin pages against public blog; verify shared color palette, typography, beach bar aesthetic |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 15s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
