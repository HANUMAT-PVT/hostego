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

type AuthController struct {
	authManager *manager.AuthManager
}

func NewAuthController(authManager *manager.AuthManager) *AuthController {
	return &AuthController{
		authManager: authManager,
	}
}

func (ac *AuthController) CreateUserBySignUp(ctx *gin.Context) {
	reqCtx := utilities.GetRCtx(ctx, "Signup")
	var userInfo dto.AuthSignUpUserRequest
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

	validate := validator.New()
	err = validate.Struct(userInfo)
	if err != nil {
		response.Error(ctx, reqCtx, http.StatusBadRequest, response.InvalidInput, "Please fill the required fields"+err.Error())
		return
	}
	jwtstring, err := ac.authManager.CreateUser(reqCtx, userInfo)

	if err != nil {
		response.Error(ctx, reqCtx, http.StatusInternalServerError, response.ErrorWhileProcessing, err.Error())
		return
	}

	response.Success(ctx, reqCtx, http.StatusOK, jwtstring)
}
