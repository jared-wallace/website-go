package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"
)

// RequireAPIToken returns a middleware that validates a bearer token from the
// Authorization header using constant-time comparison to prevent timing attacks.
// If the configured token is empty, the endpoint is effectively disabled and
// all requests receive 401.
func RequireAPIToken(token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if token == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			raw := r.Header.Get("Authorization")
			got := strings.TrimPrefix(raw, "Bearer ")

			// TrimPrefix returns the original string when the prefix is absent,
			// so raw == got means "Bearer " was not present at all.
			if raw == got || got == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			if subtle.ConstantTimeCompare([]byte(got), []byte(token)) != 1 {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
