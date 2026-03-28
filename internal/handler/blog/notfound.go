package blog

import "net/http"

// NotFound renders the themed 404 page with HTTP 404 status.
// The status code is set via the render helper before writing the body —
// Go's net/http requires headers (including status) to be written before body.
func (h *BlogHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	h.render(w, http.StatusNotFound, "404.html", nil)
}
