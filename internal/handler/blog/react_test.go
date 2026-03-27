package blog_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jared-wallace/website-go/internal/model"
	postservice "github.com/jared-wallace/website-go/internal/service/post"
)

// TestReact verifies that a first-time POST /posts/{slug}/react returns 200
// with JSON containing count and already_reacted:false.
func TestReact(t *testing.T) {
	repo := &mockRepository{
		posts:         []model.Post{{ID: 1, Slug: "test-slug", Title: "Test Post", Body: "body"}},
		reactionCount: 1,
		alreadyReacted: false,
	}
	h := newTestHandler(repo)

	req := httptest.NewRequest(http.MethodPost, "/posts/test-slug/react", nil)
	req.SetPathValue("slug", "test-slug")
	rec := httptest.NewRecorder()

	h.React(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("React: got status %d, want %d", rec.Code, http.StatusOK)
	}
	body := rec.Body.String()
	if !strings.Contains(body, `"count":1`) {
		t.Errorf("React: response missing count:1, got: %s", body)
	}
	if !strings.Contains(body, `"already_reacted":false`) {
		t.Errorf("React: response missing already_reacted:false, got: %s", body)
	}
}

// TestReactDuplicate verifies that a duplicate reaction returns already_reacted:true.
func TestReactDuplicate(t *testing.T) {
	repo := &mockRepository{
		posts:          []model.Post{{ID: 1, Slug: "test-slug", Title: "Test Post", Body: "body"}},
		reactionCount:  5,
		alreadyReacted: true,
	}
	h := newTestHandler(repo)

	req := httptest.NewRequest(http.MethodPost, "/posts/test-slug/react", nil)
	req.SetPathValue("slug", "test-slug")
	rec := httptest.NewRecorder()

	h.React(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("TestReactDuplicate: got status %d, want %d", rec.Code, http.StatusOK)
	}
	body := rec.Body.String()
	if !strings.Contains(body, `"already_reacted":true`) {
		t.Errorf("TestReactDuplicate: response missing already_reacted:true, got: %s", body)
	}
	if !strings.Contains(body, `"count":5`) {
		t.Errorf("TestReactDuplicate: response missing count:5, got: %s", body)
	}
}

// TestReactNotFound verifies POST /posts/{slug}/react returns 404 when slug is missing.
func TestReactNotFound(t *testing.T) {
	repo := &mockRepository{
		findErr: postservice.ErrNotFound,
	}
	h := newTestHandler(repo)

	req := httptest.NewRequest(http.MethodPost, "/posts/no-such/react", nil)
	req.SetPathValue("slug", "no-such")
	rec := httptest.NewRecorder()

	h.React(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("TestReactNotFound: got status %d, want %d", rec.Code, http.StatusNotFound)
	}
}

// TestPostReactionCount verifies GET /posts/{slug} includes the reaction count in the HTML.
func TestPostReactionCount(t *testing.T) {
	repo := &mockRepository{
		posts:         []model.Post{{ID: 1, Slug: "test-slug", Title: "Test Post", Body: "body content"}},
		reactionCount: 42,
	}
	h := newTestHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/posts/test-slug", nil)
	req.SetPathValue("slug", "test-slug")
	rec := httptest.NewRecorder()

	h.ShowPost(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("TestPostReactionCount: got status %d, want %d", rec.Code, http.StatusOK)
	}
	body := rec.Body.String()
	if !strings.Contains(body, `id="reaction-count"`) {
		t.Errorf("TestPostReactionCount: response missing id=\"reaction-count\", body length=%d", len(body))
	}
	if !strings.Contains(body, `>42<`) {
		t.Errorf("TestPostReactionCount: response missing reaction count 42")
	}
}
