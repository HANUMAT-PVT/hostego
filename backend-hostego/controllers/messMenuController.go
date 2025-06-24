package controllers

import (
	"backend-hostego/database"
	"backend-hostego/models"

	"github.com/gofiber/fiber/v2"
)

func UpdateMessMenu(c *fiber.Ctx) error {
	var id = c.Params("id")
	user_id := c.Locals("user_id")

	if user_id == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var messMenu models.MessMenu
	if err := c.BodyParser(&messMenu); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	database.DB.Where("id = ?", id).Save(&messMenu)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Mess menu updated successfully",
	})
}

func FetchMessMenu(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var messMenuItems []models.MessMenu
	database.DB.Find(&messMenuItems)

	return c.Status(fiber.StatusOK).JSON(messMenuItems)
}

func CreateMessMenuDate(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var messMenu models.MessMenu
	if err := c.BodyParser(&messMenu); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	database.DB.Create(&messMenu)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Mess menu created successfully",
	})
}
