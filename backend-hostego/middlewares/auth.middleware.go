package middlewares

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyUserAuthCookie(c fiber.Ctx) (int, error) {
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
