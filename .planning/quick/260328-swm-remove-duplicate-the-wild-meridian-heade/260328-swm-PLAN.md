---
phase: quick
plan: 260328-swm
type: execute
wave: 1
depends_on: []
files_modified:
  - web/templates/list.html
  - web/static/main.css
  - internal/handler/blog/handler_test.go
autonomous: true
requirements: []
must_haves:
  truths:
    - "The home page shows the site name and tagline exactly once (in the nav bar)"
    - "No duplicate hero section appears below the nav on the list page"
    - "All existing tests pass after the change"
  artifacts:
    - path: "web/templates/list.html"
      provides: "Post list without redundant hero"
    - path: "web/static/main.css"
      provides: "CSS without orphaned .list-hero rules"
    - path: "internal/handler/blog/handler_test.go"
      provides: "Updated test without hero assertions"
  key_links: []
---

<objective>
Remove the duplicate "The Wild Meridian" + "dispatches from the deep end" hero section from the home page. The nav bar in base.html (lines 27-28) already displays the site name and tagline on every page, so the `.list-hero` div in list.html is redundant visual clutter.

Purpose: Eliminate the double-header on the landing page for a cleaner first impression.
Output: Three files updated, zero new classes, one less thing saying "The Wild Meridian" at the user.
</objective>

<execution_context>
@$HOME/.claude/get-shit-done/workflows/execute-plan.md
@$HOME/.claude/get-shit-done/templates/summary.md
</execution_context>

<context>
@web/templates/list.html
@web/templates/base.html
@web/static/main.css
@internal/handler/blog/handler_test.go
</context>

<tasks>

<task type="auto">
  <name>Task 1: Remove hero div from list.html and clean up CSS and tests</name>
  <files>web/templates/list.html, web/static/main.css, internal/handler/blog/handler_test.go</files>
  <action>
1. In `web/templates/list.html`: Delete lines 19-22 (the entire `.list-hero` div including the h1 and p tags). The `{{if not .Posts}}` block on what is currently line 23 should follow directly after the opening `<div class="container">` on line 18.

2. In `web/static/main.css`: Delete the three CSS rule blocks for `.list-hero` (line 782-783), `.list-hero-title` (lines 786-792), and `.list-hero-tagline` (lines 794-802). These classes are no longer referenced anywhere.

3. In `internal/handler/blog/handler_test.go`: Remove the entire `TestListHero` function (lines 307-327). This test asserts the presence of the hero section we are intentionally removing. The test would fail and is no longer valid.
  </action>
  <verify>
    <automated>cd /Users/jaredwallace/src/jared-wallace/website-go && go test ./internal/handler/blog/... -v -count=1 2>&1 | tail -20</automated>
  </verify>
  <done>
- list.html has no `.list-hero` div
- main.css has no `.list-hero*` rules
- handler_test.go has no TestListHero function
- All remaining tests pass
  </done>
</task>

</tasks>

<verification>
- `grep -r "list-hero" web/ internal/` returns no matches
- `go test ./...` passes (or at minimum `go test ./internal/handler/blog/...`)
</verification>

<success_criteria>
The home page renders the site name and tagline exactly once, via the nav bar. No orphaned CSS or failing tests remain.
</success_criteria>

<output>
After completion, create `.planning/quick/260328-swm-remove-duplicate-the-wild-meridian-heade/260328-swm-SUMMARY.md`
</output>
