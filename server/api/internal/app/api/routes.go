package api

import (
	"github.com/lgc/pawstream/api/internal/transport/http/middleware"
)

// setupRoutes sets up all HTTP routes
func (a *App) setupRoutes() {
	// Health check (no auth required)
	a.fiber.Get("/health", a.healthHandler.Check)

	// MediaMTX auth callback (no auth required)
	a.fiber.Post("/mediamtx/auth", a.mediamtxHandler.Auth)

	// API routes
	api := a.fiber.Group("/api")

	// Public config endpoint (no auth required)
	api.Get("/config", a.configHandler.GetConfig)

	// Auth endpoints (no auth required)
	api.Post("/register", a.authHandler.Register)
	api.Post("/login", a.authHandler.Login)

	// Device auth endpoint (for edge client, no auth required)
	api.Post("/device/auth", a.deviceHandler.AuthDevice)

	// Protected routes (require JWT auth)
	protected := api.Group("", middleware.AuthUser(a.cfg.JWT.Secret))

	// User info
	protected.Get("/me", a.authHandler.GetMe)

	// Device management
	protected.Get("/devices", a.deviceHandler.List)
	protected.Post("/devices", a.deviceHandler.Create)
	protected.Get("/devices/:id", a.deviceHandler.Get)
	protected.Put("/devices/:id", a.deviceHandler.Update)
	protected.Delete("/devices/:id", a.deviceHandler.Delete)
	protected.Post("/devices/:id/rotate-secret", a.deviceHandler.RotateSecret)

	// Path query
	protected.Get("/paths", a.pathHandler.List)
}
