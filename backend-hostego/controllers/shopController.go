package controllers

import (
	"backend-hostego/database"
	"backend-hostego/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CreateShop(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}
	var shop models.Shop
	shop.OwnerId = user_id.(int)

	if err := c.BodyParser(&shop); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": err})
	}
	database.DB.Create(&shop)

	return c.Status(fiber.StatusCreated).JSON(shop)

}

func FetchShopById(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	offset := (page - 1) * limit

	var products []models.Product
	shop_id := c.Params("id")
	var shop models.Shop

	database.DB.Where("shop_id = ?", shop_id).First(&shop)

	database.DB.Where("shop_id = ?", shop_id).Limit(limit).Offset(offset).Find(&products)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"shop": shop, "products": products})
}

func FetchShops(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	offset := (page - 1) * limit

	var shops []models.Shop

	database.DB.Limit(limit).Offset(offset).Find(&shops)

	return c.Status(fiber.StatusOK).JSON(shops)
}

func UpdateShopById(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}
	var shop models.Shop
	shop_id := c.Params("id")

	if err := database.DB.Where("shop_id = ?", shop_id).First(&shop).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Shop not found!"})

	}
	if err := c.BodyParser(&shop); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	if err := database.DB.Save(&shop).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Shop updated succesfully", "shop": shop})

}

func FetchShopByOwnerId(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}
	var shop models.Shop

	database.DB.Where("owner_id = ?", user_id).First(&shop)

	return c.Status(fiber.StatusOK).JSON(shop)
}
