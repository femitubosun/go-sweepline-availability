package availability

import (
	"context"

	"github.com/femitubosun/go-sweepline-availability/internal/redis"
)

type Service interface {
	GetAvailability(ctx context.Context, input GetAvailabilityQuery) (GetAvailabilityQuery, error)
}

type svc struct {
	cache *redis.Cache
}

func NewService(cache *redis.Cache) Service {
	return &svc{
		cache: cache,
	}
}

func (s *svc) GetAvailability(ctx context.Context, input GetAvailabilityQuery) (GetAvailabilityQuery, error) {
	return input, nil
}

type GetAvailabilityInput struct {
	CourtID         string
	DurationMinutes int
	From            string
	To              string
}
