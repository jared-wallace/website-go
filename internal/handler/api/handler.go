// Package api provides the HTTP handler for the push-to-publish API.
// Posts are pushed as raw markdown with YAML front matter via bearer-token
// authenticated POST requests.
package api

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
)

// Upserter is the service-layer contract for creating or updating posts by slug.
type Upserter interface {
	UpsertBySlug(ctx context.Context, title, slug, body string) error
}

// MetaRenderer extracts YAML front matter metadata alongside HTML rendering.
type MetaRenderer interface {
	RenderWithMeta(src string) (template.HTML, map[string]interface{})
}

// APIHandler serves the push-to-publish API endpoint.
type APIHandler struct {
	svc      Upserter
	renderer MetaRenderer
}

// New creates an APIHandler with the given post service and markdown renderer.
func New(svc Upserter, renderer MetaRenderer) *APIHandler {
	return &APIHandler{svc: svc, renderer: renderer}
}

// PushPost accepts a raw markdown body with YAML front matter, extracts slug
// and title, then upserts the post as a draft. Requires a valid bearer token
// (enforced by upstream middleware).
func (h *APIHandler) PushPost(w http.ResponseWriter, r *http.Request) {
	const maxBody = 1 << 20 // 1 MB
	raw, err := io.ReadAll(io.LimitReader(r.Body, maxBody+1))
	if err != nil {
		slog.Error("api push: read body failed", "error", err)
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	if len(raw) > maxBody {
		http.Error(w, "body too large (max 1 MB)", http.StatusBadRequest)
		return
	}
	if len(raw) == 0 {
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	_, meta := h.renderer.RenderWithMeta(string(raw))

	slug, _ := meta["slug"].(string)
	if slug == "" {
		http.Error(w, "front matter must include 'slug'", http.StatusBadRequest)
		return
	}

	title, _ := meta["title"].(string)
	if title == "" {
		title = slug
	}

	if err := h.svc.UpsertBySlug(r.Context(), title, slug, string(raw)); err != nil {
		slog.Error("api push: upsert failed", "slug", slug, "error", err)
		http.Error(w, "failed to upsert post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "ok")
}
