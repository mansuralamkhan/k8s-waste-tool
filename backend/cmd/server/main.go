package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/mansuralamkhan/k8s-waste-tool/backend/internal/api"
	"github.com/mansuralamkhan/k8s-waste-tool/backend/internal/opencost"
)

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	openCostURL := envOr("OPENCOST_URL", "http://localhost:9003")
	client := opencost.NewClient(openCostURL)
	metrics := &api.Metrics{}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", api.Healthz)
	mux.HandleFunc("/metrics", metrics.Handler())
	mux.HandleFunc("/api/v1/idle", api.IdleHandler(client, metrics))
	mux.Handle("/", http.FileServer(http.Dir(envOr("DASHBOARD_DIR", "/dashboard"))))

	user := os.Getenv("AUTH_USERNAME")
	pass := os.Getenv("AUTH_PASSWORD")
	if user == "" || pass == "" {
		slog.Warn("AUTH_USERNAME/AUTH_PASSWORD not set; dashboard and API are open to anyone who can reach this port")
	}

	// Logger wraps auth, so rejected logins show up in the logs too.
	handler := api.RequestLogger(api.BasicAuth(user, pass, map[string]bool{"/healthz": true}, mux))
	srv := &http.Server{Addr: ":8080", Handler: handler, ReadHeaderTimeout: 5 * time.Second}

	slog.Info("starting server", "addr", srv.Addr, "opencost_url", openCostURL)
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("server stopped", "error", err.Error())
		os.Exit(1)
	}
}