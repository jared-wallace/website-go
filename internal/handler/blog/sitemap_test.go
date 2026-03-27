package blog_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jared-wallace/website-go/internal/model"
)

// TestServeSitemap verifies that GET /sitemap.xml returns 200 with valid
// Sitemap 0.9 XML containing the homepage URL and published post URLs.
func TestServeSitemap(t *testing.T) {
	posts := []model.Post{
		{
			ID: 1, Title: "Sitemap Post", Slug: "sitemap-post",
			Published: true,
			CreatedAt: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
		},
	}
	h := newTestHandler(&mockRepository{posts: posts, totalCount: 1})
	req := httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil)
	rec := httptest.NewRecorder()

	h.ServeSitemap(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("ServeSitemap: got status %d, want %d", rec.Code, http.StatusOK)
	}
	ct := rec.Header().Get("Content-Type")
	if ct != "application/xml; charset=utf-8" {
		t.Errorf("ServeSitemap: got Content-Type %q, want %q", ct, "application/xml; charset=utf-8")
	}
	body := rec.Body.String()
	checks := []string{
		`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`,
		"https://jared-wallace.com/",
		"https://jared-wallace.com/posts/sitemap-post",
	}
	for _, want := range checks {
		if !strings.Contains(body, want) {
			t.Errorf("ServeSitemap: body missing %q\nbody:\n%s", want, body)
		}
	}
}

// TestServeRobots verifies that GET /robots.txt returns 200 plain text
// with the required Sitemap directive.
func TestServeRobots(t *testing.T) {
	h := newTestHandler(&mockRepository{})
	req := httptest.NewRequest(http.MethodGet, "/robots.txt", nil)
	rec := httptest.NewRecorder()

	h.ServeRobots(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("ServeRobots: got status %d, want %d", rec.Code, http.StatusOK)
	}
	ct := rec.Header().Get("Content-Type")
	if ct != "text/plain; charset=utf-8" {
		t.Errorf("ServeRobots: got Content-Type %q, want %q", ct, "text/plain; charset=utf-8")
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Sitemap: https://jared-wallace.com/sitemap.xml") {
		t.Errorf("ServeRobots: body missing Sitemap directive\nbody:\n%s", body)
	}
}
