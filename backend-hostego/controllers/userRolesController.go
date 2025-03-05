package controllers

import (
	"backend-hostego/database"
	"backend-hostego/middlewares"
	"backend-hostego/models"

	"github.com/gofiber/fiber/v3"
)

func CreateUserRole(c fiber.Ctx) error {
	user_id, err := middlewares.VerifyUserAuthCookie(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if user_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User not found"})
	}
	userRole := new(models.UserRole)
	if err := c.Bind().JSON(userRole); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := database.DB.Create(&userRole).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User role created successfully"})
}

func FetchUserRoles(c fiber.Ctx) error {
	user_id, err := middlewares.VerifyUserAuthCookie(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if user_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User not found"})
	}
	var userRoles []models.UserRole
	if err := database.DB.Preload("Role").Where("user_id = ?", user_id).Find(&userRoles).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(userRoles)

}

func DeleteUserRole(c fiber.Ctx) error {
	user_id, err := middlewares.VerifyUserAuthCookie(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if user_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User not found"})
	}
	userRoleId := c.Params("id")
	if err := database.DB.Where("user_role_id = ?", userRoleId).Delete(&models.UserRole{}).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User role deleted successfully"})
}
