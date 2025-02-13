package controllers

import (
	"backend-hostego/database"
	"backend-hostego/middlewares"
	"fmt"

	"backend-hostego/models"

	"github.com/gofiber/fiber/v3"
)

func CreateUser(c fiber.Ctx) error {
	var user models.User
	if err := c.Bind().JSON(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	fmt.Printf("Parsed User: %+v\n", user)
	database.DB.Create(&user)
	return c.Status(201).JSON(user)
}

func GetUsers(c fiber.Ctx) error {
	var users []models.User
	database.DB.Find(&users)
	return c.Status(200).JSON(users)

}

func GetUserById(c fiber.Ctx) error {

	user_id, err := middlewares.VerifyUserAuthCookie(c)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	var user models.User

	if err := database.DB.First(&user, "user_id = ?", user_id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}
