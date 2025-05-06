package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func HealthHandler(ctx *gin.Context) {
	resp := map[string]string{
		"deployed_image": os.Getenv("DEPLOYED_IMAGE_TAG"),
	}
	ctx.JSON(http.StatusOK, resp)
}
