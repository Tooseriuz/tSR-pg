package http

import (
	"github.com/gin-gonic/gin"
	"github.com/tooseriuz/tsr-pg/apps/api/internal/service"
)

type RouterConfig struct {
	journeyRepository service.JourneyRepository
}

type RouterOption func(*RouterConfig)

func WithJourneyRepository(repository service.JourneyRepository) RouterOption {
	return func(config *RouterConfig) {
		config.journeyRepository = repository
	}
}

func NewRouter(options ...RouterOption) *gin.Engine {
	config := RouterConfig{
		journeyRepository: newMemoryJourneyRepository(),
	}
	for _, option := range options {
		option(&config)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	registerHealthRoutes(router)
	registerJourneyRoutes(router, config.journeyRepository)

	return router
}
