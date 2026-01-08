package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents an API client for device authentication
type Client struct {
	apiURL     string
	deviceID   string
	secret     string
	httpClient *http.Client
}

// DeviceInfo represents device information from API
type DeviceInfo struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	PublishPath string    `json:"publish_path"`
	Disabled    bool      `json:"disabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewClient creates a new API client
func NewClient(apiURL, deviceID, secret string, timeout time.Duration) *Client {
	return &Client{
		apiURL:   apiURL,
		deviceID: deviceID,
		secret:   secret,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// GetDeviceInfo fetches device information from API
func (c *Client) GetDeviceInfo() (*DeviceInfo, error) {
	url := fmt.Sprintf("%s/api/devices/%s", c.apiURL, c.deviceID)
	
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add device secret as authorization header
	req.Header.Set("X-Device-Secret", c.secret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned error status %d: %s", resp.StatusCode, string(body))
	}

	var deviceInfo DeviceInfo
	if err := json.Unmarshal(body, &deviceInfo); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &deviceInfo, nil
}

// HealthCheck performs a health check against the API
func (c *Client) HealthCheck() error {
	url := fmt.Sprintf("%s/health", c.apiURL)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check returned status %d", resp.StatusCode)
	}

	return nil
}

// ReportHeartbeat sends a heartbeat to the API (for future use)
func (c *Client) ReportHeartbeat(status string) error {
	url := fmt.Sprintf("%s/api/devices/%s/heartbeat", c.apiURL, c.deviceID)
	
	payload := map[string]string{
		"status": status,
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Device-Secret", c.secret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("heartbeat returned status %d: %s", resp.StatusCode, string(body))
	}

	// 404 is OK - API doesn't have heartbeat endpoint yet
	return nil
}
