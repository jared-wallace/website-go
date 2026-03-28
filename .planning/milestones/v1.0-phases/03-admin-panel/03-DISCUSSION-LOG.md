# Phase 3: Admin Panel - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-03-26
**Phase:** 03-admin-panel
**Areas discussed:** Subdomain routing, Editor experience, Auth & session flow, Admin dashboard

---

## Subdomain Routing

| Option | Description | Selected |
|--------|-------------|----------|
| Host-based mux | Single binary checks Host header — admin.jared-wallace.com to admin handlers, else blog. Nginx already forwards both to :8080. | ✓ |
| Path-based /admin/* | All admin routes under /admin/* on same host. Simpler routing but DNS/Nginx config unused. | |

**User's choice:** Host-based mux (Recommended)
**Notes:** None

---

### Admin Design Style

| Option | Description | Selected |
|--------|-------------|----------|
| Shared nautical theme | Reuse base.html and nautical CSS — admin as "back office" of the beach bar. | ✓ |
| Minimal/utilitarian | Clean functional admin UI, separate base template. | |

**User's choice:** Shared nautical theme
**Notes:** None

---

### Unauthenticated Admin Visits

| Option | Description | Selected |
|--------|-------------|----------|
| Branded login page | Nautical-themed login form at admin subdomain. | ✓ |
| Redirect to public blog | Bounce unauthenticated to jared-wallace.com. | |

**User's choice:** Branded login page (Recommended)
**Notes:** None

---

## Editor Experience

### Preview Technology

| Option | Description | Selected |
|--------|-------------|----------|
| Vanilla JS fetch | Debounced fetch() to server endpoint, goldmark+bluemonday renders. Preview identical to published. | ✓ |
| HTMX hx-post | Same server-rendered approach with HTMX attributes. Adds ~14KB dependency. | |
| Client-side JS library | marked.js or markdown-it for instant preview. May differ from server output. | |

**User's choice:** Vanilla JS fetch (Recommended)
**Notes:** None

---

### Autosave Behavior

| Option | Description | Selected |
|--------|-------------|----------|
| Manual save only | Explicit Save Draft / Publish buttons. Ctrl+S / Cmd+S shortcut. | ✓ |
| Periodic autosave | Auto-save to DB every 30-60 seconds. | |
| localStorage draft backup | Auto-save to localStorage, explicit save to DB. | |

**User's choice:** Manual save only (Recommended)
**Notes:** None

---

### Mobile Editor

| Option | Description | Selected |
|--------|-------------|----------|
| Stacked with toggle | Vertical stack on narrow screens with Write/Preview tab toggle. | ✓ |
| Editor only, no preview | Only textarea on mobile. Preview via explicit button. | |

**User's choice:** Stacked with toggle (Recommended)
**Notes:** None

---

### Editor Toolbar

| Option | Description | Selected |
|--------|-------------|----------|
| Plain textarea | Monospace textarea, no toolbar buttons. | ✓ |
| Minimal toolbar | Row of buttons for bold, italic, link, heading, code block. | |

**User's choice:** Plain textarea (Recommended)
**Notes:** None

---

## Auth & Session Flow

### Credential Storage

| Option | Description | Selected |
|--------|-------------|----------|
| Environment variables | ADMIN_EMAIL and ADMIN_PASSWORD_HASH as env vars. make hash-password target. | ✓ |
| Database admin table | Store in admins table with migration and seed script. | |

**User's choice:** Environment variables (Recommended)
**Notes:** None

---

### Session Lifetime

| Option | Description | Selected |
|--------|-------------|----------|
| 24 hours | Inactivity-based expiry after 24 hours. | ✓ |
| 7 days | Stay logged in for a week. | |

**User's choice:** 24 hours (Recommended)
**Notes:** None

---

### Failed Login Handling

| Option | Description | Selected |
|--------|-------------|----------|
| Generic error + rate limit | "Invalid email or password" + 5 attempts/min/IP rate limiter. | ✓ |
| Generic error only | Generic message, no rate limiting. | |

**User's choice:** Generic error + rate limit (Recommended)
**Notes:** None

---

## Admin Dashboard

### Post List Organization

| Option | Description | Selected |
|--------|-------------|----------|
| Table with status tabs | Table (title, status, date, actions) with All/Published/Drafts/Deleted filter tabs. | ✓ |
| Simple list, no filtering | Flat list with status badges. | |

**User's choice:** Table with status tabs (Recommended)
**Notes:** None

---

### Soft-Delete Recovery

| Option | Description | Selected |
|--------|-------------|----------|
| Restore button in Deleted tab | Deleted posts appear in "Deleted" tab with "Restore" button. Restores to draft. | ✓ |
| Undo toast after delete | Brief "Undo" toast/banner with ~10s timer. | |

**User's choice:** Restore button in Deleted tab (Recommended)
**Notes:** None

---

### Slug Behavior

| Option | Description | Selected |
|--------|-------------|----------|
| Auto-generate, editable | Slug auto-generates from title, admin can manually edit before saving. | ✓ |
| Always manual | Admin must type slug for every post. | |

**User's choice:** Auto-generate, editable (Recommended)
**Notes:** None

---

### Action Confirmations

| Option | Description | Selected |
|--------|-------------|----------|
| Confirm delete only | Delete requires "Are you sure?" dialog. Publish/unpublish act immediately. | ✓ |
| Confirm all destructive actions | Both delete and unpublish require confirmation. | |

**User's choice:** Confirm delete only (Recommended)
**Notes:** None

---

## Claude's Discretion

- CSRF implementation approach (Go 1.26 stdlib check first)
- Session middleware design
- Admin template structure
- Table styling and responsive behavior
- Login form layout
- Slug generation algorithm
- Debounce timing for editor preview
- Rate limiter implementation

## Deferred Ideas

None — discussion stayed within phase scope.
