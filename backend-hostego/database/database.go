package database

import (
	"log"
	"time"

	"gorm.io/driver/postgres"

	"backend-hostego/models"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDataBase() {
	//postgres for dev
	// dsn := "host=localhost user=postgres password=hostego_hanumat dbname=hostego_db_dev port=5432 sslmode=disable"
	// hostego_user_dev for prodai
	dsn := "host=localhost user=hostego_user_dev password=hostego_hanumat dbname=hostego_db_dev port=5432 sslmode=disable"
	// dsn := fmt.Sprintf(
	// 	"host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
	// 	os.Getenv("DB_HOST"),
	// 	os.Getenv("DB_USER"),
	// 	os.Getenv("DB_PASSWORD"),
	// 	os.Getenv("DB_NAME"),
	// 	os.Getenv("DB_PORT"),
	// )
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{

		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Error), // Only log errors and panics
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt: true, // Prepare statements for better performance and safety
	})

	if err != nil {
		log.Fatal("🚨 CRITICAL: Failed to connect database", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("🚨 CRITICAL: Failed to get underlying sql.DB:", err)
	}

	// Enhanced connection pool settings for stability
	sqlDB.SetMaxOpenConns(50)                  // Reduced from 100 to prevent connection exhaustion
	sqlDB.SetMaxIdleConns(25)                  // Reduced from 50 to maintain reasonable idle pool
	sqlDB.SetConnMaxLifetime(10 * time.Minute) // Reduced from 15 to 10 minutes
	sqlDB.SetConnMaxIdleTime(3 * time.Minute)
	// Reduced from 5 to 3 minutes for faster cleanup

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("🚨 CRITICAL: Database ping failed:", err)
	}

	log.Println("✅ Database connection established and tested successfully")

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
	err = db.AutoMigrate(&models.Notification{})
	err = db.AutoMigrate(&models.RestaurantPayout{})

	if err != nil {
		log.Fatal("🚨 CRITICAL: Failed to migrate database:", err)
	}

	log.Println("✅ Database migration completed successfully")
	log.Println("✅ Connected to Database with enhanced monitoring!")

	// Start database health monitoring
	StartDatabaseMonitoring()

}

// -- Switch to your database
