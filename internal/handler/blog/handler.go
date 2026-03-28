// Package blog provides HTTP handlers for the public-facing blog.
// It connects the post service layer to HTML templates and implements
// the request-response cycle for listing posts, viewing a single post,
// and rendering the themed 404 page.
package blog

import (
	"html/template"
	"log/slog"
	"net/http"
	"time"

	postservice "github.com/jared-wallace/website-go/internal/service/post"
	"github.com/jared-wallace/website-go/web"
)

// BlogHandler holds parsed template sets (one per page) and the post service.
// Templates are parsed once at startup and cached. Both fields are
// goroutine-safe after construction.
type BlogHandler struct {
	svc       *postservice.Service
	funcMap   template.FuncMap
	templates map[string]*template.Template
}

// New parses all HTML templates from the embedded FS, registers template
// helper functions, and returns a ready-to-use BlogHandler. Each page gets
// its own template set (base + page file) to prevent block name collisions
// across pages. Panics if any template fails to parse — a parse failure is a
// programmer error that should surface at startup, not at request time.
func New(svc *postservice.Service) *BlogHandler {
	funcMap := template.FuncMap{
		"formatDate": func(t time.Time) string {
			return t.Format("January 2, 2006")
		},
		"prevPage": func(page int) int { return page - 1 },
		"nextPage": func(page int) int { return page + 1 },
		"pages": func(total int) []int {
			p := make([]int, total)
			for i := range p {
				p[i] = i + 1
			}
			return p
		},
		"currentYear": func() int { return time.Now().Year() },
	}

	pages := []string{"list.html", "post.html", "404.html", "about.html"}
	templates := make(map[string]*template.Template, len(pages))
	for _, page := range pages {
		tmpl := template.Must(
			template.New("").Funcs(funcMap).ParseFS(
				web.Templates,
				"templates/base.html",
				"templates/"+page,
			),
		)
		templates[page] = tmpl
	}

	return &BlogHandler{
		svc:       svc,
		funcMap:   funcMap,
		templates: templates,
	}
}

// render writes an HTML response using the named template. It injects the
// current year for the footer copyright line and sets the Content-Type header.
// Errors during template execution are logged; the status code has already been
// sent so there is no way to return a different status at that point.
func (h *BlogHandler) render(w http.ResponseWriter, status int, page string, data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}
	data["Year"] = time.Now().Year()
	tmpl, ok := h.templates[page]
	if !ok {
		slog.Error("unknown template page", "page", page)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		slog.Error("template execution failed", "page", page, "error", err)
	}
}
