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
	"github.com/jared-wallace/website-go/internal/middleware"
	postservice "github.com/jared-wallace/website-go/internal/service/post"
	"github.com/jared-wallace/website-go/web"
)

// Renderer is the interface used by Preview to convert markdown to HTML.
// Using an interface (rather than *markdown.Renderer directly) keeps
// the handler mock-friendly in unit tests without importing the markdown package.
type Renderer interface {
	Render(src string) template.HTML
}

// AdminHandler holds template sets, session manager, and auth state.
// All fields are goroutine-safe after construction.
type AdminHandler struct {
	svc         *postservice.Service
	sessions    *scs.SessionManager
	renderer    Renderer
	rateLimiter *middleware.RateLimiter
	funcMap     template.FuncMap
	templates   map[string]*template.Template
	adminEmail  string
	adminHash   []byte
	dummyHash   []byte // pre-computed hash for timing-safe comparison when email doesn't match
	imageDir    string // filesystem path for uploaded images (EBS-backed)
}

// New parses admin templates and returns a ready-to-use AdminHandler.
// It pre-computes a dummy bcrypt hash at startup to ensure constant-time
// credential checking even when the email doesn't match (Pitfall 6).
// Panics if any template fails to parse — a programmer error, not a runtime one.
func New(svc *postservice.Service, sm *scs.SessionManager, r Renderer, rl *middleware.RateLimiter, cfg config.Config) *AdminHandler {
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
		imageDir:    cfg.ImageDir,
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

// SetImageDir overrides the image directory — used in tests with t.TempDir().
func (h *AdminHandler) SetImageDir(dir string) { h.imageDir = dir }

// NewPost, EditPost, SavePost, and Preview are implemented in editor.go and preview.go.
// UploadImage is implemented in upload.go.
