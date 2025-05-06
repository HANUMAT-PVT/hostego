package manager

import (
	"backend-hostego/internal/app/hostego-service/dto"
	"backend-hostego/internal/app/hostego-service/models"
	"backend-hostego/internal/app/hostego-service/services"
)

type DeliveryPartnerManager struct {
	DeliveryPartnerService *services.DeliveryPartnerService
}

func NewDeliveryPartnerManager(DeliveryPartnerService *services.DeliveryPartnerService) *DeliveryPartnerManager {
	return &DeliveryPartnerManager{
		DeliveryPartnerService: DeliveryPartnerService,
	}
}

func (dpm *DeliveryPartnerManager) FetchDeliveryPartnerData(reqCtx dto.ReqCtx) (models.DeliveryPartner, error) {
	data, err := dpm.DeliveryPartnerService.GetDeliveryPartnerInfo(reqCtx)

	if err != nil {
		return models.DeliveryPartner{}, err
	}

	return data, err
}

func (dpm *DeliveryPartnerManager) CreateDeliveryPartnerData(reqCtx dto.ReqCtx, request dto.DeliveryPartnerRequestDto) error {
	err := dpm.DeliveryPartnerService.CreateDeliveryPartner(reqCtx, request)

	if err != nil {
		return err
	}

	return err
}

func (dpm *DeliveryPartnerManager) UpdateDeliveryPartnerInfo(reqCtx dto.ReqCtx, request dto.DeliveryPartnerRequestDto) error {
	err := dpm.DeliveryPartnerService.UpdateDeliveryPartnerByUserId(reqCtx, request)
	if err != nil {
		return err
	}
	return err
}
