package controllers

import (
	"backend-hostego/database"
	"backend-hostego/models"

	"github.com/gofiber/fiber/v3"
)

func FetchSearchQuery(c fiber.Ctx) error {
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}
	var searchQuery []models.SearchQuery

	database.DB.Find(&searchQuery).Order("created_at DESC")

	return c.Status(fiber.StatusOK).JSON(searchQuery)
}
