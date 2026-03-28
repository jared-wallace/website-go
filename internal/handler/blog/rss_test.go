package blog_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jared-wallace/website-go/internal/model"
)

// TestServeRSS verifies that GET /rss returns 200 with the correct Content-Type
// and a valid RSS 2.0 XML body.
func TestServeRSS(t *testing.T) {
	posts := []model.Post{
		{
			ID:           1,
			Title:        "Hello World",
			Slug:         "hello-world",
			Body:         "This is a test post.",
			RenderedHTML: "<p>This is a test post.</p>",
			Tags:         "go,testing",
			CreatedAt:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			Published:    true,
		},
	}
	h := newTestHandler(&mockRepository{posts: posts, totalCount: 1})
	req := httptest.NewRequest(http.MethodGet, "/rss", nil)
	rec := httptest.NewRecorder()

	h.ServeRSS(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("ServeRSS: got status %d, want %d", rec.Code, http.StatusOK)
	}
	ct := rec.Header().Get("Content-Type")
	if ct != "application/rss+xml; charset=utf-8" {
		t.Errorf("ServeRSS: got Content-Type %q, want %q", ct, "application/rss+xml; charset=utf-8")
	}
	body := rec.Body.String()
	checks := []string{
		`<rss version="2.0">`,
		"<channel>",
		"<title>The Log</title>",
		"<description>dispatches from the deep end</description>",
		"Jared Wallace",
	}
	for _, want := range checks {
		if !strings.Contains(body, want) {
			t.Errorf("ServeRSS: body missing %q\nbody:\n%s", want, body)
		}
	}
}

// TestRSSDraftExclusion verifies the handler surfaces only what the repo returns
// (the repo already filters drafts; this confirms no phantom posts appear).
func TestRSSDraftExclusion(t *testing.T) {
	published := []model.Post{
		{
			ID: 1, Title: "Published One", Slug: "pub-1",
			RenderedHTML: "<p>one</p>", Published: true,
			CreatedAt: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			ID: 2, Title: "Published Two", Slug: "pub-2",
			RenderedHTML: "<p>two</p>", Published: true,
			CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	h := newTestHandler(&mockRepository{posts: published, totalCount: 2})
	req := httptest.NewRequest(http.MethodGet, "/rss", nil)
	rec := httptest.NewRecorder()

	h.ServeRSS(rec, req)

	body := rec.Body.String()
	if !strings.Contains(body, "Published One") {
		t.Error("TestRSSDraftExclusion: body missing 'Published One'")
	}
	if !strings.Contains(body, "Published Two") {
		t.Error("TestRSSDraftExclusion: body missing 'Published Two'")
	}
	// A draft title that was NOT loaded into the mock should not appear
	if strings.Contains(body, "Draft Post") {
		t.Error("TestRSSDraftExclusion: body unexpectedly contains 'Draft Post'")
	}
}

// TestRSSFullContent verifies that the RSS item description wraps RenderedHTML
// in a CDATA section rather than entity-escaping it.
func TestRSSFullContent(t *testing.T) {
	html := "<p>Hello <strong>World</strong></p>"
	posts := []model.Post{
		{
			ID: 1, Title: "Full Content", Slug: "full-content",
			RenderedHTML: html, Published: true,
			CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	h := newTestHandler(&mockRepository{posts: posts, totalCount: 1})
	req := httptest.NewRequest(http.MethodGet, "/rss", nil)
	rec := httptest.NewRecorder()

	h.ServeRSS(rec, req)

	body := rec.Body.String()
	if !strings.Contains(body, "<![CDATA[") {
		t.Errorf("TestRSSFullContent: body missing CDATA wrapper\nbody:\n%s", body)
	}
	// Entity-escaped HTML is the failure case from Pitfall 1
	if strings.Contains(body, "&lt;p&gt;") {
		t.Errorf("TestRSSFullContent: body contains entity-escaped HTML — CDATA not working\nbody:\n%s", body)
	}
}

// TestRSSCategories verifies that RSS items emit <category> elements for post tags.
func TestRSSCategories(t *testing.T) {
	posts := []model.Post{
		{
			ID: 1, Title: "Tagged Post", Slug: "tagged",
			RenderedHTML: "<p>tagged</p>", Tags: "go,testing,rss",
			Published: true,
			CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	h := newTestHandler(&mockRepository{posts: posts, totalCount: 1})
	req := httptest.NewRequest(http.MethodGet, "/rss", nil)
	rec := httptest.NewRecorder()

	h.ServeRSS(rec, req)

	body := rec.Body.String()
	for _, tag := range []string{"go", "testing", "rss"} {
		want := "<category>" + tag + "</category>"
		if !strings.Contains(body, want) {
			t.Errorf("TestRSSCategories: body missing %q\nbody:\n%s", want, body)
		}
	}
}
