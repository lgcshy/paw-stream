package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

// Recovery creates a panic recovery middleware
func Recovery(log zerolog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				// Get request ID
				requestID, _ := c.Locals("request_id").(string)

				// Get stack trace
				stackTrace := string(debug.Stack())

				// Log panic
				log.Error().
					Str("request_id", requestID).
					Str("panic", fmt.Sprintf("%v", r)).
					Str("stack_trace", stackTrace).
					Msg("panic recovered")

				// Return error response
				c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":      "internal_server_error",
					"message":    "Internal server error",
					"request_id": requestID,
				})
			}
		}()

		return c.Next()
	}
}
