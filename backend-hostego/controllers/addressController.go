package controllers

import (
	"backend-hostego/database"
	"backend-hostego/models"

	"github.com/gofiber/fiber/v2"
)

func CreateNewAddress(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")
	var address models.Address
	var user models.User

	if err := c.BodyParser(&address); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	if err := database.DB.First(&user, " user_id = ?", user_id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found !"})
	}
	address.UserId = user.UserId
	if err := database.DB.Create(&address).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Address added succesfully !"})
}

func UpdateAddress(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")
	var address models.Address
	var user models.User
	address_id := c.Params("id")

	if err := c.BodyParser(&address); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	if err := database.DB.Where("user_id = ?", user_id).Find(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found !"})
	}
	if err := database.DB.Where("address_id = ?", address_id).Save(&address).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Address updated succesfully !"})
}

func DeleteAddress(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")
	var address models.Address
	address_id := c.Params("id")
	var user models.User

	if err := database.DB.First(&user, "where user_id = ?", user_id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found !"})
	}

	if err := database.DB.Where("address_id = ?", address_id).Delete(&address).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Address updated succesfully !"})
}

func FetchUserAddress(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")
	var address []models.Address
	if err := database.DB.Where("user_id = ?", user_id).Order("created_at desc").Find(&address).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	return c.Status(fiber.StatusOK).JSON(address)

}
