package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tooseriuz/tsr-pg/apps/api/internal/application"
	"github.com/tooseriuz/tsr-pg/apps/api/internal/dto/openapi"
)

func registerHealthRoutes(router gin.IRoutes) {
	service := application.NewHealthService()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, openapi.HealthResponse{
			Status: service.Check(),
		})
	})
}
