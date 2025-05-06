package controller

import (
	"backend-hostego/internal/app/hostego-service/dto"
	"backend-hostego/internal/app/hostego-service/manager"
	"backend-hostego/internal/app/hostego-service/response"
	"backend-hostego/internal/app/hostego-service/utilities"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userManager *manager.UserManager
}

func NewUserController(userManager *manager.UserManager) *UserController {
	return &UserController{
		userManager: userManager,
	}
}

func (uc *UserController) GetUsers(ctx *gin.Context) {
	reqCtx := utilities.GetRCtx(ctx, "GetUsers")
	data, err := uc.userManager.GetUsers(reqCtx)

	if err != nil {
		response.Error(ctx, reqCtx, http.StatusInternalServerError, response.ErrorWhileProcessing, err.Error())
		return
	}

	response.Success(ctx, reqCtx, http.StatusOK, data)
}

func (uc *UserController) GetUserByUserId(ctx *gin.Context) {
	reqCtx := utilities.GetRCtx(ctx, "GetUserById")
	userId := reqCtx.UserId
	data, err := uc.userManager.GetUserById(reqCtx, userId)

	if err != nil {
		response.Error(ctx, reqCtx, http.StatusInternalServerError, response.ErrorWhileProcessing, err.Error())
		return
	}

	response.Success(ctx, reqCtx, http.StatusOK, data)
}

func (uc *UserController) UpdateUserById(ctx *gin.Context) {
	reqCtx := utilities.GetRCtx(ctx, "UpdateUserById")
	var userInfo dto.UserRequest
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
	err = uc.userManager.UpdateUserByModel(reqCtx, userInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusInternalServerError, response.ErrorWhileProcessing, err.Error())
		return
	}

	response.Success(ctx, reqCtx, http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (uc *UserController) UpdateUserAddress(ctx *gin.Context) {
	reqCtx := utilities.GetRCtx(ctx, "UpdateUserAddress")
	var userAddressInfo dto.UserAddressRequest
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusUnprocessableEntity, response.InvalidInput, "Error while reading the request body")
		return
	}

	err = json.Unmarshal(body, &userAddressInfo)

	if err != nil {
		response.Error(ctx, reqCtx, http.StatusUnprocessableEntity, response.InvalidInput, "Error while parsing the request body")
		return
	}

	if userAddressInfo.UserId == 0 {
		userAddressInfo.UserId = reqCtx.UserId
	}

	if userAddressInfo.UserId == 0 {
		response.Error(ctx, reqCtx, http.StatusBadRequest, response.InvalidInput, "User id is not passed and coming 0")
		return
	}

	validate := validator.New()
	err = validate.Struct(userAddressInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusBadRequest, response.InvalidInput, "Please fill the required fields"+err.Error())
		return
	}
	err = uc.userManager.UpdateUserAddressByModel(reqCtx, userAddressInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusInternalServerError, response.ErrorWhileProcessing, err.Error())
		return
	}

	response.Success(ctx, reqCtx, http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (uc *UserController) CreateUserAddress(ctx *gin.Context) {
	reqCtx := utilities.GetRCtx(ctx, "CreateUserAddress")
	var userAddressInfo dto.UserAddressCreateRequest
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusUnprocessableEntity, response.InvalidInput, "Error while reading the request body")
		return
	}

	err = json.Unmarshal(body, &userAddressInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusUnprocessableEntity, response.InvalidInput, "Error while parsing the request body")
		return
	}

	if userAddressInfo.UserId == 0 {
		userAddressInfo.UserId = reqCtx.UserId
	}

	if userAddressInfo.UserId == 0 {
		response.Error(ctx, reqCtx, http.StatusBadRequest, response.InvalidInput, "User id is not passed and coming 0")
		return
	}

	validate := validator.New()
	err = validate.Struct(userAddressInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusBadRequest, response.InvalidInput, "Please fill the required fields"+err.Error())
		return
	}
	err = uc.userManager.CreateUserAddressByModel(reqCtx, userAddressInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusInternalServerError, response.ErrorWhileProcessing, err.Error())
		return
	}

	response.Success(ctx, reqCtx, http.StatusOK, gin.H{"message": "User Address Created successfully"})
}

func (uc *UserController) DeleteUserAddress(ctx *gin.Context) {
	reqCtx := utilities.GetRCtx(ctx, "CreateUserAddress")
	var userAddressInfo dto.UserAddressRequest
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusUnprocessableEntity, response.InvalidInput, "Error while reading the request body")
		return
	}

	err = json.Unmarshal(body, &userAddressInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusUnprocessableEntity, response.InvalidInput, "Error while parsing the request body")
		return
	}

	if userAddressInfo.UserId == 0 {
		userAddressInfo.UserId = reqCtx.UserId
	}

	if userAddressInfo.UserId == 0 {
		response.Error(ctx, reqCtx, http.StatusBadRequest, response.InvalidInput, "User id is not passed and coming 0")
		return
	}

	validate := validator.New()
	err = validate.Struct(userAddressInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusBadRequest, response.InvalidInput, "Please fill the required fields"+err.Error())
		return
	}
	err = uc.userManager.DeleteUserAddressByModel(reqCtx, userAddressInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusInternalServerError, response.ErrorWhileProcessing, err.Error())
		return
	}

	response.Success(ctx, reqCtx, http.StatusOK, gin.H{"message": "User Address Deleted successfully"})
}
