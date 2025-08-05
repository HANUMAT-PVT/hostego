package middlewares

import (
	"backend-hostego/logs"
	"fmt"
	"log"
	"runtime/debug"
	"time"

	"github.com/gofiber/fiber/v2"
)

// CrashPreventionMiddleware provides comprehensive panic recovery, logging, and crash prevention
func CrashPreventionMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Start timer for request duration
		start := time.Now()

		// Capture user agent and important request info
		userAgent := c.Get("User-Agent")
		ip := c.IP()
		method := c.Method()
		path := c.Path()

		// Enhanced panic recovery
		defer func() {
			if r := recover(); r != nil {
				duration := time.Since(start)

				// Log detailed crash information using enhanced logger
				additionalData := map[string]interface{}{
					"duration":     duration,
					"method":       method,
					"path":         path,
					"ip":           ip,
					"user_agent":   userAgent,
					"query_params": string(c.Request().URI().QueryString()),
					"body_size":    len(c.Body()),
					"headers":      c.GetReqHeaders(),
				}

				logs.LogCrash(r, "HTTP Request Handler", additionalData)

				// Send appropriate error response
				if !c.Context().Response.Header.IsHTTP11() || c.Response().StatusCode() == 0 {
					err := c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error":      "Internal server error occurred. The issue has been logged.",
						"timestamp":  time.Now().Format(time.RFC3339),
						"request_id": fmt.Sprintf("%d", time.Now().UnixNano()),
					})
					if err != nil {
						log.Printf("Failed to send error response: %v", err)
					}
				}
			}
		}()

		// Process the request
		err := c.Next()

		duration := time.Since(start)
		status := c.Response().StatusCode()

		// Log access using enhanced logger (to access.log)
		logs.LogAccess(method, path, userAgent, ip, status, duration)

		// Also log to console/main log with user agent for immediate visibility
		if status >= 500 {
			log.Printf("🚨 %s %s | Status: %d | Duration: %v | IP: %s | UA: %s",
				method, path, status, duration, ip, userAgent)
			logs.LogError(fmt.Errorf("server error response"), "HTTP Response", map[string]interface{}{
				"status":     status,
				"path":       path,
				"method":     method,
				"user_agent": userAgent,
			})
		} else if status >= 400 {
			log.Printf("⚠️  %s %s | Status: %d | Duration: %v | IP: %s | UA: %s",
				method, path, status, duration, ip, userAgent)
			logs.LogWarning("client error response", "HTTP Response", map[string]interface{}{
				"status":     status,
				"path":       path,
				"method":     method,
				"user_agent": userAgent,
			})
		} else {
			// Log successful requests with user agent to console for monitoring
			log.Printf("✅ %s %s | Status: %d | Duration: %v | IP: %s | UA: %s",
				method, path, status, duration, ip, userAgent)
		}

		return err
	}
}

// SafeUserIDExtractor safely extracts user_id from context with proper error handling
func SafeUserIDExtractor(c *fiber.Ctx) (int, error) {
	userIDInterface := c.Locals("user_id")
	if userIDInterface == nil {
		return 0, fmt.Errorf("user_id not found in context")
	}

	switch v := userIDInterface.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	case string:
		// This shouldn't happen but let's handle it safely
		return 0, fmt.Errorf("user_id is string type, expected int")
	default:
		return 0, fmt.Errorf("user_id has unexpected type: %T", v)
	}
}

// DatabaseSafetyWrapper wraps database operations with enhanced error handling
func DatabaseSafetyWrapper(operation func() error, context string) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("🚨 Database operation panic in %s: %v", context, r)
			log.Printf("Stack trace: %s", debug.Stack())
		}
	}()

	err := operation()
	if err != nil {
		log.Printf("⚠️  Database operation failed in %s: %v", context, err)
	}
	return err
}
