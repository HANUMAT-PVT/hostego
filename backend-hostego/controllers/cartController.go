package controllers

import (
	"backend-hostego/database"
	"backend-hostego/middlewares"
	"backend-hostego/models"

	"github.com/gofiber/fiber/v3"
)

func AddProductInUserCart(c fiber.Ctx) error {
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}
	var user models.User
	var cartItem models.CartItem
	var product models.Product

	if err := c.Bind().JSON(&cartItem).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	if err := database.DB.First(&user, "where user_id = ?", user_id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found !"})
	}
	if err := database.DB.First(&product, "where product_id = ?", cartItem.ProductId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found !"})
	}

	cartItem.UserId = user_id
	cartItem.SubTotal = cartItem.Quantity * product.FoodPrice

	if err := database.DB.Create(&cartItem).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Product added in cart successfully !", "cart": cartItem})
}

func UpdateProductInUserCart(c fiber.Ctx) error {
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	cart_item_id := c.Params("id")
	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}
	var user models.User
	var cartItem models.CartItem

	if err := c.Bind().JSON(&cartItem).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	if err := database.DB.First(&user, "where user_id = ?", user_id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found !"})
	}

	if err := database.DB.Where("where cart_item_id = ?", cart_item_id).Save(&cartItem).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Cart updated successfully !"})
}

func FetchUserCart(c fiber.Ctx) error {
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	var cartItems []models.CartItem

	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}
	if err := database.DB.Where("user_id=?", user_id).Find(&cartItems).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	cartValue := CalculateFinalOrderValue(cartItems)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"cart_items": cartItems, "cart_value": cartValue})
}
