package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/mansuralamkhan/k8s-waste-tool/backend/internal/opencost"
)

type fakeFetcher struct {
	resp      *opencost.AllocationResponse
	err       error
	gotWindow string
}

func (f *fakeFetcher) GetAllocations(window, aggregate string) (*opencost.AllocationResponse, error) {
	f.gotWindow = window
	return f.resp, f.err
}

func sampleResponse() *opencost.AllocationResponse {
	return &opencost.AllocationResponse{
		Data: []map[string]opencost.Allocation{
			{"demo": {
				Name:                  "demo",
				CPUCoreRequestAverage: 1,
				CPUCoreUsageAverage:   0,
				CPUCost:               0.01,
			}},
		},
	}
}

func get(h http.Handler, target string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
	return rec
}

func TestIdleHandlerSuccess(t *testing.T) {
	f := &fakeFetcher{resp: sampleResponse()}
	m := &Metrics{}
	rec := get(IdleHandler(f, m), "/api/v1/idle?window=15m")

	if rec.Code != http.StatusOK {
		t.Fatalf("got %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"Namespace":"demo"`) {
		t.Errorf("response missing demo namespace: %s", rec.Body.String())
	}
	if f.gotWindow != "15m" {
		t.Errorf("window = %q, want 15m", f.gotWindow)
	}
	if m.IdleRequests.Load() != 1 {
		t.Errorf("request counter = %d, want 1", m.IdleRequests.Load())
	}
}

func TestIdleHandlerDefaultsWindowToOneHour(t *testing.T) {
	f := &fakeFetcher{resp: sampleResponse()}
	get(IdleHandler(f, &Metrics{}), "/api/v1/idle")
	if f.gotWindow != "1h" {
		t.Errorf("window = %q, want 1h", f.gotWindow)
	}
}

func TestIdleHandlerUpstreamFailure(t *testing.T) {
	f := &fakeFetcher{err: errors.New("boom")}
	m := &Metrics{}
	rec := get(IdleHandler(f, m), "/api/v1/idle")

	if rec.Code != http.StatusBadGateway {
		t.Fatalf("got %d, want 502", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "OpenCost") {
		t.Errorf("error message should mention OpenCost, got: %s", rec.Body.String())
	}
	if m.UpstreamErrors.Load() != 1 {
		t.Errorf("error counter = %d, want 1", m.UpstreamErrors.Load())
	}
}

func TestIdleHandlerRejectsBadWindow(t *testing.T) {
	bad := []string{"1h&aggregate=pod", "abc", "../etc", "9999h"}
	for _, w := range bad {
		f := &fakeFetcher{resp: sampleResponse()}
		rec := get(IdleHandler(f, &Metrics{}), "/api/v1/idle?window="+url.QueryEscape(w))

		if rec.Code != http.StatusBadRequest {
			t.Errorf("window %q: got %d, want 400", w, rec.Code)
		}
		if f.gotWindow != "" {
			t.Errorf("window %q should never reach OpenCost", w)
		}
	}
}
