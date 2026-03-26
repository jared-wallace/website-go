(function() {
  // Dark mode toggle
  var toggle = document.getElementById('dark-toggle');
  if (toggle) {
    toggle.addEventListener('click', function() {
      var html = document.documentElement;
      var current = html.getAttribute('data-theme');
      var next = current === 'dark' ? 'light' : 'dark';
      if (next === 'light') {
        html.removeAttribute('data-theme');
      } else {
        html.setAttribute('data-theme', 'dark');
      }
      localStorage.setItem('theme', next);
      // Update aria-label
      toggle.setAttribute('aria-label',
        next === 'dark' ? 'Switch to light mode' : 'Switch to dark mode');
    });
  }

  // ToC collapse
  var tocToggle = document.getElementById('toc-toggle');
  var tocList = document.getElementById('toc-list');
  if (tocToggle && tocList) {
    tocToggle.addEventListener('click', function() {
      var collapsed = tocList.classList.toggle('collapsed');
      tocToggle.textContent = collapsed ? '[ show ]' : '[ hide ]';
    });
  }
})();
