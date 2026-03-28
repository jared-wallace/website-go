# Phase 11: Template Changes - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-03-28
**Phase:** 11-template-changes
**Areas discussed:** Footer layout & voice, Rope divider design, List page hero, Nav restructuring

---

## Footer Layout & Voice

### Footer Nav Links

| Option | Description | Selected |
|--------|-------------|----------|
| About + RSS (Recommended) | Move About from top nav, keep RSS icon. Minimal and purposeful. | ✓ |
| About + RSS + Home | Add a Home link for users who scroll to bottom. | |
| About + RSS + Posts | Add a Posts link for direct access to post list. | |

**User's choice:** About + RSS
**Notes:** None

### Footer Tone

| Option | Description | Selected |
|--------|-------------|----------|
| Wry & weathered | "Built on salt air and stubbornness" style — beach bar vibe with a knowing wink. | ✓ |
| Warm & inviting | "Pull up a stool" style — welcoming, less edge. | |
| Terse & nautical | "Fair winds." style — short, atmospheric. | |

**User's choice:** Wry & weathered
**Notes:** None

### Phrase Authorship

| Option | Description | Selected |
|--------|-------------|----------|
| Claude picks | Claude crafts a wry, weathered one-liner. | ✓ |
| I have a phrase | User provides exact text. | |

**User's choice:** Claude picks
**Notes:** None

### Footer Layout

| Option | Description | Selected |
|--------|-------------|----------|
| Side-by-side (Recommended) | Nav links left, copyright right. Stacks on mobile. | ✓ |
| Stacked centered | Nav row on top, copyright below. Everything centered. | |
| Nav left, phrase + copyright right | Nav links left, personality phrase and copyright stacked right. | |

**User's choice:** Side-by-side
**Notes:** None

---

## Rope Divider Design

### Rope Style

| Option | Description | Selected |
|--------|-------------|----------|
| Twisted two-strand (Recommended) | Classic nautical dock line. Clean, recognizable, renders well small. | ✓ |
| Braided three-strand | Thicker, more textured. More complex SVG. | |
| Simple cord with knots | Thin single line with decorative knots. Lighter, more whimsical. | |

**User's choice:** Twisted two-strand
**Notes:** None

### Rope Color

| Option | Description | Selected |
|--------|-------------|----------|
| Match divider color (Recommended) | Use existing --color-divider CSS variable. | ✓ |
| Slightly warmer/tan | Rope-specific warm tone, adds a new color. | |
| Accent color | Use --color-accent. Stronger design element. | |

**User's choice:** Match divider color
**Notes:** None

### Rope Scope

| Option | Description | Selected |
|--------|-------------|----------|
| Footer only (Recommended) | Replace single footer .rope-divider. Conservative. | ✓ |
| All .rope-divider instances | Any element with class gets SVG. | |

**User's choice:** Footer only
**Notes:** None

---

## List Page Hero

### Hero Text

| Option | Description | Selected |
|--------|-------------|----------|
| Site name + tagline (Recommended) | "The Wild Meridian" h1 + "dispatches from the deep end" subtitle. | ✓ |
| Different heading | Distinct from nav tagline, e.g., "The Logbook". | |
| Claude picks | Claude crafts heading fitting nautical voice. | |

**User's choice:** Site name + tagline
**Notes:** None

### Hero Prominence

| Option | Description | Selected |
|--------|-------------|----------|
| Modest heading (Recommended) | Playfair Display h1, Lora italic tagline, standard spacing. | ✓ |
| Large hero section | Big type, generous padding, clear separation from cards. | |
| Minimal label | Small heading, tight spacing, almost a section label. | |

**User's choice:** Modest heading
**Notes:** None

---

## Nav Restructuring

### Nav Contents After About Removal

| Option | Description | Selected |
|--------|-------------|----------|
| Site name + tagline + dark toggle (Recommended) | Just remove About. Clean and minimal. | ✓ |
| Site name + dark toggle | Also remove tagline from nav. Ultra-minimal. | |
| Add Posts link | Replace About with Posts link to /posts. | |

**User's choice:** Site name + tagline + dark toggle
**Notes:** None

### Aria Labels

| Option | Description | Selected |
|--------|-------------|----------|
| "Main navigation" / "Footer navigation" (Recommended) | Standard, clear, expected by screen reader users. | ✓ |
| Branded labels | "Wild Meridian navigation" / "Wild Meridian footer". Distinctive but less clear. | |
| Claude picks | Claude chooses accessible, standards-compliant labels. | |

**User's choice:** "Main navigation" / "Footer navigation"
**Notes:** None

---

## Claude's Discretion

- Exact wording of the nautical personality phrase (wry/weathered tone, short)
- Exact CSS spacing and breakpoint for footer column stacking

## Deferred Ideas

None — discussion stayed within phase scope.
