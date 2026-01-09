package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Device              DeviceConfig  `yaml:"device"`
	API                 APIConfig     `yaml:"api"`
	Input               InputConfig   `yaml:"input"`
	Video               VideoConfig   `yaml:"video"`
	Stream              StreamConfig  `yaml:"stream"`
	Log                 LogConfig     `yaml:"log"`
	Health              HealthConfig  `yaml:"health"`
	WebUI               *WebUIConfig  `yaml:"webui"`
	ShutdownTimeout     time.Duration `yaml:"shutdown_timeout"`
	ValidateInputOnStart bool         `yaml:"validate_input_on_start"`
}

// DeviceConfig holds device identification
type DeviceConfig struct {
	ID     string `yaml:"id"`
	Secret string `yaml:"secret"`
}

// APIConfig holds API server configuration
type APIConfig struct {
	URL     string        `yaml:"url"`
	Timeout time.Duration `yaml:"timeout"`
}

// InputConfig holds video input configuration
type InputConfig struct {
	Type   string `yaml:"type"`   // v4l2, rtsp, file, test
	Source string `yaml:"source"` // device path, URL, or file path
}

// VideoConfig holds video encoding parameters
type VideoConfig struct {
	Codec     string `yaml:"codec"`
	Width     int    `yaml:"width"`
	Height    int    `yaml:"height"`
	Framerate int    `yaml:"framerate"`
	Bitrate   int    `yaml:"bitrate"`
}

// StreamConfig holds streaming configuration
type StreamConfig struct {
	Engine               string            `yaml:"engine"` // ffmpeg, gstreamer
	Preset               string            `yaml:"preset"` // low-latency, high-quality, balanced, power-save
	URL                  string            `yaml:"url"`
	ReconnectInterval    time.Duration     `yaml:"reconnect_interval"`
	MaxReconnectAttempts int               `yaml:"max_reconnect_attempts"`
	FFmpeg               FFmpegConfig      `yaml:"ffmpeg"`
	GStreamer            GStreamerConfig   `yaml:"gstreamer"`
}

// FFmpegConfig holds FFmpeg-specific configuration
type FFmpegConfig struct {
	Preset    string   `yaml:"preset"`     // ultrafast, superfast, veryfast, faster, fast, medium, slow, slower, veryslow
	Tune      string   `yaml:"tune"`       // film, animation, grain, stillimage, fastdecode, zerolatency
	HWAccel   string   `yaml:"hwaccel"`    // none, auto, vaapi, nvenc, qsv, videotoolbox
	ExtraArgs []string `yaml:"extra_args"` // Custom FFmpeg arguments
}

// GStreamerConfig holds GStreamer-specific configuration
type GStreamerConfig struct {
	LatencyMs   int  `yaml:"latency_ms"`   // Pipeline latency in milliseconds
	UseHardware bool `yaml:"use_hardware"` // Use hardware encoders
	BufferSize  int  `yaml:"buffer_size"`  // Buffer size in microseconds
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string `yaml:"level"`  // debug, info, warn, error
	File   string `yaml:"file"`   // empty for stdout
	Format string `yaml:"format"` // json, console
}

// HealthConfig holds health check configuration
type HealthConfig struct {
	Enabled bool   `yaml:"enabled"`
	Address string `yaml:"address"` // e.g., ":9090"
}

// WebUIConfig holds Web UI configuration
type WebUIConfig struct {
	Enabled bool             `yaml:"enabled"`
	Host    string           `yaml:"host"`
	Port    int              `yaml:"port"`
	Auth    *WebUIAuthConfig `yaml:"auth"`
}

