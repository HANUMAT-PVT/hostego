package controllers

import (
	"backend-hostego/database"
	"backend-hostego/middlewares"

	"backend-hostego/models"

	"github.com/gofiber/fiber/v3"
)

func GetUsers(c fiber.Ctx) error {
	user_id, err := middlewares.VerifyUserAuthCookie(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "message": "You are not Authenticated !"})
	}
	if user_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User not found"})
	}
	var users []models.User
	database.DB.Find(&users)
	
	for _,user:=range users{
		var roles []models.Role
		database.DB.Where("user_id = ?",user.UserId).Find(&roles)
	}
	return c.Status(200).JSON(users)

}

func GetUserById(c fiber.Ctx) error {

	user_id, err := middlewares.VerifyUserAuthCookie(c)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "message": "You are not Authenticated !"})
	}

	var user models.User

	if err := database.DB.First(&user, "user_id = ?", user_id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

func UpdateUserById(c fiber.Ctx) error {
	user_id, midError := middlewares.VerifyUserAuthCookie(c)

	if midError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": "Bad request"})
	}
	var user models.User

	if err := database.DB.First(&user, "user_id=?", user_id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	if err := c.Bind().JSON(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User updated successfully", "user": user})

}

func FetchUserByMobileNumber(c fiber.Ctx) error {
	mobileNumber := c.Params("mobile_number") // Get mobile number from URL params

	var user models.User

	// Query database for user with the given mobile number
	if err := database.DB.Where("mobile_number = ?", mobileNumber).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"user": user})
}

//test commit
