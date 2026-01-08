package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

// Logger creates a request logging middleware
func Logger(log zerolog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Get request ID from context
		requestID, _ := c.Locals("request_id").(string)

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log request
		logEvent := log.Info()

		if err != nil {
			logEvent = log.Error().Err(err)
		}

		logEvent.
			Str("request_id", requestID).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", c.Response().StatusCode()).
			Dur("duration", duration).
			Str("ip", c.IP()).
			Str("user_agent", c.Get("User-Agent")).
			Msg("request completed")

		return err
	}
}
