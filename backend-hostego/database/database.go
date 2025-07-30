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
	//postgres for dev
	// dsn := "host=localhost user=postgres password=hostego_hanumat dbname=hostego_db_dev port=5432 sslmode=disable"
	// hostego_user_dev for prod
	dsn := "host=localhost user=postgres password=hostego_hanumat dbname=hostego_db_dev port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// DisableForeignKeyConstraintWhenMigrating: true,
	})

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
	// db = db.Debug()
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
	err = db.AutoMigrate(&models.SearchQuery{})
	err = db.AutoMigrate(&models.OrderItem{})
	err = db.AutoMigrate(&models.DeliveryPartnerWallet{})
	err = db.AutoMigrate(&models.DeliveryPartnerWalletTransaction{})
	err = db.AutoMigrate(&models.MessMenu{})
	err = db.AutoMigrate(&models.ExtrCharge{})
	err = db.AutoMigrate(&models.ProductCategory{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Connected to Database !")

}

// -- Switch to your database
