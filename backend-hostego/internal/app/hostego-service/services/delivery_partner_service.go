package services

import (
	"backend-hostego/internal/app/hostego-service/dto"
	"backend-hostego/internal/app/hostego-service/models"
	"backend-hostego/internal/app/hostego-service/repository"
	"fmt"
	"time"
)

type DeliveryPartnerService struct {
	*repository.BaseRepo
}

func NewDeliveryPartnerService(baseRepo *repository.BaseRepo) *DeliveryPartnerService {
	return &DeliveryPartnerService{BaseRepo: baseRepo}
}

func (dps *DeliveryPartnerService) GetDeliveryPartnerInfo(reqCtx dto.ReqCtx) (models.DeliveryPartner, error) {
	var request models.DeliveryPartner

	WhereClause := fmt.Sprintf("user_id = %d", reqCtx.UserId)
	err := dps.BaseRepo.GetRecordsByCondition(&request, WhereClause, "delivery_partners")
	if err != nil {
		return models.DeliveryPartner{}, err
	}

	return request, err
}

func (dps *DeliveryPartnerService) CreateDeliveryPartner(reqCtx dto.ReqCtx, requestinfo dto.DeliveryPartnerRequestDto) error {

	err := dps.BaseRepo.Create(reqCtx, &requestinfo, "delivery_partners")
	if err != nil {
		return err
	}

	return err
}

func (dps *DeliveryPartnerService) UpdateDeliveryPartnerByUserId(reqCtx dto.ReqCtx, requestInfo dto.DeliveryPartnerRequestDto) error {
	requestInfo.UpdatedAt = &time.Time{}
	whereClause := fmt.Sprintf("user_id = %v", requestInfo.UserId)
	err := dps.BaseRepo.UpdateModelWithCondition(reqCtx, "delivery_partners", whereClause, &requestInfo)
	if err != nil {
		return err
	}
	return nil
}
