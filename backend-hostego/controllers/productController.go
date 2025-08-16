package controllers

import (
	"backend-hostego/database"
	"backend-hostego/models"
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
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
	// here how to get the products from  shops where the shop is verified
	var dbQuery *gorm.DB
	isAdmin := c.Query("admin") == "true"

	if isAdmin {
		// get all the shops no need to check the is_shop_verified
		dbQuery = database.DB.Preload("Shop")
	} else {
		dbQuery = database.DB.Preload("Shop").Where("shop_id IN (SELECT shop_id FROM shops WHERE is_shop_verified = ?)", true)
	}

	// Apply stock and availability filters for non-admin users by default
	if !isAdmin {
		dbQuery = dbQuery.Where("availability = ?", 1)
	}

	searchQuery := c.Query("search")
	tagsQuery := c.Query("tags") // Expecting tags=food or tags=chicken
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	availability := c.Query("availability")
	sort := c.Query("sort", "asc")
	queryLimit := c.Query("limit", "50")
	queryPage := c.Query("page", "1")
	var totalProducts int64

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
	database.DB.Model(&models.Product{}).Count(&totalProducts)

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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"products": products, "total": totalProducts})
}
func calculateSellingPrice(foodPrice float64) float64 {
	switch {
	case foodPrice > 25 && foodPrice < 50:
		return foodPrice + 5
	case foodPrice >= 50 && foodPrice <= 100:
		return foodPrice + 10
	case foodPrice > 100 && foodPrice <= 150:
		return foodPrice + 15
	case foodPrice > 150:
		return foodPrice + 20
	default:
		return foodPrice
	}
}

func UpdateProductById(c *fiber.Ctx) error {
	productID := c.Params("id")

	// 1. Define request struct inline
	type DiscountDTO struct {
		IsAvailable *int     `json:"is_available"`
		Percentage  *float64 `json:"percentage"`
	}

	type UpdateProductRequest struct {
		ProductName   *string              `json:"product_name"`
		FoodPrice     *float64             `json:"food_price"`
		FoodCategory  *models.FoodCategory `json:"food_category"`
		Tags          *datatypes.JSON      `json:"tags"`
		Description   *string              `json:"description"`
		StockQuantity *int                 `json:"stock_quantity"`
		AverageRating *float64             `json:"average_rating"`
		IsVeg         *int                 `json:"is_veg"`
		IsRecommended *int                 `json:"is_recommended"`
		Availability  *int                 `json:"availability"`
		ShopId        *int                 `json:"shop_id"`
		Discount      *DiscountDTO         `json:"discount"`
		ProductImgUrl *string              `json:"product_img_url"`
	}

	// 2. Parse JSON body
	var req UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON body")
	}

	// 3. Load product
	var product models.Product
	if err := database.DB.First(&product, "product_id = ?", productID).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Product not found")
	}

	// 4. Update only the provided fields
	if req.ProductName != nil {
		product.ProductName = *req.ProductName
	}
	if req.FoodPrice != nil {
		product.FoodPrice = *req.FoodPrice
		product.SellingPrice = calculateSellingPrice(*req.FoodPrice)
	}
	if req.FoodCategory != nil {
		product.FoodCategory = *req.FoodCategory
	}
	if req.Tags != nil {
		product.Tags = *req.Tags
	}

	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.StockQuantity != nil {
		product.StockQuantity = *req.StockQuantity
	}
	if req.AverageRating != nil {
		product.AverageRating = *req.AverageRating
	}
	if req.IsVeg != nil {
		product.FoodCategory.IsVeg = *req.IsVeg
	}
	if req.IsRecommended != nil {
		product.FoodCategory.IsCooked = *req.IsRecommended
	}
	if req.Availability != nil {
		product.Availability = *req.Availability
	}
	if req.ShopId != nil {
		product.ShopId = *req.ShopId
	}
	if req.Discount != nil {
		if req.Discount.IsAvailable != nil {
			product.FoodCategory.IsCooked = *req.Discount.IsAvailable
		}
		if req.Discount.Percentage != nil {
			product.FoodCategory.IsCooked = int(*req.Discount.Percentage)
		}
	}
	if req.ProductImgUrl != nil {
		product.ProductImgUrl = *req.ProductImgUrl
	}
	// 5. Save updated product
	if err := database.DB.Save(&product).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update product")
	}

	// 6. Return response
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
