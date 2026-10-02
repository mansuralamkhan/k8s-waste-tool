// internal/opencost/client.go
package opencost

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	BaseURL string
	http    *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL, http: &http.Client{}}
}

// Allocation mirrors the fields of OpenCost's /allocation/compute response
// that we actually need for idle-capacity calculation. Extra fields in the
// real response are ignored by encoding/json automatically.
type Allocation struct {
	Name                  string  `json:"name"`
	Minutes               float64 `json:"minutes"`
	CPUCoreRequestAverage float64 `json:"cpuCoreRequestAverage"`
	CPUCoreUsageAverage   float64 `json:"cpuCoreUsageAverage"`
	CPUCost               float64 `json:"cpuCost"`
	RAMByteRequestAverage float64 `json:"ramByteRequestAverage"`
	RAMByteUsageAverage   float64 `json:"ramByteUsageAverage"`
	RAMCost               float64 `json:"ramCost"`
	TotalCost             float64 `json:"totalCost"`
}

// AllocationResponse's "data" field is a list containing one map of
// namespace name -> Allocation, per the real API shape.
type AllocationResponse struct {
	Data []map[string]Allocation `json:"data"`
}

func (c *Client) GetAllocations(window, aggregate string) (*AllocationResponse, error) {
	url := fmt.Sprintf("%s/allocation/compute?window=%s&aggregate=%s", c.BaseURL, window, aggregate)
	resp, err := c.http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching allocations: %w", err)
	}
	defer resp.Body.Close()

	var out AllocationResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decoding allocations: %w", err)
	}
	return &out, nil
}