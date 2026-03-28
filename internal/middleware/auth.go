// Package middleware provides HTTP middleware for the admin panel.
package middleware

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
)

// RequireSession returns a middleware that checks for an authenticated session.
// Requests without a valid "authenticated" session boolean are redirected to
// /admin/login with a 303 See Other status. Authenticated requests are passed
// through to the next handler unchanged.
func RequireSession(sm *scs.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !sm.GetBool(r.Context(), "authenticated") {
				http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
