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
	// hostego_user_dev for prod
	dsn := "host=localhost user=hostego_user_dev password=hostego_hanumat dbname=hostego_db_dev port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// DisableForeignKeyConstraintWhenMigrating: true,
		Logger: logger.Default.LogMode(logger.Error), // Only log errors and panics
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
	sqlDB.SetMaxOpenConns(100)                 // Increase max connections
	sqlDB.SetMaxIdleConns(50)                  // Increase idle connections
	sqlDB.SetConnMaxLifetime(15 * time.Minute) // Shorter lifetime to prevent stale connections
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)  // Close idle connections faster

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

	if err != nil {
		log.Fatal("🚨 CRITICAL: Failed to migrate database:", err)
	}

	log.Println("✅ Database migration completed successfully")
	log.Println("✅ Connected to Database with enhanced monitoring!")

	// Start database health monitoring
	StartDatabaseMonitoring()

}

// -- Switch to your database
