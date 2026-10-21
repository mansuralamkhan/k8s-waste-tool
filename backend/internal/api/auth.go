package api

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
)

// BasicAuth protects every route except those in skip (e.g. /healthz,
// so Kubernetes probes keep working). Empty username or password = auth off.
func BasicAuth(username, password string, skip map[string]bool, next http.Handler) http.Handler {
	if username == "" || password == "" {
		return next
	}
	wantUser := sha256.Sum256([]byte(username))
	wantPass := sha256.Sum256([]byte(password))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if skip[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}
		u, p, ok := r.BasicAuth()
		gotUser := sha256.Sum256([]byte(u))
		gotPass := sha256.Sum256([]byte(p))
		userOK := subtle.ConstantTimeCompare(gotUser[:], wantUser[:]) == 1
		passOK := subtle.ConstantTimeCompare(gotPass[:], wantPass[:]) == 1

		if !ok || !userOK || !passOK {
			w.Header().Set("WWW-Authenticate", `Basic realm="K8s Waste Detector"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
