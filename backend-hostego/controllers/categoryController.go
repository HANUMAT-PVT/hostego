package controllers

import (
	"backend-hostego/database"
	"backend-hostego/models"

	"github.com/gofiber/fiber/v2"
)

func CreateNewCategory(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")
	var productCategory models.ProductCategory
	var user models.User

	if err := c.BodyParser(&productCategory); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	if err := database.DB.First(&user, " user_id = ?", user_id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found !"})
	}

	if err := database.DB.Create(&productCategory).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Product category added succesfully !"})
}

func FetchCategories(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func FetchCategoriesByShopId(c *fiber.Ctx) error {
	shop_id := c.Params("shop_id")
	var productCategories []models.ProductCategory

	if err := database.DB.Where("shop_id = ?", shop_id).Find(&productCategories).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No product categories found !"})
	}
	// fmt.Println(productCategories)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": productCategories})
}
