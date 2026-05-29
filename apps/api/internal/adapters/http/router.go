package http

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tooseriuz/tsr-pg/apps/api/internal/service"
)

type RouterConfig struct {
	journeyRepository service.JourneyRepository
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

func NewRouter(options ...RouterOption) *gin.Engine {
	config := RouterConfig{
		journeyRepository: newMemoryJourneyRepository(),
		adminToken:        os.Getenv("ADMIN_TOKEN"),
	}
	for _, option := range options {
		option(&config)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	registerHealthRoutes(router)
	registerAdminRoutes(router, config.adminToken)
	registerJourneyRoutes(router, config.journeyRepository, config.adminToken)

	return router
}
