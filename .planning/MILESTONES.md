# Milestones

## v1.2 Shore Leave Polish (Shipped: 2026-03-28)

**Phases completed:** 3 phases, 4 plans
**Files changed:** 5 files, +377/-48 lines (excluding .planning)
**Timeline:** 2026-03-28 (same day as v1.1)

**Key accomplishments:**

- SVG feTurbulence grain texture overlay with iOS-safe positioning and separate light/dark opacities
- Motion system with page fade-in, 6-card stagger entrance, and .theme-ready-gated dark mode transitions
- prefers-reduced-motion guard covering all animations including legacy reaction-bounce (WCAG 2.3.3)
- Footer redesigned as two-section personality block with relocated About link and ARIA landmarks
- Inline SVG twisted rope divider replacing CSS dashed border, theme-aware via custom properties
- List page hero heading with "The Wild Meridian" h1 and tagline

**Archive:** [v1.2-ROADMAP.md](milestones/v1.2-ROADMAP.md) | [v1.2-REQUIREMENTS.md](milestones/v1.2-REQUIREMENTS.md)

---

## v1.1 The Wild Meridian (Shipped: 2026-03-28)

**Phases completed:** 2 phases, 2 plans, 4 tasks
**Files changed:** 18 files, +214/-19 lines
**Timeline:** 2026-03-28 (same day as v1.0 — shipped in ~1 hour)

**Key accomplishments:**

- Rebranded all public surfaces from "The Log" to "The Wild Meridian" (header, RSS, OG meta, sitemap, admin panel)
- Updated copyright footer to "2026 Jared Wallace" with dynamic year
- Added discreet RSS broadcast icon in footer with 44px touch target and autodiscovery link tag
- Added /about page with embedded markdown rendered through Goldmark, nautical template, and nav link
- About page includes sitemap entry, focus-visible states, and mobile-responsive nav link

**Archive:** [v1.1-ROADMAP.md](milestones/v1.1-ROADMAP.md) | [v1.1-REQUIREMENTS.md](milestones/v1.1-REQUIREMENTS.md)

---

## v1.0 MVP (Shipped: 2026-03-28)

**Phases completed:** 6 phases, 17 plans, 125 commits
**Lines of Go code:** ~29,400
**Timeline:** 2026-03-26 to 2026-03-28

**Key accomplishments:**

- Go web server with pgx/v5 connection pool, embedded goose migrations, and goldmark+bluemonday markdown pipeline with XSS protection
- Nautical "weathered beach bar" public blog with dark mode, auto-generated ToC, paginated listing, and themed 404
- Admin panel at admin subdomain with bcrypt session auth, post CRUD, draft/publish workflow, and split-pane markdown editor with live preview
- Distribution layer: RSS 2.0 feed, XML sitemap, Open Graph meta tags, and IP-deduplicated thumbs-up reactions
- Image upload with magic-byte MIME validation and bearer-token authenticated API push endpoint
- Docker deployment: 9.9MB multi-stage alpine image with docker-compose Postgres sidecar and EBS bind-mounts

**Known gaps (accepted as tech debt):**

- BLOG-05 mobile responsive visual confirmation deferred to post-deploy
- 16 human verification items requiring live server/browser (all implementations complete)
- Nyquist validation not executed (6/6 phases draft)

**Archive:** [v1.0-ROADMAP.md](milestones/v1.0-ROADMAP.md) | [v1.0-REQUIREMENTS.md](milestones/v1.0-REQUIREMENTS.md) | [v1.0-MILESTONE-AUDIT.md](milestones/v1.0-MILESTONE-AUDIT.md)

---
