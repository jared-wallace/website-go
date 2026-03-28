package blog_test

import (
	"context"
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jared-wallace/website-go/internal/handler/blog"
	"github.com/jared-wallace/website-go/internal/model"
	postservice "github.com/jared-wallace/website-go/internal/service/post"
)

// mockRepository is a test double for the post repository interface.
type mockRepository struct {
	findErr        error
	posts          []model.Post
	totalCount     int
	reactionCount  int
	alreadyReacted bool
}

func (m *mockRepository) ListPublished(_ context.Context, limit, offset int) ([]model.Post, error) {
	if offset >= len(m.posts) {
		return nil, nil
	}
	end := offset + limit
	if end > len(m.posts) {
		end = len(m.posts)
	}
	return m.posts[offset:end], nil
}

func (m *mockRepository) CountPublished(_ context.Context) (int, error) {
	return m.totalCount, nil
}

func (m *mockRepository) FindBySlug(_ context.Context, _ string) (*model.Post, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	if len(m.posts) == 0 {
		return nil, errors.New("not found")
	}
	return &m.posts[0], nil
}

func (m *mockRepository) FindByID(_ context.Context, _ int64) (*model.Post, error) {
	return nil, errors.New("not implemented")
}

func (m *mockRepository) ListAll(_ context.Context) ([]model.Post, error) {
	return nil, nil
}

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

func (m *mockRepository) AddReaction(_ context.Context, _ int64, _ string) (bool, error) {
	return m.alreadyReacted, nil
}

func (m *mockRepository) CountReactions(_ context.Context, _ int64) (int, error) {
	return m.reactionCount, nil
}

// noopRenderer satisfies the postservice.Renderer interface for handler tests
// that never invoke admin write paths.
type noopRenderer struct{}

func (noopRenderer) Render(src string) template.HTML {
	return template.HTML(src)
}

// newTestHandler constructs a BlogHandler with the given repository.
func newTestHandler(repo *mockRepository) *blog.BlogHandler {
	svc := postservice.New(repo, noopRenderer{})
	return blog.New(svc)
}

// TestNotFound verifies that the NotFound handler returns HTTP 404.
func TestNotFound(t *testing.T) {
	h := newTestHandler(&mockRepository{})
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	rec := httptest.NewRecorder()

	h.NotFound(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("NotFound: got status %d, want %d", rec.Code, http.StatusNotFound)
	}
}

// TestListPostsEmpty verifies that ListPosts returns HTTP 200 when there are no posts.
func TestListPostsEmpty(t *testing.T) {
	h := newTestHandler(&mockRepository{posts: nil, totalCount: 0})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	h.ListPosts(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("ListPosts (empty): got status %d, want %d", rec.Code, http.StatusOK)
	}
}

// TestPostOGMeta verifies that ShowPost renders OG and Twitter Card meta tags for a post.
func TestPostOGMeta(t *testing.T) {
	repo := &mockRepository{
		posts: []model.Post{
			{Title: "Test Post Title", Slug: "test-post", Body: "Some test body content for excerpt generation"},
		},
		totalCount: 1,
	}
	h := newTestHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/posts/test-post", nil)
	req.SetPathValue("slug", "test-post")
	rec := httptest.NewRecorder()

	h.ShowPost(rec, req)

	body := rec.Body.String()
	checks := []string{
		`og:title" content="Test Post Title"`,
		`og:type" content="article"`,
		`og:url" content="https://jared-wallace.com/posts/test-post"`,
		`twitter:card" content="summary"`,
		`og:image" content="https://jared-wallace.com/static/og-fallback.png"`,
	}
	for _, want := range checks {
		if !strings.Contains(body, want) {
			t.Errorf("TestPostOGMeta: response missing %q", want)
		}
	}
}

// TestListOGMeta verifies that ListPosts renders site-level OG and Twitter Card meta tags.
func TestListOGMeta(t *testing.T) {
	repo := &mockRepository{
		posts: []model.Post{
			{Title: "Some Post", Slug: "some-post"},
		},
		totalCount: 1,
	}
	h := newTestHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	h.ListPosts(rec, req)

	body := rec.Body.String()
	checks := []string{
		`og:title" content="The Wild Meridian"`,
		`og:type" content="website"`,
		`og:description" content="dispatches from the deep end"`,
		`twitter:card" content="summary"`,
	}
	for _, want := range checks {
		if !strings.Contains(body, want) {
			t.Errorf("TestListOGMeta: response missing %q", want)
		}
	}
}

// ---------------------------------------------------------------------------
// Phase 11 Wave 0 — template validation tests (RED: expected to fail until
// Plan 02 updates the templates)
// ---------------------------------------------------------------------------

