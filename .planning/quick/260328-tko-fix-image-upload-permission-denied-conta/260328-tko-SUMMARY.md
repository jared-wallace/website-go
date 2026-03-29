---
phase: quick
plan: 260328-tko
subsystem: docker/infrastructure
tags: [docker, permissions, entrypoint, su-exec]
dependency_graph:
  requires: []
  provides: [production image upload fix]
  affects: [Dockerfile, entrypoint.sh]
tech_stack:
  added: [su-exec]
  patterns: [root-to-user entrypoint privilege drop]
key_files:
  created: [entrypoint.sh]
  modified: [Dockerfile]
decisions:
  - Used su-exec (not gosu) for minimal Alpine-compatible privilege drop
  - entrypoint mirrors postgres image pattern: start root, fix perms, drop to service user
metrics:
  duration: 3min
  completed: 2026-03-28
  tasks: 1
  files: 2
---

# Quick Task 260328-tko: Fix Image Upload Permission Denied Summary

**One-liner:** Root-start entrypoint with su-exec privilege drop fixes bind-mount ownership so image uploads succeed without manual host-side chown.

## What Was Done

Created `entrypoint.sh` and updated `Dockerfile` so the container starts as root, runs `chown -R 1001:1001` on the `IMAGE_DIR` bind-mount, then drops privileges via `su-exec appuser` before executing the Go server.

## Tasks Completed

| Task | Name | Commit | Files |
|------|------|--------|-------|
| 1 | Create entrypoint script and update Dockerfile | 4a64b9a | entrypoint.sh, Dockerfile |

## Decisions Made

- **su-exec over gosu:** `su-exec` is the Alpine-native single-binary privilege-drop tool; avoids importing a larger binary
- **Entrypoint pattern matches postgres image:** Start as root to own the volume, chown to service UID, exec the app — no manual host-side `chown 1001:1001` required on first deploy
- **Removed `USER appuser` directive:** The entrypoint handles privilege drop via `exec su-exec appuser`; a `USER` directive before the entrypoint would prevent the root-level chown from working

## Deviations from Plan

None — plan executed exactly as written.

## Known Stubs

None.

## Verification Results

- `docker build -t website-go-test .` succeeded
- `docker run --rm website-go-test cat /app/entrypoint.sh` shows correct chown + su-exec content
- `docker run --rm -e IMAGE_DIR=/tmp/test-images website-go-test ls -la /tmp/test-images` shows directory owned by `appuser:appuser (1001)`

## Self-Check: PASSED
