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

type CartController struct {
	CartManager *manager.CartManager
}

func NewCartController(CartManager *manager.CartManager) *CartController {
	return &CartController{
		CartManager: CartManager,
	}
}

func (cc *CartController) AddProductInUserCart(ctx *gin.Context) {
	reqCtx := utilities.GetRCtx(ctx, "AddProductInUserCart")
	var CartInfo dto.CartItemRequest
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusUnprocessableEntity, response.InvalidInput, "Error while reading the request body")
		return
	}

	err = json.Unmarshal(body, &CartInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusUnprocessableEntity, response.InvalidInput, "Error while parsing the request body")
		return
	}
	CartInfo.UserId = reqCtx.UserId
	validate := validator.New()
	err = validate.Struct(CartInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusBadRequest, response.InvalidInput, "Please fill the required fields"+err.Error())
		return
	}
	err = cc.CartManager.AddProductInCart(reqCtx, CartInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusInternalServerError, response.ErrorWhileProcessing, err.Error())
		return
	}

	response.Success(ctx, reqCtx, http.StatusOK, gin.H{"message": "Product Added Succesfully"})
}

func (cc *CartController) FetchUserCart(ctx *gin.Context) {
	reqCtx := utilities.GetRCtx(ctx, "FetchUserCart")
	var CartInfo dto.CartItemRequest
	CartInfo.UserId = reqCtx.UserId

	data, err := cc.CartManager.FetchCartByUserId(reqCtx, CartInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusInternalServerError, response.ErrorWhileProcessing, err.Error())
		return
	}

	response.Success(ctx, reqCtx, http.StatusOK, data)
}
