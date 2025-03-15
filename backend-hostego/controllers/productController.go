package controllers

import (
	"backend-hostego/database"
	"backend-hostego/middlewares"
	"backend-hostego/models"

	"strconv"

	"github.com/gofiber/fiber/v3"
)

func CreateNewProduct(c fiber.Ctx) error {
	user_id, err := middlewares.VerifyUserAuthCookie(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "message": "You are not Authenticated !"})
	}
	if user_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User not found"})
	}
	var product models.Product
	if err := c.Bind().JSON(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	database.DB.Create(&product)
	return c.Status(fiber.StatusCreated).JSON(product)

}

func FetchProducts(c fiber.Ctx) error {
	user_id, err := middlewares.VerifyUserAuthCookie(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "message": "You are not Authenticated !"})
	}
	if user_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User not found"})
	}

	var products []models.Product
	dbQuery := database.DB.Preload("Shop")

	isAdmin := c.Query("admin") == "true"

	// Apply stock and availability filters for non-admin users by default
	if !isAdmin {
		dbQuery = dbQuery.Where("stock_quantity > 0 AND availability = '1'")
	}

	searchQuery := c.Query("search")
	tagsQuery := c.Query("tags") // Expecting tags=food or tags=chicken
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	availability := c.Query("availability")
	sort := c.Query("sort", "asc")
	queryLimit := c.Query("limit", "50")
	queryPage := c.Query("page", "1")

	if searchQuery != "" {
		baseQuery := `(product_name ILIKE ? 
			OR description ILIKE ? 
			OR shop_id IN (SELECT shop_id FROM shops WHERE shop_name ILIKE ?)
			OR EXISTS (
				SELECT 1 FROM jsonb_array_elements_text(tags) tag 
				WHERE tag ILIKE ?
			))`

		dbQuery = dbQuery.Where(
			baseQuery,
			"%"+searchQuery+"%",
			"%"+searchQuery+"%",
			"%"+searchQuery+"%",
			"%"+searchQuery+"%",
		)
		database.DB.Create(&models.SearchQuery{
			Query: searchQuery,
		})
	}

	// ✅ Filtering by tags
	if tagsQuery != "" {
		dbQuery = dbQuery.Where("tags @> ?", `["`+tagsQuery+`"]`)
	}

	if minPrice != "" {
		dbQuery = dbQuery.Where("food_price >= ?", minPrice)
	}
	if maxPrice != "" {
		dbQuery = dbQuery.Where("food_price <= ?", maxPrice)
	}
	if availability != "" && !isAdmin {
		dbQuery = dbQuery.Where("availability = ?", availability)
	}

	if sort == "desc" {
		dbQuery = dbQuery.Order("food_price DESC")
	} else {
		dbQuery = dbQuery.Order("food_price ASC")
	}

	limit, err := strconv.Atoi(queryLimit)
	if err != nil || limit < 1 {
		limit = 50
	}

	page, err := strconv.Atoi(queryPage)
	if err != nil {
		page = 1
	}
	offset := (page - 1) * limit
	dbQuery = dbQuery.Offset(offset).Limit(limit)

	if err := dbQuery.Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch products", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

func UpdateProductById(c fiber.Ctx) error {
	user_id, err := middlewares.VerifyUserAuthCookie(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if user_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User not found"})
	}

	product_id := c.Params("id")
	var product models.Product

	// First find existing product
	if err := database.DB.First(&product, "product_id = ?", product_id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	// Get update data and update directly
	if err := c.Bind().JSON(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Save all changes
	if err := database.DB.Save(&product).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product updated successfully",
		"product": product,
	})
}

func FetchProductById(c fiber.Ctx) error {
	user_id, err := middlewares.VerifyUserAuthCookie(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "message": "You are not Authenticated !"})
	}
	if user_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User not found"})
	}
	product_id := c.Params("id")
	var product models.Product

	if err := database.DB.First(&product, "product_id=?", product_id).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"erorr": err})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"product": product})
}
