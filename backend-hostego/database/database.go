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
	err = db.AutoMigrate(&models.Shop{})
	err = db.AutoMigrate(&models.Product{})
	err = db.AutoMigrate(&models.Rating{})

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Connected to Database !")

}
