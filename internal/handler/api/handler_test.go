package api_test

import (
	"context"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jared-wallace/website-go/internal/handler/api"
	"github.com/jared-wallace/website-go/internal/middleware"
)

// mockMetaRenderer returns controlled metadata from RenderWithMeta.
type mockMetaRenderer struct {
	meta map[string]interface{}
}

func (m *mockMetaRenderer) RenderWithMeta(src string) (template.HTML, map[string]interface{}) {
	return template.HTML("<p>rendered</p>"), m.meta
}

// mockPostService records UpsertBySlug calls.
type mockPostService struct {
	returnErr   error
	calledTitle string
	calledSlug  string
	calledBody  string
}

func (m *mockPostService) UpsertBySlug(_ context.Context, title, slug, body string) error {
	m.calledTitle = title
	m.calledSlug = slug
	m.calledBody = body
	return m.returnErr
}

const validBody = "---\nslug: test-post\ntitle: Test Post Title\n---\n# Hello World\nThis is a test post.\n"
const noSlugBody = "---\ntitle: No Slug Here\n---\n# Oops\n"
const noTitleBody = "---\nslug: titleless-post\n---\n# No Title\n"

func TestPushPost_ValidToken(t *testing.T) {
	svc := &mockPostService{}
	mr := &mockMetaRenderer{meta: map[string]interface{}{"slug": "test-post", "title": "Test Post Title"}}
	h := api.New(svc, mr)

	handler := middleware.RequireAPIToken("test-token")(http.HandlerFunc(h.PushPost))
	req := httptest.NewRequest(http.MethodPost, "/api/push", strings.NewReader(validBody))
	req.Header.Set("Authorization", "Bearer test-token")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("valid push: got status %d, want %d; body: %s", rr.Code, http.StatusOK, rr.Body.String())
	}
	if svc.calledSlug != "test-post" {
		t.Errorf("valid push: slug = %q, want %q", svc.calledSlug, "test-post")
	}
	if svc.calledTitle != "Test Post Title" {
		t.Errorf("valid push: title = %q, want %q", svc.calledTitle, "Test Post Title")
	}
}

func TestPushPost_NoToken(t *testing.T) {
	svc := &mockPostService{}
	mr := &mockMetaRenderer{meta: map[string]interface{}{}}
	h := api.New(svc, mr)

	handler := middleware.RequireAPIToken("test-token")(http.HandlerFunc(h.PushPost))
	req := httptest.NewRequest(http.MethodPost, "/api/push", strings.NewReader(validBody))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("no token: got status %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}

func TestPushPost_InvalidToken(t *testing.T) {
	svc := &mockPostService{}
	mr := &mockMetaRenderer{meta: map[string]interface{}{}}
	h := api.New(svc, mr)

	handler := middleware.RequireAPIToken("test-token")(http.HandlerFunc(h.PushPost))
	req := httptest.NewRequest(http.MethodPost, "/api/push", strings.NewReader(validBody))
	req.Header.Set("Authorization", "Bearer wrong-token")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("invalid token: got status %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}

func TestPushPost_NoSlug(t *testing.T) {
	svc := &mockPostService{}
	mr := &mockMetaRenderer{meta: map[string]interface{}{"title": "No Slug Here"}}
	h := api.New(svc, mr)

	handler := middleware.RequireAPIToken("test-token")(http.HandlerFunc(h.PushPost))
	req := httptest.NewRequest(http.MethodPost, "/api/push", strings.NewReader(noSlugBody))
	req.Header.Set("Authorization", "Bearer test-token")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("no slug: got status %d, want %d; body: %s", rr.Code, http.StatusBadRequest, rr.Body.String())
	}
}

func TestPushPost_NoTitle(t *testing.T) {
	svc := &mockPostService{}
	mr := &mockMetaRenderer{meta: map[string]interface{}{"slug": "titleless-post"}}
	h := api.New(svc, mr)

	handler := middleware.RequireAPIToken("test-token")(http.HandlerFunc(h.PushPost))
	req := httptest.NewRequest(http.MethodPost, "/api/push", strings.NewReader(noTitleBody))
	req.Header.Set("Authorization", "Bearer test-token")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("no title: got status %d, want %d; body: %s", rr.Code, http.StatusOK, rr.Body.String())
	}
	if svc.calledTitle != "titleless-post" {
		t.Errorf("no title: should fallback to slug, got %q", svc.calledTitle)
	}
}

func TestPushPost_BodyTooLarge(t *testing.T) {
	svc := &mockPostService{}
	mr := &mockMetaRenderer{meta: map[string]interface{}{"slug": "big"}}
	h := api.New(svc, mr)

	handler := middleware.RequireAPIToken("test-token")(http.HandlerFunc(h.PushPost))
	bigBody := strings.Repeat("x", 1<<20+1) // 1 MB + 1 byte
	req := httptest.NewRequest(http.MethodPost, "/api/push", strings.NewReader(bigBody))
	req.Header.Set("Authorization", "Bearer test-token")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Should fail: either 400 or 500 is acceptable when body exceeds limit
	if rr.Code == http.StatusOK {
		t.Errorf("body too large: got status 200, want non-200")
	}
}
