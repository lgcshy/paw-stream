package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Check handles GET /health
func (h *HealthHandler) Check(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
