# Roadmap: website-go

## Milestones

- ✅ **v1.0 MVP** -- Phases 1-6 (shipped 2026-03-28) | [Archive](milestones/v1.0-ROADMAP.md)
- 🚧 **v1.1 The Wild Meridian** -- Phases 7-8 (in progress)

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

### 🚧 v1.1 The Wild Meridian (In Progress)

**Milestone Goal:** Rebrand the blog identity to "The Wild Meridian", add an about page, and surface RSS discoverability in the navigation.

- [ ] **Phase 7: Rebrand + Navigation** - Rename all "The Log" references, update copyright, and add RSS discoverability
- [ ] **Phase 8: About Page** - New about route rendering a static markdown file with full nautical design integration

## Phase Details

### Phase 7: Rebrand + Navigation
**Goal**: Every public surface reflects the "The Wild Meridian" identity and readers can discover the RSS feed
**Depends on**: Phase 6
**Requirements**: BRAND-01, BRAND-02, BRAND-03, BRAND-04, NAV-01, NAV-02
**Success Criteria** (what must be TRUE):
  1. Every public page header and browser tab shows "The Wild Meridian" instead of "The Log"
  2. The RSS feed title, description, and channel metadata read "The Wild Meridian"
  3. Open Graph meta tags and the XML sitemap reference "The Wild Meridian"
  4. The copyright footer reads "2026 Jared Wallace"
  5. A discreet RSS link or icon is visible in the nav or footer, and the HTML head includes `<link rel="alternate" type="application/rss+xml">`
**Plans:** 1 plan
Plans:
- [ ] 07-01-PLAN.md -- Rebrand all strings + add RSS icon to footer

### Phase 8: About Page
**Goal**: Readers can navigate to and read a personal about page that matches the site's nautical design
**Depends on**: Phase 7
**Requirements**: ABOUT-01, ABOUT-02, ABOUT-03
**Success Criteria** (what must be TRUE):
  1. A working "About" link appears in the main site navigation
  2. The about page renders content from a static markdown file on disk (not a database record)
  3. The about page shares the full nautical chrome: header, footer, dark mode toggle, and consistent typography
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
| 7. Rebrand + Navigation | v1.1 | 0/1 | In progress | - |
| 8. About Page | v1.1 | 0/? | Not started | - |
