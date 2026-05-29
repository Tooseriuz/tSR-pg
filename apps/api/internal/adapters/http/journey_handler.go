package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tooseriuz/tsr-pg/apps/api/internal/domain"
	"github.com/tooseriuz/tsr-pg/apps/api/internal/dto/openapi"
	"github.com/tooseriuz/tsr-pg/apps/api/internal/service"
)

type memoryJourneyRepository struct {
	journeys []domain.Journey
}

func newMemoryJourneyRepository() memoryJourneyRepository {
	return memoryJourneyRepository{
		journeys: []domain.Journey{
			{
				Name:      "Remote desk, first production habits",
				Timestamp: time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
				Location:  "Bangkok",
				Thumbnail: "https://picsum.photos/seed/tsr-remote-desk-2020/800/600",
			},
		},
	}
}

func (r memoryJourneyRepository) List(context.Context) ([]domain.Journey, error) {
	return r.journeys, nil
}

func registerJourneyRoutes(router gin.IRoutes, repository service.JourneyRepository) {
	service := service.NewJourneyService(repository)

	router.GET("/journeys", func(c *gin.Context) {
		journeys, err := service.List(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list journeys"})
			return
		}

		c.JSON(http.StatusOK, toJourneysResponse(journeys))
	})
}

func toJourneysResponse(journeys []domain.Journey) openapi.JourneysResponse {
	response := make(openapi.JourneysResponse, 0, len(journeys))

	for _, journey := range journeys {
		response = append(response, openapi.Journey{
			Name:      journey.Name,
			Timestamp: journey.Timestamp.Format(time.DateOnly),
			Location:  journey.Location,
			Thumbnail: journey.Thumbnail,
		})
	}

	return response
}
