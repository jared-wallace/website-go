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

  // Thumbs-up reaction
  var reactionBtn = document.getElementById('reaction-btn');
  if (reactionBtn) {
    var slug = reactionBtn.getAttribute('data-slug');
    var storageKey = 'reacted:' + slug;

    if (localStorage.getItem(storageKey)) {
      reactionBtn.disabled = true;
      reactionBtn.classList.add('reacted');
    }

    reactionBtn.addEventListener('click', function() {
      if (reactionBtn.disabled) return;
      reactionBtn.disabled = true;

      fetch('/posts/' + slug + '/react', { method: 'POST' })
        .then(function(res) { return res.json(); })
        .then(function(data) {
          document.getElementById('reaction-count').textContent = data.count;
          reactionBtn.classList.add('reacted');
          reactionBtn.classList.add('bounce');
          reactionBtn.setAttribute('aria-label', 'You gave a thumbs-up');
          localStorage.setItem(storageKey, '1');
          var icon = reactionBtn.querySelector('.reaction-icon');
          if (icon) {
            icon.addEventListener('animationend', function() {
              reactionBtn.classList.remove('bounce');
            }, { once: true });
          }
        })
        .catch(function() {
          reactionBtn.disabled = false;
        });
    });
  }
})();
