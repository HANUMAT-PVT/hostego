package controllers

import (
	"backend-hostego/internal/app/hostego-service/database"
	"backend-hostego/internal/app/hostego-service/models"

	"github.com/gofiber/fiber/v3"
)

func FetchUserByMobileNumber(c fiber.Ctx) error {
	mobileNumber := c.Params("mobile_number") // Get mobile number from URL params

	var user models.User

	// Query database for user with the given mobile number
	if err := database.DB.Where("mobile_number = ?", mobileNumber).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"user": user})
}

//test commit
