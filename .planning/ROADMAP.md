# Roadmap: website-go

## Overview

Start with a solid Go project scaffold wired to Postgres, then build the public reading experience before touching auth. Once readers can read, give the admin a way to write. Add distribution (RSS, OG, reactions, sitemap) after content flows. Lock down the security-sensitive surfaces (image upload, API push) once the core is stable. Wrap everything in Docker and ship.

## Phases

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (2.1, 2.2): Urgent insertions (marked with INSERTED)

Decimal phases appear between their surrounding integers in numeric order.

- [x] **Phase 1: Foundation** - Go project scaffold, Postgres + goose migrations, markdown/bluemonday pipeline, Makefile, GHA CI (completed 2026-03-26)
- [ ] **Phase 2: Public Blog** - Full public reading experience: post listing, single post, slugs, nautical design, dark mode, ToC, 404
- [ ] **Phase 3: Admin Panel** - Session auth, post CRUD, draft/publish workflow, split-pane markdown editor
- [ ] **Phase 4: Distribution** - RSS feed, Open Graph meta tags, sitemap, thumbs-up reactions
- [ ] **Phase 5: API + Images** - Image upload (magic-byte validated, EBS-stored) and API push endpoint with bearer token
- [ ] **Phase 6: Docker + Deployment** - Multi-stage Dockerfile, docker-compose with Postgres sidecar, EBS bind-mount, production-ready Makefile targets

## Phase Details

### Phase 1: Foundation
**Goal**: A compilable Go binary exists with the correct project structure, a live Postgres connection with embedded migrations, a goldmark + bluemonday markdown pipeline, an organized Makefile, and a passing GHA CI pipeline.
**Depends on**: Nothing (first phase)
**Requirements**: FOUND-01, FOUND-02, FOUND-03, FOUND-06, FOUND-07
**Success Criteria** (what must be TRUE):
  1. `go build ./...` succeeds and produces a runnable binary with standard Go layout (cmd/, internal/, db/, web/)
  2. `make migrate` runs goose migrations against a local Postgres instance without error
  3. A Go test calling the markdown pipeline with `<script>alert(1)</script>` input produces sanitized HTML (no script tag in output)
  4. `make lint` and `make test` pass; GHA CI runs both on push and reports green
  5. All Makefile targets (build, test, lint, run, docker, migrate) are documented and functional
**Plans:** 3/3 plans complete

Plans:
- [x] 01-01-PLAN.md -- Go module, directory skeleton, config, database, migrations
- [x] 01-02-PLAN.md -- Goldmark + bluemonday markdown rendering pipeline
- [x] 01-03-PLAN.md -- Makefile, dev tooling, GHA CI pipeline

### Phase 2: Public Blog
**Goal**: A reader visiting jared-wallace.com can browse, discover, and read published posts in a distinctive weathered beach bar design, on any device.
**Depends on**: Phase 1
**Requirements**: BLOG-01, BLOG-02, BLOG-03, BLOG-04, BLOG-05, BLOG-08, BLOG-11, BLOG-12
**Success Criteria** (what must be TRUE):
  1. Reader sees a paginated list of published posts sorted by date with titles and published dates
  2. Reader opens a post at /posts/my-readable-slug and sees markdown rendered with syntax-highlighted code blocks and an auto-generated table of contents
  3. Reader sees published date and estimated reading time on each post
  4. Reader sees the weathered beach bar nautical design on mobile and desktop (no layout breakage)
  5. Reader can toggle dark mode and sees an on-theme 404 page for invalid URLs
**Plans:** 1/3 plans executed
**UI hint**: yes

Plans:
- [x] 02-01-PLAN.md -- Data + service layer: migration, repository, pagination, reading time, ToC, excerpt
- [ ] 02-02-PLAN.md -- Templates, CSS nautical design, handlers, main.go wiring
- [ ] 02-03-PLAN.md -- Visual verification checkpoint

