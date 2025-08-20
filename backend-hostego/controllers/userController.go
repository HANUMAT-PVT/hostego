package controllers

import (
	"backend-hostego/database"
	"backend-hostego/models"
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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

	// Get query parameters for date filtering
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	var err error

	// Parse start date
	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid start_date format. Use YYYY-MM-DD"})
		}
	}

	// Parse end date
	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid end_date format. Use YYYY-MM-DD"})
		}
		// Set end date to end of day
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}
	var total int64
	var newUsersTotal int64
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User not found"})
	}
	var users []models.User
	userQuery := database.DB.Limit(limit).Offset(offset).Where("first_name LIKE ? OR last_name LIKE ? OR email LIKE ? OR mobile_number LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")

	// Apply date filters if provided
	if startDateStr != "" {
		userQuery = userQuery.Where("created_at >= ?", startDate)
	}
	if endDateStr != "" {
		userQuery = userQuery.Where("created_at <= ?", endDate)
	}

	userQuery.Order("created_at DESC").Find(&users)

	// Count total users with same filters (excluding pagination)
	totalQuery := database.DB.Model(&models.User{}).Where("first_name LIKE ? OR last_name LIKE ? OR email LIKE ? OR mobile_number LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")

	// Apply date filters if provided
	if startDateStr != "" {
		totalQuery = totalQuery.Where("created_at >= ?", startDate)
	}
	if endDateStr != "" {
		totalQuery = totalQuery.Where("created_at <= ?", endDate)
	}

	totalQuery.Count(&total)

	// Count new users - if date filters are provided, use them; otherwise use last 30 days
	newUsersQuery := database.DB.Model(&models.User{})
	if startDateStr != "" || endDateStr != "" {
		// If date filters are provided, use the same date range
		if startDateStr != "" {
			newUsersQuery = newUsersQuery.Where("created_at >= ?", startDate)
		}
		if endDateStr != "" {
			newUsersQuery = newUsersQuery.Where("created_at <= ?", endDate)
		}
	} else {
		// Default: users from last 30 days
		newUsersQuery = newUsersQuery.Where("created_at > ?", time.Now().AddDate(0, 0, -30))
	}
	newUsersQuery.Count(&newUsersTotal)

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
	userID := c.Locals("user_id")
	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var user models.User
	if err := database.DB.First(&user, "user_id = ?", userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	var req models.User
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// If mobile number is provided
	if req.MobileNumber != "" {
		var existingUser models.User
		err := database.DB.Where("mobile_number = ?", req.MobileNumber).First(&existingUser).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Mobile number not in use → just update current user
				user.MobileNumber = req.MobileNumber
				if err := database.DB.Save(&user).Error; err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
				}
				token, _ := generateJWT(user)
				return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User updated successfully", "user": user, "token": token})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
		}

		// Mobile number exists → merge accounts
		existingUser.AppleUserIdentifierId = user.AppleUserIdentifierId
		if err := database.DB.Save(&existingUser).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to merge accounts"})
		}

		token, err := generateJWT(existingUser)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "JWT generation failed"})
		}

		database.DB.Delete(&user)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User merged successfully", "user": existingUser, "token": token})
	}

	// No mobile number provided → update other fields
	// Create a map for partial updates to avoid overwriting with empty values
	updateMap := make(map[string]interface{})

	// Only include non-empty fields in the update
	if req.FirstName != "" {
		updateMap["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		updateMap["last_name"] = req.LastName
	}
	if req.Email != "" {
		updateMap["email"] = req.Email
	}
	if req.FCMToken != "" {
		updateMap["fcm_token"] = req.FCMToken
	}
	if req.AppleUserIdentifierId != "" {
		updateMap["apple_user_identifier_id"] = req.AppleUserIdentifierId
	}

	// Only update if there are fields to update
	if len(updateMap) > 0 {
		if err := database.DB.Model(&user).Updates(updateMap).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	// Fetch the updated user to return
	if err := database.DB.First(&user, "user_id = ?", userID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch updated user"})
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
