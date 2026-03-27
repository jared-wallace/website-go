package admin_test

import (
	"context"
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"golang.org/x/crypto/bcrypt"

	"github.com/jared-wallace/website-go/internal/config"
	"github.com/jared-wallace/website-go/internal/handler/admin"
	"github.com/jared-wallace/website-go/internal/middleware"
	"github.com/jared-wallace/website-go/internal/model"
	postservice "github.com/jared-wallace/website-go/internal/service/post"
)

// testPassword is the known plaintext used in tests.
const testPassword = "testpassword"

// noopRenderer satisfies postservice.Renderer without importing the markdown package.
type noopRenderer struct{}

func (noopRenderer) Render(src string) template.HTML { return template.HTML(src) }

// mockRepository satisfies the repository.Repository interface with no-ops.
type mockRepository struct{}

func (m *mockRepository) ListPublished(_ context.Context, _, _ int) ([]model.Post, error) {
	return nil, nil
}
func (m *mockRepository) CountPublished(_ context.Context) (int, error) { return 0, nil }
func (m *mockRepository) FindBySlug(_ context.Context, _ string) (*model.Post, error) {
	return nil, errors.New("not implemented")
}
func (m *mockRepository) FindByID(_ context.Context, _ int64) (*model.Post, error) {
	return nil, errors.New("not implemented")
}
func (m *mockRepository) ListAll(_ context.Context) ([]model.Post, error) { return nil, nil }
func (m *mockRepository) Create(_ context.Context, _ model.Post) (*model.Post, error) {
	return nil, errors.New("not implemented")
}
func (m *mockRepository) Update(_ context.Context, _ model.Post) error {
	return errors.New("not implemented")
}
func (m *mockRepository) SoftDelete(_ context.Context, _ int64) error {
	return errors.New("not implemented")
}
func (m *mockRepository) Restore(_ context.Context, _ int64) error {
	return errors.New("not implemented")
}
func (m *mockRepository) SetPublished(_ context.Context, _ int64, _ bool) error {
	return errors.New("not implemented")
}

// testSetup holds everything needed to run a single admin handler test.
type testSetup struct {
	handler *admin.AdminHandler
	sm      *scs.SessionManager
}

// newTestSetup creates an AdminHandler and its session manager from the same
// SCS instance so h.sessions.GetBool uses the same context data the middleware injects.
func newTestSetup(t *testing.T) *testSetup {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte(testPassword), 4) // cost 4 for test speed
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

	svc := postservice.New(&mockRepository{}, noopRenderer{})
	h := admin.New(svc, sm, nil, rl, cfg)
	return &testSetup{handler: h, sm: sm}
}

// serve wraps the handler in sm.LoadAndSave so session data is in the request context.
func (ts *testSetup) serve(h http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	ts.sm.LoadAndSave(h).ServeHTTP(rec, req)
	return rec
}

// TestLoginPageRendersForm verifies GET /admin/login returns 200 with "Sign In".
func TestLoginPageRendersForm(t *testing.T) {
	ts := newTestSetup(t)
	req := httptest.NewRequest(http.MethodGet, "/admin/login", nil)
	rec := ts.serve(ts.handler.LoginPage, req)

	if rec.Code != http.StatusOK {
		t.Errorf("LoginPage: got status %d, want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), "Sign In") {
		t.Error("LoginPage: response body does not contain 'Sign In'")
	}
}

// TestLoginPostInvalidCredentials verifies bad credentials return 401 with generic error.
func TestLoginPostInvalidCredentials(t *testing.T) {
	ts := newTestSetup(t)
	form := url.Values{"email": {"wrong@example.com"}, "password": {"wrongpass"}}
	req := httptest.NewRequest(http.MethodPost, "/admin/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := ts.serve(ts.handler.LoginPost, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("LoginPost (invalid): got status %d, want %d", rec.Code, http.StatusUnauthorized)
	}
	if !strings.Contains(rec.Body.String(), "Invalid email or password") {
		t.Errorf("LoginPost (invalid): body missing error; got: %s", rec.Body.String())
	}
}

// TestLoginPostValidCredentials verifies correct credentials redirect to /admin/posts.
func TestLoginPostValidCredentials(t *testing.T) {
	ts := newTestSetup(t)
	form := url.Values{"email": {"admin@example.com"}, "password": {testPassword}}
	req := httptest.NewRequest(http.MethodPost, "/admin/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := ts.serve(ts.handler.LoginPost, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("LoginPost (valid): got status %d, want %d", rec.Code, http.StatusSeeOther)
	}
	if loc := rec.Header().Get("Location"); loc != "/admin/posts" {
		t.Errorf("LoginPost (valid): Location = %q, want /admin/posts", loc)
	}
}

// TestLogout verifies POST /admin/logout destroys session and redirects to login.
func TestLogout(t *testing.T) {
	ts := newTestSetup(t)
	req := httptest.NewRequest(http.MethodPost, "/admin/logout", nil)
	rec := ts.serve(ts.handler.Logout, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("Logout: got status %d, want %d", rec.Code, http.StatusSeeOther)
	}
	if loc := rec.Header().Get("Location"); loc != "/admin/login" {
		t.Errorf("Logout: Location = %q, want /admin/login", loc)
	}
}
