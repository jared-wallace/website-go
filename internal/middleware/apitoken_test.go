package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
}

func TestRequireAPIToken_ValidToken(t *testing.T) {
	h := RequireAPIToken("test-secret")(okHandler())
	req := httptest.NewRequest(http.MethodPost, "/api/push", nil)
	req.Header.Set("Authorization", "Bearer test-secret")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("valid token: got status %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestRequireAPIToken_MissingHeader(t *testing.T) {
	h := RequireAPIToken("test-secret")(okHandler())
	req := httptest.NewRequest(http.MethodPost, "/api/push", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("missing header: got status %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}

func TestRequireAPIToken_InvalidToken(t *testing.T) {
	h := RequireAPIToken("test-secret")(okHandler())
	req := httptest.NewRequest(http.MethodPost, "/api/push", nil)
	req.Header.Set("Authorization", "Bearer wrong-token")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("invalid token: got status %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}

func TestRequireAPIToken_EmptyBearer(t *testing.T) {
	h := RequireAPIToken("test-secret")(okHandler())
	req := httptest.NewRequest(http.MethodPost, "/api/push", nil)
	req.Header.Set("Authorization", "Bearer ")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("empty bearer: got status %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}

func TestRequireAPIToken_EmptyConfigToken(t *testing.T) {
	h := RequireAPIToken("")(okHandler())
	req := httptest.NewRequest(http.MethodPost, "/api/push", nil)
	req.Header.Set("Authorization", "Bearer some-token")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("empty config token: got status %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}
