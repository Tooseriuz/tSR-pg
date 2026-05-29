package service

import (
	"context"
	"testing"
	"time"

	"github.com/tooseriuz/tsr-pg/apps/api/internal/domain"
)

type stubJourneyRepository struct {
	journeys []domain.Journey
}

func (r stubJourneyRepository) List(context.Context) ([]domain.Journey, error) {
	return r.journeys, nil
}

func TestJourneyServiceList(t *testing.T) {
	timestamp := time.Date(2026, time.May, 29, 0, 0, 0, 0, time.UTC)
	service := NewJourneyService(stubJourneyRepository{
		journeys: []domain.Journey{
			{
				Name:      "Build the journey endpoint",
				Timestamp: timestamp,
				Location:  "Bangkok",
				Thumbnail: "https://example.com/journey.jpg",
			},
		},
	})

	journeys, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("list journeys: %v", err)
	}

	if len(journeys) != 1 {
		t.Fatalf("expected 1 journey, got %d", len(journeys))
	}

	if journeys[0].Name != "Build the journey endpoint" {
		t.Fatalf("expected journey name, got %q", journeys[0].Name)
	}
}
