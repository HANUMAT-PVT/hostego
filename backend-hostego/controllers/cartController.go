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

	var cartItem models.CartItem
	var product models.Product

	// Bind the incoming cart item data
	if err := c.Bind().JSON(&cartItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	// Check if product exists
	if err := database.DB.First(&product, "product_id = ?", cartItem.ProductId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found!"})
	}

	// Check if item already exists in cart
	var existingItem models.CartItem
	result := database.DB.Where("user_id = ? AND product_id = ?", user_id, cartItem.ProductId).First(&existingItem)

	if result.Error == nil {
		// Item exists, update quantity
		existingItem.Quantity += cartItem.Quantity
		existingItem.SubTotal = float64(existingItem.Quantity) * product.FoodPrice

		if err := database.DB.Save(&existingItem).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Cart item quantity updated!",
			"cart":    existingItem,
		})
	}

	// Create new cart item if it doesn't exist
	cartItem.UserId = user_id
	cartItem.SubTotal = float64(cartItem.Quantity) * product.FoodPrice

	if err := database.DB.Create(&cartItem).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Product added to cart successfully!",
		"cart":    cartItem,
	})
}

func UpdateProductInUserCart(c fiber.Ctx) error {
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	cart_item_id := c.Params("id")
	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}

	var cartItem models.CartItem
	// First find the existing cart item
	if err := database.DB.Where("cart_item_id = ? AND user_id = ?", cart_item_id, user_id).First(&cartItem).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Cart item not found"})
	}

	// Bind updated data directly to cartItem
	if err := c.Bind().JSON(&cartItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Delete cart item if quantity is 0
	if cartItem.Quantity <= 0 {
		if err := database.DB.Where("cart_item_id = ? AND user_id = ?", cart_item_id, user_id).Delete(&cartItem).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Cart item removed successfully!",
		})
	}

	// Update subtotal if quantity changed
	var product models.Product
	if err := database.DB.First(&product, "product_id = ?", cartItem.ProductId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}
	cartItem.SubTotal = float64(cartItem.Quantity) * product.FoodPrice

	// Save the updates
	if err := database.DB.Save(&cartItem).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Cart updated successfully!",
		"cart_item": cartItem,
	})
}

func FetchUserCart(c fiber.Ctx) error {
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	var cartItems []models.CartItem
	var orderItems []models.Order
	freeDelivery := false
	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}
	if err := database.DB.Where("user_id = ?", user_id).Find(&orderItems).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	if len(orderItems) < 1 {
		freeDelivery = true
	}
	if err := database.DB.Preload("ProductItem.Shop").
		Where("user_id = ?", user_id).
		Order("created_at asc").
		Find(&cartItems).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	cartValue := CalculateFinalOrderValue(cartItems, freeDelivery)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"cart_items": cartItems,
		"cart_value":  cartValue,
		"free_delivery": freeDelivery,
	})
}
