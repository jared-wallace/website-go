package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexedwards/scs/v2"
	"github.com/jared-wallace/website-go/internal/middleware"
)

func TestRequireSession_Unauthenticated(t *testing.T) {
	sm := scs.New()

	handler := sm.LoadAndSave(middleware.RequireSession(sm)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))

	req := httptest.NewRequest(http.MethodGet, "/admin/posts", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("unauthenticated: got status %d, want %d (SeeOther)", rec.Code, http.StatusSeeOther)
	}
	loc := rec.Header().Get("Location")
	if loc != "/admin/login" {
		t.Errorf("unauthenticated: Location = %q, want /admin/login", loc)
	}
}

func TestRequireSession_Authenticated(t *testing.T) {
	sm := scs.New()

	// Wrap the auth middleware with LoadAndSave so the session context is populated.
	// We need to set the session value in a pre-handler, then check in the next handler.
	setAuth := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sm.Put(r.Context(), "authenticated", true)
	})

	var nextCalled bool
	nextHandler := middleware.RequireSession(sm)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	// Chain: LoadAndSave -> setAuth then requireSession in a single request
	// We simulate an authenticated session by making two requests:
	// 1. First request sets authenticated=true in session
	// 2. Second request uses saved session cookie to authenticate

	// Use in-memory store (default scs store)
	mux := http.NewServeMux()
	mux.Handle("/set", sm.LoadAndSave(setAuth))
	mux.Handle("/protected", sm.LoadAndSave(nextHandler))

	// Step 1: set authenticated flag
	req1 := httptest.NewRequest(http.MethodGet, "/set", nil)
	rec1 := httptest.NewRecorder()
	mux.ServeHTTP(rec1, req1)

	// Extract session cookie
	cookies := rec1.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("no session cookie set after authentication")
	}

	// Step 2: access protected route with session cookie
	req2 := httptest.NewRequest(http.MethodGet, "/protected", nil)
	for _, c := range cookies {
		req2.AddCookie(c)
	}
	rec2 := httptest.NewRecorder()
	mux.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Errorf("authenticated: got status %d, want %d", rec2.Code, http.StatusOK)
	}
	if !nextCalled {
		t.Error("authenticated: next handler was not called")
	}
}
