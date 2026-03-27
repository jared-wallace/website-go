// Package admin provides HTTP handlers for the admin panel.
// All admin routes require authentication via SCS session middleware.
package admin

import (
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/jared-wallace/website-go/internal/config"
	"github.com/jared-wallace/website-go/internal/markdown"
	"github.com/jared-wallace/website-go/internal/middleware"
	postservice "github.com/jared-wallace/website-go/internal/service/post"
	"github.com/jared-wallace/website-go/web"
)

// AdminHandler holds template sets, session manager, and auth state.
// All fields are goroutine-safe after construction.
type AdminHandler struct {
	svc         *postservice.Service
	sessions    *scs.SessionManager
	renderer    *markdown.Renderer
	rateLimiter *middleware.RateLimiter
	funcMap     template.FuncMap
	templates   map[string]*template.Template
	adminEmail  string
	adminHash   []byte
	dummyHash   []byte // pre-computed hash for timing-safe comparison when email doesn't match
}

// New parses admin templates and returns a ready-to-use AdminHandler.
// It pre-computes a dummy bcrypt hash at startup to ensure constant-time
// credential checking even when the email doesn't match (Pitfall 6).
// Panics if any template fails to parse — a programmer error, not a runtime one.
func New(svc *postservice.Service, sm *scs.SessionManager, r *markdown.Renderer, rl *middleware.RateLimiter, cfg config.Config) *AdminHandler {
	funcMap := template.FuncMap{
		"formatDate":  func(t time.Time) string { return t.Format("January 2, 2006") },
		"currentYear": func() int { return time.Now().Year() },
	}

	pages := []string{"admin-login.html", "admin-dashboard.html", "admin-editor.html"}
	templates := make(map[string]*template.Template, len(pages))
	for _, page := range pages {
		tmpl := template.Must(
			template.New("").Funcs(funcMap).ParseFS(
				web.Templates,
				"templates/admin-base.html",
				"templates/"+page,
			),
		)
		templates[page] = tmpl
	}

	// Pre-compute a dummy hash for timing-safe login (constant time even on email miss).
	dummy, err := bcrypt.GenerateFromPassword([]byte("dummy-timing-safe"), 12)
	if err != nil {
		panic("admin: failed to generate dummy bcrypt hash: " + err.Error())
	}

	return &AdminHandler{
		svc:         svc,
		sessions:    sm,
		renderer:    r,
		rateLimiter: rl,
		funcMap:     funcMap,
		templates:   templates,
		adminEmail:  cfg.AdminEmail,
		adminHash:   []byte(cfg.AdminPasswordHash),
		dummyHash:   dummy,
	}
}

// render writes an HTML response using the named admin template.
// It sets Content-Type and injects the current year for the footer.
func (h *AdminHandler) render(w http.ResponseWriter, status int, page string, data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}
	data["Year"] = time.Now().Year()
	tmpl, ok := h.templates[page]
	if !ok {
		slog.Error("unknown admin template page", "page", page)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		slog.Error("admin template execution failed", "page", page, "error", err)
	}
}

// Stub handlers — Plans 03-03 and 03-04 replace these with full implementations.

func (h *AdminHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *AdminHandler) NewPost(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *AdminHandler) EditPost(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *AdminHandler) SavePost(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *AdminHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *AdminHandler) RestorePost(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *AdminHandler) PublishPost(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *AdminHandler) UnpublishPost(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *AdminHandler) Preview(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
