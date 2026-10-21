package waste

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/mansuralamkhan/k8s-waste-tool/backend/internal/opencost"
)

func TestCalculate_DemoNamespaceIsIdle(t *testing.T) {
	raw, err := os.ReadFile("testdata/sample-allocation.json")
	if err != nil {
		t.Fatalf("reading fixture: %v", err)
	}

	var resp opencost.AllocationResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		t.Fatalf("unmarshaling fixture: %v", err)
	}

	results := Calculate(&resp)

	var demo *NamespaceIdle
	for i := range results {
		if results[i].Namespace == "demo" {
			demo = &results[i]
		}
	}
	if demo == nil {
		t.Fatal("expected a demo namespace result")
	}
	if demo.IdleCPU <= 0 {
		t.Errorf("expected demo namespace to have idle CPU, got %f", demo.IdleCPU)
	}
	if demo.IdleRAM <= 0 {
		t.Errorf("expected demo namespace to have idle RAM, got %f", demo.IdleRAM)
	}
}

func TestProratedCost(t *testing.T) {
	cases := []struct {
		name      string
		idle      float64
		requested float64
		total     float64
		want      float64
	}{
		{"negative cost from upstream", 0.05, 0.05, -0.00008, 0},
		{"nothing requested", 1, 0, 0.01, 0},
		{"half idle", 0.5, 1, 0.10, 0.05},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := proratedCost(tc.idle, tc.requested, tc.total)
			if got-tc.want > 1e-9 || tc.want-got > 1e-9 {
				t.Errorf("got %f, want %f", got, tc.want)
			}
		})
	}
}
