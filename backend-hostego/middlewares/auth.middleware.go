package middlewares

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyUserAuthCookie(c fiber.Ctx) (string, error) {
	cookie := c.Get("auth_token")
	
	if cookie == "" {
		return "", fmt.Errorf("no auth token found")
	}
	token, err := jwt.Parse(cookie, func(t *jwt.Token) (interface{}, error) {
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