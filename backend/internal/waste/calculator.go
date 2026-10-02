package waste
import "github.com/mansuralamkhan/k8s-waste-tool/backend/internal/opencost"

// NamespaceIdle is the idle-capacity result for one namespace.
type NamespaceIdle struct {
	Namespace   string
	IdleCPU     float64 // idle cores
	IdleCPUCost float64 // dollars
	IdleRAM     float64 // idle bytes
	IdleRAMCost float64 // dollars
	IdleCost    float64 // total dollars, this window
}

// Calculate derives idle capacity per namespace from an OpenCost allocation
// response. Idle = requested - used, clamped at zero, with cost prorated by
// the idle share of what was requested.
func Calculate(resp *opencost.AllocationResponse) []NamespaceIdle {
	var results []NamespaceIdle

	for _, nsMap := range resp.Data {
		for name, a := range nsMap {
			idleCPU := a.CPUCoreRequestAverage - a.CPUCoreUsageAverage
			if idleCPU < 0 {
				idleCPU = 0
			}
			idleRAM := a.RAMByteRequestAverage - a.RAMByteUsageAverage
			if idleRAM < 0 {
				idleRAM = 0
			}

			idleCPUCost := proratedCost(idleCPU, a.CPUCoreRequestAverage, a.CPUCost)
			idleRAMCost := proratedCost(idleRAM, a.RAMByteRequestAverage, a.RAMCost)

			results = append(results, NamespaceIdle{
				Namespace:   name,
				IdleCPU:     idleCPU,
				IdleCPUCost: idleCPUCost,
				IdleRAM:     idleRAM,
				IdleRAMCost: idleRAMCost,
				IdleCost:    idleCPUCost + idleRAMCost,
			})
		}
	}
	return results
}

// proratedCost returns the share of totalCost attributable to the idle
// portion of what was requested. If nothing was requested, idle cost is 0
// (can't be "wasting" money on something you never reserved).
func proratedCost(idle, requested, totalCost float64) float64 {
	if requested <= 0 {
		return 0
	}
	return (idle / requested) * totalCost
}