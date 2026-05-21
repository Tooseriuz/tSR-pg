package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tooseriuz/tsr-pg/apps/api/internal/application"
)

func registerHealthRoutes(router gin.IRoutes) {
	service := application.NewHealthService()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, service.Check())
	})
}
