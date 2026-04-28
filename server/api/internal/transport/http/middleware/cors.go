package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORS creates a CORS middleware with configurable origins
func CORS(allowOrigins string) fiber.Handler {
	if allowOrigins == "" {
		allowOrigins = "*"
	}
	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Request-ID",
		AllowCredentials: allowOrigins != "*",
		ExposeHeaders:    "X-Request-ID",
	})
}
