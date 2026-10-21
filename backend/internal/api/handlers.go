package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/mansuralamkhan/k8s-waste-tool/backend/internal/opencost"
	"github.com/mansuralamkhan/k8s-waste-tool/backend/internal/waste"
)

// AllocationFetcher is what the handler needs from OpenCost.
// An interface, so tests can swap in a fake.
type AllocationFetcher interface {
	GetAllocations(window, aggregate string) (*opencost.AllocationResponse, error)
}

var windowRe = regexp.MustCompile(`^\d{1,3}[mhd]$`)

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func IdleHandler(f AllocationFetcher, m *Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.IdleRequests.Add(1)

		window := r.URL.Query().Get("window")
		if window == "" {
			window = "1h"
		}
		if !windowRe.MatchString(window) {
			http.Error(w, "invalid window: use something like 15m, 1h, 24h or 7d", http.StatusBadRequest)
			return
		}

		resp, err := f.GetAllocations(window, "namespace")
		if err != nil {
			m.UpstreamErrors.Add(1)
			slog.Error("fetching allocations failed", "window", window, "error", err.Error())
			http.Error(w, "Couldn't fetch cluster data right now — check that OpenCost is running and reachable.", http.StatusBadGateway)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(waste.Calculate(resp))
	}
}