// WebUIAuthConfig holds Web UI authentication configuration
type WebUIAuthConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// Load reads configuration from a YAML file and environment variables
func Load(filePath string) (*Config, error) {
	cfg := &Config{
		// Default values
		Device: DeviceConfig{},
		API: APIConfig{
			Timeout: 10 * time.Second,
		},
		Video: VideoConfig{
			Codec:     "h264",
			Width:     1280,
			Height:    720,
			Framerate: 30,
			Bitrate:   2000000,
		},
		Stream: StreamConfig{
			Engine:               "ffmpeg", // Default to FFmpeg for compatibility
			Preset:               "",       // No preset by default
			ReconnectInterval:    5 * time.Second,
			MaxReconnectAttempts: 0, // infinite
			FFmpeg: FFmpegConfig{
				Preset:    "ultrafast",
				Tune:      "zerolatency",
				HWAccel:   "auto",
				ExtraArgs: []string{},
			},
			GStreamer: GStreamerConfig{
				LatencyMs:   100,
				UseHardware: true,
				BufferSize:  1000,
			},
		},
		Log: LogConfig{
			Level:  "info",
			Format: "json",
		},
		Health: HealthConfig{
			Enabled: false,
			Address: ":9090",
		},
		WebUI: &WebUIConfig{
			Enabled: true,
			Host:    "0.0.0.0", // Listen on all interfaces
			Port:    8088,
			Auth:    nil, // No auth by default
		},
		ShutdownTimeout:      10 * time.Second,
		ValidateInputOnStart: true,
	}

	// Load from file if provided
	if filePath != "" {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Override with environment variables
	cfg.applyEnvOverrides()

	// Apply preset if specified
	if cfg.Stream.Preset != "" {
		if err := cfg.ApplyPreset(cfg.Stream.Preset); err != nil {
			return nil, fmt.Errorf("failed to apply preset: %w", err)
		}
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// Save writes the configuration to a file
func Save(path string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// applyEnvOverrides applies environment variable overrides
func (c *Config) applyEnvOverrides() {
	if v := os.Getenv("PAWSTREAM_DEVICE_ID"); v != "" {
		c.Device.ID = v
	}
	if v := os.Getenv("PAWSTREAM_DEVICE_SECRET"); v != "" {
		c.Device.Secret = v
	}
	if v := os.Getenv("PAWSTREAM_API_URL"); v != "" {
		c.API.URL = v
	}
	if v := os.Getenv("PAWSTREAM_INPUT_TYPE"); v != "" {
		c.Input.Type = v
	}
	if v := os.Getenv("PAWSTREAM_INPUT_SOURCE"); v != "" {
		c.Input.Source = v
	}
	if v := os.Getenv("PAWSTREAM_LOG_LEVEL"); v != "" {
		c.Log.Level = v
	}
	if v := os.Getenv("PAWSTREAM_LOG_FILE"); v != "" {
		c.Log.File = v
	}
	if v := os.Getenv("PAWSTREAM_STREAM_ENGINE"); v != "" {
		c.Stream.Engine = v
	}
	if v := os.Getenv("PAWSTREAM_STREAM_PRESET"); v != "" {
		c.Stream.Preset = v
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Device validation
	if c.Device.ID == "" {
		return fmt.Errorf("device.id is required")
	}
	if c.Device.Secret == "" {
		return fmt.Errorf("device.secret is required")
	}

	// API validation
	if c.API.URL == "" {
		return fmt.Errorf("api.url is required")
	}

	// Input validation
	validInputTypes := map[string]bool{
		"v4l2": true,
		"rtsp": true,
		"file": true,
		"test": true,
		"auto": true, // Allow auto for zero-config startup
	}
	if !validInputTypes[c.Input.Type] {
		return fmt.Errorf("input.type must be one of: v4l2, rtsp, file, test, auto")
	}
	// Source validation - not required for 'test' and 'auto' types
	if c.Input.Type != "test" && c.Input.Type != "auto" && c.Input.Source == "" {
		return fmt.Errorf("input.source is required for input type: %s", c.Input.Type)
	}

	// Video validation
	if c.Video.Width <= 0 || c.Video.Height <= 0 {
		return fmt.Errorf("video.width and video.height must be positive")
	}
	if c.Video.Framerate <= 0 {
		return fmt.Errorf("video.framerate must be positive")
	}
	if c.Video.Bitrate <= 0 {
		return fmt.Errorf("video.bitrate must be positive")
	}

	// Log validation
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[c.Log.Level] {
		return fmt.Errorf("log.level must be one of: debug, info, warn, error")
	}

	return nil
}

// FindConfigFile searches for configuration file in standard locations
// Returns the path to the first found configuration file, or empty string if none found
func FindConfigFile() string {
	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		home = ""
	}

	// Search paths in priority order
	searchPaths := []string{
		"./config.yaml",
		"./configs/config.yaml",
	}

	// Add home directory path if available
	if home != "" {
		searchPaths = append(searchPaths, filepath.Join(home, ".pawstream", "config.yaml"))
	}

	// Search for config file
	for _, path := range searchPaths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			continue
		}

		if _, err := os.Stat(absPath); err == nil {
			return absPath
		}
	}

	return ""
}

// IsComplete checks if the configuration has all required fields
func (c *Config) IsComplete() bool {
	// Check required fields
	if c.Device.ID == "" {
		return false
	}
	if c.Device.Secret == "" {
		return false
	}
	if c.API.URL == "" {
		return false
	}

	return true
}

// MissingFields returns a list of missing required fields
func (c *Config) MissingFields() []string {
	var missing []string

	if c.Device.ID == "" {
		missing = append(missing, "device.id")
	}
	if c.Device.Secret == "" {
		missing = append(missing, "device.secret")
	}
	if c.API.URL == "" {
		missing = append(missing, "api.url")
	}

	return missing
}
