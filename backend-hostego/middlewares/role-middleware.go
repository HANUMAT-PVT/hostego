package middlewares

import (
	"backend-hostego/database"
	"backend-hostego/models"
	"github.com/gofiber/fiber/v3"
)

// Role-based middleware
func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {

		user_id := c.Locals("user_id")
		
		db := database.DB
		var userRoles []models.UserRole
		err := db.Preload("Role").Where("user_id = ?", user_id).Find(&userRoles).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch user roles",
			})
		}
		

		for _, userRole := range userRoles {
			for _, allowedRole := range allowedRoles {
				if userRole.Role.RoleName == allowedRole {
					return c.Next()
				}
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}

}
