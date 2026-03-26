package blog_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jared-wallace/website-go/internal/handler/blog"
	"github.com/jared-wallace/website-go/internal/model"
	postservice "github.com/jared-wallace/website-go/internal/service/post"
)

// mockRepository is a test double for the post repository interface.
type mockRepository struct {
	posts      []model.Post
	totalCount int
	findErr    error
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

// newTestHandler constructs a BlogHandler with the given repository.
func newTestHandler(repo *mockRepository) *blog.BlogHandler {
	svc := postservice.New(repo)
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
