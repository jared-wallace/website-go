---
phase: 1
slug: foundation
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-26
---

# Phase 1 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | `go test` (stdlib) |
| **Config file** | none — standard Go test conventions |
| **Quick run command** | `go test ./internal/markdown/... -v` |
| **Full suite command** | `go test ./... -v -race` |
| **Estimated runtime** | ~5 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go build ./...`
- **After every plan wave:** Run `go test ./... -v -race`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 10 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 01-01 | 01 | 1 | FOUND-01 | smoke | `go build ./...` | ❌ W0 | ⬜ pending |
| 01-02 | 01 | 1 | FOUND-02 | integration | `go test ./internal/database/... -v` | ❌ W0 | ⬜ pending |
| 01-03 | 01 | 1 | FOUND-03 | integration | `go test ./internal/database/... -run TestMigrations -v` | ❌ W0 | ⬜ pending |
| 01-04 | 01 | 1 | FOUND-01 | unit | `go test ./internal/markdown/... -v` | ❌ W0 | ⬜ pending |
| 01-05 | 01 | 1 | FOUND-06 | smoke | `make build && make lint && make test` | ❌ W0 | ⬜ pending |
| 01-06 | 01 | 1 | FOUND-07 | e2e | Push to branch; observe GHA Actions tab | ❌ W0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `internal/markdown/renderer_test.go` — XSS test + basic rendering test (pure function, no Postgres)
- [ ] `internal/database/database_test.go` — connection + migration integration tests (requires Postgres)
- [ ] `.github/workflows/ci.yml` — GHA CI pipeline
- [ ] `Makefile` — dev workflow targets
- [ ] `go.mod` + `go.sum` — module initialization (prerequisite for all)

*All test files are missing — greenfield project.*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| GHA CI reports green on push | FOUND-07 | Requires actual GitHub Actions run | Push to branch, verify Actions tab shows green |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 10s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
