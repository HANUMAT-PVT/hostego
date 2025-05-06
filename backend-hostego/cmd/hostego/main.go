package main

import (
	"backend-hostego/internal/app/hostego-service/database"
	"backend-hostego/internal/app/hostego-service/routes"

	"backend-hostego/api/rest"
	"backend-hostego/config"
	"backend-hostego/internal/app/hostego-service/constants/config_constants"
	"backend-hostego/internal/app/hostego-service/constants/string_constants"
	"backend-hostego/internal/pkg/db/postgres"
	"backend-hostego/internal/pkg/db/redis"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	// "backend-hostego/internal/pkg/db/postgres"
	// "backend-hostego/internal/pkg/db/redis"
	"backend-hostego/internal/pkg/logger"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Global DB variable
var db *gorm.DB

var (
	log *logger.MyLogger
)

func InitLogger() {
	loggingFormat := viper.GetString(config_constants.VKEYS_LOGGING_FORMAT)
	loggingLevel := viper.GetString(config_constants.VKEYS_LOGGING_LEVEL)
	log = logger.GetLoggerWithFormatAndLogging(loggingFormat, loggingLevel)
}

func init() {
	godotenv.Load()

	fmt.Println("loading config...") // log is initialized after config setup
	config.SetupConfig()

	fmt.Println("init logger...")
	InitLogger()

	log.Info("initializing postgres db connection...")
	postgres.InitDB()

	log.Info("initializing redis connection...")
	redis.InitRedis()
}

// Returns address - Eg: localhost - 0.0.0.0:8080
func getAddr() string {
	addr := config_constants.CONNECTION_ADDRESS
	ip := viper.GetString(config_constants.VKEYS_HOST_IP)
	port := viper.GetString(config_constants.VKEYS_HOST_PORT)
	if ip != string_constants.EMPTY_STRING || port != string_constants.EMPTY_STRING {
		addr = ip + string_constants.COLLON + port
	}
	return addr
}

func main() {
	log.Info("starting server...")
	// TODO - handle if we are not able to build server
	server := rest.HttpBuildServer(getAddr())
	err := server.ListenAndServe()
	if err != nil {
		log.Errorf("error while starting server %v", err.Error())
		return
	}
	log.Infof("http server is started")

	app := fiber.New()

	database.ConnectDataBase()

	// routes.SetupRoutes(app)
	routes.AuthRoutes(app)
	routes.ShopRoutes(app)
	routes.ProductRoutes(app)
	routes.OrderRoutes(app)
	routes.CartRoutes(app)
	routes.WalletRoutes(app)
	routes.PaymentRoutes(app)
	routes.DeliveryPartnerRoutes(app)
	// Fetch all users
	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome to the server Backend"})
	})
}
