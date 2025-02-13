package controllers

import (
	// "backend-hostego/config"
	"backend-hostego/database"
	"backend-hostego/models"
	// "context"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("hanumat")

func Signup(c fiber.Ctx) error {
	req := new(models.User)
	if err := c.Bind().JSON(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad Request"})
	}



	// ctx := context.Background()
	// authClient, err := config.FireBaseApp.Auth(ctx)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Firebase auth failed"})
	// }



	var user models.User
	if err := database.DB.Where(`mobile_number = ?`, req.MobileNumber).First(&user); err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "User Already Exists"})
	}
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

	token, err := generateJWT(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "JWT generation failed"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(24 * 30 * time.Hour),
		HTTPOnly: false,
		Secure:   false,
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"messsage": "Signup Succesfull", "token": token})

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
	return token.SignedString(jwtSecret)
}


//test commit