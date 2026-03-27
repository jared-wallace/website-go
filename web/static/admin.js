// === Admin Panel JavaScript ===
// Debounced preview, Ctrl+S, slug auto-generation, mobile tab toggle

(function() {
  'use strict';

  // --- Debounced Preview (D-04, 300ms per UI-SPEC) ---
  var editorBody = document.getElementById('editor-body');
  var previewPane = document.getElementById('preview-content');

  if (editorBody && previewPane) {
    var previewTimer;
    editorBody.addEventListener('input', function() {
      clearTimeout(previewTimer);
      previewTimer = setTimeout(function() {
        var formData = new FormData();
        formData.append('body', editorBody.value);
        fetch('/admin/preview', { method: 'POST', body: formData })
          .then(function(r) {
            if (!r.ok) throw new Error('Preview failed');
            return r.text();
          })
          .then(function(html) {
            previewPane.innerHTML = html;
          })
          .catch(function() {
            previewPane.innerHTML = '<p class="preview-error">Preview unavailable -- check your connection.</p>';
          });
      }, 300);
    });

    // Trigger initial preview if body has content (edit mode)
    if (editorBody.value.trim()) {
      editorBody.dispatchEvent(new Event('input'));
    }
  }

  // --- Ctrl+S / Cmd+S Save Shortcut (D-05) ---
  var editorForm = document.getElementById('editor-form');
  var actionField = document.getElementById('editor-action');
  if (editorForm && actionField) {
    document.addEventListener('keydown', function(e) {
      if ((e.ctrlKey || e.metaKey) && e.key === 's') {
        e.preventDefault();
        actionField.value = 'draft';
        editorForm.submit();
      }
    });
  }

  // --- Slug Auto-Generation (D-14) ---
  var titleField = document.getElementById('editor-title');
  var slugField = document.getElementById('editor-slug');
  if (titleField && slugField) {
    var slugManuallyEdited = false;

    // If slug already has a value (edit mode), consider it manually set
    if (slugField.value.trim()) {
      slugManuallyEdited = true;
    }

    slugField.addEventListener('input', function() {
      slugManuallyEdited = true;
    });

    titleField.addEventListener('input', function() {
      if (slugManuallyEdited) return;
      // Mirror GenerateSlug algorithm from slug.go
      var s = titleField.value.toLowerCase();
      s = s.replace(/[^a-z0-9]+/g, '-');
      s = s.replace(/^-+|-+$/g, '');
      slugField.value = s;
    });
  }

  // --- Mobile Tab Toggle (D-06) ---
  var writeTab = document.getElementById('tab-write');
  var previewTab = document.getElementById('tab-preview');
  var editorPane = document.getElementById('editor-pane');
  var previewPaneEl = document.getElementById('preview-pane');

  if (writeTab && previewTab && editorPane && previewPaneEl) {
    writeTab.addEventListener('click', function() {
      writeTab.classList.add('active');
      previewTab.classList.remove('active');
      editorPane.classList.remove('hidden');
      previewPaneEl.classList.add('hidden');
    });

    previewTab.addEventListener('click', function() {
      previewTab.classList.add('active');
      writeTab.classList.remove('active');
      previewPaneEl.classList.remove('hidden');
      editorPane.classList.add('hidden');
      // Trigger preview render when switching to preview tab
      if (editorBody && editorBody.value.trim()) {
        editorBody.dispatchEvent(new Event('input'));
      }
    });
  }
})();
