package database

import (
	"context"
	"log"
	"runtime/debug"
	"time"

	"gorm.io/gorm"
)

// DatabaseHealthCheck performs comprehensive database health monitoring
func DatabaseHealthCheck() error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("🚨 CRITICAL: Database health check panic: %v", r)
			log.Printf("Stack trace: %s", debug.Stack())
		}
	}()

	if DB == nil {
		log.Printf("🚨 CRITICAL: Database connection is nil")
		return gorm.ErrInvalidDB
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("🚨 CRITICAL: Failed to get underlying sql.DB: %v", err)
		return err
	}

	// Test basic connectivity
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		log.Printf("🚨 CRITICAL: Database ping failed: %v", err)
		return err
	}

	// Check connection pool stats
	stats := sqlDB.Stats()
	log.Printf("📊 DB Pool Stats - Open: %d, InUse: %d, Idle: %d",
		stats.OpenConnections, stats.InUse, stats.Idle)

	if stats.OpenConnections > 20 {
		log.Printf("⚠️  High number of open connections: %d", stats.OpenConnections)
	}

	return nil
}

// SafeTransaction wraps database transactions with enhanced error handling and recovery
func SafeTransaction(operation func(tx *gorm.DB) error, context string) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("🚨 CRITICAL: Transaction panic in %s: %v", context, r)
			log.Printf("Stack trace: %s", debug.Stack())
		}
	}()

	// Health check before transaction
	if err := DatabaseHealthCheck(); err != nil {
		log.Printf("🚨 Database health check failed before transaction %s: %v", context, err)
		return err
	}

	start := time.Now()
	tx := DB.Begin()

	if tx.Error != nil {
		log.Printf("🚨 Failed to begin transaction for %s: %v", context, tx.Error)
		return tx.Error
	}

	defer func() {
		duration := time.Since(start)
		if duration > 10*time.Second {
			log.Printf("⚠️  Long transaction detected in %s: %v", context, duration)
		}
	}()

	// Execute the operation
	err := operation(tx)
	if err != nil {
		log.Printf("⚠️  Transaction operation failed in %s: %v", context, err)
		tx.Rollback()
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		log.Printf("🚨 Failed to commit transaction for %s: %v", context, err)
		tx.Rollback()
		return err
	}

	log.Printf("✅ Transaction completed successfully for %s in %v", context, time.Since(start))
	return nil
}

// StartDatabaseMonitoring starts a background goroutine to monitor database health
func StartDatabaseMonitoring() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("🚨 CRITICAL: Database monitoring goroutine panic: %v", r)
				log.Printf("Stack trace: %s", debug.Stack())
			}
		}()

		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if err := DatabaseHealthCheck(); err != nil {
				log.Printf("🚨 Scheduled database health check failed: %v", err)
			}
		}
	}()

	log.Println("🔄 Database health monitoring started")
}
