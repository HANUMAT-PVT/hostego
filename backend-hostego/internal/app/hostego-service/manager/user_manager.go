package manager

import (
	"backend-hostego/internal/app/hostego-service/dto"
	"backend-hostego/internal/app/hostego-service/models"
	"backend-hostego/internal/app/hostego-service/services"
)

type UserManager struct {
	userService *services.UserService
}

func NewUserManager(userService *services.UserService) *UserManager {
	return &UserManager{
		userService: userService,
	}
}

func (um *UserManager) GetUsers(reqCtx dto.ReqCtx) ([]models.User, error) {
	data, err := um.userService.GetUsers(reqCtx)

	if err != nil {
		return nil, err
	}

	return data, err
}

func (um *UserManager) GetUserById(reqCtx dto.ReqCtx, userId int) (models.User, error) {
	data, err := um.userService.GetUsersById(reqCtx, userId)
	if err != nil {
		return models.User{}, err
	}
	return data, err
}

func (um *UserManager) UpdateUserByModel(reqCtx dto.ReqCtx, userInfo dto.UserRequest) error {
	err := um.userService.UpdateUsersBySchema(reqCtx, userInfo)
	if err != nil {
		return err
	}
	return err
}

func (um *UserManager) UpdateUserAddressByModel(reqCtx dto.ReqCtx, userAddrInfo dto.UserAddressRequest) error {
	err := um.userService.UpdateUsersAddressBySchema(reqCtx, userAddrInfo)
	if err != nil {
		return err
	}
	return err
}

func (um *UserManager) CreateUserAddressByModel(reqCtx dto.ReqCtx, userAddrInfo dto.UserAddressCreateRequest) error {
	err := um.userService.CreateUsersAddressBySchema(reqCtx, userAddrInfo)
	if err != nil {
		return err
	}
	return err
}

func (um *UserManager) DeleteUserAddressByModel(reqCtx dto.ReqCtx, userAddrInfo dto.UserAddressRequest) error {
	err := um.userService.DeleteUsersAddressBySchema(reqCtx, userAddrInfo)
	if err != nil {
		return err
	}
	return err
}
