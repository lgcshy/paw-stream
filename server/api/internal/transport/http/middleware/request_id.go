package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-ID"

// RequestID generates and injects a unique request ID
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if request ID already exists in header
		requestID := c.Get(RequestIDHeader)
		if requestID == "" {
			// Generate new request ID
			requestID = uuid.New().String()
		}

		// Set request ID in context and response header
		c.Locals("request_id", requestID)
		c.Set(RequestIDHeader, requestID)

		return c.Next()
	}
}
