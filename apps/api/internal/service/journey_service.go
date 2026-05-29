package service

import (
	"context"

	"github.com/tooseriuz/tsr-pg/apps/api/internal/domain"
)

type JourneyRepository interface {
	List(ctx context.Context) ([]domain.Journey, error)
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
