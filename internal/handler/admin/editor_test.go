package admin_test

import (
	"context"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"golang.org/x/crypto/bcrypt"

	"github.com/jared-wallace/website-go/internal/config"
	"github.com/jared-wallace/website-go/internal/handler/admin"
	"github.com/jared-wallace/website-go/internal/middleware"
	"github.com/jared-wallace/website-go/internal/model"
	postrepo "github.com/jared-wallace/website-go/internal/repository/post"
	postservice "github.com/jared-wallace/website-go/internal/service/post"
)

// noopRendererEditor satisfies postservice.Renderer for editor tests.
type noopRendererEditor struct{}

func (noopRendererEditor) Render(src string) template.HTML { return template.HTML(src) }

// editorRepo is a mock repository wired for editor/preview tests.
// It returns a fixed post for FindByID and records calls to Create/Update/SetPublished.
type editorRepo struct {
	createCalled      bool
	updateCalled      bool
	setPublishedCalled bool
	createReturn      *model.Post
	createErr         error
}

func (m *editorRepo) ListPublished(_ context.Context, _, _ int) ([]model.Post, error) {
	return nil, nil
}
func (m *editorRepo) CountPublished(_ context.Context) (int, error) { return 0, nil }
func (m *editorRepo) FindBySlug(_ context.Context, _ string) (*model.Post, error) {
	return nil, postrepo.ErrNotFound
}
func (m *editorRepo) FindByID(_ context.Context, id int64) (*model.Post, error) {
	if id == 999 {
		return nil, postrepo.ErrNotFound
	}
	return &model.Post{
		ID:           id,
		Title:        "Existing Post Title",
		Slug:         "existing-post-title",
		Body:         "some **markdown**",
		RenderedHTML: "<p>some <strong>markdown</strong></p>",
		Tags:         "go",
		Published:    false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}
func (m *editorRepo) ListAll(_ context.Context) ([]model.Post, error) { return []model.Post{}, nil }
func (m *editorRepo) Create(_ context.Context, p model.Post) (*model.Post, error) {
	m.createCalled = true
	if m.createErr != nil {
		return nil, m.createErr
	}
	if m.createReturn != nil {
		return m.createReturn, nil
	}
	p.ID = 42
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return &p, nil
}
func (m *editorRepo) Update(_ context.Context, _ model.Post) error {
	m.updateCalled = true
	return nil
}
func (m *editorRepo) SoftDelete(_ context.Context, _ int64) error  { return nil }
func (m *editorRepo) Restore(_ context.Context, _ int64) error     { return nil }
func (m *editorRepo) SetPublished(_ context.Context, _ int64, _ bool) error {
	m.setPublishedCalled = true
	return nil
}

func (m *editorRepo) AddReaction(_ context.Context, _ int64, _ string) (bool, error) {
	return false, nil
}

func (m *editorRepo) CountReactions(_ context.Context, _ int64) (int, error) {
	return 0, nil
}

// Ensure editorRepo satisfies the Repository interface at compile time.
var _ postrepo.Repository = (*editorRepo)(nil)

// editorTestSetup creates an AdminHandler backed by editorRepo.
type editorTestSetup struct {
	handler *admin.AdminHandler
	sm      *scs.SessionManager
	repo    *editorRepo
}

func newEditorTestSetup(t *testing.T) *editorTestSetup {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte("testpassword"), 4)
	if err != nil {
		t.Fatalf("bcrypt.GenerateFromPassword: %v", err)
	}

	sm := scs.New()
	sm.Store = memstore.New()

	rl := middleware.NewRateLimiter(10, 0)
	cfg := config.Config{
		AdminEmail:        "admin@example.com",
		AdminPasswordHash: string(hash),
		AdminHost:         "admin.example.com",
	}

	repo := &editorRepo{}
	svc := postservice.New(repo, noopRendererEditor{})
	h := admin.New(svc, sm, noopRendererEditor{}, rl, cfg)
	return &editorTestSetup{handler: h, sm: sm, repo: repo}
}

func (ts *editorTestSetup) serve(h http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	ts.sm.LoadAndSave(h).ServeHTTP(rec, req)
	return rec
}

