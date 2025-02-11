package models

import "time"

// // Food Item Schema
// {
//   "foodId": "12345",  // Unique identifier for the food item
//   "   name": "Paneer Butter Masala",
//   "category": {
//     "isVeg": true,
//     "isCooked": true
//   },
//   "price": 250.00,
//   "currency": "INR",
//   "availability": {
//     "inStock": true,
//     "quantity": 10
//   },
//   "shop": {
//     "shopId": "S6789",
//     "shopName": "Tandoori Delights"
//   },
//   "imageUrl": "https://example.com/paneer.jpg",
//   "description": "A creamy and rich paneer dish cooked in butter with tomato-based gravy.",
//   "tags": ["North Indian", "Spicy", "Creamy"],  // Optional tags for better search filtering
//   "ratings": {
//     "averageRating": 4.5,
//     "totalReviews": 120
//   },
//   "discount": {
//     "isAvailable": true,
//     "percentage": 10
//   },
//   "preparationTime": "30 min",
//   "addedOn": "2025-02-09T12:00:00Z"
// }

type FoodCategory struct {
	IsVeg    int `gorm:"type:int;not null;" json:"is_veg"`
	IsCooked int `gorm:"type:int:not null;" json:"is_cooked"`
}

type Discount struct {
	IsAvailable int `gorm:"type:int;not null;" json:"is_available"`
	Percentage  int `gorm:"type:int:not null;" json:"percentage"`
}

type Product struct {
	ProductId       string       `gorm:"type:uuid;not null;unique;default:gen_random_uuid()" json:"product_id"`
	ProductName     string       `gorm:"type:varchaar(255);not null" json:"product_name"`
	FoodCategory    FoodCategory `gorm:"embedded" json:"food_category"`
	FoodPrice       float64      `gorm:"type:float64;not null" json:"food_price"`
	Availability    int          `gorm:"type:int;not null;default:1" json:"availability"`
	ProductImg      string       `gorm:"type:varchaar(255);not null" json:"product_img_url"`
	Description     string       `gorm:"type:varchaar(255);not null" json:"description"`
	Discount        Discount     `gorm:"embedded" json:"discount"`
	CreatedAt       time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	PreparationTime string       `gorm:"type:varchaar(255);not null" json:"preparation_time"`
	ShopId          string       `gorm:"type:string;not null; json:"shop_id"`
}
