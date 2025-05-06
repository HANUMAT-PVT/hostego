package manager

import (
	"backend-hostego/internal/app/hostego-service/dto"
	"backend-hostego/internal/app/hostego-service/services"
	"backend-hostego/internal/app/hostego-service/utilities"
)

type AuthManager struct {
	authService *services.AuthService
}

func NewAuthManager(authService *services.AuthService) *AuthManager {
	return &AuthManager{
		authService: authService,
	}
}

func (am *AuthManager) CreateUser(reqCtx dto.ReqCtx, userInfo dto.AuthSignUpUserRequest) (string, error) {
	data, err := am.authService.CreateUsersAddressBySchema(reqCtx, userInfo)
	if err != nil {
		return "", err
	}
	jwtSecret, err := utilities.GenerateJWT(data)
	if err != nil {
		return "", err
	}

	return jwtSecret, nil
}
