package controllers

import (
	"backend-hostego/database"
	"backend-hostego/models"

	"github.com/gofiber/fiber/v3"
)

func CreateShop(c fiber.Ctx) error {

	var shop models.Shop

	if err := c.Bind().JSON(&shop); err!= nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": err})
	}
	database.DB.Create(&shop)

	return c.Status(fiber.StatusCreated).JSON(shop)

}

func FetchShopById(c fiber.Ctx) error {

	shop_id := c.Params("shop_id")
	var shop models.Shop

	database.DB.First(&shop, `where shop_id=?`, shop_id)

	return c.Status(fiber.StatusOK).JSON(shop)
}

func FetchShops(c fiber.Ctx) error {

	var shops []models.Shop

	database.DB.Find(&shops)

	return c.Status(fiber.StatusOK).JSON(shops)
}
