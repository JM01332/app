package app

import (
	"net/http"

	carrierrouter "github.com/JM01332/app/internal/carrier/router"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type healthResponse struct {
	Status string `json:"status"`
}

// NewRouter creates the HTTP router with all application routes.
func NewRouter(carrierService carrierrouter.CarrierService, logger *zap.Logger) *gin.Engine {
	if logger == nil {
		logger = zap.NewNop()
	}

	router := gin.New()
	router.Use(requestLoggingMiddleware(logger), recoveryMiddleware(logger))

	router.GET("/health", health)
	if carrierService != nil {
		carrierrouter.RegisterRoutes(router.Group("/api"), carrierService)
	}

	return router
}

func health(context *gin.Context) {
	context.JSON(http.StatusOK, healthResponse{Status: "ok"})
}
