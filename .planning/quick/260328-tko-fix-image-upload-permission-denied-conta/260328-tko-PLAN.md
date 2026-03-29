---
phase: quick
plan: 260328-tko
type: execute
wave: 1
depends_on: []
files_modified:
  - Dockerfile
  - entrypoint.sh
autonomous: true
must_haves:
  truths:
    - "Container starts as root, fixes image dir ownership, then drops to appuser before running the server"
    - "Image uploads succeed in production with bind-mounted /var/www/html/images owned by root on the host"
  artifacts:
    - path: "entrypoint.sh"
      provides: "Permission fix and privilege drop"
      contains: "chown"
    - path: "Dockerfile"
      provides: "Runtime with su-exec and entrypoint"
      contains: "su-exec"
  key_links:
    - from: "Dockerfile"
      to: "entrypoint.sh"
      via: "ENTRYPOINT"
      pattern: "ENTRYPOINT.*entrypoint"
---

<objective>
Fix image upload "permission denied" error in Docker production deployment.

The container runs as appuser (UID 1001) but the bind-mounted host directory
`/var/www/html/images` is owned by `root:root`. The Go app's `os.MkdirAll` at
startup cannot fix permissions on an already-existing bind mount point.

Purpose: Allow image uploads to work in production without manual host-side chown.
Output: Updated Dockerfile + new entrypoint.sh that fixes ownership at container start, then drops privileges.
</objective>

<execution_context>
@$HOME/.claude/get-shit-done/workflows/execute-plan.md
@$HOME/.claude/get-shit-done/templates/summary.md
</execution_context>

<context>
@Dockerfile
@docker-compose.yml
@cmd/server/main.go (lines 58-66 — MkdirAll for IMAGE_DIR)
@internal/config/config.go (IMAGE_DIR env var default)
</context>

<tasks>

<task type="auto">
  <name>Task 1: Create entrypoint script and update Dockerfile</name>
  <files>entrypoint.sh, Dockerfile</files>
  <action>
1. Create `entrypoint.sh` in the repo root:
   - `#!/bin/sh` with `set -e`
   - Read IMAGE_DIR from env, default to `/var/www/html/images`
   - If IMAGE_DIR directory exists, run `chown -R 1001:1001 "$IMAGE_DIR"`. If it does not exist, `mkdir -p "$IMAGE_DIR"` then chown it.
   - `exec su-exec appuser "$@"` to drop privileges and exec the server
   - Keep it minimal — no logging beyond what's necessary

2. Update `Dockerfile` runtime stage:
   - Add `su-exec` to the `apk add` line: `apk add --no-cache ca-certificates tzdata su-exec`
   - Remove the `USER appuser` line (entrypoint handles the privilege drop)
   - Add `COPY entrypoint.sh .` after the server binary COPY
   - Add `RUN chmod +x entrypoint.sh` after the COPY
   - Change `CMD ["./server"]` to `ENTRYPOINT ["./entrypoint.sh"]` and `CMD ["./server"]`
   - Keep `RUN adduser -D -u 1001 appuser` (still needed for the user to exist)

The pattern mirrors how the official postgres image handles pgdata ownership — start as root, fix perms, drop to the service user.
  </action>
  <verify>
    <automated>cd /Users/jaredwallace/src/jared-wallace/website-go && docker build -t website-go-test . 2>&1 | tail -5 && docker run --rm website-go-test cat /app/entrypoint.sh && echo "--- Build and entrypoint verified ---"</automated>
  </verify>
  <done>
    - Dockerfile builds successfully with su-exec installed
    - entrypoint.sh is present in the container at /app/entrypoint.sh
    - entrypoint.sh contains chown logic for IMAGE_DIR and exec su-exec
    - No USER directive in Dockerfile (container starts as root, entrypoint drops privileges)
  </done>
</task>

</tasks>

<verification>
- `docker build -t website-go-test .` succeeds
- `docker run --rm -e IMAGE_DIR=/tmp/test-images website-go-test ls -la /tmp/test-images` shows directory owned by appuser (1001)
- entrypoint.sh is executable and contains the chown + su-exec pattern
</verification>

<success_criteria>
Docker container can write to a bind-mounted image directory regardless of host-side ownership, by fixing permissions at startup before dropping to appuser.
</success_criteria>

<output>
After completion, create `.planning/quick/260328-tko-fix-image-upload-permission-denied-conta/260328-tko-SUMMARY.md`
</output>
