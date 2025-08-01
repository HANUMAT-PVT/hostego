package controllers

import (
	"backend-hostego/database"
	"backend-hostego/models"

	"github.com/gofiber/fiber/v2"
)

func AddProductInUserCart(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")

	var cartItem models.CartItem
	var product models.Product

	// Bind the incoming cart item data
	if err := c.BodyParser(&cartItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	// Check if product exists
	if err := database.DB.Preload("Shop").First(&product, "product_id = ?", cartItem.ProductId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found!"})
	}
	// check if the shop is closed or not
	if product.Shop.ShopStatus == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Shop is closed!"})
	}

	// Check if item already exists in cart
	var existingItem models.CartItem
	result := database.DB.Where("user_id = ? AND product_id = ?", user_id, cartItem.ProductId).First(&existingItem)

	if result.Error == nil {
		// Item exists, update quantity
		existingItem.Quantity += cartItem.Quantity
		existingItem.SubTotal = float64(existingItem.Quantity) * product.SellingPrice

		if err := database.DB.Save(&existingItem).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Cart item quantity updated!",
			"cart":    existingItem,
		})
	}

	var cartItems []models.CartItem

	//check if the product is from different shop
	database.DB.Preload("ProductItem").Where("user_id = ?", user_id).Find(&cartItems)

	if len(cartItems) > 0 {
		var existingNon4ShopID int
		hasOtherNon4Shop := false

		for _, item := range cartItems {
			shopID := item.ProductItem.ShopId
			if shopID != 4 {
				if existingNon4ShopID == 0 {
					existingNon4ShopID = shopID
				}
				if shopID != product.ShopId {
					hasOtherNon4Shop = true
				}
			}
		}

		if hasOtherNon4Shop {
			for _, item := range cartItems {
				// Delete all items that are not from Shop 4 or the new product's shop
				if item.ProductItem.ShopId != 4 && item.ProductItem.ShopId != product.ShopId {
					database.DB.Where("cart_item_id = ?", item.CartItemId).Delete(&item)
				}
			}
		}
	}

	// Create new cart item if it doesn't exist
	cartItem.UserId = user_id.(int)
	cartItem.SubTotal = float64(cartItem.Quantity) * product.SellingPrice
	cartItem.ActualSubTotal = float64(cartItem.Quantity) * product.FoodPrice

	if err := database.DB.Create(&cartItem).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Product added to cart successfully!",
		"cart":    cartItem,
	})
}

func UpdateProductInUserCart(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")
	cart_item_id := c.Params("id")

	var cartItem models.CartItem
	// First find the existing cart item
	if err := database.DB.Preload("ProductItem").Preload("ProductItem.Shop").Where("cart_item_id = ? AND user_id = ?", cart_item_id, user_id).First(&cartItem).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Cart item not found"})
	}
	// check if the shop is closed or not
	if cartItem.ProductItem.Shop.ShopStatus == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Shop is closed!"})
	}
	// Bind updated data directly to cartItem
	if err := c.BodyParser(&cartItem); err != nil {
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
	cartItem.SubTotal = float64(cartItem.Quantity) * product.SellingPrice
	cartItem.ActualSubTotal = float64(cartItem.Quantity) * product.FoodPrice

	// Save the updates
	if err := database.DB.Save(&cartItem).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Cart updated successfully!",
		"cart_item": cartItem,
	})
}

func FetchUserCart(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")
	var cartItems []models.CartItem
	var orderItems []models.Order
	freeDelivery := false

	if err := database.DB.Where("user_id = ?", user_id).Find(&orderItems).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	if len(orderItems) < 1 {
		freeDelivery = true
	}
	if err := database.DB.Preload("ProductItem.Shop").
		Where("user_id = ?", user_id).
		Order("created_at desc").
		Find(&cartItems).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	// delete the cart item if the shop is closed or not
	for _, cartItem := range cartItems {
		if cartItem.ProductItem.Shop.ShopStatus == 0 {
			database.DB.Where("cart_item_id = ? AND user_id = ?", cartItem.CartItemId, user_id).Delete(&cartItem)
		}
	}
	if err := database.DB.Preload("ProductItem.Shop").
		Where("user_id = ?", user_id).
		Order("created_at desc").
		Find(&cartItems).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	cartValue := CalculateFinalOrderValue(cartItems, freeDelivery)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"cart_items":    cartItems,
		"cart_value":    cartValue,
		"free_delivery": freeDelivery,
	})
}

func DeleteCartItem(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")

	cart_item_id := c.Params("id")
	if err := database.DB.Where("cart_item_id = ? AND user_id = ?", cart_item_id, user_id).Delete(&models.CartItem{}).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Cart item deleted successfully!"})
}
