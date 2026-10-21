package opencost

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetAllocationsSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/allocation/compute" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		w.Write([]byte(`{"data":[{"demo":{"name":"demo","cpuCoreRequestAverage":1}}]}`))
	}))
	defer srv.Close()

	got, err := NewClient(srv.URL).GetAllocations("15m", "namespace")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Data[0]["demo"].CPUCoreRequestAverage != 1 {
		t.Errorf("parsed wrong value: %+v", got.Data[0]["demo"])
	}
}

func TestGetAllocationsNon200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	_, err := NewClient(srv.URL).GetAllocations("15m", "namespace")
	if !errors.Is(err, ErrOpenCostUnreachable) {
		t.Errorf("want ErrOpenCostUnreachable, got %v", err)
	}
}

func TestGetAllocationsMalformedJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	}))
	defer srv.Close()

	_, err := NewClient(srv.URL).GetAllocations("15m", "namespace")
	if err == nil {
		t.Error("expected an error for malformed JSON")
	}
}

func TestGetAllocationsConnectionRefused(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	addr := srv.URL
	srv.Close() // nothing is listening any more

	_, err := NewClient(addr).GetAllocations("15m", "namespace")
	if !errors.Is(err, ErrOpenCostUnreachable) {
		t.Errorf("want ErrOpenCostUnreachable, got %v", err)
	}
}

func TestGetAllocationsTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(300 * time.Millisecond)
	}))
	defer srv.Close()

	c := NewClient(srv.URL)
	c.http.Timeout = 50 * time.Millisecond
	_, err := c.GetAllocations("15m", "namespace")
	if !errors.Is(err, ErrOpenCostUnreachable) {
		t.Errorf("want ErrOpenCostUnreachable, got %v", err)
	}
}
