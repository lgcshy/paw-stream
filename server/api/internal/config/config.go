package config

import (
	"time"
)

// Config holds all configuration for the API server
type Config struct {
	Server        ServerConfig   `mapstructure:"server"`
	Log           LogConfig      `mapstructure:"log"`
	DB            DBConfig       `mapstructure:"db"`
	JWT           JWTConfig      `mapstructure:"jwt"`
	MediaMTX      MediaMTXConfig `mapstructure:"mediamtx"`
	EncryptionKey string         `mapstructure:"encryption_key"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port         string `mapstructure:"port"`
	Host         string `mapstructure:"host"`
	Mode         string `mapstructure:"mode"`          // development or production
	CORSOrigins  string `mapstructure:"cors_origins"`  // comma-separated allowed origins, "*" for all
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level      string `mapstructure:"level"`       // debug, info, warn, error
	File       string `mapstructure:"file"`        // log file path
	Console    bool   `mapstructure:"console"`     // enable console output
	MaxSize    int    `mapstructure:"max_size"`    // max size in MB per file
	MaxBackups int    `mapstructure:"max_backups"` // max number of old log files
	MaxAge     int    `mapstructure:"max_age"`     // max age in days
	Compress   bool   `mapstructure:"compress"`    // compress old log files
}

// DBConfig holds database configuration
type DBConfig struct {
	Path            string        `mapstructure:"path"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// JWTConfig holds JWT token configuration
type JWTConfig struct {
	Secret        string        `mapstructure:"secret"`
	Expiry        time.Duration `mapstructure:"expiry"`
	RefreshExpiry time.Duration `mapstructure:"refresh_expiry"`
}

// MediaMTXConfig holds MediaMTX integration configuration
type MediaMTXConfig struct {
	URL        string `mapstructure:"url"`         // Base URL for API (internal use)
	WebRTCURL  string `mapstructure:"webrtc_url"`  // WebRTC URL for client access
	RTSPURL    string `mapstructure:"rtsp_url"`    // RTSP URL for client access
}
