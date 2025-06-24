package controllers

import (
	"backend-hostego/database"
	"backend-hostego/models"

	"github.com/gofiber/fiber/v2"
)

func FetchSearchQuery(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}
	var query = database.DB
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate != "" && endDate != "" {
		query = query.Where("search_queries.created_at BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59")
	}
	var searchQuery []models.SearchQuery

	query.Find(&searchQuery).Order("created_at DESC")

	return c.Status(fiber.StatusOK).JSON(searchQuery)
}
