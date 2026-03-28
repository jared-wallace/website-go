package admin

import (
	"net/http"
	"strconv"
)

// parsePostID extracts the {id} path value and returns an error message on failure.
func parsePostID(r *http.Request) (int64, bool) {
	raw := r.PathValue("id")
	id, err := strconv.ParseInt(raw, 10, 64)
	return id, err == nil
}

// DeletePost soft-deletes a post by ID, then redirects to the Deleted tab.
func (h *AdminHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	id, ok := parsePostID(r)
	if !ok {
		h.sessions.Put(r.Context(), "flash_error", "Something went wrong. Try again.")
		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
		return
	}
	if err := h.svc.SoftDelete(r.Context(), id); err != nil {
		h.sessions.Put(r.Context(), "flash_error", "Something went wrong. Try again.")
		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
		return
	}
	h.sessions.Put(r.Context(), "flash_success", "Post deleted. Restore it from the Deleted tab.")
	http.Redirect(w, r, "/admin/posts?filter=deleted", http.StatusSeeOther)
}

// RestorePost clears the deleted state and sends the post back to drafts.
func (h *AdminHandler) RestorePost(w http.ResponseWriter, r *http.Request) {
	id, ok := parsePostID(r)
	if !ok {
		h.sessions.Put(r.Context(), "flash_error", "Something went wrong. Try again.")
		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
		return
	}
	if err := h.svc.Restore(r.Context(), id); err != nil {
		h.sessions.Put(r.Context(), "flash_error", "Something went wrong. Try again.")
		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
		return
	}
	h.sessions.Put(r.Context(), "flash_success", "Post restored to drafts.")
	http.Redirect(w, r, "/admin/posts?filter=drafts", http.StatusSeeOther)
}

// PublishPost makes a post publicly visible on the blog.
func (h *AdminHandler) PublishPost(w http.ResponseWriter, r *http.Request) {
	id, ok := parsePostID(r)
	if !ok {
		h.sessions.Put(r.Context(), "flash_error", "Something went wrong. Try again.")
		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
		return
	}
	if err := h.svc.Publish(r.Context(), id); err != nil {
		h.sessions.Put(r.Context(), "flash_error", "Something went wrong. Try again.")
		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
		return
	}
	h.sessions.Put(r.Context(), "flash_success", "Post published.")
	http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
}

// UnpublishPost hides a post from the public blog without deleting it.
func (h *AdminHandler) UnpublishPost(w http.ResponseWriter, r *http.Request) {
	id, ok := parsePostID(r)
	if !ok {
		h.sessions.Put(r.Context(), "flash_error", "Something went wrong. Try again.")
		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
		return
	}
	if err := h.svc.Unpublish(r.Context(), id); err != nil {
		h.sessions.Put(r.Context(), "flash_error", "Something went wrong. Try again.")
		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
		return
	}
	h.sessions.Put(r.Context(), "flash_success", "Post unpublished.")
	http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
}
