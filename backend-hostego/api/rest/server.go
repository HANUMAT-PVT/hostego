package rest

import (
	"backend-hostego/api/rest/hostego"
	"backend-hostego/internal/app/hostego-service/constants/api_constants"
	"backend-hostego/internal/app/hostego-service/constants/config_constants"
	"backend-hostego/internal/app/hostego-service/middlewares"

	// "backend-hostego/internal/pkg/newrelic_setup"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// BuildServer adds all the configurations and returns the pointer to the server
func BuildServer() *gin.Engine {

	server := gin.New()
	//server.Use(middleware.QueryLogMiddleware())
	server.Use(customRecovery())
	server.Use(cors.New(cors.Config{
		AllowOrigins:     viper.GetStringSlice(config_constants.VKEYS_CORS_ORIGINS),
		AllowMethods:     viper.GetStringSlice(config_constants.VKEYS_ALLOW_METHODS),
		AllowHeaders:     viper.GetStringSlice(config_constants.VKEYS_SERVER_ALLOW_HEADERS),
		ExposeHeaders:    viper.GetStringSlice(config_constants.VKEYS_EXPOSED_HEADERS),
		AllowCredentials: true,
	}))

	// add new relic to all routes, it attaches the transaction context to gin context and is picked up
	// server.Use(nrgin.Middleware(newrelic_setup.GetNewRelicApp(viper.GetString(config_constants.VKEYS_NEWRELIC_APP_NAME))))
	server.Use(middlewares.RequestRequirements())

	pprof.Register(server, "/hostego/api/v1/pprof")
	hostego.RegisterRoutes(server.Group(api_constants.API_PATH))

	// c := custom_metrics.GetInstance()
	// server.GET(api_constants.METRICS, ginHandler(c.HttpHandler))
	// initialize global pseudo random generator
	rand.Seed(time.Now().Unix())

	return server
}

func ginHandler(httpHandler http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		httpHandler.ServeHTTP(c.Writer, c.Request)
	}
}

// HttpBuildServer for solving timeout issue - http server using gin router
func HttpBuildServer(addr string) *http.Server {
	server := BuildServer()
	s := &http.Server{
		Addr:         addr,
		Handler:      server,
		ReadTimeout:  time.Duration(viper.GetInt(config_constants.VKEYS_READ_WRITE_TIMEOUT_SERVER)) * time.Second,
		WriteTimeout: time.Duration(viper.GetInt(config_constants.VKEYS_READ_WRITE_TIMEOUT_SERVER)) * time.Second,
	}
	s.IdleTimeout = time.Duration(viper.GetInt(config_constants.VKEYS_IDLE_TIMEOUT_SERVER)) * time.Minute

	return s
}
