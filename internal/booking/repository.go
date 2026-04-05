package booking

import (
	"context"
	"time"

	"github.com/femitubosun/go-sweepline-availability/internal/redis"
)

type Interval struct {
	Start time.Time
	End   time.Time
}

type Repository interface {
	Reserve(ctx context.Context, input ReserveInput) (Booking, error)
	Release(ctx context.Context, bookingID string) error
	GetActive(ctx context.Context, courtID string, interval Interval) ([]Booking, error)
	GetGaps(ctx context.Context, input GetGapsInput) error
}

type repo struct {
	cache *redis.Cache
}

func (r *repo) Reserve(ctx context.Context, input ReserveInput) (Booking, error) {
	return Booking{}, nil
}

func (r *repo) Release(ctx context.Context, bookingID string) error {
	return nil
}

func (r *repo) GetActive(ctx context.Context, courtID string, interval Interval) ([]Booking, error) {
	return make([]Booking, 0), nil
}

func (r *repo) GetGaps(ctx context.Context, input GetGapsInput) error {
	return nil
}

func NewRepo(cache *redis.Cache) Repository {
	return &repo{
		cache: cache,
	}
}

type ReserveInput struct {
	ID        string
	CourtID   string
	Interval  Interval
	GuestName string
}

type GetGapsInput struct {
	CourtID      string
	DurationMins int
	Interval     Interval
}
