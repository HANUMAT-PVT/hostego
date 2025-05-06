package hostego

import (
	"backend-hostego/internal/app/hostego-service/controller"
	"backend-hostego/internal/app/hostego-service/manager"
	"backend-hostego/internal/app/hostego-service/repository"
	"backend-hostego/internal/app/hostego-service/services"
	"backend-hostego/internal/pkg/db/postgres"
)

type Container struct {
	userController           *controller.UserController
	authController           *controller.AuthController
	cartController           *controller.CartController
	delveryPartnerController *controller.DeliverPartnerController
}

func NewContainer() Container {
	baseRepo := repository.NewBaseRepo()
	baseRepo.SetDb(postgres.GetDBWithoutContext())

	userService := services.NewUserService(baseRepo)
	authService := services.NewAuthService(baseRepo)
	cartService := services.NewCartService(baseRepo)
	deliveryPartnerServie := services.NewDeliveryPartnerService(baseRepo)

	userManager := manager.NewUserManager(userService)
	authManager := manager.NewAuthManager(authService)
	cartManager := manager.NewCartManager(cartService)
	deliveryPartnerManager := manager.NewDeliveryPartnerManager(deliveryPartnerServie)

	userController := controller.NewUserController(userManager)
	authController := controller.NewAuthController(authManager)
	cartController := controller.NewCartController(cartManager)
	deliveryPartnerController := controller.NewDeliverPartnerController(deliveryPartnerManager)

	return Container{
		userController:           userController,
		authController:           authController,
		cartController:           cartController,
		delveryPartnerController: deliveryPartnerController,
	}
}
