package api

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

// Metrics holds a few counters, exposed in Prometheus text format.
type Metrics struct {
	IdleRequests   atomic.Int64
	UpstreamErrors atomic.Int64
}

func (m *Metrics) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		fmt.Fprintf(w, "# HELP wastedetector_idle_requests_total Requests to /api/v1/idle.\n")
		fmt.Fprintf(w, "# TYPE wastedetector_idle_requests_total counter\n")
		fmt.Fprintf(w, "wastedetector_idle_requests_total %d\n", m.IdleRequests.Load())
		fmt.Fprintf(w, "# HELP wastedetector_opencost_errors_total Failed calls to OpenCost.\n")
		fmt.Fprintf(w, "# TYPE wastedetector_opencost_errors_total counter\n")
		fmt.Fprintf(w, "wastedetector_opencost_errors_total %d\n", m.UpstreamErrors.Load())
	}
}