// TestNavAboutLinkRemoved verifies NAV-01: About link is no longer in the
// top nav but still reachable from the footer.
func TestNavAboutLinkRemoved(t *testing.T) {
	repo := &mockRepository{
		posts:      []model.Post{{Title: "Test Post", Slug: "test-post"}},
		totalCount: 1,
	}
	h := newTestHandler(repo)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	h.ListPosts(rec, req)
	body := rec.Body.String()

	// The About nav-link next to the dark-toggle should be gone.
	if strings.Contains(body, `<a href="/about" class="nav-link">About</a>`) {
		t.Error("TestNavAboutLinkRemoved: About link still present in top nav")
	}
	// But an About link should still exist somewhere (footer).
	if !strings.Contains(body, `href="/about"`) {
		t.Error("TestNavAboutLinkRemoved: About link missing entirely — expected in footer")
	}
}

// TestFooterTwoSection verifies NAV-02: footer has inner and copyright
// sections plus a footer navigation landmark.
func TestFooterTwoSection(t *testing.T) {
	repo := &mockRepository{
		posts:      []model.Post{{Title: "Test Post", Slug: "test-post"}},
		totalCount: 1,
	}
	h := newTestHandler(repo)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	h.ListPosts(rec, req)
	body := rec.Body.String()

	for _, want := range []string{
		`class="footer-inner"`,
		`class="footer-copyright"`,
		`aria-label="Footer navigation"`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("TestFooterTwoSection: missing %q", want)
		}
	}
}

// TestFooterPersonalityPhrase verifies NAV-03: footer contains the
// personality phrase element.
func TestFooterPersonalityPhrase(t *testing.T) {
	repo := &mockRepository{
		posts:      []model.Post{{Title: "Test Post", Slug: "test-post"}},
		totalCount: 1,
	}
	h := newTestHandler(repo)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	h.ListPosts(rec, req)
	body := rec.Body.String()

	if !strings.Contains(body, `class="footer-phrase"`) {
		t.Error("TestFooterPersonalityPhrase: missing footer-phrase element")
	}
	if !strings.Contains(body, `Still anchored. Still writing.`) {
		t.Error("TestFooterPersonalityPhrase: missing personality phrase text")
	}
}

// TestNavAriaLabels verifies NAV-04: both nav and footer nav carry proper
// ARIA labels.
func TestNavAriaLabels(t *testing.T) {
	repo := &mockRepository{
		posts:      []model.Post{{Title: "Test Post", Slug: "test-post"}},
		totalCount: 1,
	}
	h := newTestHandler(repo)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	h.ListPosts(rec, req)
	body := rec.Body.String()

	if !strings.Contains(body, `aria-label="Main navigation"`) {
		t.Error("TestNavAriaLabels: missing Main navigation aria-label")
	}
	if !strings.Contains(body, `aria-label="Footer navigation"`) {
		t.Error("TestNavAriaLabels: missing Footer navigation aria-label")
	}
}

// TestRopeDividerSVG verifies ATMO-03: rope divider is an SVG element, not
// an <hr>.
func TestRopeDividerSVG(t *testing.T) {
	repo := &mockRepository{
		posts:      []model.Post{{Title: "Test Post", Slug: "test-post"}},
		totalCount: 1,
	}
	h := newTestHandler(repo)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	h.ListPosts(rec, req)
	body := rec.Body.String()

	if !strings.Contains(body, `<svg class="rope-divider"`) {
		t.Error("TestRopeDividerSVG: missing SVG rope-divider element")
	}
	if strings.Contains(body, `<hr class="rope-divider">`) {
		t.Error("TestRopeDividerSVG: old HR rope-divider still present")
	}
}

// TestListHero verifies TYPO-03: homepage renders a hero area with the site
// title and tagline, even when there are no posts.
func TestListHero(t *testing.T) {
	h := newTestHandler(&mockRepository{posts: nil, totalCount: 0})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	h.ListPosts(rec, req)
	body := rec.Body.String()

	for _, want := range []string{
		`class="list-hero"`,
		`class="list-hero-title"`,
		`The Wild Meridian</h1>`,
		`dispatches from the deep end</p>`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("TestListHero: missing %q", want)
		}
	}
}

// TestShowPostNotFound verifies that ShowPost returns HTTP 404 for a missing slug.
func TestShowPostNotFound(t *testing.T) {
	repo := &mockRepository{
		findErr: postservice.ErrNotFound,
	}
	h := newTestHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/posts/no-such-post", nil)
	// Go 1.22 ServeMux sets PathValue; in tests we set it manually.
	req.SetPathValue("slug", "no-such-post")
	rec := httptest.NewRecorder()

	h.ShowPost(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("ShowPost (not found): got status %d, want %d", rec.Code, http.StatusNotFound)
	}
}
