package controllers

import (
	"backend-hostego/database"
	"backend-hostego/models"

	"github.com/gofiber/fiber/v3"
)

func CreateUserRole(c fiber.Ctx) error {
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var user models.UserRole
	if err := database.DB.Where("user_id = ? AND role_id = ?", user_id, 1).First(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userRole := new(models.UserRole)

	if err := c.Bind().JSON(userRole); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// add check if there is 1 role id coming abort if not
	if userRole.RoleId == 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Only 1 super admin role id is allowed"})
	}
	if user.RoleId == 2 && userRole.RoleId == 2 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User Admin can't assign to User Admin"})
	}
	if err := database.DB.Create(&userRole).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User role created successfully"})
}

func FetchUserRoles(c fiber.Ctx) error {
	user_id := c.Locals("user_id")

	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User not found"})
	}
	var userRoles []models.UserRole
	if err := database.DB.Preload("Role").Where("user_id = ?", user_id).Find(&userRoles).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(userRoles)

}

func DeleteUserRole(c fiber.Ctx) error {
	user_id := c.Locals("user_id")

	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User not found"})
	}
	userRoleId := c.Params("id")
	var userRole models.UserRole
	if err := database.DB.Where("user_role_id = ?", userRoleId).First(&userRole).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if userRole.RoleId == 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Super Admin can't be deleted"})
	}

	if err := database.DB.Where("user_role_id = ?", userRoleId).Delete(&models.UserRole{}).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User role deleted successfully"})
}
