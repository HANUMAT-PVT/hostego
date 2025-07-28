package middlewares

import (
	"backend-hostego/config"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyUserAuthCookie(c *fiber.Ctx) (int, error) {

	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return 0, fmt.Errorf("no auth token found")
	}
	// 2️⃣ Extract the token (format: "Bearer <token>")
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return 0, fmt.Errorf("invalid Authorization token format")
	}
	tokenString := splitToken[1]
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte("hanumat"), nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user_id, exists := claims["user_id"]

		if !exists {
			return 0, fmt.Errorf("user_id not found in token")
		}

		// Type assert user_id to int
		if id, ok := user_id.(int); ok {
			return id, nil
		}
		// Try float64 if int fails
		if idFloat, ok := user_id.(float64); ok {
			return int(idFloat), nil
		}
		return 0, fmt.Errorf("invalid user_id type")
	}

	return 0, fmt.Errorf("invalid token")
}

// test

func VerifyUserAuthCookieMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "No auth token found",
			})
		}

		// Extract JWT token from "Bearer <token>"
		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid Authorization token format",
			})
		}
		tokenString := splitToken[1]

		// Parse JWT token
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			jwtSecret := config.GetEnv("HOSTEGO_JWT_SECRET_")
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Extract claims from JWT
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		// Get user_id from claims
		userID, exists := claims["user_id"]
		if !exists {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User ID not found in token",
			})
		}

		// Convert user_id to int
		var uid int
		switch v := userID.(type) {
		case float64: // JWT stores numbers as float64
			uid = int(v)
		case int:
			uid = v
		default:
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user ID type",
			})
		}

		// Store user_id in request context

		c.Locals("user_id", uid)

		return c.Next()
	}
}

