---
phase: quick
plan: 260328-tea
type: execute
wave: 1
depends_on: []
files_modified:
  - web/templates/about.html
  - web/static/main.css
  - internal/handler/blog/about_test.go
autonomous: true
requirements: []
must_haves:
  truths:
    - "About page displays exactly one h1 heading (from markdown, not template)"
    - "No orphaned CSS rules for removed template elements"
    - "All existing about page tests pass"
  artifacts:
    - path: "web/templates/about.html"
      provides: "About page template without duplicate h1 or stale hr"
    - path: "internal/handler/blog/about_test.go"
      provides: "Updated test assertions matching new template"
  key_links:
    - from: "content/about.md"
      to: "web/templates/about.html"
      via: "{{.RenderedHTML}} renders markdown h1"
      pattern: "RenderedHTML"
---

<objective>
Fix duplicate "About" heading on the about page. The template (`about.html` line 8) renders
`<h1 class="about-title">About</h1>` AND the markdown (`content/about.md` line 1) starts with
`# About` which goldmark also renders as `<h1>`. Remove the template h1 and the orphaned
`<hr class="rope-divider">` below it (line 9) — the markdown heading is the single source of truth.

Purpose: Eliminate the double-heading visual bug.
Output: Clean about page with a single h1 from markdown.
</objective>

<execution_context>
@$HOME/.claude/get-shit-done/workflows/execute-plan.md
@$HOME/.claude/get-shit-done/templates/summary.md
</execution_context>

<context>
@web/templates/about.html
@web/static/main.css
@internal/handler/blog/about_test.go
@content/about.md
</context>

<tasks>

<task type="auto">
  <name>Task 1: Remove duplicate heading and stale hr from about template, clean up CSS and test</name>
  <files>web/templates/about.html, web/static/main.css, internal/handler/blog/about_test.go</files>
  <action>
1. **web/templates/about.html** — Remove lines 8 and 9 (the `<h1 class="about-title">About</h1>` and `<hr class="rope-divider">`). The resulting template content block should be:
```
<div class="container container-narrow">
  <article class="about-page">
    <div class="post-body">{{.RenderedHTML}}</div>
  </article>
</div>
```

2. **web/static/main.css** — Remove the `.about-title` rule block (lines 366-372). It targeted the template h1 that no longer exists. The `/* --- About Page --- */` section comment can stay if other about-page styles exist below, or be removed if the section is now empty.

3. **internal/handler/blog/about_test.go** — The test `"body contains about-title CSS class"` (lines 25-33) asserts a class that no longer exists. Replace it with a test that checks the markdown-rendered heading is present: assert `strings.Contains(body, "&lt;h1&gt;About&lt;/h1&gt;")` (the goldmark-rendered h1 from the markdown `# About`). Rename the subtest to `"body contains About heading from markdown"`.
  </action>
  <verify>
    <automated>cd /Users/jaredwallace/src/jared-wallace/website-go && go test ./internal/handler/blog/ -run TestAboutPage -v</automated>
  </verify>
  <done>
    - about.html has exactly zero `<h1` tags (markdown provides the h1 via RenderedHTML)
    - about.html has no `<hr class="rope-divider">` element
    - main.css has no `.about-title` rule
    - All TestAboutPage subtests pass
  </done>
</task>

</tasks>

<verification>
```bash
# No duplicate h1 in template
grep -c '<h1' web/templates/about.html  # expect 0

# No stale rope-divider hr in about template
grep -c 'rope-divider' web/templates/about.html  # expect 0

# No orphaned CSS
grep -c 'about-title' web/static/main.css  # expect 0

# Tests pass
go test ./internal/handler/blog/ -run TestAboutPage -v
```
</verification>

<success_criteria>
- About page renders a single h1 heading (from markdown)
- No `.about-title` CSS class in template or stylesheet
- No `<hr class="rope-divider">` in about template
- All about page tests pass
</success_criteria>

<output>
After completion, create `.planning/quick/260328-tea-fix-duplicate-about-heading-on-about-pag/260328-tea-SUMMARY.md`
</output>
