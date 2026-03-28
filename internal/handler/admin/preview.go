package admin

import (
	"net/http"
)

// Preview renders markdown body to HTML and returns it as a raw HTML fragment.
// It writes directly to the response without a template wrapper to avoid
// double-encoding the already-sanitized goldmark+bluemonday output.
func (h *AdminHandler) Preview(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB limit to prevent memory exhaustion
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	body := r.FormValue("body")
	rendered := h.renderer.Render(body)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(rendered)) //nolint:gosec // output is pre-sanitized by goldmark+bluemonday pipeline
}
