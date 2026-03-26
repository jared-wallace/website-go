package blog

import (
	"errors"
	"log/slog"
	"net/http"

	postservice "github.com/jared-wallace/website-go/internal/service/post"
)

// ShowPost handles GET /posts/{slug}.
// It fetches the post by slug and renders the single-post template with the
// table of contents, rendered HTML body, and reading time. Returns a themed
// 404 page when the slug does not match any published post.
func (h *BlogHandler) ShowPost(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		h.NotFound(w, r)
		return
	}

	detail, err := h.svc.GetBySlug(r.Context(), slug)
	if err != nil {
		if errors.Is(err, postservice.ErrNotFound) {
			h.NotFound(w, r)
			return
		}
		slog.Error("GetBySlug failed", "slug", slug, "error", err)
		h.render(w, http.StatusInternalServerError, "404.html", map[string]interface{}{
			"Error": "Something went sideways. Try refreshing — if the problem persists, the sea may have claimed it.",
		})
		return
	}

	h.render(w, http.StatusOK, "post.html", map[string]interface{}{
		"Post":         detail.Post,
		"RenderedHTML": detail.RenderedHTML,
		"ToC":          detail.ToC,
		"Tags":         detail.Tags,
		"ReadingTime":  detail.ReadingTime,
	})
}
