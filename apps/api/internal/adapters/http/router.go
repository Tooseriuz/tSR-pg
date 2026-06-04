package http

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tooseriuz/tsr-pg/apps/api/internal/service"
)

type RouterConfig struct {
	journeyRepository service.JourneyRepository
	imageStorage      service.ImageStorage
	adminToken        string
}

type RouterOption func(*RouterConfig)

func WithJourneyRepository(repository service.JourneyRepository) RouterOption {
	return func(config *RouterConfig) {
		config.journeyRepository = repository
	}
}

func WithAdminToken(token string) RouterOption {
	return func(config *RouterConfig) {
		config.adminToken = token
	}
}

func WithImageStorage(storage service.ImageStorage) RouterOption {
	return func(config *RouterConfig) {
		config.imageStorage = storage
	}
}

func NewRouter(options ...RouterOption) *gin.Engine {
	config := RouterConfig{
		journeyRepository: newMemoryJourneyRepository(),
		adminToken:        os.Getenv("ADMIN_TOKEN"),
	}
	for _, option := range options {
		option(&config)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), corsMiddleware())

	registerHealthRoutes(router)
	registerAdminRoutes(router, config.adminToken)
	registerJourneyRoutes(router, config.journeyRepository, config.imageStorage, config.adminToken)

	return router
}

func corsMiddleware() gin.HandlerFunc {
	allowedOrigins := map[string]bool{
		"https://www.tooseriuz.com": true,
		"http://localhost:3000":     true,
		"http://localhost:3001":     true,
	}

	return func(context *gin.Context) {
		origin := context.GetHeader("Origin")
		if allowedOrigins[origin] {
			context.Header("Access-Control-Allow-Origin", origin)
			context.Header("Access-Control-Allow-Credentials", "true")
			context.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			context.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			context.Header("Vary", "Origin")
		}

		if context.Request.Method == http.MethodOptions {
			context.AbortWithStatus(http.StatusNoContent)
			return
		}

		context.Next()
	}
}
