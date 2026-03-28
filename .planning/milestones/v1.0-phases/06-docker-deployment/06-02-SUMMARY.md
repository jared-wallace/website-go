---
phase: 06-docker-deployment
plan: 02
subsystem: infra
tags: [docker-compose, postgres-sidecar, ebs-volumes, makefile, deploy]
dependency_graph:
  requires: [06-01]
  provides: [docker-compose-prod, makefile-deploy]
  affects: [Makefile, docker-compose.yml]
tech_stack:
  added: []
  patterns: [compose-health-checks, ebs-bind-mount, env-file-injection]
key_files:
  created:
    - docker-compose.yml
  modified:
    - Makefile
decisions:
  - Absolute env_file path (/var/www/html/.env) means compose config cannot validate on dev machines -- intentional production-only design
  - --no-cache on docker compose build ensures clean images on each deploy; build time is trivial for a Go binary
  - logs target filters to app service only for cleaner output; use docker compose logs for both
metrics:
  duration: 2min
  completed: "2026-03-28T02:12:00Z"
---

# Phase 06 Plan 02: Docker Compose + Makefile Deploy Summary

Production docker-compose.yml with app + Postgres sidecar, EBS bind-mounts for pgdata and images, health checks on both services, and Makefile deploy/logs/status targets with .env existence guard.

## Tasks Completed

### Task 1: Production docker-compose.yml
**Commit:** 50755e2

Created `docker-compose.yml` with two services:
- **postgres**: postgres:16-alpine, EBS bind-mount at /var/www/html/pgdata, pg_isready health check, restart: unless-stopped
- **app**: builds from Dockerfile, port 8080, env_file from /var/www/html/.env, EBS image mount at /var/www/html/images, depends_on postgres (service_healthy), wget /health check with 30s start_period

Key design: No `version:` key (compose v2 ignores it), `$$` escaping in pg_isready for proper shell interpolation, start_period gives migrations time to run.

### Task 2: Makefile deploy, logs, and status targets
**Commit:** 306e342

Added three targets to existing Makefile:
- **deploy**: checks /var/www/html/.env exists (fails with helpful message if missing), git pull, docker compose build --no-cache, docker compose up -d
- **logs**: docker compose logs -f app
- **status**: docker compose ps

All three appear in `make help` output.

## Deviations from Plan

### Minor Adjustment

**1. [Observation] docker compose config --quiet cannot pass on dev machines**
- **Found during:** Task 1 verification
- **Issue:** The absolute env_file path `/var/www/html/.env` causes compose config to fail when that path doesn't exist (dev machines)
- **Resolution:** Verified YAML syntax is valid; compose warnings are about missing env vars, not syntax errors. This is a production-only compose file by design (D-07). YAML structure validated independently.

## Known Stubs

None -- both files are complete production configurations with no placeholder data.

## Verification Results

1. service_healthy present in docker-compose.yml -- PASS
2. /var/www/html/pgdata bind-mount present -- PASS
3. /var/www/html/images bind-mount present -- PASS
4. make help shows deploy target -- PASS
5. .env guard in Makefile -- PASS
6. go test ./... -count=1 -- PASS (all existing tests pass)
