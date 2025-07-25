package controllers

import (
	"backend-hostego/database"
	"backend-hostego/models"
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
)

func CreateNewProduct(c *fiber.Ctx) error {

	userID := c.Locals("user_id")
	if userID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var product models.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	var tags []string
	if err := json.Unmarshal(product.Tags, &tags); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid tags format"})
	}

	// Check if tags contain "food"
	hasFoodTag := false
	for _, tag := range tags {
		if tag == "food" {
			hasFoodTag = true
			break
		}
	}

	// Set selling price based on food_price if "food" tag is present
	if hasFoodTag {
		switch {
		case product.FoodPrice > 25 && product.FoodPrice < 50:
			product.SellingPrice = product.FoodPrice + 5
		case product.FoodPrice >= 50 && product.FoodPrice <= 100:
			product.SellingPrice = product.FoodPrice + 10
		case product.FoodPrice >= 100 && product.FoodPrice <= 150:
			product.SellingPrice = product.FoodPrice + 15
		case product.FoodPrice >= 150:
			product.SellingPrice = product.FoodPrice + 20
		default:
			product.SellingPrice = product.FoodPrice
		}
	} else {
		// if no "food" tag, just set selling_price same as food_price (optional)
		product.SellingPrice = product.FoodPrice
	}

	// Save product
	if err := database.DB.Create(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create product"})
	}

	return c.Status(fiber.StatusCreated).JSON(product)

}

func FetchProducts(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User not found"})
	}

	var products []models.Product
	dbQuery := database.DB.Preload("Shop")

	isAdmin := c.Query("admin") == "true"

	// Apply stock and availability filters for non-admin users by default
	if !isAdmin {
		dbQuery = dbQuery.Where("stock_quantity > 0").Where("availability = ?", 1)
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
			Query:  searchQuery,
			UserId: user_id.(int),
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
func UpdateProductById(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	productID := c.Params("id")
	var product models.Product

	// Find product
	if err := database.DB.First(&product, "product_id = ?", productID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	// Parse incoming data
	var updateData map[string]interface{}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Handle tags
	var existingTags, incomingTags []string
	hasFoodTag := false

	if product.Tags != nil {
		_ = json.Unmarshal(product.Tags, &existingTags)
	}
	if rawTags, ok := updateData["tags"]; ok {
		tagsJSON, err := json.Marshal(rawTags)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid tags format"})
		}
		_ = json.Unmarshal(tagsJSON, &incomingTags)
		updateData["tags"] = datatypes.JSON(tagsJSON)
	}

	// Check for "food" tag
	for _, tag := range append(existingTags, incomingTags...) {
		if tag == "food" {
			hasFoodTag = true
			break
		}
	}

	// Handle selling_price if food_price is provided
	if hasFoodTag {
		if fpRaw, ok := updateData["food_price"]; ok {
			// Ensure float64
			foodPrice, ok := fpRaw.(float64)
			if !ok {
				// Try converting from int
				if fpi, ok := fpRaw.(int); ok {
					foodPrice = float64(fpi)
				}
			}
			var sellingPrice float64
			switch {
			case foodPrice > 25 && foodPrice < 50:
				sellingPrice = foodPrice + 5
			case foodPrice >= 50 && foodPrice <= 100:
				sellingPrice = foodPrice + 10
			case foodPrice > 100 && foodPrice <= 150:
				sellingPrice = foodPrice + 15
			case foodPrice > 150:
				sellingPrice = foodPrice + 20
			default:
				sellingPrice = foodPrice
			}
			updateData["selling_price"] = sellingPrice
		}
	}

	// Convert numeric fields that must be integers
	numericIntFields := []string{"availability", "shop_id", "stock_quantity", "total_ratings"}
	for _, key := range numericIntFields {
		if val, ok := updateData[key]; ok {
			if f, ok := val.(float64); ok {
				updateData[key] = int(f)
			}
		}
	}

	// Now update the product
	if err := database.DB.Model(&product).Updates(updateData).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update product"})
	}

	// Reload the updated product

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product updated successfully",
		"product": product,
	})
}

func FetchProductById(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User not found"})
	}
	product_id := c.Params("id")
	var product models.Product

	if err := database.DB.First(&product, "product_id=?", product_id).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"erorr": err})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"product": product})
}

func FetchProductsByShopId(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}
	limit := c.Query("limit", "50")
	page := c.Query("page", "1")
	shop_id := c.Params("shop_id")
	search := c.Query("search")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 50
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	offset := (pageInt - 1) * limitInt
	var products []models.Product
	if err := database.DB.Where("shop_id = ?", shop_id).Where("product_name ILIKE ?", "%"+search+"%").Offset(offset).Limit(limitInt).Find(&products).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No product found !"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"products": products})
}
