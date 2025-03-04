package database

import (
	"log"

	"gorm.io/driver/postgres"

	"backend-hostego/models"

	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDataBase() {

	dsn := "host=localhost user=hostego_user_dev password=hostego_hanumat dbname=hostego_db_dev port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect database", err)
	}
	DB = db
	err = db.AutoMigrate(&models.User{})
	err = db.AutoMigrate(&models.Address{})
	err = db.AutoMigrate(&models.Wallet{})
	err = db.AutoMigrate(&models.WalletTransaction{})
	err = db.AutoMigrate(&models.PaymentTransaction{})
	err = db.AutoMigrate(&models.Shop{})
	err = db.AutoMigrate(&models.Product{})
	err = db.AutoMigrate(&models.CartItem{})
	err = db.AutoMigrate(&models.Rating{})
	err = db.AutoMigrate(&models.Order{})
	err = db.AutoMigrate(&models.Role{})
	err = db.AutoMigrate(&models.UserRole{})
	err = db.AutoMigrate(&models.DeliveryPartner{})

	// db.Create(&models.Role{
	// 	RoleName: "super_admin",
	// })
	// db.Create(&models.Role{
	// 	RoleName: "admin",
	// })
	// db.Create(&models.Role{
	// 	RoleName: "user",
	// })
	// db.Create(&models.Role{
	// 	RoleName: "delivery_partner",
	// })

	// db.Create(&models.Role{
	// 	RoleName: "admin",
	// })
	// db.Create(&models.Role{
	// 	RoleName: "payments_manager",
	// })
	// db.Create(&models.Role{
	// 	RoleName: "order_assign_manager",
	// })
	// db.Create(&models.Role{
	// 	RoleName: "delivery_partner_manager",
	// })
	// db.Create(&models.Role{
	// 	RoleName: "order_manager",
	// })
	// db.Create(&models.Role{
	// 	RoleName: "delivery_partner",
	// })
	// db.Create(&models.Role{
	// 	RoleName: "user",
	// })
	// db.Create(&models.Role{
	// 	RoleName: "customer_support",
	// })

	// db.Create(&models.UserRole{
	// 	UserId: "c5c7946a-d056-4319-bfb8-e930c91052bc",
	// 	RoleId: 1,
	// })

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Connected to Database !")

}

// -- Switch to your database
