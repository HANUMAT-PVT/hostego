package controller

import (
	"backend-hostego/internal/app/hostego-service/dto"
	"backend-hostego/internal/app/hostego-service/manager"
	"backend-hostego/internal/app/hostego-service/response"
	"backend-hostego/internal/app/hostego-service/utilities"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type DeliverPartnerController struct {
	DeliveryPartnerManager *manager.DeliveryPartnerManager
}

func NewDeliverPartnerController(DeliveryPartnerManager *manager.DeliveryPartnerManager) *DeliverPartnerController {
	return &DeliverPartnerController{
		DeliveryPartnerManager: DeliveryPartnerManager,
	}
}

func (dpc *DeliverPartnerController) CreateDelivaryPartnerById(ctx *gin.Context) {
	reqCtx := utilities.GetRCtx(ctx, "CreateDelivaryPartnerById")
	var userInfo dto.DeliveryPartnerRequestDto
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusUnprocessableEntity, response.InvalidInput, "Error while reading the request body")
		return
	}
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusUnprocessableEntity, response.InvalidInput, "Error while parsing the request body")
		return

	}
	if userInfo.UserId == 0 {
		userInfo.UserId = reqCtx.UserId
	}
	if userInfo.UserId == 0 {
		response.Error(ctx, reqCtx, http.StatusBadRequest, response.InvalidInput, "User id is not passed and coming 0")
		return
	}
	// userInfo.UserId = reqCtx.UserId
	validate := validator.New()
	err = validate.Struct(userInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusBadRequest, response.InvalidInput, "Please fill the required fields"+err.Error())
		return
	}
	err = dpc.DeliveryPartnerManager.CreateDeliveryPartnerData(reqCtx, userInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusInternalServerError, response.ErrorWhileProcessing, err.Error())
		return
	}

	response.Success(ctx, reqCtx, http.StatusOK, gin.H{"message": "Delivery Partner Created successfully"})
}

func (dpc *DeliverPartnerController) FetchDelivaryPartnerById(ctx *gin.Context) {
	reqCtx := utilities.GetRCtx(ctx, "FetchDelivaryPartnerById")

	data, err := dpc.DeliveryPartnerManager.FetchDeliveryPartnerData(reqCtx)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusInternalServerError, response.ErrorWhileProcessing, err.Error())
		return
	}

	response.Success(ctx, reqCtx, http.StatusOK, data)
}

func (dpc *DeliverPartnerController) UpdateDeiveryPartnerById(ctx *gin.Context) {
	reqCtx := utilities.GetRCtx(ctx, "UpdateDeiveryPartnerById")
	var userInfo dto.DeliveryPartnerRequestDto
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusUnprocessableEntity, response.InvalidInput, "Error while reading the request body")
		return
	}
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusUnprocessableEntity, response.InvalidInput, "Error while parsing the request body")
		return

	}
	if userInfo.UserId == 0 {
		userInfo.UserId = reqCtx.UserId
	}
	if userInfo.UserId == 0 {
		response.Error(ctx, reqCtx, http.StatusBadRequest, response.InvalidInput, "User id is not passed and coming 0")
		return
	}
	// userInfo.UserId = reqCtx.UserId
	validate := validator.New()
	err = validate.Struct(userInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusBadRequest, response.InvalidInput, "Please fill the required fields"+err.Error())
		return
	}
	err = dpc.DeliveryPartnerManager.UpdateDeliveryPartnerInfo(reqCtx, userInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusInternalServerError, response.ErrorWhileProcessing, err.Error())
		return
	}

	response.Success(ctx, reqCtx, http.StatusOK, gin.H{"message": "User updated successfully"})
}

// func CreateNewDeliveryPartner(c fiber.Ctx) error {
// 	var delivery_partner models.DeliveryPartner
// 	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
// 	if middleErr != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
// 	}
// 	delivery_partner.UserId = user_id
// 	if err := c.Bind().JSON(&delivery_partner); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": err.Error()})
// 	}
// 	if err := database.DB.Preload("User").Create(&delivery_partner).Error; err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": err.Error()})
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"delivery_partner": delivery_partner})
// }
