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

type ProductController struct {
	userManager *manager.UserManager
}

func NewProductController(userManager *manager.UserManager) *ProductController {
	return &ProductController{
		userManager: userManager,
	}
}

func (sc *ShopController) CreateProduct(ctx *gin.Context) {
	reqCtx := utilities.GetRCtx(ctx, "UpdateShop")
	var shopInfo dto.CreateProductRequestDto
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
