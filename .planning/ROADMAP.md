# Roadmap: website-go

## Milestones

- ✅ **v1.0 MVP** -- Phases 1-6 (shipped 2026-03-28) | [Archive](milestones/v1.0-ROADMAP.md)
- ✅ **v1.1 The Wild Meridian** -- Phases 7-8 (shipped 2026-03-28) | [Archive](milestones/v1.1-ROADMAP.md)
- 🚧 **v1.2 Shore Leave Polish** -- Phases 9-11 (in progress)

## Phases

<details>
<summary>✅ v1.0 MVP (Phases 1-6) -- SHIPPED 2026-03-28</summary>

- [x] Phase 1: Foundation (3/3 plans) -- completed 2026-03-26
- [x] Phase 2: Public Blog (3/3 plans) -- completed 2026-03-26
- [x] Phase 3: Admin Panel (4/4 plans) -- completed 2026-03-27
- [x] Phase 4: Distribution (3/3 plans) -- completed 2026-03-27
- [x] Phase 5: API + Images (2/2 plans) -- completed 2026-03-28
- [x] Phase 6: Docker + Deployment (2/2 plans) -- completed 2026-03-28

</details>

<details>
<summary>✅ v1.1 The Wild Meridian (Phases 7-8) -- SHIPPED 2026-03-28</summary>

- [x] Phase 7: Rebrand + Navigation (1/1 plans) -- completed 2026-03-28
- [x] Phase 8: About Page (1/1 plans) -- completed 2026-03-28

</details>

### 🚧 v1.2 Shore Leave Polish (In Progress)

**Milestone Goal:** Elevate visual cohesion and nautical atmosphere through CSS/template-level design polish — no backend changes.

- [ ] **Phase 9: CSS Foundation** - Pure CSS atmosphere and component fixes
- [ ] **Phase 10: Animations & Transitions** - Motion system with reduced-motion safety and dark mode transitions
- [ ] **Phase 11: Template Changes** - Footer redesign, rope divider, and homepage hero

## Phase Details

### Phase 9: CSS Foundation
**Goal**: The site's visual atmosphere is elevated through pure CSS changes — texture, depth, and legibility all improve without touching any HTML templates
**Depends on**: Phase 8 (existing codebase)
**Requirements**: ATMO-01, ATMO-02, TYPO-01, TYPO-02, HOUSE-01
**Success Criteria** (what must be TRUE):
  1. A subtle grain/noise texture is visible on the page background in both light and dark modes, and does not cause scroll lag on mobile
  2. Post cards display a warm two-layer resting shadow that gives them visual weight before any hover interaction
  3. Tag pills display a filled semi-transparent background that is legible in both light and dark modes
  4. Reaction button corners are visually consistent with the site's 4px design system (no pill shape)
  5. The CSS file header comment reads "The Wild Meridian" and no reference to "The Log" remains
**Plans:** 1 plan
Plans:
- [ ] 09-01-PLAN.md — CSS atmosphere and component fixes (grain, shadows, fills, radius, rebrand)
**UI hint**: yes

### Phase 10: Animations & Transitions
**Goal**: The site's motion system is complete, safe, and capable — entry animations delight users who prefer motion, dark mode transitions are smooth, and users who prefer reduced motion experience zero flashing or layout shift
**Depends on**: Phase 9
**Requirements**: ANIM-01, ANIM-02, ANIM-03, ANIM-04, ANIM-05
**Success Criteria** (what must be TRUE):
  1. Main content fades in on page load with a subtle opacity transition visible to sighted users
  2. Post cards on the list page stagger their entrance with sequential CSS animation delays
  3. Toggling dark mode produces a smooth color transition across all themed surfaces rather than an instant switch
  4. Dark mode does not flash white or produce a visible color transition on initial page load
  5. A user with `prefers-reduced-motion: reduce` set sees no animation, including the existing reaction bounce
**Plans**: TBD
**UI hint**: yes

### Phase 11: Template Changes
**Goal**: Navigation is restructured so the footer becomes a useful, personality-driven block with About and utility links, and the homepage and post list receive the visual improvements that require HTML changes
**Depends on**: Phase 10
**Requirements**: NAV-01, NAV-02, NAV-03, NAV-04, ATMO-03, TYPO-03
**Success Criteria** (what must be TRUE):
  1. The About link no longer appears in the primary navigation bar and is accessible from the footer
  2. The footer displays a two-section layout with navigation links and a copyright/utility area
  3. The footer contains a short nautical personality phrase that matches the site's voice
  4. Both the primary nav and the footer nav have distinct `aria-label` values so screen readers can distinguish landmarks
  5. Rope dividers render as a twisted rope SVG pattern rather than a dashed CSS border
  6. The list page displays a hero heading and tagline above the post card grid
**Plans**: TBD
**UI hint**: yes

## Progress

| Phase | Milestone | Plans Complete | Status | Completed |
|-------|-----------|----------------|--------|-----------|
| 1. Foundation | v1.0 | 3/3 | Complete | 2026-03-26 |
| 2. Public Blog | v1.0 | 3/3 | Complete | 2026-03-26 |
| 3. Admin Panel | v1.0 | 4/4 | Complete | 2026-03-27 |
| 4. Distribution | v1.0 | 3/3 | Complete | 2026-03-27 |
| 5. API + Images | v1.0 | 2/2 | Complete | 2026-03-28 |
| 6. Docker + Deployment | v1.0 | 2/2 | Complete | 2026-03-28 |
| 7. Rebrand + Navigation | v1.1 | 1/1 | Complete | 2026-03-28 |
| 8. About Page | v1.1 | 1/1 | Complete | 2026-03-28 |
| 9. CSS Foundation | v1.2 | 0/1 | Not started | - |
| 10. Animations & Transitions | v1.2 | 0/? | Not started | - |
| 11. Template Changes | v1.2 | 0/? | Not started | - |
