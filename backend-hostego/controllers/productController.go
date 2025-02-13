package controllers

import (
	"backend-hostego/database"
	"backend-hostego/models"

	"github.com/gofiber/fiber/v3"
)

func CreateNewProduct(c fiber.Ctx) error {
	var product models.Product

	if err := c.Bind().JSON(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	database.DB.Create(&product)
	return c.Status(fiber.StatusCreated).JSON(product)

}

func FetchProducts(c fiber.Ctx) error {
	var products []models.Product


	err := database.DB.Preload("Shop").Find(&products).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch products"})
	}

	return c.Status(fiber.StatusOK).JSON(products)
}
