# Requirements: website-go

**Defined:** 2026-03-28
**Core Value:** A reader visits jared-wallace.com and reads well-rendered markdown blog posts with images in a distinctive, memorable design.

## v1.2 Requirements

Requirements for the Shore Leave Polish milestone. Each maps to roadmap phases.

### Visual Atmosphere

- [ ] **ATMO-01**: Site displays a subtle grain/noise texture overlay on the page background in both light and dark modes
- [ ] **ATMO-02**: Post cards display a warm two-layer resting shadow that provides depth context for the existing hover lift
- [ ] **ATMO-03**: Rope dividers render as an inline SVG twisted rope pattern replacing the current dashed border

### Navigation & Footer

- [ ] **NAV-01**: About link is removed from the navigation bar and appears in the footer
- [ ] **NAV-02**: Footer displays a two-section layout with navigation links and copyright/utility area
- [ ] **NAV-03**: Footer includes a short nautical personality phrase matching the site's voice
- [ ] **NAV-04**: Footer navigation has proper `aria-label` and primary nav gets a matching label for landmark disambiguation

### Animations & Transitions

- [ ] **ANIM-01**: All animations (existing `reaction-bounce` + new) are wrapped in `prefers-reduced-motion` guards
- [ ] **ANIM-02**: Main content fades in on page load with a subtle opacity transition
- [ ] **ANIM-03**: Post cards on the list page stagger their entrance with CSS animation delays
- [ ] **ANIM-04**: Dark mode toggle produces a smooth color transition across all themed surfaces via CSS `@property`
- [ ] **ANIM-05**: Dark mode transitions are gated behind a `.theme-ready` class added by JS post-load to prevent flash

### Typography & Consistency

- [ ] **TYPO-01**: Tag pills display with a filled semi-transparent background visible in both light and dark modes
- [ ] **TYPO-02**: Reaction button uses `border-radius: 4px` matching the site's design system
- [ ] **TYPO-03**: List page displays a hero heading/tagline area above the post card grid

### Housekeeping

- [ ] **HOUSE-01**: CSS file header comment reads "The Wild Meridian" instead of "The Log"

## Future Requirements

Deferred to future release. Tracked but not in current roadmap.

### Interactivity

- **INT-01**: Clickable tag filtering (requires backend routing changes)
- **INT-02**: Tag hover styles with WCAG 4.5:1 contrast for interactive elements

## Out of Scope

Explicitly excluded. Documented to prevent scope creep.

| Feature | Reason |
|---------|--------|
| Social media links in footer | Anti-feature for this design direction per research |
| JavaScript animation libraries | CSS-only philosophy; zero new dependencies for this milestone |
| Card rotation (alternating tilt) | Adds visual complexity without proportional atmosphere gain |
| New Go routes or handler changes | Pure CSS/template milestone; no backend work |
| Tag filtering backend | Requires routing changes; deferred to v1.3+ |

## Traceability

Which phases cover which requirements. Updated during roadmap creation.

| Requirement | Phase | Status |
|-------------|-------|--------|
| ATMO-01 | Phase 9 | Pending |
| ATMO-02 | Phase 9 | Pending |
| ATMO-03 | Phase 11 | Pending |
| NAV-01 | Phase 11 | Pending |
| NAV-02 | Phase 11 | Pending |
| NAV-03 | Phase 11 | Pending |
| NAV-04 | Phase 11 | Pending |
| ANIM-01 | Phase 10 | Pending |
| ANIM-02 | Phase 10 | Pending |
| ANIM-03 | Phase 10 | Pending |
| ANIM-04 | Phase 10 | Pending |
| ANIM-05 | Phase 10 | Pending |
| TYPO-01 | Phase 9 | Pending |
| TYPO-02 | Phase 9 | Pending |
| TYPO-03 | Phase 11 | Pending |
| HOUSE-01 | Phase 9 | Pending |

**Coverage:**
- v1.2 requirements: 16 total
- Mapped to phases: 16
- Unmapped: 0 ✓

---
*Requirements defined: 2026-03-28*
*Last updated: 2026-03-28 after roadmap creation*
