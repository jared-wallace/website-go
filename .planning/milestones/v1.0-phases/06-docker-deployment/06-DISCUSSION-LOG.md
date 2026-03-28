# Phase 6: Docker + Deployment - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md -- this log preserves the alternatives considered.

**Date:** 2026-03-27
**Phase:** 06-docker-deployment
**Areas discussed:** Deployment workflow, Migration strategy, Environment & secrets, Compose file strategy

---

## Deployment Workflow

| Option | Description | Selected |
|--------|-------------|----------|
| Manual SSH + pull (Recommended) | SSH into EC2, git pull or docker pull, docker compose up -d. Simple, fits single-instance blog. | :heavy_check_mark: |
| GHA CD pipeline | Push to main triggers GHA workflow that SSHs into EC2 or pushes to registry + deploys. | |
| You decide | Claude picks the simplest approach. | |

**User's choice:** Manual SSH + pull
**Notes:** None

### Follow-up: Build Location

| Option | Description | Selected |
|--------|-------------|----------|
| Build on-box (Recommended) | git pull + docker compose build + up. No registry needed. | :heavy_check_mark: |
| Push to GHCR, pull on EC2 | GHA builds and pushes to GitHub Container Registry. EC2 pulls. | |

**User's choice:** Build on-box
**Notes:** None

### Follow-up: Deploy Target Behavior

| Option | Description | Selected |
|--------|-------------|----------|
| Full sequence (Recommended) | make deploy runs: git pull, build, migrate, restart. One command. | :heavy_check_mark: |
| Document only | make deploy prints the steps to run manually. | |
| You decide | Claude picks the safest approach. | |

**User's choice:** Full sequence
**Notes:** None

---

## Migration Strategy

| Option | Description | Selected |
|--------|-------------|----------|
| Container entrypoint (Recommended) | Entrypoint script runs goose up before starting the server. Auto-applied, can't forget. | :heavy_check_mark: |
| Separate make step | make deploy runs goose up as a distinct step before docker compose up. | |
| You decide | Claude picks the best fit for single-instance. | |

**User's choice:** Container entrypoint
**Notes:** None

### Follow-up: Migration Mechanism

| Option | Description | Selected |
|--------|-------------|----------|
| Embedded in binary (Recommended) | Add -migrate flag to server binary. Entrypoint runs ./server -migrate then ./server. No goose CLI in container. | :heavy_check_mark: |
| Goose CLI in container | Include goose binary in Docker image. Entrypoint runs goose up. Adds ~10MB. | |

**User's choice:** Embedded in binary
**Notes:** None

---

## Environment & Secrets

| Option | Description | Selected |
|--------|-------------|----------|
| .env file on server (Recommended) | .env on EBS volume, referenced by docker-compose env_file. Survives rebuilds, gitignored. | :heavy_check_mark: |
| Inline in docker-compose | Hardcode in docker-compose.prod.yml environment block. | |
| AWS Secrets Manager | Fetch from AWS at container start. More secure but adds SDK dependency. | |

**User's choice:** .env file on server
**Notes:** None

### Follow-up: Missing .env Handling

| Option | Description | Selected |
|--------|-------------|----------|
| Template + fail (Recommended) | Ship .env.example with placeholders. Deploy fails with helpful message if .env missing. | :heavy_check_mark: |
| Fail only | No template. Deploy just fails if .env missing. | |
| You decide | Claude picks safest approach. | |

**User's choice:** Template + fail
**Notes:** None

---

## Compose File Strategy

| Option | Description | Selected |
|--------|-------------|----------|
| Separate files (Recommended) | Keep docker-compose.dev.yml for dev, add docker-compose.yml for prod. Clean separation. | :heavy_check_mark: |
| Base + override | docker-compose.yml base + override files. More DRY but more cognitive overhead. | |
| Single file with profiles | One file with profiles. Compact but mixing dev/prod gets messy. | |

**User's choice:** Separate files
**Notes:** None

### Follow-up: Resilience

| Option | Description | Selected |
|--------|-------------|----------|
| Yes, both (Recommended) | restart: unless-stopped + health checks on app (/health) and Postgres (pg_isready). | :heavy_check_mark: |
| Restart only | Add restart policy but skip health checks. | |
| You decide | Claude picks appropriate resilience. | |

**User's choice:** Yes, both
**Notes:** None

---

## Claude's Discretion

- Dockerfile multi-stage build structure
- Entrypoint script implementation
- /health endpoint implementation
- Non-root user setup
- Makefile deploy target implementation details
- .env.example content
- Convenience targets (logs, status)
- Docker build cache optimization
- Container logging configuration

## Deferred Ideas

None -- discussion stayed within phase scope.