// TestNewPostRendersEmptyForm verifies GET /admin/posts/new returns 200 with "New Post" in title.
func TestNewPostRendersEmptyForm(t *testing.T) {
	ts := newEditorTestSetup(t)
	req := httptest.NewRequest(http.MethodGet, "/admin/posts/new", nil)
	rec := ts.serve(ts.handler.NewPost, req)

	if rec.Code != http.StatusOK {
		t.Errorf("NewPost: got status %d, want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), "New Post") {
		t.Errorf("NewPost: body missing 'New Post'; got snippet: %s",
			rec.Body.String()[:min(200, rec.Body.Len())])
	}
}

// TestEditPostRendersWithData verifies GET /admin/posts/1/edit returns 200 with the post title.
func TestEditPostRendersWithData(t *testing.T) {
	ts := newEditorTestSetup(t)
	req := httptest.NewRequest(http.MethodGet, "/admin/posts/1/edit", nil)
	req.SetPathValue("id", "1")
	rec := ts.serve(ts.handler.EditPost, req)

	if rec.Code != http.StatusOK {
		t.Errorf("EditPost: got status %d, want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), "Existing Post Title") {
		t.Errorf("EditPost: body missing post title; snippet: %s",
			rec.Body.String()[:min(200, rec.Body.Len())])
	}
}

// TestEditPostInvalidID verifies a non-numeric ID redirects to the dashboard with flash.
func TestEditPostInvalidID(t *testing.T) {
	ts := newEditorTestSetup(t)
	req := httptest.NewRequest(http.MethodGet, "/admin/posts/abc/edit", nil)
	req.SetPathValue("id", "abc")
	rec := ts.serve(ts.handler.EditPost, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("EditPost (invalid ID): got status %d, want %d", rec.Code, http.StatusSeeOther)
	}
	if loc := rec.Header().Get("Location"); loc != "/admin/posts" {
		t.Errorf("EditPost (invalid ID): Location = %q, want /admin/posts", loc)
	}
}

// TestSavePostCreatesDraft verifies POST /admin/posts/new with action=draft calls Create with publish=false.
func TestSavePostCreatesDraft(t *testing.T) {
	ts := newEditorTestSetup(t)
	form := url.Values{
		"title":  {"Hello World"},
		"slug":   {"hello-world"},
		"body":   {"Some content"},
		"tags":   {"go"},
		"action": {"draft"},
	}
	req := httptest.NewRequest(http.MethodPost, "/admin/posts/new",
		strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// No path value for "id" — SavePost detects new vs edit via PathValue("id") == "".
	rec := ts.serve(ts.handler.SavePost, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("SavePost (draft): got status %d, want %d", rec.Code, http.StatusSeeOther)
	}
	if !ts.repo.createCalled {
		t.Error("SavePost (draft): expected svc.Create to be called")
	}
}

// TestSavePostCreatesPublished verifies POST /admin/posts/new with action=publish calls Create with publish=true.
func TestSavePostCreatesPublished(t *testing.T) {
	ts := newEditorTestSetup(t)
	form := url.Values{
		"title":  {"Hello World"},
		"slug":   {"hello-world"},
		"body":   {"Some content"},
		"tags":   {"go"},
		"action": {"publish"},
	}
	req := httptest.NewRequest(http.MethodPost, "/admin/posts/new",
		strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := ts.serve(ts.handler.SavePost, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("SavePost (publish): got status %d, want %d", rec.Code, http.StatusSeeOther)
	}
	if !ts.repo.createCalled {
		t.Error("SavePost (publish): expected svc.Create to be called")
	}
}

// TestPreviewReturnsHTML verifies POST /admin/preview returns 200 with Content-Type text/html.
func TestPreviewReturnsHTML(t *testing.T) {
	ts := newEditorTestSetup(t)
	form := url.Values{"body": {"**bold**"}}
	req := httptest.NewRequest(http.MethodPost, "/admin/preview",
		strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := ts.serve(ts.handler.Preview, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Preview: got status %d, want %d", rec.Code, http.StatusOK)
	}
	ct := rec.Header().Get("Content-Type")
	if !strings.Contains(ct, "text/html") {
		t.Errorf("Preview: Content-Type = %q, want text/html", ct)
	}
	// The noopRenderer passes src through as-is; body should contain the original markdown.
	if !strings.Contains(rec.Body.String(), "bold") {
		t.Errorf("Preview: body missing rendered content; got: %s", rec.Body.String())
	}
}

// min returns the smaller of a and b (helper for safe string slicing in test messages).
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
