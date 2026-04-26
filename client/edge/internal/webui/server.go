package webui

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
)

// Server represents the Web UI HTTP server
type Server struct {
	app        *fiber.App
	host       string
	port       int
	authConfig *AuthConfig
	handler    *Handler
	sseManager *SSEManager
	logger     zerolog.Logger
}

// AuthConfig holds HTTP Basic Auth configuration
type AuthConfig struct {
	Enabled  bool
	Username string
	Password string
}

// Config holds Web UI server configuration
type Config struct {
	Host       string
	Port       int
	AuthConfig *AuthConfig
}

// NewServer creates a new Web UI server
func NewServer(cfg Config, handler *Handler, sseManager *SSEManager, logger zerolog.Logger) *Server {
	app := fiber.New(fiber.Config{
		ServerHeader:          "PawStream Edge Client",
		DisableStartupMessage: true,
		ErrorHandler:          errorHandler,
		ReadTimeout:           10 * time.Second,
		WriteTimeout:          0, // Disable write timeout for SSE streaming
		IdleTimeout:           120 * time.Second,
		StreamRequestBody:     true, // Enable streaming
	})

	// Middleware
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Optional: HTTP Basic Auth
	if cfg.AuthConfig != nil && cfg.AuthConfig.Enabled {
		app.Use(basicauth.New(basicauth.Config{
			Users: map[string]string{
				cfg.AuthConfig.Username: cfg.AuthConfig.Password,
			},
			Realm: "PawStream Edge Client",
		}))
	}

	s := &Server{
		app:        app,
		host:       cfg.Host,
		port:       cfg.Port,
		authConfig: cfg.AuthConfig,
		handler:    handler,
		sseManager: sseManager,
		logger:     logger,
	}

	s.setupRoutes()

	return s
}

// setupRoutes sets up all HTTP routes
func (s *Server) setupRoutes() {
	// Single page app - all pages served from index.html
	s.app.Get("/", serveIndex)
	s.app.Static("/", "./web", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        false,
		CacheDuration: 24 * time.Hour,
	})

	// API routes
	api := s.app.Group("/api")

	// Quick setup (one-click configuration)
	api.Post("/quick-setup", s.handler.QuickSetup)

	// Configuration
	api.Get("/config", s.handler.GetConfig)
	api.Post("/config", s.handler.SaveConfig)

	// Status and system info
	api.Get("/status", s.handler.GetStatus)
	api.Get("/system/info", s.handler.GetSystemInfo)

	// API server validation
	api.Post("/validate-server", s.handler.ValidateAPIServer)

	// User login (proxy to API server)
	api.Post("/login", s.handler.Login)

	// Device management (proxy to API server)
	api.Get("/devices", s.handler.GetDevices)
	api.Post("/devices", s.handler.CreateDevice)

	// Logs
	api.Get("/logs/recent", s.handler.GetRecentLogs)

	// SSE endpoint
	api.Get("/events", s.sseManager.Handler)
}

// Start starts the Web UI server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	s.logger.Info().
		Str("host", s.host).
		Int("port", s.port).
		Bool("auth", s.authConfig != nil && s.authConfig.Enabled).
		Msg("Starting Web UI server")

	go func() {
		if err := s.app.Listen(addr); err != nil {
			s.logger.Error().Err(err).Msg("Web UI server error")
		}
	}()

	// Wait a bit for server to start
	time.Sleep(100 * time.Millisecond)

	// Display appropriate URL based on host
	displayHost := s.host
	if s.host == "0.0.0.0" || s.host == "" {
		displayHost = "localhost"
		s.logger.Info().
			Str("local_url", fmt.Sprintf("http://localhost:%d", s.port)).
			Str("network_url", fmt.Sprintf("http://<your-ip>:%d", s.port)).
			Msg("Web UI available on all interfaces")
	} else {
		s.logger.Info().
			Str("url", fmt.Sprintf("http://%s:%d", displayHost, s.port)).
			Msg("Web UI available")
	}

	return nil
}

// Stop stops the Web UI server gracefully
func (s *Server) Stop() error {
	s.logger.Info().Msg("Stopping Web UI server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.app.ShutdownWithContext(ctx); err != nil {
		return err
	}

	s.logger.Info().Msg("Web UI server stopped")
	return nil
}

// serveIndex serves the index.html file
func serveIndex(c *fiber.Ctx) error {
	return c.SendFile("./web/index.html")
}

// errorHandler handles fiber errors
func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"error":   true,
		"message": message,
	})
}
