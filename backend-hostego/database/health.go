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
	log.Printf("📊 DB Pool Stats - Open: %d, InUse: %d, Idle: %d, WaitCount: %d, WaitDuration: %v",
		stats.OpenConnections, stats.InUse, stats.Idle, stats.WaitCount, stats.WaitDuration)

	// Updated warning threshold to match new pool settings
	if stats.OpenConnections > 800 {
		log.Printf("⚠️  High number of open connections: %d (max: 1000)", stats.OpenConnections)
	}

	// Check for connection wait issues
	if stats.WaitCount > 0 {
		log.Printf("⚠️  Connection pool wait detected: %d waits, total duration: %v", stats.WaitCount, stats.WaitDuration)
	}

	// Check idle connection health
	if stats.Idle < 50 && stats.OpenConnections > 200 {
		log.Printf("⚠️  Low idle connections: %d idle out of %d total", stats.Idle, stats.OpenConnections)
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

// EnsureTransactionCleanup ensures transactions are properly closed even in error scenarios
func EnsureTransactionCleanup(tx *gorm.DB, context string) {
	if tx != nil {
		if r := recover(); r != nil {
			log.Printf("🚨 CRITICAL: Transaction panic in %s: %v", context, r)
			tx.Rollback()
		}
	}
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

		healthTicker := time.NewTicker(30 * time.Second)
		cleanupTicker := time.NewTicker(5 * time.Minute) // Force cleanup every 5 minutes
		defer healthTicker.Stop()
		defer cleanupTicker.Stop()

		// Log initial stats
		LogConnectionPoolStats()

		for {
			select {
			case <-healthTicker.C:
				if err := DatabaseHealthCheck(); err != nil {
					log.Printf("🚨 Scheduled database health check failed: %v", err)
				}

				// Log detailed stats every 2 minutes
				if time.Now().Second() < 30 {
					LogConnectionPoolStats()
				}

			case <-cleanupTicker.C:
				// Force cleanup of potentially stuck connections
				log.Printf("🔄 Performing scheduled connection pool cleanup...")
				ForceConnectionCleanup()
			}
		}
	}()

	log.Println("🔄 Database health monitoring started with enhanced connection pool tracking and periodic cleanup")
}

// LogConnectionPoolStats logs detailed connection pool statistics
func LogConnectionPoolStats() {
	if DB == nil {
		return
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("🚨 Failed to get underlying sql.DB for stats: %v", err)
		return
	}

	stats := sqlDB.Stats()
	log.Printf("🔍 Connection Pool Analysis:")
	log.Printf("   📊 Open Connections: %d", stats.OpenConnections)
	log.Printf("   🔄 In Use: %d", stats.InUse)
	log.Printf("   💤 Idle: %d", stats.Idle)
	log.Printf("   ⏳ Wait Count: %d", stats.WaitCount)
	log.Printf("   ⏱️  Wait Duration: %v", stats.WaitDuration)
	log.Printf("   🚫 Max Open Connections: %d", stats.MaxOpenConnections)

	// Calculate utilization percentage
	utilization := float64(stats.InUse) / float64(stats.MaxOpenConnections) * 100
	log.Printf("   📈 Utilization: %.1f%%", utilization)

	if utilization > 80 {
		log.Printf("⚠️  High connection pool utilization: %.1f%%", utilization)
	}
}

// SafeTransactionWithCleanup provides a bulletproof transaction wrapper that ensures connections are always returned to the pool
func SafeTransactionWithCleanup(operation func(tx *gorm.DB) error, context string) error {
	var tx *gorm.DB
	var committed bool

	defer func() {
		if r := recover(); r != nil {
			log.Printf("🚨 CRITICAL: Transaction panic in %s: %v", context, r)
			log.Printf("Stack trace: %s", debug.Stack())

			// Ensure rollback on panic
			if tx != nil && !committed {
				if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
					log.Printf("🚨 Failed to rollback transaction after panic in %s: %v", context, rollbackErr)
				}
			}
		}
	}()

	// Health check before transaction
	if err := DatabaseHealthCheck(); err != nil {
		log.Printf("🚨 Database health check failed before transaction %s: %v", context, err)
		return err
	}

	start := time.Now()
	tx = DB.Begin()

	if tx.Error != nil {
		log.Printf("🚨 Failed to begin transaction for %s: %v", context, tx.Error)
		return tx.Error
	}

	// Execute the operation
	err := operation(tx)
	if err != nil {
		log.Printf("⚠️  Transaction operation failed in %s: %v", context, err)
		if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
			log.Printf("🚨 Failed to rollback transaction in %s: %v", context, rollbackErr)
		}
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		log.Printf("🚨 Failed to commit transaction for %s: %v", context, err)
		if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
			log.Printf("🚨 Failed to rollback transaction after commit failure in %s: %v", context, rollbackErr)
		}
		return err
	}

	committed = true
	log.Printf("✅ Transaction completed successfully for %s in %v", context, time.Since(start))
	return nil
}

// ForceConnectionCleanup forces cleanup of idle connections that might be stuck
func ForceConnectionCleanup() {
	if DB == nil {
		return
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("🚨 Failed to get underlying sql.DB for cleanup: %v", err)
		return
	}

	// Log stats before cleanup
	beforeStats := sqlDB.Stats()
	log.Printf("🔄 Before cleanup - Open: %d, InUse: %d, Idle: %d",
		beforeStats.OpenConnections, beforeStats.InUse, beforeStats.Idle)

	// Force close idle connections by setting max idle to 0 temporarily
	sqlDB.SetMaxIdleConns(0)

	// Wait a moment for connections to close
	time.Sleep(100 * time.Millisecond)

	// Restore original setting
	sqlDB.SetMaxIdleConns(500)

	// Log stats after cleanup
	afterStats := sqlDB.Stats()
	log.Printf("🔄 After cleanup - Open: %d, InUse: %d, Idle: %d",
		afterStats.OpenConnections, afterStats.InUse, afterStats.Idle)

	// Log improvement
	closed := beforeStats.OpenConnections - afterStats.OpenConnections
	if closed > 0 {
		log.Printf("✅ Successfully closed %d stuck connections", closed)
	} else {
		log.Printf("ℹ️  No stuck connections found to clean up")
	}
}
