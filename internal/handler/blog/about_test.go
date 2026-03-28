package blog_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestAboutPage verifies the About page handler behaviors.
func TestAboutPage(t *testing.T) {
	h := newTestHandler(&mockRepository{})

	t.Run("returns HTTP 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/about", nil)
		rec := httptest.NewRecorder()

		h.AboutPage(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("AboutPage: got status %d, want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("body contains about-title CSS class", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/about", nil)
		rec := httptest.NewRecorder()

		h.AboutPage(rec, req)

		if !strings.Contains(rec.Body.String(), "about-title") {
			t.Error("AboutPage: body missing 'about-title' CSS class")
		}
	})

	t.Run("body contains The Wild Meridian from rendered markdown", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/about", nil)
		rec := httptest.NewRecorder()

		h.AboutPage(rec, req)

		if !strings.Contains(rec.Body.String(), "The Wild Meridian") {
			t.Error("AboutPage: body missing 'The Wild Meridian' from rendered markdown")
		}
	})

	t.Run("body contains Chicago from rendered markdown", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/about", nil)
		rec := httptest.NewRecorder()

		h.AboutPage(rec, req)

		if !strings.Contains(rec.Body.String(), "Chicago") {
			t.Error("AboutPage: body missing 'Chicago' from rendered markdown")
		}
	})

	t.Run("body contains distributed compute from rendered markdown", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/about", nil)
		rec := httptest.NewRecorder()

		h.AboutPage(rec, req)

		if !strings.Contains(rec.Body.String(), "distributed compute") {
			t.Error("AboutPage: body missing 'distributed compute' from rendered markdown")
		}
	})

	t.Run("Content-Type is text/html charset=utf-8", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/about", nil)
		rec := httptest.NewRecorder()

		h.AboutPage(rec, req)

		ct := rec.Header().Get("Content-Type")
		if ct != "text/html; charset=utf-8" {
			t.Errorf("AboutPage: got Content-Type %q, want %q", ct, "text/html; charset=utf-8")
		}
	})
}
