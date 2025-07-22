package controllers

import (
	"backend-hostego/database"
	"backend-hostego/models"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

func SaveRatingAndUpdateStats(c *fiber.Ctx) error {
	rating := new(models.Rating)
	if err := c.BodyParser(rating); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	user_id := c.Locals("user_id").(int)
	rating.UserID = user_id
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// Step 1: Save the new rating
		if err := tx.Create(rating).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		// Step 2: Update Product stats
		var productStats struct {
			Average float64
			Count   int64
		}
		if err := tx.
			Model(&models.Rating{}).
			Select("AVG(rating) AS average, COUNT(*) AS count").
			Where("product_id = ?", rating.ProductID).
			Scan(&productStats).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := tx.
			Model(&models.Product{}).
			Where("product_id = ?", rating.ProductID).
			Updates(map[string]interface{}{
				"average_rating": productStats.Average,
				"total_ratings":  productStats.Count,
			}).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		// Step 3: Get ShopID from the product
		var product models.Product
		if err := tx.Select("shop_id").Where("product_id = ?", rating.ProductID).First(&product).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		// Step 4: Update Shop stats
		var shopStats struct {
			Average float64
			Count   int64
		}
		if err := tx.
			Table("ratings").
			Select("AVG(ratings.rating) AS average, COUNT(*) AS count").
			Joins("JOIN products ON ratings.product_id = products.product_id").
			Where("products.shop_id = ?", product.ShopId).
			Scan(&shopStats).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := tx.
			Model(&models.Shop{}).
			Where("shop_id = ?", product.ShopId).
			Updates(map[string]interface{}{
				"average_rating": shopStats.Average,
				"total_ratings":  shopStats.Count,
			}).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Rating saved successfully"})
	})
}

func FetchRatingsForProduct(c *fiber.Ctx) error {
	var average float64
	productID := c.Params("product_id")

	err := database.DB.Model(&models.Rating{}).
		Where("product_id = ?", productID).
		Select("AVG(rating)").
		Scan(&average).Error

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"average": average})
}

func FetchRatingsForOrder(c *fiber.Ctx) error {
	orderId := c.Params("order_id")

	var ratings []models.Rating

	if err := database.DB.Where("order_id = ? ", orderId).Find(&ratings).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"ratings": ratings})
}

// func FetchRatingsFor
