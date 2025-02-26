package middlewares

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyUserAuthCookie(c fiber.Ctx) (string, error) {
	authHeader := c.Get("Authorization")
	
	if authHeader == "" {
		return "", fmt.Errorf("no auth token found")
	}
	// 2️⃣ Extract the token (format: "Bearer <token>")
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return "", fmt.Errorf("invalid Authorization token format")
	}
	tokenString := splitToken[1]
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte("hanumat"), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user_id, exists := claims["user_id"].(string) // Ensure correct type
		if !exists {
			return "", fmt.Errorf("user_id not found in token")
		}
		return user_id, nil
	}

	return "", fmt.Errorf("invalid token")
}

// test
