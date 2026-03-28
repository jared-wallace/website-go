package admin

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	postrepo "github.com/jared-wallace/website-go/internal/repository/post"
	postservice "github.com/jared-wallace/website-go/internal/service/post"
)

// NewPost renders the blank editor form for creating a new post.
func (h *AdminHandler) NewPost(w http.ResponseWriter, r *http.Request) {
	h.render(w, http.StatusOK, "admin-editor.html", map[string]interface{}{
		"IsNew":  true,
		"Action": "/admin/posts/new",
	})
}

// EditPost loads an existing post by ID and renders the editor pre-populated with its data.
func (h *AdminHandler) EditPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.sessions.Put(r.Context(), "flash_error", "Invalid post ID.")
		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
		return
	}

	post, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		h.sessions.Put(r.Context(), "flash_error", "Post not found.")
		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
		return
	}

	// Pass RenderedHTML as template.HTML so the preview pane shows it unescaped.
	type postView struct {
		ID           int64
		Title        string
		Slug         string
		Body         string
		Tags         string
		Published    bool
		RenderedHTML template.HTML
	}
	pv := postView{
		ID:           post.ID,
		Title:        post.Title,
		Slug:         post.Slug,
		Body:         post.Body,
		Tags:         post.Tags,
		Published:    post.Published,
		RenderedHTML: template.HTML(post.RenderedHTML), //nolint:gosec // goldmark+bluemonday pipeline output
	}

	h.render(w, http.StatusOK, "admin-editor.html", map[string]interface{}{
		"IsNew":  false,
		"Action": fmt.Sprintf("/admin/posts/%d/edit", id),
		"Post":   pv,
	})
}

// SavePost handles form submission for both new-post creation and existing-post updates.
// The hidden "action" field distinguishes "draft" from "publish".
func (h *AdminHandler) SavePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.sessions.Put(r.Context(), "flash_error", "Failed to parse form.")
		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
		return
	}

	title := r.FormValue("title")
	slug := r.FormValue("slug")
	body := r.FormValue("body")
	tags := r.FormValue("tags")
	action := r.FormValue("action") // "draft" or "publish"

	if slug == "" {
		slug = postservice.GenerateSlug(title)
	}

	publish := action == "publish"

	// Determine new vs. edit from the URL path value.
	rawID := r.PathValue("id")
	if rawID == "" || rawID == "new" {
		// Create path: POST /admin/posts/new
		created, err := h.svc.Create(r.Context(), title, slug, body, tags, publish)
		if err != nil {
			if errors.Is(err, postrepo.ErrSlugExists) {
				h.sessions.Put(r.Context(), "flash_error", "Slug already in use -- please edit the slug field.")
				h.render(w, http.StatusOK, "admin-editor.html", map[string]interface{}{
					"IsNew":  true,
					"Action": "/admin/posts/new",
					"Form": map[string]string{
						"Title": title,
						"Slug":  slug,
						"Body":  body,
						"Tags":  tags,
					},
					"FlashError": "Slug already in use -- please edit the slug field.",
				})
				return
			}
			h.sessions.Put(r.Context(), "flash_error", "Failed to save post.")
			http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
			return
		}

		if publish {
			h.sessions.Put(r.Context(), "flash_success", "Post published.")
		} else {
			h.sessions.Put(r.Context(), "flash_success", "Draft saved.")
		}
		http.Redirect(w, r, fmt.Sprintf("/admin/posts/%d/edit", created.ID), http.StatusSeeOther)
		return
	}

	// Edit path: POST /admin/posts/{id}/edit
	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		h.sessions.Put(r.Context(), "flash_error", "Invalid post ID.")
		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
		return
	}

	if err := h.svc.Update(r.Context(), id, title, slug, body, tags); err != nil {
		h.sessions.Put(r.Context(), "flash_error", "Failed to update post.")
		http.Redirect(w, r, fmt.Sprintf("/admin/posts/%d/edit", id), http.StatusSeeOther)
		return
	}

	if publish {
		if err := h.svc.Publish(r.Context(), id); err != nil {
			h.sessions.Put(r.Context(), "flash_error", "Post updated but publish failed.")
			http.Redirect(w, r, fmt.Sprintf("/admin/posts/%d/edit", id), http.StatusSeeOther)
			return
		}
		h.sessions.Put(r.Context(), "flash_success", "Post published.")
	} else {
		h.sessions.Put(r.Context(), "flash_success", "Draft saved.")
	}

	http.Redirect(w, r, fmt.Sprintf("/admin/posts/%d/edit", id), http.StatusSeeOther)
}
