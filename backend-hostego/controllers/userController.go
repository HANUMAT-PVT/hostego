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
	if user_id ==0{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User not found"})
	}
	var users []models.User
	database.DB.Find(&users)

	type UserWithRoles struct {
		User  models.User       `json:"user"`
		Roles []models.UserRole `json:"roles"`
	}

	// Initialize slice with proper length
	usersWithRoles := make([]UserWithRoles, len(users))

	for i := range users {
		var roles []models.UserRole
		database.DB.Preload("Role").Where("user_id = ?", users[i].UserId).Find(&roles)
		usersWithRoles[i].Roles = roles
		usersWithRoles[i].User = users[i]
	}
	return c.Status(200).JSON(usersWithRoles)
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
