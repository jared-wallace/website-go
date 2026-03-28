# website-go

## What This Is

A personal blog platform for jared-wallace.com branded as "The Wild Meridian", built as a Go web server with a weathered nautical design. The public blog lives at jared-wallace.com, with an admin panel at admin.jared-wallace.com for writing and managing markdown posts. Deployed as a Docker container behind an existing AWS ALB + Nginx reverse proxy.

## Core Value

A reader visits jared-wallace.com and reads well-rendered markdown blog posts with images in a distinctive, memorable design.

## Current State

**Shipped:** v1.0 MVP (2026-03-28), v1.1 The Wild Meridian (2026-03-28)
**Codebase:** ~29,600 lines of Go across 8 phases
**Tech stack:** Go 1.26, pgx/v5, goldmark+bluemonday, html/template, SCS sessions, goose migrations
**Deployment:** Multi-stage Docker (9.9MB alpine), docker-compose with Postgres sidecar, EBS bind-mounts

## Requirements

### Validated

- Public blog with markdown rendering, syntax highlighting, and image support -- v1.0
- Paginated post listing with reading time, excerpts, and tags -- v1.0
- URL slugs, auto-generated ToC, dark mode toggle, themed 404 -- v1.0
- Admin panel with bcrypt session auth at admin subdomain -- v1.0
- Post CRUD with draft/publish workflow and soft-delete -- v1.0
- Split-pane markdown editor with live preview -- v1.0
- RSS 2.0 feed, XML sitemap, Open Graph meta tags -- v1.0
- Thumbs-up reactions with IP deduplication -- v1.0
- Image upload with magic-byte MIME validation -- v1.0
- API push endpoint with bearer token auth -- v1.0
- Docker multi-stage build with Postgres sidecar and EBS mounts -- v1.0
- Makefile with build, test, lint, run, deploy, logs, status targets -- v1.0
- GHA CI pipeline (lint, test, build) on free tier -- v1.0
- Standard Go project layout with minimal dependencies -- v1.0
- Weathered beach bar nautical design -- v1.0
- Rebrand from "The Log" to "The Wild Meridian" across all surfaces -- v1.1
- Copyright footer with dynamic year and "Jared Wallace" -- v1.1
- Discreet RSS icon in footer with autodiscovery link tag -- v1.1
- About page from embedded markdown with nautical chrome and nav link -- v1.1

## Current Milestone: v1.2 Shore Leave Polish

**Goal:** Elevate the site's visual cohesion and nautical atmosphere through CSS/template-level design polish — no backend changes.

**Target features:**
- Move About link to footer and expand footer into a proper navigation/personality block
- Add background texture (CSS noise/grain) for weathered atmosphere
- Improve post card depth, tag visibility, and rope divider as a design motif
- Add dark mode color transitions and subtle page entry animations
- Fix reaction button radius consistency and add homepage heading
- Fix stale "The Log" CSS comment

### Active

- [x] About link relocated to footer -- Phase 11
- [x] Footer expanded with navigation and personality -- Phase 11
- [x] Background texture/grain for weathered atmosphere -- Phase 9
- [x] Post card depth and visual weight improvements -- Phase 9
- [x] Rope divider strengthened as design motif -- Phase 11
- [x] Dark mode color transitions -- Phase 10
- [x] Tag pill visibility improvements -- Phase 9
- [x] Page entry animations (fade-in, card stagger) -- Phase 10
- [x] Reaction button radius consistency (pill → 4px) -- Phase 9
- [x] Homepage heading/hero area -- Phase 11
- [x] CSS comment rebrand fix -- Phase 9
- [x] Mobile nav cleanup (resolved by About move) -- Phase 11

### Out of Scope

