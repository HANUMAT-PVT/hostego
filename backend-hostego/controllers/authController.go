package controllers

import (
	// "backend-hostego/config"
	"backend-hostego/config"
	"backend-hostego/database"
	"backend-hostego/models"

	// "context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Signup(c *fiber.Ctx) error {
	req := new(models.User)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad Request"})
	}

	var user models.User
	// Check if user already exists
	if err := database.DB.Where("mobile_number = ?", req.MobileNumber).First(&user).Error; err == nil {
		// ✅ If user exists, generate and return a token instead of creating a new user
		token, err := generateJWT(user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "JWT generation failed"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User already exists", "token": token})
	}

	// 🆕 If user doesn't exist, create a new user
	user = models.User{
		FirstName:           req.FirstName,
		LastName:            req.LastName,
		Email:               req.Email,
		MobileNumber:        req.MobileNumber,
		FirebaseOTPVerified: 1,
		CreatedAt:           time.Now(),
		LastLoginTimestamp:  time.Now(),
	}

	database.DB.Create(&user)

	// Generate JWT token for the new user
	token, err := generateJWT(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "JWT generation failed"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Signup Successful", "token": token})
}

func generateJWT(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    user.UserId,
		"first_name": user.FirstName,
		"email":      user.Email,
		"mobile":     user.MobileNumber,
		"exp":        time.Now().Add(24 * 30 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := config.GetEnv("HOSTEGO_JWT_SECRET_")
	return token.SignedString([]byte(jwtSecret))
}

//test commit
