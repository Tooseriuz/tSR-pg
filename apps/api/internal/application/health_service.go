package application

import "github.com/tooseriuz/tsr-pg/apps/api/internal/domain"

type HealthService struct{}

func NewHealthService() HealthService {
	return HealthService{}
}

func (s HealthService) Check() string {
	return domain.HealthStatusOK
}
