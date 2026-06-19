package app

import (
	"net/http"

	carrierrouter "github.com/JM01332/app/internal/carrier/router"
	"github.com/gin-gonic/gin"
)

type healthResponse struct {
	Status string `json:"status"`
}

// NewRouter creates the HTTP router with all application routes.
func NewRouter(carrierService carrierrouter.CarrierService) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/health", health)
	if carrierService != nil {
		carrierrouter.RegisterRoutes(router.Group("/api"), carrierService)
	}

	return router
}

func health(context *gin.Context) {
	context.JSON(http.StatusOK, healthResponse{Status: "ok"})
}
