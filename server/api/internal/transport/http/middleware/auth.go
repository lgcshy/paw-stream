package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/lgc/pawstream/api/internal/pkg/jwtutil"
)

// AuthUser validates JWT token and injects user info into context
func AuthUser(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "unauthorized",
				"message": "Missing authorization header",
			})
		}

		// Check Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "unauthorized",
				"message": "Invalid authorization header format",
			})
		}

		token := parts[1]

		// Validate token
		claims, err := jwtutil.ValidateToken(token, jwtSecret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "unauthorized",
				"message": "Invalid or expired token",
			})
		}

		// Inject user info into context
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)

		return c.Next()
	}
}
