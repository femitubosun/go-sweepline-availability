package booking

import (
	"context"
)

type Service interface {
	CreateBooking(ctx context.Context, input CreateBookingRequest) error
}

type svc struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) CreateBooking(ctx context.Context, input CreateBookingRequest) error {
	return nil
}