- User registration / reader accounts -- single-admin blog, zero reader value
- Comments system -- moderation overhead; thumbs-up covers lightweight engagement
- OAuth / social login -- session + bcrypt is sufficient for single user
- Full-text search -- pagination + browser search sufficient at blog scale
- Newsletter / email delivery -- RSS covers follow-me use case; revisit on demand
- In-app analytics -- use Plausible/Fathom externally instead
- Image optimization / CDN -- overkill for single-instance personal blog
- WYSIWYG editor -- markdown philosophy; split-pane preview is sufficient
- Multi-author support -- personal blog; revisit if needed
- S3 image storage -- EBS volume is simpler and sufficient
- Managed database (RDS) -- Postgres sidecar on EBS volume instead
- LaTeX rendering -- deferred to future milestone (ENH-01)

## Context

- **Existing infrastructure**: AWS infra managed via Terraform in ../aws-infra
  - ALB with TLS termination (ACM cert for jared-wallace.com + wildcards)
  - ASG (min 1, max 2, desired 1) -- effectively single-instance with self-healing
  - Nginx on EC2 reverse-proxying to :8080
  - 10GB EBS volume at /var/www/html for persistent storage
  - Docker pre-installed via user data script
  - Route53 DNS with A/AAAA records for root, www, and admin subdomains
  - Lambda + EventBridge for dynamic DNS updates on ASG events
- **Deployment model**: `make deploy` on EC2 runs docker-compose build + up
- **Storage**: Postgres data at /var/www/html/pgdata, images at /var/www/html/images
- **Design direction**: "Weathered bar by the beach" -- warm wood tones, sandy off-whites, deep ocean blues, slightly rough/textured feel (driftwood, not yacht club), chalkboard-style elements, rope/knot dividers, anchor accents

## Constraints

- **Tech stack**: Go with minimal dependencies -- avoid large frameworks, prefer stdlib where reasonable
- **Infrastructure**: Must run as Docker container on port 8080 behind existing Nginx/ALB
- **Budget**: GHA CI must work on free tier (no paid GitHub features)
- **Storage**: All persistent data (DB + images) must live on the EBS volume at /var/www/html
- **Design**: Leverage the `frontend-design` skill for all template/UI work

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Postgres sidecar over managed DB | EBS volume available, avoids RDS cost, keeps infra simple | Good -- works well for single-instance blog |
| Session auth over OAuth | Single admin user, minimal deps philosophy | Good -- bcrypt + SCS pgxstore is clean and secure |
| EBS for images over S3 | Simpler, fewer AWS dependencies, sufficient for single-blog scale | Good -- bind-mount works, no SDK dependency |
| Subdomain admin panel (admin.jared-wallace.com) | Clean separation of public/admin, already wired in Terraform | Good -- host-based routing in stdlib ServeMux |
| Dual content workflow (web editor + API push) | Flexibility -- quick edits in browser, serious writing from local tools | Good -- both paths work, UpsertBySlug handles idempotency |
| Per-page template sets (html/template) | Avoids block name collisions across pages in Go templating | Good -- each page parses base+page pair independently |
| goldmark html.WithUnsafe() + bluemonday | Allows full HTML rendering then sanitizes; pipeline order is security-critical | Good -- XSS tests pass, rendering is correct |
| In-memory dashboard filtering over separate queries | Blog scale makes N queries unnecessary; ListAll + filter is simpler | Good -- works at blog scale |
| API push on blogMux (not adminMux) | Bearer token is the auth gate, not host routing | Good -- decouples API from admin subdomain |
| CDATA xml.Marshaler for RSS | Prevents double-escaping of HTML in RSS descriptions | Good -- standard RSS 2.0 pattern |
| Embedded about.md (content package) over DB | Static content that rarely changes; avoids admin UI complexity | Good -- go:embed string, zero-dependency |

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition:**
1. Requirements invalidated? -> Move to Out of Scope with reason
2. Requirements validated? -> Move to Validated with phase reference
3. New requirements emerged? -> Add to Active
4. Decisions to log? -> Add to Key Decisions
5. "What This Is" still accurate? -> Update if drifted

**After each milestone** (via `/gsd:complete-milestone`):
1. Full review of all sections
2. Core Value check -- still the right priority?
3. Audit Out of Scope -- reasons still valid?
4. Update Context with current state

---
*Last updated: 2026-03-28 after Phase 11 completion (v1.2 final phase)*
