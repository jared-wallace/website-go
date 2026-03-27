# website-go

## What This Is

A personal blog platform for jared-wallace.com, built as a Go web server with a "weathered beach bar" nautical design. The public blog lives at jared-wallace.com, with an admin panel at admin.jared-wallace.com for writing and managing markdown posts. Deployed as a Docker container behind an existing AWS ALB + Nginx reverse proxy.

## Core Value

A reader visits jared-wallace.com and reads well-rendered markdown blog posts with images in a distinctive, memorable design.

## Requirements

### Validated

- [x] Public blog with markdown rendering and image support — Validated in Phase 2: public-blog
- [x] Admin panel with session-based auth (single admin user) — Validated in Phase 3: admin-panel
- [x] Web-based markdown editor with preview — Validated in Phase 3: admin-panel

### Active
- [ ] API endpoint for pushing .md files from local machine
- [x] RSS feed — Validated in Phase 4: distribution
- [x] Thumbs-up reaction counter on posts — Validated in Phase 4: distribution
- [ ] Weathered beach bar nautical design (driftwood, sand, ocean blues, chalkboard vibes)
- [ ] Dockerized deployment (app + Postgres sidecar via docker-compose)
- [ ] Organized Makefile following best practices
- [ ] GHA CI pipeline (lint, test, build) on free tier
- [ ] Standard Go project structure with minimal dependencies
- [ ] LaTeX rendering in posts (stretch goal)

### Out of Scope

- User registration / reader accounts — single-admin blog, no need
- Comments system — not for v1, maybe later
- OAuth / social login — session + bcrypt is sufficient for single user
- S3 image storage — using EBS volume for simplicity
- Managed database (RDS) — Postgres sidecar on EBS volume instead

## Context

- **Existing infrastructure**: AWS infra managed via Terraform in ../aws-infra
  - ALB with TLS termination (ACM cert for jared-wallace.com + wildcards)
  - ASG (min 1, max 2, desired 1) — effectively single-instance with self-healing
  - Nginx on EC2 reverse-proxying to :8080
  - 10GB EBS volume at /var/www/html for persistent storage
  - Docker pre-installed via user data script
  - Route53 DNS with A/AAAA records for root, www, and admin subdomains
  - Lambda + EventBridge for dynamic DNS updates on ASG events
- **Deployment model**: Docker container running on :8080, Nginx forwards traffic from ALB
- **Storage**: Postgres data and uploaded images both on the EBS volume
- **Design direction**: "Weathered bar by the beach" — warm wood tones, sandy off-whites, deep ocean blues, slightly rough/textured feel (driftwood, not yacht club), chalkboard-style elements, rope/knot dividers, anchor accents. Use the `frontend-design` skill for template crafting.

## Constraints

- **Tech stack**: Go with minimal dependencies — avoid large frameworks, prefer stdlib where reasonable
- **Infrastructure**: Must run as Docker container on port 8080 behind existing Nginx/ALB
- **Budget**: GHA CI must work on free tier (no paid GitHub features)
- **Storage**: All persistent data (DB + images) must live on the EBS volume at /var/www/html
- **Design**: Leverage the `frontend-design` skill for all template/UI work

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Postgres sidecar over managed DB | EBS volume available, avoids RDS cost, keeps infra simple | — Pending |
| Session auth over OAuth | Single admin user, minimal deps philosophy | — Pending |
| EBS for images over S3 | Simpler, fewer AWS dependencies, sufficient for single-blog scale | — Pending |
| Subdomain admin panel (admin.jared-wallace.com) | Clean separation of public/admin, already wired in Terraform | — Pending |
| Dual content workflow (web editor + API push) | Flexibility — quick edits in browser, serious writing from local tools | — Pending |

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd:transition`):
1. Requirements invalidated? → Move to Out of Scope with reason
2. Requirements validated? → Move to Validated with phase reference
3. New requirements emerged? → Add to Active
4. Decisions to log? → Add to Key Decisions
5. "What This Is" still accurate? → Update if drifted

**After each milestone** (via `/gsd:complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state

---
*Last updated: 2026-03-27 after Phase 4 completion*
