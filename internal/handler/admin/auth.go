package admin

import (
	"net"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// LoginPage renders the admin login form. If the session is already
// authenticated, the user is redirected to /admin/posts immediately.
// A flash error from a previous failed attempt is passed to the template.
func (h *AdminHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	if h.sessions.GetBool(r.Context(), "authenticated") {
		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
		return
	}
	flashErr := h.sessions.PopString(r.Context(), "flash_error")
	h.render(w, http.StatusOK, "admin-login.html", map[string]interface{}{
		"Error": flashErr,
	})
}

// LoginPost processes admin login form submissions. It enforces rate limiting,
// performs constant-time credential verification (per D-11, D-08, Pitfall 6),
// and on success renews the session token before setting the authenticated flag.
func (h *AdminHandler) LoginPost(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB limit to prevent memory exhaustion
	if err := r.ParseForm(); err != nil {
		h.render(w, http.StatusBadRequest, "admin-login.html", map[string]interface{}{
			"Error": "Invalid form submission.",
		})
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	// Rate-limit by IP before any credential work.
	if !h.rateLimiter.Allow(extractIP(r)) {
		h.render(w, http.StatusTooManyRequests, "admin-login.html", map[string]interface{}{
			"Error": "Too many attempts. Try again later.",
			"Email": email,
		})
		return
	}

	// Constant-time credential check: always call bcrypt.CompareHashAndPassword
	// regardless of whether the email matches, so the response time does not
	// reveal whether the email is correct (Pitfall 6).
	hashToCheck := h.adminHash
	if email != h.adminEmail {
		hashToCheck = h.dummyHash
	}
	err := bcrypt.CompareHashAndPassword(hashToCheck, []byte(password))
	if email != h.adminEmail || err != nil {
		h.render(w, http.StatusUnauthorized, "admin-login.html", map[string]interface{}{
			"Error": "Invalid email or password.",
			"Email": email,
		})
		return
	}

	// Renew session token before setting authenticated flag (OWASP session fixation).
	if err := h.sessions.RenewToken(r.Context()); err != nil {
		h.render(w, http.StatusInternalServerError, "admin-login.html", map[string]interface{}{
			"Error": "Session error. Please try again.",
		})
		return
	}
	h.sessions.Put(r.Context(), "authenticated", true)
	http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
}

// Logout destroys the current session and redirects the admin to the login page.
func (h *AdminHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.sessions.Destroy(r.Context()); err != nil {
		http.Error(w, "session destroy failed", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}

// extractIP returns the client IP from the request. It checks X-Real-IP first
// (set by Nginx upstream), then falls back to RemoteAddr. Port is stripped.
func extractIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		if host, _, err := net.SplitHostPort(ip); err == nil {
			return host
		}
		return ip
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
