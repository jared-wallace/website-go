---
phase: 5
slug: api-images
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-27
---

# Phase 5 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | Go stdlib `testing` |
| **Config file** | none (standard `go test ./...`) |
| **Quick run command** | `go test ./internal/handler/admin/... ./internal/handler/api/... ./internal/middleware/... ./internal/service/post/... -run TestUpload\|TestPush\|TestBearer` |
| **Full suite command** | `go test ./...` |
| **Estimated runtime** | ~5 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./internal/handler/admin/... ./internal/handler/api/... ./internal/service/post/...`
- **After every plan wave:** Run `go test ./...`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 5 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 05-01-01 | 01 | 1 | ADMIN-08 | unit | `go test ./internal/handler/admin/... -run TestUploadImage_ValidJPEG` | ❌ W0 | ⬜ pending |
| 05-01-02 | 01 | 1 | ADMIN-08 | unit | `go test ./internal/handler/admin/... -run TestUploadImage_SpoofedMIME` | ❌ W0 | ⬜ pending |
| 05-01-03 | 01 | 1 | ADMIN-08 | unit | `go test ./internal/handler/admin/... -run TestUploadImage_TooLarge` | ❌ W0 | ⬜ pending |
| 05-01-04 | 01 | 1 | ADMIN-08 | unit | `go test ./internal/handler/admin/... -run TestUploadImage_RandomFilename` | ❌ W0 | ⬜ pending |
| 05-02-01 | 02 | 1 | ADMIN-09 | unit | `go test ./internal/handler/api/... -run TestPushPost_ValidToken` | ❌ W0 | ⬜ pending |
| 05-02-02 | 02 | 1 | ADMIN-09 | unit | `go test ./internal/handler/api/... -run TestPushPost_NoToken` | ❌ W0 | ⬜ pending |
| 05-02-03 | 02 | 1 | ADMIN-09 | unit | `go test ./internal/handler/api/... -run TestPushPost_InvalidToken` | ❌ W0 | ⬜ pending |
| 05-02-04 | 02 | 1 | ADMIN-09 | unit | `go test ./internal/service/post/... -run TestUpsertBySlug` | ❌ W0 | ⬜ pending |
| 05-02-05 | 02 | 1 | ADMIN-09 | unit | `go test ./internal/handler/api/... -run TestPushPost_NoSlug` | ❌ W0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `internal/handler/admin/upload_test.go` — stubs for ADMIN-08 upload cases
- [ ] `internal/handler/api/handler_test.go` — stubs for ADMIN-09 push cases
- [ ] `internal/service/post/upsert_test.go` — stubs for UpsertBySlug logic

*Existing `go test` infrastructure covers framework needs — no additional framework install required.*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Image renders in browser after upload | ADMIN-08 | Requires visual inspection of rendered page | Upload image via admin, embed in draft post, verify `<img>` renders in preview |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 5s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
