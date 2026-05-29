package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tooseriuz/tsr-pg/apps/api/internal/dto/openapi"
	"github.com/tooseriuz/tsr-pg/apps/api/internal/service"
)

func registerHealthRoutes(router gin.IRoutes) {
	service := service.NewHealthService()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, openapi.HealthResponse{
			Status: service.Check(),
		})
	})
}
