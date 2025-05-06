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

type ShopController struct {
	ShopManager *manager.ShopManager
}

func NewShopController(ShopManager *manager.ShopManager) *ShopController {
	return &ShopController{
		ShopManager: ShopManager,
	}
}

func (sc *ShopController) GetShops(ctx *gin.Context) {
	reqCtx := utilities.GetRCtx(ctx, "GetShops")
	data, err := sc.ShopManager.GetAllShops(reqCtx)

	if err != nil {
		response.Error(ctx, reqCtx, http.StatusInternalServerError, response.ErrorWhileProcessing, err.Error())
		return
	}

	response.Success(ctx, reqCtx, http.StatusOK, data)
}

func (sc *ShopController) GetShopByShopId(ctx *gin.Context) {
	reqCtx := utilities.GetRCtx(ctx, "GetShopByShopId")
	var shopInfo dto.ShopRequest
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusUnprocessableEntity, response.InvalidInput, "Error while reading the request body")
		return
	}

	err = json.Unmarshal(body, &shopInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusUnprocessableEntity, response.InvalidInput, "Error while parsing the request body")
		return

	}
	if shopInfo.UserId == 0 {
		shopInfo.UserId = reqCtx.UserId
	}
	if shopInfo.UserId == 0 {
		response.Error(ctx, reqCtx, http.StatusBadRequest, response.InvalidInput, "User id is not passed and coming 0")
		return
	}
	// userInfo.UserId = reqCtx.UserId
	validate := validator.New()
	err = validate.Struct(shopInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusBadRequest, response.InvalidInput, "Please fill the required fields"+err.Error())
		return
	}
	data, err := sc.ShopManager.GetShopByShopId(reqCtx, shopInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusInternalServerError, response.ErrorWhileProcessing, err.Error())
		return
	}
	response.Success(ctx, reqCtx, http.StatusOK, data)
}

func (sc *ShopController) GetShopByUserId(ctx *gin.Context) {
	reqCtx := utilities.GetRCtx(ctx, "GetShopByUserId")
	userId := reqCtx.UserId
	data, err := sc.ShopManager.GetShopByUserId(reqCtx, userId)

	if err != nil {
		response.Error(ctx, reqCtx, http.StatusInternalServerError, response.ErrorWhileProcessing, err.Error())
		return
	}

	response.Success(ctx, reqCtx, http.StatusOK, data)
}

func (sc *ShopController) UpdateShop(ctx *gin.Context) {
	reqCtx := utilities.GetRCtx(ctx, "UpdateShop")
	var shopInfo dto.UpdateShopRequest
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusUnprocessableEntity, response.InvalidInput, "Error while reading the request body")
		return
	}

	err = json.Unmarshal(body, &shopInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusUnprocessableEntity, response.InvalidInput, "Error while parsing the request body")
		return

	}
	if shopInfo.UserId == 0 {
		shopInfo.UserId = reqCtx.UserId
	}
	if shopInfo.UserId == 0 {
		response.Error(ctx, reqCtx, http.StatusBadRequest, response.InvalidInput, "User id is not passed and coming 0")
		return
	}
	// userInfo.UserId = reqCtx.UserId
	validate := validator.New()
	err = validate.Struct(shopInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusBadRequest, response.InvalidInput, "Please fill the required fields"+err.Error())
		return
	}
	err = sc.ShopManager.UpdateShopByShopId(reqCtx, shopInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusInternalServerError, response.ErrorWhileProcessing, err.Error())
		return
	}
	response.Success(ctx, reqCtx, http.StatusOK, gin.H{"message": "Shop Updated Succesfully"})
}

func (sc *ShopController) CreateShop(ctx *gin.Context) {
	reqCtx := utilities.GetRCtx(ctx, "UpdateShop")
	var shopInfo dto.CreateShopRequest
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusUnprocessableEntity, response.InvalidInput, "Error while reading the request body")
		return
	}

	err = json.Unmarshal(body, &shopInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusUnprocessableEntity, response.InvalidInput, "Error while parsing the request body")
		return

	}
	if shopInfo.UserId == 0 {
		shopInfo.UserId = reqCtx.UserId
	}
	if shopInfo.UserId == 0 {
		response.Error(ctx, reqCtx, http.StatusBadRequest, response.InvalidInput, "User id is not passed and coming 0")
		return
	}
	// userInfo.UserId = reqCtx.UserId
	validate := validator.New()
	err = validate.Struct(shopInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusBadRequest, response.InvalidInput, "Please fill the required fields"+err.Error())
		return
	}
	err = sc.ShopManager.CreateShop(reqCtx, shopInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusInternalServerError, response.ErrorWhileProcessing, err.Error())
		return
	}
	response.Success(ctx, reqCtx, http.StatusOK, gin.H{"message": "Shop Created Succesfully"})
}
