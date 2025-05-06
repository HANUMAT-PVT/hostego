package services

import (
	"backend-hostego/internal/app/hostego-service/dto"
	"backend-hostego/internal/app/hostego-service/models"
	"backend-hostego/internal/app/hostego-service/repository"
	"fmt"
	"time"
)

type UserService struct {
	*repository.BaseRepo
}

func NewUserService(baseRepo *repository.BaseRepo) *UserService {
	return &UserService{BaseRepo: baseRepo}
}

func (um *UserService) GetUsers(reqCtx dto.ReqCtx) ([]models.User, error) {
	var request []models.User
	filterModel := models.User{}
	err := um.BaseRepo.GetRecords(&request, filterModel, "users")

	if err != nil {
		return nil, err
	}

	return request, err
}

func (um *UserService) GetUsersById(reqCtx dto.ReqCtx, userId int) (models.User, error) {
	var request models.User

	WhereClause := fmt.Sprintf("user_id = %d", userId)
	err := um.BaseRepo.GetRecordsByCondition(&request, WhereClause, "users")

	if err != nil {
		return models.User{}, err
	}

	return request, err
}

func (um *UserService) UpdateUsersBySchema(reqCtx dto.ReqCtx, userInfo dto.UserRequest) error {
	userModel := models.User{}

	err := um.BaseRepo.GetRecordsByCondition(&userModel, fmt.Sprintf("user_id = %v", userInfo.UserId), "users")
	if err != nil {
		return err
	}

	userInfo.UpdatedAt = &time.Time{}

	whereClause := fmt.Sprintf("user_id='%v'", userInfo.UserId)
	err = um.BaseRepo.UpdateModelWithCondition(reqCtx, "users", whereClause, &userInfo)
	if err != nil {
		return err
	}
	return nil
}

func (um *UserService) UpdateUsersAddressBySchema(reqCtx dto.ReqCtx, userAddrInfo dto.UserAddressRequest) error {
	userAddressModel := models.UserAddress{}

	err := um.BaseRepo.GetRecordsByCondition(&userAddressModel, fmt.Sprintf("user_id = %v and address_id = '%v'", userAddrInfo.UserId, userAddrInfo.AddressId), "users_address")
	if err != nil {
		return err
	}
	userAddrInfo.UpdatedAt = &time.Time{}
	whereClause := fmt.Sprintf("user_id = %v and address_id = '%v'", userAddrInfo.UserId, userAddrInfo.AddressId)
	err = um.BaseRepo.UpdateModelWithCondition(reqCtx, "users_address", whereClause, &userAddrInfo)
	if err != nil {
		return err
	}
	return nil
}

func (um *UserService) CreateUsersAddressBySchema(reqCtx dto.ReqCtx, userAddrInfo dto.UserAddressCreateRequest) error {
	err := um.BaseRepo.Create(reqCtx, &userAddrInfo, "users_address")
	if err != nil {
		return err
	}
	return nil
}

func (um *UserService) DeleteUsersAddressBySchema(reqCtx dto.ReqCtx, userAddrInfo dto.UserAddressRequest) error {
	err := um.BaseRepo.Delete(&userAddrInfo, "users_address")
	if err != nil {
		return err
	}
	return nil
}
