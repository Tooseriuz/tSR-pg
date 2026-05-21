package application

import (
	"github.com/tooseriuz/tsr-pg/apps/api/internal/domain"
	"github.com/tooseriuz/tsr-pg/apps/api/internal/generated/openapi"
)

type HealthService struct{}

func NewHealthService() HealthService {
	return HealthService{}
}

func (s HealthService) Check() openapi.HealthResponse {
	return openapi.HealthResponse{
		Status: domain.HealthStatusOK,
	}
}
