package services

import (
	"backend-hostego/internal/app/hostego-service/dto"
	"backend-hostego/internal/app/hostego-service/models"
	"backend-hostego/internal/app/hostego-service/repository"
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type CartService struct {
	*repository.BaseRepo
}

func NewCartService(baseRepo *repository.BaseRepo) *CartService {
	return &CartService{BaseRepo: baseRepo}
}

func (cs *CartService) ProductValidationInCart(reqCtx dto.ReqCtx, cartRequest dto.CartItemRequest) (dto.CartItemRequest, error) {

	var request models.CartItem

	WhereClause := fmt.Sprintf("product_id = '%v' and user_id = %v", cartRequest.ProductId, cartRequest.UserId)
	err := cs.BaseRepo.GetRecordsByCondition(&request, WhereClause, "cart_items")

	if err != nil {
		return dto.CartItemRequest{}, err
	}

	cartRequest.Quantity = request.Quantity + cartRequest.Quantity
	if !(request.ProductItem.IsAvailable) || cartRequest.Quantity > request.ProductItem.MaxQuantity && cartRequest.Quantity > request.ProductItem.AvailableQuantity {
		return dto.CartItemRequest{}, errors.New("product or its more quantity is not available !")
	}

	cartRequest.SubTotal = decimal.NewFromFloat(float64(request.Quantity)).Mul(request.ProductItem.Price)

	return cartRequest, nil
}

func (cs *CartService) UpdateCartValue(reqCtx dto.ReqCtx, cartRequest dto.CartItemRequest) error {
	cartRequest.UpdatedAt = &time.Time{}
	whereClause := fmt.Sprintf("product_id = '%v' and user_id = %v", cartRequest.ProductId, cartRequest.UserId)
	err := cs.BaseRepo.UpdateModelWithCondition(reqCtx, "cart_items", whereClause, &cartRequest)
	if err != nil {
		return err
	}
	return nil
}

func (cs *CartService) GetCartItemsofUser(reqCtx dto.ReqCtx, cartRequest dto.CartItemRequest) ([]models.CartItem, error) {
	var request []models.CartItem
	whereClause := fmt.Sprintf("user_id = %v", cartRequest.UserId)
	err := cs.BaseRepo.GetRecordsByCondition(&request, "cart_items", whereClause)
	if err != nil {
		return []models.CartItem{}, err
	}
	return request, nil
}
