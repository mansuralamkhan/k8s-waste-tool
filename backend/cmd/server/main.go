package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mansuralamkhan/k8s-waste-tool/backend/internal/api"
	"github.com/mansuralamkhan/k8s-waste-tool/backend/internal/opencost"
	"github.com/mansuralamkhan/k8s-waste-tool/backend/internal/waste"
)

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	client := opencost.NewClient(envOr("OPENCOST_URL", "http://localhost:9003"))

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("/api/v1/idle", func(w http.ResponseWriter, r *http.Request) {
		window := r.URL.Query().Get("window")
		if window == "" {
			window = "1h"
		}

		resp, err := client.GetAllocations(window, "namespace")
		if err != nil {
			log.Printf("idle handler error: %v", err)
			http.Error(w, "Couldn't fetch cluster data right now — check that OpenCost is running and reachable.", http.StatusBadGateway)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(waste.Calculate(resp))
	})

	// DASHBOARD_DIR defaults to the container path; override for local `go run`.
	mux.Handle("/", http.FileServer(http.Dir(envOr("DASHBOARD_DIR", "/dashboard"))))

	user := os.Getenv("AUTH_USERNAME")
	pass := os.Getenv("AUTH_PASSWORD")
	if user == "" || pass == "" {
		log.Println("WARNING: AUTH_USERNAME/AUTH_PASSWORD not set; dashboard and API are open to anyone who can reach this port")
	}

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           api.BasicAuth(user, pass, map[string]bool{"/healthz": true}, mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("starting server on :8080")
	log.Fatal(srv.ListenAndServe())
}