package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORS creates a CORS middleware
func CORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "*", // TODO: Configure in production
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Request-ID",
		AllowCredentials: false,
		ExposeHeaders:    "X-Request-ID",
	})
}
