---
phase: 4
slug: distribution
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-27
---

# Phase 4 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | Go stdlib `testing` + `net/http/httptest` |
| **Config file** | none (standard `go test ./...`) |
| **Quick run command** | `go test ./internal/handler/blog/... ./internal/service/post/... ./internal/repository/post/...` |
| **Full suite command** | `go test ./...` |
| **Estimated runtime** | ~10 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./internal/handler/blog/... -run TestServe`
- **After every plan wave:** Run `go test ./...`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 10 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 04-01-01 | 01 | 1 | BLOG-09 | unit | `go test ./internal/handler/blog/... -run TestServeRSS` | ❌ W0 | ⬜ pending |
| 04-01-02 | 01 | 1 | BLOG-09 | unit | `go test ./internal/handler/blog/... -run TestRSSDraftExclusion` | ❌ W0 | ⬜ pending |
| 04-01-03 | 01 | 1 | BLOG-09 | unit | `go test ./internal/handler/blog/... -run TestRSSFullContent` | ❌ W0 | ⬜ pending |
| 04-02-01 | 02 | 1 | BLOG-07 | unit | `go test ./internal/handler/blog/... -run TestServeSitemap` | ❌ W0 | ⬜ pending |
| 04-02-02 | 02 | 1 | BLOG-07 | unit | `go test ./internal/handler/blog/... -run TestServeRobots` | ❌ W0 | ⬜ pending |
| 04-03-01 | 03 | 1 | BLOG-06 | unit | `go test ./internal/handler/blog/... -run TestPostOGMeta` | ❌ W0 | ⬜ pending |
| 04-03-02 | 03 | 1 | BLOG-06 | unit | `go test ./internal/handler/blog/... -run TestListOGMeta` | ❌ W0 | ⬜ pending |
| 04-04-01 | 04 | 2 | BLOG-10 | unit | `go test ./internal/handler/blog/... -run TestReact` | ❌ W0 | ⬜ pending |
| 04-04-02 | 04 | 2 | BLOG-10 | unit | `go test ./internal/handler/blog/... -run TestReactDuplicate` | ❌ W0 | ⬜ pending |
| 04-04-03 | 04 | 2 | BLOG-10 | unit | `go test ./internal/handler/blog/... -run TestPostReactionCount` | ❌ W0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `internal/handler/blog/rss_test.go` — stubs for BLOG-09 RSS tests
- [ ] `internal/handler/blog/sitemap_test.go` — stubs for BLOG-07
- [ ] `internal/handler/blog/react_test.go` — stubs for BLOG-10
- [ ] `internal/repository/post/reactions_test.go` — stubs for AddReaction and CountReactions
- [ ] Extend existing `mockRepository` in `handler_test.go` with `AddReaction` and `CountReactions` methods

*OG meta tests (BLOG-06) can be covered in the existing `handler_test.go` by asserting on rendered HTML body from `ListPosts` and `ShowPost`. No new test file needed.*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| OG preview renders correctly on Slack/Twitter | BLOG-06 | Requires external service validation | Share a post URL on Slack/Twitter and verify title, description, image render |
| RSS feed renders in a feed reader | BLOG-09 | Visual rendering verification | Open `/rss` URL in a feed reader (e.g., Feedly, NetNewsWire) |
| Thumbs-up CSS animation is visible | BLOG-10 | Visual/interaction check | Tap thumbs-up button, verify count increment animation |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 10s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
