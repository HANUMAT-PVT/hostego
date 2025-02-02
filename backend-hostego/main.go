package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/postgres"
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

// Connect to PostgreSQL
func connectDatabase() {
	dsn := "host=localhost user=hostego_user_dev password=hostego_hanumat dbname=hostego_db_dev port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	db.AutoMigrate(&User{}, &Address{})

	
}


func CreateUser(c fiber.Ctx) error {
	var user User
	if err := c.JSON(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Insert user into database
	db.Create(&user)
	return c.Status(201).JSON(user)
}

func GetUsers(c fiber.Ctx) error {
	var users []User
	db.Find(&users)
	return c.Status(200).JSON(users)
}

// API Handlers
func main() {
	connectDatabase()

	app := fiber.New()

	// Fetch all users
	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message":"Welcome to the server"})
	})

	app.Post("/users", CreateUser)

	app.Get("/users", GetUsers)
	
	

	log.Fatal(app.Listen(":3000"))
}
