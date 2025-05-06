package manager

import (
	"backend-hostego/internal/app/hostego-service/dto"
	"backend-hostego/internal/app/hostego-service/models"
	"backend-hostego/internal/app/hostego-service/services"
)

type ShopManager struct {
	ShopService *services.ShopService
}

func NewShopManager(ShopService *services.ShopService) *ShopManager {
	return &ShopManager{
		ShopService: ShopService,
	}
}

func (um *ShopManager) GetAllShops(reqCtx dto.ReqCtx) ([]models.Shop, error) {
	data, err := um.ShopService.GetShops(reqCtx)

	if err != nil {
		return []models.Shop{}, err
	}

	return data, err
}

func (sm *ShopManager) GetShopByUserId(reqCtx dto.ReqCtx, userId int) ([]models.Shop, error) {
	data, err := sm.ShopService.GetShopsByUser(reqCtx, userId)

	if err != nil {
		return []models.Shop{}, err
	}

	return data, err
}

func (sm *ShopManager) GetShopByShopId(reqCtx dto.ReqCtx, shopInfo dto.ShopRequest) ([]models.Shop, error) {
	data, err := sm.ShopService.GetShopsByShopsId(reqCtx, shopInfo)

	if err != nil {
		return []models.Shop{}, err
	}

	return data, err
}

func (sm *ShopManager) UpdateShopByShopId(reqCtx dto.ReqCtx, shopInfo dto.UpdateShopRequest) error {
	err := sm.ShopService.UpdateShopsByShopsId(reqCtx, shopInfo)

	if err != nil {
		return err
	}

	return err
}

func (sm *ShopManager) CreateShop(reqCtx dto.ReqCtx, shopInfo dto.CreateShopRequest) error {
	err := sm.ShopService.CreateShopsByShopsId(reqCtx, shopInfo)

	if err != nil {
		return err
	}

	return err
}
