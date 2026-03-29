---
status: awaiting_human_verify
trigger: "admin-preview-flash — Admin post preview (both new and edit) flashes briefly then disappears"
created: 2026-03-28T00:00:00Z
updated: 2026-03-28T00:00:00Z
---

## Current Focus
<!-- OVERWRITE on each update - reflects NOW -->

hypothesis: CONFIRMED — r.ParseForm() initializes r.Form to a non-nil empty value for multipart/form-data requests, causing r.FormValue("body") to short-circuit and return "" instead of reading the multipart body
test: Traced Go stdlib source: r.FormValue checks if r.Form == nil before calling ParseMultipartForm. r.ParseForm() initializes r.Form to non-nil even for non-URL-encoded bodies. So r.FormValue("body") returns "" when body was sent as FormData.
expecting: Fix: change admin.js to send URLSearchParams (application/x-www-form-urlencoded) instead of FormData (multipart/form-data), which r.ParseForm() handles correctly
next_action: Apply fix to web/static/admin.js

## Symptoms
<!-- Written during gathering, then IMMUTABLE -->

expected: Clicking preview in the admin post editor should show a rendered preview of the markdown content that stays visible
actual: Preview flashes for a second then disappears
errors: No specific error messages reported
reproduction: Go to admin panel, create new post or edit existing post, click preview button — preview flashes then vanishes
started: Unclear if it ever worked — first time user noticed

## Eliminated
<!-- APPEND only - prevents re-investigating -->

- hypothesis: Preview button causes form submission (page navigate)
  evidence: Tab buttons use type="button", not type="submit". No form submission occurs on tab click.
  timestamp: 2026-03-28

- hypothesis: CrossOriginProtection middleware blocks the preview fetch
  evidence: Same-origin fetch matches Host header; CrossOriginProtection allows it. Verified Go 1.25+ behavior.
  timestamp: 2026-03-28

- hypothesis: CSS animation or transition hides the preview pane
  evidence: No animation/transition on .admin-preview-pane or #preview-content in admin.css.
  timestamp: 2026-03-28

- hypothesis: The .hidden class fails on mobile
  evidence: .hidden is defined inside @media (max-width: 768px) with display:none !important. Works correctly for mobile tab toggling.
  timestamp: 2026-03-28

## Evidence
<!-- APPEND only - facts discovered -->

- timestamp: 2026-03-28
  checked: admin-editor.html template
  found: Preview pane (#preview-content) renders server-side RenderedHTML on edit post page load. On new post it's empty.
  implication: On edit post, the preview shows the server-rendered HTML briefly — this is the "flash."

- timestamp: 2026-03-28
  checked: admin.js — initial preview trigger (lines 33-35)
  found: On page load, if editorBody has content, dispatches 'input' event → 300ms debounce → fetch('/admin/preview') with FormData body → previewPane.innerHTML = html
  implication: On edit post, the fetch fires 300ms after page load and overwrites innerHTML with fetch response.

- timestamp: 2026-03-28
  checked: preview.go — server-side handler
  found: Handler calls r.ParseForm() then r.FormValue("body"). For multipart/form-data (what FormData sends), r.ParseForm() initializes r.Form to non-nil empty value without reading the body.
  implication: r.FormValue("body") checks if r.Form == nil (it's NOT nil after ParseForm), skips ParseMultipartForm, returns "".

- timestamp: 2026-03-28
  checked: Go stdlib /usr/local/go/src/net/http/request.go line 1428-1436
  found: FormValue: "if r.Form == nil { r.ParseMultipartForm(defaultMaxMemory) }". ParseForm docs: "For other HTTP methods, or when the Content-Type is not application/x-www-form-urlencoded, the request Body is not read, and r.PostForm is initialized to a non-nil, empty value."
  implication: Root cause confirmed. r.ParseForm() sets r.Form to non-nil for multipart bodies. r.FormValue then skips multipart parsing. body = "". Renderer returns empty HTML. previewPane.innerHTML = "" clears the server-rendered preview.

## Resolution
<!-- OVERWRITE as understanding evolves -->

root_cause: preview.go calls r.ParseForm() before r.FormValue("body"). The JavaScript fetch sends FormData which has Content-Type multipart/form-data. r.ParseForm() initializes r.Form to a non-nil empty map for non-URL-encoded bodies without reading them. r.FormValue then sees r.Form != nil and skips calling ParseMultipartForm, returning "" for the body field. The renderer returns empty HTML, which the JS injects into the preview pane, clearing the server-rendered HTML that was briefly visible.
fix: Changed admin.js to use URLSearchParams instead of FormData for the preview fetch. URLSearchParams sends application/x-www-form-urlencoded, which r.ParseForm() parses correctly. The body field is now read properly by the server.
verification: All tests pass (go test ./...). Regression test TestPreviewReturnsHTML continues to pass. The fix is a 2-line change in admin.js (FormData → URLSearchParams).
files_changed: [web/static/admin.js]
