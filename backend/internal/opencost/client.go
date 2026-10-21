// Package opencost wraps calls to an in-cluster OpenCost instance's
// /allocation/compute API to pull cost and usage data per namespace.
package opencost

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// ErrOpenCostUnreachable is returned (wrapped) whenever we can't get a
// usable answer from OpenCost: connection failed, timed out, or a non-200.
var ErrOpenCostUnreachable = errors.New("couldn't reach OpenCost — check it's installed and OPENCOST_URL is correct")

type Client struct {
	BaseURL string
	http    *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		http: &http.Client{
			Timeout: 8 * time.Second,
		},
	}
}

// Allocation mirrors only the fields of OpenCost's response that we need
// for idle-capacity calculation. Extra fields in the real response are
// ignored automatically by encoding/json.
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

// AllocationResponse's "data" field is a list of maps, each mapping a
// namespace name to its Allocation, per the real API shape.
type AllocationResponse struct {
	Data []map[string]Allocation `json:"data"`
}

// GetAllocations fetches per-aggregate allocations for the given window
// (e.g. "15m", "1h", "7d") and aggregation (e.g. "namespace").
func (c *Client) GetAllocations(window, aggregate string) (*AllocationResponse, error) {
	endpoint := fmt.Sprintf("%s/allocation/compute?window=%s&aggregate=%s",
		c.BaseURL, url.QueryEscape(window), url.QueryEscape(aggregate))

	resp, err := c.http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrOpenCostUnreachable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: OpenCost returned status %d", ErrOpenCostUnreachable, resp.StatusCode)
	}

	var out AllocationResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("OpenCost returned something unexpected: %w", err)
	}
	return &out, nil
}
