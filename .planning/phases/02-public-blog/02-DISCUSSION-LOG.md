# Phase 2: Public Blog - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-03-26
**Phase:** 02-public-blog
**Areas discussed:** Post listing layout, Pagination style, Nautical design details, Table of Contents

---

## Post Listing Layout

| Option | Description | Selected |
|--------|-------------|----------|
| Card grid | 2-col responsive grid with title, date, reading time, excerpt | ✓ |
| Compact list | Vertical list, title + date per row, no excerpts | |
| Full preview | Full first paragraph per post, ~200 words | |

**User's choice:** Card grid
**Notes:** Fits the "chalkboard menu" aesthetic naturally.

### Excerpt Length

| Option | Description | Selected |
|--------|-------------|----------|
| 2-3 lines | Compact, uniform card height | ✓ |
| First paragraph | ~150 chars, variable height | |
| You decide | Claude picks | |

**User's choice:** 2-3 lines

### Card Metadata

| Option | Description | Selected |
|--------|-------------|----------|
| Published date | Required by BLOG-04 | ✓ |
| Reading time | Required by BLOG-04 | ✓ |
| Post tags/categories | Visual labels, needs new field | ✓ |
| Word count | Raw count | |

**User's choice:** Published date, Reading time, Post tags/categories
**Notes:** Tags not in current Post model — will need new column + migration.

### Tag Implementation

| Option | Description | Selected |
|--------|-------------|----------|
| Simple text field | Comma-separated TEXT column, visual labels only | ✓ |
| Clickable filter tags | Links to /tags/{tag} filtered views | |
| You decide | Claude picks simplest | |

**User's choice:** Simple text field

---

## Pagination Style

| Option | Description | Selected |
|--------|-------------|----------|
| Numbered pages | Traditional /posts?page=2, SEO-friendly, no JS | ✓ |
| Older/Newer buttons | Simple prev/next, no page jumping | |
| No pagination yet | All posts on one page | |

**User's choice:** Numbered pages

### Page Size

| Option | Description | Selected |
|--------|-------------|----------|
| 10 per page | Standard blog convention | ✓ |
| 6 per page | Fewer scrolling, more clicking | |
| You decide | Claude picks | |

**User's choice:** 10 per page

### Home Page

| Option | Description | Selected |
|--------|-------------|----------|
| Posts listing directly | Root URL = card grid, name/tagline header | ✓ |
| Landing hero + posts | Hero section then cards below | |
| You decide | Claude picks | |

**User's choice:** Posts listing directly

---

## Nautical Design Details

### Theme Depth

| Option | Description | Selected |
|--------|-------------|----------|
| Tasteful accents | Nautical palette + subtle decorative elements | ✓ |
| Full immersion | Heavy textures, patterns, decorative borders throughout | |
| Colors only | Palette only, no decorative nautical elements | |

**User's choice:** Tasteful accents

### Typography

| Option | Description | Selected |
|--------|-------------|----------|
| System fonts + chalkboard headers | System body, chalk-style heading font | |
| Serif body + sans headers | Classic blog readability | |
| You decide | Claude picks | ✓ |

**User's choice:** You decide

### Dark Mode Palette

| Option | Description | Selected |
|--------|-------------|----------|
| Night beach | Deep navy/charcoal, sand-gold accents | ✓ |
| Inverted palette | Straightforward light-on-dark inversion | |
| You decide | Claude picks | |

**User's choice:** Night beach

### Dark Mode Toggle

| Option | Description | Selected |
|--------|-------------|----------|
| Nav bar corner | Sun/moon icon top-right, localStorage persistence | ✓ |
| Footer | Less prominent, relies on prefers-color-scheme | |
| You decide | Claude picks | |

**User's choice:** Nav bar corner

---

## Table of Contents

### Placement

| Option | Description | Selected |
|--------|-------------|----------|
| Inline before content | Collapsible block after title, before content | ✓ |
| Sticky sidebar | Fixed right sidebar, highlights current section | |
| You decide | Claude picks | |

**User's choice:** Inline before content

### Heading Depth

| Option | Description | Selected |
|--------|-------------|----------|
| h2 and h3 | Sections + subsections | ✓ |
| h2 only | Major sections only | |
| You decide | Claude picks | |

**User's choice:** h2 and h3

### Trigger Threshold

| Option | Description | Selected |
|--------|-------------|----------|
| 3+ headings | Show when 3+ h2/h3 headings exist | ✓ |
| Always show | Every post gets ToC | |
| You decide | Claude picks | |

**User's choice:** 3+ headings

---

## Claude's Discretion

- Typography choices (font families, sizes, line heights)
- Exact color hex values within the nautical palette
- Card hover effects and transitions
- 404 page design and copy
- Reading time calculation formula
- Excerpt extraction approach
- Pagination component styling
- Nav bar layout and content

## Deferred Ideas

None — discussion stayed within phase scope.
