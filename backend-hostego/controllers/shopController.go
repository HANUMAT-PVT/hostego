package controllers

import (
	"backend-hostego/database"
	"backend-hostego/models"

	"github.com/gofiber/fiber/v3"
)

func CreateShop(c fiber.Ctx) error {
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}
	var shop models.Shop

	if err := c.Bind().JSON(&shop); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": err})
	}
	database.DB.Create(&shop)

	return c.Status(fiber.StatusCreated).JSON(shop)

}

func FetchShopById(c fiber.Ctx) error {
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}
	var products []models.Product
	shop_id := c.Params("id")
	var shop models.Shop

	database.DB.Where("shop_id = ?", shop_id).First(&shop)

	database.DB.Where("shop_id = ?", shop_id).Find(&products)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"shop": shop, "products": products})
}

func FetchShops(c fiber.Ctx) error {
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}
	var shops []models.Shop

	database.DB.Find(&shops)

	return c.Status(fiber.StatusOK).JSON(shops)
}

func UpdateShopById(c fiber.Ctx) error {
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}
	var shop models.Shop
	shop_id := c.Params("id")

	if err := database.DB.Where("shop_id = ?", shop_id).First(&shop).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Shop not found!"})

	}
	if err := c.Bind().JSON(&shop); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	if err := database.DB.Save(&shop).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Shop updated succesfully", "shop": shop})

}
