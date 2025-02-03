package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// Global DB variable
var db *gorm.DB

// User Model
type User struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Email       string    `json:"email" gorm:"unique"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phone" gorm:"unique;not null"`
	// Addresses   []Address `gorm:"foreignKey:UserID"`
}

// Address Model
type Address struct {
	ID           uint   `gorm:"primaryKey"`
	UserID       uint   `gorm:"index"`
	Hostel       string `json:"hostel"`
	RoomNumber   string `json:"room_number"`
	FloorNumber  string `json:"floor_number"`
	OtherDetails string `json:"other_details"`
}


func main() {


	app := fiber.New()

	// Fetch all users
	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message":"Welcome to the server"})
	})

	
	

	log.Fatal(app.Listen(":8080"))
}
