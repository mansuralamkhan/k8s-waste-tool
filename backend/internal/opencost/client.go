// Package opencost wraps calls to an in-cluster OpenCost instance's
// /allocation API to pull cost and usage data per namespace/pod.
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

// AllocationResponse is a minimal placeholder shape — expand this to match
// OpenCost's actual /allocation API response once wired up against a real
// cluster.
type AllocationResponse struct {
	Data []map[string]any `json:"data"`
}

func (c *Client) GetAllocations(window string) (*AllocationResponse, error) {
	url := fmt.Sprintf("%s/allocation?window=%s", c.BaseURL, window)
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
