package controllers

import (
	"backend-hostego/database"
	"backend-hostego/models"

	"github.com/gofiber/fiber/v3"
)

func FetchSearchQuery(c fiber.Ctx) error {

	var searchQuery []models.SearchQuery

	database.DB.Find(&searchQuery).Order("created_at DESC")

	return c.Status(fiber.StatusOK).JSON(searchQuery)
}
