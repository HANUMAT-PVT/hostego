package controllers

import (
	"backend-hostego/database"
	"backend-hostego/middlewares"
	"backend-hostego/models"

	"github.com/gofiber/fiber/v3"
)

func CreateNewDeliveryPartner(c fiber.Ctx) error {
	var delivery_partner models.DeliveryPartner
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}
	if err := database.DB.Where("user_id=?", user_id).Find(&delivery_partner).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": err.Error()})
	}
	delivery_partner.UserId = user_id
	if err := c.Bind().JSON(&delivery_partner); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": err.Error()})
	}
	if err := database.DB.Preload("User").Create(&delivery_partner).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"delivery_partner": delivery_partner})
}

func UpdateDeliveryPartner(c fiber.Ctx) error {
	var delivery_partner models.DeliveryPartner
	delivery_partner_id := c.Params("id")
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}

	// First find the existing delivery partner
	if err := database.DB.Where("delivery_partner_id = ? AND user_id = ?", delivery_partner_id, user_id).First(&delivery_partner).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Delivery partner not found"})
	}

	// Bind the updated data
	if err := c.Bind().JSON(&delivery_partner); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Save the updates
	if err := database.DB.Save(&delivery_partner).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"delivery_partner": delivery_partner,
		"message":          "Delivery Partner Updated successfully",
	})
}

func FetchDeliveryPartnerByUserId(c fiber.Ctx) error {

	user_id, err := middlewares.VerifyUserAuthCookie(c)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "message": "You are not Authenticated !"})
	}
	if user_id != "" {
	}

	var delivery_partner models.DeliveryPartner

	if err := database.DB.Preload("User").First(&delivery_partner, "user_id = ?", user_id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Delivery partner not found"})
	}
	return c.Status(fiber.StatusOK).JSON(delivery_partner)
}
