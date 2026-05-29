package service

import (
	"context"
	"errors"
	"strings"

	"github.com/tooseriuz/tsr-pg/apps/api/internal/domain"
)

var ErrJourneyNotFound = errors.New("journey not found")
var ErrInvalidJourney = errors.New("invalid journey")

type JourneyRepository interface {
	List(ctx context.Context) ([]domain.Journey, error)
	Get(ctx context.Context, id int64) (domain.JourneyContent, error)
	Create(ctx context.Context, journey domain.CreateJourney) (int64, error)
}

type JourneyService struct {
	repository JourneyRepository
}

func NewJourneyService(repository JourneyRepository) JourneyService {
	return JourneyService{repository: repository}
}

func (s JourneyService) List(ctx context.Context) ([]domain.Journey, error) {
	return s.repository.List(ctx)
}

func (s JourneyService) Get(ctx context.Context, id int64) (domain.JourneyContent, error) {
	return s.repository.Get(ctx, id)
}

func (s JourneyService) Create(ctx context.Context, journey domain.CreateJourney) (int64, error) {
	journey.Name = strings.TrimSpace(journey.Name)
	journey.Location = strings.TrimSpace(journey.Location)
	journey.Content = strings.TrimSpace(journey.Content)
	if journey.Thumbnail != nil {
		trimmedThumbnail := strings.TrimSpace(*journey.Thumbnail)
		if trimmedThumbnail == "" {
			journey.Thumbnail = nil
		} else {
			journey.Thumbnail = &trimmedThumbnail
		}
	}
	if journey.Name == "" || journey.Location == "" || journey.Content == "" {
		return 0, ErrInvalidJourney
	}

	return s.repository.Create(ctx, journey)
}
