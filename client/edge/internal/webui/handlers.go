package webui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

// Handler handles Web UI API requests
type Handler struct {
	configPath  string
	statusFunc  func() interface{}
	apiURL      string
	logger      zerolog.Logger
	logBuffer   *LogBuffer
	reloadChan  chan<- bool
	mu          sync.RWMutex
}

// LogBuffer stores recent log entries
type LogBuffer struct {
	mu      sync.RWMutex
	entries []LogEntry
	maxSize int
}

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
}

// NewLogBuffer creates a new log buffer
func NewLogBuffer(maxSize int) *LogBuffer {
	return &LogBuffer{
		entries: make([]LogEntry, 0, maxSize),
		maxSize: maxSize,
	}
}

// Add adds a log entry to the buffer
func (lb *LogBuffer) Add(entry LogEntry) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	lb.entries = append(lb.entries, entry)
	if len(lb.entries) > lb.maxSize {
		lb.entries = lb.entries[1:]
	}
}

// GetAll returns all log entries
func (lb *LogBuffer) GetAll() []LogEntry {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	// Return a copy
	entries := make([]LogEntry, len(lb.entries))
	copy(entries, lb.entries)
	return entries
}

// NewHandler creates a new handler
func NewHandler(configPath string, statusFunc func() interface{}, apiURL string, logger zerolog.Logger) *Handler {
	return &Handler{
		configPath: configPath,
		statusFunc: statusFunc,
		apiURL:     apiURL,
		logger:     logger,
		logBuffer:  NewLogBuffer(1000),
	}
}

// SetReloadChan sets the reload channel for triggering config reloads
func (h *Handler) SetReloadChan(reloadChan chan<- bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.reloadChan = reloadChan
}

// GetConfig returns the current configuration
func (h *Handler) GetConfig(c *fiber.Ctx) error {
	data, err := os.ReadFile(h.configPath)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to read config file")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to read configuration",
		})
	}

	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		h.logger.Error().Err(err).Msg("Failed to parse config file")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to parse configuration",
		})
	}

	return c.JSON(config)
}

// SaveConfig saves the configuration
func (h *Handler) SaveConfig(c *fiber.Ctx) error {
	var config map[string]interface{}
	if err := c.BodyParser(&config); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// Convert to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to marshal config")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to marshal configuration",
		})
	}

	// Create backup
	backupPath := h.configPath + ".backup." + time.Now().Format("20060102-150405")
	if err := copyFile(h.configPath, backupPath); err != nil {
		h.logger.Warn().Err(err).Msg("Failed to create config backup")
	}

	// Write new config
	if err := os.WriteFile(h.configPath, data, 0644); err != nil {
		h.logger.Error().Err(err).Msg("Failed to write config file")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to save configuration",
		})
	}

	h.logger.Info().Msg("Configuration saved via Web UI")

	// Trigger reload (non-blocking)
	h.mu.RLock()
	reloadChan := h.reloadChan
	h.mu.RUnlock()

	if reloadChan != nil {
		select {
		case reloadChan <- true:
			h.logger.Info().Msg("Config reload triggered")
		default:
			h.logger.Debug().Msg("Reload already pending")
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Configuration saved and reload triggered",
	})
}

// GetStatus returns the current status
func (h *Handler) GetStatus(c *fiber.Ctx) error {
	if h.statusFunc != nil {
		status := h.statusFunc()
		return c.JSON(status)
	}

	return c.JSON(fiber.Map{
		"status": "unknown",
	})
}

// GetSystemInfo returns system information
func (h *Handler) GetSystemInfo(c *fiber.Ctx) error {
	info := GetSystemInfo()
	return c.JSON(info)
}

// ValidateAPIServer validates API server connectivity
func (h *Handler) ValidateAPIServer(c *fiber.Ctx) error {
	var req struct {
		URL string `json:"url"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// Try to connect to API server
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(req.URL + "/health")
	if err != nil {
		return c.JSON(fiber.Map{
			"valid":   false,
			"message": fmt.Sprintf("Failed to connect: %v", err),
		})
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return c.JSON(fiber.Map{
			"valid":   true,
			"message": "API server is reachable",
		})
	}

	return c.JSON(fiber.Map{
		"valid":   false,
		"message": fmt.Sprintf("API server returned status %d", resp.StatusCode),
	})
}

// Login proxies login request to API server
func (h *Handler) Login(c *fiber.Ctx) error {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&credentials); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// Forward to API server
	body, _ := json.Marshal(credentials)
	resp, err := http.Post(
		h.apiURL+"/api/auth/login",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": fmt.Sprintf("Failed to connect to API server: %v", err),
		})
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Forward response
	return c.Status(resp.StatusCode).Send(respBody)
}

// GetDevices proxies device list request to API server
func (h *Handler) GetDevices(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Missing authorization token",
		})
	}

	req, _ := http.NewRequest(http.MethodGet, h.apiURL+"/api/devices", nil)
	req.Header.Set("Authorization", token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": fmt.Sprintf("Failed to connect to API server: %v", err),
		})
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	return c.Status(resp.StatusCode).Send(respBody)
}

// CreateDevice proxies device creation to API server
func (h *Handler) CreateDevice(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Missing authorization token",
		})
	}

	req, _ := http.NewRequest(http.MethodPost, h.apiURL+"/api/devices", bytes.NewBuffer(c.Body()))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": fmt.Sprintf("Failed to connect to API server: %v", err),
		})
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	return c.Status(resp.StatusCode).Send(respBody)
}

// DetectInputSources detects available input sources
func (h *Handler) DetectInputSources(c *fiber.Ctx) error {
	sources := DetectInputSources()
	return c.JSON(sources)
}

// GetRecentLogs returns recent log entries
func (h *Handler) GetRecentLogs(c *fiber.Ctx) error {
	entries := h.logBuffer.GetAll()
	return c.JSON(entries)
}

// copyFile copies a file
func copyFile(src, dst string) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}
