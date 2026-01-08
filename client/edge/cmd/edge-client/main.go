package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/lgc/pawstream/edge-client/internal/auth"
	"github.com/lgc/pawstream/edge-client/internal/capture"
	"github.com/lgc/pawstream/edge-client/internal/config"
	"github.com/lgc/pawstream/edge-client/internal/health"
	"github.com/lgc/pawstream/edge-client/internal/stream"
	"github.com/lgc/pawstream/edge-client/internal/webui"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sevlyar/go-daemon"
)

var (
	version   = "dev"
	buildTime = "unknown"
)

const (
	pidFileName = "edge-client.pid"
	logFileName = "edge-client.log"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "start":
		startCommand()
	case "stop":
		stopCommand()
	case "status":
		statusCommand()
	case "restart":
		restartCommand()
	case "version":
		versionCommand()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf(`PawStream Edge Client %s

Usage:
  edge-client <command> [options]

Commands:
  start       Start the edge client
  stop        Stop the edge client
  status      Show client status
  restart     Restart the edge client
  version     Show version information

Start Options:
  --config PATH       Path to config file (required)
  --daemon            Run as daemon in background
  --log-level LEVEL   Log level (debug, info, warn, error)
  --input-type TYPE   Override input type (v4l2, rtsp, file, test)

Examples:
  edge-client start --config config.yaml
  edge-client start --config config.yaml --daemon
  edge-client stop
  edge-client status

`, version)
}

func versionCommand() {
	fmt.Printf("PawStream Edge Client %s (built %s)\n", version, buildTime)
}

func startCommand() {
	fs := flag.NewFlagSet("start", flag.ExitOnError)
	configFile := fs.String("config", "", "path to config file")
	daemonMode := fs.Bool("daemon", false, "run as daemon")
	logLevel := fs.String("log-level", "", "log level")
	inputType := fs.String("input-type", "", "override input type")
	fs.Parse(os.Args[2:])

	if *configFile == "" {
		fmt.Fprintf(os.Stderr, "Error: --config is required\n")
		fs.Usage()
		os.Exit(1)
	}

	// Get absolute paths for daemon
	absConfigFile, _ := filepath.Abs(*configFile)
	workDir, _ := os.Getwd()
	pidFile := filepath.Join(workDir, pidFileName)

	if *daemonMode {
		// Check if already running
		if isRunning(pidFile) {
			fmt.Println("Edge client is already running")
			os.Exit(1)
		}

		// Setup daemon context
		cntxt := &daemon.Context{
			PidFileName: pidFile,
			PidFilePerm: 0644,
			LogFileName: filepath.Join(workDir, logFileName),
			LogFilePerm: 0640,
			WorkDir:     workDir,
			Umask:       027,
		}

		// Daemonize
		d, err := cntxt.Reborn()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to daemonize: %v\n", err)
			os.Exit(1)
		}
		if d != nil {
			// Parent process
			fmt.Printf("Edge client started in background (PID: %d)\n", d.Pid)
			fmt.Printf("Log file: %s\n", filepath.Join(workDir, logFileName))
			return
		}
		defer cntxt.Release()

		// Child process continues here
		*configFile = absConfigFile
	}

	// Run the client
	runClient(*configFile, *logLevel, *inputType)
}

func stopCommand() {
	workDir, _ := os.Getwd()
	pidFile := filepath.Join(workDir, pidFileName)

	if !isRunning(pidFile) {
		fmt.Println("Edge client is not running")
		return
	}

	// Read PID
	data, err := os.ReadFile(pidFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read PID file: %v\n", err)
		os.Exit(1)
	}

	var pid int
	fmt.Sscanf(string(data), "%d", &pid)

	// Send SIGTERM
	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find process: %v\n", err)
		os.Exit(1)
	}

	if err := process.Signal(syscall.SIGTERM); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to stop process: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Edge client stopped (PID: %d)\n", pid)

	// Wait for process to exit and remove PID file
	time.Sleep(1 * time.Second)
	os.Remove(pidFile)
}

func statusCommand() {
	workDir, _ := os.Getwd()
	pidFile := filepath.Join(workDir, pidFileName)

	if !isRunning(pidFile) {
		fmt.Println("Edge client is not running")
		return
	}

	data, err := os.ReadFile(pidFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read PID file: %v\n", err)
		os.Exit(1)
	}

	var pid int
	fmt.Sscanf(string(data), "%d", &pid)

	fmt.Printf("Edge client is running (PID: %d)\n", pid)
	fmt.Printf("PID file: %s\n", pidFile)
	fmt.Printf("Log file: %s\n", filepath.Join(workDir, logFileName))
}

