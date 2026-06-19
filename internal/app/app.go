package app

import (
	"net/http"

	carrierrouter "github.com/JM01332/app/internal/carrier/router"
	"github.com/JM01332/app/internal/security"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type healthResponse struct {
	Status string `json:"status"`
}

// NewRouter creates the HTTP router with all application routes.
func NewRouter(
	carrierService carrierrouter.CarrierService,
	logger *zap.Logger,
	tokenVerifier security.Verifier,
	oidcClientID string,
) *gin.Engine {
	if logger == nil {
		logger = zap.NewNop()
	}

	router := gin.New()
	router.Use(requestLoggingMiddleware(logger), recoveryMiddleware(logger))

	router.GET("/health", health)
	if carrierService != nil {
		if tokenVerifier == nil {
			panic("token verifier is required for carrier routes")
		}

		api := router.Group("/api")
		api.Use(security.Authenticate(tokenVerifier))
		carrierrouter.RegisterRoutes(api, carrierService, carrierrouter.RouteAuthorization{
			Read:  security.RequireAnyRole(oidcClientID, security.RoleUser, security.RoleAdmin),
			Write: security.RequireAnyRole(oidcClientID, security.RoleAdmin),
		})
	}

	return router
}

func health(context *gin.Context) {
	context.JSON(http.StatusOK, healthResponse{Status: "ok"})
}
