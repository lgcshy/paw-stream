package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Load loads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	// Start with default config
	cfg := DefaultConfig()

	// Setup viper
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// Add config paths
	if configPath != "" {
		v.AddConfigPath(configPath)
	}
	v.AddConfigPath(".")
	v.AddConfigPath("./")

	// Environment variables
	v.SetEnvPrefix("PAWSTREAM")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Read config file (optional, won't fail if not found)
	if err := v.ReadInConfig(); err != nil {
		// Config file not found is OK, we'll use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal into config struct
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate critical configuration
	if err := validate(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// validate checks critical configuration values
func validate(cfg *Config) error {
	if cfg.Server.Port == "" {
		return fmt.Errorf("server.port is required")
	}

	if cfg.JWT.Secret == "change-me-in-production" && cfg.Server.Mode == "production" {
		return fmt.Errorf("JWT secret must be changed in production mode")
	}

	if cfg.DB.Path == "" {
		return fmt.Errorf("db.path is required")
	}

	return nil
}
