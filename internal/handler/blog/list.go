package blog

import (
	"log/slog"
	"net/http"
	"strconv"
)

// ListPosts handles GET / and GET /posts?page=N.
// It fetches a page of published posts from the service and renders the list
// template with pagination metadata. Invalid or missing page params default to
// page 1; the service clamps out-of-range values.
func (h *BlogHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if n, err := strconv.Atoi(p); err == nil && n > 0 {
			page = n
		}
	}

	result, err := h.svc.ListPublished(r.Context(), page)
	if err != nil {
		slog.Error("ListPublished failed", "page", page, "error", err)
		h.render(w, http.StatusInternalServerError, "list.html", map[string]interface{}{
			"Error": "Something went sideways. Try refreshing — if the problem persists, the sea may have claimed it.",
		})
		return
	}

	h.render(w, http.StatusOK, "list.html", map[string]interface{}{
		"Posts":       result.Posts,
		"CurrentPage": result.CurrentPage,
		"TotalPages":  result.TotalPages,
		"HasPrev":     result.HasPrev,
		"HasNext":     result.HasNext,
	})
}
