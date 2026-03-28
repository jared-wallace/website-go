---
phase: 6
slug: docker-deployment
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-27
---

# Phase 6 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test + docker compose |
| **Config file** | docker-compose.yml, Dockerfile |
| **Quick run command** | `go test ./...` |
| **Full suite command** | `docker compose build && docker compose up -d && curl -sf http://localhost:8080/health` |
| **Estimated runtime** | ~30 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./...`
- **After every plan wave:** Run `docker compose build && docker compose up -d && curl -sf http://localhost:8080/health`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 30 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 6-01-01 | 01 | 1 | FOUND-04 | build | `docker build -t website-go .` | ❌ W0 | ⬜ pending |
| 6-01-02 | 01 | 1 | FOUND-04 | inspect | `docker inspect website-go --format '{{.Config.User}}'` | ❌ W0 | ⬜ pending |
| 6-02-01 | 02 | 1 | FOUND-05 | integration | `docker compose up -d && curl -sf http://localhost:8080/health` | ❌ W0 | ⬜ pending |
| 6-02-02 | 02 | 1 | FOUND-05 | integration | `docker compose down && docker compose up -d && psql check` | ❌ W0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `Dockerfile` — multi-stage build for Go binary
- [ ] `docker-compose.yml` — prod stack with app + Postgres
- [ ] `.env.example` — template for production env vars

*Existing go test infrastructure covers unit tests. Docker artifacts are the new infrastructure.*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| EBS bind-mount persistence | FOUND-05 | Requires actual EBS volume | SSH to EC2, stop/start containers, verify data persists |
| Nginx/ALB proxying | FOUND-05 | Requires production infra | Verify HTTPS traffic routes through ALB -> Nginx -> :8080 |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 30s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
