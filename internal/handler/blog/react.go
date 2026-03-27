package blog

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	postservice "github.com/jared-wallace/website-go/internal/service/post"
)

type reactResponse struct {
	Count          int  `json:"count"`
	AlreadyReacted bool `json:"already_reacted"`
}

// React handles POST /posts/{slug}/react.
// It records a thumbs-up reaction (one per IP per post) and returns the
// updated count as JSON.
func (h *BlogHandler) React(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.NotFound(w, r)
		return
	}

	detail, err := h.svc.GetBySlug(r.Context(), slug)
	if err != nil {
		if errors.Is(err, postservice.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		slog.Error("React: GetBySlug failed", "slug", slug, "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	ipHash := hashIP(clientIP(r))
	alreadyExists, err := h.svc.AddReaction(r.Context(), detail.Post.ID, ipHash)
	if err != nil {
		slog.Error("React: AddReaction failed", "slug", slug, "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	count, err := h.svc.CountReactions(r.Context(), detail.Post.ID)
	if err != nil {
		slog.Error("React: CountReactions failed", "slug", slug, "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reactResponse{ //nolint:errcheck
		Count:          count,
		AlreadyReacted: alreadyExists,
	})
}

// clientIP extracts the real client IP, preferring X-Real-IP (set by Nginx)
// and falling back to RemoteAddr for local development.
func clientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

// hashIP returns a SHA-256 hex digest of the IP address for privacy-preserving storage.
func hashIP(ip string) string {
	sum := sha256.Sum256([]byte(ip))
	return fmt.Sprintf("%x", sum)
}
