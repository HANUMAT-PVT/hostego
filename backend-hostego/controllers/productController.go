package controllers

import (
	"backend-hostego/database"
	"backend-hostego/models"
	"strconv"

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

	dbQuery := database.DB

	searchQuery := c.Query("search")
	// tagsQuery := c.Query("tags")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	availability := c.Query("availability")
	sort := c.Query("sort", "asc")
	queryLimit := c.Query("limit", "10") // Default 10 products per page
	queryPage := c.Query("page", "1")    //Default page is 1

	if searchQuery != "" {
		dbQuery = dbQuery.Where("product_name ILIKE  ? OR description ILIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}
	// if tagsQuery != "" {

	// 	dbQuery = dbQuery.Where("tags @> ?", datatypes.JSON(tagsQuery))

	// }
	if minPrice != "" {
		dbQuery = dbQuery.Where("food_price >= ?", minPrice)

	}
	if maxPrice != "" {
		dbQuery = dbQuery.Where("food_price<= ?", maxPrice)
	}
	if availability != "" {
		dbQuery = dbQuery.Where("availability=?", availability)
	}

	if sort == "desc" {
		dbQuery = dbQuery.Order("food_price DESC")
	} else {
		dbQuery = dbQuery.Order("food_price ASC")
	}

	limit, err := strconv.Atoi(queryLimit)
	if err != nil || limit < 1 {
		limit = 10
	}

	// ✅ Get "page" query parameter (default to 1)

	page, err := strconv.Atoi(queryPage)
	if err != nil {
		page = 1
	}
	offset := (page - 1) * limit
	dbQuery = dbQuery.Offset(offset).Limit(limit)

	// **Preload Shop details**
	errDb := dbQuery.Preload("Shop").Find(&products).Error

	if errDb != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch products"})
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

func UpdateProductById(c fiber.Ctx) error {
	var product models.Product
	product_id := c.Params("id")

	err := database.DB.First(&product, "product_id ?", product_id).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Product not found !"})
	}
	if err := c.Bind().JSON(&product).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := database.DB.Save(&product).Error; err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Product updated succesfully", "product": product})

}
