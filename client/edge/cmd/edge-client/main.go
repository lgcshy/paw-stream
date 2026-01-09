package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/lgc/pawstream/edge-client/internal/auth"
	"github.com/lgc/pawstream/edge-client/internal/capture"
	"github.com/lgc/pawstream/edge-client/internal/config"
	"github.com/lgc/pawstream/edge-client/internal/health"
	"github.com/lgc/pawstream/edge-client/internal/input"
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
  --config PATH           Path to config file (optional, auto-search if not specified)
  --daemon                Run as daemon in background
  
  Core Configuration (overrides config file):
  --device-id ID          Device ID
  --device-secret SECRET  Device secret
  --api-url URL           API server URL
  
  Input Configuration (overrides config file):
  --input-type TYPE       Input type: v4l2, rtsp, file, test
  --input-source SOURCE   Input source path or URL
  --mediamtx-url URL      MediaMTX server URL
  
  Other Options:
  --log-level LEVEL       Log level: debug, info, warn, error

Configuration File Search Order:
  1. ./config.yaml
  2. ./configs/config.yaml
  3. ~/.pawstream/config.yaml

Examples:
  # Small user (zero-config)
  edge-client start                    # Auto-search config or enter setup wizard
  edge-client start --daemon           # Run in background
  
  # Developer (advanced)
  edge-client start --config my.yaml   # Use custom config
  edge-client start --device-id xxx --device-secret yyy --api-url http://...
  edge-client start --config my.yaml --input-type test
  
  # Management
  edge-client stop
  edge-client status
  edge-client restart

