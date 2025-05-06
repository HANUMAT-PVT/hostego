package hostego

import (
	"backend-hostego/internal/app/hostego-service/constants/api_constants"
	"backend-hostego/internal/app/hostego-service/controller/handler"
	"backend-hostego/internal/app/hostego-service/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(routerGroup *gin.RouterGroup) *gin.RouterGroup {
	container := NewContainer()
	v1 := routerGroup.Group(api_constants.V1)
	auth := v1.Group(api_constants.AUTH)
	UserRoutes(v1, &container)
	UserAuthRoutes(auth, &container)
	return v1
}

func UserRoutes(routerGroup *gin.RouterGroup, container *Container) {
	routerGroup.GET(api_constants.TEST_HEALTH, handler.HealthHandler)
	routerGroup.GET("/users/private", middlewares.InternalTokenVerifivationORAuth(), container.userController.GetUsers)
	routerGroup.GET("/user/profile", middlewares.UserAuthTokenValidation(), container.userController.GetUserByUserId)
	routerGroup.PATCH("/profile-update", middlewares.InternalTokenVerifivationORAuth(), container.userController.UpdateUserById)
	routerGroup.POST("/profile/address-update", middlewares.InternalTokenVerifivationORAuth(), container.userController.UpdateUserAddress)
	routerGroup.POST("/profile/address-create", middlewares.InternalTokenVerifivationORAuth(), container.userController.CreateUserAddress)
	routerGroup.POST("/profile/address-delete", middlewares.InternalTokenVerifivationORAuth(), container.userController.DeleteUserAddress)
}

func UserAuthRoutes(routerGroup *gin.RouterGroup, container *Container) {
	routerGroup.POST("/user-signup", container.authController.CreateUserBySignUp)
}

func UserCartRoutes(routerGroup *gin.RouterGroup, container *Container) {
	routerGroup.POST("/cart/add-product", middlewares.UserAuthTokenValidation(), container.cartController.AddProductInUserCart)
	routerGroup.GET("/", middlewares.UserAuthTokenValidation(), container.cartController.FetchUserCart)
}

// Compare this snippet from internal/app/hostego-service/constants/api_constants.go:
// cartRoutes.Post("/", controllers.AddProductInUserCart)
// cartRoutes.Get("/", controllers.FetchUserCart)
