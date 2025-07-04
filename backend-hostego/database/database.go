package database

import (
	"log"
	"time"

	"gorm.io/driver/postgres"

	"backend-hostego/models"

	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDataBase() {
	// postgres for dev
	dsn := "host=localhost user=postgres password=hostego_hanumat dbname=hostego_db_dev port=5432 sslmode=disable"
	// hostego_user_dev for prod
	// dsn := "host=localhost user=hostego_user_dev password=hostego_hanumat dbname=hostego_db_dev port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect database", err)

	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get underlying sql.DB:", err)
	}
	sqlDB.SetMaxOpenConns(10)               // max total DB connections
	sqlDB.SetMaxIdleConns(5)                // max idle connections
	sqlDB.SetConnMaxLifetime(1 * time.Hour) // max lifetime of a connection

	DB = db
	DB = DB.Debug()

	err = db.AutoMigrate(&models.User{}, &models.Address{}, &models.Wallet{}, &models.WalletTransaction{},
		&models.PaymentTransaction{}, &models.Shop{}, &models.Product{}, &models.CartItem{}, &models.Rating{},
		&models.Order{}, &models.Role{}, &models.UserRole{}, &models.DeliveryPartner{}, &models.SearchQuery{}, &models.OrderItem{}, &models.DeliveryPartnerWallet{}, &models.DeliveryPartnerWalletTransaction{}, &models.MessMenu{}, &models.ExtrCharge{}, &models.ProductCategory{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Connected to Database !")

}

// -- Switch to your database
