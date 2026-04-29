package api

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"

	"github.com/lgc/pawstream/api/internal/config"
	"github.com/lgc/pawstream/api/internal/domain/acl"
	"github.com/lgc/pawstream/api/internal/domain/device"
	"github.com/lgc/pawstream/api/internal/domain/user"
	"github.com/lgc/pawstream/api/internal/store/sqlite"
	"github.com/lgc/pawstream/api/internal/transport/http/handlers"
	"github.com/lgc/pawstream/api/internal/transport/http/middleware"
)

// App represents the API application
type App struct {
	cfg   *config.Config
	fiber *fiber.App
	db    *sqlite.DB
	log   zerolog.Logger

	// Services
	userService   *user.Service
	deviceService *device.Service
	aclService    *acl.Service
	shareRepo     *sqlite.DeviceShareRepository

	// Handlers
	healthHandler   *handlers.HealthHandler
	authHandler     *handlers.AuthHandler
	deviceHandler   *handlers.DeviceHandler
	pathHandler     *handlers.PathHandler
	mediamtxHandler *handlers.MediaMTXHandler
	configHandler   *handlers.ConfigHandler
	adminHandler    *handlers.AdminHandler
	metrics         *handlers.Metrics
}

// New creates a new API application
func New(cfg *config.Config, log zerolog.Logger) (*App, error) {
	app := &App{
		cfg: cfg,
		log: log,
	}

	// Initialize database
	if err := app.initDB(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize services
	app.initServices()

	// Initialize handlers
	app.initHandlers()

	// Initialize Fiber app
	app.initFiber()

	// Setup routes
	app.setupRoutes()

	return app, nil
}

// initDB initializes the database connection
func (a *App) initDB() error {
	db, err := sqlite.New(a.cfg.DB)
	if err != nil {
		return err
	}

	a.db = db
	a.log.Info().Str("path", a.cfg.DB.Path).Msg("database connected")

	// Run migrations
	if err := db.RunMigrations("migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	a.log.Info().Msg("database migrations completed")

	return nil
}

// initServices initializes domain services
func (a *App) initServices() {
	// Create repositories
	userRepo := sqlite.NewUserRepository(a.db)
	deviceRepo := sqlite.NewDeviceRepository(a.db)
	shareRepo := sqlite.NewDeviceShareRepository(a.db)

	// Create services
	a.userService = user.NewService(userRepo)
	a.deviceService = device.NewService(deviceRepo, a.cfg.EncryptionKey)
	a.aclService = acl.NewService(a.userService, a.deviceService, shareRepo, a.cfg.JWT.Secret)
	a.shareRepo = shareRepo

	a.log.Info().Msg("services initialized")
}

// initHandlers initializes HTTP handlers
func (a *App) initHandlers() {
	refreshRepo := sqlite.NewRefreshTokenRepository(a.db)
	a.healthHandler = handlers.NewHealthHandler()
	a.authHandler = handlers.NewAuthHandler(a.userService, refreshRepo, a.cfg.JWT.Secret, a.cfg.JWT.Expiry, a.cfg.JWT.RefreshExpiry)
	a.deviceHandler = handlers.NewDeviceHandler(a.deviceService, a.userService, a.shareRepo)
	a.pathHandler = handlers.NewPathHandler(a.deviceService)
	a.mediamtxHandler = handlers.NewMediaMTXHandler(a.aclService, a.deviceService, a.log)
	a.configHandler = handlers.NewConfigHandler(a.cfg.MediaMTX.WebRTCURL, a.cfg.MediaMTX.RTSPURL)
	a.adminHandler = handlers.NewAdminHandler(a.deviceService, a.userService)
	a.metrics = handlers.NewMetrics()

	a.log.Info().Msg("handlers initialized")
}

// initFiber initializes the Fiber app with middleware
func (a *App) initFiber() {
	app := fiber.New(fiber.Config{
		AppName:      "PawStream API",
		ServerHeader: "PawStream",
		ErrorHandler: a.errorHandler,
	})

	// Global middleware
	app.Use(middleware.Recovery(a.log))
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger(a.log))
	app.Use(middleware.CORS(a.cfg.Server.CORSOrigins))
	app.Use(a.metrics.Middleware())

	a.fiber = app
	a.log.Info().Msg("fiber app initialized")
}

// errorHandler handles Fiber errors
func (a *App) errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	requestID, _ := c.Locals("request_id").(string)

	return c.Status(code).JSON(fiber.Map{
		"error":      "error",
		"message":    err.Error(),
		"request_id": requestID,
	})
}

// Run starts the API server
func (a *App) Run() error {
	// Start server in goroutine
	go func() {
		addr := fmt.Sprintf("%s:%s", a.cfg.Server.Host, a.cfg.Server.Port)
		a.log.Info().Str("addr", addr).Msg("starting server")

		if err := a.fiber.Listen(addr); err != nil {
			a.log.Fatal().Err(err).Msg("server error")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.log.Info().Msg("shutting down server...")

	// Graceful shutdown
	return a.Shutdown()
}

// Shutdown gracefully shuts down the server
func (a *App) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown Fiber server
	if err := a.fiber.ShutdownWithContext(ctx); err != nil {
		a.log.Error().Err(err).Msg("error shutting down fiber")
	}

	// Close database connection
	if err := a.db.Close(); err != nil {
		a.log.Error().Err(err).Msg("error closing database")
	}

	a.log.Info().Msg("server shut down successfully")

	return nil
}
