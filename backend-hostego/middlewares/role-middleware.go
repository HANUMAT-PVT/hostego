package middlewares

import (
	"github.com/gofiber/fiber/v3"
)

// Role-based middleware
func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		userRole := c.Get("X-User-Role") // Assuming role is sent in header

		for _, role := range allowedRoles {
			if userRole == role {
				return c.Next() // Allow access
			}
		}

		// Deny access if role not allowed
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}
}
