package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/mansuralamkhan/k8s-waste-tool/backend/internal/opencost"
	"github.com/mansuralamkhan/k8s-waste-tool/backend/internal/waste"
)

func main() {
	openCostURL := os.Getenv("OPENCOST_URL")
	if openCostURL == "" {
		openCostURL = "http://localhost:9003"
	}
	client := opencost.NewClient(openCostURL)

	mux := http.NewServeMux()
	// mux.Handle("/", http.FileServer(http.Dir("../dashboard")))
	mux.Handle("/", http.FileServer(http.Dir("/dashboard")))

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

    results := waste.Calculate(resp)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(results)
})

	log.Println("starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}