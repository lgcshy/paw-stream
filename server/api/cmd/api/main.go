package main

import (
	"log"

	"github.com/lgc/pawstream/api/internal/app/api"
	"github.com/lgc/pawstream/api/internal/config"
	"github.com/lgc/pawstream/api/internal/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	appLogger, err := logger.Init(cfg.Log)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	appLogger.Info().Msg("Starting PawStream API Server")

	// Create and run app
	app, err := api.New(cfg, appLogger)
	if err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to create app")
	}

	if err := app.Run(); err != nil {
		appLogger.Fatal().Err(err).Msg("Server error")
	}
}