`, version)
}

func versionCommand() {
	fmt.Printf("PawStream Edge Client %s (built %s)\n", version, buildTime)
}

func startCommand() {
	fs := flag.NewFlagSet("start", flag.ExitOnError)
	
	// File and mode options
	configFile := fs.String("config", "", "path to config file (optional, will auto-search if not specified)")
	daemonMode := fs.Bool("daemon", false, "run as daemon")
	
	// Core configuration overrides
	deviceID := fs.String("device-id", "", "device ID (overrides config file)")
	deviceSecret := fs.String("device-secret", "", "device secret (overrides config file)")
	apiURL := fs.String("api-url", "", "API server URL (overrides config file)")
	
	// Input configuration overrides
	inputType := fs.String("input-type", "", "input type: v4l2, rtsp, file, test (overrides config file)")
	inputSource := fs.String("input-source", "", "input source path or URL (overrides config file)")
	
	// MediaMTX configuration
	mediamtxURL := fs.String("mediamtx-url", "", "MediaMTX server URL (overrides config file)")
	
	// Other options
	logLevel := fs.String("log-level", "", "log level: debug, info, warn, error")
	
	fs.Parse(os.Args[2:])

	// Auto-search for config file if not specified
	if *configFile == "" {
		*configFile = config.FindConfigFile()
		if *configFile != "" {
			fmt.Printf("📄 Found configuration file: %s\n", *configFile)
		}
	}

	// Check if config file exists
	if *configFile == "" {
		fmt.Println("⚙️  No configuration file found")
		fmt.Println("🌐 Starting setup wizard...")
		fmt.Println()
		runSetupWizard()
		return
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

	// Prepare config overrides
	overrides := &configOverrides{
		deviceID:     *deviceID,
		deviceSecret: *deviceSecret,
		apiURL:       *apiURL,
		inputType:    *inputType,
		inputSource:  *inputSource,
		mediamtxURL:  *mediamtxURL,
		logLevel:     *logLevel,
	}

	// Run the client
	runClientWithOverrides(*configFile, overrides)
}

// configOverrides holds command-line configuration overrides
type configOverrides struct {
	deviceID     string
	deviceSecret string
	apiURL       string
	inputType    string
	inputSource  string
	mediamtxURL  string
	logLevel     string
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

func runClientWithOverrides(configFile string, overrides *configOverrides) {
	// Load configuration
	cfg, err := config.Load(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Apply command-line overrides (highest priority)
	if overrides.deviceID != "" {
		cfg.Device.ID = overrides.deviceID
		fmt.Printf("🔧 Override: device.id = %s\n", overrides.deviceID)
	}
	if overrides.deviceSecret != "" {
		cfg.Device.Secret = overrides.deviceSecret
		fmt.Printf("🔧 Override: device.secret = ***\n")
	}
	if overrides.apiURL != "" {
		cfg.API.URL = overrides.apiURL
		fmt.Printf("🔧 Override: api.url = %s\n", overrides.apiURL)
	}
	if overrides.inputType != "" {
		cfg.Input.Type = overrides.inputType
		fmt.Printf("🔧 Override: input.type = %s\n", overrides.inputType)
	}
	if overrides.inputSource != "" {
		cfg.Input.Source = overrides.inputSource
		fmt.Printf("🔧 Override: input.source = %s\n", overrides.inputSource)
	}
	if overrides.mediamtxURL != "" {
		cfg.Stream.URL = overrides.mediamtxURL
		fmt.Printf("🔧 Override: mediamtx.url = %s\n", overrides.mediamtxURL)
	}
	if overrides.logLevel != "" {
		cfg.Log.Level = overrides.logLevel
		fmt.Printf("🔧 Override: log.level = %s\n", overrides.logLevel)
	}

	// Auto-detect input source if needed
	configModified := false
	if cfg.Input.Type == "auto" || cfg.Input.Source == "auto" {
		fmt.Printf("🔍 Auto-detecting input source...\n")
		detectedType, detectedSource, err := input.AutoDetectInput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Auto-detection failed: %v\n", err)
			fmt.Fprintf(os.Stderr, "Falling back to test pattern\n")
			cfg.Input.Type = "test"
			cfg.Input.Source = "pattern=smpte"
		} else {
			cfg.Input.Type = detectedType
			cfg.Input.Source = detectedSource
			fmt.Printf("✅ Detected: %s -> %s\n", detectedType, detectedSource)
		}
		configModified = true
	}

	// Default to GStreamer if engine not specified
	if cfg.Stream.Engine == "" {
		cfg.Stream.Engine = "gstreamer"
		fmt.Printf("🎬 Using default engine: gstreamer\n")
		configModified = true
	}

	// Save updated configuration if modified
	if configModified {
		if err := config.Save(configFile, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Warning: Failed to save updated configuration: %v\n", err)
		} else {
			fmt.Printf("💾 Configuration updated and saved\n")
		}
	}

	// Check configuration completeness (after applying overrides)
	if !cfg.IsComplete() {
		missing := cfg.MissingFields()
		fmt.Fprintf(os.Stderr, "❌ Configuration is incomplete\n")
		fmt.Fprintf(os.Stderr, "Missing required fields: %v\n", missing)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Please run without --config to use setup wizard,\n")
		fmt.Fprintf(os.Stderr, "or provide missing fields via command-line arguments.\n")
		os.Exit(1)
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
	// Extract host:port from cfg.Stream.URL (e.g., "rtsp://localhost:8554" -> "localhost:8554")
	streamHost := strings.TrimPrefix(cfg.Stream.URL, "rtsp://")
	streamHost = strings.TrimPrefix(streamHost, "rtmps://")
	outputURL := fmt.Sprintf("rtsp://%s:%s@%s/%s",
		cfg.Device.ID,
		cfg.Device.Secret,
		streamHost,
		deviceInfo.PublishPath)

	// Create stream manager
	streamCfg := stream.Config{
		Engine:                   stream.EngineType(cfg.Stream.Engine),
		VideoCodec:               cfg.Video.Codec,
		VideoBitrate:             cfg.Video.Bitrate,
		VideoWidth:               cfg.Video.Width,
		VideoHeight:              cfg.Video.Height,
		VideoFramerate:           cfg.Video.Framerate,
		FFmpegPreset:             cfg.Stream.FFmpeg.Preset,
		FFmpegTune:               cfg.Stream.FFmpeg.Tune,
		FFmpegHWAccel:            cfg.Stream.FFmpeg.HWAccel,
		GStreamerLatencyMs:       cfg.Stream.GStreamer.LatencyMs,
		GStreamerUseHardware:     cfg.Stream.GStreamer.UseHardware,
		GStreamerBufferSize:      cfg.Stream.GStreamer.BufferSize,
		ReconnectInterval:        cfg.Stream.ReconnectInterval,
		MaxReconnectAttempts:     cfg.Stream.MaxReconnectAttempts,
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

// runSetupWizard runs the setup wizard mode
func runSetupWizard() {
	fmt.Println("┌─────────────────────────────────────────┐")
	fmt.Println("│  🐾 PawStream Edge Client Setup Wizard │")
	fmt.Println("└─────────────────────────────────────────┘")
	fmt.Println()
	fmt.Println("Welcome! Let's set up your edge client.")
	fmt.Println()

	// Default config file location
	defaultConfigPath := "./config.yaml"

	fmt.Printf("Configuration will be saved to: %s\n", defaultConfigPath)
	fmt.Println()

	// Start Web UI for setup
	fmt.Println("🌐 Starting Web UI setup wizard...")
	fmt.Println()

	// Initialize basic logger for setup wizard
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Create minimal config for Web UI only
	setupConfig := &config.Config{
		WebUI: &config.WebUIConfig{
			Enabled: true,
			Host:    "0.0.0.0",
			Port:    8088,
			Auth:    nil,
		},
		Log: config.LogConfig{
			Level:  "info",
			Format: "console",
		},
	}

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		fmt.Printf("\n\n⚠️  Received %s signal, shutting down...\n", sig.String())
		cancel()
	}()

	// Create SSE manager
	sseManager := webui.NewSSEManager(log.Logger)
	defer sseManager.Close()

	// Create status function (placeholder for setup mode)
	statusFunc := func() interface{} {
		return map[string]interface{}{
			"mode":   "setup",
			"status": "waiting_for_configuration",
		}
	}

	// Create Web UI handler with setup mode
	webuiHandler := webui.NewHandler(defaultConfigPath, statusFunc, "", log.Logger)

	// Start Web UI server
	webuiCfg := webui.Config{
		Host:       setupConfig.WebUI.Host,
		Port:       setupConfig.WebUI.Port,
		AuthConfig: nil, // No auth in setup mode
	}

	webuiServer := webui.NewServer(webuiCfg, webuiHandler, sseManager, log.Logger)
	if err := webuiServer.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to start Web UI: %v\n", err)
		os.Exit(1)
	}
	defer webuiServer.Stop()

	fmt.Println("✅ Web UI started successfully!")
	fmt.Println()
	fmt.Printf("🔗 Open in browser: \033[1;36mhttp://localhost:%d/setup\033[0m\n", setupConfig.WebUI.Port)
	fmt.Println()
	fmt.Println("📋 Setup steps:")
	fmt.Println("   1. Connect to API server")
	fmt.Println("   2. Login and select device")
	fmt.Println("   3. Configure input source")
	fmt.Println("   4. Save configuration and start")
	fmt.Println()
	fmt.Println("💡 Tip: Configuration will be saved automatically")
	fmt.Println("        You can edit it later in the Web UI")
	fmt.Println()

	// Try to open browser automatically
	openBrowser(fmt.Sprintf("http://localhost:%d/setup", setupConfig.WebUI.Port))

	// Wait for configuration to be completed
	fmt.Println("⏳ Waiting for configuration...")
	fmt.Println("   (Press Ctrl+C to exit)")
	fmt.Println()

	// Watch for config file creation
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Check if config file was created
			if _, err := os.Stat(defaultConfigPath); err == nil {
				// Config file exists, check if complete
				cfg, err := config.Load(defaultConfigPath)
				if err != nil {
					continue
				}

				if cfg.IsComplete() {
					fmt.Println()
					fmt.Println("✅ Configuration completed!")
					fmt.Println()
					fmt.Print("🚀 Start streaming now? [Y/n]: ")

					var response string
					fmt.Scanln(&response)

					if response == "" || response == "Y" || response == "y" {
						fmt.Println()
						fmt.Println("Starting edge client...")
						fmt.Println()

						// Stop Web UI
						webuiServer.Stop()
						sseManager.Close()

						// Run client with new config
						runClientWithOverrides(defaultConfigPath, &configOverrides{})
						return
					} else {
						fmt.Println()
						fmt.Println("✅ Configuration saved!")
						fmt.Printf("💡 Run \033[1;36m./edge-client start\033[0m to start streaming\n")
						fmt.Println()
						return
					}
				}
			}

		case <-ctx.Done():
			fmt.Println()
			fmt.Println("Setup wizard cancelled.")
			return
		}
	}
}

// openBrowser attempts to open the default browser
func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime := os.Getenv("GOOS"); runtime {
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	default:
		// Can't open browser, just return
		return
	}

	// Try to open browser (non-blocking, ignore errors)
	go func() {
		exec.Command(cmd, args...).Start()
	}()
}
