package webui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

	// Create backup if config file exists
	if _, err := os.Stat(h.configPath); err == nil {
		backupPath := h.configPath + ".backup." + time.Now().Format("20060102-150405")
		if err := copyFile(h.configPath, backupPath); err != nil {
			h.logger.Warn().Err(err).Msg("Failed to create config backup")
		} else {
			h.logger.Info().Str("backup", backupPath).Msg("Config backup created")
		}
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
		h.logger.Error().Err(err).Msg("Failed to parse validate-server request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	h.logger.Info().Str("url", req.URL).Msg("Validating API server")

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
		// Save the validated API URL for later use (e.g., in Login)
		h.mu.Lock()
		h.apiURL = req.URL
		h.mu.Unlock()
		
		h.logger.Info().Str("api_url", req.URL).Msg("API server validated and saved")
		
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
		h.logger.Error().Err(err).Msg("Failed to parse login credentials")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// Get API URL
	h.mu.RLock()
	apiURL := h.apiURL
	h.mu.RUnlock()

	h.logger.Debug().
		Str("api_url", apiURL).
		Str("username", credentials.Username).
		Msg("Login request received")

	if apiURL == "" {
		h.logger.Warn().Msg("Login attempted without API URL configured")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "API server URL not configured. Please validate API server first (Step 1).",
		})
	}

	// Forward to API server
	body, _ := json.Marshal(credentials)
	resp, err := http.Post(
		apiURL+"/api/login",
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

	// Get API URL
	h.mu.RLock()
	apiURL := h.apiURL
	h.mu.RUnlock()

	if apiURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "API server URL not configured",
		})
	}

	req, _ := http.NewRequest(http.MethodGet, apiURL+"/api/devices", nil)
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

	// Get API URL
	h.mu.RLock()
	apiURL := h.apiURL
	h.mu.RUnlock()

	if apiURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "API server URL not configured",
		})
	}

	req, _ := http.NewRequest(http.MethodPost, apiURL+"/api/devices", bytes.NewBuffer(c.Body()))
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

// ExportConfig exports the current configuration as a downloadable file
func (h *Handler) ExportConfig(c *fiber.Ctx) error {
	// Read current config
	data, err := os.ReadFile(h.configPath)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to read config file for export")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to read configuration",
		})
	}

	// Set headers for file download
	filename := fmt.Sprintf("pawstream-config-%s.yaml", time.Now().Format("20060102-150405"))
	c.Set("Content-Type", "application/x-yaml")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	
	h.logger.Info().Str("filename", filename).Msg("Config exported")
	
	return c.Send(data)
}

// ImportConfig imports a configuration file
func (h *Handler) ImportConfig(c *fiber.Ctx) error {
	// Get uploaded file
	file, err := c.FormFile("config")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "No file uploaded",
		})
	}

	// Validate file extension
	if !strings.HasSuffix(file.Filename, ".yaml") && !strings.HasSuffix(file.Filename, ".yml") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid file type. Please upload a YAML file",
		})
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to open uploaded file")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to read uploaded file",
		})
	}
	defer src.Close()

	// Read file content
	data, err := io.ReadAll(src)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to read file content")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to read file content",
		})
	}

	// Validate YAML syntax
	var testConfig map[string]interface{}
	if err := yaml.Unmarshal(data, &testConfig); err != nil {
		h.logger.Error().Err(err).Msg("Invalid YAML in uploaded file")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid YAML format",
		})
	}

	// Create backup of current config
	backupPath := h.configPath + ".backup." + time.Now().Format("20060102-150405")
	if err := copyFile(h.configPath, backupPath); err != nil {
		h.logger.Warn().Err(err).Msg("Failed to create config backup")
	} else {
		h.logger.Info().Str("path", backupPath).Msg("Created config backup")
	}

	// Write new config
	if err := os.WriteFile(h.configPath, data, 0644); err != nil {
		h.logger.Error().Err(err).Msg("Failed to write imported config")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to save configuration",
		})
	}

	h.logger.Info().Str("filename", file.Filename).Msg("Config imported successfully")

	// Trigger reload (non-blocking)
	h.mu.RLock()
	reloadChan := h.reloadChan
	h.mu.RUnlock()

	if reloadChan != nil {
		select {
		case reloadChan <- true:
			h.logger.Info().Msg("Config reload triggered after import")
		default:
			h.logger.Debug().Msg("Reload already pending")
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Configuration imported and reload triggered",
		"backup":  backupPath,
	})
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

// GetEnginesAvailable returns available streaming engines
func (h *Handler) GetEnginesAvailable(c *fiber.Ctx) error {
	// Import stream package functions (will add these)
	engines := []map[string]interface{}{
		{
			"name":        "ffmpeg",
			"displayName": "FFmpeg",
			"description": "兼容性好，通用场景推荐",
			"available":   true, // FFmpeg is always available
			"features": []string{
				"硬件加速支持",
				"多种编码器",
				"成熟稳定",
			},
		},
	}

	// Check if GStreamer is available
	gstAvailable := isGStreamerAvailable()
	engines = append(engines, map[string]interface{}{
		"name":        "gstreamer",
		"displayName": "GStreamer",
		"description": "低延迟，专业场景推荐",
		"available":   gstAvailable,
		"features": []string{
			"超低延迟（100-200ms）",
			"丰富的硬件编码器支持",
			"灵活的 pipeline 架构",
		},
		"installCommand": "sudo apt-get install gstreamer1.0-tools gstreamer1.0-plugins-base gstreamer1.0-plugins-good gstreamer1.0-rtsp",
	})

	return c.JSON(map[string]interface{}{
		"engines": engines,
	})
}

