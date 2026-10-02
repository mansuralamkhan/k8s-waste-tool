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