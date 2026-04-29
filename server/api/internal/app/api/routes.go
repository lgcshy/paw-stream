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

	// Internal endpoints (called by MediaMTX hooks)
	a.fiber.Post("/internal/stream-closed", a.mediamtxHandler.StreamClosed)

	// API routes
	api := a.fiber.Group("/api")

	// Public config endpoint (no auth required)
	api.Get("/config", a.configHandler.GetConfig)

	// Auth endpoints (no auth required, rate limited)
	authLimiter := middleware.RateLimitAuth()
	api.Post("/register", authLimiter, a.authHandler.Register)
	api.Post("/login", authLimiter, a.authHandler.Login)
	api.Post("/refresh", authLimiter, a.authHandler.Refresh)

	// Device auth endpoint (for edge client, no auth required)
	api.Post("/device/auth", a.deviceHandler.AuthDevice)

	// Protected routes (require JWT auth)
	protected := api.Group("", middleware.AuthUser(a.cfg.JWT.Secret))

	// User info
	protected.Get("/me", a.authHandler.GetMe)
	protected.Post("/me/avatar", a.authHandler.UploadAvatar)

	// Avatar serving (public)
	api.Get("/avatars/:id", a.authHandler.GetAvatar)

	// Device management
	protected.Get("/devices", a.deviceHandler.List)
	protected.Post("/devices", a.deviceHandler.Create)
	protected.Get("/devices/:id", a.deviceHandler.Get)
	protected.Put("/devices/:id", a.deviceHandler.Update)
	protected.Delete("/devices/:id", a.deviceHandler.Delete)
	protected.Post("/devices/:id/rotate-secret", a.deviceHandler.RotateSecret)

	// Device sharing
	protected.Post("/devices/:id/share", a.deviceHandler.ShareDevice)
	protected.Delete("/devices/:id/share/:userId", a.deviceHandler.UnshareDevice)
	protected.Get("/devices/:id/shares", a.deviceHandler.ListShares)

	// Path query
	protected.Get("/paths", a.pathHandler.List)
}
