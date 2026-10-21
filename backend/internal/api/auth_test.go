package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func TestBasicAuth(t *testing.T) {
	h := BasicAuth("admin", "s3cret", map[string]bool{"/healthz": true}, okHandler())

	tests := []struct {
		name     string
		path     string
		user     string
		pass     string
		sendAuth bool
		want     int
	}{
		{"no credentials", "/", "", "", false, http.StatusUnauthorized},
		{"wrong password", "/", "admin", "nope", true, http.StatusUnauthorized},
		{"wrong username", "/", "root", "s3cret", true, http.StatusUnauthorized},
		{"right credentials", "/", "admin", "s3cret", true, http.StatusOK},
		{"healthz needs no auth", "/healthz", "", "", false, http.StatusOK},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			if tc.sendAuth {
				req.SetBasicAuth(tc.user, tc.pass)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != tc.want {
				t.Errorf("got %d, want %d", rec.Code, tc.want)
			}
		})
	}
}

func TestBasicAuthDisabledWhenUnset(t *testing.T) {
	h := BasicAuth("", "", nil, okHandler())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("auth should be off when creds unset, got %d", rec.Code)
	}
}
