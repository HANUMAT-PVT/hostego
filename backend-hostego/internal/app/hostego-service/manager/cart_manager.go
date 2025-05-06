package manager

import (
	"backend-hostego/internal/app/hostego-service/dto"
	"backend-hostego/internal/app/hostego-service/models"
	"backend-hostego/internal/app/hostego-service/services"
	"errors"

	"gorm.io/gorm"
)

type CartManager struct {
	CartService *services.CartService
}

func NewCartManager(CartService *services.CartService) *CartManager {
	return &CartManager{
		CartService: CartService,
	}
}

func (cm *CartManager) AddProductInCart(reqCtx dto.ReqCtx, userCartInfo dto.CartItemRequest) error {
	cartInfo, err := cm.CartService.ProductValidationInCart(reqCtx, userCartInfo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = cm.CartService.Create(reqCtx, &cartInfo, "cart_items")
		}
		return err
	}
	err = cm.CartService.UpdateCartValue(reqCtx, cartInfo)

	if err != nil {
		return err
	}
	return err
}

func (cm *CartManager) FetchCartByUserId(reqCtx dto.ReqCtx, userCartInfo dto.CartItemRequest) ([]models.CartItem, error) {
	request, err := cm.CartService.GetCartItemsofUser(reqCtx, userCartInfo)

	if err != nil {
		return []models.CartItem{}, err
	}
	return request, err
}
