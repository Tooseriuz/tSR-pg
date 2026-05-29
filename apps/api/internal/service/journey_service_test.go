package service

import (
	"context"
	"testing"
	"time"

	"github.com/tooseriuz/tsr-pg/apps/api/internal/domain"
)

type stubJourneyRepository struct {
	journeys       []domain.Journey
	journeyContent domain.JourneyContent
	createdID      int64
}

func (r stubJourneyRepository) List(context.Context) ([]domain.Journey, error) {
	return r.journeys, nil
}

func (r stubJourneyRepository) Get(context.Context, int64) (domain.JourneyContent, error) {
	return r.journeyContent, nil
}

func (r stubJourneyRepository) Create(context.Context, domain.CreateJourney) (int64, error) {
	return r.createdID, nil
}

func TestJourneyServiceList(t *testing.T) {
	timestamp := time.Date(2026, time.May, 29, 0, 0, 0, 0, time.UTC)
	thumbnail := "https://example.com/journey.jpg"
	service := NewJourneyService(stubJourneyRepository{
		journeys: []domain.Journey{
			{
				Name:      "Build the journey endpoint",
				Timestamp: timestamp,
				Location:  "Bangkok",
				Thumbnail: &thumbnail,
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

func TestJourneyServiceGet(t *testing.T) {
	timestamp := time.Date(2026, time.May, 29, 0, 0, 0, 0, time.UTC)
	service := NewJourneyService(stubJourneyRepository{
		journeyContent: domain.JourneyContent{
			Name:      "Build the journey endpoint",
			Timestamp: timestamp,
			Content:   "# Start",
		},
	})

	journey, err := service.Get(context.Background(), 1)
	if err != nil {
		t.Fatalf("get journey: %v", err)
	}

	if journey.Content != "# Start" {
		t.Fatalf("expected journey content, got %q", journey.Content)
	}
}

func TestJourneyServiceCreate(t *testing.T) {
	service := NewJourneyService(stubJourneyRepository{createdID: 42})

	id, err := service.Create(context.Background(), domain.CreateJourney{
		Name:      " Build the journey endpoint ",
		Location:  " Bangkok ",
		Content:   " # Start ",
		Highlight: true,
	})
	if err != nil {
		t.Fatalf("create journey: %v", err)
	}

	if id != 42 {
		t.Fatalf("expected created journey id 42, got %d", id)
	}
}

func TestJourneyServiceCreateRejectsEmptyContent(t *testing.T) {
	service := NewJourneyService(stubJourneyRepository{})

	_, err := service.Create(context.Background(), domain.CreateJourney{Name: "Build the journey endpoint", Location: "Bangkok"})
	if err != ErrInvalidJourney {
		t.Fatalf("expected ErrInvalidJourney, got %v", err)
	}
}
