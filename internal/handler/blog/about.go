package blog

import (
	"net/http"

	"github.com/jared-wallace/website-go/content"
	"github.com/jared-wallace/website-go/internal/markdown"
)

// AboutPage handles GET /about. It renders the embedded about.md through
// the goldmark pipeline and serves the result inside the nautical base template.
func (h *BlogHandler) AboutPage(w http.ResponseWriter, r *http.Request) {
	renderer := markdown.NewRenderer()
	rendered := renderer.Render(content.AboutMD)
	h.render(w, http.StatusOK, "about.html", map[string]interface{}{
		"RenderedHTML": rendered,
	})
}
