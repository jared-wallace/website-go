# Phase 4: Distribution - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-03-27
**Phase:** 04-distribution
**Areas discussed:** RSS feed shape, OG & social previews, Thumbs-up reactions, Sitemap scope

---

## RSS Feed Shape

### Feed Size

| Option | Description | Selected |
|--------|-------------|----------|
| All published posts | Every published post in the feed — simple, no pagination concern | |
| Most recent 25 | Keeps the feed lightweight. Standard practice for personal blogs | ✓ |
| Most recent 50 | Generous window — covers ~a year of weekly posting | |

**User's choice:** Most recent 25
**Notes:** None

### Author Name

| Option | Description | Selected |
|--------|-------------|----------|
| Jared Wallace | Real name in managingEditor and per-item author fields | ✓ |
| The Log | Use the blog name as author | |
| You decide | Claude picks | |

**User's choice:** Jared Wallace
**Notes:** None

### Tags in Feed

| Option | Description | Selected |
|--------|-------------|----------|
| Yes, include tags | Map each tag to an RSS category element | ✓ |
| No, skip tags | Keep feed items minimal | |

**User's choice:** Include tags as category elements
**Notes:** None

### Auto-Discovery

| Option | Description | Selected |
|--------|-------------|----------|
| Yes, in base.html | Standard practice — feed readers find it on any page | ✓ |
| Only on homepage | Lighter touch | |

**User's choice:** In base.html (all pages)
**Notes:** None

---

## OG & Social Previews

### OG Description Source

| Option | Description | Selected |
|--------|-------------|----------|
| Post excerpt | Reuse existing excerpt from card listings | ✓ |
| First 160 chars of body | Raw truncation — simpler but might cut mid-sentence | |
| You decide | Claude picks | |

**User's choice:** Post excerpt
**Notes:** None

### OG Image Strategy

| Option | Description | Selected |
|--------|-------------|----------|
| Static site-wide fallback | Branded nautical graphic used on all pages until Phase 5 | ✓ |
| No OG image for now | Skip og:image entirely | |
| Generate placeholder per post | Programmatic OG images with post title text | |

**User's choice:** Static site-wide fallback
**Notes:** None

### Twitter Card Type

| Option | Description | Selected |
|--------|-------------|----------|
| summary | Small square image + title + description. Standard for text content | ✓ |
| summary_large_image | Wide banner image. Looks better but static fallback would be repeated | |
| You decide | Claude picks | |

**User's choice:** summary
**Notes:** None

### Homepage OG Tags

| Option | Description | Selected |
|--------|-------------|----------|
| Yes, site-level OG | og:title="The Log", og:description="dispatches from the deep end" | ✓ |
| Posts only | Only add OG tags on individual post pages | |

**User's choice:** Site-level OG on homepage
**Notes:** None

---

## Thumbs-Up Reactions

### Button Placement

| Option | Description | Selected |
|--------|-------------|----------|
| Below post content | After article body, before footer. Natural "I finished reading" spot | ✓ |
| Floating sidebar | Fixed position on side, always visible while scrolling | |
| Both top and bottom | Count in metadata area + clickable button at bottom | |

**User's choice:** Below post content
**Notes:** None

### Tap Feedback

| Option | Description | Selected |
|--------|-------------|----------|
| Count increment + animation | Number ticks up with subtle CSS bounce or color fill | ✓ |
| Just count increment | Number changes, no animation | |
| You decide | Claude picks appropriate micro-interaction | |

**User's choice:** Count increment + brief CSS animation
**Notes:** None

### Rate Limiting Approach

| Option | Description | Selected |
|--------|-------------|----------|
| IP-based server-side | One per IP per post per 24h in reactions table | |
| IP + localStorage double-check | Server-side IP limit + localStorage flag for better UX | ✓ |
| You decide | Claude picks | |

**User's choice:** IP + localStorage double-check
**Notes:** Server-side is source of truth; localStorage improves UX by keeping button in "liked" state

### Multi-Tap Behavior

| Option | Description | Selected |
|--------|-------------|----------|
| One per reader | Binary like. One tap per IP per post | ✓ |
| Multiple claps (up to N) | Reader can tap multiple times for enthusiasm | |

**User's choice:** One per reader
**Notes:** None

---

## Sitemap Scope

### Pages Included

| Option | Description | Selected |
|--------|-------------|----------|
| Posts + homepage | All published post URLs plus root URL | ✓ |
| Posts only | Only /posts/{slug} URLs | |
| You decide | Claude includes whatever makes sense | |

**User's choice:** Posts + homepage
**Notes:** None

### robots.txt

| Option | Description | Selected |
|--------|-------------|----------|
| Yes, add robots.txt | Handler returning Sitemap directive + standard Allow | ✓ |
| No robots.txt | Skip it | |
| You decide | Claude handles discoverability | |

**User's choice:** Add robots.txt
**Notes:** None

---

## Claude's Discretion

- RSS channel metadata (title, link, description)
- OG fallback image design and dimensions
- Thumbs-up icon and CSS animation specifics
- Sitemap changefreq/priority values
- Reactions table schema and IP hashing approach
- robots.txt additional directives

## Deferred Ideas

None — discussion stayed within phase scope.
