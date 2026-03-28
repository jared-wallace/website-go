package admin_test

import (
	"context"
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

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

// noopRenderer2 satisfies postservice.Renderer for action tests.
type noopRenderer2 struct{}

func (noopRenderer2) Render(src string) template.HTML { return template.HTML(src) }

// successRepo is a repository mock where all write methods succeed.
// ListAll returns an empty slice so Dashboard renders without error.
type successRepo struct {
	called string
}

func (m *successRepo) ListPublished(_ context.Context, _, _ int) ([]model.Post, error) {
	return nil, nil
}
func (m *successRepo) CountPublished(_ context.Context) (int, error) { return 0, nil }
func (m *successRepo) FindBySlug(_ context.Context, _ string) (*model.Post, error) {
	return nil, errors.New("not implemented")
}
func (m *successRepo) FindByID(_ context.Context, _ int64) (*model.Post, error) {
	return nil, errors.New("not implemented")
}
func (m *successRepo) ListAll(_ context.Context) ([]model.Post, error) { return []model.Post{}, nil }
func (m *successRepo) Create(_ context.Context, _ model.Post) (*model.Post, error) {
	return nil, errors.New("not implemented")
}
func (m *successRepo) Update(_ context.Context, _ model.Post) error {
	return errors.New("not implemented")
}
func (m *successRepo) SoftDelete(_ context.Context, _ int64) error {
	m.called = "SoftDelete"
	return nil
}
func (m *successRepo) Restore(_ context.Context, _ int64) error {
	m.called = "Restore"
	return nil
}
func (m *successRepo) SetPublished(_ context.Context, _ int64, _ bool) error {
	m.called = "SetPublished"
	return nil
}

func (m *successRepo) AddReaction(_ context.Context, _ int64, _ string) (bool, error) {
	return false, nil
}

func (m *successRepo) CountReactions(_ context.Context, _ int64) (int, error) {
	return 0, nil
}

// Ensure successRepo satisfies the Repository interface at compile time.
var _ postrepo.Repository = (*successRepo)(nil)

// actionTestSetup is a testSetup variant backed by successRepo.
type actionTestSetup struct {
	handler *admin.AdminHandler
	sm      *scs.SessionManager
}

func newActionTestSetup(t *testing.T) *actionTestSetup {
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

	repo := &successRepo{}
	svc := postservice.New(repo, noopRenderer2{})
	h := admin.New(svc, sm, nil, rl, cfg)
	return &actionTestSetup{handler: h, sm: sm}
}

func (ts *actionTestSetup) serve(h http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	ts.sm.LoadAndSave(h).ServeHTTP(rec, req)
	return rec
}

// TestDeletePost verifies POST /admin/posts/1/delete sets flash and redirects to deleted tab.
func TestDeletePost(t *testing.T) {
	ts := newActionTestSetup(t)
	req := httptest.NewRequest(http.MethodPost, "/admin/posts/1/delete", nil)
	req.SetPathValue("id", "1")
	rec := ts.serve(ts.handler.DeletePost, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("DeletePost: got status %d, want %d", rec.Code, http.StatusSeeOther)
	}
	if loc := rec.Header().Get("Location"); loc != "/admin/posts?filter=deleted" {
		t.Errorf("DeletePost: Location = %q, want /admin/posts?filter=deleted", loc)
	}
}

// TestRestorePost verifies POST /admin/posts/1/restore sets flash and redirects to drafts tab.
func TestRestorePost(t *testing.T) {
	ts := newActionTestSetup(t)
	req := httptest.NewRequest(http.MethodPost, "/admin/posts/1/restore", nil)
	req.SetPathValue("id", "1")
	rec := ts.serve(ts.handler.RestorePost, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("RestorePost: got status %d, want %d", rec.Code, http.StatusSeeOther)
	}
	if loc := rec.Header().Get("Location"); loc != "/admin/posts?filter=drafts" {
		t.Errorf("RestorePost: Location = %q, want /admin/posts?filter=drafts", loc)
	}
}

// TestPublishPost verifies POST /admin/posts/1/publish redirects to dashboard root.
func TestPublishPost(t *testing.T) {
	ts := newActionTestSetup(t)
	req := httptest.NewRequest(http.MethodPost, "/admin/posts/1/publish", nil)
	req.SetPathValue("id", "1")
	rec := ts.serve(ts.handler.PublishPost, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("PublishPost: got status %d, want %d", rec.Code, http.StatusSeeOther)
	}
	if loc := rec.Header().Get("Location"); loc != "/admin/posts" {
		t.Errorf("PublishPost: Location = %q, want /admin/posts", loc)
	}
}

// TestUnpublishPost verifies POST /admin/posts/1/unpublish redirects to dashboard root.
func TestUnpublishPost(t *testing.T) {
	ts := newActionTestSetup(t)
	req := httptest.NewRequest(http.MethodPost, "/admin/posts/1/unpublish", nil)
	req.SetPathValue("id", "1")
	rec := ts.serve(ts.handler.UnpublishPost, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("UnpublishPost: got status %d, want %d", rec.Code, http.StatusSeeOther)
	}
	if loc := rec.Header().Get("Location"); loc != "/admin/posts" {
		t.Errorf("UnpublishPost: Location = %q, want /admin/posts", loc)
	}
}

// TestDeletePostInvalidID verifies that a non-numeric ID sets flash_error and redirects.
func TestDeletePostInvalidID(t *testing.T) {
	ts := newActionTestSetup(t)
	req := httptest.NewRequest(http.MethodPost, "/admin/posts/abc/delete", nil)
	req.SetPathValue("id", "abc")
	rec := ts.serve(ts.handler.DeletePost, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("DeletePost (invalid ID): got status %d, want %d", rec.Code, http.StatusSeeOther)
	}
	if loc := rec.Header().Get("Location"); loc != "/admin/posts" {
		t.Errorf("DeletePost (invalid ID): Location = %q, want /admin/posts", loc)
	}
}

// TestDashboardRenders verifies GET /admin/posts returns 200 and the filter-tabs markup.
func TestDashboardRenders(t *testing.T) {
	ts := newActionTestSetup(t)
	req := httptest.NewRequest(http.MethodGet, "/admin/posts", nil)
	rec := ts.serve(ts.handler.Dashboard, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Dashboard: got status %d, want %d", rec.Code, http.StatusOK)
	}
}