func restartCommand() {
	stopCommand()
	time.Sleep(2 * time.Second)
	startCommand()
}

func isRunning(pidFile string) bool {
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return false
	}

	var pid int
	fmt.Sscanf(string(data), "%d", &pid)

	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// Send signal 0 to check if process exists
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

func runClient(configFile, logLevel, inputType string) {
	// Load configuration
	cfg, err := config.Load(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Override log level if specified
	if logLevel != "" {
		cfg.Log.Level = logLevel
	}

	// Override input type if specified
	if inputType != "" {
		cfg.Input.Type = inputType
	}

	// Initialize logger
	initLogger(cfg.Log)

	log.Info().
		Str("version", version).
		Str("build_time", buildTime).
		Str("device_id", cfg.Device.ID).
		Msg("Starting PawStream Edge Client")

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Info().Str("signal", sig.String()).Msg("Received shutdown signal")
		cancel()
	}()

	// Create SSE manager for Web UI
	sseManager := webui.NewSSEManager(log.Logger)
	defer sseManager.Close()

	// Create status function for Web UI
	var framesPushed atomic.Int64
	startTime := time.Now()
	statusFunc := func() interface{} {
		return map[string]interface{}{
			"client_status": "running",
			"stream_status": "active",
			"uptime":        time.Since(startTime).String(),
			"frames_pushed": framesPushed.Load(),
		}
	}

	// Create Web UI handler
	webuiHandler := webui.NewHandler(configFile, statusFunc, cfg.API.URL, log.Logger)
	webuiHandler.SetReloadChan(make(chan bool, 1))

	// Start Web UI server if enabled
	var webuiServer *webui.Server
	if cfg.WebUI != nil && cfg.WebUI.Enabled {
		var authConfig *webui.AuthConfig
		if cfg.WebUI.Auth != nil && cfg.WebUI.Auth.Enabled {
			authConfig = &webui.AuthConfig{
				Enabled:  true,
				Username: cfg.WebUI.Auth.Username,
				Password: cfg.WebUI.Auth.Password,
			}
		}

		webuiCfg := webui.Config{
			Host:       cfg.WebUI.Host,
			Port:       cfg.WebUI.Port,
			AuthConfig: authConfig,
		}

		webuiServer = webui.NewServer(webuiCfg, webuiHandler, sseManager, log.Logger)
		if err := webuiServer.Start(); err != nil {
			log.Error().Err(err).Msg("Failed to start Web UI server")
		} else {
			defer webuiServer.Stop()
		}
	}

	// Start configuration file watcher
	configWatcher, err := config.NewWatcher(configFile, log.Logger)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create config watcher")
	} else {
		configWatcher.Start()
		defer configWatcher.Stop()

		// Handle config reload
		go func() {
			for range configWatcher.ReloadChan() {
				log.Info().Msg("Configuration file changed, reloading...")
				
				// Reload config
				newCfg, err := config.Load(configFile)
				if err != nil {
					log.Error().Err(err).Msg("Failed to reload configuration")
					continue
				}

				// Update config (thread-safe)
				cfg = newCfg
				log.Info().Msg("Configuration reloaded successfully")

				// Notify SSE clients
				sseManager.BroadcastConfigChange()
			}
		}()
	}

	// Create API client
	apiClient := auth.NewClient(cfg.API.URL, cfg.Device.ID, cfg.Device.Secret, cfg.API.Timeout)

	// Authenticate with API and get device info
	log.Info().Msg("Authenticating with API server...")
	deviceInfo, err := apiClient.GetDeviceInfo()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to authenticate with API")
	}

	if deviceInfo.Disabled {
		log.Fatal().Msg("Device is disabled on the server")
	}

	log.Info().
		Str("name", deviceInfo.Name).
		Str("location", deviceInfo.Location).
		Str("publish_path", deviceInfo.PublishPath).
		Msg("Device authenticated successfully")

	// Update stream URL if not configured
	if cfg.Stream.URL == "" {
		cfg.Stream.URL = "localhost:8554"
		log.Info().Str("url", cfg.Stream.URL).Msg("Using default MediaMTX address")
	}

	// Create input source
	var inputSource capture.InputSource

	switch cfg.Input.Type {
	case "test":
		inputSource = capture.NewTestSource(cfg.Video.Width, cfg.Video.Height, cfg.Video.Framerate)
	case "file":
		inputSource = capture.NewFileSource(cfg.Input.Source, true)
	case "v4l2":
		inputSource = capture.NewV4L2Source(cfg.Input.Source, cfg.Video.Width, cfg.Video.Height, cfg.Video.Framerate)
	case "rtsp":
		inputSource = capture.NewRTSPSource(cfg.Input.Source, "tcp")
	default:
		log.Fatal().Str("type", cfg.Input.Type).Msg("Unknown input type")
	}

	// Validate input source
	if err := inputSource.Validate(); err != nil {
		log.Fatal().Err(err).Msg("Invalid input source")
	}

	log.Info().Str("source", inputSource.String()).Msg("Input source configured")

	// Build output RTSP URL with authentication
	outputURL := fmt.Sprintf("rtsp://%s:%s@%s/%s",
		cfg.Device.ID,
		cfg.Device.Secret,
		cfg.Stream.URL,
		deviceInfo.PublishPath)

	// Create stream manager
	streamCfg := stream.Config{
		VideoCodec:           cfg.Video.Codec,
		VideoBitrate:         cfg.Video.Bitrate,
		Preset:               "ultrafast",
		ReconnectInterval:    cfg.Stream.ReconnectInterval,
		MaxReconnectAttempts: cfg.Stream.MaxReconnectAttempts,
	}

	streamMgr := stream.NewManager(inputSource, outputURL, streamCfg, log.Logger)

	// Start streaming
	log.Info().Msg("Starting video stream...")
	if err := streamMgr.Start(ctx); err != nil {
		log.Fatal().Err(err).Msg("Failed to start streaming")
	}
	defer streamMgr.Stop()

	log.Info().Msg("Streaming started successfully")

	// Start health check server if enabled
	var healthServer *health.Server
	if cfg.Health.Enabled {
		healthServer = health.NewServer(cfg.Health.Address, streamMgr, log.Logger)
		go func() {
			if err := healthServer.Start(); err != nil {
				log.Error().Err(err).Msg("Health server error")
			}
		}()
		log.Info().Str("address", cfg.Health.Address).Msg("Health check endpoint started")
	}

	// Periodically broadcast status to Web UI
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Update frame count (simulated)
				framesPushed.Add(60) // Assume 30 fps * 2 seconds

				// Broadcast status
				status := statusFunc()
				sseManager.BroadcastStatus(status)

			case <-ctx.Done():
				return
			}
		}
	}()

	// Main loop
	log.Info().Msg("Edge client running (press Ctrl+C to stop)")

	// Wait for shutdown signal
	<-ctx.Done()

	// Graceful shutdown
	log.Info().Msg("Shutting down...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer shutdownCancel()

	// Stop streaming
	log.Info().Msg("Stopping stream...")
	if err := streamMgr.Stop(); err != nil {
		log.Error().Err(err).Msg("Error stopping stream")
	}

	// Stop health check server
	if healthServer != nil {
		if err := healthServer.Stop(); err != nil {
			log.Error().Err(err).Msg("Error stopping health server")
		}
	}

	// Stop Web UI server
	if webuiServer != nil {
		log.Info().Msg("Stopping Web UI server...")
		if err := webuiServer.Stop(); err != nil {
			log.Error().Err(err).Msg("Error stopping Web UI server")
		}
	}

	// Stop SSE manager
	log.Info().Msg("Closing SSE connections...")
	sseManager.Close()

	// Stop config watcher
	if configWatcher != nil {
		configWatcher.Stop()
	}

	<-shutdownCtx.Done()
	log.Info().Msg("Shutdown complete")
}

// initLogger initializes the logger based on configuration
func initLogger(cfg config.LogConfig) {
	// Set log level
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Set output
	if cfg.File != "" {
		file, err := os.OpenFile(cfg.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", err)
			log.Logger = log.Output(os.Stdout)
		} else {
			log.Logger = log.Output(file)
		}
	} else {
		log.Logger = log.Output(os.Stdout)
	}

	// Set format
	if cfg.Format == "console" {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	}
}
