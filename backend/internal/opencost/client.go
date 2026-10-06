package opencost

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

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

var ErrOpenCostUnreachable = errors.New("couldn't reach OpenCost — check it's installed and OPENCOST_URL is correct")

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

type AllocationResponse struct {
	Data []map[string]Allocation `json:"data"`
}

func (c *Client) GetAllocations(window, aggregate string) (*AllocationResponse, error) {
	url := fmt.Sprintf("%s/allocation/compute?window=%s&aggregate=%s", c.BaseURL, window, aggregate)
	resp, err := c.http.Get(url)
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