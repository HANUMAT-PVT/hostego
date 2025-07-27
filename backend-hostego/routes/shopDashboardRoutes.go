package routes

import (
	"backend-hostego/controllers"

	"github.com/gofiber/fiber/v2"
)

func ShopDashboardRoutes(app *fiber.App) {
	shopDashboard := app.Group("/api/shop-dashboard/")
	shopDashboard.Get("/delivery-stats/:shop_id", controllers.GetShopDashboardStats)
	shopDashboard.Get("/top-selling-products/:shop_id", controllers.GetTopSellingProducts)
	shopDashboard.Get("/order-analytics/:shop_id", controllers.GetOrderAnalytics)
	shopDashboard.Get("/customer-insights/:shop_id", controllers.GetCustomerInsights)
	shopDashboard.Get("/performance-metrics/:shop_id", controllers.GetRestaurantPerformanceMetrics)
	shopDashboard.Get("/revenue-analytics/:shop_id", controllers.GetRestaurantRevenueAnalytics)
}
