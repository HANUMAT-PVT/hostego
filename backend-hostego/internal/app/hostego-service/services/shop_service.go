package services

import (
	"backend-hostego/internal/app/hostego-service/dto"
	"backend-hostego/internal/app/hostego-service/models"
	"backend-hostego/internal/app/hostego-service/repository"
	"fmt"
	"time"
)

type ShopService struct {
	*repository.BaseRepo
}

func NewShopService(baseRepo *repository.BaseRepo) *ShopService {
	return &ShopService{BaseRepo: baseRepo}
}

func (um *ShopService) GetShops(reqCtx dto.ReqCtx) ([]models.Shop, error) {
	var request []models.Shop
	filterModel := models.Shop{}
	err := um.BaseRepo.GetRecords(&request, filterModel, "shops")

	if err != nil {
		return nil, err
	}
	return request, err
}

func (um *ShopService) GetShopsByUser(reqCtx dto.ReqCtx, userId int) ([]models.Shop, error) {
	var request []models.Shop
	filterModel := models.Shop{UserId: (userId)}
	err := um.BaseRepo.GetRecords(&request, filterModel, "shops")

	if err != nil {
		return nil, err
	}
	return request, err
}

func (um *ShopService) GetShopsByShopsId(reqCtx dto.ReqCtx, shops dto.ShopRequest) ([]models.Shop, error) {
	var request []models.Shop
	filterModel := models.Shop{ShopId: (shops.ShopId)}
	err := um.BaseRepo.GetRecords(&request, filterModel, "shops")

	if err != nil {
		return nil, err
	}
	return request, err
}

func (um *ShopService) UpdateShopsByShopsId(reqCtx dto.ReqCtx, shopsinfo dto.UpdateShopRequest) error {
	shopModel := models.Shop{}

	err := um.BaseRepo.GetRecordsByCondition(&shopModel, fmt.Sprintf("user_id = %v and shop_id = '%v'", shopsinfo.UserId, shopsinfo.ShopId), "shops")
	if err != nil {
		return err
	}

	shopsinfo.UpdatedAt = &time.Time{}

	whereClause := fmt.Sprintf("user_id = %v and shop_id = '%v'", shopsinfo.UserId, shopsinfo.ShopId)
	err = um.BaseRepo.UpdateModelWithCondition(reqCtx, "shops", whereClause, &shopsinfo)
	if err != nil {
		return err
	}
	return nil
}

func (um *ShopService) CreateShopsByShopsId(reqCtx dto.ReqCtx, shopsinfo dto.CreateShopRequest) error {
	err := um.BaseRepo.Create(reqCtx, &shopsinfo, "shops")
	if err != nil {
		return err
	}
	return nil
}
