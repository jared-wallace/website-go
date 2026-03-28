---
phase: quick
plan: 260327-vgp
subsystem: documentation
tags: [readme, deployment, aws, documentation]
dependency_graph:
  requires: []
  provides: [production-deployment-docs]
  affects: [README.md]
tech_stack:
  added: []
  patterns: []
key_files:
  created: []
  modified:
    - README.md
decisions:
  - Framed deployment steps as "Phase 6 intent" rather than operational instructions since Dockerfile/deploy.sh do not exist yet
  - Used table format for EBS storage layout for scannability
  - Included inline SSH/build/compose commands in code blocks to match existing README style
metrics:
  duration: 1min
  completed: 2026-03-28
---

# Phase quick Plan 260327-vgp: Add Production Deployment Section to README Summary

## One-liner

Added a Production Deployment section documenting the full ALB -> Nginx -> Go app request flow, EBS storage layout, and ASG/Postgres caveats for the AWS-hosted jared-wallace.com infrastructure.

## What Was Done

Added `## Production Deployment` section to `README.md` between the Continuous Integration section and License section. The section covers:

1. **Architecture Overview** — ASCII request flow diagram: Route53 -> ALB (TLS) -> EC2 t4g.micro ARM64 -> Nginx :80 -> Go app :8080. Notes the ASG is min=max=1 (self-healing, not horizontal scaling) and that infra is managed in the `aws-infra` repo via Terraform.

2. **Deployment Steps** — Six-step SSH + git pull + docker build + docker compose sequence, framed as the intended flow once Phase 6 Dockerfile/deploy.sh deliverables exist.

3. **Persistent Storage Layout** — Table of four EBS paths: mount point, app code, postgres-data (planned), images (planned).

4. **Caveats** — Three bullets: ASG max_size=1 + delete_on_termination=false requirement, chown 999:999 prereq for Postgres container, and explicit callout that Dockerfile/deploy.sh/docker-compose.prod.yml are Phase 6 deliverables not yet created.

## Tasks Completed

| Task | Name | Commit | Files |
|------|------|--------|-------|
| 1 | Add Production Deployment section to README.md | d275590 | README.md |

## Deviations from Plan

None — plan executed exactly as written.

## Known Stubs

The deployment steps reference `docker-compose.prod.yml` and `Dockerfile` that do not yet exist. This is intentional and explicitly called out in the section with a Phase 6 note. The section itself is accurate documentation of planned behavior, not a stub that blocks the plan's goal (which is documentation, not the artifacts themselves).

## Self-Check: PASSED

- README.md exists and contains exactly one "## Production Deployment" heading
- Commit d275590 verified in git log
- Section appears between CI and License sections as specified
