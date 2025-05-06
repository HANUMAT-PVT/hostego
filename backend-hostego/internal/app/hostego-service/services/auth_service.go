package services

import (
	"backend-hostego/internal/app/hostego-service/dto"
	"backend-hostego/internal/app/hostego-service/models"
	"backend-hostego/internal/app/hostego-service/repository"
	"fmt"
)

type AuthService struct {
	*repository.BaseRepo
}

func NewAuthService(baseRepo *repository.BaseRepo) *AuthService {
	return &AuthService{BaseRepo: baseRepo}
}

func (as *AuthService) CreateUsersAddressBySchema(reqCtx dto.ReqCtx, userAddrInfo dto.AuthSignUpUserRequest) (models.User, error) {
	err := as.BaseRepo.Create(reqCtx, &userAddrInfo, "users")

	if err != nil {
		return models.User{}, err
	}

	var request models.User

	WhereClause := fmt.Sprintf("mobile_number = %s", userAddrInfo.MobileNumber)
	err = as.BaseRepo.GetRecordsByCondition(&request, WhereClause, "users")

	if err != nil {
		return models.User{}, err
	}

	return request, nil
}