### Phase 3: Admin Panel
**Goal**: The admin can securely log in at admin.jared-wallace.com and create, edit, publish, and soft-delete posts using a split-pane markdown editor with live preview.
**Depends on**: Phase 2
**Requirements**: ADMIN-01, ADMIN-02, ADMIN-03, ADMIN-04, ADMIN-05, ADMIN-06, ADMIN-07
**Success Criteria** (what must be TRUE):
  1. Admin logs in with email/password; session persists across browser refresh (Postgres-backed); session cookie carries HttpOnly + Secure + SameSite=Lax flags
  2. Admin can create a new post with title, body, and slug; post saves to the database
  3. Admin can edit an existing post and soft-delete it (post is recoverable, not destroyed)
  4. Admin can toggle a post between draft and published; only published posts appear on public routes
  5. Admin writes in a split-pane editor with live markdown preview updating as they type
**Plans:** 3/4 plans executed
**UI hint**: yes

Plans:
- [x] 03-01-PLAN.md -- Data layer: config, migration, repository writes, service writes, slug, middleware, hashpw
- [x] 03-02-PLAN.md -- Auth + admin shell: host router, SCS sessions, login flow, admin base template, CSS foundation
- [x] 03-03-PLAN.md -- Dashboard: post table with filter tabs, publish/unpublish/delete/restore actions
- [x] 03-04-PLAN.md -- Editor: split-pane markdown editor with live preview, post create/edit, visual verification

### Phase 4: Distribution
**Goal**: Published posts are discoverable via RSS, shareable with rich social previews, indexed by search engines via sitemap, and readers can express appreciation with a thumbs-up reaction.
**Depends on**: Phase 3
**Requirements**: BLOG-06, BLOG-07, BLOG-09, BLOG-10
**Success Criteria** (what must be TRUE):
  1. /rss returns a valid RSS 2.0 feed containing only published posts (draft posts never appear)
  2. /sitemap.xml lists all published post URLs and validates against the sitemap schema
  3. Sharing a post URL on Slack or Twitter renders the correct Open Graph title, description, and image
  4. Reader can tap the thumbs-up on a post; count increments and is rate-limited without requiring login
**Plans:** 2/3 plans executed

Plans:
- [x] 04-01-PLAN.md -- RSS 2.0 feed, XML sitemap, and robots.txt endpoints
- [x] 04-02-PLAN.md -- Open Graph and Twitter Card meta tags, RSS auto-discovery, fallback OG image
- [ ] 04-03-PLAN.md -- Thumbs-up reaction system (migration, handler, JS, CSS)

### Phase 5: API + Images
**Goal**: The admin can upload images for embedding in posts and push .md files from a local machine via a bearer-token-authenticated API endpoint.
**Depends on**: Phase 3
**Requirements**: ADMIN-08, ADMIN-09
**Success Criteria** (what must be TRUE):
  1. Admin uploads a JPEG or PNG via the admin panel; image is stored on the EBS volume with a server-generated random filename; client-supplied filename is never used
  2. A file with a spoofed MIME extension (e.g., .jpg with HTML content) is rejected at the magic-byte check
  3. `POST /api/push` with a valid bearer token accepts a .md file body and creates or upserts a post by slug
  4. `POST /api/push` with a missing or invalid token returns 401; no post is created
**Plans**: TBD

### Phase 6: Docker + Deployment
**Goal**: The application ships as a minimal Docker container in a docker-compose stack with a Postgres sidecar, correctly bind-mounted to the EBS volume, and deployable behind the existing Nginx + ALB without manual steps.
**Depends on**: Phase 5
**Requirements**: FOUND-04, FOUND-05
**Success Criteria** (what must be TRUE):
  1. `docker build` produces a multi-stage image with a non-root runtime; final image contains only the binary and embedded assets
  2. `docker compose up` starts the app and Postgres; app serves traffic on :8080 within 30 seconds
  3. Postgres data directory is bind-mounted from /var/www/html/pgdata; container stops and restarts without data loss
  4. Makefile `deploy` target documents (and enforces) the `chown 999:999` prerequisite for the Postgres bind-mount directory
**Plans**: TBD

## Progress

**Execution Order:**
Phases execute in numeric order: 1 -> 2 -> 3 -> 4 -> 5 -> 6

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Foundation | 3/3 | Complete   | 2026-03-26 |
| 2. Public Blog | 1/3 | In Progress|  |
| 3. Admin Panel | 3/4 | In Progress|  |
| 4. Distribution | 2/3 | In Progress|  |
| 5. API + Images | 0/? | Not started | - |
| 6. Docker + Deployment | 0/? | Not started | - |
