# Phase 9: CSS Foundation - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-03-28
**Phase:** 09-css-foundation
**Areas discussed:** Grain texture look, Post card shadows, Tag pill fill color

---

## Grain Texture Look

### Intensity

| Option | Description | Selected |
|--------|-------------|----------|
| Barely-there whisper | ~3-5% opacity. Subconscious warmth, disappears when reading. Like aged paper under glass. | ✓ |
| Visible patina | ~8-12% opacity. Clearly textured like a weathered wooden menu board. | |
| Heavy weathered | ~15-20% opacity. Strong visual texture like a photo filter. | |

**User's choice:** Barely-there whisper
**Notes:** None

### Tone

| Option | Description | Selected |
|--------|-------------|----------|
| Warm sand | Brownish-tan speckle reinforcing beach bar palette. Sun-bleached paper feel. | ✓ |
| Neutral monochrome | Black/white noise for texture without color temperature shift. Film grain feel. | |

**User's choice:** Warm sand
**Notes:** None

---

## Post Card Shadows

### Depth

| Option | Description | Selected |
|--------|-------------|----------|
| Gentle float | Cards sit slightly above page. Tight + ambient two-layer shadow. Hover lift still meaningful. | ✓ |
| Noticeable ledge | Cards clearly float. Like wooden plaques mounted on wall. Obvious layering. | |
| Flat with edge | Subtle bottom/right edge shadow only. Minimal depth, dramatic hover contrast. | |

**User's choice:** Gentle float
**Notes:** None

### Shadow Color

| Option | Description | Selected |
|--------|-------------|----------|
| Warm brown tint | Uses existing rgba(44, 36, 24, ...) in light mode. Neutral black in dark mode. | ✓ |
| Neutral gray | Standard rgba(0,0,0,...) shadows. Clean and conventional. | |

**User's choice:** Warm brown tint
**Notes:** None

---

## Tag Pill Fill Color

| Option | Description | Selected |
|--------|-------------|----------|
| Accent-tinted | Semi-transparent accent color (ocean blue light / gold dark). Subtle pop of signature color. | ✓ |
| Surface-derived | Slightly darker/lighter than card surface. Understated, like subtle indentations. | |
| Warm sand tone | Sandy parchment fill consistent in both modes. Weathered label feel. | |

**User's choice:** Accent-tinted
**Notes:** None

---

## Claude's Discretion

- CSS technique for grain generation (SVG noise vs CSS gradients vs base64 PNG)
- Exact opacity values within decided ranges
- Shadow transition timing details
- Tag pill padding adjustments after fill addition

## Deferred Ideas

None — discussion stayed within phase scope