// GetPresets returns available configuration presets
func (h *Handler) GetPresets(c *fiber.Ctx) error {
	presets := []map[string]interface{}{
		{
			"id":          "low-latency",
			"name":        "低延迟",
			"description": "优化低延迟（100-200ms）",
			"engine":      "gstreamer",
			"latency":     "100-200ms",
			"quality":     "中等",
			"resource":    "中等",
			"scenario":    "监控、实时互动",
			"icon":        "⚡",
		},
		{
			"id":          "high-quality",
			"name":        "高质量",
			"description": "优化视频质量",
			"engine":      "ffmpeg",
			"latency":     "1-2s",
			"quality":     "最高",
			"resource":    "高",
			"scenario":    "录制、存档",
			"icon":        "🎬",
		},
		{
			"id":          "balanced",
			"name":        "平衡",
			"description": "平衡质量和延迟",
			"engine":      "ffmpeg",
			"latency":     "500ms",
			"quality":     "良好",
			"resource":    "中等",
			"scenario":    "通用场景",
			"icon":        "⚖️",
		},
		{
			"id":          "power-save",
			"name":        "省电",
			"description": "优化资源占用",
			"engine":      "ffmpeg",
			"latency":     "500ms",
			"quality":     "中等",
			"resource":    "低",
			"scenario":    "边缘设备、省电模式",
			"icon":        "🔋",
		},
	}

	return c.JSON(map[string]interface{}{
		"presets": presets,
	})
}

// GetEncoders returns available video encoders
func (h *Handler) GetEncoders(c *fiber.Ctx) error {
	encoders := map[string]interface{}{
		"hardware": []map[string]interface{}{},
		"software": []map[string]interface{}{
			{
				"name":        "libx264",
				"displayName": "x264 (软件)",
				"available":   true,
				"type":        "software",
			},
		},
	}

	// Check hardware encoders
	hwEncoders := []map[string]interface{}{}

	// VAAPI (Intel)
	if isEncoderAvailable("vaapih264enc") || isEncoderAvailable("h264_vaapi") {
		hwEncoders = append(hwEncoders, map[string]interface{}{
			"name":        "vaapi",
			"displayName": "VAAPI (Intel)",
			"available":   true,
			"type":        "hardware",
			"vendor":      "Intel",
		})
	}

	// NVENC (NVIDIA)
	if isEncoderAvailable("nvh264enc") || isEncoderAvailable("h264_nvenc") {
		hwEncoders = append(hwEncoders, map[string]interface{}{
			"name":        "nvenc",
			"displayName": "NVENC (NVIDIA)",
			"available":   true,
			"type":        "hardware",
			"vendor":      "NVIDIA",
		})
	}

	// QSV (Intel Quick Sync)
	if isEncoderAvailable("h264_qsv") {
		hwEncoders = append(hwEncoders, map[string]interface{}{
			"name":        "qsv",
			"displayName": "Quick Sync (Intel)",
			"available":   true,
			"type":        "hardware",
			"vendor":      "Intel",
		})
	}

	// VideoToolbox (macOS)
	if isEncoderAvailable("vtenc_h264") || isEncoderAvailable("h264_videotoolbox") {
		hwEncoders = append(hwEncoders, map[string]interface{}{
			"name":        "videotoolbox",
			"displayName": "VideoToolbox (Apple)",
			"available":   true,
			"type":        "hardware",
			"vendor":      "Apple",
		})
	}

	encoders["hardware"] = hwEncoders

	return c.JSON(encoders)
}

// isGStreamerAvailable checks if GStreamer is installed
func isGStreamerAvailable() bool {
	// Try to run gst-launch-1.0 --version
	cmd := execCommand("gst-launch-1.0", "--version")
	_, err := cmd.CombinedOutput()
	return err == nil
}

// isEncoderAvailable checks if a specific encoder is available
func isEncoderAvailable(encoder string) bool {
	// Check GStreamer element
	if strings.Contains(encoder, "gst") || !strings.Contains(encoder, "_") {
		cmd := execCommand("gst-inspect-1.0", encoder)
		output, err := cmd.CombinedOutput()
		if err == nil && !strings.Contains(string(output), "No such element") {
			return true
		}
	}

	// Check FFmpeg encoder
	cmd := execCommand("ffmpeg", "-hide_banner", "-encoders")
	output, err := cmd.CombinedOutput()
	if err == nil && strings.Contains(string(output), encoder) {
		return true
	}

	return false
}

// execCommand is a helper to execute commands
func execCommand(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}
