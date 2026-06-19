package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type healthResponse struct {
	Status string `json:"status"`
}

// NewRouter creates the HTTP router with all application routes.
func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/health", health)

	return router
}

func health(context *gin.Context) {
	context.JSON(http.StatusOK, healthResponse{Status: "ok"})
}
