package config

import "time"

// DefaultConfig returns a Config with default values
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:        "3000",
			Host:        "0.0.0.0",
			Mode:        "development",
			CORSOrigins: "*",
		},
		Log: LogConfig{
			Level:      "info",
			File:       "logs/api.log",
			Console:    true,
			MaxSize:    100, // 100 MB
			MaxBackups: 7,
			MaxAge:     30, // 30 days
			Compress:   true,
		},
		DB: DBConfig{
			Path:            "data/pawstream.db",
			MaxOpenConns:    10,
			MaxIdleConns:    5,
			ConnMaxLifetime: 1 * time.Hour,
		},
		JWT: JWTConfig{
			Secret:        "change-me-in-production",
			Expiry:        15 * time.Minute,
			RefreshExpiry: 30 * 24 * time.Hour,
		},
		MediaMTX: MediaMTXConfig{
			URL: "http://localhost:8554",
		},
	}
}
