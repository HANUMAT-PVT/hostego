package controllers

// func FetchUserCart(c fiber.Ctx) error {
// 	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
// 	var cartItems []models.CartItem

// 	if middleErr != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
// 	}
// 	if err := database.DB.Preload("ProductItem.Shop").Where("user_id=?", user_id).Find(&cartItems).Error; err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
// 	}
// 	cartValue := CalculateFinalOrderValue(cartItems)

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{"cart_items": cartItems, "cart_value": cartValue})
// }
