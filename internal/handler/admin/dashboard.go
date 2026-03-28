package admin

import (
	"net/http"
)

// Dashboard renders the post management table with optional filter tab active.
// It fetches all posts (including drafts and deleted) and filters in-memory
// so the tab counts always reflect the current database state.
func (h *AdminHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")

	posts, err := h.svc.ListAll(r.Context())
	if err != nil {
		h.sessions.Put(r.Context(), "flash_error", "Failed to load posts.")
		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
		return
	}

	// Filter in-memory per the active tab.
	filtered := posts[:0]
	for _, p := range posts {
		switch filter {
		case "published":
			if p.DeletedAt == nil && p.Published {
				filtered = append(filtered, p)
			}
		case "drafts":
			if p.DeletedAt == nil && !p.Published {
				filtered = append(filtered, p)
			}
		case "deleted":
			if p.DeletedAt != nil {
				filtered = append(filtered, p)
			}
		default:
			// "all" or "" — everything except deleted
			if p.DeletedAt == nil {
				filtered = append(filtered, p)
			}
		}
	}

	flashSuccess := h.sessions.PopString(r.Context(), "flash_success")
	flashError := h.sessions.PopString(r.Context(), "flash_error")

	h.render(w, http.StatusOK, "admin-dashboard.html", map[string]interface{}{
		"Posts":        filtered,
		"Filter":       filter,
		"FlashSuccess": flashSuccess,
		"FlashError":   flashError,
	})
}
