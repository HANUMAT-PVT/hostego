package controllers

import (
	"backend-hostego/database"
	"backend-hostego/middlewares"
	"backend-hostego/models"

	"github.com/gofiber/fiber/v3"
)

func FetchUserRoles(c fiber.Ctx) error {
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": middleErr.Error()})
	}
	db := database.DB
	var userRoles []models.UserRole
	if err := db.Preload("User").Preload("Role").Where("user_id = ?", user_id).Find(&userRoles).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(userRoles)
}

func CreateNewRole(c fiber.Ctx) error {
	db := database.DB
	var role models.Role

	if err := c.Bind().JSON(&role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	db.Create(&role)
	return c.JSON(role)
}
