package controllers

import (
	"backend-hostego/database"
	"backend-hostego/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	search := c.Query("search", "")
	offset := (page - 1) * limit
	var total int64
	var newUsersTotal int64
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User not found"})
	}
	var users []models.User
	database.DB.Limit(limit).Offset(offset).Where("first_name LIKE ? OR last_name LIKE ? OR email LIKE ? OR mobile_number LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%").Order("created_at DESC").Find(&users)
	database.DB.Model(&models.User{}).Count(&total)
	database.DB.Model(&models.User{}).Where("created_at > ?", time.Now().AddDate(0, 0, -30)).Count(&newUsersTotal)

	type UserWithRoles struct {
		User  models.User       `json:"user"`
		Roles []models.UserRole `json:"roles"`
		Total int64             `json:"total"`
	}

	// Initialize slice with proper length
	usersWithRoles := make([]UserWithRoles, len(users))

	for i := range users {
		var roles []models.UserRole
		database.DB.Preload("Role").Where("user_id = ?", users[i].UserId).Find(&roles)
		usersWithRoles[i].Roles = roles
		usersWithRoles[i].User = users[i]

	}
	return c.Status(200).JSON(fiber.Map{"users": usersWithRoles, "total": total, "new_users": newUsersTotal})
}

func GetUserById(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var user models.User

	if err := database.DB.First(&user, "user_id = ?", user_id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func UpdateUserById(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var user models.User
	var req models.User
	var existingUser models.User

	if err := database.DB.First(&user, "user_id=?", user_id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if req.MobileNumber != "" {
		if err := database.DB.Where("mobile_number = ?", req.MobileNumber).First(&existingUser).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		existingUser.AppleUserIdentifierId = user.AppleUserIdentifierId
		database.DB.Save(&existingUser)

		token, err := generateJWT(existingUser)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "JWT generation failed"})
		}
		database.DB.Delete(&user)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User updated successfully", "user": existingUser, "token": token})
	}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User updated successfully", "user": user})

}

func FetchUserByMobileNumber(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}
	mobileNumber := c.Params("mobile_number") // Get mobile number from URL params

	var user models.User

	// Query database for user with the given mobile number
	if err := database.DB.Where("mobile_number = ?", mobileNumber).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"user": user})
}

func DeleteUserById(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")

	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}
	var user models.User
	var userRoles []models.UserRole
	database.DB.Where("user_id=?", user_id).Find(&userRoles)
	if len(userRoles) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User has roles, cannot be deleted"})
	}
	if err := database.DB.First(&user, "user_id=?", user_id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	if err := database.DB.Delete(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	database.DB.Delete(&models.UserRole{}, "user_id=?", user_id)
	database.DB.Delete(&models.Wallet{}, "user_id=?", user_id)
	database.DB.Delete(&models.Order{}, "user_id=?", user_id)
	database.DB.Delete(&models.Address{}, "user_id=?", user_id)
	database.DB.Delete(&models.DeliveryPartner{}, "user_id=?", user_id)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User deleted successf~ully"})
}

//test commit
